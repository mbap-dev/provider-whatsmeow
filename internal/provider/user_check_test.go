package provider

import "testing"

func TestCheckUserCandidatesAcceptsJID(t *testing.T) {
	got := checkUserCandidates("5511999999999@s.whatsapp.net", false)
	if len(got) != 1 || got[0] != "5511999999999" {
		t.Fatalf("checkUserCandidates() = %#v, want JID user only", got)
	}
}

func TestCheckUserCandidatesNormalizesBRNationalNumber(t *testing.T) {
	got := checkUserCandidates("11 99999-9999", false)
	if len(got) == 0 || got[0] != "+5511999999999" {
		t.Fatalf("checkUserCandidates() = %#v, want first candidate +5511999999999", got)
	}
}
