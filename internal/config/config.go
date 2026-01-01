package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Database  DatabaseConfig
	Fetcher   FetcherConfig
	Logger    LoggerConfig
	Scheduler SchedulerConfig
}

// SchedulerConfig holds scheduler-related configuration
type SchedulerConfig struct {
	ScrapeInterval time.Duration // How often to scrape all sources
	Enabled        bool          // Whether to enable automatic scraping
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Driver          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// FetcherConfig holds fetcher-related configuration
type FetcherConfig struct {
	HTTPTimeout     time.Duration
	BrowserAPIURL   string
	BrowserTimeout  time.Duration
	BrowserWaitTime int
}

// LoggerConfig holds logger-related configuration
type LoggerConfig struct {
	Level     string
	Format    string // "json" or "text"
	AddSource bool
}

// Load loads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "sqlite"),
			DSN:             getEnv("DB_DSN", "./data/db.sqlite"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Fetcher: FetcherConfig{
			HTTPTimeout:     getEnvDuration("FETCHER_HTTP_TIMEOUT", 15*time.Second),
			BrowserAPIURL:   getEnv("FETCHER_BROWSER_API_URL", "http://localhost:8000"),
			BrowserTimeout:  getEnvDuration("FETCHER_BROWSER_TIMEOUT", 60*time.Second),
			BrowserWaitTime: getEnvInt("FETCHER_BROWSER_WAIT_TIME", 15),
		},
		Logger: LoggerConfig{
			Level:     getEnv("LOG_LEVEL", "info"),
			Format:    getEnv("LOG_FORMAT", "json"),
			AddSource: getEnvBool("LOG_ADD_SOURCE", true),
		},
		Scheduler: SchedulerConfig{
			ScrapeInterval: getEnvDuration("SCHEDULER_SCRAPE_INTERVAL", 1*time.Hour),
			Enabled:        getEnvBool("SCHEDULER_ENABLED", true),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.DSN == "" {
		return fmt.Errorf("database DSN is required")
	}

	if c.Fetcher.BrowserAPIURL == "" {
		return fmt.Errorf("browser API URL is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logger.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logger.Level)
	}

	if c.Logger.Format != "json" && c.Logger.Format != "text" {
		return fmt.Errorf("invalid log format: %s (must be 'json' or 'text')", c.Logger.Format)
	}

	if c.Scheduler.Enabled && c.Scheduler.ScrapeInterval < 1*time.Minute {
		return fmt.Errorf("scrape interval must be at least 1 minute, got: %v", c.Scheduler.ScrapeInterval)
	}

	return nil
}

// Helper functions to read environment variables with defaults

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
