package provider

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		p.AlwaysOnlineIntervalSeconds != nil
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
	snapshot, err := m.GetSessionConfig(sid)
	return existed, snapshot, err
}

func (m *ClientManager) effectiveSessionConfig(sessionID string) SessionConfigValues {
	alwaysOnline, interval := m.getAlwaysOnlineConfig(sessionID)
	return SessionConfigValues{
		RejectCalls:                 m.getRejectCalls(sessionID),
		RejectCallsMessage:          m.getRejectMsg(sessionID),
		MarkReadOnMessage:           m.getAutoMarkRead(sessionID),
		IgnoreStatusBroadcast:       m.getIgnoreStatusBroadcast(sessionID),
		IgnoreNewsletters:           m.getIgnoreNewsletters(sessionID),
		OutgoingMessagesCoex:        m.getOutgoingMessagesCoex(sessionID),
		AlwaysOnline:                alwaysOnline,
		AlwaysOnlineIntervalSeconds: int(interval / time.Second),
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
		sid = strings.TrimSpace(sid)
		if sid == "" {
			continue
		}
		ov := sessionOverrides{
			rejectCalls:           cloneBoolPtr(raw.RejectCalls),
			rejectMsg:             cloneStringPtr(raw.RejectCallsMessage),
			autoMarkRead:          cloneBoolPtr(raw.MarkReadOnMessage),
			ignoreStatusBroadcast: cloneBoolPtr(raw.IgnoreStatusBroadcast),
			ignoreNewsletters:     cloneBoolPtr(raw.IgnoreNewsletters),
			outgoingMessagesCoex:  cloneBoolPtr(raw.OutgoingMessagesCoex),
			alwaysOnline:          cloneBoolPtr(raw.AlwaysOnline),
		}
		if raw.AlwaysOnlineIntervalSeconds != nil && *raw.AlwaysOnlineIntervalSeconds > 0 {
			d := time.Duration(*raw.AlwaysOnlineIntervalSeconds) * time.Second
			ov.alwaysOnlineEvery = &d
		}
		out[sid] = ov
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
		sid = strings.TrimSpace(sid)
		if sid == "" {
			continue
		}
		item := sessionOverridePayload{
			RejectCalls:           cloneBoolPtr(ov.rejectCalls),
			RejectCallsMessage:    cloneStringPtr(ov.rejectMsg),
			MarkReadOnMessage:     cloneBoolPtr(ov.autoMarkRead),
			IgnoreStatusBroadcast: cloneBoolPtr(ov.ignoreStatusBroadcast),
			IgnoreNewsletters:     cloneBoolPtr(ov.ignoreNewsletters),
			OutgoingMessagesCoex:  cloneBoolPtr(ov.outgoingMessagesCoex),
			AlwaysOnline:          cloneBoolPtr(ov.alwaysOnline),
		}
		if ov.alwaysOnlineEvery != nil && *ov.alwaysOnlineEvery > 0 {
			secs := int((*ov.alwaysOnlineEvery) / time.Second)
			item.AlwaysOnlineIntervalSeconds = &secs
		}
		payload.Sessions[sid] = item
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
	file := m.sessionCfgFile
	m.overridesMu.RUnlock()
	return saveSessionOverrideState(file, m.workerID, snapshot)
}
