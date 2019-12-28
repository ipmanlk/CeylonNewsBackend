package api

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/providers"
	"net/http"
	"strconv"
	"strings"
)

func validateGetSources(r *http.Request) ([]common.Lang, *common.Err) {
	if r.Method != "GET" {
		return nil, common.ErrInvalidRequestMethod
	}

	langs, err := getLangs(r.URL.Query().Get("languages"))
	if err != nil {
		return nil, err
	}

	if len(langs) == 0 {
		return nil, &common.Err{
			Code:    400,
			Message: "missing languages query string parameter",
		}
	}

	return langs, nil
}

func validateHandleGetNews(r *http.Request) (*getNewsData, *common.Err) {
	if r.Method != "GET" {
		return nil, common.ErrInvalidRequestMethod
	}

	langs, err := getLangs(r.URL.Query().Get("languages"))
	if err != nil {
		return nil, err
	}

	sources, err := getSources(r.URL.Query().Get("sources"), langs)
	if err != nil {
		return nil, err
	}

	// search keyword/query
	query := r.URL.Query().Get("query")

	// get pagination data
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 20

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, &common.Err{
				Code:    400,
				Message: "invalid page query string parameter",
			}
		}
		page = p
	}

	if pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, &common.Err{
				Code:    400,
				Message: "invalid pageSize query string parameter",
			}
		}

		if ps > 50 {
			return nil, &common.Err{
				Code:    400,
				Message: "pageSize cannot be greater than 50",
			}
		}

		pageSize = ps
	}

	// Start: Compatibility with old API
	// Source: https://github.com/ipmanlk/CeylonNewsBackend/blob/rewrite/src/api/v1.0/validators/get-posts.validator.ts
	skipStr := r.URL.Query().Get("skip")
	if skipStr != "" {
		skip, err := strconv.Atoi(skipStr)
		if err != nil {
			return nil, &common.Err{
				Code:    400,
				Message: "invalid skip query string parameter",
			}
		}
		// limit will always be 20 in the old api
		// calculate page from the provided skip
		page = skip/20 + 1
	}

	keyword := r.URL.Query().Get("keyword")
	if keyword != "" {
		query = keyword
	}
	// End: Compatibility with old API

	return &getNewsData{
		Langs:    langs,
		Sources:  sources,
		Page:     page,
		PageSize: pageSize,
		Query:    query,
	}, nil
}

func validateHandleGetNewsItem(r *http.Request) (uint, *common.Err) {
	if r.Method != "GET" {
		return 0, common.ErrInvalidRequestMethod
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1.0/news/")

	if idStr == "" {
		return 0, &common.Err{
			Code:    400,
			Message: "missing news item id",
		}
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, &common.Err{
			Code:    400,
			Message: "invalid news item id",
		}
	}

	return uint(id), nil
}

func getSources(sourcesStr string, langs []common.Lang) ([]string, *common.Err) {
	if sourcesStr == "" {
		return providers.ActiveProviderNames, nil
	}

	sourceStrs := strings.Split(sourcesStr, ",")

	if len(sourceStrs) > providers.ActiveProviderCount {
		return nil, &common.Err{
			Code:    400,
			Message: "too many sources",
		}
	}

	addedSources := make(map[string]struct{})
	sources := make([]string, 0)

	for _, s := range sourceStrs {
		source := strings.TrimSpace(s)

		if _, ok := addedSources[source]; ok {
			continue
		}

		for _, lang := range langs {
			if _, ok := providers.LangSourceNames[lang][source]; ok {
				sources = append(sources, source)
				addedSources[source] = struct{}{}
				break
			}
		}
	}

	if len(sources) == 0 {
		return nil, &common.Err{
			Code:    400,
			Message: "no valid sources",
		}
	}

	return sources, nil
}

func getLangs(langsStr string) ([]common.Lang, *common.Err) {
	if langsStr == "" {
		return providers.ActiveLangs, nil
	}

	langStrs := strings.Split(langsStr, ",")

	if len(langStrs) > 3 {
		return nil, &common.Err{
			Code:    400,
			Message: "too many languages",
		}
	}

	langs := make([]common.Lang, 0)

	for _, l := range langStrs {
		lang := strings.TrimSpace(l)

		if _, ok := supportedLangs[common.Lang(lang)]; !ok {
			return nil, common.ErrUnsupportedLang
		}

		langs = append(langs, common.Lang(lang))
	}

	return langs, nil
}
