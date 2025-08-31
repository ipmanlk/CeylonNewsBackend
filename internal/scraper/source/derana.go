package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
)

type DeranaScraper struct {
	fetcher *fetcher.Fetcher
}

func NewDeranaScraper(fetcher *fetcher.Fetcher) *DeranaScraper {
	return &DeranaScraper{
		fetcher: fetcher,
	}
}

func (s *DeranaScraper) Name() string {
	return "Derana"
}

func (s *DeranaScraper) Languages() []model.Language {
	return []model.Language{model.LangSi, model.LangEn, model.LangTa}
}

func (s *DeranaScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
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

func (s *DeranaScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, "https://www.adaderana.lk/hot-news/", true)
	if err != nil {
		return nil, err
	}

	articleLinks := s.fetcher.ExtractLinks(doc, ".story-text h2 a", "https://www.adaderana.lk/news/")
	articleLinks = articleLinks[:5]

	return s.scrapeArticles(ctx, articleLinks, model.LangEn)
}

func (s *DeranaScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	articles, err := s.fetcher.FetchRSS(ctx, "https://sinhala.adaderana.lk/rsshotnews.php")
	if err != nil {
		return nil, err
	}

	for i := range articles {
		articles[i].SourceName = s.Name()
		articles[i].Language = model.LangSi
	}

	return articles, nil
}

func (s *DeranaScraper) scrapeTa(ctx context.Context) ([]model.ScrapedArticle, error) {
	// TODO: Implement Tamil scraping
	return nil, nil
}

func (s *DeranaScraper) scrapeArticles(ctx context.Context, links []string, language model.Language) ([]model.ScrapedArticle, error) {
	articles := make([]model.ScrapedArticle, 0, len(links))
	seenLinks := make(map[string]bool)

	for _, link := range links {
		if seenLinks[link] {
			continue
		}
		seenLinks[link] = true

		result, err := s.fetcher.ExtractArticle(ctx, link, true)
		if err != nil {
			slog.Warn("failed to extract article", "scraper", "Derana", "url", link, "error", err)
			continue
		}

		if result == nil || result.Metadata.Title == "" || result.ContentText == "" {
			slog.Debug("skipping article with missing content", "scraper", "Derana", "url", link)
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
