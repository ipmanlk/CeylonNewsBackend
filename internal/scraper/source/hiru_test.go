package source

import (
	"context"
	
	"ipmanlk/cnapi/internal/model"
	"testing"
	"time"
)

func TestHiruScraper(t *testing.T) {
	// Create fetcher with HTTP client only for testing
	fetcher := createTestFetcher()

	scraper := NewHiruScraper(fetcher)

	// Test basic properties
	if scraper.Name() != "Hiru News" {
		t.Errorf("Expected name 'Hiru News', got '%s'", scraper.Name())
	}

	languages := scraper.Languages()
	expectedLangs := []model.Language{model.LangSi, model.LangEn, model.LangTa}
	if len(languages) != len(expectedLangs) {
		t.Errorf("Expected %d languages, got %d", len(expectedLangs), len(languages))
	}
	for i, lang := range expectedLangs {
		if languages[i] != lang {
			t.Errorf("Expected language '%s', got '%s'", lang, languages[i])
		}
	}

	// Test Sinhala scraping
	t.Run("Scrape Sinhala", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		articles, err := scraper.Scrape(ctx, model.LangSi)
		if err != nil {
			t.Fatalf("Failed to scrape Sinhala articles: %v", err)
		}

		if len(articles) == 0 {
			t.Skip("No articles found (this might be normal if the website is down or has no content)")
		}

		// Validate first article
		article := articles[0]
		if article.SourceName != "Hiru News" {
			t.Errorf("Expected source name 'Hiru News', got '%s'", article.SourceName)
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
		if article.ContentText == "" {
			t.Error("Article content is empty")
		}
		if article.PublishedAt.IsZero() {
			t.Error("Article published date is zero")
		}

		t.Logf("Successfully scraped %d Sinhala articles", len(articles))
		t.Logf("First article: %s", article.Title)
		t.Logf("First article URL: %s", article.URL)
	})

	// Test English scraping
	t.Run("Scrape English", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		articles, err := scraper.Scrape(ctx, model.LangEn)
		if err != nil {
			t.Fatalf("Failed to scrape English articles: %v", err)
		}

		if len(articles) == 0 {
			t.Skip("No articles found (this might be normal if the website is down or has no content)")
		}

		// Validate first article
		article := articles[0]
		if article.SourceName != "Hiru News" {
			t.Errorf("Expected source name 'Hiru News', got '%s'", article.SourceName)
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
		t.Logf("First article URL: %s", article.URL)
	})

	// Test Tamil scraping
	t.Run("Scrape Tamil", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		articles, err := scraper.Scrape(ctx, model.LangTa)
		if err != nil {
			t.Fatalf("Failed to scrape Tamil articles: %v", err)
		}

		if len(articles) == 0 {
			t.Skip("No articles found (this might be normal if the website is down or has no content)")
		}

		// Validate first article
		article := articles[0]
		if article.SourceName != "Hiru News" {
			t.Errorf("Expected source name 'Hiru News', got '%s'", article.SourceName)
		}
		if article.Language != model.LangTa {
			t.Errorf("Expected language 'ta', got '%s'", article.Language)
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

		t.Logf("Successfully scraped %d Tamil articles", len(articles))
		t.Logf("First article: %s", article.Title)
		t.Logf("First article URL: %s", article.URL)
	})

}
