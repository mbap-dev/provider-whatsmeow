-- +goose Up
CREATE TABLE IF NOT EXISTS session_pairing_links (
    worker_id  text        NOT NULL,
    token_hash text        NOT NULL,
    session_id text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    used_at    timestamptz,
    revoked_at timestamptz,
    PRIMARY KEY (worker_id, token_hash)
);

CREATE INDEX IF NOT EXISTS idx_session_pairing_links_session_id
    ON session_pairing_links(worker_id, session_id);

-- +goose Down
DROP TABLE IF EXISTS session_pairing_links;
