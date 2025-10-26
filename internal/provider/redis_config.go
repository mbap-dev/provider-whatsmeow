package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// unoConfigPayload holds a subset of UnoAPI config fields we care about.
type unoConfigPayload struct {
	// raw fields with presence information
	rejectCalls                 string
	readOnReceipt               bool
	hasReadOnReceipt            bool
	ignoreBroadcastStatuses     bool
	hasIgnoreBroadcastStatuses  bool
	ignoreNewsletterMessages    bool
	hasIgnoreNewsletterMessages bool
	alwaysOnline                bool
	hasAlwaysOnline             bool
	alwaysOnlineEverySeconds    int
	hasAlwaysOnlineEverySeconds bool
}

// fetchUnoConfig reads UnoAPI config for a given phone (digits) from Redis.
// It returns (payload, true) if found and parsed, otherwise (_, false).
func fetchUnoConfig(redisURL, phone string) (unoConfigPayload, bool) {
	out := unoConfigPayload{}
	if strings.TrimSpace(redisURL) == "" || strings.TrimSpace(phone) == "" {
		return out, false
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return out, false
	}
	c := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer c.Close()
	key := "unoapi-config:" + phone
	s, err := c.Get(ctx, key).Result()
	if err != nil || strings.TrimSpace(s) == "" {
		return out, false
	}
	// We only need a small subset; using a map allows presence detection
	// without defining the entire schema.
	// Map-based decode to capture presence, then extract typed values.
	m := map[string]any{}
	if err := jsonUnmarshal([]byte(s), &m); err != nil {
		return out, false
	}
	// rejectCalls as string message (non-empty means enabled)
	if v, ok := m["rejectCalls"]; ok {
		if str, ok2 := v.(string); ok2 {
			out.rejectCalls = str
		}
	}
	if v, ok := m["readOnReceipt"]; ok {
		if b, ok2 := asBool(v); ok2 {
			out.readOnReceipt = b
			out.hasReadOnReceipt = true
		}
	}
	if v, ok := m["ignoreBroadcastStatuses"]; ok {
		if b, ok2 := asBool(v); ok2 {
			out.ignoreBroadcastStatuses = b
			out.hasIgnoreBroadcastStatuses = true
		}
	}
	if v, ok := m["ignoreNewsletterMessages"]; ok {
		if b, ok2 := asBool(v); ok2 {
			out.ignoreNewsletterMessages = b
			out.hasIgnoreNewsletterMessages = true
		}
	}
	// Optional: support alwaysOnline flags if present in config
	if v, ok := m["alwaysOnline"]; ok {
		if b, ok2 := asBool(v); ok2 {
			out.alwaysOnline = b
			out.hasAlwaysOnline = true
		}
	}
	if v, ok := m["alwaysOnlineIntervalSeconds"]; ok {
		if n, ok2 := asInt(v); ok2 {
			out.alwaysOnlineEverySeconds = n
			out.hasAlwaysOnlineEverySeconds = true
		}
	}
	return out, true
}

// Helpers to decode loosely-typed JSON values
func asBool(v any) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		if s == "true" || s == "1" {
			return true, true
		}
		if s == "false" || s == "0" {
			return false, true
		}
		return false, false
	case float64:
		return t != 0, true
	default:
		return false, false
	}
}

func asInt(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return 0, false
		}
		// simple parse
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}

// small wrapper to avoid importing encoding/json everywhere in this package file
func jsonUnmarshal(b []byte, v any) error {
	return json.Unmarshal(b, v)
}
