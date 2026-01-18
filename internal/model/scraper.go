package model

import "time"

type ScrapedArticle struct {
	SourceName  string
	Title       string
	URL         string
	ContentText string
	ContentHTML string
	ImageURL    *string
	Categories  []string
	Language    Language
	PublishedAt time.Time
}

// Article represents a news article stored in the database
type Article struct {
	ID          int64     `json:"id" db:"id"`
	SourceName  string    `json:"source_name" db:"source_name"`
	Title       string    `json:"title" db:"title"`
	URL         string    `json:"url" db:"url"`
	ContentText string    `json:"content_text" db:"content_text"`
	ContentHTML string    `json:"content_html,omitempty" db:"content_html"`
	ImageURL    *string   `json:"image_url,omitempty" db:"image_url"`
	Language    string    `json:"language" db:"language"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ArticleFilter represents filters for querying articles
type ArticleFilter struct {
	Language    *string    `json:"language,omitempty"`
	SourceNames []string   `json:"source_names,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

// SearchFilter represents filters for searching articles
type SearchFilter struct {
	Query       string     `json:"query"`
	Language    *string    `json:"language,omitempty"`
	SourceNames []string   `json:"source_names,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

// SearchResult represents a search result with relevance score
type SearchResult struct {
	Article
	RelevanceScore float64 `json:"relevance_score" db:"rank"`
}

// PaginatedResult represents a paginated result set
type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// NewPaginatedResult creates a new paginated result
func NewPaginatedResult[T any](data []T, total int64, page, perPage int) *PaginatedResult[T] {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return &PaginatedResult[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
