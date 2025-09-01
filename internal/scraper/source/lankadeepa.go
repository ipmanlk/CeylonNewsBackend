package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
)

type LankadeepaScraper struct {
	fetcher *fetcher.Fetcher
}

func NewLankadeepaScraper(fetcher *fetcher.Fetcher) *LankadeepaScraper {
	return &LankadeepaScraper{
		fetcher: fetcher,
	}
}

func (s *LankadeepaScraper) Name() string {
	return "Lankadeepa"
}

func (s *LankadeepaScraper) Languages() []model.Language {
	return []model.Language{model.LangSi}
}

func (s *LankadeepaScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangSi:
		return s.scrapeSi(ctx)
	default:
		return nil, nil
	}
}

func (s *LankadeepaScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	// Fetch RSS feed
	rssItems, err := s.fetcher.FetchRSS(ctx, "https://www.lankadeepa.lk/rss/latest_news/1", 5)
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
			slog.Warn("failed to extract article", "scraper", "Lankadeepa", "url", item.Link, "error", err)
			continue
		}

		if article.Title == "" || article.Content == "" {
			slog.Debug("skipping article with missing content", "scraper", "Lankadeepa", "url", item.Link)
			continue
		}

		// Update the article with our source-specific information
		article.SourceName = s.Name()
		article.Language = model.LangSi

		articles = append(articles, article)
	}

	return articles, nil
}
