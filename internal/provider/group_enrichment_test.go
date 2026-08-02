package provider

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
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
