# Efficient Group Message Enrichment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make live group-message enrichment constant with respect to group size while preserving the Chatwoot group, author name, author picture, receipt, edit, deletion, and media contracts.

**Architecture:** Remove eager `group_participants` construction from live messages and make administrative group events explicit no-ops. Keep group and current-author enrichment, but route WhatsApp metadata lookups through cache-first, singleflight-protected helpers that perform one request per cache miss and never block webhook delivery on optional metadata.

**Tech Stack:** Go 1.26 module, whatsmeow, `golang.org/x/sync/singleflight`, existing in-memory caches, Go `testing`, `httptest`, Make.

## Global Constraints

- `HistorySync` behavior is out of scope.
- Do not enumerate a group's full participant list from a live message or administrative event.
- Only the current message author and the group may have profile pictures requested.
- Perform at most one author-picture request and one group-picture request per cache miss.
- Cache successful avatar lookups for one hour and failed/hidden/absent lookups for six hours.
- Cache successful group names for one hour.
- Missing optional metadata must not prevent webhook delivery.
- Keep receipt, deletion, edit, and media payload contracts compatible.
- Run both `make check` and `make build` after code changes, per `AGENTS.md`.

---

## File Structure

- Modify `internal/provider/events.go`: add bounded lookup seams and singleflight groups, simplify avatar/group-name lookups, stop attaching complete participant lists, and make administrative group webhooks no-ops.
- Modify `internal/provider/caches.go`: set the approved avatar and group-name TTLs.
- Create `internal/provider/group_enrichment_test.go`: focused tests for bounded lookups, cache behavior, message payload shape, and ignored administrative events.
- Modify `README.md`: document incremental group-contact behavior and the absence of eager participant synchronization.
- Modify `go.mod` and `go.sum`: make `golang.org/x/sync` a direct dependency after importing `singleflight` and running `go mod tidy`.

### Task 1: Bound WhatsApp metadata lookups

**Files:**
- Modify: `internal/provider/events.go:2520-2678,2726-2734`
- Modify: `internal/provider/caches.go:34-41`
- Create: `internal/provider/group_enrichment_test.go`
- Modify: `go.mod`
- Modify: `go.sum`

**Interfaces:**
- Produces: `getAvatarURL(sessionID string, cli *whatsmeow.Client, jid types.JID) (string, bool)` with one lookup per cache miss.
- Produces: `safeGetGroupName(cli *whatsmeow.Client, gid types.JID) (string, bool)` with one lookup per cache miss.
- Produces test seams `isWhatsAppConnected`, `fetchProfilePictureInfo`, and `fetchGroupInfo` with production defaults that call the real client.
- Uses package-level `singleflight.Group` values keyed by session/JID for avatars and group JID for names.

- [ ] **Step 1: Write failing avatar lookup tests**

Create `internal/provider/group_enrichment_test.go` with cache reset and lookup seam restoration performed directly in each test:

```go
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
```

- [ ] **Step 2: Run the avatar tests and verify RED**

Run:

```bash
rtk go test ./internal/provider -run 'TestGetAvatarURL(UsesOnePreviewLookupAndNegativeCache|SkipsDisconnectedClient)$' -count=1
```

Expected: build failure because `isWhatsAppConnected` and `fetchProfilePictureInfo` do not exist.

- [ ] **Step 3: Add failing group-name and concurrent avatar tests**

Append tests that verify one `GetGroupInfo` call is cached and concurrent avatar misses share one request:

```go
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
        return &types.GroupInfo{Name: "Grupo Comercial"}, nil
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
```

- [ ] **Step 4: Implement the bounded lookup seams and singleflight behavior**

In `internal/provider/events.go`, import `golang.org/x/sync/singleflight` and define:

```go
var (
    avatarLookupGroup    singleflight.Group
    groupNameLookupGroup singleflight.Group

    isWhatsAppConnected = func(cli *whatsmeow.Client) bool {
        return cli != nil && cli.IsConnected()
    }
    fetchProfilePictureInfo = func(ctx context.Context, cli *whatsmeow.Client, jid types.JID, params *whatsmeow.GetProfilePictureParams) (*types.ProfilePictureInfo, error) {
        return cli.GetProfilePictureInfo(ctx, jid, params)
    }
    fetchGroupInfo = func(ctx context.Context, cli *whatsmeow.Client, jid types.JID) (*types.GroupInfo, error) {
        return cli.GetGroupInfo(ctx, jid)
    }
)
```

Replace `safeGetGroupName` with one cache-first, connected, singleflight-protected lookup. Keep the existing five-second timeout and cache only a non-empty name:

```go
func safeGetGroupName(cli *whatsmeow.Client, gid types.JID) (string, bool) {
    if name, ok := cacheGetGroupName(gid.String()); ok {
        return name, true
    }
    if !isWhatsAppConnected(cli) {
        return "", false
    }

    value, err, _ := groupNameLookupGroup.Do(gid.String(), func() (any, error) {
        if name, ok := cacheGetGroupName(gid.String()); ok {
            return name, nil
        }
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        info, err := fetchGroupInfo(ctx, cli, gid)
        if err != nil || info == nil || strings.TrimSpace(info.Name) == "" {
            return "", err
        }
        name := strings.TrimSpace(info.Name)
        cachePutGroupName(gid.String(), name)
        return name, nil
    })
    name, ok := value.(string)
    return name, err == nil && ok && name != ""
}
```

Replace `getAvatarURL` with a cache-first, connected, singleflight-protected lookup using only `Preview: true`. Use a session-scoped cache key so one linked account never reuses another account's profile-picture visibility result. On every empty result or error, call `cachePutAvatarMissing` and return false; remove all sleeps and retries:

```go
func getAvatarURL(sessionID string, cli *whatsmeow.Client, jid types.JID) (string, bool) {
    jid = normalizedAvatarLookupJID(sessionID, cli, jid)
    if cli == nil || jid == (types.JID{}) {
        return "", false
    }

    key := sessionID + "|" + jid.String()
    if cached, ok := cacheGetAvatar(key); ok {
        return cached, cached != avatarMissingSentinel
    }
    if !isWhatsAppConnected(cli) {
        return "", false
    }

    value, _, _ := avatarLookupGroup.Do(key, func() (any, error) {
        if cached, ok := cacheGetAvatar(key); ok {
            return cached, nil
        }
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        info, err := fetchProfilePictureInfo(ctx, cli, jid, &whatsmeow.GetProfilePictureParams{Preview: true})
        if err != nil || info == nil || strings.TrimSpace(info.URL) == "" {
            cachePutAvatarMissing(key)
            return avatarMissingSentinel, nil
        }
        url := strings.TrimSpace(info.URL)
        cachePutAvatar(key, url)
        return url, nil
    })
    url, ok := value.(string)
    return url, ok && url != "" && url != avatarMissingSentinel
}
```

Route the existing `groupParticipantsPayload` metadata call through the same
test seam while preserving its behavior until Task 2:

```go
info, err := fetchGroupInfo(ctx, cli, groupJID)
```

This makes the eager participant lookup observable in Task 2's failing test;
Task 2 then removes that call site from live messages.

- [ ] **Step 5: Apply the approved cache TTLs**

In `internal/provider/caches.go`, set:

```go
avatarTTL         = time.Hour
avatarNegativeTTL = 6 * time.Hour
groupNameTTL      = time.Hour
```

- [ ] **Step 6: Run Task 1 tests and dependency cleanup**

Run:

```bash
rtk gofmt -w internal/provider/events.go internal/provider/caches.go internal/provider/group_enrichment_test.go
rtk go mod tidy
rtk go test ./internal/provider -run 'Test(GetAvatarURL|SafeGetGroupName)' -count=1
```

Expected: PASS; `golang.org/x/sync` is a direct dependency.

- [ ] **Step 7: Commit Task 1**

```bash
rtk git add internal/provider/events.go internal/provider/caches.go internal/provider/group_enrichment_test.go go.mod go.sum
rtk git commit -m "fix: bound WhatsApp metadata lookups"
```

### Task 2: Make group messages incremental in Chatwoot

**Files:**
- Modify: `internal/provider/events.go:991-1066`
- Modify: `internal/provider/group_enrichment_test.go`
- Modify: `README.md`

**Interfaces:**
- Consumes: bounded `safeGetGroupName` and `getAvatarURL` from Task 1.
- Produces: inbound group payloads containing the current author in `contacts[0]`, `group_id`, `group_subject`, and `group_picture`, without `group_participants`.

- [ ] **Step 1: Write the failing group webhook contract test**

Add `TestEmitCloudMessageBuildsIncrementalGroupContact` to `group_enrichment_test.go`. Use an `httptest.Server`, register a session webhook, and install the Task 1 lookup seams. Return one unrelated member from `fetchGroupInfo` so the old path exposes the regression:

Extend the test import block with the concrete packages used by the webhook
assertions:

```go
import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "reflect"
    "sort"

    "go.mau.fi/whatsmeow/proto/waE2E"
    "go.mau.fi/whatsmeow/types/events"
    "google.golang.org/protobuf/proto"
)
```

Capture the webhook with this server and manager setup before emitting the
message:

```go
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
    URLAbsolute: server.URL, Enabled: true, SendIncomingMessages: true, SendGroupMessages: true,
}); err != nil {
    t.Fatalf("register webhook: %v", err)
}
```

```go
fetchGroupInfo = func(_ context.Context, _ *whatsmeow.Client, _ types.JID) (*types.GroupInfo, error) {
    groupInfoCalls++
    return &types.GroupInfo{
        Name: "Grupo Comercial",
        Participants: []types.GroupParticipant{{
            JID: types.NewJID("5511988888888", types.DefaultUserServer),
        }},
    }, nil
}
fetchProfilePictureInfo = func(_ context.Context, _ *whatsmeow.Client, jid types.JID, _ *whatsmeow.GetProfilePictureParams) (*types.ProfilePictureInfo, error) {
    pictureJIDs = append(pictureJIDs, jid.String())
    return &types.ProfilePictureInfo{URL: "https://example.test/" + jid.User + ".jpg"}, nil
}
```

Emit a text message from `5562999999999@s.whatsapp.net` in group `120363000000000000@g.us`, then assert:

```go
event := &events.Message{
    Info: types.MessageInfo{
        MessageSource: types.MessageSource{
            Chat: types.NewJID("120363000000000000", types.GroupServer),
            Sender: types.NewJID("5562999999999", types.DefaultUserServer),
            IsGroup: true,
        },
        ID: "3A0EB6D4AD5516704383",
        PushName: "Maria",
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
```

```go
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
```

- [ ] **Step 2: Run the group webhook test and verify RED**

Run:

```bash
rtk go test ./internal/provider -run TestEmitCloudMessageBuildsIncrementalGroupContact -count=1
```

Expected: FAIL because `group_participants` is present, `GetGroupInfo` is called more than once, or the unrelated member has a picture requested.

- [ ] **Step 3: Remove eager participant enrichment from live messages**

Delete this block from `emitCloudMessage`:

```go
if isGroupJID(chatJID) {
    if participants := m.groupParticipantsPayload(sessionID, client, chatJID); len(participants) > 0 {
        contactObj["group_participants"] = participants
    }
}
```

Do not remove the current author-name logic, `buildCloudProfile`, `group_subject`, or `group_picture` construction.

- [ ] **Step 4: Document the incremental Chatwoot contract**

Add a short README subsection near webhook behavior stating:

```markdown
### Group message contacts

Inbound group-message webhooks include the current author in `contacts[0]`
along with `group_id`, `group_subject`, and `group_picture`. The provider does
not eagerly send `group_participants`; Chatwoot associates authors with the
group incrementally as they send messages. Author and group pictures are
cache-backed optional enrichment and never block message delivery.
```

- [ ] **Step 5: Run group-message and edit regression tests**

Run:

```bash
rtk gofmt -w internal/provider/events.go internal/provider/group_enrichment_test.go
rtk go test ./internal/provider -run 'TestEmitCloudMessage(BuildsIncrementalGroupContact|DeliversUnwrappedGroupEdit|DecryptsSecretEncryptedGroupEdit)' -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit Task 2**

```bash
rtk git add internal/provider/events.go internal/provider/group_enrichment_test.go README.md
rtk git commit -m "fix: enrich only active group authors"
```

### Task 3: Ignore administrative group events without lookups

**Files:**
- Modify: `internal/provider/events.go:1083-1165`
- Modify: `internal/provider/group_enrichment_test.go`

**Interfaces:**
- Produces: `emitCloudGroupInfo` and `emitCloudJoinedGroup` as explicit no-op policy boundaries returning nil.
- Preserves event-handler registration while guaranteeing `add`, `remove`, `promote`, `demote`, and joined-group events perform no enrichment or webhook delivery.

- [ ] **Step 1: Write the failing administrative-event test**

Add a test that installs lookup seams which count calls, invokes `emitCloudGroupInfo` with a remove event, and asserts zero calls:

```go
func TestEmitCloudGroupInfoIgnoresAdministrativeEvents(t *testing.T) {
    resetGroupEnrichmentCaches()
    oldConnected := isWhatsAppConnected
    oldGroup := fetchGroupInfo
    oldPicture := fetchProfilePictureInfo
    t.Cleanup(func() {
        isWhatsAppConnected = oldConnected
        fetchGroupInfo = oldGroup
        fetchProfilePictureInfo = oldPicture
    })

    isWhatsAppConnected = func(*whatsmeow.Client) bool { return true }
    lookups := 0
    fetchGroupInfo = func(context.Context, *whatsmeow.Client, types.JID) (*types.GroupInfo, error) {
        lookups++
        return &types.GroupInfo{Name: "Grupo"}, nil
    }
    fetchProfilePictureInfo = func(context.Context, *whatsmeow.Client, types.JID, *whatsmeow.GetProfilePictureParams) (*types.ProfilePictureInfo, error) {
        lookups++
        return &types.ProfilePictureInfo{URL: "https://example.test/avatar.jpg"}, nil
    }

    manager := NewClientManager(ManagerOptions{SessionStore: t.TempDir(), WorkerID: "test-worker"})
    client := whatsmeow.NewClient(&store.Device{}, nil)
    err := manager.emitCloudGroupInfo("553121159080", client, &events.GroupInfo{
        JID:   types.NewJID("120363000000000000", types.GroupServer),
        Leave: []types.JID{types.NewJID("5562999999999", types.DefaultUserServer)},
    })
    if err != nil {
        t.Fatalf("emitCloudGroupInfo: %v", err)
    }
    if lookups != 0 {
        t.Fatalf("administrative event lookups = %d, want 0", lookups)
    }
}
```

- [ ] **Step 2: Run the administrative-event test and verify RED**

Run:

```bash
rtk go test ./internal/provider -run TestEmitCloudGroupInfoIgnoresAdministrativeEvents -count=1
```

Expected: FAIL with a non-zero lookup count.

- [ ] **Step 3: Replace administrative webhook builders with explicit no-ops**

Replace both function bodies with:

```go
func (m *ClientManager) emitCloudJoinedGroup(string, *whatsmeow.Client, *events.JoinedGroup) error {
    return nil
}

func (m *ClientManager) emitCloudGroupInfo(string, *whatsmeow.Client, *events.GroupInfo) error {
    return nil
}
```

Keep the event-handler cases so the ignore policy remains explicit and can be changed later without touching whatsmeow registration. Do not change receipt, message, or history dispatch.

- [ ] **Step 4: Run administrative and payload tests**

Run:

```bash
rtk gofmt -w internal/provider/events.go internal/provider/group_enrichment_test.go
rtk go test ./internal/provider -run 'Test(EmitCloudGroupInfoIgnoresAdministrativeEvents|EmitCloudMessageBuildsIncrementalGroupContact|GroupParticipantChange)' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit Task 3**

```bash
rtk git add internal/provider/events.go internal/provider/group_enrichment_test.go
rtk git commit -m "fix: ignore administrative group webhooks"
```

### Task 4: Full verification

**Files:**
- Verify only; no planned source changes.

**Interfaces:**
- Consumes all behavior from Tasks 1-3.
- Produces evidence that tests, vet/tidy checks, and the adapter build pass.

- [ ] **Step 1: Check formatting and diff hygiene**

Run:

```bash
rtk gofmt -w internal/provider/events.go internal/provider/caches.go internal/provider/group_enrichment_test.go
rtk git diff --check
rtk git status --short
```

Expected: no formatting errors or unrelated files.

- [ ] **Step 2: Run the focused provider suite**

Run:

```bash
rtk go test ./internal/provider -count=1
```

Expected: PASS.

- [ ] **Step 3: Run repository-required checks**

Run:

```bash
rtk make check
rtk make build
```

Expected: both commands exit 0. If local Opus/CGO dependencies are unavailable, use the documented fallback and report it explicitly:

```bash
rtk go test -tags nopkgconfig ./...
rtk go build -tags nopkgconfig -o bin/whatsmeow-adapter ./cmd/whatsmeow-adapter
```

- [ ] **Step 4: Inspect the final change set**

Run:

```bash
rtk git status --short
rtk git log -4 --oneline
rtk git diff HEAD~3..HEAD --stat
```

Expected: only the approved group-enrichment implementation, tests, dependency metadata, and README documentation are present.
