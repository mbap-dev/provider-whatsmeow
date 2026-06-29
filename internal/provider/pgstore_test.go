package provider

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"your.org/provider-whatsmeow/internal/store/pg"
)

// newTestPGStore connects to POSTGRES_TEST_URL, applies migrations and returns a
// clean store. Tests are skipped when the env var is absent so `go test ./...`
// stays green without a database.
func newTestPGStore(t *testing.T) MetadataStore {
	t.Helper()
	dsn := os.Getenv("POSTGRES_TEST_URL")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_URL not set; skipping Postgres store tests")
	}
	ctx := context.Background()
	if err := pg.Migrate(ctx, dsn); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	pool, err := pg.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	if _, err := pool.Exec(ctx, "TRUNCATE worker_assignments, session_overrides, session_webhooks, global_webhooks, session_templates"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "TRUNCATE worker_assignments, session_overrides, session_webhooks, global_webhooks, session_templates")
		pool.Close()
	})
	return NewPostgresStore(pool)
}

func TestPGStoreAssignmentsRoundTrip(t *testing.T) {
	s := newTestPGStore(t)
	ctx := context.Background()
	const w = "worker-test"

	if err := s.SaveAssignments(ctx, w, []string{"5511999999999", "5511888888888"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := s.LoadAssignments(ctx, w)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 assignments, got %v", got)
	}

	// Replace semantics: saving a smaller set removes the missing entries.
	if err := s.SaveAssignments(ctx, w, []string{"5511999999999"}); err != nil {
		t.Fatalf("save2: %v", err)
	}
	got, _ = s.LoadAssignments(ctx, w)
	if !reflect.DeepEqual(got, []string{"5511999999999"}) {
		t.Fatalf("expected replace, got %v", got)
	}
}

func TestPGStoreOverridesRoundTrip(t *testing.T) {
	s := newTestPGStore(t)
	ctx := context.Background()
	const w = "worker-test"

	reject := true
	msg := "no calls"
	every := 90 * time.Second
	in := map[string]sessionOverrides{
		"5511999999999": {
			rejectCalls:       &reject,
			rejectMsg:         &msg,
			alwaysOnlineEvery: &every,
		},
	}
	if err := s.SaveOverrides(ctx, w, in); err != nil {
		t.Fatalf("save: %v", err)
	}
	out, err := s.LoadOverrides(ctx, w)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("override round trip mismatch:\n in=%+v\nout=%+v", in["5511999999999"], out["5511999999999"])
	}
}

func TestPGStoreWebhooksAndTemplatesRoundTrip(t *testing.T) {
	s := newTestPGStore(t)
	ctx := context.Background()
	const w = "worker-test"

	hooks := map[string][]SessionWebhook{
		"5511999999999": {{ID: "h1", URL: "https://x", Enabled: true, TimeoutMS: 5000}},
	}
	if err := s.SaveWebhooks(ctx, w, hooks); err != nil {
		t.Fatalf("save webhooks: %v", err)
	}
	gotHooks, err := s.LoadWebhooks(ctx, w)
	if err != nil {
		t.Fatalf("load webhooks: %v", err)
	}
	if !reflect.DeepEqual(hooks, gotHooks) {
		t.Fatalf("webhooks mismatch: %+v vs %+v", hooks, gotHooks)
	}

	globalHooks := []SessionWebhook{
		{ID: "global-1", URLAbsolute: "https://x/global", Enabled: true, TimeoutMS: 5000, SendSessionDisconnected: true},
	}
	if err := s.SaveGlobalWebhooks(ctx, w, globalHooks); err != nil {
		t.Fatalf("save global webhooks: %v", err)
	}
	gotGlobalHooks, err := s.LoadGlobalWebhooks(ctx, w)
	if err != nil {
		t.Fatalf("load global webhooks: %v", err)
	}
	if !reflect.DeepEqual(globalHooks, gotGlobalHooks) {
		t.Fatalf("global webhooks mismatch: %+v vs %+v", globalHooks, gotGlobalHooks)
	}

	tpls := map[string][]SessionTemplate{
		"5511999999999": {{ID: "t1", Name: "welcome", Body: "hi"}},
	}
	if err := s.SaveTemplates(ctx, w, tpls); err != nil {
		t.Fatalf("save templates: %v", err)
	}
	gotTpls, err := s.LoadTemplates(ctx, w)
	if err != nil {
		t.Fatalf("load templates: %v", err)
	}
	if !reflect.DeepEqual(tpls, gotTpls) {
		t.Fatalf("templates mismatch: %+v vs %+v", tpls, gotTpls)
	}
}

func TestImportFileStateIfEmpty(t *testing.T) {
	dst := newTestPGStore(t)
	ctx := context.Background()
	const w = "worker-test"

	src := NewFileStore(t.TempDir())
	if err := src.SaveAssignments(ctx, w, []string{"5511777777777"}); err != nil {
		t.Fatalf("seed src: %v", err)
	}

	imported, err := ImportFileStateIfEmpty(ctx, src, dst, w)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if !imported {
		t.Fatal("expected import to copy assignments")
	}
	got, _ := dst.LoadAssignments(ctx, w)
	if !reflect.DeepEqual(got, []string{"5511777777777"}) {
		t.Fatalf("imported assignments mismatch: %v", got)
	}

	// Second run is a no-op because dst already has data.
	imported, err = ImportFileStateIfEmpty(ctx, src, dst, w)
	if err != nil {
		t.Fatalf("import2: %v", err)
	}
	if imported {
		t.Fatal("expected second import to be a no-op")
	}
}
