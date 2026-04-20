package status

import (
	"context"
	"strings"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Manager is a minimal Redis-backed session status writer.
// It writes keys in the format: <prefix>{session}
// Accepted values: connecting, online, offline, disconnected, restart_required, standby
type Manager struct {
	client *redis.Client
	prefix string
	mu     sync.RWMutex
	cache  map[string]string
}

var mgr *Manager

// Init initialises the global status manager with the given Redis URL.
// If url is empty, status updates are no-ops.
func Init(redisURL, prefix string) {
	if strings.TrimSpace(prefix) == "" {
		prefix = "provider-whatsmeow-status:"
	}

	if strings.TrimSpace(redisURL) == "" {
		mgr = &Manager{client: nil, prefix: prefix, cache: map[string]string{}}
		return
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		// Fallback to disabled if parse fails
		mgr = &Manager{client: nil, prefix: prefix, cache: map[string]string{}}
		return
	}
	c := redis.NewClient(opt)
	// Ping with short timeout to validate
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = c.Ping(ctx).Err()
	mgr = &Manager{client: c, prefix: prefix, cache: map[string]string{}}
}

// Set writes the status value for the given session phone id.
// The phone id should be digits-only; callers may pass raw session ids,
// this helper strips non-digits.
func Set(sessionID, value string) {
	if mgr == nil {
		return
	}
	phone := normalizeSessionKey(sessionID)
	if phone == "" {
		return
	}
	value = strings.TrimSpace(value)
	mgr.mu.Lock()
	mgr.cache[phone] = value
	mgr.mu.Unlock()

	if mgr.client == nil {
		return
	}
	key := mgr.prefix + phone
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = mgr.client.Set(ctx, key, value, 0).Err()
}

// Get returns the most recent known status for a session. The value is sourced
// from in-memory cache and falls back to Redis when available.
func Get(sessionID string) string {
	if mgr == nil {
		return ""
	}
	phone := normalizeSessionKey(sessionID)
	if phone == "" {
		return ""
	}
	mgr.mu.RLock()
	if v := mgr.cache[phone]; v != "" {
		mgr.mu.RUnlock()
		return v
	}
	mgr.mu.RUnlock()

	if mgr.client == nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	v, err := mgr.client.Get(ctx, mgr.prefix+phone).Result()
	if err != nil || strings.TrimSpace(v) == "" {
		return ""
	}
	mgr.mu.Lock()
	mgr.cache[phone] = strings.TrimSpace(v)
	mgr.mu.Unlock()
	return strings.TrimSpace(v)
}

// Delete removes the cached/redis-backed status for a session.
func Delete(sessionID string) bool {
	if mgr == nil {
		return false
	}
	phone := normalizeSessionKey(sessionID)
	if phone == "" {
		return false
	}

	mgr.mu.Lock()
	_, existed := mgr.cache[phone]
	delete(mgr.cache, phone)
	mgr.mu.Unlock()

	if mgr.client == nil {
		return existed
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	removed, err := mgr.client.Del(ctx, mgr.prefix+phone).Result()
	if err != nil {
		return existed
	}
	return existed || removed > 0
}

// Snapshot returns a copy of all statuses tracked in memory.
func Snapshot() map[string]string {
	if mgr == nil {
		return map[string]string{}
	}
	out := map[string]string{}
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()
	for k, v := range mgr.cache {
		out[k] = v
	}
	return out
}

func normalizeSessionKey(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if looksLikePhoneIdentifier(s) {
		if d := digitsOnly(s); d != "" {
			return d
		}
	}
	return s
}

func digitsOnly(s string) string {
	b := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b = append(b, r)
		}
	}
	return string(b)
}

func looksLikePhoneIdentifier(v string) bool {
	hasDigit := false
	for _, r := range v {
		switch {
		case r >= '0' && r <= '9':
			hasDigit = true
		case r == '+' || r == ' ' || r == '-' || r == '(' || r == ')':
		default:
			return false
		}
	}
	return hasDigit
}
