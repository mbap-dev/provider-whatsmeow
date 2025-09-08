// events.go
package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"your.org/provider-whatsmeow/internal/log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	_ "google.golang.org/protobuf/proto"
)

// Context key used to carry the wildcard suffix (e.g., phone_number_id)
// extracted from the AMQP routing key into ClientManager.Send.
type ctxKey string

// CtxKeyPhoneNumberID is the key to read/write the phone number id in context.
const CtxKeyPhoneNumberID ctxKey = "phone_number_id"

// =========== Registro de handlers ===========

func (m *ClientManager) registerEventHandlers(client *whatsmeow.Client, sessionID string) {
	if client == nil {
		return
	}

	client.AddEventHandler(func(evt interface{}) {
		switch e := evt.(type) {
		case *events.Message:
			if err := m.emitCloudMessage(sessionID, client, e); err != nil {
				log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("webhook cloud message error: %v", err)
			}
		case *events.Receipt:
			if err := m.emitCloudReceipt(sessionID, client, e); err != nil {
				log.WithSession(sessionID).Error("webhook cloud receipt error: %v", err)
			}
		default:
			// Silencia outros eventos no webhook (Uno espera apenas messages/statuses)
		}
	})
}

// =========== Emissão de mensagens no formato Cloud ===========

func (m *ClientManager) emitCloudMessage(sessionID string, client *whatsmeow.Client, e *events.Message) error {
	phone := normalizePhone(sessionID) // Uno usa “phone” como id/metadata.* sem "+"
	msg := e.Message
	if msg == nil {
		return nil
	}

	// Calcula contrapartes
	fromMe := e.Info.IsFromMe
	chatJID := e.Info.Chat
	senderJID := e.Info.Sender

	contactPhone := jidToPhoneNumberIfUser(chatJID)    // quem aparece em contacts.wa_id
	if isGroupJID(chatJID) && (senderJID.User != "") { // em grupos, “from” é quem falou
		contactPhone = jidToPhoneNumberIfUser(senderJID) // ainda mantemos contacts.wa_id do chat; group_id vai junto
	}
	fromField := phone
	if !fromMe {
		fromField = jidToPhoneNumberIfUser(senderJID)
	}

	// Monta “message” no padrão Cloud
	wireMsg := map[string]any{
		"from":      strings.ReplaceAll(fromField, "+", ""),
		"id":        e.Info.ID,
		"timestamp": strconv.FormatInt(e.Info.Timestamp.Unix(), 10),
	}

	// Context (reply/quote) se houver
	if ctx := messageContextInfo(msg); ctx != nil {
		wireMsg["context"] = ctx
	}

	// Tipo de conteúdo
	switch {
	case msg.GetConversation() != "":
		wireMsg["type"] = "text"
		wireMsg["text"] = map[string]any{"body": msg.GetConversation()}

	case msg.GetExtendedTextMessage() != nil:
		t := strings.TrimSpace(msg.GetExtendedTextMessage().GetText())
		wireMsg["type"] = "text"
		wireMsg["text"] = map[string]any{"body": t}

	case msg.GetImageMessage() != nil:
		im := msg.GetImageMessage()
		wireMsg["type"] = "image"
		mimeType := im.GetMimetype()
		objName := mediaKey(phone, e.Info.ID) + extensionByMime(mimeType)
		image := map[string]any{
			"caption":   im.GetCaption(),
			"mime_type": splitMime(mimeType),
			"sha256":    b64(im.GetFileSHA256()),
			"id":        e.Info.ID,
		}
		if _, err := m.storeMedia(context.Background(), client, msg, objName, mimeType); err != nil {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("media upload error: %v", err)
		}
		wireMsg["image"] = image

	case msg.GetDocumentMessage() != nil:
		d := msg.GetDocumentMessage()
		wireMsg["type"] = "document"
		mimeType := d.GetMimetype()
		objName := mediaKey(phone, e.Info.ID) + extensionByMime(mimeType)
		document := map[string]any{
			"caption":   d.GetCaption(),
			"filename":  firstNonEmpty(d.GetFileName(), d.GetTitle()),
			"mime_type": splitMime(mimeType),
			"sha256":    b64(d.GetFileSHA256()),
			"id":        e.Info.ID,
		}
		if _, err := m.storeMedia(context.Background(), client, msg, objName, mimeType); err != nil {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("media upload error: %v", err)
		}
		wireMsg["document"] = document

	case msg.GetVideoMessage() != nil:
		v := msg.GetVideoMessage()
		wireMsg["type"] = "video"
		mimeType := v.GetMimetype()
		objName := mediaKey(phone, e.Info.ID) + extensionByMime(mimeType)
		video := map[string]any{
			"caption":   v.GetCaption(),
			"mime_type": splitMime(mimeType),
			"sha256":    b64(v.GetFileSHA256()),
			"id":        e.Info.ID,
		}
		if _, err := m.storeMedia(context.Background(), client, msg, objName, mimeType); err != nil {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("media upload error: %v", err)
		}
		wireMsg["video"] = video

	case msg.GetAudioMessage() != nil:
		a := msg.GetAudioMessage()
		wireMsg["type"] = "audio"
		mimeType := a.GetMimetype()
		objName := mediaKey(phone, e.Info.ID) + extensionByMime(mimeType)
		audio := map[string]any{
			"mime_type": splitMime(mimeType),
			"sha256":    b64(a.GetFileSHA256()),
			"id":        e.Info.ID,
		}
		if _, err := m.storeMedia(context.Background(), client, msg, objName, mimeType); err != nil {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("media upload error: %v", err)
		}
		if a.Seconds != nil {
			audio["seconds"] = a.GetSeconds()
		}
		if a.PTT != nil {
			audio["ptt"] = a.GetPTT()
		}
		// inclui waveform se existir
		if len(a.GetWaveform()) > 0 {
			audio["waveform"] = base64.StdEncoding.EncodeToString(a.GetWaveform())
		}
		wireMsg["audio"] = audio

	case msg.GetStickerMessage() != nil:
		s := msg.GetStickerMessage()
		wireMsg["type"] = "sticker"
		mimeType := s.GetMimetype()
		objName := mediaKey(phone, e.Info.ID) + extensionByMime(mimeType)
		sticker := map[string]any{
			"mime_type": splitMime(mimeType),
			"sha256":    b64(s.GetFileSHA256()),
			"id":        e.Info.ID,
		}
		if _, err := m.storeMedia(context.Background(), client, msg, objName, mimeType); err != nil {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("media upload error: %v", err)
		}
		wireMsg["sticker"] = sticker

	case msg.GetLocationMessage() != nil:
		l := msg.GetLocationMessage()
		wireMsg["type"] = "location"
		wireMsg["location"] = map[string]any{
			"latitude":  l.GetDegreesLatitude(),
			"longitude": l.GetDegreesLongitude(),
		}

	default:
		// fallback: se nada identificado mas tem algo, tenta extrair caption/text
		if cap := extractAnyCaption(msg); cap != "" {
			wireMsg["type"] = "text"
			wireMsg["text"] = map[string]any{"body": cap}
		} else {
			// ignora tipos não suportados
			return nil
		}
	}

	// Contatos (inclui group_id quando for grupo)
	contactObj := map[string]any{
		"profile": map[string]any{"name": contactPhone},
		"wa_id":   contactPhone,
	}
	if isGroupJID(chatJID) {
		contactObj["group_id"] = chatJID.String()
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["contacts"] = []any{contactObj}
	val["messages"] = []any{wireMsg}

	log.WithSession(sessionID).WithMessageID(e.Info.ID).
		Info("evt=message cloud payload ready type=%s", wireMsg["type"])

	return m.deliverWebhook(sessionID, payload)
}

// =========== Emissão de recibos no formato Cloud (statuses) ===========

func (m *ClientManager) emitCloudReceipt(sessionID string, client *whatsmeow.Client, e *events.Receipt) error {
	if len(e.MessageIDs) == 0 {
		return nil
	}
	phone := normalizePhone(sessionID)

	status := mapReceiptStatus(e.Type) // sent|delivered|read|played|deleted
	if status == "" {
		return nil
	}

	recipient := "" // melhor esforço: destinatário do chat (user)
	if e.Chat != (types.JID{}) {
		recipient = jidToPhoneNumberIfUser(e.Chat)
	}
	if recipient == "" && e.Sender != (types.JID{}) {
		recipient = jidToPhoneNumberIfUser(e.Sender)
	}

	states := make([]any, 0, len(e.MessageIDs))
	for _, id := range e.MessageIDs {
		st := map[string]any{
			"id":           id,
			"recipient_id": strings.ReplaceAll(recipient, "+", ""),
			"status":       status,
		}
		if !e.Timestamp.IsZero() {
			st["timestamp"] = strconv.FormatInt(e.Timestamp.Unix(), 10)
		}
		// conversation.id ajuda o Uno a agrupar
		if e.Chat != (types.JID{}) {
			st["conversation"] = map[string]any{"id": e.Chat.String()}
		}
		states = append(states, st)
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["statuses"] = states

	log.WithSession(sessionID).Info("evt=receipt cloud payload ready type=%s ids=%v", e.Type, e.MessageIDs)

	return m.deliverWebhook(sessionID, payload)
}

func (m *ClientManager) emitCloudSent(sessionID string, to types.JID, id types.MessageID) error {
	phone := normalizePhone(sessionID)
	recipient := jidToPhoneNumberIfUser(to)
	st := map[string]any{
		"id":           id,
		"recipient_id": strings.ReplaceAll(recipient, "+", ""),
		"status":       "sent",
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
		"conversation": map[string]any{"id": to.String()},
	}
	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["statuses"] = []any{st}
	log.WithSession(sessionID).WithMessageID(id).Info("evt=sent cloud payload ready id=%s", id)
	return m.deliverWebhook(sessionID, payload)
}

// =========== Helpers de mapeamento Cloud ===========

func cloudEnvelope(phone string) map[string]any {
	phone = strings.ReplaceAll(phone, "+", "")
	val := map[string]any{
		"messaging_product": "whatsapp",
		"metadata": map[string]any{
			"display_phone_number": phone,
			"phone_number_id":      phone,
		},
		"messages": []any{},
		"contacts": []any{},
		"statuses": []any{},
		"errors":   []any{},
	}
	change := map[string]any{
		"field": "messages",
		"value": val,
	}
	entry := map[string]any{
		"id":      phone,
		"changes": []any{change},
	}
	return map[string]any{
		"object": "whatsapp_business_account",
		"entry":  []any{entry},
	}
}

func normalizePhone(s string) string {
	return digitsOnly(s)
}

func jidToPhoneNumberIfUser(j any) string {
	var jid types.JID
	switch v := j.(type) {
	case types.JID:
		jid = v
	case string:
		jid, _ = types.ParseJID(v)
	default:
		return ""
	}
	// para usuários: “user@server” => “+<digits>” sem domínio
	if isGroupJID(jid) {
		return jid.User // grupo: devolvemos só o user (id do grupo) — Uno lida com group_id separado
	}
	return digitsOnly(jid.User)
}

func isGroupJID(j types.JID) bool {
	return strings.HasSuffix(j.Server, "g.us")
}

func splitMime(m string) string {
	if m == "" {
		return ""
	}
	return strings.SplitN(m, ";", 2)[0]
}

func b64(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func mediaKey(phone, waMsgID string) string {
	return fmt.Sprintf("%s/%s", strings.ReplaceAll(phone, "+", ""), waMsgID)
}

func extensionByMime(m string) string {
	exts, _ := mime.ExtensionsByType(splitMime(m))
	if len(exts) > 0 {
		return exts[0]
	}
	return ""
}

func (m *ClientManager) storeMedia(ctx context.Context, cli *whatsmeow.Client, msg *waE2E.Message, object, mimeType string) (string, error) {
	if m.storage == nil {
		return "", fmt.Errorf("storage not configured")
	}
	data, err := cli.DownloadAny(ctx, msg)
	if err != nil {
		return "", err
	}
	return m.storage.Upload(ctx, object, data, mimeType)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func extractAnyCaption(m *waE2E.Message) string {
	switch {
	case m.GetImageMessage() != nil:
		return m.GetImageMessage().GetCaption()
	case m.GetVideoMessage() != nil:
		return m.GetVideoMessage().GetCaption()
	case m.GetDocumentMessage() != nil:
		return m.GetDocumentMessage().GetCaption()
	}
	return ""
}

func messageContextInfo(m *waE2E.Message) map[string]any {
	// tenta extrair stanzaId e participant de qualquer ContextInfo disponível
	var ci *waE2E.ContextInfo
	switch {
	case m.GetExtendedTextMessage() != nil:
		ci = m.GetExtendedTextMessage().GetContextInfo()
	case m.GetImageMessage() != nil:
		ci = m.GetImageMessage().GetContextInfo()
	case m.GetDocumentMessage() != nil:
		ci = m.GetDocumentMessage().GetContextInfo()
	case m.GetVideoMessage() != nil:
		ci = m.GetVideoMessage().GetContextInfo()
	case m.GetAudioMessage() != nil:
		ci = m.GetAudioMessage().GetContextInfo()
	case m.GetStickerMessage() != nil:
		ci = m.GetStickerMessage().GetContextInfo()
	}
	if ci == nil {
		return nil
	}
	stanzaID := strings.TrimSpace(ci.GetStanzaID())
	if stanzaID == "" && ci.GetQuotedMessage() != nil {
		// se veio quotedMessage mas sem stanzaId, ainda assim envia context pra compat
		return map[string]any{
			"message_id": stanzaID,
			"id":         stanzaID,
		}
	}
	if stanzaID == "" {
		return nil
	}
	return map[string]any{
		"message_id": stanzaID,
		"id":         stanzaID,
	}
}

func mapReceiptStatus(t events.ReceiptType) string {
	switch t {
	case events.ReceiptTypeDelivered:
		return "delivered"
	case events.ReceiptTypeRead:
		return "read"
	case events.ReceiptTypePlayed:
		return "played"
	case events.ReceiptTypeSender:
		return "sent"
	default:
		return ""
	}
}

// =========== Entrega HTTP (mantida) ===========

func (m *ClientManager) deliverWebhook(sessionID string, payload any) error {
	entry := log.WithSession(sessionID)

	if m.webhookBase == "" {
		entry.Info("webhook disabled (WEBHOOK_BASE empty), dropping")
		return nil
	}

	url := strings.TrimRight(m.webhookBase, "/") + "/" + sessionID

	body, err := json.Marshal(payload)
	if err != nil {
		entry.Error("webhook marshal error: %v", err)
		return err
	}

	entry.Info("webhook post start url=%s bytes=%d", url, len(body))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		entry.Error("webhook build request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	go func() {
		resp, err := http.DefaultClient.Do(req.WithContext(context.Background()))
		if err != nil {
			entry.Error("webhook post error: %v", err)
			return
		}
		defer resp.Body.Close()

		snippet := ""
		if b, _ := io.ReadAll(io.LimitReader(resp.Body, 512)); len(b) > 0 {
			snippet = string(b)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			entry.Info("webhook post done status=%d", resp.StatusCode)
		} else {
			entry.Error("webhook post non-2xx status=%d body=%q", resp.StatusCode, snippet)
		}
	}()

	return nil
}
