package http

import "testing"

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
