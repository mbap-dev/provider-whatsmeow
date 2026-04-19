package provider

import "testing"

func TestApplyInstanceMessageIDs(t *testing.T) {
	payload := map[string]any{
		"entry": []any{
			map[string]any{
				"changes": []any{
					map[string]any{
						"value": map[string]any{
							"messages": []any{
								map[string]any{
									"id": "ABC123",
									"context": map[string]any{
										"id":         "CTX1",
										"message_id": "CTX2",
									},
									"reaction": map[string]any{
										"message_id": "REA1",
									},
								},
							},
							"statuses": []any{
								map[string]any{
									"id":         "STS1",
									"message_id": "STS2",
								},
							},
						},
					},
				},
			},
		},
	}

	applyInstanceMessageIDs("5562999777968", payload)

	value := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	msg := value["messages"].([]any)[0].(map[string]any)
	if got := msg["id"].(string); got != "5562999777968_ABC123" {
		t.Fatalf("message id = %q", got)
	}
	ctx := msg["context"].(map[string]any)
	if got := ctx["id"].(string); got != "5562999777968_CTX1" {
		t.Fatalf("context id = %q", got)
	}
	if got := ctx["message_id"].(string); got != "5562999777968_CTX2" {
		t.Fatalf("context message_id = %q", got)
	}
	reaction := msg["reaction"].(map[string]any)
	if got := reaction["message_id"].(string); got != "5562999777968_REA1" {
		t.Fatalf("reaction message_id = %q", got)
	}
	status := value["statuses"].([]any)[0].(map[string]any)
	if got := status["id"].(string); got != "5562999777968_STS1" {
		t.Fatalf("status id = %q", got)
	}
	if got := status["message_id"].(string); got != "5562999777968_STS2" {
		t.Fatalf("status message_id = %q", got)
	}
}

func TestApplyInstanceMessageIDsMessageEchoes(t *testing.T) {
	payload := map[string]any{
		"entry": []any{
			map[string]any{
				"changes": []any{
					map[string]any{
						"value": map[string]any{
							"message_echoes": []any{
								map[string]any{"id": "ECHO1"},
							},
						},
					},
				},
			},
		},
	}

	applyInstanceMessageIDs("5562999777968", payload)

	value := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	echo := value["message_echoes"].([]any)[0].(map[string]any)
	if got := echo["id"].(string); got != "5562999777968_ECHO1" {
		t.Fatalf("message echo id = %q", got)
	}
}

func TestNormalizeOutgoingMessageIDs(t *testing.T) {
	msg := OutgoingMessage{
		MessageID: "5562999777968_RAW1",
		Context: &MessageContext{
			ID:        "5562999777968_CTX1",
			MessageID: "5562999777968_CTX2",
			StanzaId:  "5562999777968_CTX3",
			Stanza_id: "5562999777968_CTX4",
		},
	}

	normalizeOutgoingMessageIDs("5562999777968", &msg)

	if msg.MessageID != "RAW1" {
		t.Fatalf("message id = %q", msg.MessageID)
	}
	if msg.Context.ID != "CTX1" {
		t.Fatalf("context id = %q", msg.Context.ID)
	}
	if msg.Context.MessageID != "CTX2" {
		t.Fatalf("context message_id = %q", msg.Context.MessageID)
	}
	if msg.Context.StanzaId != "CTX3" {
		t.Fatalf("context stanzaId = %q", msg.Context.StanzaId)
	}
	if msg.Context.Stanza_id != "CTX4" {
		t.Fatalf("context stanza_id = %q", msg.Context.Stanza_id)
	}
}

func TestCompatibleMediaIDPrefersRawMediaID(t *testing.T) {
	meta := MediaMetadata{
		SessionID: "553121159080",
		MediaID:   "3B0D861C2EE5F41B2A7A",
		ID:        "553121159080/3B0D861C2EE5F41B2A7A",
	}

	if got := CompatibleMediaID(meta); got != "3B0D861C2EE5F41B2A7A" {
		t.Fatalf("CompatibleMediaID() = %q", got)
	}
}

func TestCompatibleMediaIDFallsBackToPublicRef(t *testing.T) {
	meta := MediaMetadata{
		ID: "553121159080/3B0D861C2EE5F41B2A7A",
	}

	if got := CompatibleMediaID(meta); got != "3B0D861C2EE5F41B2A7A" {
		t.Fatalf("CompatibleMediaID() = %q", got)
	}
}
