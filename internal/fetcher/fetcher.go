package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"ipmanlk/cnapi/internal/model"
	"log/slog"

	"github.com/PuerkitoBio/goquery"
	"github.com/markusmobius/go-trafilatura"
	"github.com/mmcdole/gofeed"
)

type Fetcher struct {
	httpClient    *HTTPClient
	browserClient *BrowserAPIClient
	logger        *slog.Logger
}

func NewFetcher(httpClient *HTTPClient, browserClient *BrowserAPIClient, logger *slog.Logger) *Fetcher {
	return &Fetcher{
		httpClient:    httpClient,
		browserClient: browserClient,
		logger:        logger,
	}
}

// FetchHTML fetches HTML content, optionally using browser API
func (f *Fetcher) FetchHTML(ctx context.Context, url string, useBrowser ...bool) ([]byte, error) {
	if len(useBrowser) > 0 && useBrowser[0] {
		return f.browserClient.FetchHTML(ctx, url)
	}
	return f.httpClient.FetchHTML(ctx, url)
}

// FetchHTMLDoc fetches and parses HTML into goquery document
func (f *Fetcher) FetchHTMLDoc(ctx context.Context, url string, useBrowser ...bool) (*goquery.Document, error) {
	if len(useBrowser) > 0 && useBrowser[0] {
		return f.browserClient.FetchHTMLDoc(ctx, url)
	}
	return f.httpClient.FetchHTMLDoc(ctx, url)
}

// ExtractArticle extracts article content from URL using Trafilatura
func (f *Fetcher) ExtractArticle(ctx context.Context, url string, useBrowser ...bool) (*trafilatura.ExtractResult, error) {
	html, err := f.FetchHTML(ctx, url, useBrowser...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HTML from %s: %w", url, err)
	}

	opts := trafilatura.Options{
		IncludeLinks:    true,
		IncludeImages:   true,
		ExcludeComments: true,
		EnableFallback:  true,
		Deduplicate:     true,
	}

	result, err := trafilatura.Extract(bytes.NewReader(html), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content from %s: %w", url, err)
	}

	return result, nil
}

// FetchRSS fetches and parses RSS feed
func (f *Fetcher) FetchRSS(ctx context.Context, url string) ([]model.ScrapedArticle, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	articles := make([]model.ScrapedArticle, 0, len(feed.Items))
	articleURLs := make(map[string]struct{})

	for _, item := range feed.Items {
		if item.Link == "" {
			f.logger.Warn("skipping item with empty link", "feed_url", url, "item_title", item.Title)
			continue
		}

		if _, exists := articleURLs[item.Link]; exists {
			f.logger.Debug("skipping duplicate item", "feed_url", url, "item_link", item.Link)
			continue
		}

		article, err := f.processRSSItem(ctx, item)
		if err != nil {
			f.logger.Warn("failed to process item, skipping", "feed_url", url, "item_link", item.Link, "error", err)
			continue
		}

		articles = append(articles, article)
		articleURLs[item.Link] = struct{}{}
	}

	return articles, nil
}

func (f *Fetcher) processRSSItem(ctx context.Context, item *gofeed.Item) (model.ScrapedArticle, error) {
	result, err := f.ExtractArticle(ctx, item.Link)
	if err != nil {
		return model.ScrapedArticle{}, fmt.Errorf("failed to extract article: %w", err)
	}

	imageURL := f.getImageURL(item)

	article := model.ScrapedArticle{
		Title:       item.Title,
		URL:         item.Link,
		Content:     result.ContentText,
		ImageURL:    imageURL,
		Categories:  item.Categories,
		PublishedAt: f.getPublishedAt(item),
	}

	return article, nil
}

func (f *Fetcher) getImageURL(item *gofeed.Item) *string {
	if item.Image != nil {
		return &item.Image.URL
	}

	if url := f.getFirstImageFromHTML([]byte(item.Content)); url != nil {
		return url
	}

	if url := f.getFirstImageFromHTML([]byte(item.Description)); url != nil {
		return url
	}

	var imageTypes = map[string]struct{}{
		"image/jpeg": {},
		"image/png":  {},
		"image/webp": {},
		"image/jpg":  {},
	}

	if len(item.Enclosures) > 0 {
		for _, enclosure := range item.Enclosures {
			if _, ok := imageTypes[enclosure.Type]; ok {
				return &enclosure.URL
			}
		}
	}

	return nil
}

func (f *Fetcher) getPublishedAt(item *gofeed.Item) time.Time {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed
	}
	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed
	}
	return time.Now()
}

// HTML utility methods
func (f *Fetcher) getFirstImageFromHTML(html []byte) *string {
	var thumbnailURL string
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err == nil {
		thumbnailURL = doc.Find("img").First().AttrOr("src", "")
	}
	if thumbnailURL != "" {
		return &thumbnailURL
	}
	return nil
}

func (f *Fetcher) ExtractLinks(doc *goquery.Document, selector, urlPattern string) []string {
	var links []string
	doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if exists && href != "" && strings.HasPrefix(href, urlPattern) {
			links = append(links, href)
		}
	})
	return links
}
