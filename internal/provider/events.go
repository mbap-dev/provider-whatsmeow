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
	"your.org/provider-whatsmeow/internal/status"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Defaults used for LID->AD resolution when no manager-specific tuning is available.
var (
	defaultLIDResolveAttempts = 3
	defaultLIDResolveDelay    = 300 * time.Millisecond
)

// Shared HTTP client for webhook delivery with tuned connection pooling.
var webhookHTTPClient = &http.Client{
	Timeout: 20 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

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
		// Connection lifecycle → Redis status
		case *events.Connected:
			status.Set(sessionID, "online")
			// Start always-online presence refresher if enabled
			m.maybeStartAlwaysOnline(sessionID, client)
		case *events.Disconnected:
			status.Set(sessionID, "offline")
			m.stopAlwaysOnline(sessionID)
		case *events.LoggedOut:
			status.Set(sessionID, "disconnected")
			m.stopAlwaysOnline(sessionID)
		case *events.StreamReplaced:
			status.Set(sessionID, "restart_required")
			m.stopAlwaysOnline(sessionID)
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
		case *events.AppStateSyncComplete:
			if e.Name == "critical_block" || e.Name == "regular_high" {
				_ = client.SendPresence(types.PresenceAvailable)
			}
		case *events.PushNameSetting:
			_ = client.SendPresence(types.PresenceAvailable)
		case *events.Message:
			if err := m.emitCloudMessage(sessionID, client, e); err != nil {
				log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("webhook cloud message error: %v", err)
			}
			if m.autoMarkRead {
				go m.markReadEvent(sessionID, client, e)
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

	fromMe := e.Info.IsFromMe
	chatJID := e.Info.Chat
	senderJID := e.Info.Sender
	// Preferir JID alternativo (AD) quando a lib expõe SenderAlt não‑LID
	if e.Info.SenderAlt != (types.JID{}) && !isLIDJID(e.Info.SenderAlt) && e.Info.SenderAlt.Server != "" {
		senderJID = e.Info.SenderAlt
	}

	// Ignore optional channels per configuration
	if m.ignoreStatusBroadcast && (isStatusBroadcastJID(chatJID) || isStatusBroadcastJID(senderJID)) {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Info("evt=message skip: status_broadcast")
		return nil
	}
	if m.ignoreNewsletters && (isNewsletterJID(chatJID) || isNewsletterJID(senderJID)) {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Info("evt=message skip: newsletter")
		return nil
	}

	// Resolve contato e "from" considerando JIDs LID (memoize para evitar resolver duas vezes)
	resolvedMemo := map[string]string{}
	resolveDisplay := func(j types.JID) string {
		key := j.String()
		if v, ok := resolvedMemo[key]; ok {
			return v
		}
		v := jidDisplayMaybeResolve(sessionID, client, j)
		resolvedMemo[key] = v
		return v
	}
	contactPhone := resolveDisplay(chatJID)
	if isGroupJID(chatJID) && (senderJID.User != "") {
		contactPhone = resolveDisplay(senderJID)
	}
	fromField := phone
	if !fromMe {
		fromField = resolveDisplay(senderJID)
	}

	// Monta “message” no padrão Cloud
	wireMsg := map[string]any{
		"from":      normalizeFromField(fromField),
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
			body = body + "\n`Mensagem Editada`"
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
		} else if pm.GetType() == waE2E.ProtocolMessage_REVOKE {
			// Message deletion (revoke). Emit a status 'deleted' similar to UnoAPI.
			rid := ""
			if pm.GetKey() != nil {
				rid = strings.TrimSpace(pm.GetKey().GetId())
			}
			if rid == "" {
				// Without the original id, do nothing
				return nil
			}
			if err := m.emitCloudDeleted(sessionID, client, e, rid); err != nil {
				log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("webhook cloud deleted error: %v", err)
			}
			return nil
		} else {
			// ignore other protocol messages
			return nil
		}
	case msg.GetContactMessage() != nil:
		c := msg.GetContactMessage()
		wireMsg["type"] = "contacts"
		// Debug raw inbound vCard to help diagnose WA Business contacts
		{
			entry := log.WithSession(sessionID).WithMessageID(e.Info.ID)
			v := c.GetVcard()
			snippet := v
			if len(snippet) > 2000 {
				snippet = snippet[:2000] + "...(truncated)"
			}
			entry.Debug("in_contact_raw display_name=%q vcard_len=%d vcard=%q", strings.TrimSpace(c.GetDisplayName()), len(v), snippet)
		}

		name, phones := parseVCard(c.GetVcard())
		if name == "" {
			name = strings.TrimSpace(c.GetDisplayName())
		}
		if len(phones) == 0 {
			phones = []map[string]any{}
		}
		{
			entry := log.WithSession(sessionID).WithMessageID(e.Info.ID)
			entry.Debug("in_contact_parsed name=%q phones=%v", name, phones)
		}
		entryC := map[string]any{
			"name":   map[string]any{"formatted_name": name},
			"phones": phones,
		}
		wireMsg["contacts"] = []any{entryC}
		{
			entry := log.WithSession(sessionID).WithMessageID(e.Info.ID)
			entry.Debug("in_contact_wire contacts=%v", wireMsg["contacts"])
		}

	case msg.GetContactsArrayMessage() != nil:
		arr := msg.GetContactsArrayMessage()
		wireMsg["type"] = "contacts"
		var list []any
		for i, cm := range arr.GetContacts() {
			// Debug raw for each entry
			{
				entry := log.WithSession(sessionID).WithMessageID(e.Info.ID)
				v := cm.GetVcard()
				snippet := v
				if len(snippet) > 2000 {
					snippet = snippet[:2000] + "...(truncated)"
				}
				entry.Debug("in_contacts_raw[%d] display_name=%q vcard_len=%d vcard=%q", i, strings.TrimSpace(cm.GetDisplayName()), len(v), snippet)
			}
			cn := strings.TrimSpace(cm.GetDisplayName())
			name, phones := parseVCard(cm.GetVcard())
			if name == "" {
				name = cn
			}
			if len(phones) == 0 {
				phones = []map[string]any{}
			}
			{
				entry := log.WithSession(sessionID).WithMessageID(e.Info.ID)
				entry.Debug("in_contacts_parsed[%d] name=%q phones=%v", i, name, phones)
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
		{
			entry := log.WithSession(sessionID).WithMessageID(e.Info.ID)
			entry.Debug("in_contacts_wire contacts=%v", wireMsg["contacts"])
		}

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

	// Contato com nome amigável:
	// - Grupo: usar nome do grupo
	// - 1:1: se a mensagem é minha (fromMe), não usar PushName; usar o número (contactPhone)
	//         se não é minha, usar PushName (fallback: número)
	var dispName string
	if isGroupJID(chatJID) {
		if gname, ok := safeGetGroupName(client, chatJID); ok && strings.TrimSpace(gname) != "" {
			dispName = gname
		} else {
			dispName = chatJID.User
		}
	} else {
		if fromMe {
			dispName = contactPhone
		} else {
			dispName = strings.TrimSpace(e.Info.PushName)
			if dispName == "" {
				dispName = contactPhone
			}
		}
	}
	// Monta profile com name e, se disponível, avatar (Cloud-like)
	profile := map[string]any{"name": dispName}
	if av, ok := getAvatarURL(client, func() types.JID {
		if isGroupJID(chatJID) {
			return chatJID
		}
		return chatJID
	}()); ok && strings.TrimSpace(av) != "" {
		profile["picture"] = av
	}
	contactObj := map[string]any{
		"profile": profile,
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

// emitCloudDeleted publishes a status event for a message deletion (revoke)
// in the Cloud-like format UnoAPI expects.
func (m *ClientManager) emitCloudDeleted(sessionID string, client *whatsmeow.Client, e *events.Message, revokedID string) error {
	phone := normalizePhone(sessionID)

	// Best-effort recipient id derived from chat/sender
	recipient := ""
	if e.Info.Chat != (types.JID{}) {
		recipient = jidToPhoneNumberIfUser(e.Info.Chat)
	}
	if recipient == "" && e.Info.Sender != (types.JID{}) {
		recipient = jidToPhoneNumberIfUser(e.Info.Sender)
	}

	st := map[string]any{
		"id":           revokedID,
		"recipient_id": strings.ReplaceAll(recipient, "+", ""),
		"status":       "deleted",
	}
	if !e.Info.Timestamp.IsZero() {
		st["timestamp"] = strconv.FormatInt(e.Info.Timestamp.Unix(), 10)
	}
	if e.Info.Chat != (types.JID{}) {
		st["conversation"] = map[string]any{"id": e.Info.Chat.String()}
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["statuses"] = []any{st}

	log.WithSession(sessionID).WithMessageID(e.Info.ID).Info("evt=deleted cloud payload ready revoked_id=%s", revokedID)
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

func isStatusBroadcastJID(j types.JID) bool {
	return j.Server == types.BroadcastServer && j.User == types.StatusBroadcastJID.User
}

func isNewsletterJID(j types.JID) bool {
	return j.Server == types.NewsletterServer
}

// isLIDJID reports whether a JID points to the Local ID (lid) server.
func isLIDJID(j types.JID) bool {
	return j.Server == "lid"
}

// normalizeFromField keeps "xxxxx@lid" as-is (for unresolved LID),
// otherwise strips "+" as we do for phone digits.
func normalizeFromField(s string) string {
	if strings.HasSuffix(strings.ToLower(strings.TrimSpace(s)), "@lid") {
		return s
	}
	return strings.ReplaceAll(s, "+", "")
}

// jidDisplayMaybeResolve returns a displayable identifier for the given JID:
// - for groups: the group id (jid.User)
// - for user AD JIDs: digits-only phone
// - for LID JIDs: best-effort resolve to AD via store/sync; if not resolved, returns "user@lid"
func jidDisplayMaybeResolve(sessionID string, cli *whatsmeow.Client, j types.JID) string {
	if isGroupJID(j) {
		return j.User
	}
	if !isLIDJID(j) {
		return digitsOnly(j.User)
	}

	// Try to resolve LID -> AD (s.whatsapp.net)
	if resolved, ok := resolveLIDToAD(sessionID, cli, j, defaultLIDResolveAttempts, defaultLIDResolveDelay); ok {
		return digitsOnly(resolved.User)
	}
	// Fallback: include @lid so downstream sees it's unresolved
	return fmt.Sprintf("%s@lid", j.User)
}

// resolveLIDToAD resolve JIDs @lid usando a API pública do whatsmeow (Store.LIDs).
// Tenta algumas vezes, aguardando pequeno intervalo entre tentativas.
func resolveLIDToAD(sessionID string, cli *whatsmeow.Client, lid types.JID, attempts int, delay time.Duration) (types.JID, bool) {
	entry := log.WithSession(sessionID)
	if cli == nil || cli.Store == nil || cli.Store.LIDs == nil {
		return types.JID{}, false
	}
	if v, ok := cacheGetLID(lid.String()); ok {
		return v, true
	}
	for i := 0; i < attempts; i++ {
		ad, err := cli.Store.LIDs.GetPNForLID(context.Background(), lid)
		if err == nil {
			if !isLIDJID(ad) && ad.User != "" && ad.Server != "" {
				entry.Info("lid_resolve success lid=%s -> %s", lid.String(), ad.String())
				cachePutLID(lid.String(), ad)
				return ad, true
			}
		}
		if i < attempts-1 {
			time.Sleep(delay)
		}
	}
	entry.Info("lid_resolve failed lid=%s using_fallback", lid.String())
	return types.JID{}, false
}

// resolveLIDToADReflect attempts to find an AD JID that corresponds to the given LID JID
// by peeking the client's store via reflection and optionally forcing a contacts fetch.
// This avoids compile-time coupling to whatsmeow's store internals while still providing
// best-effort behavior. Returns (zero, false) if not resolved.
// resolveLIDToADReflect was removed in favor of resolveLIDToAD using the public API.

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
			// Extract phone after the last ':' if present; some WA Business
			// vCards may have an empty value (e.g., "TEL;waid=123:") – in
			// that case, use the waid as the phone fallback.
			c := strings.LastIndex(s, ":")
			var phone string
			if c >= 0 && c+1 < len(s) {
				phone = strings.TrimSpace(s[c+1:])
			} else {
				phone = ""
			}
			// Extract waid parameter if present
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
			// Fallback: if phone is empty but waid is present, use waid
			if strings.TrimSpace(phone) == "" && strings.TrimSpace(waid) != "" {
				phone = strings.TrimSpace(waid)
			}
			// If both are empty, skip this TEL line
			if strings.TrimSpace(phone) == "" {
				continue
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

// safeGetGroupName chama Client.GetGroupInfo com reflexão, suportando (ctx,jid) e (jid).
func safeGetGroupName(cli *whatsmeow.Client, gid types.JID) (string, bool) {
	if cli == nil {
		return "", false
	}
	// Cache first
	if n, ok := cacheGetGroupName(gid.String()); ok {
		return n, true
	}
	// Try a couple of times with a small delay to allow metadata to load
	for i := 0; i < 2; i++ {
		info, err := cli.GetGroupInfo(gid)
		if err == nil && info != nil && strings.TrimSpace(info.Name) != "" {
			name := strings.TrimSpace(info.Name)
			cachePutGroupName(gid.String(), name)
			return name, true
		}
		if i == 0 {
			time.Sleep(200 * time.Millisecond)
		}
	}
	return "", false
}

// getAvatarURL tenta obter a URL da foto de perfil (contato ou grupo).
// Usa a API pública do whatsmeow; se não disponível/erro, retorna false.
func getAvatarURL(cli *whatsmeow.Client, jid types.JID) (string, bool) {
	if cli == nil || jid == (types.JID{}) {
		return "", false
	}
	// Cache first
	if u, ok := cacheGetAvatar(jid.String()); ok {
		return u, true
	}
	// Try a few combinations (prefer preview for speed)
	try := func(p *whatsmeow.GetProfilePictureParams) (string, bool) {
		if info, err := cli.GetProfilePictureInfo(jid, p); err == nil && info != nil && strings.TrimSpace(info.URL) != "" {
			url := strings.TrimSpace(info.URL)
			cachePutAvatar(jid.String(), url)
			return url, true
		}
		return "", false
	}
	// Try immediate attempts
	if url, ok := try(nil); ok {
		return url, true
	}
	if url, ok := try(&whatsmeow.GetProfilePictureParams{Preview: true}); ok {
		return url, true
	}
	if isGroupJID(jid) {
		if url, ok := try(&whatsmeow.GetProfilePictureParams{IsCommunity: true, Preview: true}); ok {
			return url, true
		}
		if url, ok := try(&whatsmeow.GetProfilePictureParams{IsCommunity: true}); ok {
			return url, true
		}
	}
	// One quick retry after a short delay in case metadata just loaded
	time.Sleep(150 * time.Millisecond)
	if url, ok := try(&whatsmeow.GetProfilePictureParams{Preview: true}); ok {
		return url, true
	}
	return "", false
}

// markReadEvent marks the given message as read when configured to do so.
func (m *ClientManager) markReadEvent(sessionID string, client *whatsmeow.Client, e *events.Message) {
	if client == nil || e == nil || e.Info.IsFromMe {
		return
	}
	chat := e.Info.Chat
	sender := e.Info.Sender
	if e.Info.SenderAlt != (types.JID{}) && !isLIDJID(e.Info.SenderAlt) && e.Info.SenderAlt.Server != "" {
		sender = e.Info.SenderAlt
	}
	// Best-effort resolve LID to AD
	if isLIDJID(chat) {
		if r, ok := resolveLIDToAD(sessionID, client, chat, 2, 200*time.Millisecond); ok {
			chat = r
		}
	}
	if isLIDJID(sender) {
		if r, ok := resolveLIDToAD(sessionID, client, sender, 2, 200*time.Millisecond); ok {
			sender = r
		}
	}
	// Send read receipt
	// Pequeno atraso para permitir que o consumidor processe o webhook primeiro
	time.Sleep(200 * time.Millisecond)
	ids := []types.MessageID{types.MessageID(e.Info.ID)}
	if err := client.MarkRead(ids, time.Now(), chat, sender); err != nil {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("mark_read error: %v", err)
	} else {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Info("mark_read success chat=%s sender=%s", chat.String(), sender.String())
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
		resp, err := webhookHTTPClient.Do(req.WithContext(context.Background()))
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
