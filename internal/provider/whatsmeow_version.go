package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"google.golang.org/protobuf/proto"
)

const whatsAppWebServiceWorkerURL = "https://web.whatsapp.com/sw.js"

var (
	whatsAppWebVersionCacheMu  sync.Mutex
	whatsAppWebVersionApplyMu  sync.Mutex
	whatsAppWebVersionCached   *store.WAVersionContainer
	whatsAppWebVersionCachedAt time.Time
	whatsAppWebVersionCacheTTL = time.Hour
	// 2.3000.1042462245 is currently advertised by web.whatsapp.com, but is
	// rejected by the websocket with "client outdated (405)" with this schema.
	// Fall back to the known version bundled with the patched tulir/webauthn
	// branch so the WA version matches the protocol payload being sent.
	whatsAppWebVersionRejectedFallbacks = map[uint32]store.WAVersionContainer{
		1042462245: {2, 3000, 1042386815},
	}
	whatsAppWebVersionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`"client_revision":\s*(\d+)`),
		regexp.MustCompile(`\\"client_revision\\":\s*(\d+)`),
		regexp.MustCompile(`client_revision\\?\\"?:[\s]*(\d+)`),
		regexp.MustCompile(`["']client_revision["'][\s]*:[\s]*(\d+)`),
	}
)

// ApplyWhatsAppWebVersion keeps the whatsmeow global version state aligned
// across the user agent, registration hash and device props payloads.
func ApplyWhatsAppWebVersion(version store.WAVersionContainer) store.WAVersionContainer {
	if version.IsZero() {
		return store.GetWAVersion()
	}

	store.SetWAVersion(version)
	if store.BaseClientPayload != nil && store.BaseClientPayload.UserAgent != nil {
		store.BaseClientPayload.UserAgent.AppVersion = version.ProtoAppVersion()
	}
	if store.DeviceProps != nil && store.DeviceProps.Version != nil {
		store.DeviceProps.Version.Primary = proto.Uint32(version[0])
		store.DeviceProps.Version.Secondary = proto.Uint32(version[1])
		store.DeviceProps.Version.Tertiary = proto.Uint32(version[2])
	}

	return version
}

func CurrentWhatsAppWebVersion() string {
	return store.GetWAVersion().String()
}

func parseWhatsAppWebVersionContent(content string) (store.WAVersionContainer, error) {
	for _, pattern := range whatsAppWebVersionPatterns {
		matches := pattern.FindStringSubmatch(content)
		if len(matches) < 2 {
			continue
		}
		clientRevision, err := strconv.ParseUint(matches[1], 10, 32)
		if err != nil {
			continue
		}
		if clientRevision > 0 {
			return store.WAVersionContainer{2, 3000, uint32(clientRevision)}, nil
		}
	}

	preview := content
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}
	return store.WAVersionContainer{}, fmt.Errorf("client_revision not found in WhatsApp Web service worker: %s", preview)
}

func fetchLatestWhatsAppWebVersion(ctx context.Context, httpClient *http.Client) (store.WAVersionContainer, error) {
	return fetchLatestWhatsAppWebVersionCached(ctx, httpClient, false)
}

func fetchLatestWhatsAppWebVersionCached(ctx context.Context, httpClient *http.Client, useCache bool) (store.WAVersionContainer, error) {
	if override, ok, err := whatsAppWebVersionOverride(); ok || err != nil {
		return override, err
	}

	whatsAppWebVersionCacheMu.Lock()
	if useCache && whatsAppWebVersionCached != nil && time.Since(whatsAppWebVersionCachedAt) < whatsAppWebVersionCacheTTL {
		cached := *whatsAppWebVersionCached
		whatsAppWebVersionCacheMu.Unlock()
		return cached, nil
	}
	whatsAppWebVersionCacheMu.Unlock()

	version, err := fetchLatestWhatsAppWebVersionFromHTML(ctx, httpClient)
	if err == nil {
		version = normalizeWhatsAppWebVersion(version)
		whatsAppWebVersionCacheMu.Lock()
		whatsAppWebVersionCached = &version
		whatsAppWebVersionCachedAt = time.Now()
		whatsAppWebVersionCacheMu.Unlock()
		return version, nil
	}

	version, fallbackErr := fetchLatestWhatsAppWebVersionFromServiceWorker(ctx, httpClient)
	if fallbackErr != nil {
		return store.WAVersionContainer{}, fmt.Errorf("fetch WhatsApp Web version from HTML: %w; service worker fallback: %w", err, fallbackErr)
	}
	version = normalizeWhatsAppWebVersion(version)

	whatsAppWebVersionCacheMu.Lock()
	whatsAppWebVersionCached = &version
	whatsAppWebVersionCachedAt = time.Now()
	whatsAppWebVersionCacheMu.Unlock()

	return version, nil
}

func whatsAppWebVersionOverride() (store.WAVersionContainer, bool, error) {
	raw := strings.TrimSpace(os.Getenv("WHATSAPP_WEB_VERSION"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("WA_WEB_VERSION"))
	}
	if raw == "" {
		return store.WAVersionContainer{}, false, nil
	}
	version, err := store.ParseVersion(raw)
	if err != nil {
		return store.WAVersionContainer{}, true, fmt.Errorf("invalid WHATSAPP_WEB_VERSION %q: %w", raw, err)
	}
	return version, true, nil
}

func normalizeWhatsAppWebVersion(version store.WAVersionContainer) store.WAVersionContainer {
	if fallback, rejected := whatsAppWebVersionRejectedFallbacks[version[2]]; rejected {
		return fallback
	}
	return version
}

func fetchLatestWhatsAppWebVersionFromHTML(ctx context.Context, httpClient *http.Client) (store.WAVersionContainer, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{Proxy: proxyFuncFromConfigOrEnvironment},
		}
	}
	latest, err := whatsmeow.GetLatestVersion(ctx, httpClient)
	if err != nil {
		return store.WAVersionContainer{}, err
	}
	if latest == nil || latest.IsZero() {
		return store.WAVersionContainer{}, fmt.Errorf("empty WhatsApp Web version from HTML")
	}
	return *latest, nil
}

func fetchLatestWhatsAppWebVersionFromServiceWorker(ctx context.Context, httpClient *http.Client) (store.WAVersionContainer, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, whatsAppWebServiceWorkerURL, nil)
	if err != nil {
		return store.WAVersionContainer{}, fmt.Errorf("prepare WhatsApp Web version request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")

	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{Proxy: proxyFuncFromConfigOrEnvironment},
		}
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return store.WAVersionContainer{}, fmt.Errorf("fetch WhatsApp Web version: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return store.WAVersionContainer{}, fmt.Errorf("read WhatsApp Web version response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return store.WAVersionContainer{}, fmt.Errorf("unexpected WhatsApp Web version status %d", resp.StatusCode)
	}

	version, err := parseWhatsAppWebVersionContent(string(body))
	if err != nil {
		return store.WAVersionContainer{}, err
	}
	return version, nil
}

// ConfigureLatestWhatsAppWebVersion tries to fetch the latest WhatsApp Web
// version and applies it globally. On failure, the current whatsmeow version
// remains active. The fetch bypasses the cache so session startup/reload can
// pick up the current WhatsApp Web revision.
func ConfigureLatestWhatsAppWebVersion(ctx context.Context) (previous store.WAVersionContainer, applied store.WAVersionContainer, err error) {
	return ConfigureLatestWhatsAppWebVersionWithClient(ctx, nil)
}

func ConfigureLatestWhatsAppWebVersionWithClient(ctx context.Context, httpClient *http.Client) (previous store.WAVersionContainer, applied store.WAVersionContainer, err error) {
	previous = store.GetWAVersion()
	latest, err := fetchLatestWhatsAppWebVersion(ctx, httpClient)
	if err != nil {
		return previous, previous, err
	}
	whatsAppWebVersionApplyMu.Lock()
	defer whatsAppWebVersionApplyMu.Unlock()
	previous = store.GetWAVersion()
	applied = ApplyWhatsAppWebVersion(latest)
	return previous, applied, nil
}

func refreshLatestWhatsAppWebVersionForConnect(ctx context.Context, httpClient *http.Client, connect func() error) (previous store.WAVersionContainer, applied store.WAVersionContainer, refreshErr error, connectErr error) {
	latest, refreshErr := fetchLatestWhatsAppWebVersion(ctx, httpClient)

	whatsAppWebVersionApplyMu.Lock()
	defer whatsAppWebVersionApplyMu.Unlock()

	previous = store.GetWAVersion()
	if refreshErr == nil {
		applied = ApplyWhatsAppWebVersion(latest)
	} else {
		applied = previous
	}
	if connect != nil {
		connectErr = connect()
	}
	return previous, applied, refreshErr, connectErr
}
