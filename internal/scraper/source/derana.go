package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
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
	articles, err := s.fetcher.FetchRSS(ctx, "https://www.adaderana.lk/rss.php")
	if err != nil {
		return nil, err
	}

	for i := range articles {
		articles[i].SourceName = s.Name()
		articles[i].Language = model.LangEn
	}

	return articles, nil
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
	articles, err := s.fetcher.FetchRSS(ctx, "http://tamil.adaderana.lk/rss.php")
	if err != nil {
		return nil, err
	}

	for i := range articles {
		articles[i].SourceName = s.Name()
		articles[i].Language = model.LangTa
	}

	return articles, nil
}
