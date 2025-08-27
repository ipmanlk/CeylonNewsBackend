package scraper

import (
	"context"
	"errors"
	"fmt"
	"ipmanlk/cnapi/internal/model"
	"log/slog"
	"time"

	"github.com/mmcdole/gofeed"
)

type RSSScraper struct {
	logger *slog.Logger
}

func NewRSSScraper() *RSSScraper {
	return &RSSScraper{}
}

func (r *RSSScraper) Scrape(ctx context.Context, url string) ([]model.ScrapedArticle, error) {
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

		var imageURL = r.getImageURL(item)
		contentText, contentHTML, err := r.getContent(item)
		if err != nil {
			r.logger.Warn("failed to get content from item, skipping",
				"feed_url", url,
				"item_link", item.Link,
				"error", err,
			)
			continue
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

		articles = append(articles, article)
		articleURLs[item.Link] = struct{}{}
	}

	return articles, nil
}

func (r *RSSScraper) getImageURL(item *gofeed.Item) *string {
	if item.Image != nil {
		return &item.Image.URL
	}

	if url := GetFirstImageFromHTML([]byte(item.Content)); url != nil {
		return url
	}

	if url := GetFirstImageFromHTML([]byte(item.Description)); url != nil {
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

func (r *RSSScraper) getContent(item *gofeed.Item) (string, string, error) {
	var html string
	if item.Content != "" {
		html = item.Content
	} else {
		html = item.Description
	}

	if html == "" {
		return "", "", errors.New("no content available in item")
	}

	cleanedHTML, err := CleanHTMLString(html)
	if err != nil {
		return "", "", err
	}

	text, err := GetTextFromHTML(cleanedHTML)
	if err != nil {
		return "", "", err
	}

	return text, cleanedHTML, nil
}

func (r *RSSScraper) getPublishedAt(item *gofeed.Item) time.Time {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed
	}
	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed
	}
	return time.Now()
}
