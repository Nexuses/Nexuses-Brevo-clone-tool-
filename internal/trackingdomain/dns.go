package trackingdomain

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
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

// NetResolver performs real DNS CNAME lookups via net.Resolver.
type NetResolver struct {
	Resolver *net.Resolver
	Timeout  time.Duration
}

// LookupCNAME resolves the CNAME for host with a context timeout.
func (r *NetResolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	res := r.Resolver
	if res == nil {
		res = net.DefaultResolver
	}
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = DefaultDNSTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return res.LookupCNAME(ctx, host)
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
