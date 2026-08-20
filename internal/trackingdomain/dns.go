package trackingdomain

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

const (
	StatusPending  = "pending"
	StatusVerified = "verified"
	StatusFailed   = "failed"

	DefaultDNSTimeout = 5 * time.Second
)

// CNAMEResolver looks up CNAME records. Implementations must be safe for tests.
type CNAMEResolver interface {
	LookupCNAME(ctx context.Context, host string) (string, error)
}

// NetResolver performs real DNS CNAME lookups.
// It prefers a direct CNAME query (via public resolvers) because Go's
// net.LookupCNAME often returns the original hostname after following the
// chain — especially inside Docker — which falsely looks like "no CNAME".
type NetResolver struct {
	Resolver *net.Resolver
	Timeout  time.Duration
	// Servers are DNS endpoints used for direct CNAME queries.
	// Empty defaults to Google/Cloudflare public resolvers.
	Servers []string
}

// LookupCNAME resolves the CNAME target for host.
func (r *NetResolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	host = CanonicalHost(host)
	if host == "" {
		return "", &net.DNSError{Err: "empty host", Name: host, IsNotFound: true}
	}

	timeout := r.Timeout
	if timeout <= 0 {
		timeout = DefaultDNSTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if cname, err := r.lookupCNAMEDirect(ctx, host); err == nil && cname != "" {
		return cname, nil
	} else if err != nil && !isNoDataError(err) {
		// Keep going to system resolver unless this was a hard failure.
		if !isTemporaryDNSError(err) && !isNotFoundDNSError(err) {
			// Still try system resolver as fallback.
		}
	}

	res := r.Resolver
	if res == nil {
		res = net.DefaultResolver
	}
	cname, err := res.LookupCNAME(ctx, host)
	if err != nil {
		return "", err
	}
	got := CanonicalHost(cname)
	if got == "" || got == host {
		return "", &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	return got, nil
}

func (r *NetResolver) lookupCNAMEDirect(ctx context.Context, host string) (string, error) {
	servers := r.Servers
	if len(servers) == 0 {
		servers = []string{"8.8.8.8:53", "1.1.1.1:53"}
	}

	var lastErr error
	for _, server := range servers {
		cname, err := queryCNAME(ctx, server, host)
		if err == nil && cname != "" {
			return cname, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	if lastErr == nil {
		lastErr = &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
	}
	return "", lastErr
}

func queryCNAME(ctx context.Context, server, host string) (string, error) {
	var msg dnsmessage.Message
	msg.Header.RecursionDesired = true
	msg.Questions = []dnsmessage.Question{{
		Name:  mustDNSName(host),
		Type:  dnsmessage.TypeCNAME,
		Class: dnsmessage.ClassINET,
	}}

	packed, err := msg.Pack()
	if err != nil {
		return "", err
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "udp", server)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	if _, err := conn.Write(packed); err != nil {
		return "", err
	}

	buf := make([]byte, 1500)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	var resp dnsmessage.Message
	if err := resp.Unpack(buf[:n]); err != nil {
		return "", err
	}
	if resp.RCode != dnsmessage.RCodeSuccess {
		return "", &net.DNSError{
			Err:        fmt.Sprintf("dns rcode %v", resp.RCode),
			Name:       host,
			IsNotFound: resp.RCode == dnsmessage.RCodeNameError,
		}
	}

	want := CanonicalHost(host)
	for _, a := range resp.Answers {
		if a.Header.Type != dnsmessage.TypeCNAME {
			continue
		}
		cname, ok := a.Body.(*dnsmessage.CNAMEResource)
		if !ok {
			continue
		}
		owner := CanonicalHost(a.Header.Name.String())
		target := CanonicalHost(cname.CNAME.String())
		if owner == want && target != "" && target != want {
			return target, nil
		}
	}

	return "", &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
}

func mustDNSName(host string) dnsmessage.Name {
	n, err := dnsmessage.NewName(host + ".")
	if err != nil {
		// Fall back to empty; Pack will fail if invalid.
		var empty dnsmessage.Name
		return empty
	}
	return n
}

func isNoDataError(err error) bool {
	return isNotFoundDNSError(err)
}

func isNotFoundDNSError(err error) bool {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return dnsErr.IsNotFound
	}
	return false
}

// VerifyResult is the outcome of a DNS CNAME verification attempt.
type VerifyResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	CNAME   string `json:"cname,omitempty"`
}

// VerifyCNAME checks that domain CNAMEs to expectedHost (from app.tracking_url).
// Not found / temporary / propagation errors leave status pending; wrong target is failed;
// exact canonical match is verified.
func VerifyCNAME(ctx context.Context, resolver CNAMEResolver, domain, expectedHost string) VerifyResult {
	domain = CanonicalHost(domain)
	expectedHost = CanonicalHost(expectedHost)

	if domain == "" || expectedHost == "" {
		return VerifyResult{
			Status:  StatusFailed,
			Message: "missing domain or expected CNAME target",
		}
	}

	cname, err := resolver.LookupCNAME(ctx, domain)
	if err != nil {
		return classifyLookupError(err)
	}

	got := CanonicalHost(cname)
	// Some resolvers return the queried name when there is no CNAME.
	if got == "" || got == domain {
		return VerifyResult{
			Status:  StatusPending,
			Message: "CNAME record not found yet; DNS may still be propagating",
			CNAME:   got,
		}
	}

	if got != expectedHost {
		return VerifyResult{
			Status:  StatusFailed,
			Message: fmt.Sprintf("CNAME points to %q, expected %q", got, expectedHost),
			CNAME:   got,
		}
	}

	return VerifyResult{
		Status:  StatusVerified,
		Message: fmt.Sprintf("CNAME verified: %s → %s", domain, expectedHost),
		CNAME:   got,
	}
}

func classifyLookupError(err error) VerifyResult {
	msg := err.Error()
	lower := strings.ToLower(msg)

	// Temporary / not found / propagation → pending.
	if isTemporaryDNSError(err) ||
		strings.Contains(lower, "no such host") ||
		strings.Contains(lower, "server misbehaving") ||
		strings.Contains(lower, "i/o timeout") ||
		strings.Contains(lower, "timeout") ||
		strings.Contains(lower, "temporary failure") ||
		strings.Contains(lower, "try again") {
		return VerifyResult{
			Status:  StatusPending,
			Message: "DNS lookup pending or temporarily unavailable: " + msg,
		}
	}

	return VerifyResult{
		Status:  StatusPending,
		Message: "DNS lookup did not succeed yet: " + msg,
	}
}

func isTemporaryDNSError(err error) bool {
	type temporary interface {
		Temporary() bool
	}
	if t, ok := err.(temporary); ok && t.Temporary() {
		return true
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return dnsErr.Temporary() || dnsErr.IsNotFound
	}
	return false
}
