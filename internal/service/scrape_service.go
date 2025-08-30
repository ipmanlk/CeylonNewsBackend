package service

import (
	"context"
	"fmt"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper"
	"log/slog"
)

// ScrapeService orchestrates scraping operations across all available scrapers
type ScrapeService struct {
	registry *scraper.Registry
	logger   *slog.Logger
}

// NewScrapeService creates a new scraping service
func NewScrapeService(registry *scraper.Registry, logger *slog.Logger) *ScrapeService {
	return &ScrapeService{
		registry: registry,
		logger:   logger,
	}
}

// ScrapeAll scrapes articles from all available scrapers and languages
func (s *ScrapeService) ScrapeAll(ctx context.Context) ([]model.ScrapedArticle, error) {
	var allArticles []model.ScrapedArticle
	scrapers := s.registry.GetScrapers()

	for _, scraper := range scrapers {
		for _, lang := range scraper.Languages() {
			articles, err := scraper.Scrape(ctx, lang)
			if err != nil {
				s.logger.Warn("failed to scrape from source",
					"source", scraper.Name(),
					"language", lang,
					"error", err,
				)
				continue
			}

			s.logger.Info("successfully scraped articles",
				"source", scraper.Name(),
				"language", lang,
				"count", len(articles),
			)

			allArticles = append(allArticles, articles...)
		}
	}

	return allArticles, nil
}

// ScrapeBySource scrapes articles from a specific source
func (s *ScrapeService) ScrapeBySource(ctx context.Context, sourceName string) ([]model.ScrapedArticle, error) {
	scraper := s.registry.GetScraperByName(sourceName)
	if scraper == nil {
		return nil, fmt.Errorf("scraper not found: %s", sourceName)
	}

	var allArticles []model.ScrapedArticle
	for _, lang := range scraper.Languages() {
		articles, err := scraper.Scrape(ctx, lang)
		if err != nil {
			s.logger.Warn("failed to scrape from source",
				"source", sourceName,
				"language", lang,
				"error", err,
			)
			continue
		}

		allArticles = append(allArticles, articles...)
	}

	return allArticles, nil
}

// ScrapeByLanguage scrapes articles from all sources for a specific language
func (s *ScrapeService) ScrapeByLanguage(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	var allArticles []model.ScrapedArticle
	scrapers := s.registry.GetScrapersByLanguage(string(language))

	for _, scraper := range scrapers {
		articles, err := scraper.Scrape(ctx, language)
		if err != nil {
			s.logger.Warn("failed to scrape from source",
				"source", scraper.Name(),
				"language", language,
				"error", err,
			)
			continue
		}

		allArticles = append(allArticles, articles...)
	}

	return allArticles, nil
}

// GetAvailableSources returns a list of all available source names
func (s *ScrapeService) GetAvailableSources() []string {
	scrapers := s.registry.GetScrapers()
	sources := make([]string, len(scrapers))
	for i, scraper := range scrapers {
		sources[i] = scraper.Name()
	}
	return sources
}

// GetAvailableLanguages returns a list of all available languages
func (s *ScrapeService) GetAvailableLanguages() []model.Language {
	scrapers := s.registry.GetScrapers()
	languageMap := make(map[model.Language]bool)
	
	for _, scraper := range scrapers {
		for _, lang := range scraper.Languages() {
			languageMap[lang] = true
		}
	}

	languages := make([]model.Language, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}
	
	return languages
}
