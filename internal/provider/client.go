package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// clientEntry mantém refs para facilitar cleanup
type clientEntry struct {
	Client  *whatsmeow.Client
	DB      *sqlstore.Container
	KeysDB  *sqlstore.Container // opcional
	Log     waLog.Logger
	Session string
	Webhook string
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
	preferLID       bool
}

func NewClientManager(sessionStore, webhookBase string, defaultAudioPTT bool, rejectCalls bool, rejectMsg string, autoMarkRead bool, preferLID bool) *ClientManager {
	return &ClientManager{
		clients:         make(map[string]*clientEntry),
		sessionStore:    sessionStore,
		webhookBase:     webhookBase,
		defaultAudioPTT: defaultAudioPTT,
		rejectCalls:     rejectCalls,
		rejectMsg:       rejectMsg,
		autoMarkRead:    autoMarkRead,
		preferLID:       preferLID,
	}
}

// Connect cria/recupera cliente e inicia conexão (emite QR se for 1o login)
func (m *ClientManager) Connect(sessionID string) (*whatsmeow.Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ent, ok := m.clients[sessionID]; ok && ent.Client != nil {
		return ent.Client, nil
	}

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
		}
	}()

	return cli, nil
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
	return nil
}

func (m *ClientManager) Reload(sessionID string) error {
	m.mu.RLock()
	_, ok := m.clients[sessionID]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("session %s not found", sessionID)
	}
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
