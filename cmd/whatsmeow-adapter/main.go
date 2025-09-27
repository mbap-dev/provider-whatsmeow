package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	amqpconsumer "your.org/provider-whatsmeow/internal/amqp"
	"your.org/provider-whatsmeow/internal/config"
	httpserver "your.org/provider-whatsmeow/internal/http"
	ilog "your.org/provider-whatsmeow/internal/log"
	"your.org/provider-whatsmeow/internal/provider"
	"your.org/provider-whatsmeow/internal/status"
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

	// Init Redis status writer (no-op if REDIS_URL empty)
	status.Init(cfg.RedisURL)

	// Create a manager for WhatsApp clients.  The manager is responsible
	// for bootstrapping new sessions, keeping track of connected clients
	// and exposing helper functions used by the HTTP handlers and the
	// message consumer.
	clientManager := provider.NewClientManager(cfg.SessionStore, cfg.WebhookBase, cfg.AudioPTTDefault, cfg.RejectCalls, cfg.RejectCallsMessage, cfg.AutoMarkReadOnMessage, cfg.PNResolverURL)

	// Auto-restore sessions that already have a saved store (session.db)
	if restored, err := clientManager.RestoreSavedSessions(); err != nil {
		ilog.Errorf("failed to restore saved sessions: %v", err)
	} else if len(restored) > 0 {
		ilog.Infof("restoring %d saved session(s): %v", len(restored), restored)
	} else {
		ilog.Infof("no saved sessions to restore")
	}

	// Ensure the exchange and durable queue exist so that publishers can
	// send messages even if the adapter is temporarily offline.
	if err := amqpconsumer.InitExchange(cfg); err != nil {
		ilog.Errorf("failed to initialize AMQP exchange: %v", err)
		os.Exit(1)
	}

	// Ensure the exchange and durable queue exist so that publishers can
	// send messages even if the adapter is temporarily offline.
	if err := amqpconsumer.InitExchange(cfg); err != nil {
		ilog.Errorf("failed to initialize AMQP exchange: %v", err)
		os.Exit(1)
	}

	// Initialize the AMQP consumer.  The consumer will connect to the
	// broker, bind the configured queue to the exchange and routing key and
	// then dispatch all incoming messages to the provider send function.
	consumer, err := amqpconsumer.NewConsumer(cfg, clientManager)
	if err != nil {
		ilog.Errorf("failed to initialise AMQP consumer: %v", err)
		os.Exit(1)
	}

	// Spin up the HTTP server exposing health checks, QR code
	// retrieval and basic session lifecycle management.  The consumer is
	// passed into the server so it can report readiness.
	srv := httpserver.NewServer(cfg, clientManager, consumer)

	// The root context is cancelled on SIGINT or SIGTERM which signals
	// all subordinate goroutines to stop.  Each component listens for
	// context cancellation and cleans up its resources accordingly.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the AMQP consumer in a separate goroutine.  If it returns
	// an error it will be logged.  The consumer blocks until the
	// context is cancelled.
	go func() {
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
