package scraper

import (
	"context"
	"ipmanlk/cnapi/internal/model"
)

// SourceScraper defines the interface for all news source scrapers
type SourceScraper interface {
	Name() string
	Languages() []model.Language
	Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error)
}
