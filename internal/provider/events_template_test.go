package provider

import "testing"

func TestBuildTemplateWebhookWireMessageText(t *testing.T) {
	msg := buildTemplateWebhookWireMessage("553121159080", "MSG1", &renderedTemplate{
		HeaderText: "Olá",
		Body:       "Corpo",
		Footer:     "Rodapé",
	})
	if got := msg["type"]; got != "text" {
		t.Fatalf("type = %#v", got)
	}
	text := msg["text"].(map[string]any)
	if got := text["body"]; got != "Olá\n\nCorpo\n\nRodapé" {
		t.Fatalf("body = %#v", got)
	}
	if got := msg["from"]; got != "553121159080" {
		t.Fatalf("from = %#v", got)
	}
}

func TestBuildTemplateWebhookWireMessageDocument(t *testing.T) {
	msg := buildTemplateWebhookWireMessage("553121159080", "MSG1", &renderedTemplate{
		Body: "Corpo",
		HeaderDocument: &DocumentContent{
			ID:       "media-id",
			Link:     "https://example.com/boleto.pdf",
			FileName: "boleto.pdf",
		},
	})
	if got := msg["type"]; got != "document" {
		t.Fatalf("type = %#v", got)
	}
	doc := msg["document"].(map[string]any)
	if got := doc["caption"]; got != "Corpo" {
		t.Fatalf("caption = %#v", got)
	}
	if got := doc["id"]; got != "media-id" {
		t.Fatalf("id = %#v", got)
	}
	if got := doc["filename"]; got != "boleto.pdf" {
		t.Fatalf("filename = %#v", got)
	}
	if _, ok := doc["link"]; ok {
		t.Fatalf("document payload should not include link: %#v", doc)
	}
}

func TestBuildTemplateWebhookWireMessageImage(t *testing.T) {
	msg := buildTemplateWebhookWireMessage("553121159080", "MSG1", &renderedTemplate{
		Body:        "Corpo",
		HeaderImage: &ImageContent{ID: "image-id", Link: "https://example.com/image.jpg"},
	})
	if got := msg["type"]; got != "image" {
		t.Fatalf("type = %#v", got)
	}
	image := msg["image"].(map[string]any)
	if got := image["caption"]; got != "Corpo" {
		t.Fatalf("caption = %#v", got)
	}
	if got := image["id"]; got != "image-id" {
		t.Fatalf("id = %#v", got)
	}
	if _, ok := image["link"]; ok {
		t.Fatalf("image payload should not include link: %#v", image)
	}
}

func TestConvertCloudPayloadToOutgoingEchoes(t *testing.T) {
	payload := cloudEnvelope("553121159080")
	change := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)
	value := change["value"].(map[string]any)
	value["messages"] = []any{
		map[string]any{
			"from": "553121159080",
			"id":   "MSG1",
			"type": "text",
			"text": map[string]any{"body": "Corpo"},
		},
	}

	convertCloudPayloadToOutgoingEchoes(payload)

	if got := change["field"]; got != "smb_message_echoes" {
		t.Fatalf("field = %#v", got)
	}
	if _, ok := value["messages"]; ok {
		t.Fatalf("messages should be removed: %#v", value)
	}
	echoes, ok := value["message_echoes"].([]any)
	if !ok || len(echoes) != 1 {
		t.Fatalf("message_echoes = %#v", value["message_echoes"])
	}
	echo := echoes[0].(map[string]any)
	if got := echo["type"]; got != "text" {
		t.Fatalf("echo type = %#v", got)
	}
}
