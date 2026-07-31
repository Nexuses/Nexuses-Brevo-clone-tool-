package models

import (
	"testing"
)

func TestJSONScan(t *testing.T) {
	var attribs JSON
	if err := (&attribs).Scan([]byte(`{"firstname":"Jane","company":"Acme"}`)); err != nil {
		t.Fatalf("scan []byte: %v", err)
	}
	if attribs["firstname"] != "Jane" || attribs["company"] != "Acme" {
		t.Fatalf("attribs = %#v", attribs)
	}

	attribs = nil
	if err := (&attribs).Scan(`{"lastname":"Doe"}`); err != nil {
		t.Fatalf("scan string: %v", err)
	}
	if attribs["lastname"] != "Doe" {
		t.Fatalf("attribs = %#v", attribs)
	}
}

func TestJSONValue(t *testing.T) {
	var nilJSON JSON
	v, err := nilJSON.Value()
	if err != nil {
		t.Fatalf("value nil: %v", err)
	}
	if string(v.([]byte)) != "{}" {
		t.Fatalf("nil value = %s", v)
	}
}
