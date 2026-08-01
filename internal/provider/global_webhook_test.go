package provider

import (
	"context"
	"reflect"
	"testing"

	"go.mau.fi/whatsmeow/types/events"
)

func TestGlobalWebhookFileStoreRoundTrip(t *testing.T) {
	store := NewFileStore(t.TempDir())
	ctx := context.Background()
	const workerID = "worker-test"
	hooks := []SessionWebhook{
		{
			ID:                      "disconnect",
			URLAbsolute:             "https://example.com/hooks/global",
			Enabled:                 true,
			TimeoutMS:               5000,
			SendSessionDisconnected: true,
		},
	}

	if err := store.SaveGlobalWebhooks(ctx, workerID, hooks); err != nil {
		t.Fatalf("SaveGlobalWebhooks(): %v", err)
	}
	got, err := store.LoadGlobalWebhooks(ctx, workerID)
	if err != nil {
		t.Fatalf("LoadGlobalWebhooks(): %v", err)
	}
	if !reflect.DeepEqual(hooks, got) {
		t.Fatalf("global webhooks mismatch: %+v vs %+v", hooks, got)
	}
}

func TestShouldDeliverSessionDisconnectedWebhook(t *testing.T) {
	class := webhookPayloadClass{SessionDisconnected: true}
	if !shouldDeliverSessionWebhook(SessionWebhook{Enabled: true, SendSessionDisconnected: true}, class) {
		t.Fatal("expected enabled disconnect webhook to match disconnect event")
	}
	if shouldDeliverSessionWebhook(SessionWebhook{Enabled: true, SendSessionDisconnected: false}, class) {
		t.Fatal("expected webhook without disconnect filter to skip disconnect event")
	}
}

func TestShouldDeliverSessionOfflineWebhook(t *testing.T) {
	class := webhookPayloadClass{SessionOffline: true}
	if !shouldDeliverSessionWebhook(SessionWebhook{Enabled: true, SendSessionOffline: true}, class) {
		t.Fatal("expected enabled offline webhook to match offline event")
	}
	if shouldDeliverSessionWebhook(SessionWebhook{Enabled: true, SendSessionOffline: false}, class) {
		t.Fatal("expected webhook without offline filter to skip offline event")
	}
}

func TestClassifyWebhookPayloadTreatsEditedGroupMessageAsUpdate(t *testing.T) {
	payload := cloudEnvelope("553121159080")
	value := extractCloudPayloadValue(payload)
	value["messages"] = []any{
		map[string]any{
			"from":   "5562933000233",
			"id":     "5562933000233~3B6F8B989F3D465157DF",
			"type":   "text",
			"edited": true,
			"text": map[string]any{
				"body": "texto novo\n`Mensagem Editada`",
			},
		},
	}
	value["contacts"] = []any{
		map[string]any{
			"wa_id":    "5562933000233",
			"group_id": "120363404017899025@g.us",
		},
	}

	class := classifyWebhookPayload("553121159080", payload)
	if !class.Update {
		t.Fatal("expected edited message to be classified as update")
	}
	if !class.Incoming {
		t.Fatal("expected edited group member message to be classified as incoming")
	}
	if !class.Group {
		t.Fatal("expected edited group member message to be classified as group")
	}
}

func TestSessionDisconnectedReasonIgnoresTransientSocketEvents(t *testing.T) {
	for _, evt := range []interface{}{
		&events.Disconnected{},
		&events.StreamReplaced{},
		&events.ConnectFailure{Reason: events.ConnectFailureServiceUnavailable},
	} {
		if reason, ok := sessionDisconnectedReasonForEvent(evt); ok {
			t.Fatalf("sessionDisconnectedReasonForEvent(%T) = %q, true; want false", evt, reason)
		}
	}
}

func TestSessionDisconnectedReasonIncludesLogoutAndRestrictionEvents(t *testing.T) {
	cases := []interface{}{
		&events.LoggedOut{Reason: events.ConnectFailureLoggedOut},
		&events.TemporaryBan{Code: events.TempBanBlockedByUsers},
		&events.ConnectFailure{Reason: events.ConnectFailureUnknownLogout},
	}
	for _, evt := range cases {
		if reason, ok := sessionDisconnectedReasonForEvent(evt); !ok || reason == "" {
			t.Fatalf("sessionDisconnectedReasonForEvent(%T) = %q, %v; want reason,true", evt, reason, ok)
		}
	}
}
