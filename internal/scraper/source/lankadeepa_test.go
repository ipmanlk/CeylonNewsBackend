package source

import (
	"context"
	
	"ipmanlk/cnapi/internal/model"
	"testing"
	"time"
)

func TestLankadeepaScraper(t *testing.T) {
	// Create fetcher with HTTP client only for testing
	fetcher := createTestFetcher()

	scraper := NewLankadeepaScraper(fetcher)

	// Test basic properties
	if scraper.Name() != "Lankadeepa" {
		t.Errorf("Expected name 'Lankadeepa', got '%s'", scraper.Name())
	}

	languages := scraper.Languages()
	expectedLangs := []model.Language{model.LangSi}
	if len(languages) != len(expectedLangs) {
		t.Errorf("Expected %d languages, got %d", len(expectedLangs), len(languages))
	}

	// Test Sinhala scraping
	t.Run("Scrape Sinhala", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		articles, err := scraper.Scrape(ctx, model.LangSi)
		if err != nil {
			t.Fatalf("Failed to scrape Sinhala articles: %v", err)
		}

		if len(articles) == 0 {
			t.Skip("No articles found (this might be normal if the site is down or has no content)")
		}

		// Validate first article
		article := articles[0]
		if article.SourceName != "Lankadeepa" {
			t.Errorf("Expected source name 'Lankadeepa', got '%s'", article.SourceName)
		}
		if article.Language != model.LangSi {
			t.Errorf("Expected language 'si', got '%s'", article.Language)
		}
		if article.Title == "" {
			t.Error("Article title is empty")
		}
		if article.URL == "" {
			t.Error("Article URL is empty")
		}
		if article.Content == "" {
			t.Error("Article content is empty")
		}
		if article.PublishedAt.IsZero() {
			t.Error("Article published date is zero")
		}

		t.Logf("Successfully scraped %d Sinhala articles", len(articles))
		t.Logf("First article: %s", article.Title)
	})
}
