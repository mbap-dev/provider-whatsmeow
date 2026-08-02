# LID Private Recipient Resolution Design

## Problem

`whatsmeow.Client.IsOnWhatsApp` now returns the canonical private recipient as
an `@lid` JID when the account uses LID addressing. The provider currently
accepts only `@s.whatsapp.net` responses, so it reports registered contacts as
missing. A container restart clears the in-memory PN cache and exposes the
invalid rejection again.

## Design

Add one small resolver rule that accepts a successful `IsOnWhatsApp` response
only when `IsIn` is true and the canonical JID is a non-empty private PN
(`@s.whatsapp.net`) or LID (`@lid`) address. Both `resolvePNtoJID` and
`checkUserExistsOnWhatsApp` will use this rule and cache the canonical JID by
the normalized input phone number for 24 hours.

The provider will prefer the canonical `JID` returned by whatsmeow instead of
forcing `PhoneNumber`. This matches the current LID addressing mode and is
supported by `whatsmeow.SendMessage`. Group, broadcast, hosted, and malformed
JIDs remain excluded from this private-number path.

## Failure Behavior

- `IsIn=false` continues trying Brazilian number candidates.
- Network and empty-response errors keep the existing fallback behavior.
- A successful response with an unsupported or empty JID is not treated as a
  registered private recipient.
- Existing PN responses remain compatible.

## Verification

Introduce a test seam around `IsOnWhatsApp` and regression tests using phone
`5562933000233`, canonical JID `46849671094330@lid`, and PN JID
`5562933000233@s.whatsapp.net`. Verify that existence checking and legacy PN
resolution return the LID on the first request, cache it, and do not query the
alternate Brazilian candidate. Run `make check` and `make build`, then rebuild
only the local Docker worker while preserving its session volume.

## Scope

No changes to group enrichment, webhooks, media, receipts, or HistorySync.
