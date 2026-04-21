package provider

import "testing"

func TestFormatMentionDisplayIncludesNameAndNumber(t *testing.T) {
	got := formatMentionDisplay("Tal Tal", "5562999999999")
	want := "@Tal Tal (5562999999999)"
	if got != want {
		t.Fatalf("formatMentionDisplay() = %q, want %q", got, want)
	}
}

func TestApplyMentionReplacementsReplacesOpaqueMentionIdentifiers(t *testing.T) {
	got := applyMentionReplacements(
		"Oi @175556939227181, fale com @5562990000000 depois.",
		map[string]string{
			"@175556939227181": "@Tal Tal (5562999999999)",
		},
	)
	want := "Oi @Tal Tal (5562999999999), fale com @5562990000000 depois."
	if got != want {
		t.Fatalf("applyMentionReplacements() = %q, want %q", got, want)
	}
}
