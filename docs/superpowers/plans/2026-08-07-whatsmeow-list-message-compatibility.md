# WhatsMeow ListMessage Compatibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restore WhatsApp acceptance of legacy interactive lists while continuing to depend directly on official whatsmeow.

**Architecture:** Add a provider-owned startup compatibility shim that changes only the generated enum-name lookup used by official whatsmeow to derive the `biz/list` type. Keep the protobuf numeric value and all public payload contracts unchanged, and isolate the workaround in its own file with focused tests.

**Tech Stack:** Go, `go.mau.fi/whatsmeow`, protobuf-generated `waE2E` types, Docker Compose.

## Global Constraints

- Keep the official `go.mau.fi/whatsmeow` module without a local `replace`.
- Keep the list protobuf type numerically equal to `ListMessage_SINGLE_SELECT`.
- Emit `product_list` through the enum-name lookup used by whatsmeow.
- Do not change buttons, templates, payment messages, incoming parsing, or the public request contract.
- Make the shim idempotent and reject an unexpected upstream enum shape visibly.

---

### Task 1: Legacy List Biz-Type Compatibility

**Files:**
- Create: `internal/provider/whatsmeow_list_compat.go`
- Create: `internal/provider/whatsmeow_list_compat_test.go`
- Verify: `internal/provider/interactive_test.go`
- Verify: `go.mod`

**Interfaces:**
- Consumes: `waE2E.ListMessage_ListType_name map[int32]string` and numeric value `waE2E.ListMessage_SINGLE_SELECT`.
- Produces: `applyLegacyListBizType(map[int32]string) error`, called during package initialization.

- [x] **Step 1: Write the failing compatibility tests**

```go
func TestApplyLegacyListBizType(t *testing.T) {
	names := map[int32]string{int32(waE2E.ListMessage_SINGLE_SELECT): "SINGLE_SELECT"}
	if err := applyLegacyListBizType(names); err != nil {
		t.Fatalf("applyLegacyListBizType() error = %v", err)
	}
	if got := names[int32(waE2E.ListMessage_SINGLE_SELECT)]; got != "PRODUCT_LIST" {
		t.Fatalf("single-select enum name = %q, want PRODUCT_LIST", got)
	}
	if err := applyLegacyListBizType(names); err != nil {
		t.Fatalf("second applyLegacyListBizType() error = %v", err)
	}
}

func TestApplyLegacyListBizTypeRejectsUnexpectedUpstreamShape(t *testing.T) {
	names := map[int32]string{int32(waE2E.ListMessage_SINGLE_SELECT): "CHANGED_UPSTREAM"}
	if err := applyLegacyListBizType(names); err == nil {
		t.Fatal("applyLegacyListBizType() error = nil, want unsupported-shape error")
	}
}

func TestOfficialWhatsmeowUsesLegacyListBizType(t *testing.T) {
	if got := waE2E.ListMessage_ListType_name[int32(waE2E.ListMessage_SINGLE_SELECT)]; got != "PRODUCT_LIST" {
		t.Fatalf("runtime single-select enum name = %q, want PRODUCT_LIST", got)
	}
}
```

- [x] **Step 2: Run the focused tests and verify the red state**

Run:

```bash
rtk mise exec -- go test ./internal/provider -run 'TestApplyLegacyListBizType|TestOfficialWhatsmeowUsesLegacyListBizType' -count=1
```

Expected: compilation fails because `applyLegacyListBizType` is undefined.

- [x] **Step 3: Implement the startup shim**

```go
package provider

import (
	"fmt"

	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
)

const (
	officialSingleSelectName = "SINGLE_SELECT"
	legacyListBizTypeName    = "PRODUCT_LIST"
)

func init() {
	if err := applyLegacyListBizType(waE2E.ListMessage_ListType_name); err != nil {
		panic(err)
	}
}

func applyLegacyListBizType(names map[int32]string) error {
	key := int32(waE2E.ListMessage_SINGLE_SELECT)
	switch current := names[key]; current {
	case officialSingleSelectName:
		names[key] = legacyListBizTypeName
		return nil
	case legacyListBizTypeName:
		return nil
	default:
		return fmt.Errorf("unsupported whatsmeow ListMessage_SINGLE_SELECT enum name %q", current)
	}
}
```

- [x] **Step 4: Run focused compatibility and existing interactive tests**

Run:

```bash
rtk mise exec -- go test ./internal/provider -run 'TestApplyLegacyListBizType|TestOfficialWhatsmeowUsesLegacyListBizType|TestBuildInteractiveListMessage' -count=1
```

Expected: PASS. Existing list-builder assertions continue to prove that the protobuf numeric type remains `ListMessage_SINGLE_SELECT`.

- [x] **Step 5: Verify dependency and complete Go suite**

Run:

```bash
rtk go list -m -f '{{.Path}} {{.Version}} {{if .Replace}}=> {{.Replace.Path}}{{end}}' go.mau.fi/whatsmeow
rtk mise exec -- go test ./...
rtk git diff --check
```

Expected: official version is printed without `=>`; all tests and diff checks pass.

- [x] **Step 6: Rebuild and start the local Docker worker**

Run:

```bash
rtk docker compose up -d --build worker-a
rtk docker compose ps worker-a
rtk docker compose exec -T worker-a wget -qO- http://localhost:8080/healthz
```

Expected: `worker-a` is running and the health endpoint returns success.

- [x] **Step 7: Verify a real list send and inspect logs**

Submit the existing interactive-list payload through the provider's normal outgoing path, then run:

```bash
rtk docker compose logs --since=5m worker-a
```

Expected: `interactive_list_transport=legacy_list_message biz_node_type=product_list` followed by `interactive list sent`; no 405 rejection.

- [x] **Step 8: Commit the implementation**

```bash
rtk git add internal/provider/interactive_test.go internal/provider/whatsmeow_list_compat.go internal/provider/whatsmeow_list_compat_test.go docs/superpowers/plans/2026-08-07-whatsmeow-list-message-compatibility.md
rtk git commit -m "fix: restore legacy list message transport"
```
