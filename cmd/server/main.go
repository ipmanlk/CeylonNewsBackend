package main

import (
	"context"
	"fmt"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper"
	"ipmanlk/cnapi/internal/service"
	"log/slog"
	"os"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create registry and service
	registry := scraper.NewRegistry(logger)
	scrapeService := service.NewScrapeService(registry, logger)

	ctx := context.Background()

	// Scrape all articles
	articles, err := scrapeService.ScrapeAll(ctx)
	if err != nil {
		logger.Error("failed to scrape articles", "error", err)
		os.Exit(1)
	}

	logger.Info("scraping completed", "total_articles", len(articles))

	// Print some basic information about scraped articles
	for _, article := range articles {
		fmt.Printf("Source: %s | Language: %s | Title: %s\n", 
			article.SourceName, 
			article.Language, 
			article.Title)
	}

	// Example: Scrape by specific source
	bbcArticles, err := scrapeService.ScrapeBySource(ctx, "BBC")
	if err != nil {
		logger.Error("failed to scrape BBC articles", "error", err)
	} else {
		logger.Info("BBC articles scraped", "count", len(bbcArticles))
	}

	// Example: Scrape by specific language
	englishArticles, err := scrapeService.ScrapeByLanguage(ctx, model.LangEn)
	if err != nil {
		logger.Error("failed to scrape English articles", "error", err)
	} else {
		logger.Info("English articles scraped", "count", len(englishArticles))
	}

	// Print available sources and languages
	fmt.Printf("\nAvailable sources: %v\n", scrapeService.GetAvailableSources())
	fmt.Printf("Available languages: %v\n", scrapeService.GetAvailableLanguages())
}
