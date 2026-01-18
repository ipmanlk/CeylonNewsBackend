package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
	"strings"
)

type GaganaScraper struct {
	fetcher *fetcher.Fetcher
}

func NewGaganaScraper(fetcher *fetcher.Fetcher) *GaganaScraper {
	return &GaganaScraper{
		fetcher: fetcher,
	}
}

func (s *GaganaScraper) Name() string {
	return "Gagana"
}

func (s *GaganaScraper) Languages() []model.Language {
	return []model.Language{model.LangSi}
}

func (s *GaganaScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangSi:
		return s.scrapeSi(ctx)
	default:
		return nil, nil
	}
}

func (s *GaganaScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://gagana.lk/news/srilanka")
	if err != nil {
		return nil, err
	}

	// Extract article links from the news page
	// Based on the HTML structure: div.a-item h4 a
	articleLinks := s.fetcher.ExtractLinks(doc, "div.a-item h4 a", "https://gagana.lk/article/")

	// Convert relative URLs to absolute URLs
	for i, link := range articleLinks {
		if strings.HasPrefix(link, "/") {
			articleLinks[i] = "https://gagana.lk" + link
		}
	}

	// Limit to 5 articles
	if len(articleLinks) > 5 {
		articleLinks = articleLinks[:5]
	}

	return s.scrapeArticles(ctx, articleLinks, model.LangSi)
}

func (s *GaganaScraper) scrapeArticles(ctx context.Context, links []string, language model.Language) ([]model.ScrapedArticle, error) {
	articles := make([]model.ScrapedArticle, 0, len(links))
	seenLinks := make(map[string]bool)

	for _, link := range links {
		if seenLinks[link] {
			continue
		}
		seenLinks[link] = true

		result, err := s.fetcher.ExtractArticle(ctx, link)
		if err != nil {
			slog.Warn("failed to extract article", "url", link, "error", err)
			continue
		}

		if result == nil || result.Metadata.Title == "" || result.ContentText == "" {
			slog.Debug("skipping article with missing content", "url", link)
			continue
		}

		article := s.fetcher.CreateScrapedArticle(s.Name(), result, link, &result.Metadata.Image, result.Metadata.Date)
		if article.ContentText == "" || article.ContentHTML == "" {
			slog.Debug("skipping article with missing extracted content", "url", link)
			continue
		}
		article.Language = language

		articles = append(articles, article)
	}

	return articles, nil
}
