package provider

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"google.golang.org/protobuf/proto"
)

// ApplyWhatsAppWebVersion keeps the whatsmeow global version state aligned
// across the user agent, registration hash and device props payloads.
func ApplyWhatsAppWebVersion(version store.WAVersionContainer) store.WAVersionContainer {
	if version.IsZero() {
		return store.GetWAVersion()
	}

	store.SetWAVersion(version)
	if store.BaseClientPayload != nil && store.BaseClientPayload.UserAgent != nil {
		store.BaseClientPayload.UserAgent.AppVersion = version.ProtoAppVersion()
	}
	if store.DeviceProps != nil && store.DeviceProps.Version != nil {
		store.DeviceProps.Version.Primary = proto.Uint32(version[0])
		store.DeviceProps.Version.Secondary = proto.Uint32(version[1])
		store.DeviceProps.Version.Tertiary = proto.Uint32(version[2])
	}

	return version
}

// ConfigureLatestWhatsAppWebVersion tries to fetch the latest WhatsApp Web
// version and applies it globally. On failure, the bundled whatsmeow version
// remains active.
func ConfigureLatestWhatsAppWebVersion(ctx context.Context) (previous store.WAVersionContainer, applied store.WAVersionContainer, err error) {
	previous = store.GetWAVersion()
	latest, err := whatsmeow.GetLatestVersion(ctx, nil)
	if err != nil {
		return previous, previous, err
	}
	applied = ApplyWhatsAppWebVersion(*latest)
	return previous, applied, nil
}
