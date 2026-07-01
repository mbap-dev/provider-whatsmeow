# WhatsApp Passkey Pairing Debug Plan

## Context

Meta is gradually adding an extra security step after the normal WhatsApp Web QR scan. Some accounts may require a passkey/access-key confirmation before the websocket pairing is completed.

The whatsmeow pairing flow considers the device paired only after the socket receives `pair-success`, verifies the device identity, saves the device store, and sends `pair-device-sign` back to WhatsApp. Accounts that require passkeys now use the upstream `whatsmeow` WebAuthn flow added in `go.mau.fi/whatsmeow v0.0.0-20260630180629-b572e5bcb92b`.

## Goal

Capture the first affected account flow with maximum useful detail while keeping this isolated from the main provider application.

## Phase 1: Standalone Capture

Use `cmd/debug-pairing` to simulate a fresh whatsmeow connection:

- create an isolated SQLite store;
- enable `DEBUG` logs for whatsmeow database/client/QR channel;
- call `GetQRChannel()` before `Connect()`;
- save every QR as PNG for easy scanning;
- print every exported whatsmeow event with type, JSON where possible, and Go representation;
- keep the process alive after QR scan to catch `pair-success`, `PairError`, `StreamError`, `ConnectFailure`, `LoggedOut`, or any unexpected event.

Run example:

```sh
mise exec -- go run ./cmd/debug-pairing \
  --session passkey-test \
  --store ./state/passkey-debug \
  --qr-png ./state/passkey-debug/qr.png \
  --log-level DEBUG
```

Use a clean `--session` or remove the session directory to force a new QR. The QR PNG is rewritten whenever whatsmeow rotates the QR code.

## Phase 2: Identify Protocol Gap

From the captured output, determine where the passkey challenge appears:

- as an unknown `stream:error`;
- as a `connect_failure` node with raw details;
- as a `PairError`;
- as a changed `pair-success` payload;
- as an unhandled IQ/IB node only visible in debug logs.

If whatsmeow already exposes enough raw data, implement the response in our code. If the library drops the node before surfacing it, patch or update whatsmeow to expose that node safely.

Initial capture on an affected account showed the challenge as an unhandled notification:

```xml
<notification type="passkey_prologue_request">
  <passkey_request_options>{
    "challenge":"<redacted>",
    "timeout":600000,
    "rpId":"whatsapp.com",
    "allowCredentials":[],
    "userVerification":"required",
    "extensions":{"uvm":true}
  }</passkey_request_options>
</notification>
```

The official whatsmeow implementation emits `events.PairPasskeyRequest` for this notification. The provider stores the marshaled `types.WebAuthnPublicKey` as `passkey.request_options_json` so session status and the public pairing page can show that the connection is waiting for the passkey step instead of silently timing out.

## Phase 3: Browser Passkey Bridge

The public pairing page cannot directly execute the WebAuthn request because the challenge uses `rpId: "whatsapp.com"`. The browser-mediated FIDO step must run in an origin accepted for WhatsApp Web.

Planned Chrome extension:

- Content script runs on `https://web.whatsapp.com/*`.
- Extension asks the provider for the target session's `passkey.request_options_json`.
- Extension calls `navigator.credentials.get({ publicKey })` in the WhatsApp Web context.
- Extension posts the resulting assertion to `/api/{version}/sessions/{id}/passkey-prologue`.
- If `passkey.confirmation_required=true`, the user confirms the shown code through the secure pairing link or a dedicated extension UI that requires explicit user confirmation.

Implemented packaging:

- `/pair/{token}/extension.zip` returns a Chrome extension package with `config.json` scoped to the current provider URL.
- The extension runs a provider-page content script on `/pair/{token}` to store the current token in `chrome.storage.local`, so the same installed extension can be reused for later links from the same provider.
- `/pair/{token}/passkey-prologue` receives the real assertion from that extension and forwards it to upstream `Client.SendPasskeyResponse`.
- `/pair/{token}/passkey-test` receives a real browser WebAuthn assertion for transport validation only, using the currently stored token. It does not call WhatsApp or mutate pairing state.

The public pairing page remains responsible for:

- QR scan remains the first step;
- showing that the socket is waiting for passkey;
- displaying `passkey.confirmation_code` when upstream emits `PairPasskeyConfirmation` with `SkipHandoffUX=false`;
- confirming the code through `/pair/{token}/passkey-confirmation`;
- expiring the public link only after real `PairSuccess`/login.

Important constraint: the captured options use `rpId: "whatsapp.com"`. Browser WebAuthn validation normally only allows a page to request credentials for its own effective domain or a registrable-domain suffix of it. A pairing page hosted on our own domain is therefore expected to fail if it calls `navigator.credentials.get()` with `rpId: "whatsapp.com"`, unless WhatsApp explicitly supports related-origin requests that include our origin.

Product note: a normal public pairing page on our own domain cannot directly use the WhatsApp passkey. The extension is the production path unless WhatsApp later exposes a supported related-origin flow.

## Phase 4: WhatsApp Web Browser Capture

To compare the official browser flow with whatsmeow, use the Chrome DevTools Protocol capture script. It injects a hook into `navigator.credentials.get()` on WhatsApp Web and records:

- WebAuthn request options passed by WhatsApp Web;
- WebAuthn assertion returned by the browser/passkey flow;
- WebSocket send summaries immediately before/after the passkey result;
- optional raw WebSocket frames when `--include-ws-payloads` is set.

Start an isolated Chrome profile with remote debugging:

```sh
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --remote-debugging-port=9222 \
  --user-data-dir=/tmp/wa-passkey-chrome \
  --no-first-run \
  https://web.whatsapp.com/
```

Then run the capture:

```sh
node scripts/debug-whatsapp-web-passkey.mjs \
  --port 9222 \
  --out ./state/wa-web-passkey.ndjson \
  --dump-sources ./state/wa-web-sources
```

Complete the official WhatsApp Web connection, including the browser passkey prompt. The important records are `injected.passkey_debug` with `payload.kind` equal to `credentials.get.call`, `credentials.get.result`, and nearby `websocket.send`.

Raw WebSocket frames may be encrypted by WhatsApp Web, so the WebAuthn assertion captured from `navigator.credentials.get()` is expected to be the most useful artifact for implementing the whatsmeow response. Use the files in `--dump-sources` to inspect snippets from the WhatsApp Web bundles that reference `passkey`, `credentials.get`, `webauthn`, or related terms; those snippets are the next place to identify the internal node sent before socket encryption.

Captured WhatsApp Web bundle mapping:

- UI action: `WAWebShortcakeLinkingAction.startShortcakePasskeyPairing()`
- main flow: `WAWebShortcakeLinkingPasskeyPrologue.executePasskeyPrologue()`
- WebAuthn call: `performWebAuthnAssertion(requestOptionsBytes)`
- outgoing RPC: `WASmaxMdSetPasskeyPrologueRPC.sendSetPasskeyPrologueRPC(...)`
- outgoing node builder: `WASmaxOutMdSetPasskeyPrologueRequest.makeSetPasskeyPrologueRequest(...)`

The outgoing pre-encryption stanza shape is:

```xml
<iq to="s.whatsapp.net" xmlns="md" type="set" id="...">
  <passkey_prologue>
    <credential_id>...</credential_id>
    <webauthn_assertion>...</webauthn_assertion>
    <prologue_payload>...</prologue_payload>
    <pairing_handoff_proof>...</pairing_handoff_proof> <!-- optional -->
  </passkey_prologue>
</iq>
```

Field sources from the browser flow:

- `credential_id`: decoded `PublicKeyCredential.rawId` bytes; the provider accepts the JSON assertion and lets upstream whatsmeow derive this from `types.WebAuthnResponse.RawID`
- `webauthn_assertion`: JSON bytes containing `id`, `rawId`, `type`, and `response.{clientDataJSON,authenticatorData,signature,userHandle}` with binary fields base64url-encoded
- `prologue_payload`: generated by upstream `Client.SendPasskeyResponse` as `WACompanionReg.ProloguePayload`
- `pairing_handoff_proof`: generated by upstream whatsmeow when a valid handoff key exists

## Phase 5: Backend Prologue Submission

The provider uses the official upstream flow from [tulir/whatsmeow@b572e5b](https://github.com/tulir/whatsmeow/commit/b572e5bcb92bbc285b68cb6d6540da3093330e04), which is also consumed by [mautrix/whatsapp@be4adcbd](https://github.com/mautrix/whatsapp/commit/be4adcbd1712b6f51eba87411780079cb1c0ae52).

The server-side sequence is now:

1. `events.PairPasskeyRequest` exposes `*types.WebAuthnPublicKey`.
2. Browser/extension returns `types.WebAuthnResponse`.
3. Provider calls `Client.SendPasskeyResponse(ctx, &resp)`.
4. Upstream sends `passkey_prologue` with `credential_id`, `webauthn_assertion`, `prologue_payload`, and optional `pairing_handoff_proof`.
5. Upstream handles `crsc_continuation`, sends `companion_nonce`, derives the encryption key, and emits `events.PairPasskeyConfirmation`.
6. If `SkipHandoffUX=true`, provider calls `Client.SendPasskeyConfirmation(ctx)` automatically.
7. If `SkipHandoffUX=false`, the public pairing page shows `PairPasskeyConfirmation.Code`; after explicit confirmation, `/pair/{token}/passkey-confirmation` calls `Client.SendPasskeyConfirmation(ctx)`.

`SetPasskeyPrologue` still has this pre-encryption shape:

```xml
<iq to="s.whatsapp.net" xmlns="md" type="set">
  <passkey_prologue>
    <credential_id>...</credential_id>
    <webauthn_assertion>...</webauthn_assertion>
    <prologue_payload>...</prologue_payload>
    <pairing_handoff_proof>...</pairing_handoff_proof>
  </passkey_prologue>
</iq>
```

The provider accepts the WebAuthn assertion through an authenticated experimental endpoint:

```sh
curl -X POST http://localhost:8080/api/v1/sessions/5511999999999/passkey-prologue \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer wm_...' \
  -d '{
    "credential_id": "...",
    "webauthn_assertion": {
      "id": "...",
      "rawId": "...",
      "type": "public-key",
      "response": {
        "clientDataJSON": "...",
        "authenticatorData": "...",
        "signature": "...",
        "userHandle": null
      }
    }
  }'
```

Do not reuse captured assertions across sessions. The assertion is bound to the exact challenge received in `passkey_request_options`, and the endpoint only starts the server-side Shortcake flow after a valid browser/WebAuthn context has produced that assertion.

If the provider status shows `passkey.confirmation_required=true`, confirm after comparing the displayed code with the phone:

```sh
curl -X POST http://localhost:8080/api/v1/sessions/5511999999999/passkey-confirmation \
  -H 'Authorization: Bearer wm_...'
```

## Phase 6: Temporary Chrome Debug Bridge

Before building the Chrome extension, the CDP script can test the same idea manually:

1. Start the provider and create/connect the session until it reaches `passkey_required`.
2. Start Chrome with remote debugging open on WhatsApp Web:

```sh
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --remote-debugging-port=9222 \
  --user-data-dir=/tmp/wa-passkey-chrome \
  --no-first-run \
  https://web.whatsapp.com/
```

3. Run the bridge:

```sh
WHATSMEOW_API_TOKEN=wm_... \
mise exec -- node scripts/debug-whatsapp-web-passkey.mjs \
  --bridge-passkey \
  --provider-url http://localhost:8080 \
  --session 5511999999999 \
  --api-version v1 \
  --out ./state/wa-passkey-bridge.ndjson
```

The script fetches `/sessions/{id}/status` from the provider, extracts `passkey.request_options_json`, executes `navigator.credentials.get()` inside the `web.whatsapp.com` tab using Chrome DevTools Protocol, and posts the assertion to `/api/v1/sessions/{id}/passkey-prologue`. If the status then exposes `passkey.confirmation_required=true`, use the secure pairing page or `/api/v1/sessions/{id}/passkey-confirmation` to finish the pairing.

If Chrome returns `NotAllowedError`, bring the WhatsApp Web tab to the front and rerun the bridge. The assertion cannot be reused; rerun only while the same session is still waiting on the same passkey challenge.

## Capture Notes

Do not paste production QR payloads, passkey challenge responses, WebAuthn assertions, browser capture logs, or session database files into tickets or chats. The commands print raw protocol details for debugging and may store real credentials or signatures.
