package provider

import (
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
)

type qrState struct {
	pngBase64 string
	updatedAt time.Time
}

var (
	qrMu        sync.RWMutex
	qrBySession = map[string]*qrState{}
)

// startQRWatcher deve ser chamado ANTES de Connect(); ele escuta o canal de QR
// do whatsmeow e atualiza o último QR para a sessão.
func (m *ClientManager) startQRWatcher(sessionID string, ch <-chan whatsmeow.QRChannelItem) {
	go func() {
		for item := range ch {
			switch item.Event {
			case "code":
				// item.Code == string crua do QR (não é PNG)
				png, _ := qrcode.Encode(item.Code, qrcode.Medium, 256)
				b64 := base64.StdEncoding.EncodeToString(png)

				qrMu.Lock()
				qrBySession[sessionID] = &qrState{
					pngBase64: b64,
					updatedAt: time.Now(),
				}
				qrMu.Unlock()

			case "success", "timeout":
				// Limpa no timeout ou depois do sucesso (opcional)
				qrMu.Lock()
				delete(qrBySession, sessionID)
				qrMu.Unlock()

			case "error":
				// Você pode logar item.Error se quiser
			}
		}
	}()
}

// GetQR retorna o PNG/base64 do último QR da sessão.
func (m *ClientManager) GetQR(sessionID string) (string, error) {
	qrMu.RLock()
	state, ok := qrBySession[sessionID]
	qrMu.RUnlock()
	if !ok || state == nil || state.pngBase64 == "" {
		return "", errors.New("qr not ready or session not connected")
	}
	return state.pngBase64, nil
}
