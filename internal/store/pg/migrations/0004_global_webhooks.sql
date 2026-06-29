-- +goose Up
CREATE TABLE IF NOT EXISTS global_webhooks (
    worker_id  text NOT NULL,
    webhook_id text NOT NULL,
    data       jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (worker_id, webhook_id)
);

-- +goose Down
DROP TABLE IF EXISTS global_webhooks;
