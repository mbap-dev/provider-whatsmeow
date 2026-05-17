package provider

import "testing"

func TestBuildProxyAddressInfersProtocolAndAuth(t *testing.T) {
	got, err := buildProxyAddress("", "proxy.example.com", "8080", "user", "pass")
	if err != nil {
		t.Fatalf("buildProxyAddress(): %v", err)
	}
	want := "http://user:pass@proxy.example.com:8080"
	if got != want {
		t.Fatalf("buildProxyAddress() = %q, want %q", got, want)
	}
}

func TestBuildProxyAddressInfersSocks5Port(t *testing.T) {
	got, err := buildProxyAddress("", "127.0.0.1", "1080", "", "")
	if err != nil {
		t.Fatalf("buildProxyAddress(): %v", err)
	}
	want := "socks5://127.0.0.1:1080"
	if got != want {
		t.Fatalf("buildProxyAddress() = %q, want %q", got, want)
	}
}

func TestSessionConfigProxyOverride(t *testing.T) {
	manager := NewClientManager(ManagerOptions{
		SessionStore:  t.TempDir(),
		WorkerID:      "worker-test",
		ProxyProtocol: "http",
		ProxyHost:     "global.proxy",
		ProxyPort:     "8080",
		ProxyUsername: "global-user",
		ProxyPassword: "global-pass",
	})

	sessionID := "553121159080"
	snapshot, err := manager.GetSessionConfig(sessionID)
	if err != nil {
		t.Fatalf("GetSessionConfig(): %v", err)
	}
	if !snapshot.Config.ProxyEnabled || snapshot.Config.ProxyHost != "global.proxy" {
		t.Fatalf("default proxy config = %#v, want enabled global proxy", snapshot.Config)
	}

	enabled := true
	protocol := "socks"
	host := "session.proxy"
	port := "1080"
	username := "session-user"
	password := "session-pass"
	snapshot, err = manager.UpdateSessionConfig(sessionID, SessionConfigPatch{
		ProxyEnabled:  &enabled,
		ProxyProtocol: &protocol,
		ProxyHost:     &host,
		ProxyPort:     &port,
		ProxyUsername: &username,
		ProxyPassword: &password,
	})
	if err != nil {
		t.Fatalf("UpdateSessionConfig(): %v", err)
	}
	if !snapshot.Config.ProxyEnabled ||
		snapshot.Config.ProxyProtocol != "socks5" ||
		snapshot.Config.ProxyHost != "session.proxy" ||
		snapshot.Config.ProxyPort != "1080" ||
		snapshot.Config.ProxyUsername != "session-user" ||
		snapshot.Config.ProxyPassword != "session-pass" {
		t.Fatalf("session proxy config = %#v", snapshot.Config)
	}
}

func TestSessionConfigCanDisableGlobalProxy(t *testing.T) {
	manager := NewClientManager(ManagerOptions{
		SessionStore:  t.TempDir(),
		WorkerID:      "worker-test",
		ProxyProtocol: "http",
		ProxyHost:     "global.proxy",
		ProxyPort:     "8080",
	})

	enabled := false
	snapshot, err := manager.UpdateSessionConfig("553121159080", SessionConfigPatch{
		ProxyEnabled: &enabled,
	})
	if err != nil {
		t.Fatalf("UpdateSessionConfig(): %v", err)
	}
	if snapshot.Config.ProxyEnabled {
		t.Fatalf("ProxyEnabled = true, want false")
	}
}
