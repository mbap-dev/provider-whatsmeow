package provider

import (
	"context"
	"reflect"
	"testing"
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
