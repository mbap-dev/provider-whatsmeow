package provider

import (
	"fmt"
	"time"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
)

// FakeCall envia um "offer" de chamada (somente para tocar) e encerra após ringMS.
// Retorna o call-id usado.
func (m *ClientManager) FakeCall(sessionID, to string, ringMS int) (string, error) {
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
	own := cli.DangerousInternals().GetOwnID().ToNonAD()

	// IDs
	stanzaID := cli.GenerateMessageID()
	callID := string(cli.GenerateMessageID())

	// Envia nó de oferta (ring)
	offer := waBinary.Node{
		Tag:   "offer",
		Attrs: waBinary.Attrs{"call-id": callID, "call-creator": own},
	}
	node := waBinary.Node{
		Tag:   "call",
		Attrs: waBinary.Attrs{"id": stanzaID, "from": own, "to": dest, "platform": "web", "version": "2"},
		Content: []waBinary.Node{
			offer,
		},
	}
	if err := cli.DangerousInternals().SendNode(node); err != nil {
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
			Attrs: waBinary.Attrs{"call-id": callID, "call-creator": own, "reason": "timeout"},
		}
		tnode := waBinary.Node{
			Tag:     "call",
			Attrs:   waBinary.Attrs{"id": cli.GenerateMessageID(), "from": own, "to": dest},
			Content: []waBinary.Node{term},
		}
		_ = cli.DangerousInternals().SendNode(tnode)
	}()

	ent.Log.Infof("fake_call started to=%s call_id=%s ring_ms=%d", dest.String(), callID, ringMS)
	return callID, nil
}
