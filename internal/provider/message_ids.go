package provider

import "strings"

const groupParticipantIDSeparator = "_"

// WithInstanceMessageID returns the public-facing message ID used by
// UnoAPI-compatible webhooks and API responses, prefixed with the
// phone/session identifier when needed.
func WithInstanceMessageID(phone, id string) string {
	return withInstanceMessageID(phone, id)
}

// WithGroupParticipantMessageID returns the group-compatible public message ID
// that embeds the quoted message participant before the WhatsApp message ID.
func WithGroupParticipantMessageID(participant, id string) string {
	return withGroupParticipantMessageID(participant, id)
}

func instanceMessagePrefix(phone string) string {
	normalized := strings.TrimPrefix(normalizePhone(phone), "+")
	if normalized == "" {
		return ""
	}
	return normalized + "_"
}

func withInstanceMessageID(phone, id string) string {
	rawID := strings.TrimSpace(id)
	if rawID == "" {
		return ""
	}
	prefix := instanceMessagePrefix(phone)
	if prefix == "" || strings.HasPrefix(rawID, prefix) {
		return rawID
	}
	return prefix + rawID
}

func stripInstanceMessageID(phone, id string) string {
	rawID := strings.TrimSpace(id)
	if rawID == "" {
		return ""
	}
	prefix := instanceMessagePrefix(phone)
	if prefix == "" || !strings.HasPrefix(rawID, prefix) {
		return rawID
	}
	return strings.TrimPrefix(rawID, prefix)
}

func applyInstanceMessageIDs(phone string, payload map[string]any) {
	entries, ok := payload["entry"].([]any)
	if !ok {
		return
	}
	for _, entryAny := range entries {
		entry, ok := entryAny.(map[string]any)
		if !ok {
			continue
		}
		changes, ok := entry["changes"].([]any)
		if !ok {
			continue
		}
		for _, changeAny := range changes {
			change, ok := changeAny.(map[string]any)
			if !ok {
				continue
			}
			value, ok := change["value"].(map[string]any)
			if !ok {
				continue
			}
			applyInstanceIDsToMessages(phone, value["messages"])
			applyInstanceIDsToMessages(phone, value["message_echoes"])
			applyInstanceIDsToStatuses(phone, value["statuses"])
		}
	}
}

func applyInstanceIDsToMessages(phone string, raw any) {
	messages, ok := raw.([]any)
	if !ok {
		return
	}
	for _, messageAny := range messages {
		message, ok := messageAny.(map[string]any)
		if !ok {
			continue
		}
		if id, ok := message["id"].(string); ok {
			message["id"] = withInstanceMessageID(phone, id)
		}
		if context, ok := message["context"].(map[string]any); ok {
			if id, ok := context["id"].(string); ok {
				context["id"] = withInstanceMessageID(phone, id)
			}
			if messageID, ok := context["message_id"].(string); ok {
				context["message_id"] = withInstanceMessageID(phone, messageID)
			}
		}
		if reaction, ok := message["reaction"].(map[string]any); ok {
			if messageID, ok := reaction["message_id"].(string); ok {
				reaction["message_id"] = withInstanceMessageID(phone, messageID)
			}
		}
	}
}

func applyInstanceIDsToStatuses(phone string, raw any) {
	statuses, ok := raw.([]any)
	if !ok {
		return
	}
	for _, statusAny := range statuses {
		status, ok := statusAny.(map[string]any)
		if !ok {
			continue
		}
		if id, ok := status["id"].(string); ok {
			status["id"] = withInstanceMessageID(phone, id)
		}
		if messageID, ok := status["message_id"].(string); ok {
			status["message_id"] = withInstanceMessageID(phone, messageID)
		}
	}
}

func normalizeOutgoingMessageIDs(phone string, msg *OutgoingMessage) {
	if msg == nil {
		return
	}
	msg.MessageID = stripInstanceMessageID(phone, msg.MessageID)
	if msg.Context == nil {
		return
	}
	msg.Context.ID = stripInstanceMessageID(phone, msg.Context.ID)
	msg.Context.MessageID = stripInstanceMessageID(phone, msg.Context.MessageID)
	msg.Context.StanzaId = stripInstanceMessageID(phone, msg.Context.StanzaId)
	msg.Context.Stanza_id = stripInstanceMessageID(phone, msg.Context.Stanza_id)
}

func withGroupParticipantMessageID(participant, id string) string {
	msgID := strings.TrimSpace(id)
	if msgID == "" {
		return ""
	}
	if _, _, ok := splitGroupParticipantMessageID(msgID); ok {
		return msgID
	}
	user := participantMessageIDUser(participant)
	if user == "" {
		return msgID
	}
	return user + groupParticipantIDSeparator + msgID
}

func splitGroupParticipantMessageID(id string) (participant, messageID string, ok bool) {
	rawID := strings.TrimSpace(id)
	idx := strings.Index(rawID, groupParticipantIDSeparator)
	if idx <= 0 || idx >= len(rawID)-1 {
		return "", rawID, false
	}
	user := participantMessageIDUser(rawID[:idx])
	msgID := strings.TrimSpace(rawID[idx+1:])
	if user == "" || msgID == "" {
		return "", rawID, false
	}
	return user, msgID, true
}

func participantMessageIDUser(participant string) string {
	p := strings.TrimSpace(participant)
	if p == "" {
		return ""
	}
	if at := strings.Index(p, "@"); at >= 0 {
		p = p[:at]
	}
	return strings.TrimSpace(p)
}
