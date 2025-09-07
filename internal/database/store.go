package database

import (
	"database/sql"
	"ipmanlk/cnapi/internal/database/store"
	"ipmanlk/cnapi/internal/model"
)

// ArticlesStoreInterface defines the interface for article operations
type ArticlesStoreInterface interface {
	Create(scrapedArticle model.ScrapedArticle) (int64, error)
	Upsert(scrapedArticle model.ScrapedArticle) (int64, error)
	BulkCreate(scrapedArticles []model.ScrapedArticle) ([]int64, error)
	BulkUpsert(scrapedArticles []model.ScrapedArticle) ([]int64, error)
	GetByID(id int64) (*model.Article, error)
	GetByURL(url string) (*model.Article, error)
	List(filter model.ArticleFilter) ([]*model.Article, error)
	Count(filter model.ArticleFilter) (int64, error)
	ListPaginated(filter model.ArticleFilter) (*model.PaginatedResult[*model.Article], error)
	Update(article *model.Article) error
	Delete(id int64) error
	ExistsByURL(url string) (bool, error)
}

// SearchStoreInterface defines the interface for search operations
type SearchStoreInterface interface {
	Search(filter model.SearchFilter) ([]*model.SearchResult, error)
	SearchPaginated(filter model.SearchFilter) (*model.PaginatedResult[*model.SearchResult], error)
	CountSearchResults(filter model.SearchFilter) (int64, error)
	GetAvailableSources() ([]string, error)
	GetAvailableLanguages() ([]string, error)
	GetSourcesByLanguage(language string) ([]string, error)
	SearchWithHighlight(filter model.SearchFilter) ([]*model.SearchResult, error)
	GetRecentArticles(language *string, sourceNames []string, limit int) ([]*model.Article, error)
}

// Store provides access to all database operations
type Store struct {
	Articles ArticlesStoreInterface
	Search   SearchStoreInterface
}

// NewStore creates a new database store with all required stores
func NewStore(db *sql.DB) *Store {
	return &Store{
		Articles: store.NewArticlesStore(db),
		Search:   store.NewSearchStore(db),
	}
}
