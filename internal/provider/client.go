package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	ilog "your.org/provider-whatsmeow/internal/log"
	"your.org/provider-whatsmeow/internal/status"
)

// clientEntry mantém refs para facilitar cleanup
type clientEntry struct {
	Client             *whatsmeow.Client
	DB                 *sqlstore.Container
	KeysDB             *sqlstore.Container // opcional
	Log                waLog.Logger
	Session            string
	Webhook            string
	WhatsAppWebVersion string
	// cancel function for the always-online presence refresher
	presenceCancel context.CancelFunc
}

// ClientManager mantém um registro de clientes por sessão
type ClientManager struct {
	mu              sync.RWMutex
	clients         map[string]*clientEntry
	sessionStore    string
	assignmentFile  string
	sessionCfgFile  string
	sessionHookFile string
	sessionTplFile  string
	webhookBase     string
	defaultAudioPTT bool
	rejectCalls     bool
	rejectMsg       string
	autoMarkRead    bool
	pnResolverURL   string
	proxyDefault    ProxyConfig
	// simple PN->JID cache (E.164 -> JID string)
	pnCache map[string]pnCacheItem

	// always-online settings
	alwaysOnline      bool
	alwaysOnlineEvery time.Duration

	// inbound filters
	ignoreStatusBroadcast bool
	ignoreNewsletters     bool
	// OutgoingMessagesCoex switches outgoing message webhooks to coexistence
	// format (`smb_message_echoes`) for Chatwoot compatibility.
	outgoingMessagesCoex bool

	// redis url for distributed config (optional)
	redisURL string
	// key prefix used to resolve per-session config in Redis.
	sessionConfigPrefix string

	// per-session overrides loaded from redis config
	overridesMu sync.RWMutex
	overrides   map[string]sessionOverrides

	// registered webhook endpoints by session
	webhooksMu      sync.RWMutex
	sessionWebhooks map[string][]SessionWebhook

	// registered Cloud API-compatible message templates by session
	templatesMu      sync.RWMutex
	sessionTemplates map[string][]SessionTemplate

	// checkUserExists controls whether we should verify that a
	// destination number is registered on WhatsApp before sending.
	checkUserExists bool

	// worker identity and assignment controls
	workerID     string
	workerStrict bool
	assigned     map[string]struct{}

	// media storage and metadata registry (local/S3)
	mediaStore *mediaStore
}

type sessionOverrides struct {
	// if nil, use manager default
	rejectCalls           *bool
	rejectMsg             *string
	autoMarkRead          *bool
	proxyEnabled          *bool
	proxyProtocol         *string
	proxyHost             *string
	proxyPort             *string
	proxyUsername         *string
	proxyPassword         *string
	ignoreStatusBroadcast *bool
	ignoreNewsletters     *bool
	outgoingMessagesCoex  *bool
	alwaysOnline          *bool
	alwaysOnlineEvery     *time.Duration
	historySyncEnabled    *bool
	historySyncDays       *int
}

type ManagerOptions struct {
	SessionStore          string
	WebhookBase           string
	DefaultAudioPTT       bool
	RejectCalls           bool
	RejectMessage         string
	AutoMarkRead          bool
	PNResolverURL         string
	ProxyProtocol         string
	ProxyHost             string
	ProxyPort             string
	ProxyUsername         string
	ProxyPassword         string
	IgnoreStatusBroadcast bool
	IgnoreNewsletters     bool
	OutgoingMessagesCoex  bool
	RedisURL              string
	SessionConfigPrefix   string
	CheckUserExists       bool
	WorkerID              string
	AssignedSessions      []string
	WorkerStrict          bool
	AssignmentFile        string
	SessionConfigFile     string
	SessionWebhookFile    string
	SessionTemplateFile   string
	MediaStoreDir         string
	MediaRegistryFile     string
	MediaStorageMode      string
	MediaURLTTLSeconds    int
	S3Bucket              string
	S3Region              string
	S3Endpoint            string
	S3AccessKeyID         string
	S3SecretAccessKey     string
	S3ForcePathStyle      bool
	S3Prefix              string
}

var (
	ErrMissingSessionID  = errors.New("missing session id")
	ErrInvalidSessionID  = errors.New("invalid session id")
	ErrReservedSessionID = errors.New("reserved session id")
	ErrMissingPhone      = errors.New("missing phone")
	ErrAlreadyPaired     = errors.New("session already paired")
	ErrPairingNotReady   = errors.New("pairing setup not ready yet")
)

func NewClientManager(opts ManagerOptions) *ClientManager {
	assignmentFile := strings.TrimSpace(opts.AssignmentFile)
	if assignmentFile == "" {
		assignmentFile = filepath.Join(opts.SessionStore, "_worker_assignments.json")
	}
	sessionCfgFile := strings.TrimSpace(opts.SessionConfigFile)
	if sessionCfgFile == "" {
		sessionCfgFile = filepath.Join(opts.SessionStore, "_session_overrides.json")
	}
	sessionHookFile := strings.TrimSpace(opts.SessionWebhookFile)
	if sessionHookFile == "" {
		sessionHookFile = filepath.Join(opts.SessionStore, "_session_webhooks.json")
	}
	sessionTplFile := strings.TrimSpace(opts.SessionTemplateFile)
	if sessionTplFile == "" {
		sessionTplFile = filepath.Join(opts.SessionStore, "_session_templates.json")
	}
	assigned := map[string]struct{}{}
	for _, sid := range opts.AssignedSessions {
		n := normalizeSessionID(sid)
		if n != "" {
			assigned[n] = struct{}{}
		}
	}
	if persisted, err := loadAssignments(assignmentFile, opts.WorkerID); err == nil {
		for _, sid := range persisted {
			n := normalizeSessionID(sid)
			if n != "" {
				assigned[n] = struct{}{}
			}
		}
	}
	overrides := map[string]sessionOverrides{}
	if persisted, err := loadSessionOverrideState(sessionCfgFile, opts.WorkerID); err == nil && persisted != nil {
		overrides = persisted
	}
	sessionWebhooks := map[string][]SessionWebhook{}
	if persisted, err := loadSessionWebhookState(sessionHookFile, opts.WorkerID); err == nil && persisted != nil {
		sessionWebhooks = persisted
	}
	sessionTemplates := map[string][]SessionTemplate{}
	if persisted, err := loadSessionTemplateState(sessionTplFile, opts.WorkerID); err == nil && persisted != nil {
		sessionTemplates = persisted
	}
	mediaDir := strings.TrimSpace(opts.MediaStoreDir)
	if mediaDir == "" {
		mediaDir = filepath.Join(opts.SessionStore, "_media")
	}
	mediaRegistry := strings.TrimSpace(opts.MediaRegistryFile)
	if mediaRegistry == "" {
		mediaRegistry = filepath.Join(opts.SessionStore, "_media_registry.json")
	}
	media := newMediaStore(MediaStoreOptions{
		Mode:             opts.MediaStorageMode,
		LocalDir:         mediaDir,
		RegistryFile:     mediaRegistry,
		URLTTLSeconds:    opts.MediaURLTTLSeconds,
		S3Bucket:         opts.S3Bucket,
		S3Region:         opts.S3Region,
		S3Endpoint:       opts.S3Endpoint,
		S3AccessKeyID:    opts.S3AccessKeyID,
		S3SecretAccess:   opts.S3SecretAccessKey,
		S3ForcePathStyle: opts.S3ForcePathStyle,
		S3Prefix:         opts.S3Prefix,
	})
	return &ClientManager{
		clients:         make(map[string]*clientEntry),
		sessionStore:    opts.SessionStore,
		assignmentFile:  assignmentFile,
		sessionCfgFile:  sessionCfgFile,
		sessionHookFile: sessionHookFile,
		sessionTplFile:  sessionTplFile,
		webhookBase:     opts.WebhookBase,
		defaultAudioPTT: opts.DefaultAudioPTT,
		rejectCalls:     opts.RejectCalls,
		rejectMsg:       opts.RejectMessage,
		autoMarkRead:    opts.AutoMarkRead,
		pnResolverURL:   opts.PNResolverURL,
		proxyDefault: normalizeProxyConfig(ProxyConfig{
			Protocol: opts.ProxyProtocol,
			Host:     opts.ProxyHost,
			Port:     opts.ProxyPort,
			Username: opts.ProxyUsername,
			Password: opts.ProxyPassword,
			Enabled:  strings.TrimSpace(opts.ProxyHost) != "" && strings.TrimSpace(opts.ProxyPort) != "",
		}),
		pnCache:               make(map[string]pnCacheItem),
		alwaysOnlineEvery:     time.Minute,
		ignoreStatusBroadcast: opts.IgnoreStatusBroadcast,
		ignoreNewsletters:     opts.IgnoreNewsletters,
		outgoingMessagesCoex:  opts.OutgoingMessagesCoex,
		redisURL:              opts.RedisURL,
		sessionConfigPrefix:   opts.SessionConfigPrefix,
		overrides:             overrides,
		sessionWebhooks:       sessionWebhooks,
		sessionTemplates:      sessionTemplates,
		checkUserExists:       opts.CheckUserExists,
		workerID:              strings.TrimSpace(opts.WorkerID),
		workerStrict:          opts.WorkerStrict,
		assigned:              assigned,
		mediaStore:            media,
	}
}

func normalizeSessionID(sessionID string) string {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return ""
	}
	if isInternalSessionEntry(sessionID) || sessionID == "." || sessionID == ".." {
		return ""
	}
	if strings.Contains(sessionID, "/") || strings.Contains(sessionID, "\\") {
		return ""
	}
	if looksLikePhoneIdentifier(sessionID) {
		if d := digitsOnly(sessionID); len(d) >= 6 {
			return d
		}
		return ""
	}
	if len(sessionID) < 2 {
		return ""
	}
	return sessionID
}

func isInternalSessionEntry(name string) bool {
	name = strings.TrimSpace(name)
	return name == "" || strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".")
}

func looksLikePhoneIdentifier(v string) bool {
	hasDigit := false
	for _, r := range v {
		switch {
		case r >= '0' && r <= '9':
			hasDigit = true
		case r == '+' || r == ' ' || r == '-' || r == '(' || r == ')':
			// keep scanning
		default:
			return false
		}
	}
	return hasDigit
}

func (m *ClientManager) AllowsSession(sessionID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.assigned) == 0 {
		return true
	}
	norm := normalizeSessionID(sessionID)
	return m.isAssignedSession(norm) || !m.workerStrict
}

func (m *ClientManager) isAssignedSession(normalizedSessionID string) bool {
	_, ok := m.assigned[normalizedSessionID]
	return ok
}

func (m *ClientManager) IsAssignedSession(sessionID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isAssignedSession(normalizeSessionID(sessionID))
}

func (m *ClientManager) WorkerID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.workerID
}

func (m *ClientManager) WorkerStrict() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.workerStrict
}

func (m *ClientManager) CanonicalSessionID(sessionID string) (string, error) {
	norm := normalizeSessionID(sessionID)
	if norm == "" {
		return "", ErrInvalidSessionID
	}
	return norm, nil
}

func (m *ClientManager) AssignedSessions() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.assigned) == 0 {
		return nil
	}
	out := make([]string, 0, len(m.assigned))
	for sid := range m.assigned {
		out = append(out, sid)
	}
	sort.Strings(out)
	return out
}

func (m *ClientManager) AssignSession(sessionID string) (string, error) {
	norm := normalizeSessionID(sessionID)
	if norm == "" {
		return "", ErrInvalidSessionID
	}

	m.mu.Lock()
	if m.assigned == nil {
		m.assigned = map[string]struct{}{}
	}
	_, existed := m.assigned[norm]
	if !existed {
		m.assigned[norm] = struct{}{}
	}
	sessions := make([]string, 0, len(m.assigned))
	for sid := range m.assigned {
		sessions = append(sessions, sid)
	}
	sort.Strings(sessions)
	file := m.assignmentFile
	workerID := m.workerID
	m.mu.Unlock()

	if err := saveAssignments(file, workerID, sessions); err != nil {
		if !existed {
			m.mu.Lock()
			delete(m.assigned, norm)
			m.mu.Unlock()
		}
		return "", err
	}
	return norm, nil
}

func (m *ClientManager) UnassignSession(sessionID string) (bool, error) {
	norm := normalizeSessionID(sessionID)
	if norm == "" {
		return false, ErrInvalidSessionID
	}

	m.mu.Lock()
	if len(m.assigned) == 0 {
		m.mu.Unlock()
		return false, nil
	}
	_, existed := m.assigned[norm]
	if !existed {
		m.mu.Unlock()
		return false, nil
	}
	delete(m.assigned, norm)
	sessions := make([]string, 0, len(m.assigned))
	for sid := range m.assigned {
		sessions = append(sessions, sid)
	}
	sort.Strings(sessions)
	file := m.assignmentFile
	workerID := m.workerID
	m.mu.Unlock()

	if err := saveAssignments(file, workerID, sessions); err != nil {
		m.mu.Lock()
		m.assigned[norm] = struct{}{}
		m.mu.Unlock()
		return false, err
	}
	return true, nil
}

func (m *ClientManager) AutoConnectAssignedSessions() []string {
	assigned := m.AssignedSessions()
	connected := make([]string, 0, len(assigned))
	for _, sid := range assigned {
		sid := sid
		go func() {
			_, _ = m.Connect(sid)
		}()
		connected = append(connected, sid)
	}
	return connected
}

type assignmentState struct {
	WorkerID string   `json:"worker_id"`
	Sessions []string `json:"sessions"`
}

func loadAssignments(path, workerID string) ([]string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	st := assignmentState{}
	if err := jsonUnmarshal(b, &st); err != nil {
		return nil, err
	}
	if strings.TrimSpace(st.WorkerID) != "" && strings.TrimSpace(workerID) != "" && strings.TrimSpace(st.WorkerID) != strings.TrimSpace(workerID) {
		return nil, nil
	}
	return st.Sessions, nil
}

func saveAssignments(path, workerID string, sessions []string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	payload := assignmentState{
		WorkerID: strings.TrimSpace(workerID),
		Sessions: sessions,
	}
	b, err := jsonMarshalIndent(payload)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// EnableAlwaysOnline configures the manager to periodically send presence available
// while connected. Interval must be > 0.
func (m *ClientManager) EnableAlwaysOnline(interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}
	m.alwaysOnline = true
	m.alwaysOnlineEvery = interval
}

func (m *ClientManager) maybeStartAlwaysOnline(sessionID string, cli *whatsmeow.Client) {
	if cli == nil {
		return
	}
	// Resolve whether this session should run always-online.
	should, interval := m.getAlwaysOnlineConfig(sessionID)
	if !should {
		return
	}
	m.mu.Lock()
	ent, ok := m.clients[sessionID]
	if !ok || ent == nil || ent.Client == nil {
		m.mu.Unlock()
		return
	}
	// stop previous if any
	if ent.presenceCancel != nil {
		ent.presenceCancel()
		ent.presenceCancel = nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	ent.presenceCancel = cancel
	m.clients[sessionID] = ent
	// interval resolved above
	log := ent.Log
	m.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		// send immediately once
		if err := cli.SendPresence(ctx, types.PresenceAvailable); err != nil {
			log.Errorf("presence available error: %v", err)
		} else {
			log.Infof("presence set to available, interval=%s", interval)
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if !cli.IsConnected() {
					continue
				}
				if err := cli.SendPresence(ctx, types.PresenceAvailable); err != nil {
					log.Errorf("presence refresh error: %v", err)
				}
			}
		}
	}()
}

func (m *ClientManager) stopAlwaysOnline(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ent, ok := m.clients[sessionID]; ok && ent != nil && ent.presenceCancel != nil {
		ent.presenceCancel()
		ent.presenceCancel = nil
	}
}

// ===== Effective config getters (consider per-session overrides) =====

func (m *ClientManager) getRejectCalls(sessionID string) bool {
	m.overridesMu.RLock()
	ov, ok := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	if ok && ov.rejectCalls != nil {
		return *ov.rejectCalls
	}
	return m.rejectCalls
}

func (m *ClientManager) getRejectMsg(sessionID string) string {
	m.overridesMu.RLock()
	ov, ok := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	if ok && ov.rejectMsg != nil {
		return *ov.rejectMsg
	}
	return m.rejectMsg
}

func (m *ClientManager) getAutoMarkRead(sessionID string) bool {
	m.overridesMu.RLock()
	ov, ok := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	if ok && ov.autoMarkRead != nil {
		return *ov.autoMarkRead
	}
	return m.autoMarkRead
}

func (m *ClientManager) getIgnoreStatusBroadcast(sessionID string) bool {
	m.overridesMu.RLock()
	ov, ok := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	if ok && ov.ignoreStatusBroadcast != nil {
		return *ov.ignoreStatusBroadcast
	}
	return m.ignoreStatusBroadcast
}

func (m *ClientManager) getIgnoreNewsletters(sessionID string) bool {
	m.overridesMu.RLock()
	ov, ok := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	if ok && ov.ignoreNewsletters != nil {
		return *ov.ignoreNewsletters
	}
	return m.ignoreNewsletters
}

func (m *ClientManager) getOutgoingMessagesCoex(sessionID string) bool {
	m.overridesMu.RLock()
	ov, ok := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	if ok && ov.outgoingMessagesCoex != nil {
		return *ov.outgoingMessagesCoex
	}
	return m.outgoingMessagesCoex
}

func (m *ClientManager) getCheckUserExists(sessionID string) bool {
	// For now there is only a global setting. Session‑level overrides
	// can be added later if Uno config exposes this flag.
	return m.checkUserExists
}

func (m *ClientManager) getProxyConfig(sessionID string) ProxyConfig {
	cfg := m.proxyDefault
	m.overridesMu.RLock()
	if ov, ok := m.overrides[sessionID]; ok {
		if ov.proxyEnabled != nil {
			cfg.Enabled = *ov.proxyEnabled
		}
		if ov.proxyProtocol != nil {
			cfg.Protocol = *ov.proxyProtocol
		}
		if ov.proxyHost != nil {
			cfg.Host = *ov.proxyHost
		}
		if ov.proxyPort != nil {
			cfg.Port = *ov.proxyPort
		}
		if ov.proxyUsername != nil {
			cfg.Username = *ov.proxyUsername
		}
		if ov.proxyPassword != nil {
			cfg.Password = *ov.proxyPassword
		}
	}
	m.overridesMu.RUnlock()
	return normalizeProxyConfig(cfg)
}

func (m *ClientManager) applyProxyConfig(sessionID string, cli *whatsmeow.Client, logger waLog.Logger) {
	if cli == nil {
		return
	}
	cfg := m.getProxyConfig(sessionID)
	if !cfg.Enabled {
		ilog.WithSession(sessionID).Debug(
			"proxy_client_config disabled enabled=%t host_set=%t port_set=%t protocol=%s",
			cfg.Enabled,
			cfg.Host != "",
			cfg.Port != "",
			cfg.Protocol,
		)
		return
	}
	addr, err := buildProxyAddress(cfg.Protocol, cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	if err != nil {
		ilog.WithSession(sessionID).Debug(
			"proxy_client_config build_error protocol=%s host=%s port=%s username_set=%t password_set=%t err=%v",
			normalizeProxyProtocol(cfg.Protocol, cfg.Port),
			cfg.Host,
			cfg.Port,
			strings.TrimSpace(cfg.Username) != "",
			cfg.Password != "",
			err,
		)
		if logger != nil {
			logger.Warnf("proxy config ignored: %v", err)
		}
		return
	}

	protocol := normalizeProxyProtocol(cfg.Protocol, cfg.Port)
	if protocol == "http" || protocol == "https" {
		proxyURL, err := url.Parse(addr)
		if err != nil {
			ilog.WithSession(sessionID).Debug(
				"proxy_client_config parse_error proxy=%s err=%v",
				redactedProxyAddress(cfg),
				err,
			)
			if logger != nil {
				logger.Warnf("proxy config ignored: %v", err)
			}
			return
		}
		cli.SetProxy(func(req *http.Request) (*url.URL, error) {
			ilog.WithSession(sessionID).Debug(
				"proxy_whatsapp_route enabled target=%s path=%s proxy=%s username_set=%t password_set=%t",
				req.URL.Host,
				req.URL.Path,
				redactedProxyAddress(cfg),
				strings.TrimSpace(cfg.Username) != "",
				cfg.Password != "",
			)
			return proxyURL, nil
		})
	} else {
		if err := cli.SetProxyAddress(addr); err != nil {
			ilog.WithSession(sessionID).Debug(
				"proxy_client_config apply_error proxy=%s username_set=%t password_set=%t err=%v",
				redactedProxyAddress(cfg),
				strings.TrimSpace(cfg.Username) != "",
				cfg.Password != "",
				err,
			)
			if logger != nil {
				logger.Warnf("proxy config ignored: %v", err)
			}
			return
		}
		ilog.WithSession(sessionID).Debug(
			"proxy_whatsapp_route socks5_applied proxy=%s covers=websocket,media_upload,media_download username_set=%t password_set=%t",
			redactedProxyAddress(cfg),
			strings.TrimSpace(cfg.Username) != "",
			cfg.Password != "",
		)
	}
	ilog.WithSession(sessionID).Info("proxy enabled protocol=%s host=%s port=%s covers=websocket,media_upload,media_download", protocol, cfg.Host, cfg.Port)
	ilog.WithSession(sessionID).Debug(
		"proxy_client_config applied proxy=%s username_set=%t password_set=%t",
		redactedProxyAddress(cfg),
		strings.TrimSpace(cfg.Username) != "",
		cfg.Password != "",
	)
}

// loadSessionOverrides fetches session config for the given session from Redis
// (if configured) and stores per-session overrides only for flags that were not
// explicitly set via environment variables.
func (m *ClientManager) loadSessionOverrides(sessionID string) {
	if strings.TrimSpace(m.redisURL) == "" {
		return
	}
	phone := digitsOnly(sessionID)
	if phone == "" {
		return
	}
	cfg, ok := fetchSessionConfig(m.redisURL, m.sessionConfigPrefix, phone)
	if !ok {
		return
	}
	// Prepare override record
	m.overridesMu.Lock()
	ov := m.overrides[sessionID]

	// Env presence checks: only override when env var not explicitly set
	// REJECT_CALLS: only override enable/disable when env not set
	if _, set := os.LookupEnv("REJECT_CALLS"); !set && ov.rejectCalls == nil {
		val := strings.TrimSpace(cfg.rejectCalls)
		b := val != ""
		ov.rejectCalls = &b
	}
	// REJECT_CALLS_MESSAGE: allow message from Redis if env message not set,
	// regardless of REJECT_CALLS env state
	if _, setMsg := os.LookupEnv("REJECT_CALLS_MESSAGE"); !setMsg && ov.rejectMsg == nil {
		val := strings.TrimSpace(cfg.rejectCalls)
		if val != "" {
			ov.rejectMsg = &val
		}
	}
	if _, set := os.LookupEnv("MARK_READ_ON_MESSAGE"); !set && ov.autoMarkRead == nil {
		if v, okb := cfg.readOnReceipt, cfg.hasReadOnReceipt; okb {
			ov.autoMarkRead = &v
		}
	}
	if _, set := os.LookupEnv("IGNORE_STATUS_BROADCAST"); !set && ov.ignoreStatusBroadcast == nil {
		if v, okb := cfg.ignoreBroadcastStatuses, cfg.hasIgnoreBroadcastStatuses; okb {
			ov.ignoreStatusBroadcast = &v
		}
	}
	if _, set := os.LookupEnv("IGNORE_NEWSLETTERS"); !set && ov.ignoreNewsletters == nil {
		if v, okb := cfg.ignoreNewsletterMessages, cfg.hasIgnoreNewsletterMessages; okb {
			ov.ignoreNewsletters = &v
		}
	}
	if _, set := os.LookupEnv("OUTGOING_MESSAGES_COEX"); !set && ov.outgoingMessagesCoex == nil {
		if v, okb := cfg.outgoingMessagesCoex, cfg.hasOutgoingMessagesCoex; okb {
			ov.outgoingMessagesCoex = &v
		}
	}
	if _, set := os.LookupEnv("ALWAYS_ONLINE"); !set {
		if v, okb := cfg.alwaysOnline, cfg.hasAlwaysOnline; okb && ov.alwaysOnline == nil {
			ov.alwaysOnline = &v
		}
		if v, okb := cfg.alwaysOnlineEverySeconds, cfg.hasAlwaysOnlineEverySeconds; okb && v > 0 && ov.alwaysOnlineEvery == nil {
			d := time.Duration(v) * time.Second
			ov.alwaysOnlineEvery = &d
		}
	}
	if _, set := os.LookupEnv("PROXY_HOST"); !set && ov.proxyHost == nil {
		if v := strings.TrimSpace(cfg.proxyHost); v != "" {
			ov.proxyHost = &v
		}
	}
	if _, set := os.LookupEnv("PROXY_PORT"); !set && ov.proxyPort == nil {
		if v := strings.TrimSpace(cfg.proxyPort); v != "" {
			ov.proxyPort = &v
		}
	}
	if _, set := os.LookupEnv("PROXY_PROTOCOL"); !set && ov.proxyProtocol == nil {
		if v := strings.TrimSpace(cfg.proxyProtocol); v != "" {
			ov.proxyProtocol = &v
		}
	}
	if _, set := os.LookupEnv("PROXY_USERNAME"); !set && ov.proxyUsername == nil {
		if v := strings.TrimSpace(cfg.proxyUsername); v != "" {
			ov.proxyUsername = &v
		}
	}
	if _, set := os.LookupEnv("PROXY_PASSWORD"); !set && ov.proxyPassword == nil {
		if cfg.proxyPassword != "" {
			v := cfg.proxyPassword
			ov.proxyPassword = &v
		}
	}
	if ov.proxyEnabled == nil && cfg.hasProxyEnabled {
		v := cfg.proxyEnabled
		ov.proxyEnabled = &v
	}
	m.overrides[sessionID] = ov
	m.overridesMu.Unlock()
}

// Connect cria/recupera cliente e inicia conexão (emite QR se for 1o login)
func (m *ClientManager) Connect(sessionID string) (*whatsmeow.Client, error) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return nil, ErrInvalidSessionID
	}
	if !m.AllowsSession(sessionID) {
		return nil, fmt.Errorf("session %s is not assigned to worker %s", sessionID, m.workerID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if ent, ok := m.clients[sessionID]; ok && ent.Client != nil {
		return ent.Client, nil
	}

	// Mark status as connecting
	status.Set(sessionID, "connecting")

	// Garante diretório por sessão
	if err := os.MkdirAll(m.SessionPath(sessionID), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir session dir: %w", err)
	}

	ctx := context.Background()

	// SQLite por sessão: SESSION_STORE/<sessionID>/session.db  (modernc)
	dbPath := filepath.Join(m.SessionPath(sessionID), "session.db")

	waLogLevel := strings.ToUpper(strings.TrimSpace(os.Getenv("WHATSMEOW_LOG_LEVEL")))
	if waLogLevel == "" {
		waLogLevel = "INFO"
	}
	dbLog := waLog.Stdout("Database", waLogLevel, true)
	clientLog := waLog.Stdout("Client", waLogLevel, true)

	// PRAGMAs para reduzir SQLITE_BUSY e melhorar concorrência
	dsn := fmt.Sprintf(
		"file:%s?cache=shared&_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)",
		dbPath,
	)
	storeContainer, err := sqlstore.New(ctx, "sqlite", dsn, dbLog)
	if err != nil {
		return nil, fmt.Errorf("sqlstore.New: %w", err)
	}

	// Pega primeiro device (cria se não existir) ou cria manualmente
	device, err := storeContainer.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("get first device: %w", err)
	}
	if device == nil {
		// NewDevice NÃO tem argumentos e NÃO retorna erro
		device = storeContainer.NewDevice()
	}

	// (Opcional) ajustar props do device; deixamos simples
	_ = store.DeviceProps

	// (Opcional) KeysDB separado para cache crypto (mantido comentado)
	var keysContainer *sqlstore.Container = nil
	// Exemplo:
	// keysPath := filepath.Join(m.SessionPath(sessionID), "keys.db")
	// keysDSN := fmt.Sprintf(
	//     "file:%s?cache=shared&_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)",
	//     keysPath,
	// )
	// if kc, err := sqlstore.New(ctx, "sqlite", keysDSN, dbLog); err == nil && device.ID != nil {
	//     inner := sqlstore.NewSQLStore(kc, *device.ID)
	//     device.Identities    = inner
	//     device.Sessions      = inner
	//     device.PreKeys       = inner
	//     device.SenderKeys    = inner
	//     device.MsgSecrets    = inner
	//     device.PrivacyTokens = inner
	//     keysContainer = kc
	// }

	cli := whatsmeow.NewClient(device, clientLog)
	cli.EnableAutoReconnect = true
	cli.AutoTrustIdentity = true

	// Load per-session overrides before any connect/QR flow so proxy settings
	// are already attached to websocket/media clients.
	m.loadSessionOverrides(sessionID)
	m.logSessionConfig(sessionID)
	m.applyProxyConfig(sessionID, cli, clientLog)

	// Se já existe ID no store, a sessão já foi pareada antes -> conecta direto
	// Caso contrário, liga o canal de QR ANTES de conectar.
	if cli.Store != nil && cli.Store.ID != nil {
		// handlers (eventos → webhook)
		m.registerEventHandlers(cli, sessionID)

		entry := &clientEntry{
			Client:  cli,
			DB:      storeContainer,
			KeysDB:  keysContainer,
			Log:     clientLog,
			Session: sessionID,
			Webhook: m.webhookBase,
		}
		m.clients[sessionID] = entry

		// conecta (assíncrono)
		go m.connectClientWithLatestWhatsAppWebVersion(sessionID, cli, clientLog)
		return cli, nil
	}

	// Ainda não pareado: abre canal de QR ANTES do Connect e inicia watcher
	qrCh, err := cli.GetQRChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("get QR channel: %w", err)
	}
	m.startQRWatcher(sessionID, qrCh)

	// handlers (eventos → webhook)
	m.registerEventHandlers(cli, sessionID)

	entry := &clientEntry{
		Client:  cli,
		DB:      storeContainer,
		KeysDB:  keysContainer,
		Log:     clientLog,
		Session: sessionID,
		Webhook: m.webhookBase,
	}
	m.clients[sessionID] = entry

	// Conecta (assíncrono)
	go m.connectClientWithLatestWhatsAppWebVersion(sessionID, cli, clientLog)

	return cli, nil
}

func (m *ClientManager) connectClientWithLatestWhatsAppWebVersion(sessionID string, cli *whatsmeow.Client, clientLog waLog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	proxyCfg := m.getProxyConfig(sessionID)
	ilog.WithSession(sessionID).Debug(
		"wa_web_version_refresh proxy_cfg enabled=%t proxy=%s username_set=%t password_set=%t",
		proxyCfg.Enabled,
		redactedProxyAddress(proxyCfg),
		strings.TrimSpace(proxyCfg.Username) != "",
		proxyCfg.Password != "",
	)
	httpClient := httpClientForProxyConfig(sessionID, proxyCfg)
	previous, applied, refreshErr, connectErr := refreshLatestWhatsAppWebVersionForConnect(ctx, httpClient, cli.Connect)
	if refreshErr != nil {
		ilog.WithSession(sessionID).Warn(
			"failed to refresh WhatsApp Web version for session connect, keeping %s: %v",
			applied.String(),
			refreshErr,
		)
	} else if previous != applied {
		ilog.WithSession(sessionID).Info(
			"updated WhatsApp Web version for session connect from %s to %s",
			previous.String(),
			applied.String(),
		)
	} else {
		ilog.WithSession(sessionID).Info("using WhatsApp Web version %s for session connect", applied.String())
	}

	m.mu.Lock()
	if ent := m.clients[sessionID]; ent != nil {
		ent.WhatsAppWebVersion = applied.String()
	}
	m.mu.Unlock()

	if connectErr != nil {
		clientLog.Errorf("connect error: %v", connectErr)
		status.Set(sessionID, "offline")
	}
}

func (m *ClientManager) logSessionConfig(sessionID string) {
	entry := ilog.WithSession(sessionID)
	rc := m.getRejectCalls(sessionID)
	msg := strings.TrimSpace(m.getRejectMsg(sessionID))
	amr := m.getAutoMarkRead(sessionID)
	isb := m.getIgnoreStatusBroadcast(sessionID)
	inl := m.getIgnoreNewsletters(sessionID)
	coex := m.getOutgoingMessagesCoex(sessionID)
	proxyCfg := m.getProxyConfig(sessionID)
	// Resolve always-online effective
	should, interval := m.getAlwaysOnlineConfig(sessionID)
	m.overridesMu.RLock()
	ov := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	// Sources: env set?
	envRC, _ := os.LookupEnv("REJECT_CALLS")
	envRCM, _ := os.LookupEnv("REJECT_CALLS_MESSAGE")
	envAMR, _ := os.LookupEnv("MARK_READ_ON_MESSAGE")
	envISB, _ := os.LookupEnv("IGNORE_STATUS_BROADCAST")
	envINL, _ := os.LookupEnv("IGNORE_NEWSLETTERS")
	envCOEX, _ := os.LookupEnv("OUTGOING_MESSAGES_COEX")
	envAO, _ := os.LookupEnv("ALWAYS_ONLINE")
	envProxyHost, _ := os.LookupEnv("PROXY_HOST")
	entry.Info(
		"session_config reject_calls=%t source_reject_calls=%s reject_msg_set=%t source_reject_msg=%s mark_read_on_message=%t source_mark_read=%s ignore_status_broadcast=%t source_isb=%s ignore_newsletters=%t source_inl=%s outgoing_messages_coex=%t source_coex=%s always_online=%t source_always_online=%s always_online_interval_s=%d proxy_enabled=%t proxy_protocol=%s proxy_host_set=%t source_proxy=%s",
		rc,
		source(envRC, ov.rejectCalls != nil),
		msg != "",
		source(envRCM, ov.rejectMsg != nil),
		amr,
		source(envAMR, ov.autoMarkRead != nil),
		isb,
		source(envISB, ov.ignoreStatusBroadcast != nil),
		inl,
		source(envINL, ov.ignoreNewsletters != nil),
		coex,
		source(envCOEX, ov.outgoingMessagesCoex != nil),
		should,
		source(envAO, ov.alwaysOnline != nil || ov.alwaysOnlineEvery != nil),
		int(interval/time.Second),
		proxyCfg.Enabled,
		proxyCfg.Protocol,
		strings.TrimSpace(proxyCfg.Host) != "",
		source(envProxyHost, ov.proxyEnabled != nil || ov.proxyHost != nil || ov.proxyPort != nil || ov.proxyProtocol != nil),
	)
}

func source(envVal string, hasOverride bool) string {
	if strings.TrimSpace(envVal) != "" {
		return "env"
	}
	if hasOverride {
		return "session_override"
	}
	return "default"
}

func (m *ClientManager) Disconnect(sessionID string) error {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return ErrInvalidSessionID
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if ent, ok := m.clients[sessionID]; ok {
		if ent.presenceCancel != nil {
			ent.presenceCancel()
			ent.presenceCancel = nil
		}
		if ent.Client != nil {
			ent.Client.Disconnect()
		}
		delete(m.clients, sessionID)
	}
	// Mark as disconnected
	status.Set(sessionID, "disconnected")
	return nil
}

func (m *ClientManager) Reload(sessionID string) error {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return ErrInvalidSessionID
	}
	m.mu.RLock()
	_, ok := m.clients[sessionID]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("session %s not found", sessionID)
	}
	status.Set(sessionID, "connecting")
	if err := m.Disconnect(sessionID); err != nil {
		return err
	}
	_, err := m.Connect(sessionID)
	return err
}

// SessionStatus expõe estado básico da sessão
type SessionStatus struct {
	SessionID string `json:"session_id"`
	Connected bool   `json:"connected"`
	LoggedIn  bool   `json:"logged_in"`
	DeviceJID string `json:"device_jid,omitempty"`
}

type SessionOverview struct {
	SessionID          string `json:"session_id"`
	Assigned           bool   `json:"assigned"`
	Status             string `json:"status"`
	Connected          bool   `json:"connected"`
	LoggedIn           bool   `json:"logged_in"`
	DeviceJID          string `json:"device_jid,omitempty"`
	HasQR              bool   `json:"has_qr"`
	WhatsAppWebVersion string `json:"whatsapp_web_version,omitempty"`
}

func (m *ClientManager) Status(sessionID string) (SessionStatus, error) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return SessionStatus{}, ErrInvalidSessionID
	}
	m.mu.RLock()
	ent, ok := m.clients[sessionID]
	m.mu.RUnlock()
	if !ok || ent == nil || ent.Client == nil {
		return SessionStatus{}, fmt.Errorf("session %s not connected", sessionID)
	}
	cli := ent.Client
	st := SessionStatus{
		SessionID: sessionID,
		Connected: cli.IsConnected(),
		LoggedIn:  cli.IsLoggedIn(),
	}
	if cli.Store != nil && cli.Store.ID != nil {
		st.DeviceJID = cli.Store.ID.String()
	}
	return st, nil
}

func (m *ClientManager) GetClient(sessionID string) (*whatsmeow.Client, bool) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return nil, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	ent, ok := m.clients[sessionID]
	if !ok || ent == nil {
		return nil, false
	}
	return ent.Client, ent.Client != nil
}

func (m *ClientManager) ListSessions() []SessionOverview {
	m.mu.RLock()
	active := make(map[string]*clientEntry, len(m.clients))
	for sid, ent := range m.clients {
		norm := normalizeSessionID(sid)
		if norm == "" {
			continue
		}
		active[norm] = ent
	}
	m.mu.RUnlock()

	allIDs := map[string]struct{}{}
	for sid := range active {
		allIDs[sid] = struct{}{}
	}
	for _, sid := range m.KnownSessions() {
		allIDs[sid] = struct{}{}
	}
	for _, sid := range m.AssignedSessions() {
		allIDs[sid] = struct{}{}
	}

	ids := make([]string, 0, len(allIDs))
	for sid := range allIDs {
		ids = append(ids, sid)
	}
	sort.Strings(ids)

	out := make([]SessionOverview, 0, len(ids))
	for _, sid := range ids {
		ov := SessionOverview{
			SessionID: sid,
			Assigned:  m.IsAssignedSession(sid),
			Status:    status.Get(sid),
			HasQR:     hasQR(sid),
		}
		if ent := active[sid]; ent != nil && ent.Client != nil {
			ov.Connected = ent.Client.IsConnected()
			ov.LoggedIn = ent.Client.IsLoggedIn()
			ov.WhatsAppWebVersion = ent.WhatsAppWebVersion
			if ent.Client.Store != nil && ent.Client.Store.ID != nil {
				ov.DeviceJID = ent.Client.Store.ID.String()
			}
		}
		if ov.Status == "" {
			if ov.Connected {
				ov.Status = "online"
			} else {
				ov.Status = "offline"
			}
		}
		out = append(out, ov)
	}
	return out
}

func (m *ClientManager) KnownSessions() []string {
	entries, err := os.ReadDir(m.sessionStore)
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		sid := normalizeSessionID(e.Name())
		if sid == "" {
			continue
		}
		out = append(out, sid)
	}
	sort.Strings(out)
	return out
}

func (m *ClientManager) sessionDirPath(sessionID string) (string, error) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return "", ErrInvalidSessionID
	}

	root := filepath.Clean(strings.TrimSpace(m.sessionStore))
	target := filepath.Clean(filepath.Join(root, sessionID))
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if rel == "." || rel == "" || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("%w: %s", ErrInvalidSessionID, sessionID)
	}
	return target, nil
}

func (m *ClientManager) clearAllSessionWebhooks(sessionID string) (bool, error) {
	sid := normalizeSessionID(sessionID)
	if sid == "" {
		return false, ErrInvalidSessionID
	}

	m.webhooksMu.Lock()
	prev := cloneSessionWebhookMap(m.sessionWebhooks)
	_, existed := m.sessionWebhooks[sid]
	if existed {
		delete(m.sessionWebhooks, sid)
	}
	m.webhooksMu.Unlock()

	if !existed {
		return false, nil
	}
	if err := m.persistSessionWebhooks(); err != nil {
		m.webhooksMu.Lock()
		m.sessionWebhooks = prev
		m.webhooksMu.Unlock()
		return false, err
	}
	return true, nil
}

func (m *ClientManager) DeleteSession(sessionID string) (bool, error) {
	sid := normalizeSessionID(sessionID)
	if sid == "" {
		return false, ErrInvalidSessionID
	}

	sessionDir, err := m.sessionDirPath(sid)
	if err != nil {
		return false, err
	}

	removed := false

	m.mu.RLock()
	_, connected := m.clients[sid]
	m.mu.RUnlock()
	if connected {
		removed = true
	}

	if err := m.Disconnect(sid); err != nil {
		return removed, err
	}

	if unassigned, err := m.UnassignSession(sid); err != nil {
		return removed, err
	} else if unassigned {
		removed = true
	}

	if cleared, _, err := m.ClearSessionConfig(sid); err != nil {
		return removed, err
	} else if cleared {
		removed = true
	}

	if cleared, err := m.clearAllSessionWebhooks(sid); err != nil {
		return removed, err
	} else if cleared {
		removed = true
	}

	if clearQR(sid) {
		removed = true
	}
	if status.Delete(sid) {
		removed = true
	}

	if _, err := os.Stat(sessionDir); err == nil {
		removed = true
		if err := os.RemoveAll(sessionDir); err != nil {
			return removed, err
		}
	} else if !os.IsNotExist(err) {
		return removed, err
	}

	return removed, nil
}

func (m *ClientManager) waitForPairingReady(sessionID string, timeout time.Duration) error {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return ErrInvalidSessionID
	}
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	deadline := time.Now().Add(timeout)
	for {
		if cli, ok := m.GetClient(sessionID); ok && cli != nil {
			if cli.Store != nil && cli.Store.ID != nil {
				return ErrAlreadyPaired
			}
			if cli.IsConnected() && hasQR(sessionID) {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return ErrPairingNotReady
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func (m *ClientManager) PairPhone(sessionID, phone string) (string, error) {
	sid := normalizeSessionID(sessionID)
	if sid == "" {
		return "", ErrInvalidSessionID
	}
	if _, err := m.sessionDirPath(sid); err != nil {
		return "", err
	}

	normalizedPhone := digitsOnly(phone)
	if normalizedPhone == "" {
		return "", ErrMissingPhone
	}

	cli, ok := m.GetClient(sid)
	if ok && cli != nil && cli.Store != nil && cli.Store.ID != nil {
		return "", ErrAlreadyPaired
	}
	if ok && cli != nil && !cli.IsConnected() {
		if err := m.Disconnect(sid); err != nil {
			return "", err
		}
		cli = nil
		ok = false
	}
	if !ok || cli == nil {
		var err error
		cli, err = m.Connect(sid)
		if err != nil {
			return "", err
		}
	}

	if err := m.waitForPairingReady(sid, 20*time.Second); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return cli.PairPhone(ctx, normalizedPhone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
}

// ResolveDest applies overrides/heuristics and LID lookup to determine the destination JID.
// It returns the normalized PN JID (s.whatsapp.net), the resolved LID (if any),
// and the final destination JID that would be used for sending.
type ResolveResult struct {
	Input        string `json:"input"`
	NormalizedPN string `json:"normalized_pn"`
	PNJID        string `json:"pn_jid"`
	DestJID      string `json:"dest_jid"`
}

func (m *ClientManager) ResolveDest(sessionID, to string) (ResolveResult, error) {
	res := ResolveResult{Input: to}
	ent, ok := m.clients[sessionID]
	if !ok || ent == nil || ent.Client == nil {
		return res, errors.New("session not connected")
	}
	// normalize digits
	rawTo := strings.TrimSpace(to)
	digits := phoneNumberToJIDDigits(rawTo)
	res.NormalizedPN = digits
	pnJ := types.NewJID(digits, types.DefaultUserServer)
	res.PNJID = pnJ.String()
	dest := pnJ
	res.DestJID = dest.String()
	return res, nil
}

func (m *ClientManager) SessionPath(sessionID string) string {
	if norm := normalizeSessionID(sessionID); norm != "" {
		sessionID = norm
	}
	return filepath.Join(m.sessionStore, sessionID)
}

func (m *ClientManager) mustHave(sessionID string) (*clientEntry, error) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return nil, ErrInvalidSessionID
	}
	if ent, ok := m.clients[sessionID]; ok && ent != nil && ent.Client != nil {
		return ent, nil
	}
	return nil, errors.New("session not connected")
}

// RestoreSavedSessions varre o SESSION_STORE e reconecta sessões que tiverem session.db
func (m *ClientManager) RestoreSavedSessions() ([]string, error) {
	candidates := m.KnownSessions()
	if len(candidates) == 0 {
		return nil, nil
	}

	if len(m.assigned) > 0 {
		filtered := make([]string, 0, len(candidates))
		for _, sid := range candidates {
			if m.IsAssignedSession(sid) {
				filtered = append(filtered, sid)
			}
		}
		candidates = filtered
	}

	var restored []string
	for _, sid := range candidates {
		dbFile := filepath.Join(m.SessionPath(sid), "session.db")
		if _, err := os.Stat(dbFile); err == nil {
			go func(id string) {
				if _, err := m.Connect(id); err != nil {
					// opcional: logar erro
				}
			}(sid)
			restored = append(restored, sid)
		}
	}
	return restored, nil
}
