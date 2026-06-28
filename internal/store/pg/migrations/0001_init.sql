-- +goose Up
CREATE TABLE IF NOT EXISTS worker_assignments (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

CREATE TABLE IF NOT EXISTS session_overrides (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    data       jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

CREATE TABLE IF NOT EXISTS session_webhooks (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    data       jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

CREATE TABLE IF NOT EXISTS session_templates (
    worker_id  text NOT NULL,
    session_id text NOT NULL,
    data       jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, session_id)
);

-- +goose Down
DROP TABLE IF EXISTS session_templates;
DROP TABLE IF EXISTS session_webhooks;
DROP TABLE IF EXISTS session_overrides;
DROP TABLE IF EXISTS worker_assignments;
