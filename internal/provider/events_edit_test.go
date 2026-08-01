package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"go.mau.fi/whatsmeow/util/gcmutil"
	"go.mau.fi/whatsmeow/util/hkdfutil"
	"google.golang.org/protobuf/proto"
)

type editMessageSecretStore struct {
	secret []byte
	sender types.JID
}

func (s *editMessageSecretStore) PutMessageSecrets(context.Context, []store.MessageSecretInsert) error {
	return nil
}

func (s *editMessageSecretStore) PutMessageSecret(context.Context, types.JID, types.JID, types.MessageID, []byte) error {
	return nil
}

func (s *editMessageSecretStore) GetMessageSecret(context.Context, types.JID, types.JID, types.MessageID) ([]byte, types.JID, error) {
	return s.secret, s.sender, nil
}

func TestEmitCloudMessageDeliversUnwrappedGroupEdit(t *testing.T) {
	const (
		sessionID = "553121159080"
		messageID = "3A0EB6D4AD5516704383"
	)

	received := make(chan map[string]any, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode webhook payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		received <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager := NewClientManager(ManagerOptions{
		SessionStore: t.TempDir(),
		WorkerID:     "test-worker",
	})
	_, err := manager.UpsertSessionWebhook(sessionID, SessionWebhook{
		URLAbsolute:          server.URL,
		Enabled:              true,
		SendIncomingMessages: true,
		SendUpdateMessages:   true,
		SendGroupMessages:    true,
	})
	if err != nil {
		t.Fatalf("register webhook: %v", err)
	}

	group := types.NewJID("120363369993948203", types.GroupServer)
	sender := types.NewJID("5562933000233", types.DefaultUserServer)
	evt := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{
				Chat:    group,
				Sender:  sender,
				IsGroup: true,
			},
			ID:        types.MessageID(messageID),
			PushName:  "Moises",
			Timestamp: time.Unix(1785622982, 0),
		},
		Message: &waE2E.Message{Conversation: proto.String("texto novo")},
		IsEdit:  true,
	}

	if err := manager.emitCloudMessage(sessionID, nil, evt); err != nil {
		t.Fatalf("emitCloudMessage: %v", err)
	}

	select {
	case payload := <-received:
		value := extractCloudPayloadValue(payload)
		messages, _ := value["messages"].([]any)
		if len(messages) != 1 {
			t.Fatalf("messages = %#v", value["messages"])
		}
		message, _ := messages[0].(map[string]any)
		if edited, _ := message["edited"].(bool); !edited {
			t.Fatalf("edited = %#v", message["edited"])
		}
		text, _ := message["text"].(map[string]any)
		if got := text["body"]; got != "texto novo\n`Mensagem Editada`" {
			t.Fatalf("text.body = %#v", got)
		}
		context, _ := message["context"].(map[string]any)
		wantContextID := sessionID + "_5562933000233~" + messageID
		if got := context["id"]; got != wantContextID {
			t.Fatalf("context.id = %#v, want %q", got, wantContextID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for edit webhook")
	}
}

func TestEmitCloudMessageDecryptsSecretEncryptedGroupEdit(t *testing.T) {
	const (
		sessionID  = "553121159080"
		originalID = "3AAFB6FE84132A75B2A7"
		editID     = "3AF1DD5878BEC525FD88"
	)

	received := make(chan map[string]any, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode webhook payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		received <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager := NewClientManager(ManagerOptions{
		SessionStore: t.TempDir(),
		WorkerID:     "test-worker",
	})
	_, err := manager.UpsertSessionWebhook(sessionID, SessionWebhook{
		URLAbsolute:          server.URL,
		Enabled:              true,
		SendIncomingMessages: true,
		SendUpdateMessages:   true,
		SendGroupMessages:    true,
	})
	if err != nil {
		t.Fatalf("register webhook: %v", err)
	}

	group := types.NewJID("120363369993948203", types.GroupServer)
	sender := types.NewJID("46849671094330", types.HiddenUserServer)
	baseSecret := []byte("0123456789abcdef0123456789abcdef")
	plaintext, err := proto.Marshal(&waE2E.Message{Conversation: proto.String("Testes")})
	if err != nil {
		t.Fatalf("marshal edited message: %v", err)
	}
	useCase := make([]byte, 0, 128)
	useCase = append(useCase, originalID...)
	useCase = append(useCase, sender.ToNonAD().String()...)
	useCase = append(useCase, sender.ToNonAD().String()...)
	useCase = append(useCase, "Message Edit"...)
	secretKey := hkdfutil.SHA256(baseSecret, nil, useCase, 32)
	iv := []byte("0123456789ab")
	ciphertext, err := gcmutil.Encrypt(secretKey, iv, plaintext, nil)
	if err != nil {
		t.Fatalf("encrypt edited message: %v", err)
	}

	targetKey := &waCommon.MessageKey{
		RemoteJID:   proto.String(group.String()),
		FromMe:      proto.Bool(false),
		ID:          proto.String(originalID),
		Participant: proto.String(sender.String()),
	}
	evt := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{
				Chat:    group,
				Sender:  sender,
				IsGroup: true,
			},
			ID:        types.MessageID(editID),
			PushName:  "Moises",
			Timestamp: time.Unix(1785623984, 0),
			Edit:      types.EditAttributeMessageEdit,
		},
		Message: &waE2E.Message{SecretEncryptedMessage: &waE2E.SecretEncryptedMessage{
			TargetMessageKey: targetKey,
			EncPayload:       ciphertext,
			EncIV:            iv,
			SecretEncType:    waE2E.SecretEncryptedMessage_MESSAGE_EDIT.Enum(),
		}},
	}
	client := whatsmeow.NewClient(&store.Device{
		MsgSecrets: &editMessageSecretStore{secret: baseSecret, sender: sender},
	}, nil)

	if err := manager.emitCloudMessage(sessionID, client, evt); err != nil {
		t.Fatalf("emitCloudMessage: %v", err)
	}

	select {
	case payload := <-received:
		value := extractCloudPayloadValue(payload)
		messages, _ := value["messages"].([]any)
		if len(messages) != 1 {
			t.Fatalf("messages = %#v", value["messages"])
		}
		message, _ := messages[0].(map[string]any)
		if edited, _ := message["edited"].(bool); !edited {
			t.Fatalf("edited = %#v", message["edited"])
		}
		text, _ := message["text"].(map[string]any)
		if got := text["body"]; got != "Testes\n`Mensagem Editada`" {
			t.Fatalf("text.body = %#v", got)
		}
		context, _ := message["context"].(map[string]any)
		wantContextID := sessionID + "_46849671094330~" + originalID
		if got := context["id"]; got != wantContextID {
			t.Fatalf("context.id = %#v, want %q", got, wantContextID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for secret-encrypted edit webhook")
	}
}
