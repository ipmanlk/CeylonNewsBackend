package handlers

import (
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/service"
	"ipmanlk/cnapi/pkg/httpx"
	"log/slog"
	"net/http"
)

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
	sourceName := httpx.ParseQueryStringPtr(r, "source_name")

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

	filter := model.ArticleFilter{
		Language:   language,
		SourceName: sourceName,
		StartDate:  startDate,
		EndDate:    endDate,
		Limit:      limit,
		Offset:     offset,
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

	article, err := h.articleService.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get article", "id", id, "error", err)
		httpx.RespondInternalError(w, "failed to retrieve article")
		return
	}

	if article == nil {
		httpx.RespondNotFound(w, "article not found")
		return
	}

	httpx.RespondJSON(w, http.StatusOK, article)
}
