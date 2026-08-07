# WhatsMeow ListMessage Compatibility Design

## Problem

The provider builds interactive lists as legacy `ListMessage` protobufs. The
removed local whatsmeow fork forced the automatically generated `biz/list`
node to use `type="product_list"` while keeping the protobuf list type as
`SINGLE_SELECT`.

After returning to the official whatsmeow module, the library derives the biz
attribute from the protobuf enum and emits `type="single_select"`. WhatsApp
rejects this transport shape with status 405. Existing tests only inspect the
protobuf and therefore did not detect the regression.

## Design

Keep the official `go.mau.fi/whatsmeow` module without a local `replace`.
Install a provider-owned compatibility shim before clients can send messages.
The shim changes the generated enum name used by whatsmeow for
`ListMessage_SINGLE_SELECT` from `SINGLE_SELECT` to `PRODUCT_LIST`. The numeric
protobuf value remains unchanged, so list contents and receiving behavior stay
the same; only the derived biz-node attribute returns to the previously working
wire shape.

The shim must be idempotent and process-wide because the generated enum-name
map is package-global. It must validate the expected upstream enum entry before
changing it, so a future dependency update fails visibly instead of silently
applying the workaround to an incompatible version.

Buttons, templates, payment messages, incoming interactive parsing, and the
public request contract are outside this change.

## Verification

Add focused tests that verify:

- the upstream enum starts from the supported shape;
- applying the shim maps `SINGLE_SELECT` to `PRODUCT_LIST`;
- applying it repeatedly is safe;
- `buildInteractiveListMessage` still produces numeric
  `ListMessage_SINGLE_SELECT` and preserves sections, rows, header, body, and
  footer.

Run the focused provider tests and the complete Go test suite. Rebuild and
start the local Docker worker, confirm health, then send a real interactive
list and verify that it is accepted and the provider logs a successful send.

## Rollback And Follow-up

The rollback is removal of the shim and its initialization call. Long term,
the workaround can be removed if official whatsmeow exposes a supported biz
attribute override or emits a WhatsApp-accepted node for this legacy list
shape.
