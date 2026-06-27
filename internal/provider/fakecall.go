package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// FakeCall envia um "offer" de chamada (somente para tocar) e encerra após ringMS.
// Retorna o call-id usado.
func (m *ClientManager) FakeCall(sessionID, to string, ringMS int, deviceMode string) (string, error) {
	ent, err := m.mustHave(sessionID)
	if err != nil {
		return "", err
	}
	cli := ent.Client
	if cli == nil || !cli.IsConnected() || !cli.IsLoggedIn() {
		return "", fmt.Errorf("session not connected")
	}

	// Resolve destino usando a mesma lógica das mensagens
	res, err := m.ResolveDest(sessionID, to)
	if err != nil {
		return "", err
	}
	destJID, err := types.ParseJID(res.DestJID)
	if err != nil {
		return "", fmt.Errorf("invalid destination jid: %w", err)
	}
	// WhatsApp usa JID não-AD para eventos de call
	dest := destJID.ToNonAD()
	own := cli.DangerousInternals().GetOwnID()
	deviceMode = strings.ToLower(strings.TrimSpace(deviceMode))
	if deviceMode == "" {
		deviceMode = "direct"
	}
	termDest := dest
	termOwn := own
	if deviceMode == "direct" {
		termDest = cli.ResolveCallRecipient(context.Background(), dest)
		if ownLID := cli.DangerousInternals().GetOwnLID().ToNonAD(); !ownLID.IsEmpty() {
			termOwn = ownLID
		} else {
			termOwn = own.ToNonAD()
		}
	}

	callID, err := cli.OfferCallWithDeviceMode(context.Background(), dest, false, deviceMode)
	if err != nil {
		return "", fmt.Errorf("send offer failed: %w", err)
	}

	// Agenda o término após ringMS
	if ringMS <= 0 {
		ringMS = 15000
	}
	go func() {
		timer := time.NewTimer(time.Duration(ringMS) * time.Millisecond)
		defer timer.Stop()
		<-timer.C
		// Envia terminate para parar o toque
		term := waBinary.Node{
			Tag:   "terminate",
			Attrs: waBinary.Attrs{"call-id": callID, "call-creator": termOwn, "reason": "timeout"},
		}
		tnode := waBinary.Node{
			Tag:     "call",
			Attrs:   waBinary.Attrs{"id": cli.GenerateMessageID(), "from": termOwn, "to": termDest},
			Content: []waBinary.Node{term},
		}
		if err := cli.DangerousInternals().SendNode(context.Background(), tnode); err != nil {
			ent.Log.Errorf("fake_call terminate error to=%s wire_to=%s call_id=%s err=%v", dest.String(), termDest.String(), callID, err)
		}
	}()

	return callID, nil
}
