package scraper

import (
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/scraper/source"
	"log/slog"
)

// Registry manages all available scrapers
type Registry struct {
	scrapers []SourceScraper
}

// NewRegistry creates a new scraper registry with all available scrapers
func NewRegistry(logger *slog.Logger) *Registry {
	// Create shared fetchers
	httpClient := fetcher.NewHTTPClient()
	htmlProcessor := fetcher.NewHTMLProcessor()
	contentExtractor := fetcher.NewContentExtractor(httpClient)

	// Create RSS fetcher with extractor fallback
	rssFetcher := fetcher.NewRSSFetcher(httpClient, htmlProcessor, logger)
	rssFetcher.SetContentExtractor(contentExtractor)

	// Create scrapers
	scrapers := []SourceScraper{
		source.NewBBCScraper(httpClient, htmlProcessor, contentExtractor, logger),
		source.NewDivainaScraper(rssFetcher, logger),
	}

	return &Registry{
		scrapers: scrapers,
	}
}

// GetScrapers returns all registered scrapers
func (r *Registry) GetScrapers() []SourceScraper {
	return r.scrapers
}

// GetScraperByName returns a specific scraper by name
func (r *Registry) GetScraperByName(name string) SourceScraper {
	for _, scraper := range r.scrapers {
		if scraper.Name() == name {
			return scraper
		}
	}
	return nil
}

// GetScrapersByLanguage returns scrapers that support a specific language
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
