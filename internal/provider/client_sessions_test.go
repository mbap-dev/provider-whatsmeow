package provider

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"your.org/provider-whatsmeow/internal/status"
)

func TestKnownSessionsIgnoresInternalDirectories(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{"553121159080", "1", "_media", ".cache"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "README.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}

	manager := NewClientManager(ManagerOptions{
		SessionStore: root,
		WorkerID:     "worker-test",
	})

	got := manager.KnownSessions()
	want := []string{"553121159080"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("KnownSessions() = %#v, want %#v", got, want)
	}
}

func TestNewClientManagerIgnoresInvalidAssignedSessions(t *testing.T) {
	root := t.TempDir()
	assignmentFile := filepath.Join(root, "_worker_assignments.json")
	if err := os.WriteFile(assignmentFile, []byte(`{"worker_id":"worker-test","sessions":["1","553121159080"]}`), 0o644); err != nil {
		t.Fatalf("WriteFile(_worker_assignments.json): %v", err)
	}

	manager := NewClientManager(ManagerOptions{
		SessionStore: root,
		WorkerID:     "worker-test",
	})

	got := manager.AssignedSessions()
	want := []string{"553121159080"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AssignedSessions() = %#v, want %#v", got, want)
	}
}

func TestConnectRejectsInvalidSessionID(t *testing.T) {
	manager := NewClientManager(ManagerOptions{
		SessionStore: t.TempDir(),
		WorkerID:     "worker-test",
	})

	if _, err := manager.Connect("1"); !errors.Is(err, ErrInvalidSessionID) {
		t.Fatalf("Connect(1) error = %v, want %v", err, ErrInvalidSessionID)
	}
}

func TestDeleteSessionRemovesPersistedState(t *testing.T) {
	root := t.TempDir()
	status.Init("", "provider-whatsmeow-status-test:")

	manager := NewClientManager(ManagerOptions{
		SessionStore: root,
		WorkerID:     "worker-test",
	})

	sessionID := "553121159080"
	sessionDir := manager.SessionPath(sessionID)
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(sessionDir): %v", err)
	}
	if err := os.WriteFile(filepath.Join(sessionDir, "session.db"), []byte("db"), 0o644); err != nil {
		t.Fatalf("WriteFile(session.db): %v", err)
	}

	if _, err := manager.AssignSession(sessionID); err != nil {
		t.Fatalf("AssignSession(): %v", err)
	}

	rejectCalls := true
	rejectMessage := "nao atendemos ligacoes"
	alwaysOnline := true
	alwaysOnlineIntervalSeconds := 90
	if _, err := manager.UpdateSessionConfig(sessionID, SessionConfigPatch{
		RejectCalls:                 &rejectCalls,
		RejectCallsMessage:          &rejectMessage,
		AlwaysOnline:                &alwaysOnline,
		AlwaysOnlineIntervalSeconds: &alwaysOnlineIntervalSeconds,
	}); err != nil {
		t.Fatalf("UpdateSessionConfig(): %v", err)
	}

	if _, err := manager.UpsertSessionWebhook(sessionID, SessionWebhook{
		ID:                   "main-hook",
		URLAbsolute:          "https://example.com/webhook",
		Enabled:              true,
		AppendSession:        true,
		SendIncomingMessages: true,
		SendOutgoingMessages: true,
		SendUpdateMessages:   true,
		SendGroupMessages:    true,
	}); err != nil {
		t.Fatalf("UpsertSessionWebhook(): %v", err)
	}

	qrMu.Lock()
	qrBySession[sessionID] = &qrState{pngBase64: "abc123"}
	qrMu.Unlock()
	t.Cleanup(func() {
		clearQR(sessionID)
	})

	status.Set(sessionID, "online")

	removed, err := manager.DeleteSession(sessionID)
	if err != nil {
		t.Fatalf("DeleteSession(): %v", err)
	}
	if !removed {
		t.Fatal("DeleteSession() = false, want true")
	}

	if _, err := os.Stat(sessionDir); !os.IsNotExist(err) {
		t.Fatalf("session dir still exists after delete: %v", err)
	}
	if got := manager.AssignedSessions(); len(got) != 0 {
		t.Fatalf("AssignedSessions() = %#v, want empty", got)
	}
	if snapshot, err := manager.GetSessionConfig(sessionID); err != nil {
		t.Fatalf("GetSessionConfig() after delete: %v", err)
	} else if snapshot.HasOverride {
		t.Fatalf("GetSessionConfig().HasOverride = true, want false")
	}
	if hooks := manager.ListSessionWebhooks(sessionID); len(hooks) != 0 {
		t.Fatalf("ListSessionWebhooks() = %#v, want empty", hooks)
	}
	if hasQR(sessionID) {
		t.Fatal("hasQR() = true, want false")
	}
	if got := status.Get(sessionID); got != "" {
		t.Fatalf("status.Get() = %q, want empty", got)
	}

	reloaded := NewClientManager(ManagerOptions{
		SessionStore: root,
		WorkerID:     "worker-test",
	})
	if got := reloaded.AssignedSessions(); len(got) != 0 {
		t.Fatalf("reloaded AssignedSessions() = %#v, want empty", got)
	}
	if snapshot, err := reloaded.GetSessionConfig(sessionID); err != nil {
		t.Fatalf("reloaded GetSessionConfig(): %v", err)
	} else if snapshot.HasOverride {
		t.Fatalf("reloaded GetSessionConfig().HasOverride = true, want false")
	}
	if hooks := reloaded.ListSessionWebhooks(sessionID); len(hooks) != 0 {
		t.Fatalf("reloaded ListSessionWebhooks() = %#v, want empty", hooks)
	}
	if got := reloaded.KnownSessions(); len(got) != 0 {
		t.Fatalf("reloaded KnownSessions() = %#v, want empty", got)
	}
}
