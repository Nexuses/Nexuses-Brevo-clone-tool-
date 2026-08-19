package trackingdomain

import (
	"context"
	"errors"
	"net"
	"testing"
)

func TestNormalizeDomain(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr error
	}{
		{"  Track.Example.COM. ", "track.example.com", nil},
		{"https://track.example.com", "track.example.com", nil},
		{"http://track.example.com/", "track.example.com", nil},
		{"https://track.example.com/path", "", ErrDomainHasPath},
		{"track.example.com?x=1", "", ErrDomainHasPath},
		{"track.example.com#frag", "", ErrDomainHasPath},
		{"track.example.com:443", "", ErrDomainHasPort},
		{"127.0.0.1", "", ErrDomainIsIP},
		{"localhost", "", ErrDomainLocalhost},
		{"*.example.com", "", ErrDomainWildcard},
		{"ftp://track.example.com", "", ErrUnsupportedScheme},
		{"singlelabel", "", ErrInvalidDomain},
		{"", "", ErrEmptyDomain},
		{"-bad.example.com", "", ErrInvalidDomain},
	}

	for _, tt := range tests {
		got, err := NormalizeDomain(tt.in)
		if tt.wantErr != nil {
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NormalizeDomain(%q) err=%v, want %v", tt.in, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Fatalf("NormalizeDomain(%q) unexpected err: %v", tt.in, err)
		}
		if got != tt.want {
			t.Fatalf("NormalizeDomain(%q)=%q, want %q", tt.in, got, tt.want)
		}
	}
}

type fakeResolver struct {
	cname string
	err   error
}

func (f fakeResolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	return f.cname, f.err
}

func TestVerifyCNAME(t *testing.T) {
	ctx := context.Background()

	ok := VerifyCNAME(ctx, fakeResolver{cname: "track.platform.example."}, "clicks.brand.com", "track.platform.example")
	if ok.Status != StatusVerified {
		t.Fatalf("verified: %+v", ok)
	}

	wrong := VerifyCNAME(ctx, fakeResolver{cname: "other.example."}, "clicks.brand.com", "track.platform.example")
	if wrong.Status != StatusFailed {
		t.Fatalf("failed wrong target: %+v", wrong)
	}

	pending := VerifyCNAME(ctx, fakeResolver{err: &net.DNSError{Err: "no such host", IsNotFound: true}}, "clicks.brand.com", "track.platform.example")
	if pending.Status != StatusPending {
		t.Fatalf("pending not found: %+v", pending)
	}

	same := VerifyCNAME(ctx, fakeResolver{cname: "clicks.brand.com."}, "clicks.brand.com", "track.platform.example")
	if same.Status != StatusPending {
		t.Fatalf("pending self: %+v", same)
	}
}

func TestResolveTrackingBase(t *testing.T) {
	if got := ResolveTrackingBase("Clicks.Brand.com.", "https://platform.example", "http://root.example"); got != "https://clicks.brand.com" {
		t.Fatalf("verified domain: %s", got)
	}
	if got := ResolveTrackingBase("", "https://platform.example/", "http://root.example"); got != "https://platform.example" {
		t.Fatalf("tracking_url: %s", got)
	}
	if got := ResolveTrackingBase("", "", "http://root.example/"); got != "http://root.example" {
		t.Fatalf("root_url: %s", got)
	}
}

func TestNormalizeTrackingURL(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"tracking.platform.example", "https://tracking.platform.example", false},
		{"https://tracking.platform.example/", "https://tracking.platform.example", false},
		{"http://localhost:9000", "http://localhost:9000", false},
		{"https://tracking.platform.example/path", "", true},
		{"https://user@tracking.platform.example", "", true},
		{"ftp://tracking.platform.example", "", true},
	}
	for _, tt := range tests {
		got, err := NormalizeTrackingURL(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("NormalizeTrackingURL(%q) expected an error", tt.in)
			}
			continue
		}
		if err != nil || got != tt.want {
			t.Fatalf("NormalizeTrackingURL(%q)=(%q, %v), want %q", tt.in, got, err, tt.want)
		}
	}
}

func TestOwnerIDFromAuthIgnoresBody(t *testing.T) {
	if got := OwnerIDFromAuth(42, 99); got != 42 {
		t.Fatalf("got %d", got)
	}
}

func TestBlocksDelete(t *testing.T) {
	if !BlocksDelete([]string{"draft", "running"}) {
		t.Fatal("running should block")
	}
	if BlocksDelete([]string{"draft", "finished"}) {
		t.Fatal("finished should not block")
	}
}

func TestExpectedCNAMETarget(t *testing.T) {
	if got := ExpectedCNAMETarget("https://Track.Example.com/path", "http://root.example"); got != "track.example.com" {
		t.Fatalf("got %q", got)
	}
	if got := ExpectedCNAMETarget("", "https://Root.Example.com"); got != "root.example.com" {
		t.Fatalf("fallback got %q", got)
	}
}

func TestDNSHostLabel(t *testing.T) {
	tests := []struct {
		host, base, want string
	}{
		{"emailtrack.eguardian.in", "eguardian.in", "emailtrack"},
		{"click.example.com", "example.com", "click"},
		{"example.com", "example.com", "@"},
	}
	for _, tt := range tests {
		if got := DNSHostLabel(tt.host, tt.base); got != tt.want {
			t.Fatalf("DNSHostLabel(%q, %q)=%q, want %q", tt.host, tt.base, got, tt.want)
		}
	}
}

func TestIsUnderBase(t *testing.T) {
	if !IsUnderBase("emailtrack.eguardian.in", "eguardian.in") {
		t.Fatal("subdomain should be under base")
	}
	if IsUnderBase("other.com", "eguardian.in") {
		t.Fatal("unrelated domain should not match")
	}
}
