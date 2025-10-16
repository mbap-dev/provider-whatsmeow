package provider

import (
	"strings"
	"testing"
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
