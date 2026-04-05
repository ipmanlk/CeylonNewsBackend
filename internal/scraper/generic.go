package scraper

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
	"strings"
)

type GenericScraper struct {
	config  SourceConfig
	fetcher *fetcher.Fetcher
}

func NewGenericScraper(cfg SourceConfig, f *fetcher.Fetcher) *GenericScraper {
	return &GenericScraper{config: cfg, fetcher: f}
}

func (s *GenericScraper) Name() string { return s.config.Name }

func (s *GenericScraper) Languages() []model.Language {
	langs := make([]model.Language, 0, len(s.config.Languages))
	for _, lc := range s.config.Languages {
		langs = append(langs, model.Language(lc.Language))
	}
	return langs
}

func (s *GenericScraper) Scrape(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	for _, lc := range s.config.Languages {
		if model.Language(lc.Language) == language {
			return s.scrapeLanguage(ctx, lc)
		}
	}
	return nil, nil
}

func (s *GenericScraper) scrapeLanguage(ctx context.Context, lc LangConfig) ([]model.ScrapedArticle, error) {
	switch lc.Strategy {
	case StrategyRSS:
		return s.scrapeRSS(ctx, lc, false, false, "")
	case StrategyRSSBrowser:
		return s.scrapeRSSBrowser(ctx, lc)
	case StrategyRSSArticleBrowser:
		return s.scrapeRSS(ctx, lc, false, true, "")
	case StrategyRSSContentSelector:
		return s.scrapeRSS(ctx, lc, false, true, lc.ContentSelector)
	case StrategyHTML:
		return s.scrapeHTML(ctx, lc, false)
	case StrategyHTMLBrowser:
		return s.scrapeHTML(ctx, lc, true)
	default:
		return nil, nil
	}
}

func (s *GenericScraper) scrapeRSS(ctx context.Context, lc LangConfig, feedBrowser bool, articleBrowser bool, contentSelector string) ([]model.ScrapedArticle, error) {
	items, err := s.fetcher.FetchRSS(ctx, lc.FeedURL, lc.MaxItems)
	if err != nil {
		return nil, err
	}

	articles := make([]model.ScrapedArticle, 0, len(items))
	seen := make(map[string]bool)
	tt := s.config.TitleTransform

	for _, item := range items {
		if seen[item.Link] {
			continue
		}
		seen[item.Link] = true

		title := item.Title

		// Apply skip rules
		shouldSkip := false
		for _, rule := range tt.Skip {
			if rule.CaseSensitive {
				if strings.Contains(title, rule.Contains) {
					shouldSkip = true
					break
				}
			} else {
				if strings.Contains(strings.ToLower(title), strings.ToLower(rule.Contains)) {
					shouldSkip = true
					break
				}
			}
		}

		if shouldSkip {
			slog.Debug("skipping article by title filter", "source", s.Name(), "url", item.Link)
			continue
		}

		var article model.ScrapedArticle
		var err error

		if contentSelector != "" {
			article, err = s.fetcher.ExtractArticleFromRSSItemWithSelector(ctx, item, contentSelector, articleBrowser)
		} else {
			article, err = s.fetcher.ExtractArticleFromRSSItem(ctx, item, articleBrowser)
		}

		if err != nil {
			slog.Warn("failed to extract article", "source", s.Name(), "url", item.Link, "error", err)
			continue
		}

		if article.Title == "" || article.ContentText == "" || article.ContentHTML == "" {
			slog.Debug("skipping article with missing content", "source", s.Name(), "url", item.Link)
			continue
		}

		// Apply replacement rules
		for _, rule := range tt.Replace {
			if rule.CaseSensitive {
				article.Title = strings.ReplaceAll(article.Title, rule.Pattern, rule.With)
			} else {
				article.Title = replaceAllCaseInsensitive(article.Title, rule.Pattern, rule.With)
			}
		}

		article.SourceName = s.Name()
		article.Language = model.Language(lc.Language)
		articles = append(articles, article)
	}

	return articles, nil
}

func replaceAllCaseInsensitive(s, old, new string) string {
	if old == "" {
		return s
	}
	lowerS := strings.ToLower(s)
	lowerOld := strings.ToLower(old)

	var result strings.Builder
	start := 0
	for {
		idx := strings.Index(lowerS[start:], lowerOld)
		if idx == -1 {
			result.WriteString(s[start:])
			break
		}
		idx += start
		result.WriteString(s[start:idx])
		result.WriteString(new)
		start = idx + len(old)
	}
	return result.String()
}

func (s *GenericScraper) scrapeRSSBrowser(ctx context.Context, lc LangConfig) ([]model.ScrapedArticle, error) {
	items, err := s.fetcher.FetchRSSWithBrowser(ctx, lc.FeedURL, lc.MaxItems)
	if err != nil {
		return nil, err
	}

	articles := make([]model.ScrapedArticle, 0, len(items))
	seen := make(map[string]bool)

	for _, item := range items {
		if seen[item.Link] {
			continue
		}
		seen[item.Link] = true

		article, err := s.fetcher.ExtractArticleFromRSSItem(ctx, item)
		if err != nil {
			slog.Warn("failed to extract article", "source", s.Name(), "url", item.Link, "error", err)
			continue
		}

		if article.Title == "" || article.ContentText == "" || article.ContentHTML == "" {
			slog.Debug("skipping article with missing content", "source", s.Name(), "url", item.Link)
			continue
		}

		article.SourceName = s.Name()
		article.Language = model.Language(lc.Language)
		articles = append(articles, article)
	}
	return articles, nil
}

func (s *GenericScraper) scrapeHTML(ctx context.Context, lc LangConfig, useBrowser bool) ([]model.ScrapedArticle, error) {
	doc, err := s.fetcher.FetchHTMLDoc(ctx, lc.PageURL, useBrowser)
	if err != nil {
		return nil, err
	}

	var allLinks []string
	for _, sel := range lc.LinkSelectors {
		links := s.fetcher.ExtractLinks(doc, sel, lc.URLPrefix)
		allLinks = append(allLinks, links...)
	}

	if lc.BaseURL != "" {
		for i, link := range allLinks {
			if strings.HasPrefix(link, "/") {
				allLinks[i] = lc.BaseURL + link
			}
		}
	}

	if lc.MaxItems > 0 && len(allLinks) > lc.MaxItems {
		allLinks = allLinks[:lc.MaxItems]
	}

	articles := make([]model.ScrapedArticle, 0, len(allLinks))
	seen := make(map[string]bool)

	for _, link := range allLinks {
		if seen[link] {
			continue
		}
		seen[link] = true

		result, err := s.fetcher.ExtractArticle(ctx, link, useBrowser)
		if err != nil {
			slog.Warn("failed to extract article", "source", s.Name(), "url", link, "error", err)
			continue
		}

		if result == nil || result.Metadata.Title == "" || result.ContentText == "" {
			slog.Debug("skipping article with missing content", "source", s.Name(), "url", link)
			continue
		}

		article := s.fetcher.CreateScrapedArticle(s.Name(), result, link, &result.Metadata.Image, result.Metadata.Date)
		if article.ContentText == "" || article.ContentHTML == "" {
			slog.Debug("skipping article with missing extracted content", "source", s.Name(), "url", link)
			continue
		}
		article.Language = model.Language(lc.Language)
		articles = append(articles, article)
	}

	return articles, nil
}
