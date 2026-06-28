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
	dirFlag := flag.String("dir", "", "directory containing per-session subdirectories (each with session.db)")
	sqliteFlag := flag.String("sqlite", "", "path to a single session.db file")
	postgresFlag := flag.String("postgres", os.Getenv("POSTGRES_URL"), "PostgreSQL URL")
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
	fmt.Printf("\ndone — total devices migrated: %d\n", totalDevices)
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
