package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.mau.fi/whatsmeow/store"
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

func TestWorkerEndpointIncludesWhatsAppWebVersion(t *testing.T) {
	previous := store.GetWAVersion()
	t.Cleanup(func() {
		provider.ApplyWhatsAppWebVersion(previous)
	})

	target := provider.ApplyWhatsAppWebVersion(store.WAVersionContainer{2, 3000, 1037793778})

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
	if !strings.Contains(rec.Body.String(), `"whatsapp_web_version":"`+target.String()+`"`) {
		t.Fatalf("response body = %s, want whatsapp_web_version=%s", rec.Body.String(), target.String())
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
