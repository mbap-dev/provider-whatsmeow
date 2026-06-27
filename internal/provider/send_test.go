package provider

import (
	"context"
	"strings"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
)

func TestToUserJIDDigits(t *testing.T) {
	j, err := toUserJID("5562991728088")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.User != "5562991728088" {
		t.Fatalf("user mismatch: %s", j.User)
	}
	if j.Server != "s.whatsapp.net" {
		t.Fatalf("server mismatch: %s", j.Server)
	}
}

func TestToJIDParsing(t *testing.T) {
	j1, err := toJID("5562991728088@s.whatsapp.net")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j1.Server != "s.whatsapp.net" || j1.User != "5562991728088" {
		t.Fatalf("parsed JID mismatch: %#v", j1)
	}
	j2, err := toJID("5562991728088")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j2.Server != "s.whatsapp.net" || j2.User != "5562991728088" {
		t.Fatalf("digits JID mismatch: %#v", j2)
	}
}

func TestNormalizeAudioMime(t *testing.T) {
	if got := normalizeAudioMime("x", "audio/mpga", false); got != "audio/mpeg" {
		t.Fatalf("normalize mpga => %s", got)
	}
	if got := normalizeAudioMime("x", "application/octet-stream", true); !strings.HasPrefix(got, "audio/ogg") {
		t.Fatalf("normalize ptt => %s", got)
	}
}

func TestShouldSendAsPTT(t *testing.T) {
	if !shouldSendAsPTT("https://x/y.ogg", "audio/ogg") {
		t.Fatal("should detect ogg extension")
	}
	if shouldSendAsPTT("https://x/y.mp3", "audio/mpeg") {
		t.Fatal("mp3 should not default to PTT")
	}
}

func TestDigitsOnly(t *testing.T) {
	if got := digitsOnly(" +55 (62) 91728-088 "); got != "556291728088" {
		t.Fatalf("digitsOnly => %s", got)
	}
}

func TestFilenameByMime(t *testing.T) {
	if got := filenameByMime("abc123", "audio/ogg; codecs=opus"); got != "abc123.ogg" {
		t.Fatalf("audio filename => %s", got)
	}
	if got := filenameByMime("abc123", "application/vnd.ms-excel"); got != "abc123.xls" {
		t.Fatalf("xls filename => %s", got)
	}
	if got := filenameByMime("abc123", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"); got != "abc123.xlsx" {
		t.Fatalf("xlsx filename => %s", got)
	}
	if got := filenameByMime("abc123", "video/mp4"); got != "abc123.mp4" {
		t.Fatalf("video filename => %s", got)
	}
	if got := filenameByMime("abc123", "application/pdf"); got != "abc123.pdf" {
		t.Fatalf("document filename => %s", got)
	}
}

func TestBuildWaveformBounds(t *testing.T) {
	// 1 segundo de silêncio a 16 kHz
	samples := make([]int16, 16000)
	wf := buildWaveform(samples)
	if len(wf) != 64 {
		t.Fatalf("waveform length = %d, want 64", len(wf))
	}
	for i, v := range wf {
		if v > 100 {
			t.Fatalf("wf[%d]=%d > 100", i, v)
		}
		if v < 0 {
			t.Fatalf("wf[%d]=%d < 0", i, v)
		}
	}
}

func TestOutgoingMessageRequestsMentionAllFromExplicitFlag(t *testing.T) {
	msg := OutgoingMessage{MentionAll: true}

	if !outgoingMessageRequestsMentionAll(msg) {
		t.Fatal("expected mention_all flag to be honored")
	}
}

func TestOutgoingMessageRequestsMentionAllFromPrefix(t *testing.T) {
	msg := OutgoingMessage{
		Type: "text",
		Text: &TextContent{Body: "@todos Reuniao as 15h"},
	}

	if !outgoingMessageRequestsMentionAll(msg) {
		t.Fatal("expected @todos prefix to enable mention_all")
	}
}

func TestOutgoingMessageRequestsMentionAllIgnoresInnerWord(t *testing.T) {
	msg := OutgoingMessage{
		Type: "text",
		Text: &TextContent{Body: "Avisar @todos depois"},
	}

	if outgoingMessageRequestsMentionAll(msg) {
		t.Fatal("did not expect inner @todos to enable mention_all")
	}
}

func TestNormalizeMentionTargetsDeduplicates(t *testing.T) {
	got := normalizeMentionTargets([]string{
		"5562999999999",
		"5562999999999@s.whatsapp.net",
		"+5562996303132@s.whatsapp.net",
		" 5562888888888 ",
		"",
	})

	want := []string{
		"5562999999999@s.whatsapp.net",
		"5562996303132@s.whatsapp.net",
		"5562888888888@s.whatsapp.net",
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("normalizeMentionTargets() = %v, want %v", got, want)
	}
}

func TestMentionTargetPhoneLookupExtractsUserJIDPhone(t *testing.T) {
	got := mentionTargetPhoneLookup("+5562996303132@s.whatsapp.net")

	if got != "+5562996303132" {
		t.Fatalf("mentionTargetPhoneLookup() = %q", got)
	}
}

func TestMentionTargetPhoneLookupIgnoresLID(t *testing.T) {
	if got := mentionTargetPhoneLookup("123456@lid"); got != "" {
		t.Fatalf("mentionTargetPhoneLookup() = %q", got)
	}
}

func TestResolveMentionTargetsUsesResolvedJIDForE164Phone(t *testing.T) {
	manager := &ClientManager{
		pnCache: map[string]pnCacheItem{
			"+5562996303132": {
				val: "556296303132@s.whatsapp.net",
				exp: time.Now().Add(time.Hour),
			},
		},
	}

	got := manager.resolveMentionTargets(context.Background(), &whatsmeow.Client{}, []string{"+5562996303132"})
	want := []string{"556296303132@s.whatsapp.net"}

	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("resolveMentionTargets() = %v, want %v", got, want)
	}
}
