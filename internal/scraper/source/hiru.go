package source

import (
	"context"
	"errors"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
)

type HiruScraper struct {
	fetcher *fetcher.Fetcher
}

func NewHiruScraper(fetcher *fetcher.Fetcher) *HiruScraper {
	return &HiruScraper{
		fetcher: fetcher,
	}
}

func (s *HiruScraper) Name() string {
	return "Hiru News"
}

func (s *HiruScraper) Languages() []model.Language {
	return []model.Language{model.LangSi, model.LangEn, model.LangTa}
}

func (s *HiruScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangSi:
		return s.scrapeSi(ctx)
	case model.LangEn:
		return s.scrapeEn(ctx)
	case model.LangTa:
		return s.scrapeTa(ctx)
	default:
		return nil, errors.New("unsupported language for Hiru News scraper")
	}
}

func (s *HiruScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://www.hirunews.lk")
	if err != nil {
		return nil, err
	}

	// Extract featured article links (card-featured class)
	featuredLinks := s.fetcher.ExtractLinks(doc, "a.card-featured", "https://hirunews.lk/")

	// Extract latest news article links (card-v1 class)
	latestLinks := s.fetcher.ExtractLinks(doc, "a.card-v1", "https://hirunews.lk/")

	// Combine both types of links
	allLinks := append(featuredLinks, latestLinks...)

	// Limit to 5 articles total
	if len(allLinks) > 5 {
		allLinks = allLinks[:5]
	}

	return s.scrapeArticles(ctx, allLinks, model.LangSi)
}

func (s *HiruScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://www.hirunews.lk/en/")
	if err != nil {
		return nil, err
	}

	// Extract featured article links (card-featured class)
	featuredLinks := s.fetcher.ExtractLinks(doc, "a.card-featured", "https://hirunews.lk/")

	// Extract latest news article links (card-v1 class)
	latestLinks := s.fetcher.ExtractLinks(doc, "a.card-v1", "https://hirunews.lk/")

	// Combine both types of links
	allLinks := append(featuredLinks, latestLinks...)

	// Limit to 5 articles total
	if len(allLinks) > 5 {
		allLinks = allLinks[:5]
	}

	return s.scrapeArticles(ctx, allLinks, model.LangEn)
}

func (s *HiruScraper) scrapeTa(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://www.hirunews.lk/tm/")
	if err != nil {
		return nil, err
	}

	// Extract featured article links (card-featured class)
	featuredLinks := s.fetcher.ExtractLinks(doc, "a.card-featured", "https://hirunews.lk/")

	// Extract latest news article links (card-v1 class)
	latestLinks := s.fetcher.ExtractLinks(doc, "a.card-v1", "https://hirunews.lk/")

	// Combine both types of links
	allLinks := append(featuredLinks, latestLinks...)

	// Limit to 5 articles total
	if len(allLinks) > 5 {
		allLinks = allLinks[:5]
	}

	return s.scrapeArticles(ctx, allLinks, model.LangTa)
}

func (s *HiruScraper) scrapeArticles(ctx context.Context, links []string, language model.Language) ([]model.ScrapedArticle, error) {
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
