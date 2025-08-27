package source

import (
	"context"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BBCScraper struct {
}

func NewBBCScraper(
	rssScraper *scraper.RSSScraper,
) *BBCScraper {
	return &BBCScraper{}
}

func (s *BBCScraper) Name() string {
	return "BBC"
}

func (s *BBCScraper) Languages() []model.Language {
	return []model.Language{model.LangSi, model.LangEn, model.LangTa}
}

func (s *BBCScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	return s.scrapeEn(ctx)
}

func (s *BBCScraper) scrapeEn(ctx context.Context) ([]model.ScrapedArticle, error) {
	doc, err := scraper.GetGoQueryDocFromURL(ctx, "https://www.bbc.com/news/topics/cywd23g0gxgt")
	if err != nil {
		return nil, err
	}

	articleLinks := []string{}
	doc.Find("a[class*='hMvGwj']").Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if exists && href != "" && (strings.HasPrefix(href, "/news/articles/")) {
			href = "https://www.bbc.com" + href
			articleLinks = append(articleLinks, href)
		}
	})

	articles := make([]model.ScrapedArticle, 0, len(articleLinks))
	seenLinks := make(map[string]bool)
	for _, link := range articleLinks {
		if seenLinks[link] {
			continue
		}
		seenLinks[link] = true

		result, err := scraper.ScrapeArticleFromURL(ctx, link)
		if err != nil {
			continue
		}

		if result == nil || result.Metadata.Title == "" || result.ContentText == "" {
			continue
		}

		article := model.ScrapedArticle{
			SourceName:  s.Name(),
			Title:       result.Metadata.Title,
			ContentText: result.ContentText,
			ContentHTML: "",
			URL:         link,
			ImageURL:    &result.Metadata.Image,
			Language:    model.LangEn,
			PublishedAt: result.Metadata.Date,
		}

		articles = append(articles, article)
	}
	return articles, nil
}
