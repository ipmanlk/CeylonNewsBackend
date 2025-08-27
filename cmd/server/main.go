package main

import (
	"context"
	"fmt"
	"ipmanlk/cnapi/internal/model"
	"ipmanlk/cnapi/internal/scraper"
	"ipmanlk/cnapi/internal/scraper/source"
)

func main() {
	fmt.Println("Hello, World!")

	rssScraper := &scraper.RSSScraper{}
	bbcScraper := source.NewBBCScraper(rssScraper)

	articles, err := bbcScraper.Scrape(context.TODO(), model.LangEn)
	if err != nil {
		fmt.Printf("Error scraping BBC: %v\n", err)
		return
	}
	fmt.Printf("Scraped %d articles from BBC\n", len(articles))
	for _, article := range articles {
		fmt.Printf("Title: %s\nURL: %s\n\n", article.ContentText, article.URL)
	}
}
