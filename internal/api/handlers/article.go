package handlers

import (
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/service"
	"ipmanlk/cnapi/pkg/httpx"
	"log/slog"
	"net/http"
)

type ArticleResponse struct {
	ID          int64   `json:"id"`
	SourceName  string  `json:"source_name"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	ContentHTML string  `json:"content_html,omitempty"`
	ContentText *string `json:"content_text,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	Language    string  `json:"language"`
	PublishedAt string  `json:"published_at"`
}

type ArticleHandler struct {
	articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

func (h *ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
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

	includeText := r.URL.Query().Get("include_text") == "true"

	filter := model.ArticleFilter{
		Language:    language,
		SourceNames: sourceNames,
		StartDate:   startDate,
		EndDate:     endDate,
		Limit:       limit,
		Offset:      offset,
		IncludeText: includeText,
	}

	result, err := h.articleService.ListPaginated(r.Context(), filter)
	if err != nil {
		slog.Error("failed to list articles", "error", err)
		httpx.RespondInternalError(w, "failed to retrieve articles")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, result)
}

func (h *ArticleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParsePathInt64(r, "id")
	if err != nil {
		httpx.RespondBadRequest(w, err.Error())
		return
	}

	includeText := r.URL.Query().Get("include_text") == "true"
	filter := model.ArticleFilter{IncludeText: includeText}

	article, err := h.articleService.GetByIDWithFilter(r.Context(), id, filter)
	if err != nil {
		slog.Error("failed to get article", "id", id, "error", err)
		httpx.RespondInternalError(w, "failed to retrieve article")
		return
	}

	if article == nil {
		httpx.RespondNotFound(w, "article not found")
		return
	}

	response := ArticleResponse{
		ID:          article.ID,
		SourceName:  article.SourceName,
		Title:       article.Title,
		URL:         article.URL,
		ContentHTML: article.ContentHTML,
		ImageURL:    article.ImageURL,
		Language:    article.Language,
		PublishedAt: article.PublishedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if includeText {
		response.ContentText = &article.ContentText
	}

	httpx.RespondJSON(w, http.StatusOK, response)
}
