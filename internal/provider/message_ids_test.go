package provider

import (
	"testing"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

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

func TestGroupParticipantMessageID(t *testing.T) {
	got := withGroupParticipantMessageID("553121159080@s.whatsapp.net", "3EB07BAED2DCAC54DE03")
	want := "553121159080~3EB07BAED2DCAC54DE03"
	if got != want {
		t.Fatalf("withGroupParticipantMessageID() = %q, want %q", got, want)
	}

	participant, msgID, ok := splitGroupParticipantMessageID(got)
	if !ok {
		t.Fatal("expected group participant message id to split")
	}
	if participant != "553121159080" {
		t.Fatalf("participant = %q", participant)
	}
	if msgID != "3EB07BAED2DCAC54DE03" {
		t.Fatalf("message id = %q", msgID)
	}
}

func TestGroupParticipantMessageIDKeepsPlainIDs(t *testing.T) {
	participant, msgID, ok := splitGroupParticipantMessageID("3EB07BAED2DCAC54DE03")
	if ok {
		t.Fatal("did not expect plain id to split")
	}
	if participant != "" || msgID != "3EB07BAED2DCAC54DE03" {
		t.Fatalf("split plain = participant:%q msgID:%q", participant, msgID)
	}
}

func TestNormalizeOutgoingMessageIDsKeepsLegacyGroupReplyIDCompatible(t *testing.T) {
	msg := OutgoingMessage{
		Context: &MessageContext{
			MessageID: "553121159080_3EB07BAED2DCAC54DE03",
		},
	}

	normalizeOutgoingMessageIDs("553121159080", &msg)

	if msg.Context.MessageID != "3EB07BAED2DCAC54DE03" {
		t.Fatalf("legacy context message_id = %q", msg.Context.MessageID)
	}
	if participant, msgID, ok := splitGroupParticipantMessageID(msg.Context.MessageID); ok || participant != "" || msgID != "3EB07BAED2DCAC54DE03" {
		t.Fatalf("legacy id parsed as new format: ok=%t participant=%q msgID=%q", ok, participant, msgID)
	}
}

func TestGroupParticipantMessageIDDoesNotCollideWithInstancePrefix(t *testing.T) {
	groupID := withGroupParticipantMessageID("553121159080", "3EB07BAED2DCAC54DE03")
	publicID := withInstanceMessageID("553121159080", groupID)
	want := "553121159080_553121159080~3EB07BAED2DCAC54DE03"
	if publicID != want {
		t.Fatalf("public group id = %q, want %q", publicID, want)
	}

	stripped := stripInstanceMessageID("553121159080", publicID)
	participant, msgID, ok := splitGroupParticipantMessageID(stripped)
	if !ok {
		t.Fatal("expected stripped group id to split")
	}
	if participant != "553121159080" || msgID != "3EB07BAED2DCAC54DE03" {
		t.Fatalf("split stripped = participant:%q msgID:%q", participant, msgID)
	}
}

func TestGroupParticipantMessageIDStripsDeviceSuffix(t *testing.T) {
	got := withGroupParticipantMessageID("5562933000233:3@s.whatsapp.net", "3B4D9307CBB2F3D99403")
	want := "5562933000233~3B4D9307CBB2F3D99403"
	if got != want {
		t.Fatalf("withGroupParticipantMessageID() = %q, want %q", got, want)
	}

	participant, msgID, ok := splitGroupParticipantMessageID(got)
	if !ok {
		t.Fatal("expected group participant message id to split")
	}
	if participant != "5562933000233" || msgID != "3B4D9307CBB2F3D99403" {
		t.Fatalf("split = participant:%q msgID:%q", participant, msgID)
	}
}

func TestMessageContextInfoKeepsPrivateReplyIDPlain(t *testing.T) {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String("testando"),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    proto.String("3B8C1C3F75F1EC9F56DC"),
				Participant: proto.String("46849671094330@lid"),
			},
		},
	}

	ctx := messageContextInfo("553121159080", msg, false)
	if ctx == nil {
		t.Fatal("expected context")
	}
	if got := ctx["message_id"]; got != "3B8C1C3F75F1EC9F56DC" {
		t.Fatalf("private context message_id = %q", got)
	}
}

func TestMessageContextInfoKeepsOwnPrivateReplyParticipant(t *testing.T) {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String("testando"),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    proto.String("3B8C1C3F75F1EC9F56DC"),
				Participant: proto.String("553121159080@s.whatsapp.net"),
			},
		},
	}

	ctx := messageContextInfo("553121159080", msg, false)
	if ctx == nil {
		t.Fatal("expected context")
	}
	if got := ctx["message_id"]; got != "553121159080~3B8C1C3F75F1EC9F56DC" {
		t.Fatalf("own private context message_id = %q", got)
	}
}

func TestMessageContextInfoAddsParticipantForGroupReplyID(t *testing.T) {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String("testando"),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    proto.String("3B8C1C3F75F1EC9F56DC"),
				Participant: proto.String("5562933000233:3@s.whatsapp.net"),
			},
		},
	}

	ctx := messageContextInfo("553121159080", msg, true)
	if ctx == nil {
		t.Fatal("expected context")
	}
	if got := ctx["message_id"]; got != "5562933000233~3B8C1C3F75F1EC9F56DC" {
		t.Fatalf("group context message_id = %q", got)
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
