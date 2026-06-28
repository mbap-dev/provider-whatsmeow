-- Schema used by sqlc for type inference. Keep in sync with migrations/.
-- All metadata is keyed by worker_id (mirrors the historical per-worker JSON
-- files). Payloads are stored as jsonb matching the existing on-disk shapes.

CREATE TABLE worker_assignments (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

CREATE TABLE session_overrides (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    data       jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

CREATE TABLE session_webhooks (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    data       jsonb NOT NULL, -- JSON array of webhooks for the session
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

CREATE TABLE session_templates (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    data       jsonb NOT NULL, -- JSON array of templates for the session
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

-- Media registry (migration 0003) --

CREATE TABLE media_records (
    session_id        text        NOT NULL,
    media_id          text        NOT NULL,
    messaging_product text        NOT NULL DEFAULT 'whatsapp',
    mime_type         text,
    sha256            text,
    file_size         bigint,
    id                text        NOT NULL,
    filename          text,
    relative_path     text        NOT NULL,
    storage_key       text        NOT NULL,
    created_at        timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (session_id, media_id)
);

-- Auth tables (migration 0002) --

CREATE TABLE users (
    id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         text        NOT NULL UNIQUE,
    password_hash text        NOT NULL,
    role          text        NOT NULL DEFAULT 'user'
                              CHECK (role IN ('owner', 'admin', 'user')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE user_sessions (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  text        NOT NULL UNIQUE,
    expires_at  timestamptz NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    ip          text,
    user_agent  text
);

CREATE TABLE api_keys (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash     text        NOT NULL UNIQUE,
    name         text        NOT NULL,
    last_used_at timestamptz,
    created_at   timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz
);
