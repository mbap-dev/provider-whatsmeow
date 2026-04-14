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
	receiptCache   = map[string]time.Time{}
	syntheticCache = map[string]time.Time{}

	// Tuned TTLs: names/avatars can change, but not frequently.
	avatarTTL    = 10 * time.Minute
	groupNameTTL = 10 * time.Minute
	lidMapTTL    = 30 * time.Minute
	// Keep receipt dedupe entries long enough to collapse repeated WA receipt
	// fan-out, while eventually releasing memory.
	receiptTTL = 6 * time.Hour
	// Synthetic outgoing message webhook entries are short-lived and only used
	// to skip the equivalent native event once.
	syntheticTTL = 15 * time.Minute
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

// cacheTryPutReceipt returns true when the receipt key is new/expired and should
// be emitted. It returns false when the same receipt was seen recently.
func cacheTryPutReceipt(key string) bool {
	now := time.Now()

	cacheMu.Lock()
	defer cacheMu.Unlock()

	if exp, ok := receiptCache[key]; ok && now.Before(exp) {
		return false
	}
	receiptCache[key] = now.Add(receiptTTL)

	// Opportunistic cleanup to avoid unbounded growth.
	if len(receiptCache)%1024 == 0 {
		for k, exp := range receiptCache {
			if now.After(exp) {
				delete(receiptCache, k)
			}
		}
	}
	return true
}

// cacheMarkSyntheticOutgoing stores a short-lived marker for a synthetic
// outgoing message payload that was already emitted to webhook.
func cacheMarkSyntheticOutgoing(sessionID, msgID string) {
	if sessionID == "" || msgID == "" {
		return
	}
	now := time.Now()
	key := sessionID + "|" + msgID

	cacheMu.Lock()
	syntheticCache[key] = now.Add(syntheticTTL)
	// Opportunistic cleanup.
	if len(syntheticCache)%256 == 0 {
		for k, exp := range syntheticCache {
			if now.After(exp) {
				delete(syntheticCache, k)
			}
		}
	}
	cacheMu.Unlock()
}

// cacheConsumeSyntheticOutgoing returns true when a matching synthetic marker
// exists and consumes it (single use).
func cacheConsumeSyntheticOutgoing(sessionID, msgID string) bool {
	if sessionID == "" || msgID == "" {
		return false
	}
	key := sessionID + "|" + msgID
	now := time.Now()

	cacheMu.Lock()
	defer cacheMu.Unlock()

	exp, ok := syntheticCache[key]
	if !ok {
		return false
	}
	delete(syntheticCache, key)
	return now.Before(exp)
}
