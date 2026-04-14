package provider

import (
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
