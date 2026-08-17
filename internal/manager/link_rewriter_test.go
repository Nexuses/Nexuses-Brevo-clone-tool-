package manager

import (
	"bytes"
	"strings"
	"testing"

	"github.com/knadh/listmonk/models"
	"golang.org/x/net/html"
)

func TestRewriteHTMLAnchors(t *testing.T) {
	input := []byte(`<html><body>
<a class="button" href="https://example.com/pricing?plan=pro&amp;src=email">Pricing</a>
<a href="mailto:sales@example.com">Mail</a>
<a href="tel:+15551234567">Call</a>
<a href="#details">Details</a>
<img src="https://cdn.example.com/logo.png" alt="Logo">
</body></html>`)

	output, err := rewriteHTMLAnchors(input, func(value string) string {
		if value == "https://example.com/pricing?plan=pro&src=email" {
			return "https://click.customer.com/link/token/campaign/contact"
		}
		return value
	})
	if err != nil {
		t.Fatal(err)
	}

	got := string(output)
	if !strings.Contains(got, `href="https://click.customer.com/link/token/campaign/contact"`) {
		t.Fatalf("HTTP link was not rewritten: %s", got)
	}
	for _, unchanged := range []string{
		`href="mailto:sales@example.com"`,
		`href="tel:+15551234567"`,
		`href="#details"`,
		`src="https://cdn.example.com/logo.png"`,
	} {
		if !strings.Contains(got, unchanged) {
			t.Fatalf("expected %q to remain intact: %s", unchanged, got)
		}
	}

	if _, err := html.Parse(bytes.NewReader(output)); err != nil {
		t.Fatalf("rewritten HTML is invalid: %v", err)
	}
}

func TestTrackableHrefFiltering(t *testing.T) {
	m := &Manager{cfg: Config{
		RootURL:     "https://app.platform.example",
		TrackingURL: "https://tracking.platform.example",
	}}
	campaign := &models.Campaign{}

	tests := []struct {
		href string
		want bool
	}{
		{"https://customer.example/pricing?plan=pro", true},
		{"mailto:sales@example.com", false},
		{"tel:+15551234567", false},
		{"#section", false},
		{"https://tracking.platform.example/link/a/b/c", false},
		{"https://app.platform.example/subscription/a/b", false},
		{"https://app.platform.example/archive/campaign", false},
		{"not a URL", false},
	}
	for _, tt := range tests {
		if got := m.isTrackableHref(tt.href, campaign); got != tt.want {
			t.Fatalf("isTrackableHref(%q)=%v, want %v", tt.href, got, tt.want)
		}
	}
}
