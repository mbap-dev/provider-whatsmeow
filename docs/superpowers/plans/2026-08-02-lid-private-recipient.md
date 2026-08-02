# LID Private Recipient Resolution Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Allow private messages to registered contacts whose canonical WhatsApp address is an `@lid` JID.

**Architecture:** Centralize the accepted canonical private-JID rule in `resolver.go` and use it in both outbound resolution paths. Route `IsOnWhatsApp` through a narrow test seam so regression tests can reproduce the exact LID response seen in local Docker logs without making network calls.

**Tech Stack:** Go 1.26, whatsmeow, Go `testing`, Docker Compose.

## Global Constraints

- Require `IsIn=true` before accepting any recipient.
- Accept only non-empty `@s.whatsapp.net` and `@lid` canonical JIDs.
- Prefer the canonical `JID` returned by whatsmeow.
- Preserve existing Brazilian number candidates, error handling, and 24-hour PN cache.
- Do not change groups, webhooks, media, receipts, or HistorySync.
- Run `make check` and `make build` before rebuilding Docker.

---

### Task 1: Accept canonical private LID recipients

**Files:**
- Modify: `internal/provider/resolver.go`
- Create: `internal/provider/resolver_lid_test.go`

**Interfaces:**
- Produces: `canonicalPrivateRecipient(types.IsOnWhatsAppResponse) (types.JID, bool)`.
- Produces test seam: `isOnWhatsApp(context.Context, *whatsmeow.Client, []string) ([]types.IsOnWhatsAppResponse, error)`.
- Updates: `resolvePNtoJID` and `checkUserExistsOnWhatsApp`.

- [ ] **Step 1: Write the failing LID regression tests**

Create `internal/provider/resolver_lid_test.go` with tests that replace and restore the lookup seam. Return this exact response on the first query:

```go
types.IsOnWhatsAppResponse{
    Query:       "+5562933000233",
    IsIn:        true,
    JID:         types.NewJID("46849671094330", types.HiddenUserServer),
    PhoneNumber: types.NewJID("5562933000233", types.DefaultUserServer),
}
```

`TestCheckUserExistsOnWhatsAppAcceptsCanonicalLID` must assert that:

```go
jid, found, err := manager.checkUserExistsOnWhatsApp(context.Background(), nil, "5562933000233")
// err == nil
// found == true
// jid == 46849671094330@lid
// lookup calls == 1, so the alternate BR candidate is not queried
// manager.pnCache["+5562933000233"] contains the canonical LID
```

`TestResolvePNtoJIDAcceptsCanonicalLID` must assert the same canonical JID and one lookup call through the legacy resolver.

Add a table test for `canonicalPrivateRecipient` that accepts PN and LID only when `IsIn=true`, and rejects a group JID, an empty JID, and `IsIn=false`.

- [ ] **Step 2: Run the focused tests and verify RED**

Run:

```bash
rtk mise exec -- go test ./internal/provider -run 'Test(CheckUserExistsOnWhatsAppAcceptsCanonicalLID|ResolvePNtoJIDAcceptsCanonicalLID|CanonicalPrivateRecipient)' -count=1
```

Expected: build failure because the seam and `canonicalPrivateRecipient` do not exist.

- [ ] **Step 3: Add the lookup seam and canonical private-JID rule**

In `resolver.go`, define:

```go
var isOnWhatsApp = func(ctx context.Context, cli *whatsmeow.Client, phones []string) ([]types.IsOnWhatsAppResponse, error) {
    return cli.IsOnWhatsApp(ctx, phones)
}

func canonicalPrivateRecipient(info types.IsOnWhatsAppResponse) (types.JID, bool) {
    if !info.IsIn || info.JID.User == "" {
        return types.EmptyJID, false
    }
    switch info.JID.Server {
    case types.DefaultUserServer, types.HiddenUserServer:
        return info.JID, true
    default:
        return types.EmptyJID, false
    }
}
```

Replace both direct `cli.IsOnWhatsApp` calls in `resolvePNtoJID` and `checkUserExistsOnWhatsApp` with `isOnWhatsApp`. Replace the strict `DefaultUserServer` conditions with `canonicalPrivateRecipient`, then cache and return the selected canonical JID.

- [ ] **Step 4: Run focused and package tests**

Run:

```bash
rtk mise exec -- gofmt -w internal/provider/resolver.go internal/provider/resolver_lid_test.go
rtk mise exec -- go test ./internal/provider -run 'Test(CheckUserExistsOnWhatsAppAcceptsCanonicalLID|ResolvePNtoJIDAcceptsCanonicalLID|CanonicalPrivateRecipient|CandidatesBR)' -count=1
rtk mise exec -- go test ./internal/provider -count=1
```

Expected: PASS.

- [ ] **Step 5: Verify the repository and commit**

Run:

```bash
rtk make GO='mise exec -- go' check
rtk make GO='mise exec -- go' build
rtk git diff --check
```

Commit:

```bash
rtk git add internal/provider/resolver.go internal/provider/resolver_lid_test.go
rtk git commit -m "fix: accept LID private recipients"
```

### Task 2: Deploy to local Docker

**Files:**
- No source changes.

**Interfaces:**
- Consumes the committed provider image from Task 1.
- Produces a running local `worker-a` on port `8080` with existing session data.

- [ ] **Step 1: Build the local image from the worktree**

Run:

```bash
rtk docker build -t ghcr.io/mbap-dev/provider-whatsmeow:latest .
```

Expected: Dockerfile tests and both Go builds pass.

- [ ] **Step 2: Recreate only worker-a from the main compose project**

From `/Users/mbppereira/projetos/Envolve-conectores/provider-whatsmeow`, run:

```bash
rtk docker compose up -d --no-deps --no-build --force-recreate worker-a
```

- [ ] **Step 3: Verify readiness and session restoration**

Run:

```bash
rtk curl -i -fsS http://localhost:8080/readyz
rtk docker compose logs --since=2m worker-a
```

Expected: HTTP 200 with `ready`, no startup errors, and session `553121159080` successfully authenticated.
