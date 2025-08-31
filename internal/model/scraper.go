package model

import "time"

type ScrapedArticle struct {
	SourceName  string
	Title       string
	URL         string
	Content     string
	ImageURL    *string
	Categories  []string
	Language    Language
	PublishedAt time.Time
}
