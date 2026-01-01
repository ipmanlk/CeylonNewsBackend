package source

import (
	"context"
	
	"ipmanlk/cnapi/internal/model"
	"strings"
	"testing"
	"time"
)

func TestDailyMirrorScraper(t *testing.T) {
	// Create fetcher with HTTP client only for testing
	fetcher := createTestFetcher()

	scraper := NewDailyMirrorScraper(fetcher)

	// Test basic properties
	if scraper.Name() != "Daily Mirror" {
		t.Errorf("Expected name 'Daily Mirror', got '%s'", scraper.Name())
	}

	languages := scraper.Languages()
	expectedLangs := []model.Language{model.LangEn}
	if len(languages) != len(expectedLangs) {
		t.Errorf("Expected %d languages, got %d", len(expectedLangs), len(languages))
	}
	if languages[0] != model.LangEn {
		t.Errorf("Expected language 'en', got '%s'", languages[0])
	}

	// Test English scraping
	t.Run("Scrape English", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		articles, err := scraper.Scrape(ctx, model.LangEn)
		if err != nil {
			t.Fatalf("Failed to scrape English articles: %v", err)
		}

		if len(articles) == 0 {
			t.Skip("No articles found (this might be normal if the RSS feed is down or has no content)")
		}

		// Validate first article
		article := articles[0]
		if article.SourceName != "Daily Mirror" {
			t.Errorf("Expected source name 'Daily Mirror', got '%s'", article.SourceName)
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
		if article.Content == "" {
			t.Error("Article content is empty")
		}
		if article.PublishedAt.IsZero() {
			t.Error("Article published date is zero")
		}

		// Check that title doesn't contain the suffix we're supposed to remove
		if strings.Contains(article.Title, " - Breaking News | Daily Mirror") {
			t.Error("Article title still contains the suffix that should have been removed")
		}

		t.Logf("Successfully scraped %d English articles", len(articles))
		t.Logf("First article: %s", article.Title)
		t.Logf("First article URL: %s", article.URL)
	})

}
