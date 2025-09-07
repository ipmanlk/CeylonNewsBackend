package store

import (
	"database/sql"
	"fmt"
	"strings"

	"ipmanlk/cnapi/internal/model"
)

type SearchStore struct {
	db *sql.DB
}

func NewSearchStore(db *sql.DB) *SearchStore {
	return &SearchStore{db: db}
}

// Search performs a full-text search on articles with optional filtering
func (s *SearchStore) Search(filter model.SearchFilter) ([]*model.SearchResult, error) {
	query, args := s.buildSearchQuery(filter)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search articles: %w", err)
	}
	defer rows.Close()

	var results []*model.SearchResult
	for rows.Next() {
		result := &model.SearchResult{}
		err := rows.Scan(
			&result.ID,
			&result.SourceName,
			&result.Title,
			&result.URL,
			&result.Content,
			&result.ImageURL,
			&result.Language,
			&result.PublishedAt,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.RelevanceScore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search rows: %w", err)
	}

	return results, nil
}

// SearchPaginated performs a paginated full-text search on articles
func (s *SearchStore) SearchPaginated(filter model.SearchFilter) (*model.PaginatedResult[*model.SearchResult], error) {
	// Get total count
	total, err := s.CountSearchResults(filter)
	if err != nil {
		return nil, err
	}

	// Get search results
	results, err := s.Search(filter)
	if err != nil {
		return nil, err
	}

	page := (filter.Offset / filter.Limit) + 1
	if filter.Limit == 0 {
		page = 1
	}

	return model.NewPaginatedResult(results, total, page, filter.Limit), nil
}

// CountSearchResults returns the total number of search results matching the filter
func (s *SearchStore) CountSearchResults(filter model.SearchFilter) (int64, error) {
	query, args := s.buildSearchCountQuery(filter)

	var count int64
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count search results: %w", err)
	}

	return count, nil
}

// GetAvailableSources returns all unique source names in the database
func (s *SearchStore) GetAvailableSources() ([]string, error) {
	query := `SELECT DISTINCT source_name FROM articles ORDER BY source_name`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get available sources: %w", err)
	}
	defer rows.Close()

	var sources []string
	for rows.Next() {
		var source string
		err := rows.Scan(&source)
		if err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating source rows: %w", err)
	}

	return sources, nil
}

// GetAvailableLanguages returns all unique languages in the database
func (s *SearchStore) GetAvailableLanguages() ([]string, error) {
	query := `SELECT DISTINCT language FROM articles ORDER BY language`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get available languages: %w", err)
	}
	defer rows.Close()

	var languages []string
	for rows.Next() {
		var language string
		err := rows.Scan(&language)
		if err != nil {
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}
		languages = append(languages, language)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating language rows: %w", err)
	}

	return languages, nil
}

// GetSourcesByLanguage returns all sources that have articles in the specified language
func (s *SearchStore) GetSourcesByLanguage(language string) ([]string, error) {
	query := `SELECT DISTINCT source_name FROM articles WHERE language = ? ORDER BY source_name`

	rows, err := s.db.Query(query, language)
	if err != nil {
		return nil, fmt.Errorf("failed to get sources by language: %w", err)
	}
	defer rows.Close()

	var sources []string
	for rows.Next() {
		var source string
		err := rows.Scan(&source)
		if err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating source rows: %w", err)
	}

	return sources, nil
}

// buildSearchQuery builds the SQL query for full-text search
func (s *SearchStore) buildSearchQuery(filter model.SearchFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// FTS search condition
	if filter.Query != "" {
		// Escape special characters in the query for FTS
		escapedQuery := s.escapeFTSQuery(filter.Query)
		conditions = append(conditions, "articles_fts MATCH ?")
		args = append(args, escapedQuery)
	}

	// Additional filters
	if filter.Language != nil {
		conditions = append(conditions, "a.language = ?")
		args = append(args, *filter.Language)
	}

	if len(filter.SourceNames) > 0 {
		placeholders := make([]string, len(filter.SourceNames))
		for i, source := range filter.SourceNames {
			placeholders[i] = "?"
			args = append(args, source)
		}
		conditions = append(conditions, fmt.Sprintf("a.source_name IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.StartDate != nil {
		conditions = append(conditions, "a.published_at >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		conditions = append(conditions, "a.published_at <= ?")
		args = append(args, *filter.EndDate)
	}

	query := `
		SELECT a.id, a.source_name, a.title, a.url, a.content, a.image_url, a.language, 
		       a.published_at, a.created_at, a.updated_at,
		       articles_fts.rank
		FROM articles a
		JOIN articles_fts ON a.id = articles_fts.rowid
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Order by relevance score (rank) for FTS, then by published date
	if filter.Query != "" {
		query += " ORDER BY articles_fts.rank, a.published_at DESC"
	} else {
		query += " ORDER BY a.published_at DESC"
	}

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filter.Offset)
	}

	return query, args
}

// buildSearchCountQuery builds the SQL query for counting search results
func (s *SearchStore) buildSearchCountQuery(filter model.SearchFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// FTS search condition
	if filter.Query != "" {
		escapedQuery := s.escapeFTSQuery(filter.Query)
		conditions = append(conditions, "articles_fts MATCH ?")
		args = append(args, escapedQuery)
	}

	// Additional filters
	if filter.Language != nil {
		conditions = append(conditions, "a.language = ?")
		args = append(args, *filter.Language)
	}

	if len(filter.SourceNames) > 0 {
		placeholders := make([]string, len(filter.SourceNames))
		for i, source := range filter.SourceNames {
			placeholders[i] = "?"
			args = append(args, source)
		}
		conditions = append(conditions, fmt.Sprintf("a.source_name IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.StartDate != nil {
		conditions = append(conditions, "a.published_at >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		conditions = append(conditions, "a.published_at <= ?")
		args = append(args, *filter.EndDate)
	}

	query := `
		SELECT COUNT(*)
		FROM articles a
		JOIN articles_fts ON a.id = articles_fts.rowid
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}

// escapeFTSQuery escapes special characters in the query for FTS5
func (s *SearchStore) escapeFTSQuery(query string) string {
	// FTS5 special characters that need to be escaped: " ' * + - : ^ ~ ( ) [ ] { } ,
	// We'll wrap the entire query in quotes to treat it as a phrase search
	// and escape any internal quotes
	escaped := strings.ReplaceAll(query, `"`, `""`)
	return `"` + escaped + `"`
}

// SearchWithHighlight performs a search and returns results with highlighted snippets
func (s *SearchStore) SearchWithHighlight(filter model.SearchFilter) ([]*model.SearchResult, error) {
	// For now, we'll use the regular search. In the future, we could implement
	// snippet highlighting using FTS5's snippet() function
	return s.Search(filter)
}

// GetRecentArticles returns the most recent articles with optional filtering
func (s *SearchStore) GetRecentArticles(language *string, sourceNames []string, limit int) ([]*model.Article, error) {
	var conditions []string
	var args []interface{}

	if language != nil {
		conditions = append(conditions, "language = ?")
		args = append(args, *language)
	}

	if len(sourceNames) > 0 {
		placeholders := make([]string, len(sourceNames))
		for i, source := range sourceNames {
			placeholders[i] = "?"
			args = append(args, source)
		}
		conditions = append(conditions, fmt.Sprintf("source_name IN (%s)", strings.Join(placeholders, ",")))
	}

	query := `
		SELECT id, source_name, title, url, content, image_url, language, published_at, created_at, updated_at
		FROM articles
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY published_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent articles: %w", err)
	}
	defer rows.Close()

	var articles []*model.Article
	for rows.Next() {
		article := &model.Article{}
		err := rows.Scan(
			&article.ID,
			&article.SourceName,
			&article.Title,
			&article.URL,
			&article.Content,
			&article.ImageURL,
			&article.Language,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating article rows: %w", err)
	}

	return articles, nil
}
