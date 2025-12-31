package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
	"strings"
)

type NewsLKScraper struct {
	fetcher *fetcher.Fetcher
}

func NewNewsLKScraper(fetcher *fetcher.Fetcher) *NewsLKScraper {
	return &NewsLKScraper{
		fetcher: fetcher,
	}
}

func (s *NewsLKScraper) Name() string {
	return "News.lk"
}

func (s *NewsLKScraper) Languages() []model.Language {
	return []model.Language{model.LangEn, model.LangSi, model.LangTa}
}

func (s *NewsLKScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
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

func (s *NewsLKScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://news.lk/news/", true)
	if err != nil {
		return nil, err
	}

	var articleLinks []string

	articleLinks = s.fetcher.ExtractLinks(doc, "article.item h2 a", "/news/")

	// Convert relative URLs to absolute URLs
	for i, link := range articleLinks {
		if strings.HasPrefix(link, "/") {
			articleLinks[i] = "https://news.lk" + link
		}
	}

	if len(articleLinks) > 2 {
		articleLinks = articleLinks[:2]
	}

	return s.scrapeArticles(ctx, articleLinks, model.LangEn)
}

func (s *NewsLKScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://sinhala.news.lk/news/", true)
	if err != nil {
		return nil, err
	}

	var articleLinks []string

	articleLinks = s.fetcher.ExtractLinks(doc, "article.item h2 a", "/news/")

	// Convert relative URLs to absolute URLs
	for i, link := range articleLinks {
		if strings.HasPrefix(link, "/") {
			articleLinks[i] = "https://sinhala.news.lk" + link
		}
	}

	if len(articleLinks) > 2 {
		articleLinks = articleLinks[:2]
	}

	return s.scrapeArticles(ctx, articleLinks, model.LangSi)
}

func (s *NewsLKScraper) scrapeTa(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://tamil.news.lk/news/", true)
	if err != nil {
		return nil, err
	}

	var articleLinks []string

	articleLinks = s.fetcher.ExtractLinks(doc, "article.item h2 a", "/news/")

	// Convert relative URLs to absolute URLs
	for i, link := range articleLinks {
		if strings.HasPrefix(link, "/") {
			articleLinks[i] = "https://tamil.news.lk" + link
		}
	}

	if len(articleLinks) > 2 {
		articleLinks = articleLinks[:2]
	}

	return s.scrapeArticles(ctx, articleLinks, model.LangTa)
}

func (s *NewsLKScraper) scrapeArticles(ctx context.Context, links []string, language model.Language) ([]model.ScrapedArticle, error) {
	articles := make([]model.ScrapedArticle, 0, len(links))
	seenLinks := make(map[string]bool)

	for _, link := range links {
		if seenLinks[link] {
			continue
		}
		seenLinks[link] = true

		result, err := s.fetcher.ExtractArticle(ctx, link, true)
		if err != nil {
			slog.Warn("failed to extract article", "scraper", "News.lk", "url", link, "error", err)
			continue
		}

		if result == nil || result.Metadata.Title == "" || result.ContentText == "" {
			slog.Debug("skipping article with missing content", "scraper", "News.lk", "url", link)
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
