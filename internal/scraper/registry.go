package scraper

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper/source"
	"log/slog"
)

type SourceScraper interface {
	Name() string
	Languages() []model.Language
	Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error)
}

type Registry struct {
	scrapers []SourceScraper
}

func NewRegistry(logger *slog.Logger) *Registry {
	// Create shared fetchers
	httpClient := fetcher.NewHTTPClient()
	htmlProcessor := fetcher.NewHTMLProcessor()
	contentExtractor := fetcher.NewContentExtractor(httpClient)

	// Create RSS fetcher with extractor fallback
	rssFetcher := fetcher.NewRSSFetcher(httpClient, htmlProcessor, contentExtractor, logger)

	// Create scrapers with browser API support
	browserClient := fetcher.NewBrowserAPIClient()
	htmlFetcher := fetcher.NewHTMLFetcher(browserClient)

	scrapers := []SourceScraper{
		source.NewBBCScraper(htmlFetcher, htmlProcessor, contentExtractor, logger),
		source.NewDivainaScraper(rssFetcher, logger),
	}

	return &Registry{
		scrapers: scrapers,
	}
}

func (r *Registry) GetScrapers() []SourceScraper {
	return r.scrapers
}

func (r *Registry) GetScraperByName(name string) SourceScraper {
	for _, scraper := range r.scrapers {
		if scraper.Name() == name {
			return scraper
		}
	}
	return nil
}

func (r *Registry) GetScrapersByLanguage(language string) []SourceScraper {
	var filtered []SourceScraper
	for _, scraper := range r.scrapers {
		for _, lang := range scraper.Languages() {
			if string(lang) == language {
				filtered = append(filtered, scraper)
				break
			}
		}
	}
	return filtered
}
