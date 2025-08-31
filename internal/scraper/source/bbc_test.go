package source

import (
	"context"
	"ipmanlk/cnapi/internal/fetcher"
	"ipmanlk/cnapi/internal/model"
	"testing"
	"time"
)

func TestBBCScraper(t *testing.T) {
	// Create fetcher with HTTP client only for testing
	httpClient := fetcher.NewHTTPClient()
	browserClient := fetcher.NewBrowserAPIClient()
	fetcher := fetcher.NewFetcher(httpClient, browserClient)

	scraper := NewBBCScraper(fetcher)

	// Test basic properties
	if scraper.Name() != "BBC" {
		t.Errorf("Expected name 'BBC', got '%s'", scraper.Name())
	}

	languages := scraper.Languages()
	expectedLangs := []model.Language{model.LangSi, model.LangEn, model.LangTa}
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
		if article.SourceName != "BBC" {
			t.Errorf("Expected source name 'BBC', got '%s'", article.SourceName)
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

		t.Logf("Successfully scraped %d English articles", len(articles))
		t.Logf("First article: %s", article.Title)
	})

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
		if article.SourceName != "BBC" {
			t.Errorf("Expected source name 'BBC', got '%s'", article.SourceName)
		}
		if article.Language != model.LangSi {
			t.Errorf("Expected language 'si', got '%s'", article.Language)
		}

		t.Logf("Successfully scraped %d Sinhala articles", len(articles))
		t.Logf("First article: %s", article.Title)
	})

	// Test Tamil scraping
	t.Run("Scrape Tamil", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		articles, err := scraper.Scrape(ctx, model.LangTa)
		if err != nil {
			t.Fatalf("Failed to scrape Tamil articles: %v", err)
		}

		if len(articles) == 0 {
			t.Skip("No articles found (this might be normal if the site is down or has no content)")
		}

		// Validate first article
		article := articles[0]
		if article.SourceName != "BBC" {
			t.Errorf("Expected source name 'BBC', got '%s'", article.SourceName)
		}
		if article.Language != model.LangTa {
			t.Errorf("Expected language 'ta', got '%s'", article.Language)
		}

		t.Logf("Successfully scraped %d Tamil articles", len(articles))
		t.Logf("First article: %s", article.Title)
	})

}
