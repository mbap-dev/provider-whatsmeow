package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all configurable settings for the adapter.  Each field
// corresponds to an environment variable.  Defaults are applied where
// reasonable so the service can run locally with minimal setup.
type Config struct {
	// AMQPURL is the connection string used to connect to the RabbitMQ
	// broker.  Example: amqp://guest:guest@localhost:5672/.
	AMQPURL string
	// AMQPExchange is the name of the exchange to bind to for outgoing
	// messages.  By default this is set to "provider.whatsmeow.outgoing".
	AMQPExchange string
	// AMQPBinding is the routing key pattern used when binding the
	// queue.  It is set to "provider.whatsmeow.*" to receive all
	// messages destined for this provider.
	AMQPBinding string
	// AMQPBindings allows multiple explicit bindings. When set, it
	// takes precedence over AMQPBinding.
	AMQPBindings []string
	// AMQPQueue is the name of the queue that will be declared and
	// consumed by this service.  If empty, a server-generated name is
	// used.
	AMQPQueue string
	// AMQPWebhookExchange is the exchange to which Cloud webhook payloads
	// are published. Default: "provider.whatsmeow.webhooks" (direct).
	AMQPWebhookExchange string
	// AMQPConsumerEnabled toggles AMQP consumption on/off.
	AMQPConsumerEnabled bool
	// WebhookBase is the base URL to which webhook payloads should be
	// delivered.  If left empty the service will not emit webhooks.
	WebhookBase string
	// SessionStore is the directory on disk where session files are
	// persisted.  A separate subdirectory will be created for each
	// session identifier.
	SessionStore string
	// HTTPAddr is the host:port on which to expose the HTTP API and
	// health checks.  The default is ":8080" which listens on all
	// interfaces.
	HTTPAddr string

	// AudioPTTDefault controls whether audio messages are sent as PTT by default.
	// If not set, defaults to true.
	AudioPTTDefault bool

	// RejectCalls enables automatic rejection of incoming calls. Default: true.
	RejectCalls bool
	// RejectCallsMessage is the optional text sent to the caller after rejecting.
	// Supports escaped newlines (\n) which are converted to real newlines.
	RejectCallsMessage string

	// AutoMarkReadOnMessage, when true, marks chats as read upon receiving messages.
	AutoMarkReadOnMessage bool
	// PNResolverURL: optional HTTP endpoint to resolve an MSISDN to the correct WA PN (digits)
	PNResolverURL string

	// RedisURL points to an optional Redis instance used by status/config.
	// Example: redis://user:pass@localhost:6379/0
	RedisURL string
	// StatusPrefix defines the key prefix used for session status in Redis.
	StatusPrefix string
	// SessionConfigPrefix defines the key prefix used to fetch per-session
	// configuration from Redis.
	SessionConfigPrefix string

	// AlwaysOnline, when true, keeps the presence as available by periodically
	// calling SendPresence(PresenceAvailable) while connected.
	AlwaysOnline bool
	// AlwaysOnlineIntervalSeconds controls how often to refresh presence.
	// Default: 60 seconds.
	AlwaysOnlineIntervalSeconds int

	// IgnoreStatusBroadcast drops incoming WhatsApp Status (status@broadcast) messages.
	// Default: true.
	IgnoreStatusBroadcast bool
	// IgnoreNewsletters drops incoming Newsletter messages (server=newsletter).
	// Default: true.
	IgnoreNewsletters bool
	// OutgoingMessagesCoex publishes own/outgoing webhook messages in
	// coexistence format (`field=smb_message_echoes`, `value.message_echoes`),
	// compatible with newer Chatwoot behavior for sent messages.
	OutgoingMessagesCoex bool

	// CheckUserExists, when true, verifies if the destination number
	// is registered on WhatsApp before sending a message. If the
	// check fails due to network or API errors, the adapter falls
	// back to the previous behaviour and still attempts delivery.
	// Default: true.
	CheckUserExists bool

	// WorkerID identifies the current worker instance.
	WorkerID string
	// WorkerNumbers optionally limits this worker to a fixed list of sessions
	// (usually phone numbers, digits only).
	WorkerNumbers []string
	// WorkerStrict, when true, blocks session actions for sessions that are
	// not assigned to this worker.
	WorkerStrict bool
	// WorkerAutoConnectAssigned connects all assigned sessions during startup.
	WorkerAutoConnectAssigned bool

	// UIEnabled controls whether the built-in dashboard and API are exposed.
	UIEnabled bool
	// UITitle customises the dashboard title.
	UITitle string
	// UIAuthEnabled enables HTTP Basic authentication for dashboard routes.
	UIAuthEnabled bool
	// UIAuthUser is the HTTP Basic username for protected dashboard/API routes.
	UIAuthUser string
	// UIAuthPass is the HTTP Basic password for protected dashboard/API routes.
	UIAuthPass string
	// UIAuthRealm customises the HTTP Basic realm shown by browsers.
	UIAuthRealm string

	// APIAuthEnabled enables access-token verification on API endpoints.
	APIAuthEnabled bool
	// APIAuthToken is the bearer token used to authorize API requests.
	APIAuthToken string
	// APIAuthAllowQuery allows access token via query string (access_token),
	// useful for compatibility with Graph API style requests.
	APIAuthAllowQuery bool

	// Media storage mode. Supports: auto, local, s3.
	MediaStorageMode string
	// MediaStoreDir is the local path used to store media blobs when mode=local
	// or when auto mode falls back from S3.
	MediaStoreDir string
	// MediaRegistryFile stores media metadata index.
	MediaRegistryFile string
	// MediaURLTTLSeconds controls generated URL expiry for S3 pre-signed links.
	MediaURLTTLSeconds int

	// S3 media storage configuration (used when mode=s3 or auto+configured).
	S3Bucket         string
	S3Region         string
	S3Endpoint       string
	S3AccessKeyID    string
	S3SecretAccess   string
	S3ForcePathStyle bool
	S3Prefix         string
}

// NewConfig reads configuration from the environment and returns a
// populated Config instance.  Missing variables fall back to sensible
// defaults as documented on the struct fields.
func NewConfig() *Config {
	cfg := &Config{}
	cfg.AMQPURL = getEnv("AMQP_URL", "")
	cfg.AMQPExchange = getEnv("AMQP_EXCHANGE", "provider.whatsmeow.outgoing")
	cfg.AMQPBinding = getEnv("AMQP_BINDING", "provider.whatsmeow.*")
	cfg.AMQPWebhookExchange = getEnv("AMQP_WEBHOOK_EXCHANGE", "provider.whatsmeow.webhooks")
	cfg.AMQPConsumerEnabled = getEnvBool("AMQP_CONSUMER_ENABLED", true)
	cfg.WebhookBase = getEnv("WEBHOOK_BASE", "")
	cfg.SessionStore = getEnv("SESSION_STORE", "./state/whatsmeow")
	cfg.HTTPAddr = getEnv("HTTP_ADDR", ":8080")

	// Default: send audio as PTT unless explicitly disabled
	cfg.AudioPTTDefault = getEnvBool("AUDIO_PTT_DEFAULT", true)
	cfg.PNResolverURL = getEnv("PN_RESOLVER_URL", "")

	// Auto-reject incoming calls + optional reply message
	cfg.RejectCalls = getEnvBool("REJECT_CALLS", true)
	cfg.RejectCallsMessage = getEnv("REJECT_CALLS_MESSAGE", "")
	// Allow escaped newlines in env var
	if strings.Contains(cfg.RejectCallsMessage, "\\n") {
		cfg.RejectCallsMessage = strings.ReplaceAll(cfg.RejectCallsMessage, "\\n", "\n")
	}

	// Mark chat as read automatically on inbound messages (default: false)
	cfg.AutoMarkReadOnMessage = getEnvBool("MARK_READ_ON_MESSAGE", false)

	// Optional Redis for distributed status/config
	cfg.RedisURL = getEnv("REDIS_URL", "")
	cfg.StatusPrefix = getEnv("STATUS_PREFIX", "provider-whatsmeow-status:")
	cfg.SessionConfigPrefix = getEnv("SESSION_CONFIG_PREFIX", "provider-whatsmeow-config:")

	// Presence keep-alive / always-online
	cfg.AlwaysOnline = getEnvBool("ALWAYS_ONLINE", false)
	if v := getEnv("ALWAYS_ONLINE_INTERVAL_SECONDS", ""); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.AlwaysOnlineIntervalSeconds = n
		}
	}
	if cfg.AlwaysOnlineIntervalSeconds == 0 {
		cfg.AlwaysOnlineIntervalSeconds = 60
	}

	// Inbound filters (defaults: ignore both)
	cfg.IgnoreStatusBroadcast = getEnvBool("IGNORE_STATUS_BROADCAST", true)
	cfg.IgnoreNewsletters = getEnvBool("IGNORE_NEWSLETTERS", true)
	cfg.OutgoingMessagesCoex = getEnvBool("OUTGOING_MESSAGES_COEX", false)

	// Outbound behaviour: validate destination on WhatsApp before sending
	cfg.CheckUserExists = getEnvBool("CHECK_USER_EXISTS", true)

	hostname, _ := os.Hostname()
	workerDefault := sanitizeWorkerID(hostname)
	if workerDefault == "" {
		workerDefault = "worker-1"
	}
	cfg.WorkerID = sanitizeWorkerID(getEnv("WORKER_ID", workerDefault))
	if cfg.WorkerID == "" {
		cfg.WorkerID = workerDefault
	}

	cfg.WorkerNumbers = normalizeWorkerSessions(splitCSV(getEnv("WORKER_NUMBERS", "")))
	cfg.WorkerStrict = getEnvBool("WORKER_STRICT_ASSIGNMENT", true)
	cfg.WorkerAutoConnectAssigned = getEnvBool("WORKER_AUTOCONNECT_ASSIGNED", true)

	rawBindings := splitCSV(getEnv("AMQP_BINDINGS", ""))
	if len(rawBindings) > 0 {
		cfg.AMQPBindings = rawBindings
	} else if len(cfg.WorkerNumbers) > 0 {
		prefix := getEnv("AMQP_BINDING_PREFIX", deriveBindingPrefix(cfg.AMQPBinding))
		if prefix == "" {
			prefix = "provider.whatsmeow"
		}
		cfg.AMQPBindings = make([]string, 0, len(cfg.WorkerNumbers))
		for _, sess := range cfg.WorkerNumbers {
			cfg.AMQPBindings = append(cfg.AMQPBindings, prefix+"."+sess)
		}
	}

	defaultQueue := "outgoing.whatsmeow"
	if len(cfg.WorkerNumbers) > 0 {
		defaultQueue += "." + cfg.WorkerID
	}
	cfg.AMQPQueue = getEnv("AMQP_QUEUE", defaultQueue)

	cfg.UIEnabled = getEnvBool("UI_ENABLED", true)
	cfg.UITitle = getEnv("UI_TITLE", "WhatsMeow Control Center")
	cfg.UIAuthUser = getEnv("UI_AUTH_USER", "")
	cfg.UIAuthPass = getEnv("UI_AUTH_PASS", "")
	cfg.UIAuthRealm = getEnv("UI_AUTH_REALM", "WhatsMeow Dashboard")
	uiAuthDefault := strings.TrimSpace(cfg.UIAuthUser) != "" && strings.TrimSpace(cfg.UIAuthPass) != ""
	cfg.UIAuthEnabled = getEnvBool("UI_AUTH_ENABLED", uiAuthDefault)
	if strings.TrimSpace(cfg.UIAuthUser) == "" || strings.TrimSpace(cfg.UIAuthPass) == "" {
		cfg.UIAuthEnabled = false
	}

	cfg.APIAuthToken = firstNonEmpty(
		getEnv("API_AUTH_TOKEN", ""),
		getEnv("UNOAPI_AUTH_TOKEN", ""),
	)
	cfg.APIAuthEnabled = getEnvBool("API_AUTH_ENABLED", strings.TrimSpace(cfg.APIAuthToken) != "")
	cfg.APIAuthAllowQuery = getEnvBool("API_AUTH_ALLOW_QUERY", true)
	if strings.TrimSpace(cfg.APIAuthToken) == "" {
		cfg.APIAuthEnabled = false
	}

	cfg.MediaStorageMode = strings.ToLower(getEnv("MEDIA_STORAGE_MODE", "auto"))
	cfg.MediaStoreDir = getEnv("MEDIA_STORE_DIR", "")
	cfg.MediaRegistryFile = getEnv("MEDIA_REGISTRY_FILE", "")
	cfg.MediaURLTTLSeconds = 60 * 60 * 24 * 3
	if v := getEnv("MEDIA_URL_TTL_SECONDS", ""); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MediaURLTTLSeconds = n
		}
	}

	// S3 env aliases keep compatibility with UnoAPI-style naming.
	cfg.S3Bucket = firstNonEmpty(
		getEnv("S3_BUCKET", ""),
		getEnv("STORAGE_BUCKET_NAME", ""),
	)
	cfg.S3Region = firstNonEmpty(
		getEnv("S3_REGION", ""),
		getEnv("STORAGE_REGION", ""),
	)
	cfg.S3Endpoint = firstNonEmpty(
		getEnv("S3_ENDPOINT", ""),
		getEnv("STORAGE_ENDPOINT", ""),
	)
	cfg.S3AccessKeyID = firstNonEmpty(
		getEnv("S3_ACCESS_KEY_ID", ""),
		getEnv("STORAGE_ACCESS_KEY_ID", ""),
	)
	cfg.S3SecretAccess = firstNonEmpty(
		getEnv("S3_SECRET_ACCESS_KEY", ""),
		getEnv("STORAGE_SECRET_ACCESS_KEY", ""),
	)
	if _, ok := os.LookupEnv("S3_FORCE_PATH_STYLE"); ok {
		cfg.S3ForcePathStyle = getEnvBool("S3_FORCE_PATH_STYLE", false)
	} else {
		cfg.S3ForcePathStyle = getEnvBool("STORAGE_FORCE_PATH_STYLE", false)
	}
	cfg.S3Prefix = firstNonEmpty(
		getEnv("S3_PREFIX", ""),
		getEnv("STORAGE_PREFIX", ""),
	)

	return cfg
}

// getEnv returns the value of the environment variable named by key.  If
// the variable is not present or empty then defaultVal is returned.
func getEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}

func splitCSV(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func sanitizeWorkerID(v string) string {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		}
	}
	return strings.Trim(b.String(), "-_")
}

func normalizeWorkerSessions(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, raw := range in {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		norm := digitsOnly(raw)
		if norm == "" {
			norm = raw
		}
		if _, exists := seen[norm]; exists {
			continue
		}
		seen[norm] = struct{}{}
		out = append(out, norm)
	}
	return out
}

func deriveBindingPrefix(binding string) string {
	binding = strings.TrimSpace(binding)
	if binding == "" {
		return ""
	}
	if i := strings.IndexAny(binding, "*#"); i >= 0 {
		binding = strings.TrimSuffix(binding[:i], ".")
	}
	return strings.TrimSpace(binding)
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
