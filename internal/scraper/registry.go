package scraper

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper/source"
)

type SourceScraper interface {
	Name() string
	Languages() []model.Language
	Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error)
}

type Registry struct {
	scrapers []SourceScraper
}

func NewRegistry() *Registry {
	httpClient := fetcher.NewHTTPClient()
	browserClient := fetcher.NewBrowserAPIClient()
	fetcher := fetcher.NewFetcher(httpClient, browserClient)

	scrapers := []SourceScraper{
		source.NewBBCScraper(fetcher),
		source.NewDeranaScraper(fetcher),
		source.NewDailyMirrorScraper(fetcher),
		source.NewDivainaScraper(fetcher),
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
