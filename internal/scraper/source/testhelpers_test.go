package source

import (
	"ipmanlk/cnapi/internal/fetcher"
	"time"
)

// createTestFetcher creates a fetcher for testing purposes
func createTestFetcher() *fetcher.Fetcher {
	httpClient := fetcher.NewHTTPClient(15 * time.Second)
	browserClient := fetcher.NewBrowserAPIClient("http://localhost:8000", 60*time.Second, 15)
	return fetcher.NewFetcher(httpClient, browserClient)
}
