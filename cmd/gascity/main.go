// Package main is the entry point for the gascity application.
// gascity is a fork of gastownhall/gascity, providing gas price monitoring
// and analysis tooling for EVM-compatible networks.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gascity/gascity/internal/config"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"
	// Commit is the git commit hash set at build time.
	Commit = "none"
	// BuildDate is the build timestamp set at build time.
	BuildDate = "unknown"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // changed from LevelInfo to see more detail during local dev
	}))
	slog.SetDefault(logger)

	slog.Info("starting gascity",
		"version", Version,
		"commit", Commit,
		"build_date", BuildDate,
	)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	if err := run(ctx, cfg); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}

	slog.Info("gascity shut down cleanly")
}

// run initialises and starts all application components, blocking until
// the provided context is cancelled or a fatal error occurs.
func run(ctx context.Context, cfg *config.Config) error {
	slog.Info("configuration loaded",
		"rpc_url", cfg.RPCURL,
		"poll_interval", cfg.PollInterval,
		"listen_addr", cfg.ListenAddr,
	)

	// TODO: wire up collector, server, and storage components.
	<-ctx.Done()

	slog.Info("shutdown signal received")

	// NOTE: context.Canceled on clean shutdown is expected; not really an
	// error worth logging as one. Keeping the return for now to satisfy the
	// error-check in main, but this should be revisited when run() grows.
	return fmt.Errorf("context cancelled: %w", ctx.Err())
}
