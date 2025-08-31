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
	items, err := s.fetcher.FetchRSS(ctx, "https://www.adaderana.lk/rss.php", 5)
	if err != nil {
		return nil, err
	}

	articles := make([]model.ScrapedArticle, 0, len(items))
	for _, item := range items {
		article, err := s.fetcher.ExtractArticleFromRSSItem(ctx, item)
		if err != nil {
			continue
		}
		article.SourceName = s.Name()
		article.Language = model.LangEn
		articles = append(articles, article)
	}

	return articles, nil
}
func (s *DeranaScraper) scrapeSi(ctx context.Context) ([]model.ScrapedArticle, error) {
	items, err := s.fetcher.FetchRSS(ctx, "https://sinhala.adaderana.lk/rsshotnews.php", 5)
	if err != nil {
		return nil, err
	}

	articles := make([]model.ScrapedArticle, 0, len(items))
	for _, item := range items {
		article, err := s.fetcher.ExtractArticleFromRSSItem(ctx, item)
		if err != nil {
			continue
		}
		article.SourceName = s.Name()
		article.Language = model.LangSi
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *DeranaScraper) scrapeTa(ctx context.Context) ([]model.ScrapedArticle, error) {
	items, err := s.fetcher.FetchRSSWithBrowser(ctx, "http://tamil.adaderana.lk/rss.php", 5)
	if err != nil {
		return nil, err
	}

	articles := make([]model.ScrapedArticle, 0, len(items))
	for _, item := range items {
		article, err := s.fetcher.ExtractArticleFromRSSItem(ctx, item)
		if err != nil {
			continue
		}
		article.SourceName = s.Name()
		article.Language = model.LangTa
		articles = append(articles, article)
	}

	return articles, nil
}
