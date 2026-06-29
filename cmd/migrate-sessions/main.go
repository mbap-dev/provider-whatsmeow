// migrate-sessions copies whatsmeow SQLite session databases into a shared
// Postgres store. Run it once per legacy environment before switching the
// adapter to Postgres mode (POSTGRES_URL set).
//
// Usage:
//
//	migrate-sessions --dir ./state/whatsmeow --postgres "$POSTGRES_URL"
//	migrate-sessions --sqlite ./state/whatsmeow/553199999999/session.db --postgres "$POSTGRES_URL"
//
// The tool is idempotent: rows that already exist in Postgres are silently
// skipped (ON CONFLICT DO NOTHING).
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	pgxv5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	"your.org/provider-whatsmeow/internal/provider"
)

// wmTables lists all whatsmeow tables in FK-safe order (parents before children).
// whatsmeow_version is intentionally excluded — Postgres tracks its own schema version.
var wmTables = []string{
	"whatsmeow_device",
	"whatsmeow_lid_map",
	"whatsmeow_identity_keys",
	"whatsmeow_pre_keys",
	"whatsmeow_sessions",
	"whatsmeow_sender_keys",
	"whatsmeow_app_state_sync_keys",
	"whatsmeow_app_state_version",
	"whatsmeow_app_state_mutation_macs",
	"whatsmeow_contacts",
	"whatsmeow_chat_settings",
	"whatsmeow_message_secrets",
	"whatsmeow_privacy_tokens",
	"whatsmeow_nct_salt",
	"whatsmeow_event_buffer",
	"whatsmeow_retry_buffer",
}

func main() {
	dirFlag := flag.String("dir", "", "directory containing per-session subdirectories (each with session.db) plus the _*.json metadata files")
	sqliteFlag := flag.String("sqlite", "", "path to a single session.db file")
	postgresFlag := flag.String("postgres", os.Getenv("POSTGRES_URL"), "PostgreSQL URL")
	workerIDFlag := flag.String("worker-id", "", "target worker_id for metadata (webhooks/templates/overrides/assignments). Empty = keep the worker_id stored in the JSON files. Set a fixed value to consolidate several containers under one worker.")
	dryRun := flag.Bool("dry-run", false, "print what would be migrated without writing to Postgres")
	flag.Parse()

	if strings.TrimSpace(*postgresFlag) == "" {
		log.Fatal("--postgres (or POSTGRES_URL env var) is required")
	}
	if *dirFlag == "" && *sqliteFlag == "" {
		log.Fatal("one of --dir or --sqlite is required")
	}

	ctx := context.Background()

	// Initialize Postgres schema via whatsmeow's own migration system.
	// This creates all whatsmeow_* tables if they don't exist yet.
	if !*dryRun {
		pgxCfg, err := pgxv5.ParseConfig(*postgresFlag)
		if err != nil {
			log.Fatalf("parse postgres url: %v", err)
		}
		sqlDB := stdlib.OpenDB(*pgxCfg)
		storeLog := waLog.Stdout("WAStore", "ERROR", false)
		c := sqlstore.NewWithDB(sqlDB, "pgx", storeLog)
		if err := c.Upgrade(ctx); err != nil {
			log.Fatalf("init postgres schema: %v", err)
		}
		c.Close()
		log.Println("postgres schema ready")
	}

	// Open Postgres pool for the actual data inserts.
	pgPool, err := pgxpool.New(ctx, *postgresFlag)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pgPool.Close()

	// Collect SQLite files.
	var sqliteFiles []string
	if *sqliteFlag != "" {
		sqliteFiles = []string{*sqliteFlag}
	} else {
		entries, err := os.ReadDir(*dirFlag)
		if err != nil {
			log.Fatalf("read dir %s: %v", *dirFlag, err)
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			p := filepath.Join(*dirFlag, e.Name(), "session.db")
			if _, err := os.Stat(p); err == nil {
				sqliteFiles = append(sqliteFiles, p)
			}
		}
	}

	if len(sqliteFiles) == 0 {
		log.Println("no session.db files found — nothing to migrate")
		return
	}

	log.Printf("found %d session(s) to migrate", len(sqliteFiles))

	totalDevices := 0
	for _, f := range sqliteFiles {
		n, err := migrateFile(ctx, f, pgPool, *dryRun)
		if err != nil {
			log.Printf("ERROR %s: %v", f, err)
		} else {
			totalDevices += n
		}
	}

	// Import the per-worker JSON metadata (webhooks, templates, overrides,
	// assignments, global webhooks) that lives alongside the session.db files.
	// migrate-sessions only handled the whatsmeow_* tables before; these JSON
	// files were never touched, which is why they appeared "missing" after a
	// SQLite→Postgres switch. We only know where they are in --dir mode.
	if metaDir := metadataDir(*dirFlag, *sqliteFlag); metaDir != "" {
		if err := importMetadata(ctx, metaDir, pgPool, *workerIDFlag, *dryRun); err != nil {
			log.Printf("ERROR metadata import (%s): %v", metaDir, err)
		}
	} else if *sqliteFlag != "" {
		log.Printf("note: metadata (webhooks/templates/...) import skipped in --sqlite mode; use --dir pointing at the SESSION_STORE root to include it")
	}

	fmt.Printf("\ndone — total devices migrated: %d\n", totalDevices)
}

// metadataDir returns the directory that holds the _*.json metadata files.
// In --dir mode that is the directory itself. In --sqlite mode the layout is
// <root>/<number>/session.db, so the root is two levels up — but that is only a
// best-effort guess, so we return "" and let the caller warn instead.
func metadataDir(dirFlag, sqliteFlag string) string {
	if strings.TrimSpace(dirFlag) != "" {
		return dirFlag
	}
	return ""
}

// importMetadata merges the per-worker JSON metadata into Postgres. It loads the
// file state ignoring the stored worker_id (so a mismatched WORKER_ID can never
// silently drop data), then merges it with whatever already exists in Postgres
// for the target worker and writes the union back. This makes it safe to run for
// several containers in a row, consolidating them under a single --worker-id.
func importMetadata(ctx context.Context, dir string, pool *pgxpool.Pool, targetWorkerID string, dryRun bool) error {
	src := provider.NewFileStore(dir)
	detected := detectWorkerID(dir)
	if strings.TrimSpace(targetWorkerID) == "" {
		targetWorkerID = detected
	}

	// Load file state with workerID="" to bypass the stored-worker_id filter.
	fileAssign, err := src.LoadAssignments(ctx, "")
	if err != nil {
		return fmt.Errorf("load assignments: %w", err)
	}
	fileOverrides, err := src.LoadOverrides(ctx, "")
	if err != nil {
		return fmt.Errorf("load overrides: %w", err)
	}
	fileWebhooks, err := src.LoadWebhooks(ctx, "")
	if err != nil {
		return fmt.Errorf("load webhooks: %w", err)
	}
	fileGlobal, err := src.LoadGlobalWebhooks(ctx, "")
	if err != nil {
		return fmt.Errorf("load global webhooks: %w", err)
	}
	fileTemplates, err := src.LoadTemplates(ctx, "")
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	fmt.Printf("\nmetadata @ %s (file worker_id=%q → postgres worker_id=%q)\n", dir, detected, targetWorkerID)
	fmt.Printf("  assignments=%d overrides=%d webhooks=%d global_webhooks=%d templates=%d\n",
		len(fileAssign), len(fileOverrides), len(fileWebhooks), len(fileGlobal), len(fileTemplates))

	if dryRun {
		fmt.Printf("  [dry-run] skipping metadata write\n")
		return nil
	}

	dst := provider.NewPostgresStore(pool)

	// Assignments: union of existing + file (dedup).
	curAssign, err := dst.LoadAssignments(ctx, targetWorkerID)
	if err != nil {
		return fmt.Errorf("load pg assignments: %w", err)
	}
	mergedAssign := unionStrings(curAssign, fileAssign)
	if len(mergedAssign) > 0 {
		if err := dst.SaveAssignments(ctx, targetWorkerID, mergedAssign); err != nil {
			return fmt.Errorf("save assignments: %w", err)
		}
	}

	// Overrides: combine maps, file wins on per-session conflict.
	curOverrides, err := dst.LoadOverrides(ctx, targetWorkerID)
	if err != nil {
		return fmt.Errorf("load pg overrides: %w", err)
	}
	for sid, v := range fileOverrides {
		curOverrides[sid] = v
	}
	if len(curOverrides) > 0 {
		if err := dst.SaveOverrides(ctx, targetWorkerID, curOverrides); err != nil {
			return fmt.Errorf("save overrides: %w", err)
		}
	}

	// Webhooks: combine maps, file wins on per-session conflict.
	curWebhooks, err := dst.LoadWebhooks(ctx, targetWorkerID)
	if err != nil {
		return fmt.Errorf("load pg webhooks: %w", err)
	}
	for sid, v := range fileWebhooks {
		curWebhooks[sid] = v
	}
	if len(curWebhooks) > 0 {
		if err := dst.SaveWebhooks(ctx, targetWorkerID, curWebhooks); err != nil {
			return fmt.Errorf("save webhooks: %w", err)
		}
	}

	// Global webhooks: union by webhook ID.
	curGlobal, err := dst.LoadGlobalWebhooks(ctx, targetWorkerID)
	if err != nil {
		return fmt.Errorf("load pg global webhooks: %w", err)
	}
	mergedGlobal := mergeGlobalWebhooks(curGlobal, fileGlobal)
	if len(mergedGlobal) > 0 {
		if err := dst.SaveGlobalWebhooks(ctx, targetWorkerID, mergedGlobal); err != nil {
			return fmt.Errorf("save global webhooks: %w", err)
		}
	}

	// Templates: combine maps, file wins on per-session conflict.
	curTemplates, err := dst.LoadTemplates(ctx, targetWorkerID)
	if err != nil {
		return fmt.Errorf("load pg templates: %w", err)
	}
	for sid, v := range fileTemplates {
		curTemplates[sid] = v
	}
	if len(curTemplates) > 0 {
		if err := dst.SaveTemplates(ctx, targetWorkerID, curTemplates); err != nil {
			return fmt.Errorf("save templates: %w", err)
		}
	}

	fmt.Printf("  imported into worker_id=%q (assignments=%d overrides=%d webhooks=%d global_webhooks=%d templates=%d)\n",
		targetWorkerID, len(mergedAssign), len(curOverrides), len(curWebhooks), len(mergedGlobal), len(curTemplates))
	return nil
}

// detectWorkerID reads the worker_id stored at the top of the metadata JSON
// files. All five share the same `worker_id` field; the first non-empty wins.
func detectWorkerID(dir string) string {
	files := []string{
		"_worker_assignments.json",
		"_session_webhooks.json",
		"_session_overrides.json",
		"_session_templates.json",
		"_global_webhooks.json",
	}
	for _, name := range files {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		var probe struct {
			WorkerID string `json:"worker_id"`
		}
		if err := json.Unmarshal(b, &probe); err != nil {
			continue
		}
		if strings.TrimSpace(probe.WorkerID) != "" {
			return strings.TrimSpace(probe.WorkerID)
		}
	}
	return ""
}

// unionStrings returns the deduplicated union of two string slices, preserving
// order (existing first, then new entries).
func unionStrings(existing, incoming []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(existing)+len(incoming))
	for _, list := range [][]string{existing, incoming} {
		for _, s := range list {
			s = strings.TrimSpace(s)
			if s == "" || seen[s] {
				continue
			}
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// mergeGlobalWebhooks returns the union of two global-webhook slices, deduped by
// webhook ID (existing entries win).
func mergeGlobalWebhooks(existing, incoming []provider.SessionWebhook) []provider.SessionWebhook {
	seen := make(map[string]bool)
	out := make([]provider.SessionWebhook, 0, len(existing)+len(incoming))
	for _, list := range [][]provider.SessionWebhook{existing, incoming} {
		for _, w := range list {
			if w.ID != "" {
				if seen[w.ID] {
					continue
				}
				seen[w.ID] = true
			}
			out = append(out, w)
		}
	}
	return out
}

func migrateFile(ctx context.Context, sqlitePath string, pg *pgxpool.Pool, dryRun bool) (int, error) {
	dsn := fmt.Sprintf(
		"file:%s?cache=shared&_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)",
		sqlitePath,
	)
	sqliteDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return 0, fmt.Errorf("open sqlite %s: %w", sqlitePath, err)
	}
	defer sqliteDB.Close()

	// Count devices so we can report progress.
	var deviceCount int
	if err := sqliteDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM whatsmeow_device").Scan(&deviceCount); err != nil {
		return 0, fmt.Errorf("count devices: %w", err)
	}

	fmt.Printf("\n[%s] %d device(s)\n", sqlitePath, deviceCount)

	if dryRun {
		fmt.Printf("  [dry-run] skipping actual migration\n")
		return deviceCount, nil
	}

	for _, table := range wmTables {
		n, err := copyTable(ctx, sqliteDB, pg, table)
		if err != nil {
			log.Printf("  WARN table %s: %v", table, err)
		} else if n > 0 {
			fmt.Printf("  %-45s %d rows\n", table, n)
		}
	}
	return deviceCount, nil
}

type colInfo struct {
	name    string
	typName string
}

func copyTable(ctx context.Context, src *sql.DB, dst *pgxpool.Pool, table string) (int, error) {
	// Get column metadata from SQLite.
	pragmaRows, err := src.QueryContext(ctx, "PRAGMA table_info("+table+")")
	if err != nil {
		return 0, fmt.Errorf("pragma: %w", err)
	}
	var cols []colInfo
	for pragmaRows.Next() {
		var cid, notnull, pk int
		var name, typName string
		var dflt sql.NullString
		if err := pragmaRows.Scan(&cid, &name, &typName, &notnull, &dflt, &pk); err != nil {
			continue
		}
		cols = append(cols, colInfo{name, strings.ToUpper(typName)})
	}
	pragmaRows.Close()

	if len(cols) == 0 {
		// Table doesn't exist in this (older-schema) SQLite — skip silently.
		return 0, nil
	}

	names := make([]string, len(cols))
	placeholders := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.name
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING",
		table,
		strings.Join(names, ", "),
		strings.Join(placeholders, ", "),
	)

	dataRows, err := src.QueryContext(ctx,
		"SELECT "+strings.Join(names, ", ")+" FROM "+table,
	)
	if err != nil {
		return 0, fmt.Errorf("select: %w", err)
	}
	defer dataRows.Close()

	raw := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range raw {
		ptrs[i] = &raw[i]
	}

	inserted := 0
	for dataRows.Next() {
		if err := dataRows.Scan(ptrs...); err != nil {
			log.Printf("    WARN scan in %s: %v", table, err)
			continue
		}

		args := make([]interface{}, len(cols))
		for i, col := range cols {
			v := raw[i]
			// SQLite stores BOOLEAN as INTEGER (0/1); Postgres needs bool.
			if col.typName == "BOOLEAN" {
				if vi, ok := v.(int64); ok {
					v = vi != 0
				}
			}
			args[i] = v
		}

		if _, err := dst.Exec(ctx, insertSQL, args...); err != nil {
			log.Printf("    WARN insert in %s: %v", table, err)
			continue
		}
		inserted++
	}
	return inserted, nil
}
