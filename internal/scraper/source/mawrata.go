package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
)

type MawrataScraper struct {
	fetcher *fetcher.Fetcher
}

func NewMawrataScraper(fetcher *fetcher.Fetcher) *MawrataScraper {
	return &MawrataScraper{
		fetcher: fetcher,
	}
}

func (s *MawrataScraper) Name() string {
	return "Mawrata"
}

func (s *MawrataScraper) Languages() []model.Language {
	return []model.Language{model.LangEn, model.LangSi}
}

func (s *MawrataScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangEn:
		return s.scrapeEn(ctx)
	case model.LangSi:
		return s.scrapeSi(ctx)
	default:
		return nil, nil
	}
}

func (s *MawrataScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	// Fetch RSS feed
	rssItems, err := s.fetcher.FetchRSS(ctx, "https://mawratanews.lk/feed/", 5)
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
			slog.Warn("failed to extract article", "url", item.Link, "error", err)
			continue
		}

		if article.Title == "" || article.ContentText == "" {
			slog.Debug("skipping article with missing content", "url", item.Link)
			continue
		}

		if article.ContentText == "" || article.ContentHTML == "" {
			slog.Debug("skipping article with missing extracted content", "url", item.Link)
			continue
		}

		// Update the article with our source-specific information
		article.SourceName = s.Name()
		article.Language = model.LangEn

		articles = append(articles, article)
	}

	return articles, nil
}

func (s *MawrataScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	// Fetch RSS feed
	rssItems, err := s.fetcher.FetchRSS(ctx, "https://sinhala.mawratanews.lk/feed/", 5)
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
			slog.Warn("failed to extract article", "url", item.Link, "error", err)
			continue
		}

		if article.Title == "" || article.ContentText == "" {
			slog.Debug("skipping article with missing content", "url", item.Link)
			continue
		}

		if article.ContentText == "" || article.ContentHTML == "" {
			slog.Debug("skipping article with missing extracted content", "url", item.Link)
			continue
		}

		// Update the article with our source-specific information
		article.SourceName = s.Name()
		article.Language = model.LangSi

		articles = append(articles, article)
	}

	return articles, nil
}
