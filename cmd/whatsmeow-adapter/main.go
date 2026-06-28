package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	amqpconsumer "your.org/provider-whatsmeow/internal/amqp"
	broker "your.org/provider-whatsmeow/internal/broker"
	"your.org/provider-whatsmeow/internal/auth"
	"your.org/provider-whatsmeow/internal/config"
	httpserver "your.org/provider-whatsmeow/internal/http"
	ilog "your.org/provider-whatsmeow/internal/log"
	"your.org/provider-whatsmeow/internal/provider"
	"your.org/provider-whatsmeow/internal/status"
	"your.org/provider-whatsmeow/internal/store/pg"
)

// main is the entrypoint for the WhatsMeow adapter.  It wires together the
// configuration loader, the WhatsApp client manager, the AMQP consumer and
// the HTTP API server.  All long‑running components are started
// concurrently and the application will shut down gracefully when an
// interrupt signal (SIGINT or SIGTERM) is received.
func main() {
	// Load configuration from environment variables.  If required values are
	// missing a sensible default is used instead.  See config.NewConfig
	// for details on each field.
	cfg := config.NewConfig()
	provider.SetDefaultProxyConfig(provider.ProxyConfig{
		Protocol: cfg.ProxyProtocol,
		Host:     cfg.ProxyHost,
		Port:     cfg.ProxyPort,
		Username: cfg.ProxyUsername,
		Password: cfg.ProxyPassword,
	})

	// Init status writer (in-memory + optional Redis mirror)
	status.Init(cfg.RedisURL, cfg.StatusPrefix)

	// Optional Postgres-backed metadata store. When POSTGRES_URL is set we
	// connect, apply embedded migrations and import any existing JSON state on
	// first boot. When unset, the manager keeps using the JSON files.
	var metadataStore provider.MetadataStore
	var authHandler *auth.Handler
	var pgPool *pgxpool.Pool
	if cfg.PostgresURL != "" {
		var err error
		pgPool, err = connectPostgresWithRetry(cfg.PostgresURL)
		if err != nil {
			ilog.Errorf("evt=postgres_connect failed err=%v", err)
			os.Exit(1)
		}
		defer pgPool.Close()
		pool := pgPool
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		metadataStore = provider.NewPostgresStore(pool)
		if imported, err := provider.ImportFileStateIfEmpty(ctx, provider.NewFileStore(cfg.SessionStore), metadataStore, cfg.WorkerID); err != nil {
			ilog.Errorf("evt=postgres_import_json failed err=%v", err)
		} else if imported {
			ilog.Infof("evt=postgres_import_json success worker_id=%s", cfg.WorkerID)
		}
		// Wire up Postgres-backed auth.
		authStore := auth.NewStore(pool)
		if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
			if created, err := authStore.BootstrapOwner(ctx, cfg.AdminEmail, cfg.AdminPassword); err != nil {
				ilog.Warnf("evt=auth_bootstrap err=%v", err)
			} else if created {
				ilog.Infof("evt=auth_bootstrap created owner email=%s", cfg.AdminEmail)
			}
		}
		authHandler = auth.NewHandler(authStore, cfg.AuthSecureCookie, time.Duration(cfg.SessionTTLHours)*time.Hour)
		// Start background session cleanup — runs every hour, independent of request lifecycle.
		go func(store *auth.Store) {
			ticker := time.NewTicker(time.Hour)
			defer ticker.Stop()
			for range ticker.C {
				if err := store.PurgeExpiredSessions(context.Background()); err != nil {
					ilog.Warnf("evt=auth_session_purge err=%v", err)
				}
			}
		}(authStore)
		cancel()
		ilog.Infof("evt=metadata_store backend=postgres")
	}

	// Create a manager for WhatsApp clients.  The manager is responsible
	// for bootstrapping new sessions, keeping track of connected clients
	// and exposing helper functions used by the HTTP handlers and the
	// message consumer.
	clientManager := provider.NewClientManager(provider.ManagerOptions{
		MetadataStore:         metadataStore,
		PGPool:                pgPool,
		PostgresURL:           cfg.PostgresURL,
		SessionStore:          cfg.SessionStore,
		WebhookBase:           cfg.WebhookBase,
		DefaultAudioPTT:       cfg.AudioPTTDefault,
		RejectCalls:           cfg.RejectCalls,
		RejectMessage:         cfg.RejectCallsMessage,
		AutoMarkRead:          cfg.AutoMarkReadOnMessage,
		PNResolverURL:         cfg.PNResolverURL,
		ProxyProtocol:         cfg.ProxyProtocol,
		ProxyHost:             cfg.ProxyHost,
		ProxyPort:             cfg.ProxyPort,
		ProxyUsername:         cfg.ProxyUsername,
		ProxyPassword:         cfg.ProxyPassword,
		IgnoreStatusBroadcast: cfg.IgnoreStatusBroadcast,
		IgnoreNewsletters:     cfg.IgnoreNewsletters,
		OutgoingMessagesCoex:  cfg.OutgoingMessagesCoex,
		RedisURL:              cfg.RedisURL,
		SessionConfigPrefix:   cfg.SessionConfigPrefix,
		CheckUserExists:       cfg.CheckUserExists,
		WorkerID:              cfg.WorkerID,
		AssignedSessions:      cfg.WorkerNumbers,
		WorkerStrict:          cfg.WorkerStrict,
		MediaStoreDir:         cfg.MediaStoreDir,
		MediaRegistryFile:     cfg.MediaRegistryFile,
		MediaStorageMode:      cfg.MediaStorageMode,
		MediaURLTTLSeconds:    cfg.MediaURLTTLSeconds,
		S3Bucket:              cfg.S3Bucket,
		S3Region:              cfg.S3Region,
		S3Endpoint:            cfg.S3Endpoint,
		S3AccessKeyID:         cfg.S3AccessKeyID,
		S3SecretAccessKey:     cfg.S3SecretAccess,
		S3ForcePathStyle:      cfg.S3ForcePathStyle,
		S3Prefix:              cfg.S3Prefix,
	})
	if cfg.AlwaysOnline {
		clientManager.EnableAlwaysOnline(time.Duration(cfg.AlwaysOnlineIntervalSeconds) * time.Second)
	}

	// Auto-restore sessions that already have a saved store (session.db)
	if restored, err := clientManager.RestoreSavedSessions(); err != nil {
		ilog.Errorf("failed to restore saved sessions: %v", err)
	} else if len(restored) > 0 {
		ilog.Infof("restoring %d saved session(s): %v", len(restored), restored)
	} else {
		ilog.Infof("no saved sessions to restore")
	}
	if cfg.WorkerAutoConnectAssigned && len(cfg.WorkerNumbers) > 0 {
		auto := clientManager.AutoConnectAssignedSessions()
		if len(auto) > 0 {
			ilog.Infof("worker auto-connect enabled: %d assigned session(s) scheduled: %v", len(auto), auto)
		}
	}

	// Ensure the exchange and durable queue exist so that publishers can
	// send messages even if the adapter is temporarily offline.
	if err := initAMQPExchangeWithRetry(cfg); err != nil {
		ilog.Errorf("failed to initialize AMQP exchange: %v", err)
		if cfg.AMQPConsumerEnabled {
			os.Exit(1)
		}
	}

	// Initialize AMQP webhook publisher (optional sink)
	if cfg.AMQPURL != "" {
		if err := broker.InitWebhookPublisher(cfg); err != nil {
			ilog.Errorf("failed to initialize AMQP webhook publisher: %v", err)
			// Do not exit: keep running; HTTP fallback may be used.
		}
	} else {
		ilog.Infof("AMQP webhook publisher skipped (AMQP_URL empty)")
	}

	// Initialize the AMQP consumer.  The consumer will connect to the
	// broker, bind the configured queue to the exchange and routing key and
	// then dispatch all incoming messages to the provider send function.
	var consumer *amqpconsumer.Consumer
	var err error
	if cfg.AMQPConsumerEnabled {
		consumer, err = amqpconsumer.NewConsumer(cfg, clientManager)
		if err != nil {
			ilog.Errorf("failed to initialise AMQP consumer: %v", err)
			os.Exit(1)
		}
	} else {
		ilog.Infof("running in standalone mode (AMQP consumer disabled)")
	}

	// Spin up the HTTP server exposing health checks, QR code
	// retrieval and basic session lifecycle management.  The consumer is
	// passed into the server so it can report readiness.
	srv := httpserver.NewServer(cfg, clientManager, consumer, authHandler)

	// The root context is cancelled on SIGINT or SIGTERM which signals
	// all subordinate goroutines to stop.  Each component listens for
	// context cancellation and cleans up its resources accordingly.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the AMQP consumer in a separate goroutine.  If it returns
	// an error it will be logged.  The consumer blocks until the
	// context is cancelled.
	go func() {
		if consumer == nil {
			return
		}
		if err := consumer.Start(ctx); err != nil {
			ilog.Errorf("AMQP consumer stopped: %v", err)
		}
	}()

	// Start the HTTP API server.  ListenAndServe blocks so it is
	// executed in its own goroutine.  If the server returns an error
	// before the context is cancelled it is logged.
	go func() {
		if err := srv.Start(); err != nil {
			ilog.Errorf("HTTP server stopped: %v", err)
		}
	}()

	// Wait for a termination signal and initiate a graceful shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	ilog.Infof("Shutting down…")

	// Cancel the root context which causes the consumer to exit.
	cancel()

	// Gracefully shut down the HTTP server.  A fresh context with a
	// timeout could be provided here to bound the shutdown period.
	if err := srv.Shutdown(context.Background()); err != nil {
		ilog.Errorf("failed to shutdown HTTP server: %v", err)
	}
}

// connectPostgresWithRetry retries pg.Migrate + pg.Connect for up to 60 s so
// the adapter survives Docker Swarm starting it before Postgres is ready.
func connectPostgresWithRetry(dsn string) (*pgxpool.Pool, error) {
	const maxAttempts = 30
	const delay = 2 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		migrateErr := pg.Migrate(ctx, dsn)
		cancel()
		if migrateErr != nil {
			lastErr = migrateErr
			ilog.Warnf("evt=postgres_not_ready attempt=%d/%d err=%v", attempt, maxAttempts, migrateErr)
			time.Sleep(delay)
			continue
		}
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		pool, connectErr := pg.Connect(ctx, dsn)
		cancel()
		if connectErr != nil {
			lastErr = connectErr
			ilog.Warnf("evt=postgres_connect_failed attempt=%d/%d err=%v", attempt, maxAttempts, connectErr)
			time.Sleep(delay)
			continue
		}
		if attempt > 1 {
			ilog.Infof("evt=postgres_ready attempts=%d", attempt)
		}
		return pool, nil
	}
	return nil, lastErr
}

func initAMQPExchangeWithRetry(cfg *config.Config) error {
	const attempts = 30
	const delay = 2 * time.Second

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := amqpconsumer.InitExchange(cfg); err != nil {
			lastErr = err
			if cfg == nil || !cfg.AMQPConsumerEnabled {
				return err
			}
			ilog.Warnf("AMQP exchange initialization failed attempt=%d/%d err=%v", attempt, attempts, err)
			time.Sleep(delay)
			continue
		}
		if attempt > 1 {
			ilog.Infof("AMQP exchange initialized after %d attempt(s)", attempt)
		}
		return nil
	}
	return lastErr
}
