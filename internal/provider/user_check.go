package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type CheckUser struct {
	Query        string  `json:"Query"`
	IsInWhatsapp bool    `json:"IsInWhatsapp"`
	JID          string  `json:"JID"`
	RemoteJID    string  `json:"RemoteJID"`
	LID          *string `json:"LID"`
	VerifiedName string  `json:"VerifiedName"`
}

type CheckUserCollection struct {
	Users []CheckUser `json:"Users"`
}

func (m *ClientManager) CheckUsers(ctx context.Context, sessionID string, numbers []string, formatJID bool) (*CheckUserCollection, error) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return nil, ErrInvalidSessionID
	}
	if len(numbers) == 0 {
		return nil, errors.New("phone number is required")
	}

	cli, ok := m.GetClient(sessionID)
	if !ok || cli == nil {
		return nil, fmt.Errorf("session %s not connected", sessionID)
	}
	if !cli.IsConnected() {
		return nil, errors.New("client disconnected")
	}

	out := &CheckUserCollection{Users: make([]CheckUser, 0, len(numbers))}
	for _, raw := range numbers {
		query := strings.TrimSpace(raw)
		if query == "" {
			return nil, errors.New("phone number is required")
		}

		user, err := m.checkOneUser(ctx, cli, query, formatJID)
		if err != nil {
			return nil, err
		}
		out.Users = append(out.Users, user)
	}

	return out, nil
}

func (m *ClientManager) checkOneUser(ctx context.Context, cli *whatsmeow.Client, query string, formatJID bool) (CheckUser, error) {
	result := CheckUser{
		Query:     query,
		RemoteJID: query,
	}

	candidates := checkUserCandidates(query, formatJID)
	var lastErr error
	for _, candidate := range candidates {
		resp, err := cli.IsOnWhatsApp(ctx, []string{candidate})
		if err != nil {
			lastErr = fmt.Errorf("failed to check if number %s exists on WhatsApp: %w", candidate, err)
			continue
		}
		if len(resp) == 0 {
			lastErr = fmt.Errorf("number %s not found in WhatsApp response", candidate)
			continue
		}

		info := resp[0]
		if !info.IsIn {
			continue
		}
		if info.JID.User == "" || info.JID.Server == "" {
			continue
		}

		jid := info.JID.String()
		result.IsInWhatsapp = true
		result.JID = jid
		result.RemoteJID = jid
		result.VerifiedName = verifiedNameString(info.VerifiedName)
		result.LID = lookupLID(ctx, cli, info.JID)

		cacheKey := normalizeE164Local(query)
		if cacheKey == "" {
			cacheKey = query
		}
		if m != nil && m.pnCache != nil && info.JID.Server == types.DefaultUserServer {
			m.pnCache[cacheKey] = pnCacheItem{val: jid, exp: time.Now().Add(24 * time.Hour)}
		}
		return result, nil
	}

	if lastErr != nil {
		return CheckUser{}, lastErr
	}
	return result, nil
}

func lookupLID(ctx context.Context, cli *whatsmeow.Client, pn types.JID) *string {
	if cli == nil || cli.Store == nil || cli.Store.LIDs == nil || pn.IsEmpty() {
		return nil
	}
	lid, err := cli.Store.LIDs.GetLIDForPN(ctx, pn)
	if err != nil || lid.IsEmpty() {
		return nil
	}
	lidStr := lid.String()
	return &lidStr
}

func checkUserCandidates(query string, formatJID bool) []string {
	if jid, err := types.ParseJID(strings.TrimSpace(query)); err == nil && jid.User != "" {
		return []string{jid.User}
	}
	var out []string
	seen := map[string]bool{}
	add := func(values ...string) {
		for _, value := range values {
			value = strings.TrimSpace(value)
			if value == "" || seen[value] {
				continue
			}
			seen[value] = true
			out = append(out, value)
		}
	}
	if formatJID {
		if digits := phoneNumberToJIDDigits(query); digits != "" {
			if strings.HasPrefix(digits, "55") || strings.HasPrefix(strings.TrimSpace(query), "+") {
				add(candidatesBR("+" + digits)...)
			}
		}
	}
	e164 := normalizeE164Local(query)
	if e164 == "" {
		e164 = strings.TrimSpace(query)
	}
	add(candidatesBR(e164)...)
	return out
}

func verifiedNameString(vn *types.VerifiedName) string {
	if vn == nil || vn.Details == nil {
		return ""
	}
	return vn.Details.GetVerifiedName()
}
