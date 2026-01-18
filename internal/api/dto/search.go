package dto

import (
	"errors"
	"net/http"
	"time"

	"ipmanlk/cnapi/pkg/httpx"
)

var (
	ErrQueryRequired = errors.New("query parameter 'q' is required")
	ErrQueryTooShort = errors.New("query must be at least 2 characters")
)

// SearchFilterRequest represents search filtering parameters from HTTP request
type SearchFilterRequest struct {
	Query       string
	Language    *string
	SourceNames []string
	StartDate   *time.Time
	EndDate     *time.Time
}

// Validate validates search filter parameters
func (f *SearchFilterRequest) Validate() error {
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

// ParseSearchFilterRequest parses and validates search filter parameters from HTTP request
func ParseSearchFilterRequest(r *http.Request) (*SearchFilterRequest, error) {
	query := httpx.ParseQueryString(r, "q", "")

	startDate, err := httpx.ParseQueryTime(r, "start_date")
	if err != nil {
		return nil, err
	}

	endDate, err := httpx.ParseQueryTime(r, "end_date")
	if err != nil {
		return nil, err
	}

	req := &SearchFilterRequest{
		Query:       query,
		Language:    httpx.ParseQueryStringPtr(r, "language"),
		SourceNames: httpx.ParseQueryStrings(r, "source_names"),
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}
