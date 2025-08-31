package source

import (
	"context"
	"errors"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
)

type DivainaScraper struct {
	fetcher *fetcher.Fetcher
}

func NewDivainaScraper(fetcher *fetcher.Fetcher) *DivainaScraper {
	return &DivainaScraper{
		fetcher: fetcher,
	}
}

func (s *DivainaScraper) Name() string {
	return "Divaina"
}

func (s *DivainaScraper) Languages() []model.Language {
	return []model.Language{model.LangSi}
}

func (s *DivainaScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	switch language {
	case model.LangSi:
		return s.scrapeSi(ctx)
	default:
		return nil, errors.New("unsupported language for Divaina scraper")
	}
}

func (s *DivainaScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	articles, err := s.fetcher.FetchRSS(ctx, "https://www.divaina.com/rss.php")
	if err != nil {
		return nil, err
	}

	// Set source name and language for all articles
	for i := range articles {
		articles[i].SourceName = s.Name()
		articles[i].Language = model.LangSi
	}

	slog.Info("scraped Divaina articles",
		"scraper", "Divaina",
		"count", len(articles),
		"language", model.LangSi,
	)

	return articles, nil
}
