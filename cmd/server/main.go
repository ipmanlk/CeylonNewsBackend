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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	registry := scraper.NewRegistry()
	scrapeService := service.NewScrapeService(registry)

	ctx := context.Background()

	articles, err := scrapeService.ScrapeAll(ctx)
	if err != nil {
		slog.Error("failed to scrape articles", "error", err)
		os.Exit(1)
	}

	slog.Info("scraping completed", "total_articles", len(articles))

	for _, article := range articles {
		fmt.Printf("Source: %s | Language: %s | Title: %s\n",
			article.SourceName,
			article.Language,
			article.Title)
	}

	bbcArticles, err := scrapeService.ScrapeBySource(ctx, "BBC")
	if err != nil {
		slog.Error("failed to scrape BBC articles", "error", err)
	} else {
		slog.Info("BBC articles scraped", "count", len(bbcArticles))
	}

	englishArticles, err := scrapeService.ScrapeByLanguage(ctx, model.LangEn)
	if err != nil {
		slog.Error("failed to scrape English articles", "error", err)
	} else {
		slog.Info("English articles scraped", "count", len(englishArticles))
	}

	fmt.Printf("\nAvailable sources: %v\n", scrapeService.GetAvailableSources())
	fmt.Printf("Available languages: %v\n", scrapeService.GetAvailableLanguages())
}
