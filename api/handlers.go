package api

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/providers"
	"ipmanlk/cnapi/sqldb"
	"log"
	"net/http"
)

// TODO: Configure CORS properly
func HandleGetSources(w http.ResponseWriter, r *http.Request) {
	langs, err := validateGetSources(r)
	if err != nil {
		http.Error(w, err.Message, err.Code)
		return
	}

	sources := make([]newsSource, 0)
	addedSources := make(map[string]struct{})

	for _, lang := range langs {
		for source := range providers.LangSourceNames[lang] {
			if _, ok := addedSources[source]; ok {
				continue
			}
			addedSources[source] = struct{}{}

			sources = append(sources, newsSource{
				Name: source,
			})
		}
	}

	responseData, er := common.JSONMarshal(sources)
	if er != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func HandleGetNews(w http.ResponseWriter, r *http.Request) {
	data, err := validateHandleGetNews(r)
	if err != nil {
		http.Error(w, err.Message, err.Code)
		return
	}

	items, er := sqldb.SearchItems(data.Langs, data.Sources, data.Query, data.Cursor, data.PageSize)
	if er != nil {
		log.Printf("Error getting data: %v", er)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	responseData, er := common.JSONMarshal(items)
	if er != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func HandleGetNewsItem(w http.ResponseWriter, r *http.Request) {
	id, err := validateHandleGetNewsItem(r)
	if err != nil {
		http.Error(w, err.Message, err.Code)
		return
	}

	item, er := sqldb.GetItemByID(id)

	if er != nil {
		log.Printf("Error getting data: %v", er)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	responseData, er := common.JSONMarshal(item)
	if er != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
