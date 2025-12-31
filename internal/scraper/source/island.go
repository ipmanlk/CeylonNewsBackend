package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
)

type IslandScraper struct {
	fetcher *fetcher.Fetcher
}

func NewIslandScraper(fetcher *fetcher.Fetcher) *IslandScraper {
	return &IslandScraper{
		fetcher: fetcher,
	}
}

func (s *IslandScraper) Name() string {
	return "The Island"
}

func (s *IslandScraper) Languages() []model.Language {
	return []model.Language{model.LangEn}
}

func (s *IslandScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangEn:
		return s.scrapeEn(ctx)
	default:
		return nil, nil
	}
}

func (s *IslandScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	// Fetch RSS feed
	rssItems, err := s.fetcher.FetchRSS(ctx, "https://island.lk/feed/", 5)
	if err != nil {
		return nil, err
	}

	articles := make([]model.ScrapedArticle, 0, len(rssItems))
	seenLinks := make(map[string]bool)

	for _, item := range rssItems {
		if seenLinks[item.Link] {
			continue
		}
		seenLinks[item.Link] = true

		// Extract article content using the RSS item
		article, err := s.fetcher.ExtractArticleFromRSSItem(ctx, item)
		if err != nil {
			slog.Warn("failed to extract article", "scraper", "The Island", "url", item.Link, "error", err)
			continue
		}

		if article.Title == "" || article.Content == "" {
			slog.Debug("skipping article with missing content", "scraper", "The Island", "url", item.Link)
			continue
		}

		// Update the article with our source-specific information
		article.SourceName = s.Name()
		article.Language = model.LangEn

		articles = append(articles, article)
	}

	return articles, nil
}
