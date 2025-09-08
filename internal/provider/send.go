package provider

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"your.org/provider-whatsmeow/internal/log"

	"go.mau.fi/whatsmeow"
	goE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// TextContent representa {"text":{"body":"..."}}
type TextContent struct {
	Body string `json:"body"`
}

type ImageContent struct {
	Link    string `json:"link"`
	Caption string `json:"caption,omitempty"`
}

type DocumentContent struct {
	Link     string `json:"link"`
	Caption  string `json:"caption,omitempty"`
	FileName string `json:"filename,omitempty"`
}

type AudioContent struct {
	Link    string  `json:"link"`
	Caption string  `json:"caption,omitempty"`
	PTT     *bool   `json:"ptt,omitempty"`
	Seconds *uint32 `json:"seconds,omitempty"`
}

type StickerContent struct {
	Link string `json:"link"`
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
	MessagingProduct string           `json:"messaging_product,omitempty"`
	Type             string           `json:"type"`
	To               string           `json:"to"`
	Body             string           `json:"body,omitempty"`
	Text             *TextContent     `json:"text,omitempty"`
	Image            *ImageContent    `json:"image,omitempty"`
	Document         *DocumentContent `json:"document,omitempty"`
	Audio            *AudioContent    `json:"audio,omitempty"`
	Sticker          *StickerContent  `json:"sticker,omitempty"`
	MediaURL         string           `json:"media_url,omitempty"`
	Caption          string           `json:"caption,omitempty"`
	SessionID        string           `json:"session_id"`
	MessageID        string           `json:"message_id,omitempty"`
	Context          *MessageContext  `json:"context,omitempty"`
	Mentions         []string         `json:"mentions,omitempty"`
}

func sendWithID(ctx context.Context, cli *whatsmeow.Client, jid types.JID, msg *goE2E.Message, id string) (whatsmeow.SendResponse, error) {
	if strings.TrimSpace(id) != "" {
		return cli.SendMessage(ctx, jid, msg, whatsmeow.SendRequestExtra{ID: types.MessageID(id)})
	}
	return cli.SendMessage(ctx, jid, msg)
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

	cli, ok := m.GetClient(sessionID)
	if !ok || cli == nil {
		return fmt.Errorf("session %s not connected", sessionID)
	}
	entry := log.WithSession(sessionID).WithMessageID(msg.MessageID)

	jid, err := toJID(msg.To)
	if err != nil {
		return fmt.Errorf("invalid recipient %q: %w", msg.To, err)
	}

	var ctxInfo *goE2E.ContextInfo
	if msg.Context != nil || len(msg.Mentions) > 0 {
		var stanzaID string
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
		}
		var mList []string
		if msg.Context != nil {
			mList = append(mList, msg.Context.MentionedJid...)
			mList = append(mList, msg.Context.Mentions...)
		}
		mList = append(mList, msg.Mentions...)
		var mentioned []string
		for _, mnt := range mList {
			mnt = strings.TrimSpace(mnt)
			if mnt == "" {
				continue
			}
			j, err := toUserJID(mnt)
			if err == nil {
				mentioned = append(mentioned, j.String())
			}
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
				qt := ""
				if msg.Context != nil {
					qt = strings.TrimSpace(msg.Context.QuotedText)
				}
				ci.QuotedMessage = &goE2E.Message{Conversation: proto.String(qt)}
			}
			if len(mentioned) > 0 {
				ci.MentionedJID = mentioned
			}
			ctxInfo = ci
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
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send text failed: %v", err)
			return err
		}
		entry.Info("text sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "image":
		var link, caption string
		if msg.Image != nil {
			link = strings.TrimSpace(msg.Image.Link)
			caption = strings.TrimSpace(msg.Image.Caption)
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		if caption == "" {
			caption = strings.TrimSpace(msg.Caption)
		}
		if link == "" {
			return fmt.Errorf("image requires a link or media_url")
		}
		data, mime, err := downloadBytes(ctx, link)
		if err != nil {
			return fmt.Errorf("failed to download media: %w", err)
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
		resp, err := sendWithID(ctx, cli, jid, &goE2E.Message{ImageMessage: imgMsg}, msg.MessageID)
		if err != nil {
			entry.Error("send image failed: %v", err)
			return err
		}
		entry.Info("image sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "document":
		var link, caption, filename string
		if msg.Document != nil {
			link = strings.TrimSpace(msg.Document.Link)
			caption = strings.TrimSpace(msg.Document.Caption)
			filename = strings.TrimSpace(msg.Document.FileName)
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		if caption == "" {
			caption = strings.TrimSpace(msg.Caption)
		}
		if filename == "" {
			parts := strings.Split(strings.TrimRight(link, "/"), "/")
			filename = parts[len(parts)-1]
		}
		if link == "" {
			return fmt.Errorf("document requires a link or media_url")
		}
		data, mime, err := downloadBytes(ctx, link)
		if err != nil {
			return fmt.Errorf("failed to download media: %w", err)
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
		resp, err := sendWithID(ctx, cli, jid, &goE2E.Message{DocumentMessage: docMsg}, msg.MessageID)
		if err != nil {
			entry.Error("send document failed: %v", err)
			return err
		}
		entry.Info("document sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "audio":
		var link string
		var forcePTT *bool
		var seconds *uint32
		if msg.Audio != nil {
			link = strings.TrimSpace(msg.Audio.Link)
			forcePTT = msg.Audio.PTT
			seconds = msg.Audio.Seconds
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		if link == "" {
			return fmt.Errorf("audio requires a link or media_url")
		}

		data, mime, hdr, err := downloadBytesWithHeader(ctx, link)
		if err != nil {
			return fmt.Errorf("failed to download media: %w", err)
		}

		ptt := shouldSendAsPTT(link, mime)
		if forcePTT != nil {
			ptt = *forcePTT
		}
		mime = normalizeAudioMime(link, mime, ptt)
		secs := pickSeconds(seconds, hdr, ptt)

		entry.Info("audio meta (pre-upload): link=%s mime=%s bytes=%d ptt=%v seconds=%d dest=%s",
			link, mime, len(data), ptt, secs, jid.String())

		uploaded, err := cli.Upload(ctx, data, whatsmeow.MediaAudio)
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
			FileLength:    proto.Uint64(uint64(len(data))),
			Mimetype:      proto.String(mime),
			PTT:           &ptt,
		}
		if secs > 0 {
			audioMsg.Seconds = proto.Uint32(secs)
		}
		if ctxInfo != nil {
			audioMsg.ContextInfo = ctxInfo
		}
		if wf := makeWaveform(secs); len(wf) > 0 {
			audioMsg.Waveform = wf
		}

		resp, err := sendWithID(ctx, cli, jid, &goE2E.Message{AudioMessage: audioMsg}, msg.MessageID)
		if err != nil {
			entry.Error("send audio failed: %v", err)
			return err
		}
		entry.Info("audio sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "sticker":
		var link string
		if msg.Sticker != nil {
			link = strings.TrimSpace(msg.Sticker.Link)
		}
		if link == "" {
			link = strings.TrimSpace(msg.MediaURL)
		}
		if link == "" {
			return fmt.Errorf("sticker requires a link or media_url")
		}
		data, mime, err := downloadBytes(ctx, link)
		if err != nil {
			return fmt.Errorf("failed to download media: %w", err)
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
		resp, err := sendWithID(ctx, cli, jid, &goE2E.Message{StickerMessage: stickerMsg}, msg.MessageID)
		if err != nil {
			entry.Error("send sticker failed: %v", err)
			return err
		}
		entry.Info("sticker sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	default:
		return fmt.Errorf("unsupported message type: %s", msg.Type)
	}
}

func makeWaveform(seconds uint32) []byte {
	const n = 32
	wf := make([]byte, n)
	minA, maxA := 6.0, 31.0
	for i := 0; i < n; i++ {
		phase := 2.0 * math.Pi * float64(i) / float64(n)
		s := (math.Sin(phase) + 1.0) / 2.0
		v := minA + s*(maxA-minA)
		if v < 0 {
			v = 0
		}
		if v > 31 {
			v = 31
		}
		wf[i] = byte(int(v + 0.5))
	}
	return wf
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
