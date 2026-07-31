package cloudflare

import (
	"encoding/json"
	"testing"
)

func TestToCloudflareAddr(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"noreply@example.com", `"noreply@example.com"`},
		{"listmonk <noreply@example.com>", `{"address":"noreply@example.com","name":"listmonk"}`},
		{`"Mailing list" <hello@nexuses.in>`, `{"address":"hello@nexuses.in","name":"Mailing list"}`},
	}
	for _, tt := range tests {
		got, err := toCloudflareAddr(tt.in)
		if err != nil {
			t.Fatalf("toCloudflareAddr(%q): %v", tt.in, err)
		}
		b, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != tt.want {
			t.Fatalf("toCloudflareAddr(%q) = %s, want %s", tt.in, b, tt.want)
		}
	}
}
