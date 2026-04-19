package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"your.org/provider-whatsmeow/internal/log"
)

const (
	defaultWebhookTimeoutMS = 5000
)

// SessionWebhook follows the same core semantics used in UnoAPI:
// multiple webhooks per session with event filters.
type SessionWebhook struct {
	ID                     string `json:"id"`
	Name                   string `json:"name,omitempty"`
	URL                    string `json:"url,omitempty"`
	URLAbsolute            string `json:"url_absolute,omitempty"`
	Header                 string `json:"header,omitempty"`
	Token                  string `json:"token,omitempty"`
	TimeoutMS              int    `json:"timeout_ms"`
	Enabled                bool   `json:"enabled"`
	AppendSession          bool   `json:"append_session"`
	SendIncomingMessages   bool   `json:"send_incoming_messages"`
	SendOutgoingMessages   bool   `json:"send_outgoing_messages"`
	SendUpdateMessages     bool   `json:"send_update_messages"`
	SendGroupMessages      bool   `json:"send_group_messages"`
	SendNewsletterMessages bool   `json:"send_newsletter_messages"`
}

type sessionWebhookState struct {
	WorkerID string                      `json:"worker_id"`
	Sessions map[string][]SessionWebhook `json:"sessions"`
}

type webhookPayloadClass struct {
	Incoming   bool
	Outgoing   bool
	Update     bool
	Group      bool
	Newsletter bool
}

func (m *ClientManager) ListSessionWebhooks(sessionID string) []SessionWebhook {
	sid := strings.TrimSpace(sessionID)
	m.webhooksMu.RLock()
	hooks := cloneSessionWebhookList(m.sessionWebhooks[sid])
	m.webhooksMu.RUnlock()
	sort.SliceStable(hooks, func(i, j int) bool {
		return hooks[i].ID < hooks[j].ID
	})
	return hooks
}

func (m *ClientManager) UpsertSessionWebhook(sessionID string, hook SessionWebhook) (SessionWebhook, error) {
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		return SessionWebhook{}, errors.New("missing session id")
	}
	norm, err := normalizeSessionWebhook(sid, hook)
	if err != nil {
		return SessionWebhook{}, err
	}

	m.webhooksMu.Lock()
	prev := cloneSessionWebhookMap(m.sessionWebhooks)
	list := m.sessionWebhooks[sid]
	replaced := false
	for i := range list {
		if list[i].ID == norm.ID {
			list[i] = norm
			replaced = true
			break
		}
	}
	if !replaced {
		list = append(list, norm)
	}
	m.sessionWebhooks[sid] = list
	m.webhooksMu.Unlock()

	if err := m.persistSessionWebhooks(); err != nil {
		m.webhooksMu.Lock()
		m.sessionWebhooks = prev
		m.webhooksMu.Unlock()
		return SessionWebhook{}, err
	}
	return norm, nil
}

func (m *ClientManager) DeleteSessionWebhook(sessionID, webhookID string) (bool, error) {
	sid := strings.TrimSpace(sessionID)
	id := sanitizeWebhookID(webhookID)
	if sid == "" {
		return false, errors.New("missing session id")
	}
	if id == "" {
		return false, errors.New("missing webhook id")
	}

	m.webhooksMu.Lock()
	prev := cloneSessionWebhookMap(m.sessionWebhooks)
	list := m.sessionWebhooks[sid]
	if len(list) == 0 {
		m.webhooksMu.Unlock()
		return false, nil
	}
	next := make([]SessionWebhook, 0, len(list))
	removed := false
	for _, item := range list {
		if item.ID == id {
			removed = true
			continue
		}
		next = append(next, item)
	}
	if removed {
		if len(next) == 0 {
			delete(m.sessionWebhooks, sid)
		} else {
			m.sessionWebhooks[sid] = next
		}
	}
	m.webhooksMu.Unlock()

	if !removed {
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

func (m *ClientManager) dispatchRegisteredWebhooks(sessionID string, payload any) bool {
	hooks := m.ListSessionWebhooks(sessionID)
	if len(hooks) == 0 {
		return false
	}
	class := classifyWebhookPayload(sessionID, payload)
	posted := 0
	for _, hook := range hooks {
		if !shouldDeliverSessionWebhook(hook, class) {
			continue
		}
		posted++
		m.postSessionWebhookAsync(sessionID, hook, payload)
	}
	entry := log.WithSession(sessionID)
	if posted == 0 {
		entry.Info("webhook http skipped: no registered webhook matched event filters")
	} else {
		entry.Info("webhook http dispatch queued count=%d", posted)
	}
	return true
}

func normalizeSessionWebhook(sessionID string, hook SessionWebhook) (SessionWebhook, error) {
	out := SessionWebhook{
		ID:                     sanitizeWebhookID(hook.ID),
		Name:                   strings.TrimSpace(hook.Name),
		URL:                    strings.TrimSpace(hook.URL),
		URLAbsolute:            strings.TrimSpace(hook.URLAbsolute),
		Header:                 strings.TrimSpace(hook.Header),
		Token:                  strings.TrimSpace(hook.Token),
		TimeoutMS:              hook.TimeoutMS,
		Enabled:                hook.Enabled,
		AppendSession:          hook.AppendSession,
		SendIncomingMessages:   hook.SendIncomingMessages,
		SendOutgoingMessages:   hook.SendOutgoingMessages,
		SendUpdateMessages:     hook.SendUpdateMessages,
		SendGroupMessages:      hook.SendGroupMessages,
		SendNewsletterMessages: hook.SendNewsletterMessages,
	}
	if out.ID == "" {
		out.ID = autoWebhookID(sessionID)
	}
	if out.TimeoutMS <= 0 {
		out.TimeoutMS = defaultWebhookTimeoutMS
	}
	if out.URLAbsolute == "" && out.URL == "" {
		return SessionWebhook{}, errors.New("url or url_absolute is required")
	}
	if out.URLAbsolute != "" && !isValidWebhookURL(out.URLAbsolute) {
		return SessionWebhook{}, errors.New("url_absolute must start with http:// or https://")
	}
	if out.URL != "" && !isValidWebhookURL(out.URL) {
		return SessionWebhook{}, errors.New("url must start with http:// or https://")
	}
	return out, nil
}

func autoWebhookID(sessionID string) string {
	base := sanitizeWebhookID(sessionID)
	if base == "" {
		base = "webhook"
	}
	return fmt.Sprintf("%s-%x", base, time.Now().UnixNano())
}

func sanitizeWebhookID(v string) string {
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

func isValidWebhookURL(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")
}

func cloneSessionWebhookList(in []SessionWebhook) []SessionWebhook {
	if len(in) == 0 {
		return nil
	}
	out := make([]SessionWebhook, len(in))
	copy(out, in)
	return out
}

func cloneSessionWebhookMap(in map[string][]SessionWebhook) map[string][]SessionWebhook {
	out := make(map[string][]SessionWebhook, len(in))
	for sid, hooks := range in {
		out[sid] = cloneSessionWebhookList(hooks)
	}
	return out
}

func loadSessionWebhookState(path, workerID string) (map[string][]SessionWebhook, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return map[string][]SessionWebhook{}, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]SessionWebhook{}, nil
		}
		return nil, err
	}
	st := sessionWebhookState{}
	if err := jsonUnmarshal(b, &st); err != nil {
		return nil, err
	}
	if strings.TrimSpace(st.WorkerID) != "" && strings.TrimSpace(workerID) != "" && strings.TrimSpace(st.WorkerID) != strings.TrimSpace(workerID) {
		return map[string][]SessionWebhook{}, nil
	}
	out := map[string][]SessionWebhook{}
	for sid, hooks := range st.Sessions {
		sid = strings.TrimSpace(sid)
		if sid == "" {
			continue
		}
		list := make([]SessionWebhook, 0, len(hooks))
		for _, raw := range hooks {
			norm, err := normalizeSessionWebhook(sid, raw)
			if err != nil {
				continue
			}
			list = append(list, norm)
		}
		if len(list) > 0 {
			out[sid] = list
		}
	}
	return out, nil
}

func saveSessionWebhookState(path, workerID string, webhooks map[string][]SessionWebhook) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	payload := sessionWebhookState{
		WorkerID: strings.TrimSpace(workerID),
		Sessions: map[string][]SessionWebhook{},
	}
	for sid, hooks := range webhooks {
		sid = strings.TrimSpace(sid)
		if sid == "" {
			continue
		}
		if len(hooks) == 0 {
			continue
		}
		payload.Sessions[sid] = cloneSessionWebhookList(hooks)
	}
	b, err := jsonMarshalIndent(payload)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func (m *ClientManager) persistSessionWebhooks() error {
	m.webhooksMu.RLock()
	snapshot := cloneSessionWebhookMap(m.sessionWebhooks)
	file := m.sessionHookFile
	m.webhooksMu.RUnlock()
	return saveSessionWebhookState(file, m.workerID, snapshot)
}

func (m *ClientManager) postSessionWebhookAsync(sessionID string, hook SessionWebhook, payload any) {
	entry := log.WithSession(sessionID)
	body, err := json.Marshal(payload)
	if err != nil {
		entry.Error("webhook[%s] marshal error: %v", hook.ID, err)
		return
	}
	url := resolveSessionWebhookURL(sessionID, hook)
	if url == "" {
		entry.Error("webhook[%s] ignored: empty url", hook.ID)
		return
	}
	timeout := hook.TimeoutMS
	if timeout <= 0 {
		timeout = defaultWebhookTimeoutMS
	}
	entry.Info("webhook[%s] post queued url=%s bytes=%d", hook.ID, url, len(body))
	if len(body) > 0 {
		snippet := string(body)
		if len(snippet) > 2000 {
			snippet = snippet[:2000] + "...(truncated)"
		}
		entry.Debug("webhook[%s] payload=%s", hook.ID, snippet)
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			entry.Error("webhook[%s] build request error: %v", hook.ID, err)
			return
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		if hook.Header != "" && hook.Token != "" {
			req.Header.Set(hook.Header, hook.Token)
		}
		resp, err := webhookHTTPClient.Do(req)
		if err != nil {
			entry.Error("webhook[%s] post error: %v", hook.ID, err)
			return
		}
		defer resp.Body.Close()
		snippet := ""
		if b, _ := io.ReadAll(io.LimitReader(resp.Body, 512)); len(b) > 0 {
			snippet = string(b)
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			entry.Info("webhook[%s] post done status=%d", hook.ID, resp.StatusCode)
		} else {
			entry.Error("webhook[%s] post non-2xx status=%d body=%q", hook.ID, resp.StatusCode, snippet)
		}
	}()
}

func resolveSessionWebhookURL(sessionID string, hook SessionWebhook) string {
	if strings.TrimSpace(hook.URLAbsolute) != "" {
		return strings.TrimSpace(hook.URLAbsolute)
	}
	base := strings.TrimRight(strings.TrimSpace(hook.URL), "/")
	if base == "" {
		return ""
	}
	if hook.AppendSession {
		return base + "/" + sessionID
	}
	return base
}

func classifyWebhookPayload(sessionID string, payload any) webhookPayloadClass {
	out := webhookPayloadClass{}
	raw, _ := json.Marshal(payload)
	if strings.Contains(strings.ToLower(string(raw)), "@newsletter") {
		out.Newsletter = true
	}

	data, ok := payload.(map[string]any)
	if !ok {
		return out
	}
	value := extractCloudPayloadValue(data)
	if len(value) == 0 {
		return out
	}
	if statuses, ok := value["statuses"].([]any); ok && len(statuses) > 0 {
		out.Update = true
	}
	if messages, ok := value["messages"].([]any); ok && len(messages) > 0 {
		msg, _ := messages[0].(map[string]any)
		from := anyString(msg["from"])
		if digitsOnly(from) != "" && digitsOnly(from) == digitsOnly(sessionID) {
			out.Outgoing = true
		} else {
			out.Incoming = true
		}
	}
	if echoes, ok := value["message_echoes"].([]any); ok && len(echoes) > 0 {
		out.Outgoing = true
		msg, _ := echoes[0].(map[string]any)
		groupID := anyString(msg["group_id"])
		if groupID != "" {
			out.Group = true
		}
		if strings.Contains(strings.ToLower(groupID), "newsletter") {
			out.Newsletter = true
		}
	}
	if contacts, ok := value["contacts"].([]any); ok && len(contacts) > 0 {
		contact, _ := contacts[0].(map[string]any)
		groupID := anyString(contact["group_id"])
		if groupID != "" {
			out.Group = true
		}
		if strings.Contains(strings.ToLower(groupID), "newsletter") {
			out.Newsletter = true
		}
	}
	return out
}

func shouldDeliverSessionWebhook(hook SessionWebhook, class webhookPayloadClass) bool {
	if !hook.Enabled {
		return false
	}
	if class.Update && !hook.SendUpdateMessages {
		return false
	}
	if class.Incoming && !hook.SendIncomingMessages {
		return false
	}
	if class.Outgoing && !hook.SendOutgoingMessages {
		return false
	}
	if class.Group && !hook.SendGroupMessages {
		return false
	}
	if class.Newsletter && !hook.SendNewsletterMessages {
		return false
	}
	return true
}

func extractCloudPayloadValue(payload map[string]any) map[string]any {
	entryAny, ok := payload["entry"].([]any)
	if !ok || len(entryAny) == 0 {
		return nil
	}
	entry, ok := entryAny[0].(map[string]any)
	if !ok {
		return nil
	}
	changesAny, ok := entry["changes"].([]any)
	if !ok || len(changesAny) == 0 {
		return nil
	}
	change, ok := changesAny[0].(map[string]any)
	if !ok {
		return nil
	}
	value, _ := change["value"].(map[string]any)
	return value
}

func anyString(v any) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}
