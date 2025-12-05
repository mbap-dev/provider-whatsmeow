package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type pnCacheItem struct {
	val string
	exp time.Time
}

// resolvePNtoJID queries WhatsApp for the canonical JID of a phone number, trying BR candidates.
// Returns (jid, true) if resolved, otherwise (empty, false).
func (m *ClientManager) resolvePNtoJID(ctx context.Context, cli *whatsmeow.Client, phone string) (types.JID, bool) {
	// Check cache by E.164 input key
	e164 := normalizeE164Local(phone)
	if e164 == "" {
		e164 = phone
	}
	if it, ok := m.pnCache[e164]; ok && time.Now().Before(it.exp) {
		if j, err := types.ParseJID(it.val); err == nil {
			return j, true
		}
	}
	candidates := candidatesBR(e164)
	for _, cand := range candidates {
		// whatsmeow provides IsOnWhatsApp to check registration and canonical JID
		res, err := cli.IsOnWhatsApp(ctx, []string{cand})
		if err != nil || len(res) == 0 {
			continue
		}
		info := res[0]
		if info.IsIn && info.JID.Server == types.DefaultUserServer && info.JID.User != "" {
			m.pnCache[e164] = pnCacheItem{val: info.JID.String(), exp: time.Now().Add(24 * time.Hour)}
			return info.JID, true
		}
	}
	return types.EmptyJID, false
}

// checkUserExistsOnWhatsApp verifies whether a given phone is registered on WhatsApp.
// It mirrors the behaviour from evolution-go:
//   - on network/API errors it returns an error so the caller can decide to
//     proceed without blocking the send;
//   - when the check succeeds and the number is not registered, it returns
//     (empty, false, nil) so the caller can reject the send explicitly.
func (m *ClientManager) checkUserExistsOnWhatsApp(ctx context.Context, cli *whatsmeow.Client, phone string) (types.JID, bool, error) {
	e164 := normalizeE164Local(phone)
	cacheKey := e164
	if e164 == "" {
		e164 = phone
		cacheKey = phone
	}

	// Reuse cache when available
	if it, ok := m.pnCache[cacheKey]; ok && time.Now().Before(it.exp) {
		if j, err := types.ParseJID(it.val); err == nil {
			return j, true, nil
		}
	}

	candidates := candidatesBR(e164)
	var lastErr error
	for _, cand := range candidates {
		res, err := cli.IsOnWhatsApp(ctx, []string{cand})
		if err != nil {
			lastErr = fmt.Errorf("failed to check if number %s exists on WhatsApp: %w", cand, err)
			continue
		}
		if len(res) == 0 {
			lastErr = fmt.Errorf("number %s not found in WhatsApp response", cand)
			continue
		}
		info := res[0]
		if !info.IsIn {
			// Explicitly not registered for this candidate; try next, if any.
			continue
		}
		if info.JID.Server == types.DefaultUserServer && info.JID.User != "" {
			m.pnCache[cacheKey] = pnCacheItem{val: info.JID.String(), exp: time.Now().Add(24 * time.Hour)}
			return info.JID, true, nil
		}
	}

	// If all candidates failed due to errors or empty responses, surface the last error.
	if lastErr != nil {
		return types.EmptyJID, false, lastErr
	}

	// Check succeeded but no candidate is registered on WhatsApp.
	return types.EmptyJID, false, nil
}

func normalizeE164Local(input string) string {
	in := strings.TrimSpace(input)
	// Region BR helps when input is national format without +55
	region := "BR"
	if strings.HasPrefix(in, "+") {
		region = ""
	}
	num, err := phonenumbers.Parse(in, region)
	if err != nil {
		return ""
	}
	return phonenumbers.Format(num, phonenumbers.E164)
}

// candidatesBR: for +55 DDD local, try original, with 9, and without 9.
func candidatesBR(pnE164 string) []string {
	if !strings.HasPrefix(pnE164, "+55") || len(pnE164) < 5 {
		return []string{pnE164}
	}
	rest := pnE164[3:]
	if len(rest) < 10 {
		return []string{pnE164}
	}
	ddd := rest[:2]
	local := rest[2:]
	with9 := pnE164
	if !strings.HasPrefix(local, "9") {
		with9 = "+55" + ddd + "9" + local
	}
	without9 := pnE164
	if strings.HasPrefix(local, "9") && len(local) >= 9 {
		without9 = "+55" + ddd + local[1:]
	}
	seen := map[string]bool{}
	out := []string{}
	for _, n := range []string{pnE164, with9, without9} {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	return out
}
