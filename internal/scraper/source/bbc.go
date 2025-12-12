package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
	"strings"
)

type BBCScraper struct {
	fetcher *fetcher.Fetcher
}

func NewBBCScraper(fetcher *fetcher.Fetcher) *BBCScraper {
	return &BBCScraper{
		fetcher: fetcher,
	}
}

func (s *BBCScraper) Name() string {
	return "BBC"
}

func (s *BBCScraper) Languages() []model.Language {
	return []model.Language{model.LangSi, model.LangEn, model.LangTa}
}

func (s *BBCScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangEn:
		return s.scrapeEn(ctx)

	case model.LangSi:
		return s.scrapeSi(ctx)

	case model.LangTa:
		return s.scrapeTa(ctx)
	default:
		return nil, nil
	}
}

func (s *BBCScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://www.bbc.com/news/topics/cywd23g0gxgt")
	if err != nil {
		return nil, err
	}

	articleLinks := s.fetcher.ExtractLinks(doc, "a", "/news/articles/")
	if len(articleLinks) > 5 {
		articleLinks = articleLinks[:5]
	}

	for i, link := range articleLinks {
		if strings.HasPrefix(link, "/") {
			articleLinks[i] = "https://www.bbc.com" + link
		}
	}

	return s.scrapeArticles(ctx, articleLinks, model.LangEn)
}

func (s *BBCScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://www.bbc.com/sinhala/topics/cg7267dz901t")
	if err != nil {
		return nil, err
	}

	articleLinks := s.fetcher.ExtractLinks(doc, "a", "https://www.bbc.com/sinhala/articles/")
	if len(articleLinks) > 5 {
		articleLinks = articleLinks[:5]
	}

	return s.scrapeArticles(ctx, articleLinks, model.LangSi)
}

func (s *BBCScraper) scrapeTa(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://www.bbc.com/tamil/topics/cz74k7p3qw7t")
	if err != nil {
		return nil, err
	}

	articleLinks := s.fetcher.ExtractLinks(doc, "a", "https://www.bbc.com/tamil/articles/")
	if len(articleLinks) > 5 {
		articleLinks = articleLinks[:5]
	}

	return s.scrapeArticles(ctx, articleLinks, model.LangTa)
}

func (s *BBCScraper) scrapeArticles(ctx context.Context, links []string, language model.Language) ([]model.ScrapedArticle, error) {
	articles := make([]model.ScrapedArticle, 0, len(links))
	seenLinks := make(map[string]bool)

	for _, link := range links {
		if seenLinks[link] {
			continue
		}
		seenLinks[link] = true

		result, err := s.fetcher.ExtractArticle(ctx, link)
		if err != nil {
			slog.Warn("failed to extract article", "scraper", "BBC", "url", link, "error", err)
			continue
		}

		if result == nil || result.Metadata.Title == "" || result.ContentText == "" {
			slog.Debug("skipping article with missing content", "scraper", "BBC", "url", link)
			continue
		}

		article := model.ScrapedArticle{
			SourceName:  s.Name(),
			Title:       result.Metadata.Title,
			Content:     result.ContentText,
			URL:         link,
			ImageURL:    &result.Metadata.Image,
			Language:    language,
			PublishedAt: result.Metadata.Date,
		}

		articles = append(articles, article)
	}

	return articles, nil
}
