package source

import (
	"context"
	"ipmanlk/cnapi/internal/model"
	"testing"
	"time"
)

func TestIslandScraper(t *testing.T) {
	// Create fetcher with HTTP client only for testing
	fetcher := createTestFetcher()

	scraper := NewIslandScraper(fetcher)

	// Test basic properties
	if scraper.Name() != "The Island" {
		t.Errorf("Expected name 'The Island', got '%s'", scraper.Name())
	}

	languages := scraper.Languages()
	expectedLangs := []model.Language{model.LangEn}
	if len(languages) != len(expectedLangs) {
		t.Errorf("Expected %d languages, got %d", len(expectedLangs), len(languages))
	}

	// Test English scraping
	t.Run("Scrape English", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		articles, err := scraper.Scrape(ctx, model.LangEn)
		if err != nil {
			t.Fatalf("Failed to scrape English articles: %v", err)
		}

		if len(articles) == 0 {
			t.Skip("No articles found (this might be normal if the site is down or has no content)")
		}

		// Validate first article
		article := articles[0]
		if article.SourceName != "The Island" {
			t.Errorf("Expected source name 'The Island', got '%s'", article.SourceName)
		}
		if article.Language != model.LangEn {
			t.Errorf("Expected language 'en', got '%s'", article.Language)
		}
		if article.Title == "" {
			t.Error("Article title is empty")
		}
		if article.URL == "" {
			t.Error("Article URL is empty")
		}
		if article.ContentText == "" {
			t.Error("Article content is empty")
		}
		if article.PublishedAt.IsZero() {
			t.Error("Article published date is zero")
		}

		t.Logf("Successfully scraped %d English articles", len(articles))
		t.Logf("First article: %s", article.Title)
	})
}
