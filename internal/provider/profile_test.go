package provider

import (
	"reflect"
	"testing"

	"go.mau.fi/whatsmeow/types"
)

func TestBuildCloudProfileDefaultsPictureField(t *testing.T) {
	profile := buildCloudProfile("553121159080", nil, types.JID{}, "5562933000233", "MSG1")

	if got := profile["name"]; got != "5562933000233" {
		t.Fatalf("name = %#v", got)
	}
	if got := profile["picture"]; got != "" {
		t.Fatalf("picture = %#v", got)
	}
}

func TestNormalizedAvatarLookupJIDStripsDeviceSuffix(t *testing.T) {
	jid := types.JID{User: "5562999999999:26", Server: "s.whatsapp.net", Device: 26}

	got := normalizedAvatarLookupJID("553121159080", nil, jid)

	if got.User != "5562999999999" {
		t.Fatalf("user = %#v", got.User)
	}
	if got.Server != "s.whatsapp.net" {
		t.Fatalf("server = %#v", got.Server)
	}
	if got.Device != 0 {
		t.Fatalf("device = %#v", got.Device)
	}
}

func TestNormalizedAvatarLookupJIDClearsDeviceField(t *testing.T) {
	jid := types.JID{User: "5562999999999", Server: "s.whatsapp.net", Device: 26}

	got := normalizedAvatarLookupJID("553121159080", nil, jid)

	if got.User != "5562999999999" {
		t.Fatalf("user = %#v", got.User)
	}
	if got.Server != "s.whatsapp.net" {
		t.Fatalf("server = %#v", got.Server)
	}
	if got.Device != 0 {
		t.Fatalf("device = %#v", got.Device)
	}
}

func TestContactInfoDisplayNamePrefersPushName(t *testing.T) {
	got := contactInfoDisplayName(types.ContactInfo{
		PushName:     "Maria Grupo",
		BusinessName: "Empresa Maria",
		FullName:     "Maria da Silva",
		FirstName:    "Maria",
	})
	if got != "Maria Grupo" {
		t.Fatalf("contactInfoDisplayName() = %q, want %q", got, "Maria Grupo")
	}
}

func TestBuildCloudContactObjectForGroupKeepsChatwootShape(t *testing.T) {
	groupJID := types.JID{User: "120363000000000000", Server: "g.us"}
	profile := map[string]any{
		"name":    "Maria Grupo",
		"picture": "https://sender.example/avatar.jpg",
	}

	contact := buildCloudContactObject(
		profile,
		"5562999999999",
		groupJID,
		"Equipe Comercial",
		"https://group.example/avatar.jpg",
	)

	if got := contact["wa_id"]; got != "5562999999999" {
		t.Fatalf("wa_id = %#v", got)
	}
	if got := contact["group_id"]; got != groupJID.String() {
		t.Fatalf("group_id = %#v", got)
	}
	if got := contact["group_subject"]; got != "Equipe Comercial" {
		t.Fatalf("group_subject = %#v", got)
	}
	if got := contact["group_picture"]; got != "https://group.example/avatar.jpg" {
		t.Fatalf("group_picture = %#v", got)
	}
	if got := contact["profile"]; !reflect.DeepEqual(got, profile) {
		t.Fatalf("profile = %#v", got)
	}
}
