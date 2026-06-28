package provider

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	ilog "your.org/provider-whatsmeow/internal/log"
)

// SessionConfigValues represents the effective configuration used by a session.
type SessionConfigValues struct {
	RejectCalls                 bool   `json:"reject_calls"`
	RejectCallsMessage          string `json:"reject_calls_message"`
	MarkReadOnMessage           bool   `json:"mark_read_on_message"`
	IgnoreStatusBroadcast       bool   `json:"ignore_status_broadcast"`
	IgnoreNewsletters           bool   `json:"ignore_newsletters"`
	OutgoingMessagesCoex        bool   `json:"outgoing_messages_coex"`
	AlwaysOnline                bool   `json:"always_online"`
	AlwaysOnlineIntervalSeconds int    `json:"always_online_interval_seconds"`
	HistorySyncEnabled          bool   `json:"history_sync_enabled"`
	HistorySyncDays             int    `json:"history_sync_days"`
	ProxyEnabled                bool   `json:"proxy_enabled"`
	ProxyProtocol               string `json:"proxy_protocol"`
	ProxyHost                   string `json:"proxy_host"`
	ProxyPort                   string `json:"proxy_port"`
	ProxyUsername               string `json:"proxy_username"`
	ProxyPassword               string `json:"proxy_password"`
}

// SessionConfigPatch contains mutable per-session options.
// Nil fields mean "keep current value".
type SessionConfigPatch struct {
	RejectCalls                 *bool   `json:"reject_calls,omitempty"`
	RejectCallsMessage          *string `json:"reject_calls_message,omitempty"`
	MarkReadOnMessage           *bool   `json:"mark_read_on_message,omitempty"`
	IgnoreStatusBroadcast       *bool   `json:"ignore_status_broadcast,omitempty"`
	IgnoreNewsletters           *bool   `json:"ignore_newsletters,omitempty"`
	OutgoingMessagesCoex        *bool   `json:"outgoing_messages_coex,omitempty"`
	AlwaysOnline                *bool   `json:"always_online,omitempty"`
	AlwaysOnlineIntervalSeconds *int    `json:"always_online_interval_seconds,omitempty"`
	HistorySyncEnabled          *bool   `json:"history_sync_enabled,omitempty"`
	HistorySyncDays             *int    `json:"history_sync_days,omitempty"`
	ProxyEnabled                *bool   `json:"proxy_enabled,omitempty"`
	ProxyProtocol               *string `json:"proxy_protocol,omitempty"`
	ProxyHost                   *string `json:"proxy_host,omitempty"`
	ProxyPort                   *string `json:"proxy_port,omitempty"`
	ProxyUsername               *string `json:"proxy_username,omitempty"`
	ProxyPassword               *string `json:"proxy_password,omitempty"`
}

// SessionConfigSnapshot is returned by the session config API.
type SessionConfigSnapshot struct {
	SessionID   string              `json:"session_id"`
	Config      SessionConfigValues `json:"config"`
	Override    SessionConfigPatch  `json:"override"`
	HasOverride bool                `json:"has_override"`
}

func (p SessionConfigPatch) hasChanges() bool {
	return p.RejectCalls != nil ||
		p.RejectCallsMessage != nil ||
		p.MarkReadOnMessage != nil ||
		p.IgnoreStatusBroadcast != nil ||
		p.IgnoreNewsletters != nil ||
		p.OutgoingMessagesCoex != nil ||
		p.AlwaysOnline != nil ||
		p.AlwaysOnlineIntervalSeconds != nil ||
		p.HistorySyncEnabled != nil ||
		p.HistorySyncDays != nil ||
		p.ProxyEnabled != nil ||
		p.ProxyProtocol != nil ||
		p.ProxyHost != nil ||
		p.ProxyPort != nil ||
		p.ProxyUsername != nil ||
		p.ProxyPassword != nil
}

func (p SessionConfigPatch) hasProxyChanges() bool {
	return p.ProxyEnabled != nil ||
		p.ProxyProtocol != nil ||
		p.ProxyHost != nil ||
		p.ProxyPort != nil ||
		p.ProxyUsername != nil ||
		p.ProxyPassword != nil
}

func validateProxyPatch(fallback ProxyConfig, patch SessionConfigPatch) error {
	if !patch.hasProxyChanges() {
		return nil
	}
	enabled := fallback.Enabled
	if patch.ProxyEnabled != nil {
		enabled = *patch.ProxyEnabled
	}
	host := fallback.Host
	port := fallback.Port
	protocol := fallback.Protocol
	if patch.ProxyHost != nil {
		host = strings.TrimSpace(*patch.ProxyHost)
	}
	if patch.ProxyPort != nil {
		port = strings.TrimSpace(*patch.ProxyPort)
	}
	if patch.ProxyProtocol != nil {
		protocol = strings.TrimSpace(*patch.ProxyProtocol)
	}
	if !enabled {
		return nil
	}
	if strings.TrimSpace(host) == "" {
		return errors.New("proxy_host is required when proxy is enabled")
	}
	if strings.TrimSpace(port) == "" {
		return errors.New("proxy_port is required when proxy is enabled")
	}
	if _, err := buildProxyAddress(protocol, host, port, "", ""); err != nil {
		return err
	}
	return nil
}

func (m *ClientManager) GetSessionConfig(sessionID string) (SessionConfigSnapshot, error) {
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		return SessionConfigSnapshot{}, errors.New("missing session id")
	}
	cfg := m.effectiveSessionConfig(sid)
	override, hasOverride := m.sessionOverridePatch(sid)
	return SessionConfigSnapshot{
		SessionID:   sid,
		Config:      cfg,
		Override:    override,
		HasOverride: hasOverride,
	}, nil
}

func (m *ClientManager) UpdateSessionConfig(sessionID string, patch SessionConfigPatch) (SessionConfigSnapshot, error) {
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		return SessionConfigSnapshot{}, errors.New("missing session id")
	}
	if !patch.hasChanges() {
		return SessionConfigSnapshot{}, errors.New("empty config patch")
	}
	if patch.AlwaysOnlineIntervalSeconds != nil && *patch.AlwaysOnlineIntervalSeconds <= 0 {
		return SessionConfigSnapshot{}, errors.New("always_online_interval_seconds must be > 0")
	}
	if patch.HistorySyncDays != nil && *patch.HistorySyncDays < 0 {
		return SessionConfigSnapshot{}, errors.New("history_sync_days must be >= 0")
	}
	proxyChanged := patch.hasProxyChanges()
	if err := validateProxyPatch(m.proxyDefault, patch); err != nil {
		return SessionConfigSnapshot{}, err
	}

	m.overridesMu.Lock()
	prev, existed := m.overrides[sid]
	next := cloneSessionOverrides(prev)
	applySessionConfigPatch(&next, patch)
	m.overrides[sid] = next
	m.overridesMu.Unlock()

	if err := m.persistSessionOverrides(); err != nil {
		m.overridesMu.Lock()
		if existed {
			m.overrides[sid] = prev
		} else {
			delete(m.overrides, sid)
		}
		m.overridesMu.Unlock()
		return SessionConfigSnapshot{}, err
	}

	m.refreshAlwaysOnlineForSession(sid)
	if proxyChanged {
		m.reloadActiveSessionAsync(sid)
	}
	return m.GetSessionConfig(sid)
}

func (m *ClientManager) ClearSessionConfig(sessionID string) (bool, SessionConfigSnapshot, error) {
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		return false, SessionConfigSnapshot{}, errors.New("missing session id")
	}

	m.overridesMu.Lock()
	prev, existed := m.overrides[sid]
	if existed {
		delete(m.overrides, sid)
	}
	m.overridesMu.Unlock()

	if existed {
		if err := m.persistSessionOverrides(); err != nil {
			m.overridesMu.Lock()
			m.overrides[sid] = prev
			m.overridesMu.Unlock()
			return false, SessionConfigSnapshot{}, err
		}
	}

	m.refreshAlwaysOnlineForSession(sid)
	if existed {
		m.reloadActiveSessionAsync(sid)
	}
	snapshot, err := m.GetSessionConfig(sid)
	return existed, snapshot, err
}

func (m *ClientManager) effectiveSessionConfig(sessionID string) SessionConfigValues {
	alwaysOnline, interval := m.getAlwaysOnlineConfig(sessionID)
	proxyCfg := m.getProxyConfig(sessionID)
	return SessionConfigValues{
		RejectCalls:                 m.getRejectCalls(sessionID),
		RejectCallsMessage:          m.getRejectMsg(sessionID),
		MarkReadOnMessage:           m.getAutoMarkRead(sessionID),
		IgnoreStatusBroadcast:       m.getIgnoreStatusBroadcast(sessionID),
		IgnoreNewsletters:           m.getIgnoreNewsletters(sessionID),
		OutgoingMessagesCoex:        m.getOutgoingMessagesCoex(sessionID),
		AlwaysOnline:                alwaysOnline,
		AlwaysOnlineIntervalSeconds: int(interval / time.Second),
		HistorySyncEnabled:          m.getHistorySyncEnabled(sessionID),
		HistorySyncDays:             m.getHistorySyncDays(sessionID),
		ProxyEnabled:                proxyCfg.Enabled,
		ProxyProtocol:               proxyCfg.Protocol,
		ProxyHost:                   proxyCfg.Host,
		ProxyPort:                   proxyCfg.Port,
		ProxyUsername:               proxyCfg.Username,
		ProxyPassword:               proxyCfg.Password,
	}
}

func (m *ClientManager) sessionOverridePatch(sessionID string) (SessionConfigPatch, bool) {
	m.overridesMu.RLock()
	ov, ok := m.overrides[sessionID]
	m.overridesMu.RUnlock()
	if !ok {
		return SessionConfigPatch{}, false
	}
	patch := SessionConfigPatch{
		RejectCalls:           cloneBoolPtr(ov.rejectCalls),
		RejectCallsMessage:    cloneStringPtr(ov.rejectMsg),
		MarkReadOnMessage:     cloneBoolPtr(ov.autoMarkRead),
		IgnoreStatusBroadcast: cloneBoolPtr(ov.ignoreStatusBroadcast),
		IgnoreNewsletters:     cloneBoolPtr(ov.ignoreNewsletters),
		OutgoingMessagesCoex:  cloneBoolPtr(ov.outgoingMessagesCoex),
		AlwaysOnline:          cloneBoolPtr(ov.alwaysOnline),
		HistorySyncEnabled:    cloneBoolPtr(ov.historySyncEnabled),
		HistorySyncDays:       cloneIntPtr(ov.historySyncDays),
		ProxyEnabled:          cloneBoolPtr(ov.proxyEnabled),
		ProxyProtocol:         cloneStringPtr(ov.proxyProtocol),
		ProxyHost:             cloneStringPtr(ov.proxyHost),
		ProxyPort:             cloneStringPtr(ov.proxyPort),
		ProxyUsername:         cloneStringPtr(ov.proxyUsername),
		ProxyPassword:         cloneStringPtr(ov.proxyPassword),
	}
	if ov.alwaysOnlineEvery != nil && *ov.alwaysOnlineEvery > 0 {
		secs := int((*ov.alwaysOnlineEvery) / time.Second)
		patch.AlwaysOnlineIntervalSeconds = &secs
	}
	return patch, patch.hasChanges()
}

func applySessionConfigPatch(ov *sessionOverrides, patch SessionConfigPatch) {
	if ov == nil {
		return
	}
	if patch.RejectCalls != nil {
		v := *patch.RejectCalls
		ov.rejectCalls = &v
	}
	if patch.RejectCallsMessage != nil {
		v := strings.TrimSpace(*patch.RejectCallsMessage)
		ov.rejectMsg = &v
	}
	if patch.MarkReadOnMessage != nil {
		v := *patch.MarkReadOnMessage
		ov.autoMarkRead = &v
	}
	if patch.IgnoreStatusBroadcast != nil {
		v := *patch.IgnoreStatusBroadcast
		ov.ignoreStatusBroadcast = &v
	}
	if patch.IgnoreNewsletters != nil {
		v := *patch.IgnoreNewsletters
		ov.ignoreNewsletters = &v
	}
	if patch.OutgoingMessagesCoex != nil {
		v := *patch.OutgoingMessagesCoex
		ov.outgoingMessagesCoex = &v
	}
	if patch.AlwaysOnline != nil {
		v := *patch.AlwaysOnline
		ov.alwaysOnline = &v
	}
	if patch.AlwaysOnlineIntervalSeconds != nil && *patch.AlwaysOnlineIntervalSeconds > 0 {
		d := time.Duration(*patch.AlwaysOnlineIntervalSeconds) * time.Second
		ov.alwaysOnlineEvery = &d
	}
	if patch.HistorySyncEnabled != nil {
		v := *patch.HistorySyncEnabled
		ov.historySyncEnabled = &v
	}
	if patch.HistorySyncDays != nil && *patch.HistorySyncDays >= 0 {
		v := *patch.HistorySyncDays
		ov.historySyncDays = &v
	}
	if patch.ProxyEnabled != nil {
		v := *patch.ProxyEnabled
		ov.proxyEnabled = &v
	}
	if patch.ProxyProtocol != nil {
		v := strings.TrimSpace(*patch.ProxyProtocol)
		ov.proxyProtocol = &v
	}
	if patch.ProxyHost != nil {
		v := strings.TrimSpace(*patch.ProxyHost)
		ov.proxyHost = &v
	}
	if patch.ProxyPort != nil {
		v := strings.TrimSpace(*patch.ProxyPort)
		ov.proxyPort = &v
	}
	if patch.ProxyUsername != nil {
		v := strings.TrimSpace(*patch.ProxyUsername)
		ov.proxyUsername = &v
	}
	if patch.ProxyPassword != nil {
		v := *patch.ProxyPassword
		ov.proxyPassword = &v
	}
}

func cloneSessionOverrides(in sessionOverrides) sessionOverrides {
	out := sessionOverrides{
		rejectCalls:           cloneBoolPtr(in.rejectCalls),
		rejectMsg:             cloneStringPtr(in.rejectMsg),
		autoMarkRead:          cloneBoolPtr(in.autoMarkRead),
		ignoreStatusBroadcast: cloneBoolPtr(in.ignoreStatusBroadcast),
		ignoreNewsletters:     cloneBoolPtr(in.ignoreNewsletters),
		outgoingMessagesCoex:  cloneBoolPtr(in.outgoingMessagesCoex),
		alwaysOnline:          cloneBoolPtr(in.alwaysOnline),
		proxyEnabled:          cloneBoolPtr(in.proxyEnabled),
		historySyncEnabled:    cloneBoolPtr(in.historySyncEnabled),
		historySyncDays:       cloneIntPtr(in.historySyncDays),
		proxyProtocol:         cloneStringPtr(in.proxyProtocol),
		proxyHost:             cloneStringPtr(in.proxyHost),
		proxyPort:             cloneStringPtr(in.proxyPort),
		proxyUsername:         cloneStringPtr(in.proxyUsername),
		proxyPassword:         cloneStringPtr(in.proxyPassword),
	}
	if in.alwaysOnlineEvery != nil {
		v := *in.alwaysOnlineEvery
		out.alwaysOnlineEvery = &v
	}
	return out
}

func cloneBoolPtr(in *bool) *bool {
	if in == nil {
		return nil
	}
	v := *in
	return &v
}

func cloneStringPtr(in *string) *string {
	if in == nil {
		return nil
	}
	v := *in
	return &v
}

func cloneIntPtr(in *int) *int {
	if in == nil {
		return nil
	}
	v := *in
	return &v
}

func (m *ClientManager) getHistorySyncEnabled(sessionID string) bool {
	m.overridesMu.RLock()
	defer m.overridesMu.RUnlock()
	if ov, ok := m.overrides[sessionID]; ok && ov.historySyncEnabled != nil {
		return *ov.historySyncEnabled
	}
	return false
}

func (m *ClientManager) getHistorySyncDays(sessionID string) int {
	m.overridesMu.RLock()
	defer m.overridesMu.RUnlock()
	if ov, ok := m.overrides[sessionID]; ok && ov.historySyncDays != nil && *ov.historySyncDays > 0 {
		return *ov.historySyncDays
	}
	return 0
}

func (m *ClientManager) refreshAlwaysOnlineForSession(sessionID string) {
	m.mu.RLock()
	ent := m.clients[sessionID]
	m.mu.RUnlock()
	if ent == nil || ent.Client == nil {
		return
	}
	alwaysOnline, _ := m.getAlwaysOnlineConfig(sessionID)
	if alwaysOnline {
		m.maybeStartAlwaysOnline(sessionID, ent.Client)
		return
	}
	m.stopAlwaysOnline(sessionID)
}

func (m *ClientManager) reloadActiveSessionAsync(sessionID string) {
	m.mu.RLock()
	_, active := m.clients[sessionID]
	m.mu.RUnlock()
	if !active {
		return
	}
	go func() {
		if err := m.Reload(sessionID); err != nil {
			ilog.WithSession(sessionID).Error("proxy config reload failed err=%v", err)
		}
	}()
}

func (m *ClientManager) getAlwaysOnlineConfig(sessionID string) (bool, time.Duration) {
	alwaysOnline := m.alwaysOnline
	interval := m.alwaysOnlineEvery
	m.overridesMu.RLock()
	if ov, ok := m.overrides[sessionID]; ok {
		if ov.alwaysOnline != nil {
			alwaysOnline = *ov.alwaysOnline
		}
		if ov.alwaysOnlineEvery != nil && *ov.alwaysOnlineEvery > 0 {
			interval = *ov.alwaysOnlineEvery
		}
	}
	m.overridesMu.RUnlock()
	if interval <= 0 {
		interval = time.Minute
	}
	return alwaysOnline, interval
}

type sessionOverrideState struct {
	WorkerID string                            `json:"worker_id"`
	Sessions map[string]sessionOverridePayload `json:"sessions"`
}

type sessionOverridePayload struct {
	RejectCalls                 *bool   `json:"reject_calls,omitempty"`
	RejectCallsMessage          *string `json:"reject_calls_message,omitempty"`
	MarkReadOnMessage           *bool   `json:"mark_read_on_message,omitempty"`
	IgnoreStatusBroadcast       *bool   `json:"ignore_status_broadcast,omitempty"`
	IgnoreNewsletters           *bool   `json:"ignore_newsletters,omitempty"`
	OutgoingMessagesCoex        *bool   `json:"outgoing_messages_coex,omitempty"`
	AlwaysOnline                *bool   `json:"always_online,omitempty"`
	AlwaysOnlineIntervalSeconds *int    `json:"always_online_interval_seconds,omitempty"`
	HistorySyncEnabled          *bool   `json:"history_sync_enabled,omitempty"`
	HistorySyncDays             *int    `json:"history_sync_days,omitempty"`
	ProxyEnabled                *bool   `json:"proxy_enabled,omitempty"`
	ProxyProtocol               *string `json:"proxy_protocol,omitempty"`
	ProxyHost                   *string `json:"proxy_host,omitempty"`
	ProxyPort                   *string `json:"proxy_port,omitempty"`
	ProxyUsername               *string `json:"proxy_username,omitempty"`
	ProxyPassword               *string `json:"proxy_password,omitempty"`
}

// payloadToOverrides converts the serialized payload form into the in-memory
// override representation. Shared by the file and Postgres stores.
func payloadToOverrides(raw sessionOverridePayload) sessionOverrides {
	ov := sessionOverrides{
		rejectCalls:           cloneBoolPtr(raw.RejectCalls),
		rejectMsg:             cloneStringPtr(raw.RejectCallsMessage),
		autoMarkRead:          cloneBoolPtr(raw.MarkReadOnMessage),
		ignoreStatusBroadcast: cloneBoolPtr(raw.IgnoreStatusBroadcast),
		ignoreNewsletters:     cloneBoolPtr(raw.IgnoreNewsletters),
		outgoingMessagesCoex:  cloneBoolPtr(raw.OutgoingMessagesCoex),
		alwaysOnline:          cloneBoolPtr(raw.AlwaysOnline),
		historySyncEnabled:    cloneBoolPtr(raw.HistorySyncEnabled),
		historySyncDays:       cloneIntPtr(raw.HistorySyncDays),
		proxyEnabled:          cloneBoolPtr(raw.ProxyEnabled),
		proxyProtocol:         cloneStringPtr(raw.ProxyProtocol),
		proxyHost:             cloneStringPtr(raw.ProxyHost),
		proxyPort:             cloneStringPtr(raw.ProxyPort),
		proxyUsername:         cloneStringPtr(raw.ProxyUsername),
		proxyPassword:         cloneStringPtr(raw.ProxyPassword),
	}
	if raw.AlwaysOnlineIntervalSeconds != nil && *raw.AlwaysOnlineIntervalSeconds > 0 {
		d := time.Duration(*raw.AlwaysOnlineIntervalSeconds) * time.Second
		ov.alwaysOnlineEvery = &d
	}
	return ov
}

// overridesToPayload converts the in-memory override into its serialized form.
func overridesToPayload(ov sessionOverrides) sessionOverridePayload {
	item := sessionOverridePayload{
		RejectCalls:           cloneBoolPtr(ov.rejectCalls),
		RejectCallsMessage:    cloneStringPtr(ov.rejectMsg),
		MarkReadOnMessage:     cloneBoolPtr(ov.autoMarkRead),
		IgnoreStatusBroadcast: cloneBoolPtr(ov.ignoreStatusBroadcast),
		IgnoreNewsletters:     cloneBoolPtr(ov.ignoreNewsletters),
		OutgoingMessagesCoex:  cloneBoolPtr(ov.outgoingMessagesCoex),
		AlwaysOnline:          cloneBoolPtr(ov.alwaysOnline),
		HistorySyncEnabled:    cloneBoolPtr(ov.historySyncEnabled),
		HistorySyncDays:       cloneIntPtr(ov.historySyncDays),
		ProxyEnabled:          cloneBoolPtr(ov.proxyEnabled),
		ProxyProtocol:         cloneStringPtr(ov.proxyProtocol),
		ProxyHost:             cloneStringPtr(ov.proxyHost),
		ProxyPort:             cloneStringPtr(ov.proxyPort),
		ProxyUsername:         cloneStringPtr(ov.proxyUsername),
		ProxyPassword:         cloneStringPtr(ov.proxyPassword),
	}
	if ov.alwaysOnlineEvery != nil && *ov.alwaysOnlineEvery > 0 {
		secs := int((*ov.alwaysOnlineEvery) / time.Second)
		item.AlwaysOnlineIntervalSeconds = &secs
	}
	return item
}

func loadSessionOverrideState(path, workerID string) (map[string]sessionOverrides, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return map[string]sessionOverrides{}, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]sessionOverrides{}, nil
		}
		return nil, err
	}
	st := sessionOverrideState{}
	if err := jsonUnmarshal(b, &st); err != nil {
		return nil, err
	}
	if strings.TrimSpace(st.WorkerID) != "" && strings.TrimSpace(workerID) != "" && strings.TrimSpace(st.WorkerID) != strings.TrimSpace(workerID) {
		return map[string]sessionOverrides{}, nil
	}
	out := map[string]sessionOverrides{}
	for sid, raw := range st.Sessions {
		sid = normalizeSessionID(sid)
		if sid == "" {
			continue
		}
		out[sid] = payloadToOverrides(raw)
	}
	return out, nil
}

func saveSessionOverrideState(path, workerID string, overrides map[string]sessionOverrides) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	payload := sessionOverrideState{
		WorkerID: strings.TrimSpace(workerID),
		Sessions: map[string]sessionOverridePayload{},
	}
	for sid, ov := range overrides {
		sid = normalizeSessionID(sid)
		if sid == "" {
			continue
		}
		payload.Sessions[sid] = overridesToPayload(ov)
	}
	b, err := jsonMarshalIndent(payload)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func (m *ClientManager) persistSessionOverrides() error {
	m.overridesMu.RLock()
	snapshot := make(map[string]sessionOverrides, len(m.overrides))
	for sid, ov := range m.overrides {
		snapshot[sid] = cloneSessionOverrides(ov)
	}
	m.overridesMu.RUnlock()
	return m.store.SaveOverrides(context.Background(), m.workerID, snapshot)
}
