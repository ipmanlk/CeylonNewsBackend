package scraper

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/providers"
	"ipmanlk/cnapi/sqldb"
	"log"
	"sync"
	"time"
)

func Start() {
	runScraper()

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runScraper()
		}
	}
}

func runScraper() {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(providers.ActiveProviders))

	for _, provider := range providers.ActiveProviders {
		wg.Add(1)
		go func(p common.NewsProvider) {
			defer wg.Done()

			err := processProvider(p)
			if err != nil {
				errorChan <- err
			}
		}(provider)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Collect errors from goroutines
	for err := range errorChan {
		log.Println(err)
	}
}

func processProvider(p common.NewsProvider) error {
	sourceName := p.SourceName()
	langs := p.SupportedLangs()

	for lang := range langs {
		log.Printf("Processing %s for %v\n", sourceName, lang)

		articles, err := p.FetchArticles(lang)
		if err != nil {
			return err
		}

		// log the number of articles fetched
		log.Printf("Fetched %d articles from %s for %v\n", len(articles), sourceName, lang)

		newsItems := make([]common.NewsItem, len(articles))

		for i, article := range articles {
			newsItems[i] = common.NewsItem{
				Title:        article.Title,
				ContentText:  article.ContentText,
				ContentHTML:  article.ContentHTML,
				URL:          article.URL,
				ThumbnailURL: article.ThumbnailURL,
				Language:     lang,
				SourceName:   sourceName,
				CreatedAt:    article.CreatedAt,
			}
		}

		err = sqldb.InsertItems(newsItems)
		if err != nil {
			return err
		}
	}

	return nil
}
