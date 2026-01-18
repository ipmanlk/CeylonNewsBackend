package service

import (
	"context"

	"ipmanlk/cnapi/internal/database"
	"ipmanlk/cnapi/internal/model"
)

type searchService struct {
	store database.SearchStore
}

// NewSearchService creates a new search service
func NewSearchService(store database.SearchStore) SearchService {
	return &searchService{
		store: store,
	}
}

func (s *searchService) Search(ctx context.Context, filter model.SearchFilter) (*model.PaginatedResult[*model.SearchResult], error) {
	return s.store.Search(filter)
}

func (s *searchService) GetAvailableSources() ([]string, error) {
	return s.store.GetAvailableSources()
}

func (s *searchService) GetAvailableLanguages() ([]string, error) {
	return s.store.GetAvailableLanguages()
}

func (s *searchService) GetSourcesByLanguage(language string) ([]string, error) {
	return s.store.GetSourcesByLanguage(language)
}

func (s *searchService) GetRecentArticles(language *string, sourceNames []string, limit int) ([]*model.Article, error) {
	return s.store.GetRecentArticles(language, sourceNames, limit)
}
