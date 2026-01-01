package service

import (
	"context"
	"fmt"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper"
	"log/slog"
	"sync"
)

type scrapeService struct {
	registry *scraper.Registry
}

// NewScrapeService creates a new scrape service
func NewScrapeService(registry *scraper.Registry) ScrapeService {
	return &scrapeService{
		registry: registry,
	}
}

// scrapeTask represents a single scraping task
type scrapeTask struct {
	scraper  scraper.SourceScraper
	language model.Language
}

// ScrapeAllConcurrent scrapes all sources concurrently with streaming results
// workerCount: number of concurrent scrapers
// batchSize: number of articles to batch before sending (0 = no batching)
func (s *scrapeService) ScrapeAllConcurrent(ctx context.Context, workerCount int, batchSize int) <-chan ScrapeResult {
	resultChan := make(chan ScrapeResult, workerCount)

	// Create task queue
	tasks := make([]scrapeTask, 0)
	scrapers := s.registry.GetScrapers()

	// Build all scraping tasks (scraper + language combinations)
	for _, scraper := range scrapers {
		for _, lang := range scraper.Languages() {
			tasks = append(tasks, scrapeTask{
				scraper:  scraper,
				language: lang,
			})
		}
	}

	// Start worker pool
	go func() {
		defer close(resultChan)

		taskChan := make(chan scrapeTask, len(tasks))
		var wg sync.WaitGroup

		// Start workers
		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go s.worker(ctx, taskChan, resultChan, batchSize, &wg)
		}

		// Send tasks to workers
		for _, task := range tasks {
			select {
			case taskChan <- task:
			case <-ctx.Done():
				close(taskChan)
				wg.Wait()
				return
			}
		}

		close(taskChan)
		wg.Wait()
	}()

	return resultChan
}

// worker processes scraping tasks from the task channel
func (s *scrapeService) worker(ctx context.Context, tasks <-chan scrapeTask, results chan<- ScrapeResult, batchSize int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
			// Scrape from this source
			articles, err := task.scraper.Scrape(ctx, task.language)

			result := ScrapeResult{
				Source:   task.scraper.Name(),
				Language: task.language,
				Success:  err == nil,
				Error:    err,
			}

			if err != nil {
				slog.Warn("failed to scrape from source",
					"source", task.scraper.Name(),
					"language", task.language,
					"error", err,
				)
				result.ArticleCount = 0
				result.Articles = nil
				results <- result
				continue
			}

			slog.Info("successfully scraped articles",
				"source", task.scraper.Name(),
				"language", task.language,
				"count", len(articles),
			)

			// Send results in batches to avoid overwhelming the consumer
			if batchSize <= 0 || len(articles) <= batchSize {
				// Send all at once if batch size is disabled or articles fit in one batch
				result.Articles = articles
				result.ArticleCount = len(articles)
				results <- result
			} else {
				// Send in batches
				for i := 0; i < len(articles); i += batchSize {
					end := i + batchSize
					if end > len(articles) {
						end = len(articles)
					}

					batchResult := ScrapeResult{
						Source:       task.scraper.Name(),
						Language:     task.language,
						Articles:     articles[i:end],
						ArticleCount: end - i,
						Success:      true,
						Error:        nil,
					}

					select {
					case results <- batchResult:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}
}

// ScrapeBySource scrapes articles from a specific source
func (s *scrapeService) ScrapeBySource(ctx context.Context, sourceName string) ([]model.ScrapedArticle, error) {
	scraper := s.registry.GetScraperByName(sourceName)
	if scraper == nil {
		return nil, fmt.Errorf("scraper not found: %s", sourceName)
	}

	var allArticles []model.ScrapedArticle
	for _, lang := range scraper.Languages() {
		articles, err := scraper.Scrape(ctx, lang)
		if err != nil {
			slog.Warn("failed to scrape from source",
				"source", sourceName,
				"language", lang,
				"error", err,
			)
			continue
		}

		allArticles = append(allArticles, articles...)
	}

	return allArticles, nil
}

// ScrapeByLanguage scrapes articles from all sources for a specific language
func (s *scrapeService) ScrapeByLanguage(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error) {
	var allArticles []model.ScrapedArticle
	scrapers := s.registry.GetScrapersByLanguage(string(language))

	for _, scraper := range scrapers {
		articles, err := scraper.Scrape(ctx, language)
		if err != nil {
			slog.Warn("failed to scrape from source",
				"source", scraper.Name(),
				"language", language,
				"error", err,
			)
			continue
		}

		allArticles = append(allArticles, articles...)
	}

	return allArticles, nil
}

// GetAvailableSources returns a list of all available source names
func (s *scrapeService) GetAvailableSources() []string {
	scrapers := s.registry.GetScrapers()
	sources := make([]string, len(scrapers))
	for i, scraper := range scrapers {
		sources[i] = scraper.Name()
	}
	return sources
}

// GetAvailableLanguages returns a list of all available languages
func (s *scrapeService) GetAvailableLanguages() []model.Language {
	scrapers := s.registry.GetScrapers()
	languageMap := make(map[model.Language]bool)

	for _, scraper := range scrapers {
		for _, lang := range scraper.Languages() {
			languageMap[lang] = true
		}
	}

	languages := make([]model.Language, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}

	return languages
}
