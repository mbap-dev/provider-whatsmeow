package provider

import (
	"context"
	"testing"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func registeredLIDResponse() types.IsOnWhatsAppResponse {
	return types.IsOnWhatsAppResponse{
		Query:       "+5562933000233",
		IsIn:        true,
		JID:         types.NewJID("46849671094330", types.HiddenUserServer),
		PhoneNumber: types.NewJID("5562933000233", types.DefaultUserServer),
	}
}

func TestCheckUserExistsOnWhatsAppAcceptsCanonicalLID(t *testing.T) {
	oldLookup := isOnWhatsApp
	t.Cleanup(func() { isOnWhatsApp = oldLookup })

	calls := 0
	isOnWhatsApp = func(_ context.Context, _ *whatsmeow.Client, phones []string) ([]types.IsOnWhatsAppResponse, error) {
		calls++
		if len(phones) != 1 || phones[0] != "+5562933000233" {
			t.Fatalf("lookup phones = %#v", phones)
		}
		return []types.IsOnWhatsAppResponse{registeredLIDResponse()}, nil
	}

	manager := &ClientManager{pnCache: make(map[string]pnCacheItem)}
	want := types.NewJID("46849671094330", types.HiddenUserServer)
	got, found, err := manager.checkUserExistsOnWhatsApp(context.Background(), nil, "5562933000233")
	if err != nil {
		t.Fatalf("checkUserExistsOnWhatsApp: %v", err)
	}
	if !found {
		t.Fatal("registered LID recipient reported as missing")
	}
	if got != want {
		t.Fatalf("recipient JID = %s, want %s", got, want)
	}
	if calls != 1 {
		t.Fatalf("lookup calls = %d, want 1", calls)
	}
	if cached := manager.pnCache["+5562933000233"].val; cached != want.String() {
		t.Fatalf("cached JID = %q, want %q", cached, want.String())
	}
}

func TestResolvePNtoJIDAcceptsCanonicalLID(t *testing.T) {
	oldLookup := isOnWhatsApp
	t.Cleanup(func() { isOnWhatsApp = oldLookup })

	calls := 0
	isOnWhatsApp = func(_ context.Context, _ *whatsmeow.Client, phones []string) ([]types.IsOnWhatsAppResponse, error) {
		calls++
		if len(phones) != 1 || phones[0] != "+5562933000233" {
			t.Fatalf("lookup phones = %#v", phones)
		}
		return []types.IsOnWhatsAppResponse{registeredLIDResponse()}, nil
	}

	manager := &ClientManager{pnCache: make(map[string]pnCacheItem)}
	want := types.NewJID("46849671094330", types.HiddenUserServer)
	got, found := manager.resolvePNtoJID(context.Background(), nil, "5562933000233")
	if !found {
		t.Fatal("registered LID recipient was not resolved")
	}
	if got != want {
		t.Fatalf("recipient JID = %s, want %s", got, want)
	}
	if calls != 1 {
		t.Fatalf("lookup calls = %d, want 1", calls)
	}
}

func TestCanonicalPrivateRecipient(t *testing.T) {
	pn := types.NewJID("5562933000233", types.DefaultUserServer)
	lid := types.NewJID("46849671094330", types.HiddenUserServer)
	group := types.NewJID("120363369993948203", types.GroupServer)

	tests := []struct {
		name string
		info types.IsOnWhatsAppResponse
		want types.JID
		ok   bool
	}{
		{name: "phone number", info: types.IsOnWhatsAppResponse{IsIn: true, JID: pn}, want: pn, ok: true},
		{name: "LID", info: types.IsOnWhatsAppResponse{IsIn: true, JID: lid}, want: lid, ok: true},
		{name: "not registered", info: types.IsOnWhatsAppResponse{IsIn: false, JID: lid}},
		{name: "empty JID", info: types.IsOnWhatsAppResponse{IsIn: true}},
		{name: "group", info: types.IsOnWhatsAppResponse{IsIn: true, JID: group}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := canonicalPrivateRecipient(tt.info)
			if ok != tt.ok || got != tt.want {
				t.Fatalf("canonicalPrivateRecipient() = %s, %t, want %s, %t", got, ok, tt.want, tt.ok)
			}
		})
	}
}
