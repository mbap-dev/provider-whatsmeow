package provider

import (
	"context"
	"errors"
	"fmt"
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
	Client  *whatsmeow.Client
	DB      *sqlstore.Container
	KeysDB  *sqlstore.Container // opcional
	Log     waLog.Logger
	Session string
	Webhook string
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
	webhookBase     string
	defaultAudioPTT bool
	rejectCalls     bool
	rejectMsg       string
	autoMarkRead    bool
	pnResolverURL   string
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
	ignoreStatusBroadcast *bool
	ignoreNewsletters     *bool
	outgoingMessagesCoex  *bool
	alwaysOnline          *bool
	alwaysOnlineEvery     *time.Duration
}

type ManagerOptions struct {
	SessionStore          string
	WebhookBase           string
	DefaultAudioPTT       bool
	RejectCalls           bool
	RejectMessage         string
	AutoMarkRead          bool
	PNResolverURL         string
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
		clients:               make(map[string]*clientEntry),
		sessionStore:          opts.SessionStore,
		assignmentFile:        assignmentFile,
		sessionCfgFile:        sessionCfgFile,
		sessionHookFile:       sessionHookFile,
		webhookBase:           opts.WebhookBase,
		defaultAudioPTT:       opts.DefaultAudioPTT,
		rejectCalls:           opts.RejectCalls,
		rejectMsg:             opts.RejectMessage,
		autoMarkRead:          opts.AutoMarkRead,
		pnResolverURL:         opts.PNResolverURL,
		pnCache:               make(map[string]pnCacheItem),
		alwaysOnlineEvery:     time.Minute,
		ignoreStatusBroadcast: opts.IgnoreStatusBroadcast,
		ignoreNewsletters:     opts.IgnoreNewsletters,
		outgoingMessagesCoex:  opts.OutgoingMessagesCoex,
		redisURL:              opts.RedisURL,
		sessionConfigPrefix:   opts.SessionConfigPrefix,
		overrides:             overrides,
		sessionWebhooks:       sessionWebhooks,
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
	if looksLikePhoneIdentifier(sessionID) {
		if d := digitsOnly(sessionID); d != "" {
			return d
		}
	}
	return sessionID
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
		return "", errors.New("invalid session id")
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
		return false, errors.New("invalid session id")
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
	m.overrides[sessionID] = ov
	m.overridesMu.Unlock()
}

// Connect cria/recupera cliente e inicia conexão (emite QR se for 1o login)
func (m *ClientManager) Connect(sessionID string) (*whatsmeow.Client, error) {
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

	dbLog := waLog.Stdout("Database", "INFO", true)
	clientLog := waLog.Stdout("Client", "INFO", true)

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

	// Se já existe ID no store, a sessão já foi pareada antes -> conecta direto
	// Caso contrário, liga o canal de QR ANTES de conectar.
	if cli.Store != nil && cli.Store.ID != nil {
		// Load per-session overrides (best-effort) before registering handlers
		m.loadSessionOverrides(sessionID)
		m.logSessionConfig(sessionID)

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
		go func() {
			if err := cli.Connect(); err != nil {
				clientLog.Errorf("connect error: %v", err)
				status.Set(sessionID, "offline")
			}
		}()
		return cli, nil
	}

	// Ainda não pareado: abre canal de QR ANTES do Connect e inicia watcher
	qrCh, err := cli.GetQRChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("get QR channel: %w", err)
	}
	m.startQRWatcher(sessionID, qrCh)

	// Load per-session overrides (best-effort) before registering handlers
	m.loadSessionOverrides(sessionID)
	m.logSessionConfig(sessionID)

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
	go func() {
		if err := cli.Connect(); err != nil {
			clientLog.Errorf("connect error: %v", err)
			status.Set(sessionID, "offline")
		}
	}()

	return cli, nil
}

func (m *ClientManager) logSessionConfig(sessionID string) {
	entry := ilog.WithSession(sessionID)
	rc := m.getRejectCalls(sessionID)
	msg := strings.TrimSpace(m.getRejectMsg(sessionID))
	amr := m.getAutoMarkRead(sessionID)
	isb := m.getIgnoreStatusBroadcast(sessionID)
	inl := m.getIgnoreNewsletters(sessionID)
	coex := m.getOutgoingMessagesCoex(sessionID)
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
	entry.Info(
		"session_config reject_calls=%t source_reject_calls=%s reject_msg_set=%t source_reject_msg=%s mark_read_on_message=%t source_mark_read=%s ignore_status_broadcast=%t source_isb=%s ignore_newsletters=%t source_inl=%s outgoing_messages_coex=%t source_coex=%s always_online=%t source_always_online=%s always_online_interval_s=%d",
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
	m.mu.Lock()
	defer m.mu.Unlock()
	if ent, ok := m.clients[sessionID]; ok {
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
	SessionID string `json:"session_id"`
	Assigned  bool   `json:"assigned"`
	Status    string `json:"status"`
	Connected bool   `json:"connected"`
	LoggedIn  bool   `json:"logged_in"`
	DeviceJID string `json:"device_jid,omitempty"`
	HasQR     bool   `json:"has_qr"`
}

func (m *ClientManager) Status(sessionID string) (SessionStatus, error) {
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
		active[sid] = ent
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
		sid := strings.TrimSpace(e.Name())
		if sid == "" {
			continue
		}
		out = append(out, sid)
	}
	sort.Strings(out)
	return out
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
	return filepath.Join(m.sessionStore, sessionID)
}

func (m *ClientManager) mustHave(sessionID string) (*clientEntry, error) {
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
