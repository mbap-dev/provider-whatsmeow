package provider

import (
	"context"
	"path/filepath"
)

// MetadataStore persists per-worker session metadata: worker assignments,
// per-session config overrides, registered webhooks and message templates.
//
// The default implementation (fileStore) keeps the historical JSON files under
// SESSION_STORE, so behaviour is unchanged when no external store is wired. A
// Postgres-backed implementation can be injected via ManagerOptions.MetadataStore.
//
// The model is snapshot-oriented (load the full per-worker map at startup, save
// the full map on each mutation), mirroring how the manager already works.
type MetadataStore interface {
	LoadAssignments(ctx context.Context, workerID string) ([]string, error)
	SaveAssignments(ctx context.Context, workerID string, sessions []string) error

	LoadOverrides(ctx context.Context, workerID string) (map[string]sessionOverrides, error)
	SaveOverrides(ctx context.Context, workerID string, overrides map[string]sessionOverrides) error

	LoadWebhooks(ctx context.Context, workerID string) (map[string][]SessionWebhook, error)
	SaveWebhooks(ctx context.Context, workerID string, webhooks map[string][]SessionWebhook) error

	LoadTemplates(ctx context.Context, workerID string) (map[string][]SessionTemplate, error)
	SaveTemplates(ctx context.Context, workerID string, templates map[string][]SessionTemplate) error
}

// fileStore is the default JSON-file-backed MetadataStore. It delegates to the
// existing load*/save* helpers, preserving the on-disk format under SESSION_STORE.
type fileStore struct {
	assignmentFile string
	overrideFile   string
	webhookFile    string
	templateFile   string
}

func newFileStore(assignmentFile, overrideFile, webhookFile, templateFile string) *fileStore {
	return &fileStore{
		assignmentFile: assignmentFile,
		overrideFile:   overrideFile,
		webhookFile:    webhookFile,
		templateFile:   templateFile,
	}
}

// NewFileStore builds the default JSON-file MetadataStore for a session store
// directory, using the historical filenames.
func NewFileStore(sessionStore string) MetadataStore {
	return newFileStore(
		filepath.Join(sessionStore, "_worker_assignments.json"),
		filepath.Join(sessionStore, "_session_overrides.json"),
		filepath.Join(sessionStore, "_session_webhooks.json"),
		filepath.Join(sessionStore, "_session_templates.json"),
	)
}

// ImportFileStateIfEmpty copies each metadata category from src into dst for the
// given worker, but only when dst has no data for that category yet. It is used
// for a one-time, non-destructive migration of JSON state into Postgres on first
// boot. Returns true if anything was imported.
func ImportFileStateIfEmpty(ctx context.Context, src, dst MetadataStore, workerID string) (bool, error) {
	imported := false

	if cur, err := dst.LoadAssignments(ctx, workerID); err == nil && len(cur) == 0 {
		if from, err := src.LoadAssignments(ctx, workerID); err == nil && len(from) > 0 {
			if err := dst.SaveAssignments(ctx, workerID, from); err != nil {
				return imported, err
			}
			imported = true
		}
	}
	if cur, err := dst.LoadOverrides(ctx, workerID); err == nil && len(cur) == 0 {
		if from, err := src.LoadOverrides(ctx, workerID); err == nil && len(from) > 0 {
			if err := dst.SaveOverrides(ctx, workerID, from); err != nil {
				return imported, err
			}
			imported = true
		}
	}
	if cur, err := dst.LoadWebhooks(ctx, workerID); err == nil && len(cur) == 0 {
		if from, err := src.LoadWebhooks(ctx, workerID); err == nil && len(from) > 0 {
			if err := dst.SaveWebhooks(ctx, workerID, from); err != nil {
				return imported, err
			}
			imported = true
		}
	}
	if cur, err := dst.LoadTemplates(ctx, workerID); err == nil && len(cur) == 0 {
		if from, err := src.LoadTemplates(ctx, workerID); err == nil && len(from) > 0 {
			if err := dst.SaveTemplates(ctx, workerID, from); err != nil {
				return imported, err
			}
			imported = true
		}
	}
	return imported, nil
}

func (s *fileStore) LoadAssignments(_ context.Context, workerID string) ([]string, error) {
	return loadAssignments(s.assignmentFile, workerID)
}

func (s *fileStore) SaveAssignments(_ context.Context, workerID string, sessions []string) error {
	return saveAssignments(s.assignmentFile, workerID, sessions)
}

func (s *fileStore) LoadOverrides(_ context.Context, workerID string) (map[string]sessionOverrides, error) {
	return loadSessionOverrideState(s.overrideFile, workerID)
}

func (s *fileStore) SaveOverrides(_ context.Context, workerID string, overrides map[string]sessionOverrides) error {
	return saveSessionOverrideState(s.overrideFile, workerID, overrides)
}

func (s *fileStore) LoadWebhooks(_ context.Context, workerID string) (map[string][]SessionWebhook, error) {
	return loadSessionWebhookState(s.webhookFile, workerID)
}

func (s *fileStore) SaveWebhooks(_ context.Context, workerID string, webhooks map[string][]SessionWebhook) error {
	return saveSessionWebhookState(s.webhookFile, workerID, webhooks)
}

func (s *fileStore) LoadTemplates(_ context.Context, workerID string) (map[string][]SessionTemplate, error) {
	return loadSessionTemplateState(s.templateFile, workerID)
}

func (s *fileStore) SaveTemplates(_ context.Context, workerID string, templates map[string][]SessionTemplate) error {
	return saveSessionTemplateState(s.templateFile, workerID, templates)
}
