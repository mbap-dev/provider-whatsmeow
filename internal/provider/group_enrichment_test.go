package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func resetGroupEnrichmentCaches() {
	cacheMu.Lock()
	avatarCache = map[string]stringEntry{}
	groupNameCache = map[string]stringEntry{}
	cacheMu.Unlock()
}

func TestGetAvatarURLUsesOnePreviewLookupAndNegativeCache(t *testing.T) {
	resetGroupEnrichmentCaches()
	oldConnected := isWhatsAppConnected
	oldFetch := fetchProfilePictureInfo
	t.Cleanup(func() {
		isWhatsAppConnected = oldConnected
		fetchProfilePictureInfo = oldFetch
		resetGroupEnrichmentCaches()
	})

	isWhatsAppConnected = func(*whatsmeow.Client) bool { return true }
	calls := 0
	fetchProfilePictureInfo = func(_ context.Context, _ *whatsmeow.Client, _ types.JID, params *whatsmeow.GetProfilePictureParams) (*types.ProfilePictureInfo, error) {
		calls++
		if params == nil || !params.Preview {
			t.Fatalf("params = %#v, want preview lookup", params)
		}
		return nil, errors.New("profile picture hidden")
	}

	client := whatsmeow.NewClient(&store.Device{}, nil)
	jid := types.NewJID("5562999999999", types.DefaultUserServer)
	if _, ok := getAvatarURL("553121159080", client, jid); ok {
		t.Fatal("expected missing avatar")
	}
	if _, ok := getAvatarURL("553121159080", client, jid); ok {
		t.Fatal("expected cached missing avatar")
	}
	if calls != 1 {
		t.Fatalf("lookup calls = %d, want 1", calls)
	}
}

func TestGetAvatarURLSkipsDisconnectedClient(t *testing.T) {
	resetGroupEnrichmentCaches()
	oldConnected := isWhatsAppConnected
	oldFetch := fetchProfilePictureInfo
	t.Cleanup(func() {
		isWhatsAppConnected = oldConnected
		fetchProfilePictureInfo = oldFetch
		resetGroupEnrichmentCaches()
	})

	isWhatsAppConnected = func(*whatsmeow.Client) bool { return false }
	fetchProfilePictureInfo = func(context.Context, *whatsmeow.Client, types.JID, *whatsmeow.GetProfilePictureParams) (*types.ProfilePictureInfo, error) {
		t.Fatal("disconnected client started a lookup")
		return nil, nil
	}

	client := whatsmeow.NewClient(&store.Device{}, nil)
	if _, ok := getAvatarURL("553121159080", client, types.NewJID("5562999999999", types.DefaultUserServer)); ok {
		t.Fatal("expected missing avatar")
	}
}

func TestSafeGetGroupNameUsesOneLookupAndCache(t *testing.T) {
	resetGroupEnrichmentCaches()
	oldConnected := isWhatsAppConnected
	oldFetch := fetchGroupInfo
	t.Cleanup(func() {
		isWhatsAppConnected = oldConnected
		fetchGroupInfo = oldFetch
		resetGroupEnrichmentCaches()
	})

	isWhatsAppConnected = func(*whatsmeow.Client) bool { return true }
	calls := 0
	fetchGroupInfo = func(context.Context, *whatsmeow.Client, types.JID) (*types.GroupInfo, error) {
		calls++
		return &types.GroupInfo{GroupName: types.GroupName{Name: "Grupo Comercial"}}, nil
	}

	client := whatsmeow.NewClient(&store.Device{}, nil)
	group := types.NewJID("120363000000000000", types.GroupServer)
	for range 2 {
		name, ok := safeGetGroupName(client, group)
		if !ok || name != "Grupo Comercial" {
			t.Fatalf("safeGetGroupName() = %q, %t", name, ok)
		}
	}
	if calls != 1 {
		t.Fatalf("group info calls = %d, want 1", calls)
	}
}

func TestGetAvatarURLDeduplicatesConcurrentLookups(t *testing.T) {
	resetGroupEnrichmentCaches()
	oldConnected := isWhatsAppConnected
	oldFetch := fetchProfilePictureInfo
	t.Cleanup(func() {
		isWhatsAppConnected = oldConnected
		fetchProfilePictureInfo = oldFetch
		resetGroupEnrichmentCaches()
	})

	isWhatsAppConnected = func(*whatsmeow.Client) bool { return true }
	release := make(chan struct{})
	var calls atomic.Int32
	fetchProfilePictureInfo = func(context.Context, *whatsmeow.Client, types.JID, *whatsmeow.GetProfilePictureParams) (*types.ProfilePictureInfo, error) {
		calls.Add(1)
		<-release
		return &types.ProfilePictureInfo{URL: "https://example.test/avatar.jpg"}, nil
	}

	client := whatsmeow.NewClient(&store.Device{}, nil)
	jid := types.NewJID("5562999999999", types.DefaultUserServer)
	const workers = 8
	start := make(chan struct{})
	var ready sync.WaitGroup
	var done sync.WaitGroup
	ready.Add(workers)
	done.Add(workers)
	for range workers {
		go func() {
			defer done.Done()
			ready.Done()
			<-start
			if url, ok := getAvatarURL("553121159080", client, jid); !ok || url == "" {
				t.Errorf("getAvatarURL() = %q, %t", url, ok)
			}
		}()
	}
	ready.Wait()
	close(start)
	deadline := time.After(time.Second)
	for calls.Load() == 0 {
		select {
		case <-deadline:
			t.Fatal("lookup did not start")
		default:
		}
	}
	time.Sleep(20 * time.Millisecond)
	close(release)
	done.Wait()
	if calls.Load() != 1 {
		t.Fatalf("lookup calls = %d, want 1", calls.Load())
	}
}

func TestEmitCloudMessageBuildsIncrementalGroupContact(t *testing.T) {
	resetGroupEnrichmentCaches()
	oldConnected := isWhatsAppConnected
	oldGroup := fetchGroupInfo
	oldPicture := fetchProfilePictureInfo
	t.Cleanup(func() {
		isWhatsAppConnected = oldConnected
		fetchGroupInfo = oldGroup
		fetchProfilePictureInfo = oldPicture
		resetGroupEnrichmentCaches()
	})

	isWhatsAppConnected = func(*whatsmeow.Client) bool { return true }
	groupInfoCalls := 0
	fetchGroupInfo = func(_ context.Context, _ *whatsmeow.Client, _ types.JID) (*types.GroupInfo, error) {
		groupInfoCalls++
		return &types.GroupInfo{
			GroupName: types.GroupName{Name: "Grupo Comercial"},
			Participants: []types.GroupParticipant{{
				JID: types.NewJID("5511988888888", types.DefaultUserServer),
			}},
		}, nil
	}
	pictureJIDs := make([]string, 0, 3)
	fetchProfilePictureInfo = func(_ context.Context, _ *whatsmeow.Client, jid types.JID, _ *whatsmeow.GetProfilePictureParams) (*types.ProfilePictureInfo, error) {
		pictureJIDs = append(pictureJIDs, jid.String())
		return &types.ProfilePictureInfo{URL: "https://example.test/" + jid.User + ".jpg"}, nil
	}

	received := make(chan map[string]any, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode webhook: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		received <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	manager := NewClientManager(ManagerOptions{SessionStore: t.TempDir(), WorkerID: "test-worker"})
	if _, err := manager.UpsertSessionWebhook("553121159080", SessionWebhook{
		URLAbsolute:          server.URL,
		Enabled:              true,
		SendIncomingMessages: true,
		SendGroupMessages:    true,
	}); err != nil {
		t.Fatalf("register webhook: %v", err)
	}

	event := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{
				Chat:    types.NewJID("120363000000000000", types.GroupServer),
				Sender:  types.NewJID("5562999999999", types.DefaultUserServer),
				IsGroup: true,
			},
			ID:        "3A0EB6D4AD5516704383",
			PushName:  "Maria",
			Timestamp: time.Unix(1785622982, 0),
		},
		Message: &waE2E.Message{Conversation: proto.String("Ola")},
	}
	client := whatsmeow.NewClient(&store.Device{}, nil)
	if err := manager.emitCloudMessage("553121159080", client, event); err != nil {
		t.Fatalf("emitCloudMessage: %v", err)
	}

	var payload map[string]any
	select {
	case payload = <-received:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for webhook")
	}
	value := extractCloudPayloadValue(payload)
	contacts := value["contacts"].([]any)
	contact := contacts[0].(map[string]any)
	if _, exists := contact["group_participants"]; exists {
		t.Fatalf("group_participants must be omitted: %#v", contact["group_participants"])
	}
	profile := contact["profile"].(map[string]any)
	if profile["name"] != "Maria" {
		t.Fatalf("profile.name = %#v", profile["name"])
	}
	if profile["picture"] != "https://example.test/5562999999999.jpg" {
		t.Fatalf("profile.picture = %#v", profile["picture"])
	}
	if contact["group_subject"] != "Grupo Comercial" || contact["group_picture"] != "https://example.test/120363000000000000.jpg" {
		t.Fatalf("group metadata = %#v", contact)
	}
	if groupInfoCalls != 1 {
		t.Fatalf("GetGroupInfo calls = %d, want 1", groupInfoCalls)
	}
	sort.Strings(pictureJIDs)
	wantPictures := []string{"120363000000000000@g.us", "5562999999999@s.whatsapp.net"}
	if !reflect.DeepEqual(pictureJIDs, wantPictures) {
		t.Fatalf("picture JIDs = %#v, want %#v", pictureJIDs, wantPictures)
	}
}
