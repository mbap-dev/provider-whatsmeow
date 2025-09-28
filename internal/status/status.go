package status

import (
	"context"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Manager is a minimal Redis-backed session status writer used by UNO.
// It writes keys in the format: unoapi-status:{phone}
// Accepted values: connecting, online, offline, disconnected, restart_required, standby
type Manager struct {
	client *redis.Client
	prefix string
}

var mgr *Manager

// Init initialises the global status manager with the given Redis URL.
// If url is empty, status updates are no-ops.
func Init(redisURL string) {
	if strings.TrimSpace(redisURL) == "" {
		mgr = &Manager{client: nil, prefix: "unoapi-status:"}
		return
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		// Fallback to disabled if parse fails
		mgr = &Manager{client: nil, prefix: "unoapi-status:"}
		return
	}
	c := redis.NewClient(opt)
	// Ping with short timeout to validate
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = c.Ping(ctx).Err()
	mgr = &Manager{client: c, prefix: "unoapi-status:"}
}

// Set writes the status value for the given session phone id.
// The phone id should be digits-only; callers may pass raw session ids,
// this helper strips non-digits.
func Set(sessionID, value string) {
	if mgr == nil || mgr.client == nil {
		return
	}
	phone := digitsOnly(sessionID)
	if phone == "" {
		return
	}
	key := mgr.prefix + phone
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = mgr.client.Set(ctx, key, strings.TrimSpace(value), 0).Err()
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
