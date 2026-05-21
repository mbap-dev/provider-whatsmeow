package provider

import (
	"reflect"
	"testing"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func TestBuildInboundListResponsePayload(t *testing.T) {
	payload := buildInboundListResponsePayload(&waE2E.ListResponseMessage{
		Title:       proto.String("Aro 13"),
		Description: proto.String("Linha 13"),
		SingleSelectReply: &waE2E.ListResponseMessage_SingleSelectReply{
			SelectedRowID: proto.String("option_1"),
		},
	})

	want := map[string]any{
		"type": "list_reply",
		"list_reply": map[string]any{
			"id":          "option_1",
			"title":       "Aro 13",
			"description": "Linha 13",
		},
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("payload = %#v, want %#v", payload, want)
	}
}

func TestBuildInboundButtonsResponsePayload(t *testing.T) {
	payload := buildInboundButtonsResponsePayload(&waE2E.ButtonsResponseMessage{
		SelectedButtonID: proto.String("option_2"),
		Response: &waE2E.ButtonsResponseMessage_SelectedDisplayText{
			SelectedDisplayText: "Por descrição",
		},
	})

	want := map[string]any{
		"type": "button_reply",
		"button_reply": map[string]any{
			"id":    "option_2",
			"title": "Por descrição",
		},
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("payload = %#v, want %#v", payload, want)
	}
}

func TestBuildInboundInteractiveResponsePayloadQuickReply(t *testing.T) {
	payload := buildInboundInteractiveResponsePayload(&waE2E.InteractiveResponseMessage{
		Body: &waE2E.InteractiveResponseMessage_Body{
			Text: proto.String("Buscar por aro"),
		},
		InteractiveResponseMessage: &waE2E.InteractiveResponseMessage_NativeFlowResponseMessage_{
			NativeFlowResponseMessage: &waE2E.InteractiveResponseMessage_NativeFlowResponseMessage{
				Name:       proto.String("quick_reply"),
				ParamsJSON: proto.String(`{"id":"option_2","display_text":"Buscar por aro"}`),
				Version:    proto.Int32(1),
			},
		},
	})

	want := map[string]any{
		"type": "button_reply",
		"button_reply": map[string]any{
			"id":    "option_2",
			"title": "Buscar por aro",
		},
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("payload = %#v, want %#v", payload, want)
	}
}

func TestExtractMessageContextInfoFromListResponse(t *testing.T) {
	ctx := extractMessageContextInfo(&waE2E.Message{
		ListResponseMessage: &waE2E.ListResponseMessage{
			ContextInfo: &waE2E.ContextInfo{
				StanzaID: proto.String("MSG-123"),
			},
		},
	})
	if ctx == nil || ctx.GetStanzaID() != "MSG-123" {
		t.Fatalf("context = %#v", ctx)
	}
}
