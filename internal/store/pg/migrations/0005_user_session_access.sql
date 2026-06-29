-- +goose Up
CREATE TABLE IF NOT EXISTS user_session_access (
    user_id    uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, session_id)
);

CREATE INDEX IF NOT EXISTS idx_user_session_access_session_id ON user_session_access(session_id);

-- +goose Down
DROP TABLE IF EXISTS user_session_access;
