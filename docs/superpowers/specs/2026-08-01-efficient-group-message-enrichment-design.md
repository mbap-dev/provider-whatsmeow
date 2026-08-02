# Efficient Group Message Enrichment Design

## Goal

Keep the Chatwoot group-message contract while removing work that scales with
the total number of group participants. A live group message must enrich only
the group and the message author.

## Context

The current inbound group-message path calls `groupParticipantsPayload` for
every message. That fetches the full group metadata, iterates every member, and
requests every member's profile picture. Group administrative events perform
the same full-group enrichment. This creates `messages x participants` work and
can produce thousands of WhatsApp profile-picture requests after pairing.

The Chatwoot contract does not require `group_participants`. When the field is
absent, Chatwoot creates or updates the current author from `contacts[0]`,
creates or updates the group from the `group_*` fields, and incrementally adds
the author to the conversation's `group_contacts` association.

## Scope

This change applies to live webhook handling in `provider-whatsmeow`.

In scope:

- inbound group messages, including edits and media messages;
- group name and group picture enrichment;
- current message-author name and picture enrichment;
- ignoring group membership and role-change events;
- bounded, cached profile-picture and group-metadata lookups.

Out of scope:

- `HistorySync` behavior;
- eager or periodic synchronization of complete group rosters;
- changes to the Chatwoot application;
- outbound message behavior unrelated to webhook contact enrichment.

## Webhook Contract

An inbound group message will continue to send one contact entry with:

```json
{
  "contacts": [
    {
      "wa_id": "5562999999999",
      "profile": {
        "name": "Message Author",
        "picture": "https://example.test/author.jpg"
      },
      "group_id": "120363000000000000@g.us",
      "group_subject": "Group Name",
      "group_picture": "https://example.test/group.jpg"
    }
  ]
}
```

`group_participants` will be omitted. Authors who have never sent a message will
not be created in Chatwoot. Each author will be created or associated with the
group when Chatwoot receives that author's first message.

Missing names or pictures must not block message delivery:

- author name falls back to the resolved phone number or JID;
- group name falls back to the group JID user portion;
- missing author or group pictures are sent as empty/omitted values according
  to the existing payload shape.

## Event Flow

For each live inbound group message:

1. Parse the message and resolve the group and sender JIDs.
2. Resolve the author name from the message push name, local contact store, or
   identifier fallback.
3. Resolve the group name from cache; perform one `GetGroupInfo` request only on
   cache miss.
4. Resolve the group picture from cache; perform at most one profile-picture
   request on cache miss.
5. Resolve the current author's picture from cache; perform at most one
   profile-picture request on cache miss.
6. Build the Chatwoot-compatible contact and message payload without calling
   `groupParticipantsPayload`.
7. Deliver the webhook even if any optional enrichment fails.

Receipt and deletion events retain their current lightweight paths. Edited
messages keep using the normal message path, but therefore receive the same
bounded enrichment behavior.

## Lookup Controls

Profile-picture lookup behavior will be changed from multiple immediate
attempts to one request per cache miss:

- use the preview lookup suitable for contact and group avatars;
- do not perform `default`, `preview`, and `preview_retry` sequences;
- cache successful lookups for one hour;
- cache hidden, absent, and failed lookups for six hours;
- deduplicate concurrent lookups for the same JID;
- skip new lookups when the client is disconnected;
- return no picture immediately when the bounded lookup fails.

Group-name lookup behavior will use one `GetGroupInfo` request per cache miss,
cache successful names for one hour, and fall back without retrying when the
lookup fails. Concurrent cache misses for the same group must be deduplicated.

These controls bound live WhatsApp requests by active authors and active
groups, rather than by all group members.

## Administrative Group Events

`GroupInfo` membership and role changes (`add`, `remove`, `promote`, and
`demote`) and `JoinedGroup` events will not generate Chatwoot webhooks in this
scope. They must not trigger group metadata, participant, or picture lookups.

## Error Handling

- Optional enrichment errors are logged at debug level and do not fail the
  incoming message webhook.
- A disconnected client is treated as an immediate cache miss, with no network
  attempt.
- Webhook delivery errors retain the existing behavior.
- Media download behavior retains the existing behavior and timeout.

## Testing

Tests will verify that:

- an inbound group message omits `group_participants`;
- the author name and optional picture remain in `contacts[0].profile`;
- group subject and optional picture remain in the contact payload;
- group participant enumeration is not invoked by normal, edited, or media
  group messages;
- a profile-picture cache miss performs at most one lookup;
- successful and failed lookups are cached;
- concurrent requests for the same JID are deduplicated;
- disconnected clients do not start profile-picture lookups;
- group administrative events do not deliver webhooks or perform enrichment;
- receipt and deletion behavior remains unchanged.

## Acceptance Criteria

- Processing one live group message is constant with respect to group size.
- No live message or administrative event enumerates all group participants.
- No participant other than the current message author has a profile picture
  requested.
- At most one author-picture request and one group-picture request occur on
  their respective cache misses.
- Chatwoot creates the group conversation and incrementally associates each
  author who sends a message.
- Existing receipt, edit, deletion, and media webhook contracts continue to
  work.
