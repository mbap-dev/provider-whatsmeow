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
	"sort"
	"strconv"
	"strings"
	"time"

	amqppub "your.org/provider-whatsmeow/internal/broker"
	"your.org/provider-whatsmeow/internal/log"
	"your.org/provider-whatsmeow/internal/status"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
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
		Proxy:                 http.ProxyFromEnvironment,
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

// handleCallRejection rejects an incoming call and optionally sends a reply message.
// It handles both LID and AD JIDs correctly to avoid bans.
func (m *ClientManager) handleCallRejection(sessionID string, client *whatsmeow.Client, callCreator types.JID, callID string) {
	if !m.getRejectCalls(sessionID) {
		return
	}

	// Reject the call
	if err := client.RejectCall(context.Background(), callCreator, callID); err != nil {
		log.WithSession(sessionID).Error("call reject error: %v", err)
	} else {
		log.WithSession(sessionID).Info("evt=call_reject from=%s call_id=%s", callCreator.String(), callID)
	}

	// Optionally send a message to the caller
	msgText := strings.TrimSpace(m.getRejectMsg(sessionID))
	if msgText == "" {
		return
	}

	// Support literal \n already handled in config; just send
	wire := &waE2E.Message{Conversation: proto.String(msgText)}

	// Resolve LID to AD if necessary to avoid bans
	to := callCreator
	shouldSend := true

	// Strip device suffix if present (e.g., "1234567890:10@s.whatsapp.net" -> "1234567890@s.whatsapp.net")
	// This is important for companion device calls - we need to create a new JID without the device part
	if strings.Contains(to.User, ":") {
		parts := strings.Split(to.User, ":")
		// Create a new JID with the user part without device suffix
		to = types.JID{
			User:       parts[0],
			RawAgent:   to.RawAgent,
			Device:     0, // No device suffix
			Integrator: to.Integrator,
			Server:     to.Server,
		}
		log.WithSession(sessionID).Info("call_reply_strip_device original=%s stripped=%s", callCreator.String(), to.String())
	}

	if isLIDJID(to) {
		if resolved, ok := resolveLIDToAD(sessionID, client, to, defaultLIDResolveAttempts, defaultLIDResolveDelay); ok {
			to = resolved
			log.WithSession(sessionID).Info("call_reply_lid_resolved from=%s to=%s", callCreator.String(), to.String())
		} else {
			log.WithSession(sessionID).Error("call_reply_lid_resolve_failed from=%s, skipping message send", callCreator.String())
			// Skip sending message if we can't resolve LID to avoid bans
			shouldSend = false
		}
	}

	if shouldSend {
		resp, err := client.SendMessage(context.Background(), to, wire)
		if err != nil {
			log.WithSession(sessionID).Error("call reply send error: %v", err)
		} else {
			// UnoAPI behavior: when reject-call auto-reply is sent, emit the
			// outgoing message payload (messages[]) in addition to statuses.
			if err := m.emitCloudOutgoingText(sessionID, client, to, resp.ID, msgText); err != nil {
				log.WithSession(sessionID).WithMessageID(string(resp.ID)).Error("call reply webhook message emit error: %v", err)
			}
			if err := m.emitCloudSent(sessionID, to, resp.ID); err != nil {
				log.WithSession(sessionID).WithMessageID(string(resp.ID)).Error("call reply webhook status emit error: %v", err)
			}
			log.WithSession(sessionID).WithMessageID(string(resp.ID)).Info("evt=call_reply_sent to=%s", to.String())
		}
	}
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
			log.WithSession(sessionID).Info("evt=call_offer from=%s call_id=%s", e.CallCreator.String(), e.CallID)
			m.handleCallRejection(sessionID, client, e.CallCreator, e.CallID)
		case *events.CallOfferNotice:
			// Group calls or calls from companion devices
			log.WithSession(sessionID).Info("evt=call_offer_notice from=%s call_id=%s media=%s type=%s", e.CallCreator.String(), e.CallID, e.Media, e.Type)
			m.handleCallRejection(sessionID, client, e.CallCreator, e.CallID)
		case *events.CallAccept:
			// Log when a call is accepted (might be from another device)
			log.WithSession(sessionID).Info("evt=call_accept from=%s call_id=%s", e.CallCreator.String(), e.CallID)
		case *events.CallTerminate:
			// Log when a call terminates
			log.WithSession(sessionID).Info("evt=call_terminate call_id=%s reason=%s", e.CallID, e.Reason)
		case *events.UnknownCallEvent:
			// Log unknown call events for debugging
			log.WithSession(sessionID).Info("evt=unknown_call_event")
		case *events.AppStateSyncComplete:
			if e.Name == "critical_block" || e.Name == "regular_high" {
				_ = client.SendPresence(context.Background(), types.PresenceAvailable)
			}
		case *events.PushNameSetting:
			_ = client.SendPresence(context.Background(), types.PresenceAvailable)
		case *events.Message:
			if err := m.emitCloudMessage(sessionID, client, e); err != nil {
				log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("webhook cloud message error: %v", err)
			}
			if m.getAutoMarkRead(sessionID) {
				go m.markReadEvent(sessionID, client, e)
			}
		case *events.GroupInfo:
			if err := m.emitCloudGroupInfo(sessionID, client, e); err != nil {
				log.WithSession(sessionID).Error("webhook cloud group info error: %v", err)
			}
		case *events.JoinedGroup:
			if err := m.emitCloudJoinedGroup(sessionID, client, e); err != nil {
				log.WithSession(sessionID).Error("webhook cloud joined group error: %v", err)
			}
		case *events.HistorySync:
			go m.emitCloudHistorySync(sessionID, client, e)
		case *events.Receipt:
			if err := m.emitCloudReceipt(sessionID, client, e); err != nil {
				log.WithSession(sessionID).Error("webhook cloud receipt error: %v", err)
			}
		default:
			// Silencia outros eventos no webhook (Uno espera apenas messages/statuses)
		}
	})
}

func (m *ClientManager) emitCloudHistorySync(sessionID string, client *whatsmeow.Client, e *events.HistorySync) {
	entry := log.WithSession(sessionID)
	if e == nil || e.Data == nil {
		return
	}
	if !m.getHistorySyncEnabled(sessionID) {
		return
	}
	if len(m.ListSessionWebhooks(sessionID)) == 0 {
		return
	}
	days := m.getHistorySyncDays(sessionID)
	var cutoff time.Time
	if days > 0 {
		cutoff = time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	}
	for _, conv := range e.Data.GetConversations() {
		if conv == nil {
			continue
		}
		chat, err := types.ParseJID(conv.GetID())
		if err != nil {
			entry.Error("evt=history_sync conversation parse failed id=%q err=%v", conv.GetID(), err)
			continue
		}
		for _, historyMsg := range conv.GetMessages() {
			if historyMsg == nil || historyMsg.GetMessage() == nil {
				continue
			}
			msgTS := time.Unix(int64(historyMsg.GetMessage().GetMessageTimestamp()), 0)
			if !cutoff.IsZero() && msgTS.Before(cutoff) {
				continue
			}
			evt, err := client.ParseWebMessage(chat, historyMsg.GetMessage())
			if err != nil {
				entry.Error("evt=history_sync message parse failed chat=%s err=%v", chat.String(), err)
				continue
			}
			if err := m.emitCloudMessage(sessionID, client, evt); err != nil {
				entry.WithMessageID(evt.Info.ID).Error("evt=history_sync webhook emit failed chat=%s err=%v", chat.String(), err)
				continue
			}
		}
	}
}

// =========== Emissão de mensagens no formato Cloud ===========

func (m *ClientManager) emitCloudMessage(sessionID string, client *whatsmeow.Client, e *events.Message) error {
	phone := normalizePhone(sessionID)
	msg := e.Message
	if msg == nil {
		return nil
	}
	if cacheConsumeSyntheticOutgoing(sessionID, e.Info.ID) {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Info("evt=message skip: synthetic_outgoing_already_emitted")
		return nil
	}

	fromMe := e.Info.IsFromMe
	chatJID := e.Info.Chat
	senderJID := e.Info.Sender
	ownMessage := fromMe || e.Info.DeviceSentMeta != nil
	if !ownMessage {
		ownMessage = isOwnJID(sessionID, client, senderJID) || isOwnJID(sessionID, client, e.Info.SenderAlt)
	}
	if ownMessage && e.Info.DeviceSentMeta != nil {
		if dst := strings.TrimSpace(e.Info.DeviceSentMeta.DestinationJID); dst != "" {
			if parsed, err := types.ParseJID(dst); err == nil && parsed.User != "" && parsed.Server != "" {
				chatJID = parsed
			}
		}
	}
	// When the message is ours, WhatsMeow may keep Chat as own JID and expose
	// the real recipient in RecipientAlt. Prefer RecipientAlt in this case.
	if ownMessage && e.Info.RecipientAlt != (types.JID{}) && e.Info.RecipientAlt.Server != "" {
		if chatJID.User == "" || isOwnJID(sessionID, client, chatJID) || isLIDJID(chatJID) {
			chatJID = e.Info.RecipientAlt
		}
	}
	// Preferir JID alternativo (AD) quando a lib expõe SenderAlt não‑LID
	if e.Info.SenderAlt != (types.JID{}) && !isLIDJID(e.Info.SenderAlt) && e.Info.SenderAlt.Server != "" {
		senderJID = e.Info.SenderAlt
	}
	// Ignore optional channels per configuration
	if m.getIgnoreStatusBroadcast(sessionID) && (isStatusBroadcastJID(chatJID) || isStatusBroadcastJID(senderJID)) {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Info("evt=message skip: status_broadcast")
		return nil
	}
	if m.getIgnoreNewsletters(sessionID) && (isNewsletterJID(chatJID) || isNewsletterJID(senderJID)) {
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
		// For group messages the payload keeps:
		// - contacts[0].wa_id as the participant phone/JID
		// - contacts[0].group_id as the group conversation id
		// Chatwoot uses that split to register the sender contact first and
		// then switch the conversation contact to the group.
		contactPhone = resolveDisplay(senderJID)
	}
	fromField := phone
	if !ownMessage {
		fromField = resolveDisplay(senderJID)
	}

	// Monta “message” no padrão Cloud
	wireMsg := map[string]any{
		"from":      normalizeFromField(fromField),
		"id":        e.Info.ID,
		"timestamp": strconv.FormatInt(e.Info.Timestamp.Unix(), 10),
	}
	if isGroupJID(chatJID) && senderJID.User != "" {
		wireMsg["id"] = withGroupParticipantMessageID(senderJID.String(), e.Info.ID)
	}

	// Context (reply/quote) se houver
	if ctx := messageContextInfo(sessionID, client, msg, isGroupJID(chatJID)); ctx != nil {
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
			body = decorateMentionText(sessionID, client, body, edited)
			body = body + "\n`Mensagem Editada`"
			wireMsg["type"] = "text"
			wireMsg["text"] = map[string]any{"body": body}
			if pm.GetKey() != nil {
				rid := strings.TrimSpace(pm.GetKey().GetID())
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
				rid = strings.TrimSpace(pm.GetKey().GetID())
			}
			if rid == "" {
				// Without the original id, do nothing
				return nil
			}
			revokedID := deletedMessagePublicID(sessionID, client, rid, e.Info.Chat, e.Info.Sender, pm.GetKey())
			if err := m.emitCloudDeleted(sessionID, client, e, revokedID); err != nil {
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
		if ctx := reactionContextInfo(sessionID, client, r, isGroupJID(chatJID), senderJID); ctx != nil {
			wireMsg["context"] = ctx
		}
	case msg.GetButtonsResponseMessage() != nil:
		if interactive := buildInboundButtonsResponsePayload(msg.GetButtonsResponseMessage()); interactive != nil {
			wireMsg["type"] = "interactive"
			wireMsg["interactive"] = interactive
		} else {
			return nil
		}
	case msg.GetTemplateButtonReplyMessage() != nil:
		if interactive := buildInboundTemplateButtonReplyPayload(msg.GetTemplateButtonReplyMessage()); interactive != nil {
			wireMsg["type"] = "interactive"
			wireMsg["interactive"] = interactive
		} else {
			return nil
		}
	case msg.GetListResponseMessage() != nil:
		if interactive := buildInboundListResponsePayload(msg.GetListResponseMessage()); interactive != nil {
			wireMsg["type"] = "interactive"
			wireMsg["interactive"] = interactive
		} else {
			return nil
		}
	case msg.GetInteractiveResponseMessage() != nil:
		if interactive := buildInboundInteractiveResponsePayload(msg.GetInteractiveResponseMessage()); interactive != nil {
			wireMsg["type"] = "interactive"
			wireMsg["interactive"] = interactive
		} else {
			return nil
		}
	case msg.GetConversation() != "":
		body := decorateMentionText(sessionID, client, msg.GetConversation(), msg)
		wireMsg["type"] = "text"
		wireMsg["text"] = map[string]any{"body": body}

	case msg.GetExtendedTextMessage() != nil:
		t := strings.TrimSpace(msg.GetExtendedTextMessage().GetText())
		t = decorateMentionText(sessionID, client, t, msg)
		wireMsg["type"] = "text"
		wireMsg["text"] = map[string]any{"body": t}

	case msg.GetImageMessage() != nil:
		im := msg.GetImageMessage()
		caption := decorateMentionText(sessionID, client, im.GetCaption(), msg)
		wireMsg["type"] = "image"
		mimeType := im.GetMimetype()
		fileName := filenameByMime(e.Info.ID, mimeType)
		image := map[string]any{
			"caption":     caption,
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(im.GetFileSHA256()),
			"media_key":   b64(im.GetMediaKey()),
			"direct_path": im.GetDirectPath(),
		}
		if fileName != "" {
			image["filename"] = fileName
		}
		fileSize := int64(0)
		if im.FileLength != nil {
			image["file_length"] = im.GetFileLength()
			fileSize = int64(im.GetFileLength())
		}
		if media, err := m.persistInboundMediaFromWA(sessionID, client, e.Info.ID, im, MediaSaveInput{
			MessagingProduct: "whatsapp",
			MIMEType:         splitMime(mimeType),
			SHA256:           b64(im.GetFileSHA256()),
			FileSize:         fileSize,
			Filename:         fileName,
		}); err == nil {
			image["id"] = CompatibleMediaID(media)
		} else {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Debug("media_store image skipped: %v", err)
		}
		wireMsg["image"] = image

	case msg.GetDocumentMessage() != nil:
		d := msg.GetDocumentMessage()
		caption := decorateMentionText(sessionID, client, d.GetCaption(), msg)
		wireMsg["type"] = "document"
		mimeType := d.GetMimetype()
		fileName := strings.TrimSpace(firstNonEmpty(d.GetFileName(), d.GetTitle()))
		if fileName == "" {
			fileName = filenameByMime(e.Info.ID, mimeType)
		}
		document := map[string]any{
			"caption":     caption,
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(d.GetFileSHA256()),
			"media_key":   b64(d.GetMediaKey()),
			"direct_path": d.GetDirectPath(),
		}
		if fileName != "" {
			document["filename"] = fileName
		}
		fileSize := int64(0)
		if d.FileLength != nil {
			document["file_length"] = d.GetFileLength()
			fileSize = int64(d.GetFileLength())
		}
		if media, err := m.persistInboundMediaFromWA(sessionID, client, e.Info.ID, d, MediaSaveInput{
			MessagingProduct: "whatsapp",
			MIMEType:         splitMime(mimeType),
			SHA256:           b64(d.GetFileSHA256()),
			FileSize:         fileSize,
			Filename:         fileName,
		}); err == nil {
			document["id"] = CompatibleMediaID(media)
		} else {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Debug("media_store document skipped: %v", err)
		}
		wireMsg["document"] = document

	case msg.GetVideoMessage() != nil:
		v := msg.GetVideoMessage()
		caption := decorateMentionText(sessionID, client, v.GetCaption(), msg)
		wireMsg["type"] = "video"
		mimeType := v.GetMimetype()
		videoFileName := filenameByMime(e.Info.ID, mimeType)
		video := map[string]any{
			"caption":     caption,
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(v.GetFileSHA256()),
			"media_key":   b64(v.GetMediaKey()),
			"direct_path": v.GetDirectPath(),
		}
		if videoFileName != "" {
			video["filename"] = videoFileName
		}
		fileSize := int64(0)
		if v.FileLength != nil {
			video["file_length"] = v.GetFileLength()
			fileSize = int64(v.GetFileLength())
		}
		if media, err := m.persistInboundMediaFromWA(sessionID, client, e.Info.ID, v, MediaSaveInput{
			MessagingProduct: "whatsapp",
			MIMEType:         splitMime(mimeType),
			SHA256:           b64(v.GetFileSHA256()),
			FileSize:         fileSize,
			Filename:         videoFileName,
		}); err == nil {
			video["id"] = CompatibleMediaID(media)
		} else {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Debug("media_store video skipped: %v", err)
		}
		wireMsg["video"] = video

	case msg.GetAudioMessage() != nil:
		a := msg.GetAudioMessage()
		wireMsg["type"] = "audio"
		mimeType := a.GetMimetype()
		audioFileName := filenameByMime(e.Info.ID, mimeType)
		audio := map[string]any{
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(a.GetFileSHA256()),
			"media_key":   b64(a.GetMediaKey()),
			"direct_path": a.GetDirectPath(),
		}
		if audioFileName != "" {
			audio["filename"] = audioFileName
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
		fileSize := int64(0)
		if a.FileLength != nil {
			fileSize = int64(a.GetFileLength())
		}
		if media, err := m.persistInboundMediaFromWA(sessionID, client, e.Info.ID, a, MediaSaveInput{
			MessagingProduct: "whatsapp",
			MIMEType:         splitMime(mimeType),
			SHA256:           b64(a.GetFileSHA256()),
			FileSize:         fileSize,
			Filename:         audioFileName,
		}); err == nil {
			audio["id"] = CompatibleMediaID(media)
		} else {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Debug("media_store audio skipped: %v", err)
		}
		wireMsg["audio"] = audio

	case msg.GetStickerMessage() != nil:
		s := msg.GetStickerMessage()
		wireMsg["type"] = "sticker"
		mimeType := s.GetMimetype()
		filename := filenameByMime(e.Info.ID, mimeType)
		sticker := map[string]any{
			"mime_type":   splitMime(mimeType),
			"sha256":      b64(s.GetFileSHA256()),
			"media_key":   b64(s.GetMediaKey()),
			"direct_path": s.GetDirectPath(),
		}
		if filename != "" {
			sticker["filename"] = filename
		}
		if s.FileLength != nil {
			sticker["file_length"] = s.GetFileLength()
		}
		fileSize := int64(0)
		if s.FileLength != nil {
			fileSize = int64(s.GetFileLength())
		}
		if media, err := m.persistInboundMediaFromWA(sessionID, client, e.Info.ID, s, MediaSaveInput{
			MessagingProduct: "whatsapp",
			MIMEType:         splitMime(mimeType),
			SHA256:           b64(s.GetFileSHA256()),
			FileSize:         fileSize,
			Filename:         filename,
		}); err == nil {
			sticker["id"] = CompatibleMediaID(media)
		} else {
			log.WithSession(sessionID).WithMessageID(e.Info.ID).Debug("media_store sticker skipped: %v", err)
		}
		wireMsg["sticker"] = sticker

	case msg.GetLocationMessage() != nil:
		wireMsg["type"] = "location"
		wireMsg["location"] = buildInboundLocationPayload(msg.GetLocationMessage())

	case msg.GetLiveLocationMessage() != nil:
		wireMsg["type"] = "location"
		wireMsg["location"] = buildInboundLiveLocationPayload(msg.GetLiveLocationMessage())

	default:
		// fallback: se nada identificado mas tem algo, tenta extrair caption/text
		if cap := extractAnyCaption(msg); cap != "" {
			cap = decorateMentionText(sessionID, client, cap, msg)
			wireMsg["type"] = "text"
			wireMsg["text"] = map[string]any{"body": cap}
		} else {
			// ignora tipos não suportados
			return nil
		}
	}
	// Mantém o shape do webhook compatível com Chatwoot:
	// - profile.* representa o contato/remetente
	// - group_* representa o grupo quando a mensagem vem de um grupo
	var (
		dispName      string
		groupSubject  string
		groupPicture  string
		profileAvatar = chatJID
	)
	if isGroupJID(chatJID) {
		if gname, ok := safeGetGroupName(client, chatJID); ok && strings.TrimSpace(gname) != "" {
			groupSubject = gname
		} else {
			groupSubject = chatJID.User
		}
		dispName = strings.TrimSpace(e.Info.PushName)
		if dispName == "" {
			if cachedName, ok := safeGetContactName(sessionID, client, senderJID, e.Info.Sender, e.Info.SenderAlt); ok {
				dispName = cachedName
			}
		}
		if dispName == "" {
			dispName = contactPhone
		}
		if senderJID.User != "" && senderJID.Server != "" {
			profileAvatar = senderJID
		}
		if pic, ok := getAvatarURL(sessionID, client, chatJID); ok && strings.TrimSpace(pic) != "" {
			groupPicture = pic
		}
	} else {
		if ownMessage {
			dispName = contactPhone
		} else {
			dispName = strings.TrimSpace(e.Info.PushName)
			if dispName == "" {
				if cachedName, ok := safeGetContactName(sessionID, client, senderJID, e.Info.Sender, e.Info.SenderAlt); ok {
					dispName = cachedName
				}
			}
			if dispName == "" {
				dispName = contactPhone
			}
		}
	}
	profile := buildCloudProfile(sessionID, client, profileAvatar, dispName, e.Info.ID)
	contactObj := buildCloudContactObject(profile, contactPhone, chatJID, groupSubject, groupPicture)
	if isGroupJID(chatJID) {
		if participants := m.groupParticipantsPayload(sessionID, client, chatJID); len(participants) > 0 {
			contactObj["group_participants"] = participants
		}
	}
	if isGroupJID(chatJID) {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Debug(
			"contact_group built group_id=%s subject=%q wa_id=%s sender_name=%q", chatJID.String(), groupSubject, contactPhone, dispName,
		)
	} else {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Debug(
			"contact_user built wa_id=%s name=%q", contactPhone, dispName,
		)
	}

	payload := cloudEnvelope(phone)
	change := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)
	val := change["value"].(map[string]any)
	val["contacts"] = []any{contactObj}
	val["messages"] = []any{wireMsg}
	if ownMessage && m.getOutgoingMessagesCoex(sessionID) {
		convertCloudPayloadToOutgoingEchoes(payload)
	}

	log.WithSession(sessionID).WithMessageID(e.Info.ID).
		Info("evt=message cloud payload ready field=messages type=%s", wireMsg["type"])

	return m.deliverWebhook(sessionID, payload)
}

func (m *ClientManager) emitCloudJoinedGroup(sessionID string, client *whatsmeow.Client, e *events.JoinedGroup) error {
	if e == nil {
		return nil
	}
	evt := &events.GroupInfo{
		JID:       e.JID,
		Notify:    e.Notify,
		Sender:    e.Sender,
		SenderPN:  e.SenderPN,
		Timestamp: time.Now(),
		Join:      []types.JID{},
	}
	if e.JID.User == "" || e.JID.Server == "" {
		return nil
	}
	return m.emitCloudGroupInfo(sessionID, client, evt)
}

func (m *ClientManager) emitCloudGroupInfo(sessionID string, client *whatsmeow.Client, e *events.GroupInfo) error {
	if e == nil || !isGroupJID(e.JID) {
		return nil
	}

	action, affected := groupParticipantChange(e)
	if action == "" && len(affected) == 0 {
		return nil
	}

	phone := normalizePhone(sessionID)
	groupSubject := strings.TrimSpace(e.JID.User)
	if groupName, ok := safeGetGroupName(client, e.JID); ok && strings.TrimSpace(groupName) != "" {
		groupSubject = groupName
	}
	groupPicture := ""
	if pic, ok := getAvatarURL(sessionID, client, e.JID); ok && strings.TrimSpace(pic) != "" {
		groupPicture = pic
	}

	actorJID := groupEventActorJID(e)
	actorID := jidDisplayMaybeResolve(sessionID, client, actorJID)
	if actorID == "" {
		actorID = phone
	}
	actorName := actorID
	if name, ok := safeGetContactName(sessionID, client, actorJID); ok && strings.TrimSpace(name) != "" {
		actorName = name
	}
	profile := buildCloudProfile(sessionID, client, actorJID, actorName, "")
	contactObj := buildCloudContactObject(profile, actorID, e.JID, groupSubject, groupPicture)

	participants := m.groupParticipantsPayload(sessionID, client, e.JID)
	if len(participants) > 0 {
		contactObj["group_participants"] = participants
	}
	contactObj["group_event"] = map[string]any{
		"action":       action,
		"actor":        groupEventParticipantPayload(sessionID, client, actorJID),
		"participants": groupEventParticipantsPayload(sessionID, client, affected),
		"timestamp":    groupEventTimestamp(e).Unix(),
	}

	msgID := fmt.Sprintf("group-event-%s-%s-%d", e.JID.User, action, groupEventTimestamp(e).UnixNano())
	wireMsg := map[string]any{
		"from":      normalizeFromField(actorID),
		"id":        msgID,
		"timestamp": strconv.FormatInt(groupEventTimestamp(e).Unix(), 10),
		"type":      "group_event",
		"text":      map[string]any{"body": action},
	}

	payload := cloudEnvelope(phone)
	change := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)
	val := change["value"].(map[string]any)
	val["contacts"] = []any{contactObj}
	val["messages"] = []any{wireMsg}

	log.WithSession(sessionID).Info(
		"evt=group_info cloud payload ready group=%s action=%s affected=%d participants=%d",
		e.JID.String(),
		action,
		len(affected),
		len(participants),
	)
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
	recipient = strings.ReplaceAll(recipient, "+", "")
	conversationID := ""
	if e.Chat != (types.JID{}) {
		conversationID = e.Chat.String()
	}

	states := make([]any, 0, len(e.MessageIDs))
	for _, id := range e.MessageIDs {
		msgID := string(id)
		receiptKey := strings.Join([]string{
			sessionID,
			status,
			msgID,
			recipient,
			conversationID,
		}, "|")
		if !cacheTryPutReceipt(receiptKey) {
			log.WithSession(sessionID).WithMessageID(msgID).
				Debug("evt=receipt duplicate skipped type=%s recipient=%s conversation=%s", e.Type, recipient, conversationID)
			continue
		}

		statusID := msgID
		statusID = withGroupParticipantMessageID(phone, msgID)
		st := map[string]any{
			"id":           statusID,
			"recipient_id": recipient,
			"status":       status,
		}
		if !e.Timestamp.IsZero() {
			st["timestamp"] = strconv.FormatInt(e.Timestamp.Unix(), 10)
		}
		// conversation.id ajuda o Uno a agrupar
		if e.Chat != (types.JID{}) {
			st["conversation"] = map[string]any{"id": conversationID}
		}
		states = append(states, st)
	}
	if len(states) == 0 {
		return nil
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["statuses"] = states

	log.WithSession(sessionID).Info("evt=receipt cloud payload ready type=%s ids=%v", e.Type, e.MessageIDs)

	return m.deliverWebhook(sessionID, payload)
}

func (m *ClientManager) emitCloudSent(sessionID string, to types.JID, id types.MessageID) error {
	phone := normalizePhone(sessionID)
	statusID := string(id)
	statusID = withGroupParticipantMessageID(phone, string(id))
	recipient := jidToPhoneNumberIfUser(to)
	st := map[string]any{
		"id":           statusID,
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

// emitCloudFailed publishes a status event for a failed outgoing message,
// following the same shape used by the Baileys provider in unoapi-cloud.
func (m *ClientManager) emitCloudFailed(sessionID string, to string, messageID string, code int, title string) error {
	phone := normalizePhone(sessionID)
	recipient := digitsOnly(strings.ReplaceAll(strings.TrimSpace(to), "+", ""))
	if recipient == "" {
		recipient = strings.TrimSpace(to)
	}

	msgID := strings.TrimSpace(messageID)
	if msgID == "" {
		msgID = fmt.Sprintf("WARN-%d", time.Now().UnixNano())
	}

	status := map[string]any{
		"id":           msgID,
		"recipient_id": recipient,
		"status":       "failed",
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
	}
	if code != 0 || strings.TrimSpace(title) != "" {
		status["errors"] = []any{map[string]any{
			"code":  code,
			"title": title,
		}}
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["statuses"] = []any{status}

	log.WithSession(sessionID).WithMessageID(msgID).Info("evt=failed cloud payload ready to=%s code=%d", to, code)
	return m.deliverWebhook(sessionID, payload)
}

// emitCloudDeleted publishes a status event for a message deletion (revoke)
// in a Cloud-like webhook format.
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

func deletedMessagePublicID(sessionID string, cli *whatsmeow.Client, messageID string, chat, sender types.JID, key *waCommon.MessageKey) string {
	revokedID := strings.TrimSpace(messageID)
	if revokedID == "" {
		return ""
	}
	if !isGroupJID(chat) {
		return revokedID
	}

	participant := strings.TrimSpace(key.GetParticipant())
	if participant == "" && key.GetFromMe() {
		participant = sender.String()
	}
	if participant == "" {
		participant = sender.String()
	}

	participantID := contextParticipantMessageIDUser(sessionID, cli, participant)
	if participantID == "" {
		participantID = participant
	}

	return withGroupParticipantMessageID(participantID, revokedID)
}

// emitCloudOutgoingText publishes an outgoing text message event in the same
// Cloud-like format used for incoming messages, so downstreams can see the
// content that was sent automatically (e.g., call rejection reply).
func (m *ClientManager) emitCloudOutgoingText(sessionID string, client *whatsmeow.Client, to types.JID, id types.MessageID, body string) error {
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
		"profile": buildCloudProfile(sessionID, client, to, contactPhone, string(id)),
		"wa_id":   contactPhone,
	}
	if isGroupJID(to) {
		contactObj["group_id"] = to.String()
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["contacts"] = []any{contactObj}
	val["messages"] = []any{wireMsg}
	if m.getOutgoingMessagesCoex(sessionID) {
		convertCloudPayloadToOutgoingEchoes(payload)
	}

	log.WithSession(sessionID).WithMessageID(string(id)).
		Info("evt=message_out cloud payload ready field=messages type=text to=%s", to.String())

	cacheMarkSyntheticOutgoing(sessionID, string(id))
	return m.deliverWebhook(sessionID, payload)
}

func (m *ClientManager) emitCloudOutgoingTemplateMessage(sessionID string, client *whatsmeow.Client, to types.JID, id types.MessageID, rendered *renderedTemplate) error {
	if rendered == nil {
		return nil
	}
	phone := normalizePhone(sessionID)
	recipient := jidToPhoneNumberIfUser(to)

	mediaMeta, hasMediaMeta := m.resolveTemplateWebhookMedia(context.Background(), sessionID, string(id), rendered)
	wireMsg := buildTemplateWebhookWireMessage(phone, string(id), rendered, mediaMeta, hasMediaMeta)

	contactObj := map[string]any{
		"profile": buildCloudProfile(sessionID, client, to, recipient, string(id)),
		"wa_id":   recipient,
	}
	if isGroupJID(to) {
		contactObj["group_id"] = to.String()
	}

	payload := cloudEnvelope(phone)
	val := payload["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
	val["contacts"] = []any{contactObj}
	val["messages"] = []any{wireMsg}
	if m.getOutgoingMessagesCoex(sessionID) {
		convertCloudPayloadToOutgoingEchoes(payload)
	}

	log.WithSession(sessionID).WithMessageID(string(id)).
		Info("evt=message_out cloud payload ready field=messages type=%s source=template", wireMsg["type"])

	cacheMarkSyntheticOutgoing(sessionID, string(id))
	return m.deliverWebhook(sessionID, payload)
}

func (m *ClientManager) resolveTemplateWebhookMedia(ctx context.Context, sessionID, messageID string, rendered *renderedTemplate) (MediaMetadata, bool) {
	if m == nil || rendered == nil {
		return MediaMetadata{}, false
	}
	switch {
	case rendered.HeaderDocument != nil:
		doc := rendered.HeaderDocument
		if mediaID := strings.TrimSpace(doc.ID); mediaID != "" {
			if meta, ok := m.GetMediaMetadata(sessionID, mediaID); ok {
				return meta, true
			}
		}
		link := strings.TrimSpace(doc.Link)
		if link == "" {
			return MediaMetadata{}, false
		}
		data, mimeType, header, err := downloadBytesWithHeader(ctx, link)
		if err != nil {
			log.WithSession(sessionID).WithMessageID(messageID).Error("template webhook document download failed: %v", err)
			return MediaMetadata{}, false
		}
		filename := strings.TrimSpace(doc.FileName)
		if filename == "" {
			filename = filenameFromHTTPHeader(header)
		}
		if filename == "" {
			filename = documentFilenameFromLink(link, messageID)
		}
		meta, err := m.StoreUploadedMedia(ctx, sessionID, filename, mimeType, data)
		if err != nil {
			log.WithSession(sessionID).WithMessageID(messageID).Error("template webhook document persist failed: %v", err)
			return MediaMetadata{}, false
		}
		doc.ID = CompatibleMediaID(meta)
		if strings.TrimSpace(doc.FileName) == "" {
			doc.FileName = meta.Filename
		}
		return meta, true
	case rendered.HeaderImage != nil:
		img := rendered.HeaderImage
		if mediaID := strings.TrimSpace(img.ID); mediaID != "" {
			if meta, ok := m.GetMediaMetadata(sessionID, mediaID); ok {
				return meta, true
			}
		}
		link := strings.TrimSpace(img.Link)
		if link == "" {
			return MediaMetadata{}, false
		}
		data, mimeType, header, err := downloadBytesWithHeader(ctx, link)
		if err != nil {
			log.WithSession(sessionID).WithMessageID(messageID).Error("template webhook image download failed: %v", err)
			return MediaMetadata{}, false
		}
		filename := filenameFromHTTPHeader(header)
		if filename == "" {
			filename = documentFilenameFromLink(link, messageID)
		}
		if filename == "" {
			filename = filenameByMime(messageID, mimeType)
		}
		meta, err := m.StoreUploadedMedia(ctx, sessionID, filename, mimeType, data)
		if err != nil {
			log.WithSession(sessionID).WithMessageID(messageID).Error("template webhook image persist failed: %v", err)
			return MediaMetadata{}, false
		}
		img.ID = CompatibleMediaID(meta)
		return meta, true
	default:
		return MediaMetadata{}, false
	}
}

func buildTemplateWebhookWireMessage(phone, messageID string, rendered *renderedTemplate, mediaMeta MediaMetadata, hasMediaMeta bool) map[string]any {
	wireMsg := map[string]any{
		"from":      strings.ReplaceAll(normalizePhone(phone), "+", ""),
		"id":        messageID,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}
	body := ""
	if rendered != nil {
		body = rendered.text()
	}
	switch {
	case rendered != nil && rendered.HeaderDocument != nil:
		doc := map[string]any{"caption": body}
		if filename := strings.TrimSpace(rendered.HeaderDocument.FileName); filename != "" {
			doc["filename"] = filename
		}
		if hasMediaMeta {
			applyMediaMetadataToWebhookPart(doc, mediaMeta)
		} else if mediaID := strings.TrimSpace(rendered.HeaderDocument.ID); mediaID != "" {
			doc["id"] = mediaID
		}
		wireMsg["type"] = "document"
		wireMsg["document"] = doc
	case rendered != nil && rendered.HeaderImage != nil:
		image := map[string]any{"caption": body}
		if hasMediaMeta {
			applyMediaMetadataToWebhookPart(image, mediaMeta)
		} else if mediaID := strings.TrimSpace(rendered.HeaderImage.ID); mediaID != "" {
			image["id"] = mediaID
		}
		wireMsg["type"] = "image"
		wireMsg["image"] = image
	default:
		wireMsg["type"] = "text"
		wireMsg["text"] = map[string]any{"body": body}
	}
	return wireMsg
}

func applyMediaMetadataToWebhookPart(part map[string]any, meta MediaMetadata) {
	if id := CompatibleMediaID(meta); id != "" {
		part["id"] = id
	}
	if mimeType := strings.TrimSpace(meta.MIMEType); mimeType != "" {
		part["mime_type"] = splitMime(mimeType)
	}
	if sha := strings.TrimSpace(meta.SHA256); sha != "" {
		part["sha256"] = sha
	}
	if meta.FileSize > 0 {
		part["file_length"] = meta.FileSize
	}
	if filename := strings.TrimSpace(meta.Filename); filename != "" {
		part["filename"] = filename
	}
}

func filenameFromHTTPHeader(header http.Header) string {
	if header == nil {
		return ""
	}
	cd := strings.TrimSpace(header.Get("Content-Disposition"))
	if cd == "" {
		return ""
	}
	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(params["filename"])
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

func convertCloudPayloadToOutgoingEchoes(payload map[string]any) {
	entries, _ := payload["entry"].([]any)
	for _, entryAny := range entries {
		entry, _ := entryAny.(map[string]any)
		changes, _ := entry["changes"].([]any)
		for _, changeAny := range changes {
			change, _ := changeAny.(map[string]any)
			value, _ := change["value"].(map[string]any)
			messages, ok := value["messages"].([]any)
			if !ok || len(messages) == 0 {
				continue
			}
			change["field"] = "smb_message_echoes"
			value["message_echoes"] = messages
			delete(value, "messages")
		}
	}
}

func normalizePhone(s string) string {
	if looksLikePhoneIdentifier(s) {
		if d := digitsOnly(s); d != "" {
			return d
		}
	}
	return strings.TrimSpace(s)
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

// isCommunityJID reports whether the JID belongs to a Community entity.
// Communities are not the same as regular groups (g.us) and require
// different query parameters for profile picture lookups in whatsmeow.
// Keep this conservative to avoid breaking regular group lookups.
func isCommunityJID(j types.JID) bool {
	// Known server for communities in WhatsApp.
	return j.Server == "communities"
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

// isOwnJID reports whether j represents the current session account.
// We compare against:
// - normalized session id digits
// - connected device PN user
// - connected device LID user
func isOwnJID(sessionID string, cli *whatsmeow.Client, j types.JID) bool {
	if j.User == "" {
		return false
	}

	if sid := digitsOnly(normalizePhone(sessionID)); sid != "" && digitsOnly(j.User) == sid {
		return true
	}

	if cli != nil && cli.Store != nil && cli.Store.ID != nil {
		own := cli.Store.ID.ToNonAD()
		if own.User != "" && (j.User == own.User || digitsOnly(j.User) == digitsOnly(own.User)) {
			return true
		}
	}

	if cli != nil {
		if ownLID := cli.DangerousInternals().GetOwnLID(); ownLID.User != "" && j.User == ownLID.User {
			return true
		}
	}

	return false
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
		// Strip device suffix (":<device>") from AD JIDs before deriving phone.
		user := j.User
		if i := strings.Index(user, ":"); i >= 0 {
			user = user[:i]
		}
		return digitsOnly(user)
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
	if v, ok := cacheGetLID(lid.String()); ok {
		return v, true
	}
	if cli == nil || cli.Store == nil || cli.Store.LIDs == nil {
		return types.JID{}, false
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

func (m *ClientManager) persistInboundMediaFromWA(sessionID string, client *whatsmeow.Client, mediaID string, msg whatsmeow.DownloadableMessage, input MediaSaveInput) (MediaMetadata, error) {
	if m == nil || client == nil || msg == nil {
		return MediaMetadata{}, fmt.Errorf("media persistence unavailable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	log.WithSession(sessionID).WithMessageID(mediaID).Debug(
		"evt=media_inbound download_start mime=%s filename=%q declared_size=%d",
		strings.TrimSpace(input.MIMEType),
		strings.TrimSpace(input.Filename),
		input.FileSize,
	)
	data, err := client.Download(ctx, msg)
	if err != nil {
		log.WithSession(sessionID).WithMessageID(mediaID).Error("evt=media_inbound download_failed err=%v", err)
		return MediaMetadata{}, err
	}
	log.WithSession(sessionID).WithMessageID(mediaID).Info(
		"evt=media_inbound download_ok bytes=%d mime=%s filename=%q",
		len(data),
		strings.TrimSpace(input.MIMEType),
		strings.TrimSpace(input.Filename),
	)
	meta, err := m.StoreInboundMedia(ctx, sessionID, mediaID, data, input)
	if err != nil {
		log.WithSession(sessionID).WithMessageID(mediaID).Error("evt=media_inbound persist_failed err=%v", err)
		return MediaMetadata{}, err
	}
	log.WithSession(sessionID).WithMessageID(mediaID).Info(
		"evt=media_inbound persist_ok public_id=%s relative_path=%s mode=%s",
		meta.ID,
		meta.RelativePath,
		m.MediaMode(),
	)
	return meta, nil
}

func extensionByMime(m string) string {
	exts, _ := mime.ExtensionsByType(splitMime(m))
	if len(exts) > 0 {
		return exts[0]
	}
	return ""
}

func filenameByMime(messageID, mimeType string) string {
	if strings.TrimSpace(messageID) == "" {
		return ""
	}
	switch splitMime(strings.ToLower(mimeType)) {
	case "audio/ogg":
		return messageID + ".ogg"
	case "video/mp4":
		return messageID + ".mp4"
	case "application/vnd.ms-excel":
		return messageID + ".xls"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return messageID + ".xlsx"
	case "application/msword":
		return messageID + ".doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return messageID + ".docx"
	case "application/vnd.ms-powerpoint":
		return messageID + ".ppt"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return messageID + ".pptx"
	}
	ext := extensionByMime(mimeType)
	if ext == "" {
		return ""
	}
	return messageID + ext
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
		// Handle TEL lines with optional itemN. prefix (e.g., "item1.TEL;...")
		// Detect property name before ':' and look for "TEL" within it.
		upper := strings.ToUpper(s)
		prop := upper
		if k := strings.Index(upper, ":"); k >= 0 {
			prop = upper[:k]
		}
		if strings.Contains(prop, "TEL") {
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
			// Normalize digits and apply BR mobile 9th-digit heuristic when necessary.
			dPhone := digitsOnly(phone)
			dWAID := digitsOnly(waid)
			if len(dPhone) >= 2 && strings.HasPrefix(dPhone, "55") {
				dPhone = jidToPhoneNumberDigits(dPhone, true)
			}
			if len(dWAID) >= 2 && strings.HasPrefix(dWAID, "55") {
				dWAID = jidToPhoneNumberDigits(dWAID, true)
			}
			// Fallback: if phone is empty but waid is present, use waid
			if strings.TrimSpace(dPhone) == "" && strings.TrimSpace(dWAID) != "" {
				dPhone = dWAID
			}
			// If both are empty, skip this TEL line
			if strings.TrimSpace(dPhone) == "" {
				continue
			}
			entry := map[string]any{"phone": dPhone}
			if dWAID != "" {
				entry["wa_id"] = dWAID
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

func extractMessageContextInfo(m *waE2E.Message) *waE2E.ContextInfo {
	if m == nil {
		return nil
	}
	switch {
	case m.GetButtonsResponseMessage() != nil:
		return m.GetButtonsResponseMessage().GetContextInfo()
	case m.GetTemplateButtonReplyMessage() != nil:
		return m.GetTemplateButtonReplyMessage().GetContextInfo()
	case m.GetListResponseMessage() != nil:
		return m.GetListResponseMessage().GetContextInfo()
	case m.GetInteractiveResponseMessage() != nil:
		return m.GetInteractiveResponseMessage().GetContextInfo()
	case m.GetExtendedTextMessage() != nil:
		return m.GetExtendedTextMessage().GetContextInfo()
	case m.GetImageMessage() != nil:
		return m.GetImageMessage().GetContextInfo()
	case m.GetDocumentMessage() != nil:
		return m.GetDocumentMessage().GetContextInfo()
	case m.GetVideoMessage() != nil:
		return m.GetVideoMessage().GetContextInfo()
	case m.GetAudioMessage() != nil:
		return m.GetAudioMessage().GetContextInfo()
	case m.GetStickerMessage() != nil:
		return m.GetStickerMessage().GetContextInfo()
	case m.GetLocationMessage() != nil:
		return m.GetLocationMessage().GetContextInfo()
	case m.GetLiveLocationMessage() != nil:
		return m.GetLiveLocationMessage().GetContextInfo()
	}
	return nil
}

func buildInboundLocationPayload(in *waE2E.LocationMessage) map[string]any {
	if in == nil {
		return nil
	}
	return map[string]any{
		"latitude":  in.GetDegreesLatitude(),
		"longitude": in.GetDegreesLongitude(),
	}
}

func buildInboundLiveLocationPayload(in *waE2E.LiveLocationMessage) map[string]any {
	if in == nil {
		return nil
	}
	return map[string]any{
		"latitude":  in.GetDegreesLatitude(),
		"longitude": in.GetDegreesLongitude(),
	}
}

func buildInboundButtonsResponsePayload(in *waE2E.ButtonsResponseMessage) map[string]any {
	if in == nil {
		return nil
	}
	return buildInboundButtonReplyPayload(in.GetSelectedButtonID(), in.GetSelectedDisplayText(), nil)
}

func buildInboundTemplateButtonReplyPayload(in *waE2E.TemplateButtonReplyMessage) map[string]any {
	if in == nil {
		return nil
	}
	extra := map[string]any{}
	if idx := in.GetSelectedIndex(); idx > 0 {
		extra["index"] = idx
	}
	if idx := in.GetSelectedCarouselCardIndex(); idx > 0 {
		extra["carousel_card_index"] = idx
	}
	return buildInboundButtonReplyPayload(in.GetSelectedID(), in.GetSelectedDisplayText(), extra)
}

func buildInboundListResponsePayload(in *waE2E.ListResponseMessage) map[string]any {
	if in == nil {
		return nil
	}
	selected := in.GetSingleSelectReply()
	if selected == nil {
		return nil
	}
	listReply := map[string]any{}
	if id := strings.TrimSpace(selected.GetSelectedRowID()); id != "" {
		listReply["id"] = id
	}
	if title := strings.TrimSpace(in.GetTitle()); title != "" {
		listReply["title"] = title
	}
	if description := strings.TrimSpace(in.GetDescription()); description != "" {
		listReply["description"] = description
	}
	if len(listReply) == 0 {
		return nil
	}
	return map[string]any{
		"type":       "list_reply",
		"list_reply": listReply,
	}
}

func buildInboundInteractiveResponsePayload(in *waE2E.InteractiveResponseMessage) map[string]any {
	if in == nil {
		return nil
	}

	bodyText := strings.TrimSpace(in.GetBody().GetText())
	native := in.GetNativeFlowResponseMessage()
	if native == nil {
		return nil
	}

	name := strings.TrimSpace(native.GetName())
	paramsJSON := strings.TrimSpace(native.GetParamsJSON())
	paramsMap := parseJSONObject(paramsJSON)

	switch name {
	case "quick_reply":
		title := firstNonEmpty(
			jsonObjectString(paramsMap, "title"),
			jsonObjectString(paramsMap, "display_text"),
			bodyText,
		)
		return buildInboundButtonReplyPayload(jsonObjectString(paramsMap, "id"), title, nil)
	case "single_select":
		return buildInboundListReplyPayload(
			firstNonEmpty(
				jsonObjectString(paramsMap, "id"),
				jsonObjectString(paramsMap, "selected_row_id"),
				jsonObjectString(paramsMap, "row_id"),
			),
			firstNonEmpty(
				jsonObjectString(paramsMap, "title"),
				jsonObjectString(paramsMap, "display_text"),
				bodyText,
			),
			jsonObjectString(paramsMap, "description"),
		)
	}

	nfmReply := map[string]any{}
	if name != "" {
		nfmReply["name"] = name
	}
	if bodyText != "" {
		nfmReply["body"] = bodyText
	}
	if paramsMap != nil {
		nfmReply["response_json"] = paramsMap
	} else if paramsJSON != "" {
		nfmReply["response_json"] = paramsJSON
	}
	if version := native.GetVersion(); version > 0 {
		nfmReply["version"] = version
	}
	if len(nfmReply) == 0 {
		return nil
	}
	return map[string]any{
		"type":      "nfm_reply",
		"nfm_reply": nfmReply,
	}
}

func buildInboundButtonReplyPayload(id, title string, extra map[string]any) map[string]any {
	buttonReply := map[string]any{}
	if id = strings.TrimSpace(id); id != "" {
		buttonReply["id"] = id
	}
	if title = strings.TrimSpace(title); title != "" {
		buttonReply["title"] = title
	}
	for key, value := range extra {
		if value != nil {
			buttonReply[key] = value
		}
	}
	if len(buttonReply) == 0 {
		return nil
	}
	return map[string]any{
		"type":         "button_reply",
		"button_reply": buttonReply,
	}
}

func buildInboundListReplyPayload(id, title, description string) map[string]any {
	listReply := map[string]any{}
	if id = strings.TrimSpace(id); id != "" {
		listReply["id"] = id
	}
	if title = strings.TrimSpace(title); title != "" {
		listReply["title"] = title
	}
	if description = strings.TrimSpace(description); description != "" {
		listReply["description"] = description
	}
	if len(listReply) == 0 {
		return nil
	}
	return map[string]any{
		"type":       "list_reply",
		"list_reply": listReply,
	}
}

func parseJSONObject(raw string) map[string]any {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func jsonObjectString(in map[string]any, key string) string {
	if in == nil {
		return ""
	}
	value, ok := in[key]
	if !ok {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case fmt.Stringer:
		return strings.TrimSpace(typed.String())
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func messageContextInfo(sessionID string, cli *whatsmeow.Client, m *waE2E.Message, groupChat bool) map[string]any {
	// tenta extrair stanzaId e participant de qualquer ContextInfo disponível
	ci := extractMessageContextInfo(m)
	if ci == nil {
		return nil
	}

	ctx := make(map[string]any)
	if stanzaID := strings.TrimSpace(ci.GetStanzaID()); stanzaID != "" {
		contextID := stanzaID
		if participant := strings.TrimSpace(ci.GetParticipant()); shouldKeepContextParticipant(sessionID, groupChat, participant) {
			contextID = withGroupParticipantMessageID(contextParticipantMessageIDUser(sessionID, cli, participant), stanzaID)
		}
		ctx["message_id"] = contextID
		ctx["id"] = contextID
	}

	if ci.GetIsForwarded() {
		ctx["forwarded"] = true
	}

	if len(ctx) == 0 {
		if ci.GetQuotedMessage() != nil {
			// se veio quotedMessage mas sem stanzaId, ainda assim envia context pra compat
			return map[string]any{
				"message_id": "",
				"id":         "",
			}
		}
		return nil
	}

	return ctx
}

func reactionContextInfo(sessionID string, cli *whatsmeow.Client, r *waE2E.ReactionMessage, groupChat bool, sender types.JID) map[string]any {
	if r == nil || r.GetKey() == nil {
		return nil
	}
	key := r.GetKey()
	messageID := strings.TrimSpace(key.GetID())
	if messageID == "" {
		return nil
	}
	contextID := messageID
	participant := strings.TrimSpace(key.GetParticipant())
	if participant == "" && groupChat && key.GetFromMe() {
		participant = sender.String()
	}
	if participant == "" && !groupChat && key.GetFromMe() {
		participant = normalizePhone(sessionID)
	}
	if participant == "" && !groupChat && looksLikeOwnWhatsAppMessageID(messageID) {
		participant = normalizePhone(sessionID)
	}
	if shouldKeepContextParticipant(sessionID, groupChat, participant) {
		contextID = withGroupParticipantMessageID(contextParticipantMessageIDUser(sessionID, cli, participant), messageID)
	}
	return map[string]any{
		"message_id": contextID,
		"id":         contextID,
	}
}

func looksLikeOwnWhatsAppMessageID(messageID string) bool {
	messageID = strings.ToUpper(strings.TrimSpace(messageID))
	return strings.HasPrefix(messageID, "3EB")
}

func shouldKeepContextParticipant(sessionID string, groupChat bool, participant string) bool {
	if strings.TrimSpace(participant) == "" {
		return false
	}
	if groupChat {
		return true
	}
	return participantMessageIDUser(participant) == normalizePhone(sessionID)
}

func contextParticipantMessageIDUser(sessionID string, cli *whatsmeow.Client, participant string) string {
	raw := strings.TrimSpace(participant)
	if raw == "" {
		return ""
	}
	jid, err := types.ParseJID(raw)
	if err != nil {
		return participantMessageIDUser(raw)
	}
	if isLIDJID(jid) {
		if resolved, ok := resolveLIDToAD(sessionID, cli, jid, defaultLIDResolveAttempts, defaultLIDResolveDelay); ok {
			return participantMessageIDUser(resolved.ToNonAD().String())
		}
	}
	return participantMessageIDUser(jid.ToNonAD().String())
}

func safeGetContactInfo(cli *whatsmeow.Client, candidates ...types.JID) (types.ContactInfo, bool) {
	if cli == nil || cli.Store == nil || cli.Store.Contacts == nil {
		return types.ContactInfo{}, false
	}
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		lookup := normalizedContactLookupJID(candidate)
		if lookup == (types.JID{}) {
			continue
		}
		key := lookup.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		info, err := cli.Store.Contacts.GetContact(ctx, lookup)
		cancel()
		if err == nil && info.Found {
			return info, true
		}
	}
	return types.ContactInfo{}, false
}

func resolvedMentionJID(sessionID string, cli *whatsmeow.Client, jid types.JID) types.JID {
	if !isLIDJID(jid) {
		return types.JID{}
	}
	if resolved, ok := resolveLIDToAD(sessionID, cli, jid, defaultLIDResolveAttempts, defaultLIDResolveDelay); ok {
		return resolved
	}
	return types.JID{}
}

func mentionReplacementCandidates(jid, resolved types.JID, info types.ContactInfo) []string {
	seen := map[string]struct{}{}
	var out []string
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		if i := strings.Index(value, ":"); i >= 0 {
			value = value[:i]
		}
		if value == "" {
			return
		}
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			out = append(out, value)
		}
		if digits := digitsOnly(value); digits != "" {
			if _, ok := seen[digits]; !ok {
				seen[digits] = struct{}{}
				out = append(out, digits)
			}
		}
	}

	add(jid.User)
	add(resolved.User)
	if strings.TrimSpace(info.RedactedPhone) != "" {
		add(info.RedactedPhone)
	}
	return out
}

func mentionDisplayNumber(jid, resolved types.JID, info types.ContactInfo) string {
	if resolved.User != "" && resolved.Server != "" {
		if digits := digitsOnly(resolved.User); digits != "" {
			return digits
		}
	}
	if !isLIDJID(jid) {
		if digits := digitsOnly(jid.User); digits != "" {
			return digits
		}
	}
	number := strings.TrimSpace(firstNonEmpty(
		digitsOnly(info.RedactedPhone),
		strings.TrimSpace(info.RedactedPhone),
	))
	if number != "" && !strings.Contains(number, "@") {
		return number
	}
	return firstNonEmpty(
		digitsOnly(jid.User),
		strings.TrimSpace(jid.User),
	)
}

func formatMentionDisplay(name, number string) string {
	name = strings.TrimSpace(name)
	number = strings.TrimSpace(number)
	switch {
	case name != "" && number != "" && name != number:
		return fmt.Sprintf("@%s (%s)", name, number)
	case name != "":
		return "@" + name
	case number != "":
		return "@" + number
	default:
		return ""
	}
}

func applyMentionReplacements(text string, replacements map[string]string) string {
	if text == "" || len(replacements) == 0 {
		return text
	}
	keys := make([]string, 0, len(replacements))
	for key := range replacements {
		if strings.Contains(text, key) {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	for _, key := range keys {
		text = strings.ReplaceAll(text, key, replacements[key])
	}
	return text
}

func decorateMentionText(sessionID string, cli *whatsmeow.Client, text string, m *waE2E.Message) string {
	if strings.TrimSpace(text) == "" || m == nil {
		return text
	}
	ci := extractMessageContextInfo(m)
	if ci == nil || len(ci.GetMentionedJID()) == 0 {
		return text
	}

	replacements := make(map[string]string)
	for _, raw := range ci.GetMentionedJID() {
		jid, err := types.ParseJID(strings.TrimSpace(raw))
		if err != nil {
			continue
		}
		resolved := resolvedMentionJID(sessionID, cli, jid)
		info, _ := safeGetContactInfo(cli, jid, resolved)
		label := formatMentionDisplay(
			contactInfoDisplayName(info),
			mentionDisplayNumber(jid, resolved, info),
		)
		if label == "" {
			continue
		}
		for _, candidate := range mentionReplacementCandidates(jid, resolved, info) {
			replacements["@"+candidate] = label
		}
	}
	return applyMentionReplacements(text, replacements)
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

// safeGetGroupName obtém o nome do grupo com um pequeno cache.
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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		info, err := cli.GetGroupInfo(ctx, gid)
		cancel()
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

func contactInfoDisplayName(info types.ContactInfo) string {
	return strings.TrimSpace(firstNonEmpty(
		info.PushName,
		info.BusinessName,
		info.FullName,
		info.FirstName,
	))
}

func normalizedContactLookupJID(j types.JID) types.JID {
	if j.User == "" || j.Server == "" {
		return types.JID{}
	}
	if isGroupJID(j) || isLIDJID(j) {
		return j
	}
	if i := strings.Index(j.User, ":"); i >= 0 {
		j.User = j.User[:i]
	}
	// User profile/contact lookups expect the base user JID, not a device-scoped
	// variant like 5511999999999:3@s.whatsapp.net.
	j.Device = 0
	return j
}

func normalizedAvatarLookupJID(sessionID string, cli *whatsmeow.Client, j types.JID) types.JID {
	if j.User == "" || j.Server == "" {
		return types.JID{}
	}

	if isLIDJID(j) {
		if resolved, ok := resolveLIDToAD(sessionID, cli, j, defaultLIDResolveAttempts, defaultLIDResolveDelay); ok {
			return normalizedContactLookupJID(resolved)
		}
	}

	return normalizedContactLookupJID(j)
}

func safeGetContactName(sessionID string, cli *whatsmeow.Client, candidates ...types.JID) (string, bool) {
	if cli == nil || cli.Store == nil || cli.Store.Contacts == nil {
		return "", false
	}
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		lookup := normalizedContactLookupJID(candidate)
		if lookup == (types.JID{}) {
			continue
		}
		key := lookup.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		if name, ok := cacheGetContactName(key); ok {
			return name, true
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		info, err := cli.Store.Contacts.GetContact(ctx, lookup)
		cancel()
		if err != nil || !info.Found {
			continue
		}

		name := contactInfoDisplayName(info)
		if name == "" {
			continue
		}
		cachePutContactName(key, name)
		return name, true
	}
	return "", false
}

// getAvatarURL tenta obter a URL da foto de perfil (contato ou grupo).
// Usa a API pública do whatsmeow; se não disponível/erro, retorna false.
func getAvatarURL(sessionID string, cli *whatsmeow.Client, jid types.JID) (string, bool) {
	jid = normalizedAvatarLookupJID(sessionID, cli, jid)
	if cli == nil || jid == (types.JID{}) {
		return "", false
	}
	// Cache first
	if u, ok := cacheGetAvatar(jid.String()); ok {
		if u == avatarMissingSentinel {
			log.WithSession(sessionID).Debug("avatar_cache_miss_hit jid=%s", jid.String())
			return "", false
		}
		log.WithSession(sessionID).Debug("avatar_cache_hit jid=%s url=%s", jid.String(), u)
		return u, true
	}
	// Try a few combinations (prefer preview for speed)
	try := func(label string, p *whatsmeow.GetProfilePictureParams) (string, bool) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		info, err := cli.GetProfilePictureInfo(ctx, jid, p)
		cancel()
		if err != nil {
			log.WithSession(sessionID).Debug("avatar_try label=%s jid=%s err=%v", label, jid.String(), err)
			return "", false
		}
		if info != nil && strings.TrimSpace(info.URL) != "" {
			url := strings.TrimSpace(info.URL)
			cachePutAvatar(jid.String(), url)
			log.WithSession(sessionID).Debug("avatar_success label=%s jid=%s url=%s", label, jid.String(), url)
			return url, true
		}
		log.WithSession(sessionID).Debug("avatar_try_empty label=%s jid=%s", label, jid.String())
		return "", false
	}
	// Try immediate attempts
	if url, ok := try("default", nil); ok {
		return url, true
	}
	if url, ok := try("preview", &whatsmeow.GetProfilePictureParams{Preview: true}); ok {
		return url, true
	}
	// Communities (not regular groups) require a different parameter
	if isCommunityJID(jid) {
		if url, ok := try("community_preview", &whatsmeow.GetProfilePictureParams{IsCommunity: true, Preview: true}); ok {
			return url, true
		}
		if url, ok := try("community", &whatsmeow.GetProfilePictureParams{IsCommunity: true}); ok {
			return url, true
		}
	}
	// One quick retry after a short delay in case metadata just loaded
	time.Sleep(150 * time.Millisecond)
	if url, ok := try("preview_retry", &whatsmeow.GetProfilePictureParams{Preview: true}); ok {
		return url, true
	}
	cachePutAvatarMissing(jid.String())
	log.WithSession(sessionID).Debug("avatar_not_available jid=%s", jid.String())
	return "", false
}

func buildCloudProfile(sessionID string, cli *whatsmeow.Client, avatarJID types.JID, displayName string, messageID string) map[string]any {
	profile := map[string]any{
		"name":    displayName,
		"picture": "",
	}
	avatarLookupJID := normalizedAvatarLookupJID(sessionID, cli, avatarJID)
	if avatarLookupJID == (types.JID{}) {
		return profile
	}
	if av, ok := getAvatarURL(sessionID, cli, avatarLookupJID); ok && strings.TrimSpace(av) != "" {
		profile["picture"] = av
		log.WithSession(sessionID).WithMessageID(messageID).Debug(
			"avatar_set url=%s chat=%s lookup=%s",
			av,
			avatarJID.String(),
			avatarLookupJID.String(),
		)
	} else {
		log.WithSession(sessionID).WithMessageID(messageID).Debug(
			"avatar_not_found chat=%s lookup=%s is_group=%t",
			avatarJID.String(),
			avatarLookupJID.String(),
			isGroupJID(avatarLookupJID),
		)
	}
	return profile
}

func buildCloudContactObject(profile map[string]any, waID string, chatJID types.JID, groupSubject, groupPicture string) map[string]any {
	contactObj := map[string]any{
		"profile": profile,
		"wa_id":   waID,
	}
	if isGroupJID(chatJID) {
		contactObj["group_id"] = chatJID.String()
		if strings.TrimSpace(groupSubject) != "" {
			contactObj["group_subject"] = groupSubject
		}
		if strings.TrimSpace(groupPicture) != "" {
			contactObj["group_picture"] = groupPicture
		}
	}
	return contactObj
}

func (m *ClientManager) groupParticipantsPayload(sessionID string, cli *whatsmeow.Client, groupJID types.JID) []any {
	if cli == nil || !isGroupJID(groupJID) {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	info, err := cli.GetGroupInfo(ctx, groupJID)
	cancel()
	if err != nil || info == nil {
		log.WithSession(sessionID).Debug("group_participants fetch failed group=%s err=%v", groupJID.String(), err)
		return nil
	}

	participants := make([]any, 0, len(info.Participants))
	for _, participant := range info.Participants {
		cacheGroupParticipantLID(participant.PhoneNumber, participant.LID)
		payload := groupParticipantPayload(sessionID, cli, participant.JID, participant.PhoneNumber, participant.LID)
		if payload == nil {
			continue
		}
		payload["is_admin"] = participant.IsAdmin
		payload["is_super_admin"] = participant.IsSuperAdmin
		if strings.TrimSpace(participant.DisplayName) != "" {
			payload["name"] = strings.TrimSpace(participant.DisplayName)
		}
		participants = append(participants, payload)
	}
	return participants
}

func cacheGroupParticipantLID(phoneNumber, lid types.JID) {
	if phoneNumber.User == "" || lid.User == "" || !isLIDJID(lid) {
		return
	}
	cachePutLID(lid.String(), normalizedContactLookupJID(phoneNumber))
}

func groupParticipantPayload(sessionID string, cli *whatsmeow.Client, candidates ...types.JID) map[string]any {
	jid := firstUsableJID(candidates...)
	if jid == (types.JID{}) {
		return nil
	}
	displayID := jidDisplayMaybeResolve(sessionID, cli, jid)
	if displayID == "" {
		displayID = jid.String()
	}
	name := displayID
	if contactName, ok := safeGetContactName(sessionID, cli, candidates...); ok && strings.TrimSpace(contactName) != "" {
		name = contactName
	}
	payload := map[string]any{
		"wa_id": displayID,
		"name":  name,
	}
	if phoneNumber := normalizeCloudContactPhoneNumber(displayID); phoneNumber != "" {
		payload["phone_number"] = phoneNumber
	}
	if pic, ok := getAvatarURL(sessionID, cli, jid); ok && strings.TrimSpace(pic) != "" {
		payload["picture"] = pic
	}
	return payload
}

func normalizeCloudContactPhoneNumber(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasSuffix(strings.ToLower(value), "@lid") {
		return ""
	}
	digits := digitsOnly(value)
	if digits == "" {
		return ""
	}

	e164 := normalizeE164Local("+" + digits)
	if e164 == "" {
		e164 = normalizeE164Local(digits)
	}
	if e164 == "" {
		return ""
	}
	return preferredCloudContactE164(e164)
}

func preferredCloudContactE164(e164 string) string {
	candidates := candidatesBR(e164)
	if strings.HasPrefix(e164, "+55") {
		best := ""
		for _, candidate := range candidates {
			if len(digitsOnly(candidate)) > len(digitsOnly(best)) {
				best = candidate
			}
		}
		if best != "" {
			return best
		}
	}
	for _, candidate := range candidates {
		if validE164PhoneNumber(candidate) {
			return candidate
		}
	}
	return e164
}

func firstUsableJID(candidates ...types.JID) types.JID {
	for _, candidate := range candidates {
		if candidate.User != "" && candidate.Server != "" && !isLIDJID(candidate) {
			return candidate
		}
	}
	for _, candidate := range candidates {
		if candidate.User != "" && candidate.Server != "" {
			return candidate
		}
	}
	return types.JID{}
}

func groupEventActorJID(e *events.GroupInfo) types.JID {
	if e == nil {
		return types.JID{}
	}
	if e.SenderPN != nil && e.SenderPN.User != "" {
		return *e.SenderPN
	}
	if e.Sender != nil {
		return *e.Sender
	}
	return types.JID{}
}

func groupEventTimestamp(e *events.GroupInfo) time.Time {
	if e == nil || e.Timestamp.IsZero() {
		return time.Now()
	}
	return e.Timestamp
}

func groupParticipantChange(e *events.GroupInfo) (string, []types.JID) {
	if e == nil {
		return "", nil
	}
	switch {
	case len(e.Join) > 0:
		return "add", e.Join
	case len(e.Leave) > 0:
		return "remove", e.Leave
	case len(e.Promote) > 0:
		return "promote", e.Promote
	case len(e.Demote) > 0:
		return "demote", e.Demote
	default:
		return "", nil
	}
}

func groupEventParticipantPayload(sessionID string, cli *whatsmeow.Client, jid types.JID) map[string]any {
	payload := groupParticipantPayload(sessionID, cli, jid)
	if payload == nil {
		return map[string]any{}
	}
	return payload
}

func groupEventParticipantsPayload(sessionID string, cli *whatsmeow.Client, participants []types.JID) []any {
	out := make([]any, 0, len(participants))
	for _, participant := range participants {
		payload := groupParticipantPayload(sessionID, cli, participant)
		if payload == nil {
			continue
		}
		out = append(out, payload)
	}
	return out
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
	if err := client.MarkRead(context.Background(), ids, time.Now(), chat, sender, types.ReceiptTypeRead); err != nil {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Error("mark_read error: %v", err)
	} else {
		log.WithSession(sessionID).WithMessageID(e.Info.ID).Info("mark_read success chat=%s sender=%s", chat.String(), sender.String())
	}
}

// =========== Entrega HTTP (mantida) ===========

func (m *ClientManager) deliverWebhook(sessionID string, payload any) error {
	entry := log.WithSession(sessionID)
	rk := normalizePhone(sessionID)
	if data, ok := payload.(map[string]any); ok {
		applyInstanceMessageIDs(rk, data)
	}
	amqpErr := error(nil)
	// Prefer AMQP publication; fallback to HTTP if AMQP is unavailable
	if err := amqppub.PublishWebhook(rk, payload); err == nil {
		entry.Info("webhook amqp publish ok rk=%s", rk)
	} else {
		amqpErr = err
		entry.Error("webhook amqp publish error: %v", err)
	}

	// Registered per-session webhooks (UnoAPI-like) have priority for HTTP delivery.
	if m.dispatchRegisteredWebhooks(sessionID, payload) {
		return nil
	}

	// If AMQP worked and there are no registered webhooks we are done.
	if amqpErr == nil {
		return nil
	}

	// HTTP fallback when AMQP not available
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
