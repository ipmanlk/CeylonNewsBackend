package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	MinLimit     = 1
	MaxLimit     = 100
	DefaultLimit = 20
	MinOffset    = 0
)

var (
	ErrLimitTooSmall    = fmt.Errorf("limit must be at least %d", MinLimit)
	ErrLimitTooLarge    = fmt.Errorf("limit cannot exceed %d", MaxLimit)
	ErrOffsetNegative   = errors.New("offset cannot be negative")
	ErrInvalidLanguage  = errors.New("language must be one of: en, si, ta")
	ErrInvalidDateRange = errors.New("start_date must be before end_date")
	ErrQueryRequired    = errors.New("query parameter 'q' is required")
	ErrQueryTooShort    = errors.New("query must be at least 2 characters")
)

type PaginationParams struct {
	Limit  int
	Offset int
}

func (p *PaginationParams) Validate() error {
	if p.Limit < MinLimit {
		return ErrLimitTooSmall
	}
	if p.Limit > MaxLimit {
		return ErrLimitTooLarge
	}
	if p.Offset < MinOffset {
		return ErrOffsetNegative
	}
	return nil
}

type ArticleFilterParams struct {
	Language    *string
	SourceNames []string
	StartDate   *time.Time
	EndDate     *time.Time
}

func (f *ArticleFilterParams) Validate() error {
	if f.Language != nil {
		if !isValidLanguage(*f.Language) {
			return ErrInvalidLanguage
		}
	}

	if f.StartDate != nil && f.EndDate != nil {
		if f.StartDate.After(*f.EndDate) {
			return ErrInvalidDateRange
		}
	}

	return nil
}

type SearchFilterParams struct {
	Query       string
	Language    *string
	SourceNames []string
	StartDate   *time.Time
	EndDate     *time.Time
}

func (f *SearchFilterParams) Validate() error {
	if f.Query == "" {
		return ErrQueryRequired
	}

	if len(f.Query) < 2 {
		return ErrQueryTooShort
	}

	if f.Language != nil {
		if !isValidLanguage(*f.Language) {
			return ErrInvalidLanguage
		}
	}

	if f.StartDate != nil && f.EndDate != nil {
		if f.StartDate.After(*f.EndDate) {
			return ErrInvalidDateRange
		}
	}

	return nil
}

func isValidLanguage(lang string) bool {
	return lang == "en" || lang == "si" || lang == "ta"
}

func ParsePaginationParams(r *http.Request) (*PaginationParams, error) {
	limit, err := ParseQueryInt(r, "limit", DefaultLimit)
	if err != nil {
		return nil, err
	}

	offset, err := ParseQueryInt(r, "offset", 0)
	if err != nil {
		return nil, err
	}

	params := &PaginationParams{
		Limit:  limit,
		Offset: offset,
	}

	if err := params.Validate(); err != nil {
		return nil, err
	}

	return params, nil
}

func ParseArticleFilterParams(r *http.Request) (*ArticleFilterParams, error) {
	startDate, err := ParseQueryTime(r, "start_date")
	if err != nil {
		return nil, err
	}

	endDate, err := ParseQueryTime(r, "end_date")
	if err != nil {
		return nil, err
	}

	params := &ArticleFilterParams{
		Language:    ParseQueryStringPtr(r, "language"),
		SourceNames: ParseQueryStrings(r, "source_names"),
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := params.Validate(); err != nil {
		return nil, err
	}

	return params, nil
}

func ParseSearchFilterParams(r *http.Request) (*SearchFilterParams, error) {
	query := ParseQueryString(r, "q", "")

	startDate, err := ParseQueryTime(r, "start_date")
	if err != nil {
		return nil, err
	}

	endDate, err := ParseQueryTime(r, "end_date")
	if err != nil {
		return nil, err
	}

	params := &SearchFilterParams{
		Query:       query,
		Language:    ParseQueryStringPtr(r, "language"),
		SourceNames: ParseQueryStrings(r, "source_names"),
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := params.Validate(); err != nil {
		return nil, err
	}

	return params, nil
}
