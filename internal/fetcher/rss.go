package fetcher

import (
	"context"
	"errors"
	"fmt"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// RSSFetcher handles RSS feed fetching and parsing
type RSSFetcher struct {
	httpClient       *HTTPClient
	htmlProcessor    *HTMLProcessor
	contentExtractor *ContentExtractor
	logger           *slog.Logger
}

// NewRSSFetcher creates a new RSS fetcher
func NewRSSFetcher(httpClient *HTTPClient, htmlProcessor *HTMLProcessor, logger *slog.Logger) *RSSFetcher {
	return &RSSFetcher{
		httpClient:    httpClient,
		htmlProcessor: htmlProcessor,
		logger:        logger,
	}
}

// SetContentExtractor sets the content extractor for fallback content extraction
func (r *RSSFetcher) SetContentExtractor(extractor *ContentExtractor) {
	r.contentExtractor = extractor
}

// FetchArticles fetches and parses articles from an RSS feed
func (r *RSSFetcher) FetchArticles(ctx context.Context, url string) ([]model.ScrapedArticle, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	articles := make([]model.ScrapedArticle, 0, len(feed.Items))
	articleURLs := make(map[string]struct{}) // to skip duplicates present in some feeds

	for _, item := range feed.Items {
		if item.Link == "" {
			r.logger.Warn("skipping item with empty link",
				"feed_url", url,
				"item_title", item.Title,
			)
			continue
		}

		if _, exists := articleURLs[item.Link]; exists {
			r.logger.Debug("skipping duplicate item",
				"feed_url", url,
				"item_link", item.Link,
			)
			continue
		}

		article, err := r.processItem(ctx, item, url)
		if err != nil {
			r.logger.Warn("failed to process item, skipping",
				"feed_url", url,
				"item_link", item.Link,
				"error", err,
			)
			continue
		}

		articles = append(articles, article)
		articleURLs[item.Link] = struct{}{}
	}

	return articles, nil
}

// processItem processes a single RSS item, using extractor if needed
func (r *RSSFetcher) processItem(ctx context.Context, item *gofeed.Item, feedURL string) (model.ScrapedArticle, error) {
	var imageURL = r.getImageURL(item)
	contentText, contentHTML, err := r.getContent(item)
	if err != nil {
		return model.ScrapedArticle{}, err
	}

	// If RSS content is too short and we have an extractor, try to get full content
	if r.contentExtractor != nil && len(strings.TrimSpace(contentText)) < 200 {
		r.logger.Debug("RSS content too short, trying extractor",
			"item_link", item.Link,
			"content_length", len(contentText),
		)

		extractedResult, err := r.contentExtractor.ExtractArticleFromURL(ctx, item.Link)
		if err != nil {
			r.logger.Warn("failed to extract content from URL, using RSS content",
				"item_link", item.Link,
				"error", err,
			)
		} else if extractedResult != nil && extractedResult.ContentText != "" {
			// Use extracted content if available
			contentText = extractedResult.ContentText
			if extractedResult.Metadata.Title != "" {
				item.Title = extractedResult.Metadata.Title
			}
			if extractedResult.Metadata.Image != "" {
				imageURL = &extractedResult.Metadata.Image
			}
			if !extractedResult.Metadata.Date.IsZero() {
				item.PublishedParsed = &extractedResult.Metadata.Date
			}
		}
	}

	article := model.ScrapedArticle{
		Title:       item.Title,
		URL:         item.Link,
		ContentText: contentText,
		ContentHTML: contentHTML,
		ImageURL:    imageURL,
		Categories:  item.Categories,
		PublishedAt: r.getPublishedAt(item),
	}

	return article, nil
}

func (r *RSSFetcher) getImageURL(item *gofeed.Item) *string {
	if item.Image != nil {
		return &item.Image.URL
	}

	if url := r.htmlProcessor.GetFirstImageFromHTML([]byte(item.Content)); url != nil {
		return url
	}

	if url := r.htmlProcessor.GetFirstImageFromHTML([]byte(item.Description)); url != nil {
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

func (r *RSSFetcher) getContent(item *gofeed.Item) (string, string, error) {
	var html string
	if item.Content != "" {
		html = item.Content
	} else {
		html = item.Description
	}

	if html == "" {
		return "", "", errors.New("no content available in item")
	}

	cleanedHTML, err := r.htmlProcessor.CleanHTMLString(html)
	if err != nil {
		return "", "", err
	}

	text, err := r.htmlProcessor.GetTextFromHTML(cleanedHTML)
	if err != nil {
		return "", "", err
	}

	return text, cleanedHTML, nil
}

func (r *RSSFetcher) getPublishedAt(item *gofeed.Item) time.Time {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed
	}
	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed
	}
	return time.Now()
}
