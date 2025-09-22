// events.go
package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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
	"google.golang.org/protobuf/proto"
)

// Context key used to carry the wildcard suffix (e.g., phone_number_id)
// extracted from the AMQP routing key into ClientManager.Send.
type ctxKey string

// CtxKeyPhoneNumberID is the key to read/write the phone number id in context.
const CtxKeyPhoneNumberID ctxKey = "phone_number_id"

// =========== Registro de handlers ===========

// cloudMediaID compõe o ID de mídia no formato esperado pelo consumidor:
// "<phone_number_id>/<message_uuid>"
func cloudMediaID(phone, msgID string) string {
	p := strings.ReplaceAll(phone, "+", "")
	if p == "" || msgID == "" {
		return msgID
	}
	return p + "/" + msgID
}

func (m *ClientManager) registerEventHandlers(client *whatsmeow.Client, sessionID string) {
	if client == nil {
		return
	}

	client.AddEventHandler(func(evt interface{}) {
		switch e := evt.(type) {
		case *events.CallOffer:
			// Auto-reject incoming calls and optionally send a reply text
			if m.rejectCalls {
				if err := client.RejectCall(e.CallCreator, e.CallID); err != nil {
					log.WithSession(sessionID).Error("call reject error: %v", err)
				} else {
					log.WithSession(sessionID).Info("evt=call_reject from=%s call_id=%s", e.CallCreator.String(), e.CallID)
				}
				// Optionally send a message to the caller
				msgText := strings.TrimSpace(m.rejectMsg)
				if msgText != "" {
					// Support literal \n already handled in config; just send
					wire := &waE2E.Message{Conversation: proto.String(msgText)}
					// Send to the caller JID (non-AD if applicable)
					to := e.CallCreator
					resp, err := client.SendMessage(context.Background(), to, wire)
					if err != nil {
						log.WithSession(sessionID).Error("call reply send error: %v", err)
					} else {
						// Emit webhook 'message' (outgoing) to UnoAPI
						_ = m.emitCloudOutgoingText(sessionID, to, resp.ID, msgText)
						// Emit webhook 'sent' status to UnoAPI
						_ = m.emitCloudSent(sessionID, to, resp.ID)
						log.WithSession(sessionID).WithMessageID(string(resp.ID)).Info("evt=call_reply_sent to=%s", to.String())
					}
				}
			}
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
	case msg.GetProtocolMessage() != nil:
		pm := msg.GetProtocolMessage()
		if pm.GetType() == waE2E.ProtocolMessage_MESSAGE_EDIT && pm.GetEditedMessage() != nil {
			edited := pm.GetEditedMessage()
			// Build as a text reply to the original message that was edited
			body := strings.TrimSpace(edited.GetConversation())
			if body == "" && edited.GetExtendedTextMessage() != nil {
				body = strings.TrimSpace(edited.GetExtendedTextMessage().GetText())
			}
			if body == "" {
				// fallback to any caption-like text
				body = strings.TrimSpace(extractAnyCaption(edited))
			}
			if body == "" {
				// if we don't have a body, ignore
				return nil
			}
			wireMsg["type"] = "text"
			wireMsg["text"] = map[string]any{"body": body}
			if pm.GetKey() != nil {
				rid := strings.TrimSpace(pm.GetKey().GetId())
				if rid != "" {
					// Log edit detection with snippet
					snippet := body
					if len(snippet) > 120 {
						snippet = snippet[:120]
					}
					log.WithSession(sessionID).WithMessageID(e.Info.ID).
						Info("evt=message_edit reply_to=%s body_snippet=%q", rid, snippet)
					wireMsg["context"] = map[string]any{
						"message_id": rid,
						"id":         rid,
					}
				}
			}
		} else {
			// ignore other protocol messages
			return nil
		}
	case msg.GetContactMessage() != nil:
		c := msg.GetContactMessage()
		wireMsg["type"] = "contacts"
		name, phones := parseVCard(c.GetVcard())
		if name == "" {
			name = strings.TrimSpace(c.GetDisplayName())
		}
		if len(phones) == 0 {
			phones = []map[string]any{}
		}
		entryC := map[string]any{
			"name":   map[string]any{"formatted_name": name},
			"phones": phones,
		}
		wireMsg["contacts"] = []any{entryC}

	case msg.GetContactsArrayMessage() != nil:
		arr := msg.GetContactsArrayMessage()
		wireMsg["type"] = "contacts"
		var list []any
		for _, cm := range arr.GetContacts() {
			cn := strings.TrimSpace(cm.GetDisplayName())
			name, phones := parseVCard(cm.GetVcard())
			if name == "" {
				name = cn
			}
			if len(phones) == 0 {
				phones = []map[string]any{}
			}
			list = append(list, map[string]any{
				"name":   map[string]any{"formatted_name": name},
				"phones": phones,
			})
		}
		if len(list) == 0 {
			return nil
		}
		wireMsg["contacts"] = list

	case msg.GetReactionMessage() != nil:
		r := msg.GetReactionMessage()
		emoji := strings.TrimSpace(r.GetText())
		if emoji == "" {
			// ignore empty reaction (removal)
			return nil
		}
		wireMsg["type"] = "text"
		wireMsg["text"] = map[string]any{"body": emoji}
		// reply to the original message that was reacted
		if rk := r.GetKey(); rk != nil {
			rid := strings.TrimSpace(rk.GetId())
			if rid != "" {
				wireMsg["context"] = map[string]any{
					"message_id": rid,
					"id":         rid,
				}
			}
		}
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
		imgExt := extensionByMime(mimeType)
		image := map[string]any{
			"caption":     im.GetCaption(),
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(im.GetFileSHA256()),
			"media_key":   b64(im.GetMediaKey()),
			"direct_path": im.GetDirectPath(),
		}
		if strings.TrimSpace(imgExt) != "" {
			image["filename"] = e.Info.ID + imgExt
		}
		if im.FileLength != nil {
			image["file_length"] = im.GetFileLength()
		}
		wireMsg["image"] = image

	case msg.GetDocumentMessage() != nil:
		d := msg.GetDocumentMessage()
		wireMsg["type"] = "document"
		mimeType := d.GetMimetype()
		document := map[string]any{
			"caption":     d.GetCaption(),
			"filename":    firstNonEmpty(d.GetFileName(), d.GetTitle()),
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(d.GetFileSHA256()),
			"media_key":   b64(d.GetMediaKey()),
			"direct_path": d.GetDirectPath(),
		}
		if d.FileLength != nil {
			document["file_length"] = d.GetFileLength()
		}
		wireMsg["document"] = document

	case msg.GetVideoMessage() != nil:
		v := msg.GetVideoMessage()
		wireMsg["type"] = "video"
		mimeType := v.GetMimetype()
		video := map[string]any{
			"caption":     v.GetCaption(),
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(v.GetFileSHA256()),
			"media_key":   b64(v.GetMediaKey()),
			"direct_path": v.GetDirectPath(),
		}
		if v.FileLength != nil {
			video["file_length"] = v.GetFileLength()
		}
		wireMsg["video"] = video

	case msg.GetAudioMessage() != nil:
		a := msg.GetAudioMessage()
		wireMsg["type"] = "audio"
		mimeType := a.GetMimetype()
		audio := map[string]any{
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(a.GetFileSHA256()),
			"media_key":   b64(a.GetMediaKey()),
			"direct_path": a.GetDirectPath(),
		}
		if a.FileLength != nil {
			audio["file_length"] = a.GetFileLength()
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
		ext := extensionByMime(mimeType)
		filename := e.Info.ID + ext
		sticker := map[string]any{
			"filename":    filename,
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(s.GetFileSHA256()),
			"media_key":   b64(s.GetMediaKey()),
			"direct_path": s.GetDirectPath(),
		}
		if s.FileLength != nil {
			sticker["file_length"] = s.GetFileLength()
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
	// For group messages, set profile.name to the sender phone
	profileName := contactPhone
	contactObj := map[string]any{
		"profile": map[string]any{"name": profileName},
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

// emitCloudOutgoingText publishes an outgoing text message event in the same
// Cloud-like format used for incoming messages, so downstreams can see the
// content that was sent automatically (e.g., call rejection reply).
func (m *ClientManager) emitCloudOutgoingText(sessionID string, to types.JID, id types.MessageID, body string) error {
	phone := normalizePhone(sessionID)
	recipient := jidToPhoneNumberIfUser(to)

	// Build message payload (from = our phone)
	wireMsg := map[string]any{
		"from":      strings.ReplaceAll(phone, "+", ""),
		"id":        id,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"type":      "text",
		"text":      map[string]any{"body": body},
	}

	// Contact entry for the recipient (group or user)
	contactPhone := recipient
	contactObj := map[string]any{
		"profile": map[string]any{"name": contactPhone},
		"wa_id":   contactPhone,
	}
	if isGroupJID(to) {
		contactObj["group_id"] = to.String()
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["contacts"] = []any{contactObj}
	val["messages"] = []any{wireMsg}

	log.WithSession(sessionID).WithMessageID(string(id)).
		Info("evt=message_out cloud payload ready type=text to=%s", to.String())

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

func extensionByMime(m string) string {
	exts, _ := mime.ExtensionsByType(splitMime(m))
	if len(exts) > 0 {
		return exts[0]
	}
	return ""
}

// parseVCard extracts formatted name and phone entries from a vCard string.
// It supports lines like:
// N:John Doe
// TEL;type=CELL;type=VOICE;waid=5511999999999:5511999999999
// TEL;type=CELL;type=VOICE:+55 62 93300-0233
func parseVCard(v string) (string, []map[string]any) {
	if strings.TrimSpace(v) == "" {
		return "", nil
	}
	var fn string
	var nFormatted string
	var phones []map[string]any
	lines := strings.Split(v, "\n")
	for _, ln := range lines {
		s := strings.TrimSpace(ln)
		if s == "" {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(s), "N:") {
			raw := strings.TrimSpace(s[2:])
			// N: Family;Given;Additional;Prefix;Suffix
			parts := strings.Split(raw, ";")
			// Prefer Given + Family if present
			var given, family string
			if len(parts) > 1 {
				family = strings.TrimSpace(parts[0])
				given = strings.TrimSpace(parts[1])
			} else if len(parts) == 1 {
				given = strings.TrimSpace(parts[0])
			}
			if given != "" && family != "" {
				nFormatted = strings.TrimSpace(given + " " + family)
			} else if given != "" {
				nFormatted = given
			} else if family != "" {
				nFormatted = family
			} else {
				nFormatted = raw
			}
			continue
		}
		if strings.HasPrefix(strings.ToUpper(s), "FN:") {
			fn = strings.TrimSpace(s[3:])
			continue
		}
		if strings.HasPrefix(strings.ToUpper(s), "TEL") {
			// phone is after last ':'
			c := strings.LastIndex(s, ":")
			if c <= 0 || c+1 >= len(s) {
				continue
			}
			phone := strings.TrimSpace(s[c+1:])
			waid := ""
			if i := strings.Index(strings.ToLower(s), "waid="); i >= 0 {
				// take until next ';' or ':'
				rest := s[i+5:]
				j := strings.IndexAny(rest, ";:")
				if j >= 0 {
					waid = strings.TrimSpace(rest[:j])
				} else {
					waid = strings.TrimSpace(rest)
				}
			}
			entry := map[string]any{"phone": phone}
			if waid != "" {
				entry["wa_id"] = waid
			}
			phones = append(phones, entry)
		}
	}
	// prefer FN if present, else formatted N
	name := strings.TrimSpace(fn)
	if name == "" {
		name = strings.TrimSpace(nFormatted)
	}
	return name, phones
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
