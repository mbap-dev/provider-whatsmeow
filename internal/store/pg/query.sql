-- Media registry -------------------------------------------------------------

-- name: CountMediaRecords :one
SELECT COUNT(*) FROM media_records;

-- name: UpsertMediaRecord :exec
INSERT INTO media_records (
    session_id, media_id, messaging_product, mime_type, sha256,
    file_size, id, filename, relative_path, storage_key, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (session_id, media_id) DO UPDATE SET
    messaging_product = EXCLUDED.messaging_product,
    mime_type         = EXCLUDED.mime_type,
    sha256            = EXCLUDED.sha256,
    file_size         = EXCLUDED.file_size,
    id                = EXCLUDED.id,
    filename          = EXCLUDED.filename,
    relative_path     = EXCLUDED.relative_path,
    storage_key       = EXCLUDED.storage_key;

-- name: GetAllMediaRecords :many
SELECT session_id, media_id, messaging_product, mime_type, sha256,
       file_size, id, filename, relative_path, storage_key, created_at
FROM media_records
ORDER BY created_at;

-- Worker assignments -------------------------------------------------------

-- name: ListAssignments :many
SELECT session_id FROM worker_assignments WHERE worker_id = $1 ORDER BY session_id;

-- name: DeleteAssignmentsByWorker :exec
DELETE FROM worker_assignments WHERE worker_id = $1;

-- name: InsertAssignment :exec
INSERT INTO worker_assignments (worker_id, session_id)
VALUES ($1, $2)
ON CONFLICT (worker_id, session_id) DO NOTHING;

-- Session config overrides --------------------------------------------------

-- name: ListOverrides :many
SELECT session_id, data FROM session_overrides WHERE worker_id = $1;

-- name: DeleteOverridesByWorker :exec
DELETE FROM session_overrides WHERE worker_id = $1;

-- name: UpsertOverride :exec
INSERT INTO session_overrides (worker_id, session_id, data)
VALUES ($1, $2, $3)
ON CONFLICT (worker_id, session_id) DO UPDATE SET data = EXCLUDED.data, updated_at = now();

-- Session webhooks ----------------------------------------------------------

-- name: ListWebhooks :many
SELECT session_id, data FROM session_webhooks WHERE worker_id = $1;

-- name: DeleteWebhooksByWorker :exec
DELETE FROM session_webhooks WHERE worker_id = $1;

-- name: UpsertWebhooks :exec
INSERT INTO session_webhooks (worker_id, session_id, data)
VALUES ($1, $2, $3)
ON CONFLICT (worker_id, session_id) DO UPDATE SET data = EXCLUDED.data, updated_at = now();

-- Session templates ---------------------------------------------------------

-- name: ListTemplates :many
SELECT session_id, data FROM session_templates WHERE worker_id = $1;

-- name: DeleteTemplatesByWorker :exec
DELETE FROM session_templates WHERE worker_id = $1;

-- name: UpsertTemplates :exec
INSERT INTO session_templates (worker_id, session_id, data)
VALUES ($1, $2, $3)
ON CONFLICT (worker_id, session_id) DO UPDATE SET data = EXCLUDED.data, updated_at = now();

-- Auth: users -----------------------------------------------------------------

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- Auth: user_sessions ---------------------------------------------------------

-- name: CreateSession :one
INSERT INTO user_sessions (user_id, token_hash, expires_at, ip, user_agent)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT * FROM user_sessions
WHERE token_hash = $1 AND expires_at > now()
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM user_sessions WHERE token_hash = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM user_sessions WHERE expires_at <= now();

-- Auth: api_keys --------------------------------------------------------------

-- name: CreateAPIKey :one
INSERT INTO api_keys (user_id, key_hash, name, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAPIKeyByHash :one
SELECT * FROM api_keys
WHERE key_hash = $1
  AND (expires_at IS NULL OR expires_at > now())
LIMIT 1;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys SET last_used_at = now() WHERE id = $1;

-- name: ListAPIKeysByUser :many
SELECT * FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE id = $1 AND user_id = $2;
