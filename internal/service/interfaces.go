package service

import (
	"context"

	"ipmanlk/cnapi/internal/model"
)

// ScrapeResult represents the result of a scraping operation
type ScrapeResult struct {
	Articles     []model.ScrapedArticle
	Source       string
	Language     model.Language
	ArticleCount int
	Success      bool
	Error        error
}

// ScrapeService defines the interface for scraping operations
type ScrapeService interface {
	// ScrapeAllConcurrent scrapes all sources concurrently and streams results
	ScrapeAllConcurrent(ctx context.Context, workerCount int, batchSize int) <-chan ScrapeResult

	ScrapeBySource(ctx context.Context, sourceName string) ([]model.ScrapedArticle, error)
	ScrapeByLanguage(ctx context.Context, language model.Language) ([]model.ScrapedArticle, error)
	GetAvailableSources() []string
	GetAvailableLanguages() []model.Language
}

// ArticleService defines the interface for article operations
type ArticleService interface {
	Create(ctx context.Context, article model.ScrapedArticle) (int64, error)
	BulkCreate(ctx context.Context, articles []model.ScrapedArticle) ([]int64, error)
	BulkUpsert(ctx context.Context, articles []model.ScrapedArticle) ([]int64, error)
	GetByID(ctx context.Context, id int64) (*model.Article, error)
	GetByIDWithFilter(ctx context.Context, id int64, filter model.ArticleFilter) (*model.Article, error)
	GetByURL(ctx context.Context, url string) (*model.Article, error)
	List(ctx context.Context, filter model.ArticleFilter) ([]*model.Article, error)
	ListPaginated(ctx context.Context, filter model.ArticleFilter) (*model.PaginatedResult[*model.Article], error)
	Update(ctx context.Context, article *model.Article) error
	Delete(ctx context.Context, id int64) error
	ExistsByURL(ctx context.Context, url string) (bool, error)
}

// SearchService defines the interface for search operations
type SearchService interface {
	Search(ctx context.Context, filter model.SearchFilter) (*model.PaginatedResult[*model.SearchResult], error)
	GetAvailableSources() ([]string, error)
	GetAvailableLanguages() ([]string, error)
	GetSourcesByLanguage(language string) ([]string, error)
	GetRecentArticles(language *string, sourceNames []string, limit int) ([]*model.Article, error)
}
