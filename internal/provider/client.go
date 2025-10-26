package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

	// redis url for UNO config (optional)
	redisURL string

	// per-session overrides loaded from UNO redis config
	overrides map[string]sessionOverrides
}

type sessionOverrides struct {
	// if nil, use manager default
	rejectCalls           *bool
	rejectMsg             *string
	autoMarkRead          *bool
	ignoreStatusBroadcast *bool
	ignoreNewsletters     *bool
	alwaysOnline          *bool
	alwaysOnlineEvery     *time.Duration
}

func NewClientManager(sessionStore, webhookBase string, defaultAudioPTT bool, rejectCalls bool, rejectMsg string, autoMarkRead bool, pnResolverURL string, ignoreStatusBroadcast, ignoreNewsletters bool, redisURL string) *ClientManager {
	return &ClientManager{
		clients:               make(map[string]*clientEntry),
		sessionStore:          sessionStore,
		webhookBase:           webhookBase,
		defaultAudioPTT:       defaultAudioPTT,
		rejectCalls:           rejectCalls,
		rejectMsg:             rejectMsg,
		autoMarkRead:          autoMarkRead,
		pnResolverURL:         pnResolverURL,
		pnCache:               make(map[string]pnCacheItem),
		ignoreStatusBroadcast: ignoreStatusBroadcast,
		ignoreNewsletters:     ignoreNewsletters,
		redisURL:              redisURL,
		overrides:             make(map[string]sessionOverrides),
	}
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
	// Resolve whether this session should run always-online
	should := m.alwaysOnline
	interval := m.alwaysOnlineEvery
	if ov, ok := m.overrides[sessionID]; ok {
		if ov.alwaysOnline != nil {
			should = *ov.alwaysOnline
		}
		if ov.alwaysOnlineEvery != nil && *ov.alwaysOnlineEvery > 0 {
			interval = *ov.alwaysOnlineEvery
		}
	}
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
		if err := cli.SendPresence(types.PresenceAvailable); err != nil {
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
				if err := cli.SendPresence(types.PresenceAvailable); err != nil {
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
	if ov, ok := m.overrides[sessionID]; ok && ov.rejectCalls != nil {
		return *ov.rejectCalls
	}
	return m.rejectCalls
}

func (m *ClientManager) getRejectMsg(sessionID string) string {
	if ov, ok := m.overrides[sessionID]; ok && ov.rejectMsg != nil {
		return *ov.rejectMsg
	}
	return m.rejectMsg
}

func (m *ClientManager) getAutoMarkRead(sessionID string) bool {
	if ov, ok := m.overrides[sessionID]; ok && ov.autoMarkRead != nil {
		return *ov.autoMarkRead
	}
	return m.autoMarkRead
}

func (m *ClientManager) getIgnoreStatusBroadcast(sessionID string) bool {
	if ov, ok := m.overrides[sessionID]; ok && ov.ignoreStatusBroadcast != nil {
		return *ov.ignoreStatusBroadcast
	}
	return m.ignoreStatusBroadcast
}

func (m *ClientManager) getIgnoreNewsletters(sessionID string) bool {
	if ov, ok := m.overrides[sessionID]; ok && ov.ignoreNewsletters != nil {
		return *ov.ignoreNewsletters
	}
	return m.ignoreNewsletters
}

// loadSessionOverrides fetches UnoAPI config for the given session from Redis
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
	cfg, ok := fetchUnoConfig(m.redisURL, phone)
	if !ok {
		return
	}
	// Prepare override record
	ov := m.overrides[sessionID]

	// Env presence checks: only override when env var not explicitly set
	// REJECT_CALLS and REJECT_CALLS_MESSAGE
	if _, set := os.LookupEnv("REJECT_CALLS"); !set {
		val := strings.TrimSpace(cfg.rejectCalls)
		b := val != ""
		ov.rejectCalls = &b
		if _, setMsg := os.LookupEnv("REJECT_CALLS_MESSAGE"); !setMsg && val != "" {
			// Prefer Uno message when available
			ov.rejectMsg = &val
		}
	}
	if _, set := os.LookupEnv("MARK_READ_ON_MESSAGE"); !set {
		if v, okb := cfg.readOnReceipt, cfg.hasReadOnReceipt; okb {
			ov.autoMarkRead = &v
		}
	}
	if _, set := os.LookupEnv("IGNORE_STATUS_BROADCAST"); !set {
		if v, okb := cfg.ignoreBroadcastStatuses, cfg.hasIgnoreBroadcastStatuses; okb {
			ov.ignoreStatusBroadcast = &v
		}
	}
	if _, set := os.LookupEnv("IGNORE_NEWSLETTERS"); !set {
		if v, okb := cfg.ignoreNewsletterMessages, cfg.hasIgnoreNewsletterMessages; okb {
			ov.ignoreNewsletters = &v
		}
	}
	if _, set := os.LookupEnv("ALWAYS_ONLINE"); !set {
		if v, okb := cfg.alwaysOnline, cfg.hasAlwaysOnline; okb {
			ov.alwaysOnline = &v
		}
		if v, okb := cfg.alwaysOnlineEverySeconds, cfg.hasAlwaysOnlineEverySeconds; okb && v > 0 {
			d := time.Duration(v) * time.Second
			ov.alwaysOnlineEvery = &d
		}
	}
	m.overrides[sessionID] = ov
}

// Connect cria/recupera cliente e inicia conexão (emite QR se for 1o login)
func (m *ClientManager) Connect(sessionID string) (*whatsmeow.Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ent, ok := m.clients[sessionID]; ok && ent.Client != nil {
		return ent.Client, nil
	}

	// Mark status as connecting for UNO
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
	// Resolve always-online effective
	should := m.alwaysOnline
	interval := m.alwaysOnlineEvery
	if ov, ok := m.overrides[sessionID]; ok {
		if ov.alwaysOnline != nil {
			should = *ov.alwaysOnline
		}
		if ov.alwaysOnlineEvery != nil && *ov.alwaysOnlineEvery > 0 {
			interval = *ov.alwaysOnlineEvery
		}
	}
	// Sources: env set?
	envRC, _ := os.LookupEnv("REJECT_CALLS")
	envRCM, _ := os.LookupEnv("REJECT_CALLS_MESSAGE")
	envAMR, _ := os.LookupEnv("MARK_READ_ON_MESSAGE")
	envISB, _ := os.LookupEnv("IGNORE_STATUS_BROADCAST")
	envINL, _ := os.LookupEnv("IGNORE_NEWSLETTERS")
	envAO, _ := os.LookupEnv("ALWAYS_ONLINE")
	entry.Info(
		"session_config reject_calls=%t source_reject_calls=%s reject_msg_set=%t source_reject_msg=%s mark_read_on_message=%t source_mark_read=%s ignore_status_broadcast=%t source_isb=%s ignore_newsletters=%t source_inl=%s always_online=%t source_always_online=%s always_online_interval_s=%d",
		rc,
		source(envRC, m.overrides[sessionID].rejectCalls != nil),
		msg != "",
		source(envRCM, m.overrides[sessionID].rejectMsg != nil),
		amr,
		source(envAMR, m.overrides[sessionID].autoMarkRead != nil),
		isb,
		source(envISB, m.overrides[sessionID].ignoreStatusBroadcast != nil),
		inl,
		source(envINL, m.overrides[sessionID].ignoreNewsletters != nil),
		should,
		source(envAO, m.overrides[sessionID].alwaysOnline != nil || m.overrides[sessionID].alwaysOnlineEvery != nil),
		int(interval/time.Second),
	)
}

func source(envVal string, hasOverride bool) string {
	if strings.TrimSpace(envVal) != "" {
		return "env"
	}
	if hasOverride {
		return "redis"
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
	// Mark as disconnected in UNO
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
	entries, err := os.ReadDir(m.sessionStore)
	if err != nil {
		return nil, fmt.Errorf("read session store: %w", err)
	}
	var restored []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		sid := e.Name()
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
