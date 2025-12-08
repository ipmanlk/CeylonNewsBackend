package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
	"strings"
)

type DailyMirrorScraper struct {
	fetcher *fetcher.Fetcher
}

func NewDailyMirrorScraper(fetcher *fetcher.Fetcher) *DailyMirrorScraper {
	return &DailyMirrorScraper{
		fetcher: fetcher,
	}
}

func (s *DailyMirrorScraper) Name() string {
	return "Daily Mirror"
}

func (s *DailyMirrorScraper) Languages() []model.Language {
	return []model.Language{model.LangEn}
}

func (s *DailyMirrorScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangEn:
		return s.scrapeEn(ctx)
	default:
		return nil, nil
	}
}

func (s *DailyMirrorScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	// Fetch RSS feed
	rssItems, err := s.fetcher.FetchRSS(ctx, "https://www.dailymirror.lk/rss/todays_headlines/419", 5)
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

		// Skip articles with error messages in title
		if strings.Contains(item.Title, "An Error Was") {
			slog.Debug("skipping article with error in title", "scraper", "Daily Mirror", "url", item.Link)
			continue
		}

		// Clean up title by removing the suffix
		title := strings.ReplaceAll(item.Title, " - Breaking News | Daily Mirror", "")

		// Extract article content using the RSS item
		article, err := s.fetcher.ExtractArticleFromRSSItem(ctx, item, true)
		if err != nil {
			slog.Warn("failed to extract article", "scraper", "Daily Mirror", "url", item.Link, "error", err)
			continue
		}

		// Update the article with our source-specific information
		article.SourceName = s.Name()
		article.Title = title
		article.Language = model.LangEn

		articles = append(articles, article)
	}

	return articles, nil
}
