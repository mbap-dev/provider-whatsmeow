package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"your.org/provider-whatsmeow/internal/store/pg/pgdb"
)

// pgStore is the Postgres-backed MetadataStore. It mirrors the snapshot model of
// the file store: loads the full per-worker set, and on save replaces it inside a
// transaction (delete-by-worker + upsert each session). Payloads are stored as
// jsonb matching the historical on-disk shapes, so migrating from JSON is lossless.
type pgStore struct {
	pool *pgxpool.Pool
	q    *pgdb.Queries
}

// NewPostgresStore builds a MetadataStore backed by the given pgx pool.
func NewPostgresStore(pool *pgxpool.Pool) MetadataStore {
	return &pgStore{pool: pool, q: pgdb.New(pool)}
}

func (s *pgStore) withTx(ctx context.Context, fn func(*pgdb.Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := fn(s.q.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// --- Assignments ------------------------------------------------------------

func (s *pgStore) LoadAssignments(ctx context.Context, workerID string) ([]string, error) {
	return s.q.ListAssignments(ctx, workerID)
}

func (s *pgStore) SaveAssignments(ctx context.Context, workerID string, sessions []string) error {
	return s.withTx(ctx, func(q *pgdb.Queries) error {
		if err := q.DeleteAssignmentsByWorker(ctx, workerID); err != nil {
			return err
		}
		for _, sid := range sessions {
			sid = normalizeSessionID(sid)
			if sid == "" {
				continue
			}
			if err := q.InsertAssignment(ctx, pgdb.InsertAssignmentParams{WorkerID: workerID, SessionID: sid}); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- Overrides --------------------------------------------------------------

func (s *pgStore) LoadOverrides(ctx context.Context, workerID string) (map[string]sessionOverrides, error) {
	rows, err := s.q.ListOverrides(ctx, workerID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]sessionOverrides, len(rows))
	for _, row := range rows {
		sid := normalizeSessionID(row.SessionID)
		if sid == "" {
			continue
		}
		var raw sessionOverridePayload
		if err := json.Unmarshal(row.Data, &raw); err != nil {
			return nil, fmt.Errorf("decode override %s: %w", sid, err)
		}
		out[sid] = payloadToOverrides(raw)
	}
	return out, nil
}

func (s *pgStore) SaveOverrides(ctx context.Context, workerID string, overrides map[string]sessionOverrides) error {
	return s.withTx(ctx, func(q *pgdb.Queries) error {
		if err := q.DeleteOverridesByWorker(ctx, workerID); err != nil {
			return err
		}
		for sid, ov := range overrides {
			sid = normalizeSessionID(sid)
			if sid == "" {
				continue
			}
			data, err := json.Marshal(overridesToPayload(ov))
			if err != nil {
				return err
			}
			if err := q.UpsertOverride(ctx, pgdb.UpsertOverrideParams{WorkerID: workerID, SessionID: sid, Data: data}); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- Webhooks ---------------------------------------------------------------

func (s *pgStore) LoadWebhooks(ctx context.Context, workerID string) (map[string][]SessionWebhook, error) {
	rows, err := s.q.ListWebhooks(ctx, workerID)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]SessionWebhook, len(rows))
	for _, row := range rows {
		sid := normalizeSessionID(row.SessionID)
		if sid == "" {
			continue
		}
		var hooks []SessionWebhook
		if err := json.Unmarshal(row.Data, &hooks); err != nil {
			return nil, fmt.Errorf("decode webhooks %s: %w", sid, err)
		}
		out[sid] = hooks
	}
	return out, nil
}

func (s *pgStore) SaveWebhooks(ctx context.Context, workerID string, webhooks map[string][]SessionWebhook) error {
	return s.withTx(ctx, func(q *pgdb.Queries) error {
		if err := q.DeleteWebhooksByWorker(ctx, workerID); err != nil {
			return err
		}
		for sid, hooks := range webhooks {
			sid = normalizeSessionID(sid)
			if sid == "" {
				continue
			}
			data, err := json.Marshal(hooks)
			if err != nil {
				return err
			}
			if err := q.UpsertWebhooks(ctx, pgdb.UpsertWebhooksParams{WorkerID: workerID, SessionID: sid, Data: data}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *pgStore) LoadGlobalWebhooks(ctx context.Context, workerID string) ([]SessionWebhook, error) {
	rows, err := s.q.ListGlobalWebhooks(ctx, workerID)
	if err != nil {
		return nil, err
	}
	out := make([]SessionWebhook, 0, len(rows))
	for _, row := range rows {
		var hook SessionWebhook
		if err := json.Unmarshal(row.Data, &hook); err != nil {
			return nil, fmt.Errorf("decode global webhook %s: %w", row.WebhookID, err)
		}
		norm, err := normalizeGlobalWebhook(hook)
		if err != nil {
			continue
		}
		out = append(out, norm)
	}
	return out, nil
}

func (s *pgStore) SaveGlobalWebhooks(ctx context.Context, workerID string, webhooks []SessionWebhook) error {
	return s.withTx(ctx, func(q *pgdb.Queries) error {
		if err := q.DeleteGlobalWebhooksByWorker(ctx, workerID); err != nil {
			return err
		}
		for _, hook := range webhooks {
			norm, err := normalizeGlobalWebhook(hook)
			if err != nil {
				continue
			}
			data, err := json.Marshal(norm)
			if err != nil {
				return err
			}
			if err := q.UpsertGlobalWebhook(ctx, pgdb.UpsertGlobalWebhookParams{WorkerID: workerID, WebhookID: norm.ID, Data: data}); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- Templates --------------------------------------------------------------

func (s *pgStore) LoadTemplates(ctx context.Context, workerID string) (map[string][]SessionTemplate, error) {
	rows, err := s.q.ListTemplates(ctx, workerID)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]SessionTemplate, len(rows))
	for _, row := range rows {
		sid := normalizeSessionID(row.SessionID)
		if sid == "" {
			continue
		}
		var tpls []SessionTemplate
		if err := json.Unmarshal(row.Data, &tpls); err != nil {
			return nil, fmt.Errorf("decode templates %s: %w", sid, err)
		}
		out[sid] = tpls
	}
	return out, nil
}

func (s *pgStore) SaveTemplates(ctx context.Context, workerID string, templates map[string][]SessionTemplate) error {
	return s.withTx(ctx, func(q *pgdb.Queries) error {
		if err := q.DeleteTemplatesByWorker(ctx, workerID); err != nil {
			return err
		}
		for sid, tpls := range templates {
			sid = normalizeSessionID(sid)
			if sid == "" {
				continue
			}
			data, err := json.Marshal(tpls)
			if err != nil {
				return err
			}
			if err := q.UpsertTemplates(ctx, pgdb.UpsertTemplatesParams{WorkerID: workerID, SessionID: sid, Data: data}); err != nil {
				return err
			}
		}
		return nil
	})
}

// Ensure pgStore satisfies the interface at compile time.
var _ MetadataStore = (*pgStore)(nil)
