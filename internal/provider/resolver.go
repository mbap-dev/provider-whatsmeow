package provider

import (
	"context"
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
		res, err := cli.IsOnWhatsApp([]string{cand})
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
