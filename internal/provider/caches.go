package provider

import (
	"sync"
	"time"

	"go.mau.fi/whatsmeow/types"
)

// Small in-memory caches to avoid repeated slow lookups while keeping
// existing behavior intact. Values expire after a TTL to ensure freshness.

type stringEntry struct {
	val    string
	expiry time.Time
}

type jidEntry struct {
	val    types.JID
	expiry time.Time
}

var (
	cacheMu        sync.RWMutex
	avatarCache    = map[string]stringEntry{}
	groupNameCache = map[string]stringEntry{}
	lidADCache     = map[string]jidEntry{}

	// Tuned TTLs: names/avatars can change, but not frequently.
	avatarTTL    = 10 * time.Minute
	groupNameTTL = 10 * time.Minute
	lidMapTTL    = 30 * time.Minute
)

func cacheGetAvatar(key string) (string, bool) {
	now := time.Now()
	cacheMu.RLock()
	e, ok := avatarCache[key]
	cacheMu.RUnlock()
	if ok && now.Before(e.expiry) && e.val != "" {
		return e.val, true
	}
	return "", false
}

func cachePutAvatar(key, url string) {
	cacheMu.Lock()
	avatarCache[key] = stringEntry{val: url, expiry: time.Now().Add(avatarTTL)}
	cacheMu.Unlock()
}

func cacheGetGroupName(key string) (string, bool) {
	now := time.Now()
	cacheMu.RLock()
	e, ok := groupNameCache[key]
	cacheMu.RUnlock()
	if ok && now.Before(e.expiry) && e.val != "" {
		return e.val, true
	}
	return "", false
}

func cachePutGroupName(key, name string) {
	cacheMu.Lock()
	groupNameCache[key] = stringEntry{val: name, expiry: time.Now().Add(groupNameTTL)}
	cacheMu.Unlock()
}

func cacheGetLID(key string) (types.JID, bool) {
	now := time.Now()
	cacheMu.RLock()
	e, ok := lidADCache[key]
	cacheMu.RUnlock()
	if ok && now.Before(e.expiry) && e.val.User != "" && e.val.Server != "" {
		return e.val, true
	}
	return types.JID{}, false
}

func cachePutLID(key string, v types.JID) {
	cacheMu.Lock()
	lidADCache[key] = jidEntry{val: v, expiry: time.Now().Add(lidMapTTL)}
	cacheMu.Unlock()
}
