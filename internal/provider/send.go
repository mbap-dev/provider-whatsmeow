package provider

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"your.org/provider-whatsmeow/internal/log"

	"go.mau.fi/whatsmeow"
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
	MessagingProduct string           `json:"messaging_product,omitempty"`
	Type             string           `json:"type"`
	To               string           `json:"to"`
	Body             string           `json:"body,omitempty"`
	Text             *TextContent     `json:"text,omitempty"`
	Image            *ImageContent    `json:"image,omitempty"`
	Document         *DocumentContent `json:"document,omitempty"`
	Audio            *AudioContent    `json:"audio,omitempty"`
	Sticker          *StickerContent  `json:"sticker,omitempty"`
	Contacts         []ContactEntry   `json:"contacts,omitempty"`
	MediaURL         string           `json:"media_url,omitempty"`
	Caption          string           `json:"caption,omitempty"`
	SessionID        string           `json:"session_id"`
	MessageID        string           `json:"message_id,omitempty"`
	Context          *MessageContext  `json:"context,omitempty"`
	Mentions         []string         `json:"mentions,omitempty"`
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

	// Log input summary early for diagnostics
	entry.Debug("send_input type=%q to=%q media_url=%q has_text=%v has_image=%v has_doc=%v has_audio=%v",
		strings.TrimSpace(msg.Type), strings.TrimSpace(msg.To), strings.TrimSpace(msg.MediaURL),
		msg.Text != nil && strings.TrimSpace(msg.Text.Body) != "",
		msg.Image != nil && strings.TrimSpace(msg.Image.Link) != "",
		msg.Document != nil && strings.TrimSpace(msg.Document.Link) != "",
		msg.Audio != nil && strings.TrimSpace(msg.Audio.Link) != "",
	)

	// Normalize destination: resolve via server when possible, fallback to local PN JID
	rawTo := strings.TrimSpace(msg.To)
	// Try resolving PN -> canonical JID via server (handles BR 9th digit cases robustly)
	var jid types.JID
	var err error
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
	entry.Info("dest_jid=%s", jid.String())

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
		entry.Debug(
			"reply_ctx_select stanza_id=%q ctx_in={id:%q message_id:%q stanzaId:%q stanza_id:%q participant:%q}",
			stanzaID,
			msg.Context.ID,
			msg.Context.MessageID,
			msg.Context.StanzaId,
			msg.Context.Stanza_id,
			msg.Context.Participant,
		)
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
				// Include quoted content when available; otherwise include empty to trigger reply UI (Wuzapi behavior)
				if qt != "" {
					ci.QuotedMessage = &goE2E.Message{Conversation: proto.String(qt)}
				} else if ci.QuotedMessage == nil {
					ci.QuotedMessage = &goE2E.Message{Conversation: proto.String("")}
				}
				// Ensure participant presence for 1:1 replies. For groups, if not provided by client, leave empty (cannot infer safely without context).
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
		wire := &goE2E.Message{DocumentMessage: docMsg}
		entry.Debug("wire_payload=%s", marshalForLog(wire))
		resp, err := sendWithID(ctx, cli, jid, wire, msg.MessageID)
		if err != nil {
			entry.Error("send document failed: %v", err)
			return err
		}
		entry.Info("document sent to %s (wa_msg_id=%s)", jid.String(), resp.ID)
		return m.emitCloudSent(sessionID, jid, resp.ID)

	case "audio":
		// Mirror Evolution API: fetch URL content and convert PTT voice notes to ogg/opus with metadata.
		var src string
		if msg.Audio != nil {
			src = strings.TrimSpace(msg.Audio.Link)
		}
		if src == "" {
			src = strings.TrimSpace(msg.MediaURL)
		}
		if src == "" {
			return fmt.Errorf("audio requires a URL (http/https)")
		}
		if !(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
			return fmt.Errorf("audio link must be an http(s) URL")
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
		if err != nil {
			return err
		}
		httpResp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if httpResp.StatusCode != http.StatusOK {
			defer httpResp.Body.Close()
			return fmt.Errorf("failed to fetch audio: %s", httpResp.Status)
		}
		defer httpResp.Body.Close()

		original, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return fmt.Errorf("reading audio content: %w", err)
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
			secs = pickSeconds(payloadSeconds, httpResp.Header, ptt)
		}

		mime := httpResp.Header.Get("Content-Type")
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

	default:
		return fmt.Errorf("unsupported message type: %s", msg.Type)
	}
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
