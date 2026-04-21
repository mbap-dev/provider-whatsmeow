package provider

import (
	"testing"

	"go.mau.fi/whatsmeow/store"
)

func TestParseWhatsAppWebVersionContent(t *testing.T) {
	content := `self.__WB_MANIFEST=[];{"client_revision":1037793778,"foo":"bar"}`

	got, err := parseWhatsAppWebVersionContent(content)
	if err != nil {
		t.Fatalf("parseWhatsAppWebVersionContent() error = %v", err)
	}

	want := store.WAVersionContainer{2, 3000, 1037793778}
	if got != want {
		t.Fatalf("parseWhatsAppWebVersionContent() = %v, want %v", got, want)
	}
}

func TestParseWhatsAppWebVersionContentEscapedPattern(t *testing.T) {
	content := `const payload="{\"client_revision\":1037793778}";`

	got, err := parseWhatsAppWebVersionContent(content)
	if err != nil {
		t.Fatalf("parseWhatsAppWebVersionContent() error = %v", err)
	}

	want := store.WAVersionContainer{2, 3000, 1037793778}
	if got != want {
		t.Fatalf("parseWhatsAppWebVersionContent() = %v, want %v", got, want)
	}
}

func TestApplyWhatsAppWebVersionUpdatesGlobalPayloads(t *testing.T) {
	previous := store.GetWAVersion()
	t.Cleanup(func() {
		ApplyWhatsAppWebVersion(previous)
	})

	target := store.WAVersionContainer{2, 3000, 1037753511}
	ApplyWhatsAppWebVersion(target)

	if got := store.GetWAVersion(); got != target {
		t.Fatalf("GetWAVersion() = %v, want %v", got, target)
	}

	appVersion := store.BaseClientPayload.GetUserAgent().GetAppVersion()
	if appVersion.GetPrimary() != target[0] || appVersion.GetSecondary() != target[1] || appVersion.GetTertiary() != target[2] {
		t.Fatalf(
			"BaseClientPayload app version = %d.%d.%d, want %d.%d.%d",
			appVersion.GetPrimary(),
			appVersion.GetSecondary(),
			appVersion.GetTertiary(),
			target[0],
			target[1],
			target[2],
		)
	}

	if store.DeviceProps.GetVersion().GetPrimary() != target[0] ||
		store.DeviceProps.GetVersion().GetSecondary() != target[1] ||
		store.DeviceProps.GetVersion().GetTertiary() != target[2] {
		t.Fatalf(
			"DeviceProps version = %d.%d.%d, want %d.%d.%d",
			store.DeviceProps.GetVersion().GetPrimary(),
			store.DeviceProps.GetVersion().GetSecondary(),
			store.DeviceProps.GetVersion().GetTertiary(),
			target[0],
			target[1],
			target[2],
		)
	}
}
