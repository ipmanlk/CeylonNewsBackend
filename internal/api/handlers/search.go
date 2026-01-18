package handlers

import (
	"log/slog"
	"net/http"

	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/service"
	"ipmanlk/cnapi/pkg/httpx"
)

type SearchHandler struct {
	searchService service.SearchService
}

type SearchResultResponse struct {
	ID             int64   `json:"id"`
	SourceName     string  `json:"source_name"`
	Title          string  `json:"title"`
	URL            string  `json:"url"`
	ImageURL       *string `json:"image_url,omitempty"`
	Language       string  `json:"language"`
	PublishedAt    string  `json:"published_at"`
	RelevanceScore float64 `json:"relevance_score"`
}

func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := httpx.ParseQueryString(r, "q", "")
	if query == "" {
		httpx.RespondBadRequest(w, "query parameter 'q' is required")
		return
	}

	limit, err := httpx.ParseQueryInt(r, "limit", 20)
	if err != nil {
		httpx.RespondBadRequest(w, err.Error())
		return
	}

	offset, err := httpx.ParseQueryInt(r, "offset", 0)
	if err != nil {
		httpx.RespondBadRequest(w, err.Error())
		return
	}

	language := httpx.ParseQueryStringPtr(r, "language")
	sourceNames := httpx.ParseQueryStrings(r, "source_names")

	startDate, err := httpx.ParseQueryTime(r, "start_date")
	if err != nil {
		httpx.RespondBadRequest(w, err.Error())
		return
	}

	endDate, err := httpx.ParseQueryTime(r, "end_date")
	if err != nil {
		httpx.RespondBadRequest(w, err.Error())
		return
	}

	filter := model.SearchFilter{
		Query:       query,
		Language:    language,
		SourceNames: sourceNames,
		StartDate:   startDate,
		EndDate:     endDate,
		Limit:       limit,
		Offset:      offset,
	}

	paginatedResult, err := h.searchService.SearchPaginated(r.Context(), filter)
	if err != nil {
		slog.Error("failed to search articles", "query", query, "error", err)
		httpx.RespondInternalError(w, "failed to search articles")
		return
	}

	// Convert to lightweight response format
	searchResponses := make([]SearchResultResponse, len(paginatedResult.Data))
	for i, article := range paginatedResult.Data {
		searchResponses[i] = SearchResultResponse{
			ID:             article.ID,
			SourceName:     article.SourceName,
			Title:          article.Title,
			URL:            article.URL,
			ImageURL:       article.ImageURL,
			Language:       article.Language,
			PublishedAt:    article.PublishedAt.Format("2006-01-02T15:04:05Z07:00"),
			RelevanceScore: article.RelevanceScore,
		}
	}

	response := struct {
		Data       []SearchResultResponse `json:"data"`
		Total      int64                  `json:"total"`
		Page       int                    `json:"page"`
		PerPage    int                    `json:"per_page"`
		TotalPages int                    `json:"total_pages"`
	}{
		Data:       searchResponses,
		Total:      paginatedResult.Total,
		Page:       paginatedResult.Page,
		PerPage:    paginatedResult.PerPage,
		TotalPages: paginatedResult.TotalPages,
	}

	httpx.RespondJSON(w, http.StatusOK, response)
}

func (h *SearchHandler) GetAvailableSources(w http.ResponseWriter, r *http.Request) {
	sources, err := h.searchService.GetAvailableSources()
	if err != nil {
		slog.Error("failed to get available sources", "error", err)
		httpx.RespondInternalError(w, "failed to retrieve sources")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"sources": sources,
	})
}

func (h *SearchHandler) GetAvailableLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := h.searchService.GetAvailableLanguages()
	if err != nil {
		slog.Error("failed to get available languages", "error", err)
		httpx.RespondInternalError(w, "failed to retrieve languages")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"languages": languages,
	})
}

func (h *SearchHandler) GetSourcesByLanguage(w http.ResponseWriter, r *http.Request) {
	language := httpx.ParseQueryString(r, "language", "")
	if language == "" {
		httpx.RespondBadRequest(w, "query parameter 'language' is required")
		return
	}

	sources, err := h.searchService.GetSourcesByLanguage(language)
	if err != nil {
		slog.Error("failed to get sources by language", "language", language, "error", err)
		httpx.RespondInternalError(w, "failed to retrieve sources")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"language": language,
		"sources":  sources,
	})
}

func (h *SearchHandler) GetRecentArticles(w http.ResponseWriter, r *http.Request) {
	language := httpx.ParseQueryStringPtr(r, "language")
	sourceNames := httpx.ParseQueryStrings(r, "sourceNames")

	limit, err := httpx.ParseQueryInt(r, "limit", 20)
	if err != nil {
		httpx.RespondBadRequest(w, err.Error())
		return
	}

	articles, err := h.searchService.GetRecentArticles(language, sourceNames, limit)
	if err != nil {
		slog.Error("failed to get recent articles", "error", err)
		httpx.RespondInternalError(w, "failed to retrieve recent articles")
		return
	}

	// Convert to lightweight response format
	searchResponses := make([]SearchResultResponse, len(articles))
	for i, article := range articles {
		searchResponses[i] = SearchResultResponse{
			ID:             article.ID,
			SourceName:     article.SourceName,
			Title:          article.Title,
			URL:            article.URL,
			ImageURL:       article.ImageURL,
			Language:       article.Language,
			PublishedAt:    article.PublishedAt.Format("2006-01-02T15:04:05Z07:00"),
			RelevanceScore: 0, // No relevance score for recent articles
		}
	}

	httpx.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"articles": searchResponses,
		"count":    len(searchResponses),
	})
}
