package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"ipmanlk/cnapi/internal/api"
	"ipmanlk/cnapi/internal/config"
	"ipmanlk/cnapi/internal/database"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/scheduler"
	"ipmanlk/cnapi/internal/scraper"
	"ipmanlk/cnapi/internal/service"
)

// App holds all application dependencies and services
type App struct {
	Config     *config.Config
	DB         *sql.DB
	Store      *database.Store
	Services   *Services
	Scheduler  *scheduler.Scheduler
	HTTPServer *api.Server
	Logger     *slog.Logger
}

// Services holds all application services
type Services struct {
	Scrape  service.ScrapeService
	Article service.ArticleService
	Search  service.SearchService
}

// New creates and initializes a new application instance
func New(ctx context.Context) (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup logger
	logger := setupLogger(cfg.Logger)
	slog.SetDefault(logger)

	slog.Info("starting application initialization")

	// Initialize database
	db, err := initDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run migrations
	if err := database.InitializeDatabase(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Create store
	store := database.NewStore(db)

	// Initialize services
	services := initServices(cfg, store)

	// Initialize scheduler
	sched := scheduler.New(
		services.Scrape,
		services.Article,
		cfg.Scheduler.ScrapeInterval,
	)

	httpConfig := api.Config{
		Host:            cfg.HTTP.Host,
		Port:            cfg.HTTP.Port,
		ReadTimeout:     cfg.HTTP.ReadTimeout,
		WriteTimeout:    cfg.HTTP.WriteTimeout,
		IdleTimeout:     cfg.HTTP.IdleTimeout,
		ShutdownTimeout: cfg.HTTP.ShutdownTimeout,
	}
	httpServer := api.NewServer(services.Article, services.Search, httpConfig)

	app := &App{
		Config:     cfg,
		DB:         db,
		Store:      store,
		Services:   services,
		Scheduler:  sched,
		HTTPServer: httpServer,
		Logger:     logger,
	}

	slog.Info("application initialized successfully")

	return app, nil
}

// Close gracefully shuts down the application
func (a *App) Close(ctx context.Context) error {
	slog.Info("shutting down application")

	// Stop HTTP server if running
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Shutdown(ctx); err != nil {
			slog.Error("error stopping HTTP server", "error", err)
		}
	}

	// Stop scheduler if running
	if a.Scheduler != nil && a.Scheduler.IsRunning() {
		if err := a.Scheduler.Stop(); err != nil {
			slog.Error("error stopping scheduler", "error", err)
		}
	}

	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}

	slog.Info("application shutdown complete")
	return nil
}

// setupLogger creates and configures a structured logger
func setupLogger(cfg config.LoggerConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// initDatabase initializes the database connection
func initDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("database connection established",
		"driver", cfg.Driver,
		"max_open_conns", cfg.MaxOpenConns,
		"max_idle_conns", cfg.MaxIdleConns,
	)

	return db, nil
}

// initServices initializes all application services with their dependencies
func initServices(cfg *config.Config, store *database.Store) *Services {
	// Initialize fetcher dependencies
	httpClient := fetcher.NewHTTPClient(cfg.Fetcher.HTTPTimeout)
	browserClient := fetcher.NewBrowserAPIClient(
		cfg.Fetcher.BrowserAPIURL,
		cfg.Fetcher.BrowserTimeout,
		cfg.Fetcher.BrowserWaitTime,
	)
	fetch := fetcher.NewFetcher(httpClient, browserClient)

	// Initialize scraper registry
	scraperRegistry := scraper.NewRegistry(fetch)

	// Initialize services
	return &Services{
		Scrape:  service.NewScrapeService(scraperRegistry),
		Article: service.NewArticleService(store.Articles),
		Search:  service.NewSearchService(store.Search),
	}
}
