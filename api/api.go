package api

import (
	"log"
	"net/http"
	"os"
)

func Start() {
	port := os.Getenv("API_PORT")

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/v1.0/news", HandleGetNews)
	http.HandleFunc("/api/v1.0/news/", HandleGetNewsItem)
	http.HandleFunc("/api/v1.0/sources", HandleGetSources)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
