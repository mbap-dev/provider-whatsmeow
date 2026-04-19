package http

import (
	"net/http/httptest"
	"testing"
)

func TestAPIResponseMessageIDPrefixesSession(t *testing.T) {
	got := apiResponseMessageID("553121159080", "INC-30F6716468955937")
	want := "553121159080_INC-30F6716468955937"
	if got != want {
		t.Fatalf("apiResponseMessageID() = %q, want %q", got, want)
	}
}

func TestAPIResponseMessageIDKeepsExistingPrefix(t *testing.T) {
	got := apiResponseMessageID("553121159080", "553121159080_INC-30F6716468955937")
	want := "553121159080_INC-30F6716468955937"
	if got != want {
		t.Fatalf("apiResponseMessageID() = %q, want %q", got, want)
	}
}

func TestMediaRequestSessionIDFromPath(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/v19.0/553121159080/A1B2C3", nil)
	got := mediaRequestSessionID(req, map[string]string{"id": "553121159080"})
	if got != "553121159080" {
		t.Fatalf("mediaRequestSessionID() = %q, want %q", got, "553121159080")
	}
}

func TestMediaRequestSessionIDFromQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/v19.0/A1B2C3?phone_number_id=553121159080", nil)
	got := mediaRequestSessionID(req, map[string]string{"media_id": "A1B2C3"})
	if got != "553121159080" {
		t.Fatalf("mediaRequestSessionID() = %q, want %q", got, "553121159080")
	}
}
