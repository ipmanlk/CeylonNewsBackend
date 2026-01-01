package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ipmanlk/cnapi/internal/app"
)

func main() {
	// Create context that listens for termination signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize application
	application, err := app.New(ctx)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := application.Close(); err != nil {
			slog.Error("error during application shutdown", "error", err)
		}
	}()

	slog.Info("Ceylon News Backend started successfully")

	// Start scheduler if enabled
	if application.Config.Scheduler.Enabled {
		if err := application.Scheduler.Start(ctx); err != nil {
			slog.Error("failed to start scheduler", "error", err)
			os.Exit(1)
		}
		slog.Info("periodic scraping enabled", "interval", application.Config.Scheduler.ScrapeInterval)
	} else {
		slog.Info("periodic scraping disabled")
	}

	// Wait for termination signal
	<-ctx.Done()
	slog.Info("received shutdown signal")
}
