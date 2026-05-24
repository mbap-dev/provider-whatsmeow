package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"your.org/provider-whatsmeow/internal/config"
	"your.org/provider-whatsmeow/internal/provider"
)

func TestAPIResponseMessageIDPrefixesSession(t *testing.T) {
	got := apiResponseMessageID("553121159080", "INC-30F6716468955937")
	want := "553121159080_INC-30F6716468955937"
	if got != want {
		t.Fatalf("apiResponseMessageID() = %q, want %q", got, want)
	}
}

func TestAPIResponseMessageIDKeepsExistingPrefix(t *testing.T) {
	got := apiResponseMessageID("553121159080", "553121159080_INC-30F6716468955937")
	want := "553121159080_INC-30F6716468955937"
	if got != want {
		t.Fatalf("apiResponseMessageID() = %q, want %q", got, want)
	}
}

func TestAPIResponseMessageIDPrefixesAuthoredID(t *testing.T) {
	authoredID := provider.WithGroupParticipantMessageID("553121159080", "3EB06FA40FC0AA6D317F")
	got := apiResponseMessageID("553121159080", authoredID)
	want := "553121159080_553121159080~3EB06FA40FC0AA6D317F"
	if got != want {
		t.Fatalf("apiResponseMessageID() = %q, want %q", got, want)
	}
}

func TestGenerateAPIMessageIDUsesWhatsAppNativeFormat(t *testing.T) {
	got := generateAPIMessageID("INC")
	if strings.Contains(got, "INC-") {
		t.Fatalf("generateAPIMessageID() = %q, should not contain INC-", got)
	}
	if !strings.HasPrefix(got, "3EB0") {
		t.Fatalf("generateAPIMessageID() = %q, want whatsmeow 3EB0 prefix", got)
	}
}

func TestMediaRequestSessionIDFromPath(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/v19.0/553121159080/A1B2C3", nil)
	got := mediaRequestSessionID(req, map[string]string{"id": "553121159080"})
	if got != "553121159080" {
		t.Fatalf("mediaRequestSessionID() = %q, want %q", got, "553121159080")
	}
}

func TestMediaRequestSessionIDFromQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.test/v19.0/A1B2C3?phone_number_id=553121159080", nil)
	got := mediaRequestSessionID(req, map[string]string{"media_id": "A1B2C3"})
	if got != "553121159080" {
		t.Fatalf("mediaRequestSessionID() = %q, want %q", got, "553121159080")
	}
}

func TestDeleteSessionEndpoint(t *testing.T) {
	root := t.TempDir()
	manager := provider.NewClientManager(provider.ManagerOptions{
		SessionStore: root,
		WorkerID:     "worker-test",
	})

	sessionID := "553121159080"
	sessionDir := filepath.Join(root, sessionID)
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(): %v", err)
	}
	if err := os.WriteFile(filepath.Join(sessionDir, "session.db"), []byte("db"), 0o644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}

	srv := NewServer(&config.Config{}, manager, nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/sessions/"+sessionID, nil)
	rec := httptest.NewRecorder()

	srv.httpSrv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"removed":true`) {
		t.Fatalf("response body = %s, want removed=true", rec.Body.String())
	}
	if _, err := os.Stat(sessionDir); !os.IsNotExist(err) {
		t.Fatalf("session dir still exists after endpoint delete: %v", err)
	}
}

func TestPairingCodeEndpointRequiresPhone(t *testing.T) {
	manager := provider.NewClientManager(provider.ManagerOptions{
		SessionStore: t.TempDir(),
		WorkerID:     "worker-test",
	})

	srv := NewServer(&config.Config{}, manager, nil)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/sessions/553121159080/pairing-code",
		bytes.NewBufferString(`{"phone":"   "}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.httpSrv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "missing phone") {
		t.Fatalf("response body = %s, want missing phone", rec.Body.String())
	}
}

func TestCreateSessionEndpointRejectsInvalidID(t *testing.T) {
	manager := provider.NewClientManager(provider.ManagerOptions{
		SessionStore: t.TempDir(),
		WorkerID:     "worker-test",
	})

	srv := NewServer(&config.Config{}, manager, nil)
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/sessions",
		bytes.NewBufferString(`{"id":"1","assign":true,"connect":true}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.httpSrv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid id") {
		t.Fatalf("response body = %s, want invalid id", rec.Body.String())
	}
}

func TestWorkerEndpointOmitsGlobalWhatsAppWebVersion(t *testing.T) {
	manager := provider.NewClientManager(provider.ManagerOptions{
		SessionStore: t.TempDir(),
		WorkerID:     "worker-test",
	})

	srv := NewServer(&config.Config{}, manager, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/worker", nil)
	rec := httptest.NewRecorder()

	srv.httpSrv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "whatsapp_web_version") {
		t.Fatalf("response body = %s, did not expect global whatsapp_web_version", rec.Body.String())
	}
}

func TestExtractSendMentionAllFromContentAttributes(t *testing.T) {
	raw := map[string]any{
		"content_attributes": map[string]any{
			"unoapi": map[string]any{
				"mention_all": true,
			},
		},
	}

	if !extractSendMentionAll(raw) {
		t.Fatal("expected mention_all from content_attributes.unoapi")
	}
}

func TestExtractSendMentionAllFromTopLevelCamelCase(t *testing.T) {
	raw := map[string]any{
		"mentionAll": "true",
	}

	if !extractSendMentionAll(raw) {
		t.Fatal("expected mentionAll top-level value to be accepted")
	}
}

func TestDecodeOutgoingMessageBodyAcceptsSingleItemArray(t *testing.T) {
	msg, err := decodeOutgoingMessageBody([]byte(`[
		{
			"type":"interactive",
			"to":"5562933000233",
			"interactive":{
				"type":"list",
				"body":{"text":"teste"},
				"action":{"button":"Ver opções","sections":[{"title":"Opções","rows":[{"id":"1","title":"Aro 13"}]}]}
			}
		}
	]`))
	if err != nil {
		t.Fatalf("decodeOutgoingMessageBody() error = %v", err)
	}
	if got := msg.Type; got != "interactive" {
		t.Fatalf("type = %q", got)
	}
	if msg.Interactive == nil || msg.Interactive.Type != "list" {
		t.Fatalf("interactive = %+v, want list", msg.Interactive)
	}
}

func TestDecodeOutgoingMessageBodyAcceptsInteractiveButtonObject(t *testing.T) {
	msg, err := decodeOutgoingMessageBody([]byte(`{
		"messaging_product":"whatsapp",
		"recipient_type":"individual",
		"to":"5562933000233",
		"type":"interactive",
		"interactive":{
			"type":"button",
			"body":{"text":"Opção invalida, tente novamente."},
			"action":{"buttons":[
				{"type":"reply","reply":{"id":"option_1","title":"Sim"}},
				{"type":"reply","reply":{"id":"option_2","title":"Não, pode encerrar."}}
			]}
		}
	}`))
	if err != nil {
		t.Fatalf("decodeOutgoingMessageBody() error = %v", err)
	}
	if got := msg.Type; got != "interactive" {
		t.Fatalf("type = %q", got)
	}
	if msg.Interactive == nil || msg.Interactive.Type != "button" {
		t.Fatalf("interactive = %+v, want button", msg.Interactive)
	}
	if got := msg.Interactive.Action.Buttons[1].Reply.Title; got != "Não, pode encerrar." {
		t.Fatalf("button title = %q", got)
	}
}

func TestDecodeOutgoingMessageBodyAcceptsCloudTemplatePayload(t *testing.T) {
	msg, err := decodeOutgoingMessageBody([]byte(`{
		"messaging_product":"whatsapp",
		"to":"5562933000233",
		"type":"template",
		"template_webhook_message":true,
		"template":{
			"name":"boleto_disponivel",
			"language":{"code":"pt_BR"},
			"components":[
				{"type":"body","parameters":[
					{"type":"text","parameter_name":"nome","text":"Maria"},
					{"type":"text","text":"12345"}
				]}
			],
			"variables":{"vencimento":"10/06/2026"}
		}
	}`))
	if err != nil {
		t.Fatalf("decodeOutgoingMessageBody() error = %v", err)
	}
	if got := msg.Type; got != "template" {
		t.Fatalf("type = %q", got)
	}
	if msg.Template == nil || msg.Template.Name != "boleto_disponivel" {
		t.Fatalf("template = %+v", msg.Template)
	}
	if !msg.TemplateWebhookMessage {
		t.Fatal("expected template_webhook_message")
	}
	if got := msg.Template.Components[0].Parameters[0].ParameterName; got != "nome" {
		t.Fatalf("parameter_name = %q", got)
	}
}
