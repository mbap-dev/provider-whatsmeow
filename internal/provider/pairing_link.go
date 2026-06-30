package provider

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	ErrPairingLinkNotFound = errors.New("pairing link not found")
	ErrPairingLinkExpired  = errors.New("pairing link expired")
)

type PairingLink struct {
	SessionID string     `json:"session_id"`
	TokenHash string     `json:"token_hash"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func (l PairingLink) Active() bool {
	return strings.TrimSpace(l.SessionID) != "" &&
		strings.TrimSpace(l.TokenHash) != "" &&
		l.UsedAt == nil &&
		l.RevokedAt == nil
}

func (l PairingLink) Status() string {
	switch {
	case l.UsedAt != nil:
		return "used"
	case l.RevokedAt != nil:
		return "revoked"
	default:
		return "active"
	}
}

func (m *ClientManager) CreatePairingLink(ctx context.Context, sessionID string) (PairingLink, string, error) {
	sid := normalizeSessionID(sessionID)
	if sid == "" {
		return PairingLink{}, "", ErrMissingSessionID
	}
	token, err := generatePairingToken()
	if err != nil {
		return PairingLink{}, "", err
	}
	now := time.Now().UTC()
	link := PairingLink{
		SessionID: sid,
		TokenHash: hashPairingToken(token),
		CreatedAt: now,
	}

	m.pairingLinksMu.Lock()
	for i := range m.pairingLinks {
		if m.pairingLinks[i].SessionID == sid && m.pairingLinks[i].Active() {
			m.pairingLinks[i].RevokedAt = &now
		}
	}
	m.pairingLinks = append(m.pairingLinks, link)
	links := append([]PairingLink(nil), m.pairingLinks...)
	m.pairingLinksMu.Unlock()

	if err := m.store.SavePairingLinks(ctx, m.workerID, links); err != nil {
		return PairingLink{}, "", err
	}
	return link, token, nil
}

func (m *ClientManager) ListPairingLinks(sessionID string) []PairingLink {
	sid := normalizeSessionID(sessionID)
	m.pairingLinksMu.RLock()
	defer m.pairingLinksMu.RUnlock()
	out := make([]PairingLink, 0)
	for _, link := range m.pairingLinks {
		if sid == "" || link.SessionID == sid {
			out = append(out, link)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func (m *ClientManager) ActivePairingLink(sessionID string) (PairingLink, bool) {
	for _, link := range m.ListPairingLinks(sessionID) {
		if link.Active() {
			return link, true
		}
	}
	return PairingLink{}, false
}

func (m *ClientManager) RevokePairingLinks(ctx context.Context, sessionID string) (int, error) {
	sid := normalizeSessionID(sessionID)
	if sid == "" {
		return 0, ErrMissingSessionID
	}
	now := time.Now().UTC()
	revoked := 0
	m.pairingLinksMu.Lock()
	for i := range m.pairingLinks {
		if m.pairingLinks[i].SessionID == sid && m.pairingLinks[i].Active() {
			m.pairingLinks[i].RevokedAt = &now
			revoked++
		}
	}
	links := append([]PairingLink(nil), m.pairingLinks...)
	m.pairingLinksMu.Unlock()
	if revoked == 0 {
		return 0, nil
	}
	if err := m.store.SavePairingLinks(ctx, m.workerID, links); err != nil {
		return 0, err
	}
	return revoked, nil
}

func (m *ClientManager) ResolvePairingLink(token string) (PairingLink, error) {
	link, ok := m.LookupPairingLink(token)
	if !ok {
		return PairingLink{}, ErrPairingLinkNotFound
	}
	if !link.Active() {
		return PairingLink{}, ErrPairingLinkExpired
	}
	return link, nil
}

func (m *ClientManager) LookupPairingLink(token string) (PairingLink, bool) {
	tokenHash := hashPairingToken(token)
	if tokenHash == "" {
		return PairingLink{}, false
	}
	m.pairingLinksMu.RLock()
	defer m.pairingLinksMu.RUnlock()
	for _, link := range m.pairingLinks {
		if link.TokenHash != tokenHash {
			continue
		}
		return link, true
	}
	return PairingLink{}, false
}

func (m *ClientManager) MarkPairingLinksUsed(ctx context.Context, sessionID string) {
	sid := normalizeSessionID(sessionID)
	if sid == "" {
		return
	}
	now := time.Now().UTC()
	changed := false
	m.pairingLinksMu.Lock()
	for i := range m.pairingLinks {
		if m.pairingLinks[i].SessionID == sid && m.pairingLinks[i].Active() {
			m.pairingLinks[i].UsedAt = &now
			changed = true
		}
	}
	links := append([]PairingLink(nil), m.pairingLinks...)
	m.pairingLinksMu.Unlock()
	if changed {
		_ = m.store.SavePairingLinks(ctx, m.workerID, links)
	}
}

func generatePairingToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate pairing token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func hashPairingToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

type pairingLinkFileState map[string][]PairingLink

func loadPairingLinkState(path, workerID string) ([]PairingLink, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil
	}
	var state pairingLinkFileState
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	if err := jsonUnmarshal(b, &state); err != nil {
		return nil, err
	}
	if state == nil {
		return nil, nil
	}
	return state[workerID], nil
}

func savePairingLinkState(path, workerID string, links []PairingLink) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	state := pairingLinkFileState{}
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		_ = jsonUnmarshal(data, &state)
	}
	state[workerID] = links
	b, err := jsonMarshalIndent(state)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}
