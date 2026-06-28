-- +goose Up
CREATE TABLE IF NOT EXISTS media_records (
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

CREATE INDEX IF NOT EXISTS idx_media_records_media_id ON media_records(media_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_media_records_path ON media_records(relative_path);

-- +goose Down
DROP TABLE IF EXISTS media_records;
