package trackingdomain

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode"
)

var (
	ErrEmptyDomain       = errors.New("domain is required")
	ErrInvalidDomain     = errors.New("invalid domain")
	ErrDomainHasPath     = errors.New("domain must not include a path, query, or fragment")
	ErrDomainHasPort     = errors.New("domain must not include a port")
	ErrDomainIsIP        = errors.New("IP addresses are not allowed")
	ErrDomainLocalhost   = errors.New("localhost is not allowed")
	ErrDomainWildcard    = errors.New("wildcard domains are not allowed")
	ErrUnsupportedScheme = errors.New("only http and https schemes are accepted")
)

// NormalizeDomain trims, lowercases, and validates a tracking domain hostname.
// Optional http:// or https:// prefixes are accepted and stripped. Paths, query
// strings, fragments, ports, IPs, localhost, and wildcards are rejected.
func NormalizeDomain(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", ErrEmptyDomain
	}

	s = strings.ToLower(s)

	// Reject wildcards early.
	if strings.Contains(s, "*") {
		return "", ErrDomainWildcard
	}

	// If a scheme is present, only http/https are allowed.
	if strings.Contains(s, "://") {
		u, err := url.Parse(s)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrInvalidDomain, err)
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return "", ErrUnsupportedScheme
		}
		if u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
			return "", ErrDomainHasPath
		}
		if u.Port() != "" {
			return "", ErrDomainHasPort
		}
		s = u.Hostname()
	} else {
		// Bare host: reject path/query/fragment/userinfo-looking input.
		if strings.ContainsAny(s, "/?#") {
			return "", ErrDomainHasPath
		}
		if strings.Contains(s, "@") {
			return "", ErrInvalidDomain
		}
		// Host:port form.
		if host, port, err := net.SplitHostPort(s); err == nil {
			if port != "" {
				return "", ErrDomainHasPort
			}
			s = host
		} else if strings.Count(s, ":") == 1 {
			// SplitHostPort fails without brackets for IPv6; for host:port with one colon treat as port.
			parts := strings.SplitN(s, ":", 2)
			if parts[1] != "" && isAllDigits(parts[1]) {
				return "", ErrDomainHasPort
			}
		}
	}

	s = strings.TrimSuffix(strings.TrimSpace(s), ".")
	if s == "" {
		return "", ErrEmptyDomain
	}

	if s == "localhost" || strings.HasSuffix(s, ".localhost") {
		return "", ErrDomainLocalhost
	}

	if ip := net.ParseIP(s); ip != nil {
		return "", ErrDomainIsIP
	}
	// Bracketed IPv6 accidentally left behind.
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return "", ErrDomainIsIP
	}

	if err := validateHostname(s); err != nil {
		return "", err
	}

	return s, nil
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func validateHostname(host string) error {
	if len(host) > 253 {
		return ErrInvalidDomain
	}

	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return ErrInvalidDomain
	}

	for _, label := range labels {
		if len(label) < 1 || len(label) > 63 {
			return ErrInvalidDomain
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return ErrInvalidDomain
		}
		for _, r := range label {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return ErrInvalidDomain
			}
		}
	}

	return nil
}

// CanonicalHost lowercases and strips a trailing dot from a DNS name.
func CanonicalHost(host string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
}

// HostFromURL extracts the hostname from a URL or returns a bare hostname as-is.
func HostFromURL(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "://") {
		u, err := url.Parse(s)
		if err != nil {
			return ""
		}
		return CanonicalHost(u.Hostname())
	}
	return CanonicalHost(s)
}

// ExpectedCNAMETarget returns the hostname CTD CNAMEs must point to.
// Empty trackingURL falls back to rootURL.
func ExpectedCNAMETarget(trackingURL, rootURL string) string {
	host := HostFromURL(trackingURL)
	if host == "" {
		host = HostFromURL(rootURL)
	}
	return host
}
