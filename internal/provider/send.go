package provider

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	"your.org/provider-whatsmeow/internal/log"

	"go.mau.fi/whatsmeow"
	waBinary "go.mau.fi/whatsmeow/binary"
	goE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/hajimehoshi/go-mp3"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"layeh.com/gopus"
)

// (no external resolver)

// TextContent representa {"text":{"body":"..."}}
type TextContent struct {
	Body string `json:"body"`
}

type ImageContent struct {
	Link    string `json:"link"`
	ID      string `json:"id,omitempty"`
	Caption string `json:"caption,omitempty"`
}

type DocumentContent struct {
	Link     string `json:"link"`
	ID       string `json:"id,omitempty"`
	Caption  string `json:"caption,omitempty"`
	FileName string `json:"filename,omitempty"`
}

type VideoContent struct {
	Link    string `json:"link"`
	ID      string `json:"id,omitempty"`
	Caption string `json:"caption,omitempty"`
}

type AudioContent struct {
	Link    string  `json:"link"`
	ID      string  `json:"id,omitempty"`
	Caption string  `json:"caption,omitempty"`
	PTT     *bool   `json:"ptt,omitempty"`
	Seconds *uint32 `json:"seconds,omitempty"`
}

type StickerContent struct {
	Link string `json:"link"`
	ID   string `json:"id,omitempty"`
}

type InteractiveText struct {
	Text string `json:"text"`
}

type InteractiveReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type InteractiveButton struct {
	Type  string           `json:"type"`
	Reply InteractiveReply `json:"reply"`
}

type InteractiveRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type InteractiveSection struct {
	Title string           `json:"title"`
	Rows  []InteractiveRow `json:"rows"`
}

type InteractiveAction struct {
	Button   string               `json:"button,omitempty"`
	Buttons  []InteractiveButton  `json:"buttons,omitempty"`
	Sections []InteractiveSection `json:"sections,omitempty"`
}

type InteractiveContent struct {
	Type   string             `json:"type"`
	Header *InteractiveText   `json:"header,omitempty"`
	Body   *InteractiveText   `json:"body,omitempty"`
	Footer *InteractiveText   `json:"footer,omitempty"`
	Action *InteractiveAction `json:"action,omitempty"`
}

type PaymentAmount struct {
	Value  int64 `json:"value"`
	Offset int64 `json:"offset"`
}

type PaymentOrderItem struct {
	RetailerID string        `json:"retailer_id,omitempty"`
	Name       string        `json:"name"`
	Amount     PaymentAmount `json:"amount"`
	Quantity   int64         `json:"quantity"`
}

type PaymentOrderTax struct {
	Value       int64  `json:"value"`
	Offset      int64  `json:"offset"`
	Description string `json:"description,omitempty"`
}

type PaymentOrder struct {
	Status   string             `json:"status,omitempty"`
	Items    []PaymentOrderItem `json:"items,omitempty"`
	Subtotal PaymentAmount      `json:"subtotal"`
	Tax      *PaymentOrderTax   `json:"tax,omitempty"`
	Discount *PaymentOrderTax   `json:"discount,omitempty"`
	Shipping *PaymentOrderTax   `json:"shipping,omitempty"`
}

type PIXStaticCodeInfo struct {
	MerchantName string `json:"merchant_name"`
	Key          string `json:"key"`
	KeyType      string `json:"key_type"`
}

type BoletoInfo struct {
	DigitableLine string `json:"digitable_line"`
}

type PaymentContent struct {
	ReferenceID        string             `json:"reference_id,omitempty"`
	Currency           string             `json:"currency,omitempty"`
	Type               string             `json:"type,omitempty"`
	PaymentType        string             `json:"payment_type,omitempty"`
	Body               string             `json:"body,omitempty"`
	HeaderDocument     *DocumentContent   `json:"header_document,omitempty"`
	TotalAmount        PaymentAmount      `json:"total_amount"`
	Order              *PaymentOrder      `json:"order,omitempty"`
	ItemName           string             `json:"item_name,omitempty"`
	Quantity           int64              `json:"quantity,omitempty"`
	PIXStaticCode      *PIXStaticCodeInfo `json:"pix_static_code,omitempty"`
	Boleto             *BoletoInfo        `json:"boleto,omitempty"`
	SharePaymentStatus bool               `json:"share_payment_status,omitempty"`
	Referral           string             `json:"referral,omitempty"`
}

// Contacts payload (Cloud-compatible)
type ContactName struct {
	FormattedName string `json:"formatted_name"`
}

type ContactPhone struct {
	Phone string `json:"phone"`
	WAID  string `json:"wa_id,omitempty"`
}

type ContactEntry struct {
	Name   ContactName    `json:"name"`
	Phones []ContactPhone `json:"phones"`
}

type MessageContext struct {
	ID           string   `json:"id,omitempty"`
	MessageID    string   `json:"message_id,omitempty"`
	StanzaId     string   `json:"stanzaId,omitempty"`
	Stanza_id    string   `json:"stanza_id,omitempty"`
	Participant  string   `json:"participant,omitempty"`
	QuotedText   string   `json:"quoted_text,omitempty"`
	MentionedJid []string `json:"mentioned_jid,omitempty"`
	Mentions     []string `json:"mentions,omitempty"`
}

type OutgoingMessage struct {
	MessagingProduct string              `json:"messaging_product,omitempty"`
	RecipientType    string              `json:"recipient_type,omitempty"`
	Type             string              `json:"type"`
	To               string              `json:"to"`
	Body             string              `json:"body,omitempty"`
	Text             *TextContent        `json:"text,omitempty"`
	Image            *ImageContent       `json:"image,omitempty"`
	Document         *DocumentContent    `json:"document,omitempty"`
	Video            *VideoContent       `json:"video,omitempty"`
	Audio            *AudioContent       `json:"audio,omitempty"`
	Sticker          *StickerContent     `json:"sticker,omitempty"`
	Interactive      *InteractiveContent `json:"interactive,omitempty"`
	Payment          *PaymentContent     `json:"payment,omitempty"`
	OrderDetails     *PaymentContent     `json:"order_details,omitempty"`
	Contacts         []ContactEntry      `json:"contacts,omitempty"`
	MediaURL         string              `json:"media_url,omitempty"`
	Caption          string              `json:"caption,omitempty"`
	SessionID        string              `json:"session_id"`
	MessageID        string              `json:"message_id,omitempty"`
	Context          *MessageContext     `json:"context,omitempty"`
	Mentions         []string            `json:"mentions,omitempty"`
	MentionAll       bool                `json:"mention_all,omitempty"`
	MentionAllCompat bool                `json:"mentionAll,omitempty"`
}

// marshalForLog pretty-prints a proto message for logging, truncating if too long.
func marshalForLog(m proto.Message) string {
	if m == nil {
		return "<nil>"
	}
	b, err := (protojson.MarshalOptions{EmitUnpopulated: true}).Marshal(m)
	if err != nil {
		return fmt.Sprintf("<protojson error: %v>", err)
	}
	s := string(b)
	const max = 2000
	if len(s) > max {
		return s[:max] + "...(truncated)"
	}
	return s
}

func sendWithExtra(ctx context.Context, cli *whatsmeow.Client, jid types.JID, msg *goE2E.Message, id string, extra *whatsmeow.SendRequestExtra) (whatsmeow.SendResponse, error) {
	if extra == nil {
		return sendWithID(ctx, cli, jid, msg, id)
	}
	req := *extra
	if strings.TrimSpace(id) != "" {
		req.ID = types.MessageID(id)
	}
	return cli.SendMessage(ctx, jid, msg, req)
}

func sendWithID(ctx context.Context, cli *whatsmeow.Client, jid types.JID, msg *goE2E.Message, id string) (whatsmeow.SendResponse, error) {
	if strings.TrimSpace(id) != "" {
		return cli.SendMessage(ctx, jid, msg, whatsmeow.SendRequestExtra{ID: types.MessageID(id)})
	}
	return cli.SendMessage(ctx, jid, msg)
}

func (m OutgoingMessage) paymentContent() *PaymentContent {
	if m.Payment != nil {
		return m.Payment
	}
	return m.OrderDetails
}

// Send envia a mensagem; sessão vem do contexto (routing key) ou do payload.
func (m *ClientManager) Send(ctx context.Context, msg OutgoingMessage) error {
	sessionID := strings.TrimSpace(msg.SessionID)
	if v := ctx.Value(CtxKeyPhoneNumberID); v != nil {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			sessionID = strings.TrimSpace(s)
		}
	}
	if sessionID == "" {
		return fmt.Errorf("missing session_id (not in context nor payload)")
	}
	if !m.AllowsSession(sessionID) {
		return fmt.Errorf("session %s is not assigned to worker %s", sessionID, m.WorkerID())
	}
	normalizeOutgoingMessageIDs(sessionID, &msg)

	cli, ok := m.GetClient(sessionID)
	if !ok || cli == nil {
		return fmt.Errorf("session %s not connected", sessionID)
	}
	entry := log.WithSession(sessionID).WithMessageID(msg.MessageID)

	// Log input summary early for diagnostics
	entry.Debug("send_input type=%q to=%q media_url=%q has_text=%v has_image=%v has_doc=%v has_audio=%v",
		strings.TrimSpace(msg.Type), strings.TrimSpace(msg.To), strings.TrimSpace(msg.MediaURL),
		msg.Text != nil && strings.TrimSpace(msg.Text.Body) != "",
		msg.Image != nil && strings.TrimSpace(msg.Image.Link) != "",
		msg.Document != nil && strings.TrimSpace(msg.Document.Link) != "",
		msg.Audio != nil && strings.TrimSpace(msg.Audio.Link) != "",
	)

	// Normalize destination: optionally verify user existence via server, then
	// resolve to a canonical JID. If the existence check fails due to network
	// issues we fall back to the previous behaviour and still attempt delivery.
	rawTo := strings.TrimSpace(msg.To)
	if rawTo == "" {
		entry.Error(
			"send_input_missing_recipient type=%q media_url=%q has_text=%v has_image=%v has_doc=%v has_audio=%v",
			strings.TrimSpace(msg.Type),
			strings.TrimSpace(msg.MediaURL),
			msg.Text != nil && strings.TrimSpace(msg.Text.Body) != "",
			msg.Image != nil && strings.TrimSpace(msg.Image.Link) != "",
			msg.Document != nil && strings.TrimSpace(msg.Document.Link) != "",
			msg.Audio != nil && strings.TrimSpace(msg.Audio.Link) != "",
		)
		return fmt.Errorf("missing recipient: use to/recipient_id/source_id")
	}
	var jid types.JID
	var err error

	// Decide whether to enforce WhatsApp existence check for this destination.
	shouldCheck := m.getCheckUserExists(sessionID)
	if strings.Contains(rawTo, "@g.us") || strings.Contains(rawTo, "@broadcast") || strings.Contains(rawTo, "@lid") {
		// Groups, broadcast lists and LIDs are not checked with IsOnWhatsApp
		shouldCheck = false
	}

	if shouldCheck {
		rjid, found, checkErr := m.checkUserExistsOnWhatsApp(ctx, cli, rawTo)
		if checkErr != nil {
			entry.Info("whatsapp_check_failed to=%q err=%v; falling back to legacy normalisation", rawTo, checkErr)
		} else if !found {
			entry.Info("whatsapp_check_not_found to=%q", rawTo)
			failTitle := fmt.Sprintf("The phone number %s does not have WhatsApp account!", rawTo)
			if err := m.emitCloudFailed(sessionID, rawTo, msg.MessageID, 2, failTitle); err != nil {
				entry.Error("emit_failed_status_error to=%q err=%v", rawTo, err)
			}
			return fmt.Errorf("number %s is not registered on WhatsApp", rawTo)
		} else {
			jid = rjid
		}
	}

	// If we still don't have a JID (check disabled or failed), use the existing
	// resolution logic: server‑side PN resolver when possible, with local
	// normalisation as fallback.
	if jid == (types.JID{}) {
		if rjid, ok := m.resolvePNtoJID(ctx, cli, rawTo); ok {
			jid = rjid
		} else {
			// Fallback: local normalization -> PN JID
			digits := phoneNumberToJIDDigits(rawTo)
			if strings.Contains(rawTo, "@") {
				jid, err = toJID(rawTo)
			} else {
				jid, err = toUserJID(digits)
			}
			if err != nil {
				entry.Error("recipient_parse_error to=%q err=%v", msg.To, err)
				return fmt.Errorf("invalid recipient %q: %w", msg.To, err)
			}
		}
	}
	entry.Info("dest_jid=%s", jid.String())

	requestsMentionAll := outgoingMessageRequestsMentionAll(msg)
	var ctxInfo *goE2E.ContextInfo
	if msg.Context != nil || len(msg.Mentions) > 0 || requestsMentionAll {
		var stanzaID string
		var participantFromID string
		if msg.Context != nil {
			switch {
			case msg.Context.StanzaId != "":
				stanzaID = msg.Context.StanzaId
			case msg.Context.Stanza_id != "":
				stanzaID = msg.Context.Stanza_id
			case msg.Context.ID != "":
				stanzaID = msg.Context.ID
			case msg.Context.MessageID != "":
				stanzaID = msg.Context.MessageID
			}
			if participant, messageID, ok := splitGroupParticipantMessageID(stanzaID); ok {
				participantFromID = participant
				stanzaID = messageID
			}
			entry.Debug(
				"reply_ctx_select stanza_id=%q ctx_in={id:%q message_id:%q stanzaId:%q stanza_id:%q participant:%q}",
				stanzaID,
				msg.Context.ID,
				msg.Context.MessageID,
				msg.Context.StanzaId,
				msg.Context.Stanza_id,
				msg.Context.Participant,
			)
		}
		var mList []string
		if msg.Context != nil {
			mList = append(mList, msg.Context.MentionedJid...)
			mList = append(mList, msg.Context.Mentions...)
		}
		mList = append(mList, msg.Mentions...)
		mentioned := normalizeMentionTargets(mList)
		if requestsMentionAll && strings.HasSuffix(jid.Server, "g.us") {
			participants, err := resolveMentionAllJIDs(ctx, cli, jid)
			if err != nil {
				entry.Error("mention_all resolve failed group=%s err=%v", jid.String(), err)
				return fmt.Errorf("failed to resolve group participants for mention_all: %w", err)
			}
			mentioned = appendUniqueStrings(mentioned, participants...)
			entry.Info("mention_all enabled group=%s participants=%d", jid.String(), len(participants))
		} else if requestsMentionAll {
			entry.Debug("mention_all ignored for non-group destination=%s", jid.String())
		}
		if stanzaID != "" || len(mentioned) > 0 {
			ci := &goE2E.ContextInfo{}
			if stanzaID != "" {
				ci.StanzaID = proto.String(stanzaID)
				if msg.Context != nil && msg.Context.Participant != "" {
					if pj, err := toUserJID(msg.Context.Participant); err == nil {
						ci.Participant = proto.String(pj.String())
					}
				}
				if (ci.Participant == nil || *ci.Participant == "") && participantFromID != "" {
					if pj, err := toUserJID(participantFromID); err == nil {
						ci.Participant = proto.String(pj.String())
					}
				}
				qt := ""
				if msg.Context != nil {
					qt = strings.TrimSpace(msg.Context.QuotedText)
				}
				// Include quoted content when available; otherwise include empty to trigger reply UI (Wuzapi behavior)
				if qt != "" {
					ci.QuotedMessage = &goE2E.Message{Conversation: proto.String(qt)}
				} else if ci.QuotedMessage == nil {
					ci.QuotedMessage = &goE2E.Message{Conversation: proto.String("")}
				}
				// Ensure participant presence for 1:1 replies.
				if (ci.Participant == nil || *ci.Participant == "") && !strings.HasSuffix(jid.Server, "g.us") {
					ci.Participant = proto.String(jid.String())
				}
			}
			if len(mentioned) > 0 {
				ci.MentionedJID = mentioned
			}
			ctxInfo = ci
			partOut := ""
			if ci.Participant != nil {
				partOut = *ci.Participant
			}
			qtLen := 0
			if msg.Context != nil {
				qtLen = len(strings.TrimSpace(msg.Context.QuotedText))
			}
			entry.Debug(
				"reply_ctx_built stanza_id=%q participant=%q mentioned=%d quoted_text_len=%d dest_is_group=%t dest_jid=%s",
				stanzaID,
				partOut,
				len(mentioned),
				qtLen,
				strings.HasSuffix(jid.Server, "g.us"),
				jid.String(),
			)
		}
	}

	switch strings.ToLower(msg.Type) {
	case "text":
		body := strings.TrimSpace(msg.Body)
		if body == "" && msg.Text != nil {
			body = strings.TrimSpace(msg.Text.Body)
		}
		if body == "" {
			return fmt.Errorf("empty text body")
		}
		var wire *goE2E.Message
		if ctxInfo != nil {
			wire = &goE2E.Message{
				ExtendedTextMessage: &goE2E.ExtendedTextMessage{
					Text:        proto.String(body),
					ContextInfo: ctxInfo,
				},
			}
		} else {
			wire = &goE2E.Message{Conversation: proto.String(body)}
		}
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send text failed: %v", err)
			return err
		}
		entry.Info("text sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "image":
		var link, caption, mediaID string
		if msg.Image != nil {
			link = strings.TrimSpace(msg.Image.Link)
			mediaID = strings.TrimSpace(msg.Image.ID)
			caption = strings.TrimSpace(msg.Image.Caption)
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		if caption == "" {
			caption = strings.TrimSpace(msg.Caption)
		}
		var (
			data []byte
			mime string
		)
		if link == "" && mediaID != "" {
			stored, meta, err := m.LoadMediaContent(ctx, sessionID, mediaID)
			if err != nil {
				return fmt.Errorf("failed to load media id %q: %w", mediaID, err)
			}
			data = stored
			mime = meta.MIMEType
			if strings.TrimSpace(mime) == "" {
				mime = http.DetectContentType(data)
			}
		} else {
			if link == "" {
				return fmt.Errorf("image requires a link/media_url or image.id")
			}
			var err error
			data, mime, err = downloadBytes(ctx, link)
			if err != nil {
				return fmt.Errorf("failed to download media: %w", err)
			}
		}
		uploaded, err := cli.Upload(ctx, data, whatsmeow.MediaImage)
		if err != nil {
			entry.Error("image upload failed: %v", err)
			return err
		}
		imgMsg := &goE2E.ImageMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			Mimetype:      proto.String(mime),
		}
		if caption != "" {
			imgMsg.Caption = proto.String(caption)
		}
		if ctxInfo != nil {
			imgMsg.ContextInfo = ctxInfo
		}
		wire := &goE2E.Message{ImageMessage: imgMsg}
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send image failed: %v", err)
			return err
		}
		entry.Info("image sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "document":
		var link, caption, filename, mediaID string
		if msg.Document != nil {
			link = strings.TrimSpace(msg.Document.Link)
			mediaID = strings.TrimSpace(msg.Document.ID)
			caption = strings.TrimSpace(msg.Document.Caption)
			filename = strings.TrimSpace(msg.Document.FileName)
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		if caption == "" {
			caption = strings.TrimSpace(msg.Caption)
		}
		var (
			data []byte
			mime string
		)
		if link == "" && mediaID != "" {
			stored, meta, err := m.LoadMediaContent(ctx, sessionID, mediaID)
			if err != nil {
				return fmt.Errorf("failed to load media id %q: %w", mediaID, err)
			}
			data = stored
			mime = meta.MIMEType
			if strings.TrimSpace(mime) == "" {
				mime = http.DetectContentType(data)
			}
			if filename == "" {
				filename = strings.TrimSpace(meta.Filename)
			}
		} else {
			if link == "" {
				return fmt.Errorf("document requires a link/media_url or document.id")
			}
			var err error
			data, mime, err = downloadBytes(ctx, link)
			if err != nil {
				return fmt.Errorf("failed to download media: %w", err)
			}
		}
		if filename == "" {
			parts := strings.Split(strings.TrimRight(link, "/"), "/")
			filename = parts[len(parts)-1]
			if filename == "" && mediaID != "" {
				filename = mediaID
			}
		}
		uploaded, err := cli.Upload(ctx, data, whatsmeow.MediaDocument)
		if err != nil {
			entry.Error("document upload failed: %v", err)
			return err
		}
		docMsg := &goE2E.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			Mimetype:      proto.String(mime),
			FileName:      proto.String(filename),
		}
		if caption != "" {
			docMsg.Caption = proto.String(caption)
		}
		if ctxInfo != nil {
			docMsg.ContextInfo = ctxInfo
		}
		wire := &goE2E.Message{DocumentMessage: docMsg}
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send document failed: %v", err)
			return err
		}
		entry.Info("document sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "video":
		var link, caption, mediaID string
		if msg.Video != nil {
			link = strings.TrimSpace(msg.Video.Link)
			mediaID = strings.TrimSpace(msg.Video.ID)
			caption = strings.TrimSpace(msg.Video.Caption)
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		if caption == "" {
			caption = strings.TrimSpace(msg.Caption)
		}
		var (
			data []byte
			mime string
		)
		if link == "" && mediaID != "" {
			stored, meta, err := m.LoadMediaContent(ctx, sessionID, mediaID)
			if err != nil {
				return fmt.Errorf("failed to load media id %q: %w", mediaID, err)
			}
			data = stored
			mime = meta.MIMEType
			if strings.TrimSpace(mime) == "" {
				mime = http.DetectContentType(data)
			}
		} else {
			if link == "" {
				return fmt.Errorf("video requires a link/media_url or video.id")
			}
			var err error
			data, mime, err = downloadBytes(ctx, link)
			if err != nil {
				return fmt.Errorf("failed to download media: %w", err)
			}
		}
		uploaded, err := cli.Upload(ctx, data, whatsmeow.MediaVideo)
		if err != nil {
			entry.Error("video upload failed: %v", err)
			return err
		}
		vidMsg := &goE2E.VideoMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			Mimetype:      proto.String(mime),
		}
		if caption != "" {
			vidMsg.Caption = proto.String(caption)
		}
		if ctxInfo != nil {
			vidMsg.ContextInfo = ctxInfo
		}
		wire := &goE2E.Message{VideoMessage: vidMsg}
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send video failed: %v", err)
			return err
		}
		entry.Info("video sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "audio":
		// Mirror Evolution API: accepts URL (link/media_url) or previously uploaded media id.
		var src, mediaID string
		if msg.Audio != nil {
			src = strings.TrimSpace(msg.Audio.Link)
			mediaID = strings.TrimSpace(msg.Audio.ID)
		}
		if src == "" {
			src = strings.TrimSpace(msg.MediaURL)
		}
		var (
			original   []byte
			header     http.Header
			sourceMime string
		)
		if src == "" && mediaID != "" {
			data, meta, err := m.LoadMediaContent(ctx, sessionID, mediaID)
			if err != nil {
				return fmt.Errorf("failed to load media id %q: %w", mediaID, err)
			}
			original = data
			sourceMime = strings.TrimSpace(meta.MIMEType)
		} else {
			if src == "" {
				return fmt.Errorf("audio requires a URL (http/https) or audio.id")
			}
			if !(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
				return fmt.Errorf("audio link must be an http(s) URL")
			}
			data, mime, hdr, err := downloadBytesWithHeader(ctx, src)
			if err != nil {
				return fmt.Errorf("failed to fetch audio: %w", err)
			}
			original = data
			sourceMime = strings.TrimSpace(mime)
			header = hdr
		}

		ptt := m.defaultAudioPTT
		if msg.Audio != nil && msg.Audio.PTT != nil {
			ptt = *msg.Audio.PTT
		}

		var (
			uploadBytes []byte
			secs        uint32
			waveform    []byte
		)

		if ptt {
			// Prefer ffmpeg-based conversion with waveform for iOS WA
			uploadBytes, secs, waveform, err = convertToOggOpusFFMPEGWithWaveform(original)
			if err != nil {
				entry.Error("ffmpeg convert+waveform failed, falling back to gopus: %v", err)
				// Fallback to Go encoder
				uploadBytes, secs, waveform, err = convertMP3ToOGG(original)
				if err != nil {
					entry.Error("audio convert error: %v", err)
					return fmt.Errorf("audio conversion failed: %w", err)
				}
			}
		} else {
			uploadBytes = original
		}

		if secs == 0 {
			var payloadSeconds *uint32
			if msg.Audio != nil {
				payloadSeconds = msg.Audio.Seconds
			}
			secs = pickSeconds(payloadSeconds, header, ptt)
		}

		mime := sourceMime
		if ptt {
			mime = "audio/ogg; codecs=opus"
		} else {
			if mime == "" || mime == "application/octet-stream" {
				mime = http.DetectContentType(uploadBytes)
			}
			mime = normalizeAudioMime(src, mime, ptt)
		}

		entry.Debug("audio meta (pre-upload): bytes=%d secs=%d ptt=%v dest=%s", len(uploadBytes), secs, ptt, jid.String())

		uploaded, err := cli.Upload(ctx, uploadBytes, whatsmeow.MediaAudio)
		if err != nil {
			entry.Error("audio upload failed: %v", err)
			return err
		}

		audioMsg := &goE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(uploadBytes))),
			Mimetype:      proto.String(mime),
			PTT:           &ptt,
		}
		if secs > 0 {
			audioMsg.Seconds = proto.Uint32(secs)
		}
		if len(waveform) > 0 {
			audioMsg.Waveform = waveform
		}
		if ctxInfo != nil {
			audioMsg.ContextInfo = ctxInfo
		}

		wire := &goE2E.Message{AudioMessage: audioMsg}
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send audio failed: %v", err)
			return err
		}
		entry.Info("audio sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "sticker":
		var link, mediaID string
		if msg.Sticker != nil {
			link = strings.TrimSpace(msg.Sticker.Link)
			mediaID = strings.TrimSpace(msg.Sticker.ID)
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		var (
			data []byte
			mime string
		)
		if link == "" && mediaID != "" {
			stored, meta, err := m.LoadMediaContent(ctx, sessionID, mediaID)
			if err != nil {
				return fmt.Errorf("failed to load media id %q: %w", mediaID, err)
			}
			data = stored
			mime = meta.MIMEType
			if strings.TrimSpace(mime) == "" {
				mime = http.DetectContentType(data)
			}
		} else {
			if link == "" {
				return fmt.Errorf("sticker requires a link/media_url or sticker.id")
			}
			var err error
			data, mime, err = downloadBytes(ctx, link)
			if err != nil {
				return fmt.Errorf("failed to download media: %w", err)
			}
		}
		uploaded, err := cli.Upload(ctx, data, whatsmeow.MediaStickerPack)
		if err != nil {
			entry.Error("sticker upload failed: %v", err)
			return err
		}
		stickerMsg := &goE2E.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			Mimetype:      proto.String(mime),
		}
		if ctxInfo != nil {
			stickerMsg.ContextInfo = ctxInfo
		}
		wire := &goE2E.Message{StickerMessage: stickerMsg}
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send sticker failed: %v", err)
			return err
		}
		entry.Info("sticker sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "contacts":
		if len(msg.Contacts) == 0 {
			return fmt.Errorf("contacts requires at least one contact entry")
		}
		// Build vCards
		var displayName string
		contactMsgs := make([]*goE2E.ContactMessage, 0, len(msg.Contacts))
		for i, c := range msg.Contacts {
			name := strings.TrimSpace(c.Name.FormattedName)
			if name == "" {
				name = "Contact"
			}
			if i == 0 {
				displayName = name
			}
			// Build TEL lines
			var telLines []string
			for _, p := range c.Phones {
				phone := strings.TrimSpace(p.Phone)
				waid := strings.TrimSpace(p.WAID)
				if phone == "" {
					continue
				}
				if waid != "" {
					telLines = append(telLines, fmt.Sprintf("TEL;type=CELL;type=VOICE;waid=%s:%s", waid, phone))
				} else {
					telLines = append(telLines, fmt.Sprintf("TEL;type=CELL;type=VOICE:%s", phone))
				}
			}
			if len(telLines) == 0 {
				continue
			}
			vcard := "BEGIN:VCARD\n" +
				"VERSION:3.0\n" +
				fmt.Sprintf("N:%s\n", name) +
				strings.Join(telLines, "\n") + "\n" +
				"END:VCARD"
			contactMsgs = append(contactMsgs, &goE2E.ContactMessage{
				DisplayName: proto.String(name),
				Vcard:       proto.String(vcard),
			})
		}
		if len(contactMsgs) == 0 {
			return fmt.Errorf("contacts payload has no valid phone entries")
		}
		arr := &goE2E.ContactsArrayMessage{
			DisplayName: proto.String(displayName),
			Contacts:    contactMsgs,
		}
		if ctxInfo != nil {
			// ContactsArrayMessage does not have ContextInfo; attach at Message level via future if supported.
		}
		wire := &goE2E.Message{ContactsArrayMessage: arr}
		entry.Info("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send contacts failed: %v", err)
			return err
		}
		entry.Info("contacts sent to %s (wa_msg_id=%s) count=%d", jid.String(), resp.ID, len(contactMsgs))
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "interactive":
		if msg.Interactive == nil {
			return fmt.Errorf("interactive requires an interactive payload")
		}
		switch strings.ToLower(strings.TrimSpace(msg.Interactive.Type)) {
		case "button":
			wire, err := buildInteractiveButtonMessage(msg.Interactive, ctxInfo)
			if err != nil {
				return err
			}
			sendExtra := interactiveButtonSendExtra()
			entry.Debug("interactive_button_transport=interactive_native_flow biz_node=interactive/native_flow name=mixed")
			entry.Debug("wire_payload=%s", marshalForLog(wire))
			resp, err := sendWithExtra(ctx, cli, jid, wire, msg.MessageID, sendExtra)
			if err != nil {
				entry.Error("send interactive button failed: %v", err)
				return err
			}
			entry.Info("interactive button sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
			return m.emitCloudSent(sessionID, jid, resp.ID)
		case "list":
			wire, err := buildInteractiveListMessage(msg.Interactive, ctxInfo)
			if err != nil {
				return err
			}
			entry.Debug("interactive_list_transport=legacy_list_message biz_node_type=product_list")
			entry.Debug("wire_payload=%s", marshalForLog(wire))
			resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
			if err != nil {
				entry.Error("send interactive list failed: %v", err)
				return err
			}
			entry.Info("interactive list sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
			return m.emitCloudSent(sessionID, jid, resp.ID)
		default:
			return fmt.Errorf("unsupported interactive type: %s", msg.Interactive.Type)
		}

	case "payment", "order_details":
		payment := msg.paymentContent()
		if payment == nil {
			return fmt.Errorf("%s requires a payment payload", strings.ToLower(strings.TrimSpace(msg.Type)))
		}
		entry.Debug(
			"payment_input reference_id=%q body_len=%d has_header_document=%t has_pix=%t has_boleto=%t",
			strings.TrimSpace(payment.ReferenceID),
			len(strings.TrimSpace(payment.Body)),
			payment.HeaderDocument != nil,
			payment.PIXStaticCode != nil,
			payment.Boleto != nil,
		)

		var headerDoc *goE2E.DocumentMessage
		if payment.HeaderDocument != nil {
			var err error
			headerDoc, err = m.buildPaymentHeaderDocument(ctx, cli, sessionID, payment.HeaderDocument)
			if err != nil {
				return err
			}
		}

		wire, sendExtra, err := buildOrderDetailsMessage(payment, headerDoc, ctxInfo)
		if err != nil {
			return err
		}
		entry.Debug("payment_transport=interactive_native_flow native_flow_name=order_details")
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithExtra(ctx, cli, jid, wire, msg.MessageID, sendExtra)
		if err != nil {
			entry.Error("send payment failed: %v", err)
			return err
		}
		entry.Info("payment sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	default:
		return fmt.Errorf("unsupported message type: %s", msg.Type)
	}
}

func (m *ClientManager) buildPaymentHeaderDocument(ctx context.Context, cli *whatsmeow.Client, sessionID string, in *DocumentContent) (*goE2E.DocumentMessage, error) {
	if in == nil {
		return nil, nil
	}

	link := strings.TrimSpace(in.Link)
	mediaID := strings.TrimSpace(in.ID)
	filename := strings.TrimSpace(in.FileName)

	var data []byte
	if link == "" && mediaID != "" {
		stored, meta, err := m.LoadMediaContent(ctx, sessionID, mediaID)
		if err != nil {
			return nil, fmt.Errorf("failed to load payment header document id %q: %w", mediaID, err)
		}
		data = stored
		if filename == "" {
			filename = strings.TrimSpace(meta.Filename)
		}
	} else {
		if link == "" {
			return nil, fmt.Errorf("payment header_document requires a link or id")
		}
		var err error
		data, _, err = downloadBytes(ctx, link)
		if err != nil {
			return nil, fmt.Errorf("failed to download payment header document: %w", err)
		}
	}

	if filename == "" {
		filename = documentFilenameFromLink(link, mediaID)
	}
	if filename == "" {
		filename = "boleto.pdf"
	}

	uploaded, err := cli.Upload(ctx, data, whatsmeow.MediaDocument)
	if err != nil {
		return nil, fmt.Errorf("payment header document upload failed: %w", err)
	}

	return &goE2E.DocumentMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(data))),
		Mimetype:      proto.String("application/pdf"),
		Title:         proto.String(filename),
		FileName:      proto.String(filename),
	}, nil
}

func buildOrderDetailsMessage(in *PaymentContent, headerDoc *goE2E.DocumentMessage, ctxInfo *goE2E.ContextInfo) (*goE2E.Message, *whatsmeow.SendRequestExtra, error) {
	if in == nil {
		return nil, nil, fmt.Errorf("missing payment payload")
	}

	params, err := normalizePaymentContent(in)
	if err != nil {
		return nil, nil, err
	}

	paymentSettings, err := buildPaymentSettings(params)
	if err != nil {
		return nil, nil, err
	}

	payload := map[string]any{
		"reference_id":         params.ReferenceID,
		"type":                 params.Type,
		"payment_type":         params.PaymentType,
		"currency":             params.Currency,
		"total_amount":         params.TotalAmount,
		"order":                params.Order,
		"payment_settings":     paymentSettings,
		"share_payment_status": params.SharePaymentStatus,
		"referral":             params.Referral,
	}
	buttonParamsJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal order_details payload: %w", err)
	}

	messageSecret := make([]byte, 32)
	if _, err := rand.Read(messageSecret); err != nil {
		return nil, nil, fmt.Errorf("generate order_details message secret: %w", err)
	}

	interactiveMsg := &goE2E.InteractiveMessage{
		Body: &goE2E.InteractiveMessage_Body{
			Text: proto.String(params.Body),
		},
		InteractiveMessage: &goE2E.InteractiveMessage_NativeFlowMessage_{
			NativeFlowMessage: &goE2E.InteractiveMessage_NativeFlowMessage{
				Buttons: []*goE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
					{
						Name:             proto.String("review_and_pay"),
						ButtonParamsJSON: proto.String(string(buttonParamsJSON)),
					},
				},
			},
		},
		ContextInfo: ctxInfo,
	}
	if headerDoc != nil {
		interactiveMsg.Header = &goE2E.InteractiveMessage_Header{
			HasMediaAttachment: proto.Bool(true),
			Media: &goE2E.InteractiveMessage_Header_DocumentMessage{
				DocumentMessage: headerDoc,
			},
		}
	}

	bizNodes := []waBinary.Node{
		{
			Tag: "biz",
			Attrs: waBinary.Attrs{
				"native_flow_name": "order_details",
			},
		},
	}

	return &goE2E.Message{
			MessageContextInfo: &goE2E.MessageContextInfo{
				MessageSecret: messageSecret,
			},
			InteractiveMessage: interactiveMsg,
		}, &whatsmeow.SendRequestExtra{
			AdditionalNodes: &bizNodes,
		}, nil
}

func normalizePaymentContent(in *PaymentContent) (*PaymentContent, error) {
	out := *in
	out.ReferenceID = strings.TrimSpace(out.ReferenceID)
	out.Currency = strings.TrimSpace(out.Currency)
	out.Type = strings.TrimSpace(out.Type)
	out.PaymentType = strings.TrimSpace(out.PaymentType)
	out.Body = strings.TrimSpace(out.Body)
	out.ItemName = strings.TrimSpace(out.ItemName)
	out.Referral = strings.TrimSpace(out.Referral)
	if out.HeaderDocument != nil {
		docCopy := *out.HeaderDocument
		docCopy.Link = strings.TrimSpace(docCopy.Link)
		docCopy.ID = strings.TrimSpace(docCopy.ID)
		docCopy.FileName = strings.TrimSpace(docCopy.FileName)
		out.HeaderDocument = &docCopy
	}

	if out.TotalAmount.Value <= 0 {
		return nil, fmt.Errorf("payment requires total_amount.value > 0")
	}
	if out.TotalAmount.Offset <= 0 {
		out.TotalAmount.Offset = 100
	}
	if out.ReferenceID == "" {
		out.ReferenceID = fmt.Sprintf("ORDER-%d", time.Now().Unix())
	}
	if out.Currency == "" {
		out.Currency = "BRL"
	}
	if out.Type == "" {
		out.Type = "physical-goods"
	}
	if out.PaymentType == "" {
		out.PaymentType = "br"
	}
	if out.Body == "" {
		out.Body = "Segue a cobrança. Você pode pagar pelo botão abaixo."
	}
	if out.Referral == "" {
		out.Referral = "chat_attachment"
	}

	if out.Order != nil {
		orderCopy := *out.Order
		orderCopy.Status = strings.TrimSpace(orderCopy.Status)
		items := make([]PaymentOrderItem, 0, len(orderCopy.Items))
		for _, item := range orderCopy.Items {
			item.Name = strings.TrimSpace(item.Name)
			item.RetailerID = strings.TrimSpace(item.RetailerID)
			if item.Name == "" {
				continue
			}
			if item.Quantity <= 0 {
				item.Quantity = 1
			}
			if item.Amount.Value <= 0 {
				item.Amount = out.TotalAmount
			}
			if item.Amount.Offset <= 0 {
				item.Amount.Offset = out.TotalAmount.Offset
			}
			items = append(items, item)
		}
		orderCopy.Items = items
		if orderCopy.Status == "" {
			orderCopy.Status = "pending"
		}
		if orderCopy.Subtotal.Value <= 0 {
			orderCopy.Subtotal = out.TotalAmount
		}
		if orderCopy.Subtotal.Offset <= 0 {
			orderCopy.Subtotal.Offset = out.TotalAmount.Offset
		}
		out.Order = &orderCopy
	}

	if out.Order == nil || len(out.Order.Items) == 0 {
		itemName := out.ItemName
		if itemName == "" {
			itemName = "Cobrança"
		}
		quantity := out.Quantity
		if quantity <= 0 {
			quantity = 1
		}
		out.Order = &PaymentOrder{
			Status: "pending",
			Items: []PaymentOrderItem{
				{
					Name:     itemName,
					Amount:   out.TotalAmount,
					Quantity: quantity,
				},
			},
			Subtotal: out.TotalAmount,
		}
	}

	return &out, nil
}

func buildPaymentSettings(in *PaymentContent) ([]any, error) {
	settings := make([]any, 0, 2)
	displayName := paymentDisplayName(in)

	if in.PIXStaticCode != nil {
		merchant := strings.TrimSpace(in.PIXStaticCode.MerchantName)
		key := strings.TrimSpace(in.PIXStaticCode.Key)
		keyType := strings.ToUpper(strings.TrimSpace(in.PIXStaticCode.KeyType))
		if key == "" {
			return nil, fmt.Errorf("payment pix_static_code requires key")
		}
		if merchant == "" {
			merchant = displayName
		}
		if keyType == "" {
			keyType = "EVP"
		}
		settings = append(settings, map[string]any{
			"type": "pix_static_code",
			"pix_static_code": map[string]string{
				"merchant_name": merchant,
				"key":           key,
				"key_type":      keyType,
			},
		})
	}

	if in.Boleto != nil {
		digitableLine := strings.TrimSpace(in.Boleto.DigitableLine)
		if digitableLine == "" {
			return nil, fmt.Errorf("payment boleto requires digitable_line")
		}
		settings = append(settings, map[string]any{
			"type": "boleto",
			"boleto": map[string]string{
				"digitable_line": digitableLine,
			},
		})
	}

	if len(settings) == 0 {
		return nil, fmt.Errorf("payment requires pix_static_code and/or boleto")
	}
	return settings, nil
}

func paymentDisplayName(in *PaymentContent) string {
	if in == nil {
		return "Cobrança"
	}
	if name := strings.TrimSpace(in.ItemName); name != "" {
		return name
	}
	if in.Order != nil {
		for _, item := range in.Order.Items {
			if name := strings.TrimSpace(item.Name); name != "" {
				return name
			}
		}
	}
	return "Cobrança"
}

func documentFilenameFromLink(link, fallback string) string {
	link = strings.TrimSpace(link)
	if link != "" {
		parts := strings.Split(strings.TrimRight(link, "/"), "/")
		if len(parts) > 0 {
			name := strings.TrimSpace(parts[len(parts)-1])
			if idx := strings.IndexByte(name, '?'); idx >= 0 {
				name = name[:idx]
			}
			if name != "" {
				return name
			}
		}
	}
	return strings.TrimSpace(fallback)
}

func buildInteractiveButtonMessage(in *InteractiveContent, ctxInfo *goE2E.ContextInfo) (*goE2E.Message, error) {
	if in == nil {
		return nil, fmt.Errorf("missing interactive payload")
	}
	if in.Action == nil || len(in.Action.Buttons) == 0 {
		return nil, fmt.Errorf("interactive button requires action.buttons")
	}
	if len(in.Action.Buttons) > 3 {
		return nil, fmt.Errorf("interactive button supports at most 3 reply buttons")
	}

	body := interactiveTextValue(in.Body)
	if body == "" {
		return nil, fmt.Errorf("interactive button requires body.text")
	}

	buttons := make([]*goE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton, 0, len(in.Action.Buttons))
	for i, button := range in.Action.Buttons {
		if strings.ToLower(strings.TrimSpace(button.Type)) != "reply" {
			return nil, fmt.Errorf("unsupported interactive button type: %s", button.Type)
		}
		title := strings.TrimSpace(button.Reply.Title)
		if title == "" {
			return nil, fmt.Errorf("interactive button %d requires reply.title", i)
		}
		id := strings.TrimSpace(button.Reply.ID)
		if id == "" {
			id = title
		}
		params, err := json.Marshal(map[string]string{
			"display_text": title,
			"id":           id,
		})
		if err != nil {
			return nil, err
		}
		buttons = append(buttons, &goE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
			Name:             proto.String("quick_reply"),
			ButtonParamsJSON: proto.String(string(params)),
		})
	}

	messageParamsJSON := `{"from":"api","templateId":""}`
	interactiveMsg := &goE2E.InteractiveMessage{
		Body: &goE2E.InteractiveMessage_Body{
			Text: proto.String(body),
		},
		InteractiveMessage: &goE2E.InteractiveMessage_NativeFlowMessage_{
			NativeFlowMessage: &goE2E.InteractiveMessage_NativeFlowMessage{
				Buttons:           buttons,
				MessageParamsJSON: proto.String(messageParamsJSON),
				MessageVersion:    proto.Int32(1),
			},
		},
		ContextInfo: ctxInfo,
	}
	if footer := interactiveTextValue(in.Footer); footer != "" {
		interactiveMsg.Footer = &goE2E.InteractiveMessage_Footer{Text: proto.String(footer)}
	}
	if header := interactiveTextValue(in.Header); header != "" {
		interactiveMsg.Header = &goE2E.InteractiveMessage_Header{
			Title:              proto.String(header),
			HasMediaAttachment: proto.Bool(false),
		}
	}
	return &goE2E.Message{InteractiveMessage: interactiveMsg}, nil
}

func interactiveButtonSendExtra() *whatsmeow.SendRequestExtra {
	nodes := []waBinary.Node{
		{
			Tag: "biz",
			Content: []waBinary.Node{
				{
					Tag: "interactive",
					Attrs: waBinary.Attrs{
						"type": "native_flow",
						"v":    "1",
					},
					Content: []waBinary.Node{
						{
							Tag: "native_flow",
							Attrs: waBinary.Attrs{
								"v":    "9",
								"name": "mixed",
							},
						},
					},
				},
			},
		},
	}
	return &whatsmeow.SendRequestExtra{AdditionalNodes: &nodes}
}

func buildInteractiveListMessage(in *InteractiveContent, ctxInfo *goE2E.ContextInfo) (*goE2E.Message, error) {
	if in == nil {
		return nil, fmt.Errorf("missing interactive payload")
	}
	if in.Action == nil || len(in.Action.Sections) == 0 {
		return nil, fmt.Errorf("interactive list requires action.sections")
	}

	body := interactiveTextValue(in.Body)
	if body == "" {
		return nil, fmt.Errorf("interactive list requires body.text")
	}
	buttonText := strings.TrimSpace(in.Action.Button)
	if buttonText == "" {
		return nil, fmt.Errorf("interactive list requires action.button")
	}

	sections := make([]*goE2E.ListMessage_Section, 0, len(in.Action.Sections))
	for sectionIndex, section := range in.Action.Sections {
		if len(section.Rows) == 0 {
			return nil, fmt.Errorf("interactive list section %d requires rows", sectionIndex)
		}
		sectionTitle := strings.TrimSpace(section.Title)
		if sectionTitle == "" {
			sectionTitle = " "
		}
		rows := make([]*goE2E.ListMessage_Row, 0, len(section.Rows))
		for rowIndex, row := range section.Rows {
			rowTitle := strings.TrimSpace(row.Title)
			if rowTitle == "" {
				return nil, fmt.Errorf("interactive list section %d row %d requires title", sectionIndex, rowIndex)
			}
			rowID := strings.TrimSpace(row.ID)
			if rowID == "" {
				rowID = fmt.Sprintf("row_%d_%d", sectionIndex, rowIndex)
			}
			rows = append(rows, &goE2E.ListMessage_Row{
				Title:       proto.String(rowTitle),
				Description: proto.String(strings.TrimSpace(row.Description)),
				RowID:       proto.String(rowID),
			})
		}
		sections = append(sections, &goE2E.ListMessage_Section{
			Title: proto.String(sectionTitle),
			Rows:  rows,
		})
	}

	title := interactiveTextValue(in.Header)
	if title == "" {
		title = body
	}
	footer := interactiveTextValue(in.Footer)
	if footer == "" {
		footer = buttonText
	}
	listType := goE2E.ListMessage_SINGLE_SELECT
	listMessage := &goE2E.ListMessage{
		Title:       proto.String(title),
		Description: proto.String(body),
		ButtonText:  proto.String(buttonText),
		FooterText:  proto.String(footer),
		ListType:    &listType,
		Sections:    sections,
		ContextInfo: ctxInfo,
	}
	return &goE2E.Message{ListMessage: listMessage}, nil
}

func interactiveTextValue(text *InteractiveText) string {
	if text == nil {
		return ""
	}
	return strings.TrimSpace(text.Text)
}

func (m OutgoingMessage) wantsMentionAll() bool {
	return m.MentionAll || m.MentionAllCompat
}

func outgoingMessageRequestsMentionAll(msg OutgoingMessage) bool {
	if msg.wantsMentionAll() {
		return true
	}
	for _, candidate := range mentionAllTextCandidates(msg) {
		if hasMentionAllPrefix(candidate) {
			return true
		}
	}
	return false
}

func mentionAllTextCandidates(msg OutgoingMessage) []string {
	out := []string{
		msg.Body,
		msg.Caption,
	}
	if msg.Text != nil {
		out = append(out, msg.Text.Body)
	}
	if msg.Image != nil {
		out = append(out, msg.Image.Caption)
	}
	if msg.Document != nil {
		out = append(out, msg.Document.Caption)
	}
	if msg.Video != nil {
		out = append(out, msg.Video.Caption)
	}
	if msg.Audio != nil {
		out = append(out, msg.Audio.Caption)
	}
	return out
}

func hasMentionAllPrefix(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	lower := strings.ToLower(text)
	for _, prefix := range []string{"@todos", "@all", "@everyone"} {
		if lower == prefix {
			return true
		}
		if !strings.HasPrefix(lower, prefix) {
			continue
		}
		if len(lower) == len(prefix) {
			return true
		}
		next := rune(lower[len(prefix)])
		if unicode.IsSpace(next) || strings.ContainsRune(".,:;!?-_/\\", next) {
			return true
		}
	}
	return false
}

func normalizeMentionTargets(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		jid, err := toUserJID(value)
		if err != nil {
			continue
		}
		key := jid.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}
	return out
}

func appendUniqueStrings(base []string, values ...string) []string {
	if len(values) == 0 {
		return base
	}
	seen := make(map[string]struct{}, len(base)+len(values))
	for _, item := range base {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		seen[item] = struct{}{}
	}
	for _, item := range values {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		base = append(base, item)
	}
	return base
}

func resolveMentionAllJIDs(ctx context.Context, cli *whatsmeow.Client, groupJID types.JID) ([]string, error) {
	if cli == nil {
		return nil, fmt.Errorf("nil client")
	}
	groupInfo, err := cli.GetGroupInfo(ctx, groupJID)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(groupInfo.Participants))
	for _, participant := range groupInfo.Participants {
		jid := normalizedContactLookupJID(participant.JID)
		if jid.User == "" || jid.Server == "" {
			continue
		}
		out = appendUniqueStrings(out, jid.String())
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("group %s has no participants available for mentions", groupJID.String())
	}
	return out, nil
}

const (
	targetSampleRate = 16000
	opusFrameMs      = 20
	opusFrameSamples = targetSampleRate * opusFrameMs / 1000
)

// convertMP3ToOGG decodes MP3 bytes, converts them to mono Opus-in-OGG and returns the encoded audio.
// It also returns the audio duration in seconds and a WhatsApp-style waveform derived from the PCM data.
func convertMP3ToOGG(data []byte) ([]byte, uint32, []byte, error) {
	if len(data) == 0 {
		return nil, 0, nil, errors.New("empty audio payload")
	}

	dec, err := mp3.NewDecoder(bytes.NewReader(data))
	if err != nil {
		return nil, 0, nil, fmt.Errorf("mp3 decoder: %w", err)
	}
	srcRate := dec.SampleRate()
	if srcRate <= 0 {
		return nil, 0, nil, errors.New("invalid mp3 sample rate")
	}

	var pcmBuf bytes.Buffer
	if _, err := io.Copy(&pcmBuf, dec); err != nil {
		return nil, 0, nil, fmt.Errorf("mp3 decode: %w", err)
	}
	pcmBytes := pcmBuf.Bytes()
	if len(pcmBytes) < 2 {
		return nil, 0, nil, errors.New("decoded audio too short")
	}

	channels := 1
	if len(pcmBytes)%4 == 0 {
		channels = 2
	}

	samples := bytesToInt16(pcmBytes)
	if channels == 2 {
		samples = downmixStereoToMono(samples)
	}
	if len(samples) == 0 {
		return nil, 0, nil, errors.New("no pcm samples")
	}

	pcm16k := resampleLinearInt16(samples, srcRate, targetSampleRate)
	if len(pcm16k) < opusFrameSamples {
		pcm16k = padZeros(pcm16k, opusFrameSamples)
	}

	secs := uint32(int(math.Round(float64(len(pcm16k)) / float64(targetSampleRate))))
	if secs == 0 {
		secs = 1
	}
	waveform := buildWaveform(pcm16k)

	var oggBuf bytes.Buffer
	ogg, err := oggwriter.NewWith(&oggBuf, uint32(targetSampleRate), 1)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("oggwriter: %w", err)
	}
	enc, err := gopus.NewEncoder(targetSampleRate, 1, gopus.Voip)
	if err != nil {
		ogg.Close()
		return nil, 0, nil, fmt.Errorf("opus encoder: %w", err)
	}

	const maxOpusPacketSize = 4000
	sequence := uint16(0)
	timestamp := uint32(opusFrameSamples)
	for off := 0; off < len(pcm16k); off += opusFrameSamples {
		end := off + opusFrameSamples
		if end > len(pcm16k) {
			end = len(pcm16k)
		}
		frame := pcm16k[off:end]
		if len(frame) < opusFrameSamples {
			frame = padZeros(frame, opusFrameSamples)
		}
		pkt, err := enc.Encode(frame, opusFrameSamples, maxOpusPacketSize)
		if err != nil {
			ogg.Close()
			return nil, 0, nil, fmt.Errorf("opus encode: %w", err)
		}
		rtpPkt := &rtp.Packet{
			Header: rtp.Header{
				SequenceNumber: sequence,
				Timestamp:      timestamp,
				PayloadType:    111,
				Version:        2,
			},
			Payload: pkt,
		}
		if err := ogg.WriteRTP(rtpPkt); err != nil {
			ogg.Close()
			return nil, 0, nil, fmt.Errorf("ogg write: %w", err)
		}
		sequence++
		timestamp += uint32(opusFrameSamples)
	}
	if err := ogg.Close(); err != nil {
		return nil, 0, nil, fmt.Errorf("ogg close: %w", err)
	}

	out := make([]byte, oggBuf.Len())
	copy(out, oggBuf.Bytes())
	return out, secs, waveform, nil
}

// buildWaveform computes a 32-bar WhatsApp waveform (0..31) from PCM samples
func buildWaveform(samples []int16) []byte {
	const bars = 64
	if len(samples) == 0 {
		return nil
	}
	blockSize := len(samples) / bars
	if blockSize == 0 {
		blockSize = 1
	}
	filtered := make([]float64, bars)
	var maxAmp float64

	for i := 0; i < bars; i++ {
		start := i * blockSize
		end := start + blockSize

		if start >= len(samples) {
			filtered[i] = 0
			continue
		}
		if end > len(samples) {
			end = len(samples)
		}

		var sum float64
		chunkLen := end - start
		for _, s := range samples[start:end] {
			v := float64(s)
			if v < 0 {
				v = -v
			}
			sum += v
		}
		var avg float64
		if chunkLen > 0 {
			avg = sum / float64(chunkLen)
		}
		filtered[i] = avg
		if avg > maxAmp {
			maxAmp = avg
		}
	}
	out := make([]byte, bars)
	if maxAmp == 0 {
		return out
	}
	for i := 0; i < bars; i++ {
		val := (filtered[i] / maxAmp) * 100.0
		if val < 0 {
			val = 0
		}
		if val > 100 {
			val = 100
		}
		out[i] = byte(math.Round(val))
	}

	return out
}

// bytesToInt16 converts little-endian bytes into int16 samples.
func bytesToInt16(b []byte) []int16 {
	if len(b)%2 != 0 {
		b = b[:len(b)-1]
	}
	n := len(b) / 2
	out := make([]int16, n)
	for i := 0; i < n; i++ {
		out[i] = int16(binary.LittleEndian.Uint16(b[2*i:]))
	}
	return out
}

// downmixStereoToMono averages left/right channels to produce mono audio.
func downmixStereoToMono(stereo []int16) []int16 {
	if len(stereo) < 2 {
		return stereo
	}
	mono := make([]int16, 0, len(stereo)/2)
	for i := 0; i+1 < len(stereo); i += 2 {
		m := int(stereo[i])/2 + int(stereo[i+1])/2
		if m > math.MaxInt16 {
			m = math.MaxInt16
		}
		if m < math.MinInt16 {
			m = math.MinInt16
		}
		mono = append(mono, int16(m))
	}
	return mono
}

// resampleLinearInt16 performs simple linear interpolation resampling between arbitrary sample rates.
func resampleLinearInt16(in []int16, srcRate, dstRate int) []int16 {
	if srcRate <= 0 || dstRate <= 0 || len(in) == 0 || srcRate == dstRate {
		out := make([]int16, len(in))
		copy(out, in)
		return out
	}
	ratio := float64(dstRate) / float64(srcRate)
	outLen := int(math.Round(float64(len(in)) * ratio))
	if outLen < 1 {
		outLen = 1
	}
	out := make([]int16, outLen)
	for i := 0; i < outLen; i++ {
		pos := float64(i) / ratio
		j := int(math.Floor(pos))
		t := pos - float64(j)
		s0 := float64(sampleAt(in, j))
		s1 := float64(sampleAt(in, j+1))
		val := (1.0-t)*s0 + t*s1
		if val > math.MaxInt16 {
			val = math.MaxInt16
		}
		if val < math.MinInt16 {
			val = math.MinInt16
		}
		out[i] = int16(val)
	}
	return out
}

func sampleAt(samples []int16, idx int) int16 {
	if len(samples) == 0 {
		return 0
	}
	if idx < 0 {
		return samples[0]
	}
	if idx >= len(samples) {
		return samples[len(samples)-1]
	}
	return samples[idx]
}

func padZeros(in []int16, want int) []int16 {
	if len(in) >= want {
		return in
	}
	out := make([]int16, want)
	copy(out, in)
	return out
}

func pickSeconds(payload *uint32, hdr http.Header, ptt bool) uint32 {
	if payload != nil && *payload > 0 {
		return *payload
	}
	tryKeys := []string{"X-Duration-Seconds", "X-Content-Duration", "Content-Duration", "X-Audio-Duration"}
	for _, k := range tryKeys {
		if s := hdr.Get(k); s != "" {
			if v, err := strconv.ParseFloat(s, 64); err == nil && v > 0 {
				return uint32(math.Round(v))
			}
		}
	}
	if ptt {
		return 0
	}
	return 0
}

func toUserJID(to string) (types.JID, error) {
	to = strings.TrimSpace(to)
	if to == "" {
		return types.JID{}, fmt.Errorf("empty recipient")
	}
	if strings.Contains(to, "@") {
		j, err := types.ParseJID(to)
		if err != nil {
			return types.JID{}, err
		}
		if j.Server == "" {
			j.Server = types.DefaultUserServer
		}
		return j, nil
	}
	digits := digitsOnly(to)
	if digits == "" {
		return types.JID{}, fmt.Errorf("empty msisdn after normalization")
	}
	return types.NewJID(digits, types.DefaultUserServer), nil
}

func toJID(to string) (types.JID, error) {
	to = strings.TrimSpace(to)
	if to == "" {
		return types.JID{}, fmt.Errorf("empty recipient")
	}
	if strings.Contains(to, "@") {
		return types.ParseJID(to)
	}
	return toUserJID(to)
}

func downloadBytesWithHeader(ctx context.Context, url string) ([]byte, string, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", resp.Header, fmt.Errorf("failed to download media, status: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", resp.Header, err
	}
	mime := http.DetectContentType(data)
	if ct := resp.Header.Get("Content-Type"); ct != "" && ct != "application/octet-stream" {
		mime = ct
	}
	return data, mime, resp.Header, nil
}

func shouldSendAsPTT(link, mime string) bool {
	l := strings.ToLower(link)
	m := strings.ToLower(mime)
	if strings.HasSuffix(l, ".ogg") || strings.HasSuffix(l, ".oga") || strings.HasSuffix(l, ".opus") {
		return true
	}
	if strings.Contains(m, "audio/ogg") || strings.Contains(m, "application/ogg") || strings.Contains(m, "opus") {
		return true
	}
	return false
}

func normalizeAudioMime(link, mime string, ptt bool) string {
	m := strings.ToLower(mime)
	if strings.Contains(m, "mpga") {
		return "audio/mpeg"
	}
	if ptt {
		return "audio/ogg; codecs=opus"
	}
	return mime
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func downloadBytes(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download media, status: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	mime := http.DetectContentType(data)
	return data, mime, nil
}

// phoneNumberToJIDDigits normalizes an input phone (E.164-like or JID) into PN digits used in JIDs.
func phoneNumberToJIDDigits(value string) string {
	// If already a JID, extract user part before '@' or ':'
	in := strings.TrimSpace(value)
	if i := strings.IndexAny(in, "@:"); i >= 0 {
		in = in[:i]
	}
	number := digitsOnly(in)
	if len(number) < 3 {
		return number
	}
	// Apply Brazil mobile missing '9' heuristic once
	return jidToPhoneNumberDigits(number, true)
}

func jidToPhoneNumberDigits(number string, retry bool) string {
	if len(number) < 4 { // require at least 55 + AA
		return number
	}
	country := number[:2]
	if country != "55" {
		return number
	}
	area := number[2:4]
	subs := number[4:]
	// If BR and subscriber has 9 digits and starts with 9, drop leading '9' (WA PN often uses 8-digit subscriber)
	if len(subs) == 9 && subs[0] == '9' {
		return "55" + area + subs[1:]
	}
	// If BR and subscriber has only 8 digits and looks like mobile (starts 6..9), insert leading '9' once
	if len(subs) == 8 && subs[0] >= '6' && subs[0] <= '9' && retry {
		return "55" + area + "9" + subs
	}
	return number
}

// convertToOggOpusFFMPEG converts arbitrary audio input to mono OGG/Opus using ffmpeg.
// This matches common Baileys-compatible settings and maximizes iOS WhatsApp compatibility.
func convertToOggOpusFFMPEG(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, errors.New("empty audio payload")
	}
	args := []string{
		"-hide_banner", "-loglevel", "error",
		"-i", "pipe:0",
		"-vn",
		"-ac", "1",
		"-ar", "48000",
		"-c:a", "libopus",
		"-b:a", "24k",
		"-application", "voip",
		"-avoid_negative_ts", "make_zero",
		"-f", "ogg",
		"pipe:1",
	}
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg run: %w: %s", err, strings.TrimSpace(errBuf.String()))
	}
	if out.Len() == 0 {
		return nil, errors.New("ffmpeg produced empty output")
	}
	// Note: duration/waveform are not derived here; Seconds stays 0; Waveform omitted.
	return out.Bytes(), nil
}

// convertToOggOpusFFMPEGWithWaveform also derives waveform by decoding PCM via ffmpeg.
func convertToOggOpusFFMPEGWithWaveform(in []byte) ([]byte, uint32, []byte, error) {
	if len(in) == 0 {
		return nil, 0, nil, errors.New("empty audio payload")
	}
	// First pass: decode to mono 16k s16le to compute waveform and duration
	pcmArgs := []string{
		"-hide_banner", "-loglevel", "error",
		"-i", "pipe:0",
		"-vn",
		"-ac", "1",
		"-ar", "16000",
		"-f", "s16le",
		"pipe:1",
	}
	pcm, stderr, err := runFFMPEG(in, pcmArgs)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("ffmpeg pcm: %w: %s", err, strings.TrimSpace(stderr))
	}
	if len(pcm) < 2 {
		return nil, 0, nil, errors.New("ffmpeg pcm output empty")
	}
	samples := bytesToInt16(pcm)
	secs := uint32(int(math.Round(float64(len(samples)) / float64(16000))))
	if secs == 0 {
		secs = 1
	}
	waveform := buildWaveform(samples)

	// Second pass: encode to OGG/Opus at 48k mono VOIP
	oggArgs := []string{
		"-hide_banner", "-loglevel", "error",
		"-i", "pipe:0",
		"-vn",
		"-ac", "1",
		"-ar", "48000",
		"-c:a", "libopus",
		"-b:a", "24k",
		"-application", "voip",
		"-avoid_negative_ts", "make_zero",
		"-f", "ogg",
		"pipe:1",
	}
	ogg, stderr2, err := runFFMPEG(in, oggArgs)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("ffmpeg ogg: %w: %s", err, strings.TrimSpace(stderr2))
	}
	if len(ogg) == 0 {
		return nil, 0, nil, errors.New("ffmpeg produced empty ogg output")
	}
	return ogg, secs, waveform, nil
}

func runFFMPEG(in []byte, args []string) ([]byte, string, error) {
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.Bytes(), errBuf.String(), err
}
