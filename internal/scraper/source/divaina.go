package source

import (
	"context"
	"errors"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper"
)

type DivainaScraper struct {
	rssScraper *scraper.RSSScraper
}

func NewDivainaScraper(
	rssScraper *scraper.RSSScraper,
) *DivainaScraper {
	return &DivainaScraper{
		rssScraper: rssScraper,
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
	articles, err := s.rssScraper.Scrape(ctx, "https://www.divaina.com/rss.php")
	if err != nil {
		return nil, err
	}

	return articles, nil
}
