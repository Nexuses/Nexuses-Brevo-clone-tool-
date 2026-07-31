package subimporter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knadh/listmonk/internal/i18n"
	"github.com/knadh/listmonk/models"
)

func TestComposeSubscriberName(t *testing.T) {
	tests := []struct {
		name, first, last, want string
	}{
		{"Jane Doe", "Bob", "Smith", "Jane Doe"},
		{"", "Jane", "Doe", "Jane Doe"},
		{"", "Jane", "", "Jane"},
		{"", "", "", ""},
	}
	for _, tt := range tests {
		got := composeSubscriberName(tt.name, tt.first, tt.last)
		if got != tt.want {
			t.Errorf("composeSubscriberName(%q, %q, %q) = %q, want %q", tt.name, tt.first, tt.last, got, tt.want)
		}
	}
}

func TestCanonicalCSVHeader(t *testing.T) {
	tests := map[string]string{
		"First Name":    "firstname",
		"last_name":     "lastname",
		"Company Name":  "company",
		"Email Address": "email",
		"unknown":       "",
	}
	for in, want := range tests {
		if got := canonicalCSVHeader(in); got != want {
			t.Fatalf("canonicalCSVHeader(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSplitSubscriberName(t *testing.T) {
	first, last := splitSubscriberName("Diksha Birla")
	if first != "Diksha" || last != "Birla" {
		t.Fatalf("splitSubscriberName = %q %q", first, last)
	}
	first, last = splitSubscriberName("Diksha")
	if first != "Diksha" || last != "" {
		t.Fatalf("splitSubscriberName single = %q %q", first, last)
	}
}

func TestApplyCSVFields(t *testing.T) {
	s := &Session{}
	sub := &SubReq{}
	row := map[string]string{
		"email":      "a@b.com",
		"firstname":  "Jane",
		"lastname":   "Doe",
		"company":    "Acme",
		"attributes": `{"city":"London"}`,
	}
	s.applyCSVFields(sub, row, 1)

	if sub.Name != "Jane Doe" {
		t.Fatalf("name = %q, want Jane Doe", sub.Name)
	}
	if sub.Attribs["firstname"] != "Jane" || sub.Attribs["lastname"] != "Doe" || sub.Attribs["company"] != "Acme" {
		t.Fatalf("attribs = %#v", sub.Attribs)
	}
	if sub.Attribs["city"] != "London" {
		t.Fatalf("attributes merge failed: %#v", sub.Attribs)
	}
}

func TestApplyCSVFieldsLegacyName(t *testing.T) {
	s := &Session{}
	sub := &SubReq{}
	row := map[string]string{
		"email": "a@b.com",
		"name":  "Diksha Birla",
	}
	s.applyCSVFields(sub, row, 1)

	if sub.Attribs["firstname"] != "Diksha" || sub.Attribs["lastname"] != "Birla" {
		t.Fatalf("legacy name split failed: %#v", sub.Attribs)
	}
}

func TestLoadCSVImportsDedicatedFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subs.csv")
	csv := "email,firstname,lastname,company\n" +
		"fresh@example.com,Jane,Doe,Nexuses\n"
	if err := os.WriteFile(path, []byte(csv), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	lang, err := os.ReadFile("../../i18n/en.json")
	if err != nil {
		t.Fatalf("read i18n: %v", err)
	}
	i, err := i18n.New(lang)
	if err != nil {
		t.Fatalf("i18n: %v", err)
	}
	im := New(Options{}, nil, i)
	sess, err := im.NewSession(SessionOpt{Filename: "subs.csv", Mode: ModeSubscribe, SubStatus: models.SubscriptionStatusConfirmed})
	if err != nil {
		t.Fatalf("new session: %v", err)
	}

	done := make(chan SubReq, 1)
	go func() {
		for sub := range sess.subQueue {
			done <- sub
		}
		close(done)
	}()

	if err := sess.LoadCSV(path, ','); err != nil {
		t.Fatalf("load csv: %v", err)
	}

	sub, ok := <-done
	if !ok {
		t.Fatalf("no subscriber imported")
	}
	if sub.Email != "fresh@example.com" || sub.Name != "Jane Doe" {
		t.Fatalf("subscriber = %+v", sub)
	}
	if sub.Attribs["firstname"] != "Jane" || sub.Attribs["lastname"] != "Doe" || sub.Attribs["company"] != "Nexuses" {
		t.Fatalf("attribs = %#v", sub.Attribs)
	}
}

func TestSubscriberNameHelpers(t *testing.T) {
	sub := models.Subscriber{
		Name: "Fallback Name",
		Attribs: models.JSON{
			"firstname": "Jane",
			"lastname":  "Doe",
			"company":   "Acme Corp",
		},
	}
	if sub.FirstName() != "Jane" || sub.LastName() != "Doe" || sub.Company() != "Acme Corp" {
		t.Fatalf("FirstName=%q LastName=%q Company=%q", sub.FirstName(), sub.LastName(), sub.Company())
	}
}
