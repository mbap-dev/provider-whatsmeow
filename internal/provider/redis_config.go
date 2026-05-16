package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// sessionConfigPayload holds a subset of distributed config fields we care about.
type sessionConfigPayload struct {
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
	outgoingMessagesCoex        bool
	hasOutgoingMessagesCoex     bool
	proxyEnabled                bool
	hasProxyEnabled             bool
	proxyProtocol               string
	proxyHost                   string
	proxyPort                   string
	proxyUsername               string
	proxyPassword               string
}

// fetchSessionConfig reads session config for a given phone/session key from Redis.
// It returns (payload, true) if found and parsed, otherwise (_, false).
func fetchSessionConfig(redisURL, keyPrefix, phone string) (sessionConfigPayload, bool) {
	out := sessionConfigPayload{}
	if strings.TrimSpace(redisURL) == "" || strings.TrimSpace(phone) == "" {
		return out, false
	}
	if strings.TrimSpace(keyPrefix) == "" {
		keyPrefix = "provider-whatsmeow-config:"
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return out, false
	}
	c := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer c.Close()
	key := strings.TrimSpace(keyPrefix) + phone
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
	if v, ok := m["outgoingMessagesCoex"]; ok {
		if b, ok2 := asBool(v); ok2 {
			out.outgoingMessagesCoex = b
			out.hasOutgoingMessagesCoex = true
		}
	}
	if v, ok := m["proxyEnabled"]; ok {
		if b, ok2 := asBool(v); ok2 {
			out.proxyEnabled = b
			out.hasProxyEnabled = true
		}
	}
	if v, ok := m["proxyProtocol"]; ok {
		if str, ok2 := v.(string); ok2 {
			out.proxyProtocol = str
		}
	}
	if v, ok := m["proxyHost"]; ok {
		if str, ok2 := v.(string); ok2 {
			out.proxyHost = str
		}
	}
	if v, ok := m["proxyPort"]; ok {
		if str, ok2 := v.(string); ok2 {
			out.proxyPort = str
		}
	}
	if v, ok := m["proxyUsername"]; ok {
		if str, ok2 := v.(string); ok2 {
			out.proxyUsername = str
		}
	}
	if v, ok := m["proxyPassword"]; ok {
		if str, ok2 := v.(string); ok2 {
			out.proxyPassword = str
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

func jsonMarshalIndent(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
