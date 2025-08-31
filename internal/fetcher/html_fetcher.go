package fetcher

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type HTMLFetcher struct {
	httpClient    *HTTPClient
	browserClient *BrowserAPIClient
}

func NewHTMLFetcher(browserClient *BrowserAPIClient) *HTMLFetcher {
	return &HTMLFetcher{
		httpClient:    NewHTTPClient(),
		browserClient: browserClient,
	}
}

func (h *HTMLFetcher) FetchHTML(ctx context.Context, url string, useBrowser ...bool) ([]byte, error) {
	if len(useBrowser) > 0 && useBrowser[0] {
		return h.browserClient.FetchHTML(ctx, url)
	}
	return h.httpClient.FetchHTML(ctx, url)
}

func (h *HTMLFetcher) FetchHTMLDoc(ctx context.Context, url string, useBrowser ...bool) (*goquery.Document, error) {
	if len(useBrowser) > 0 && useBrowser[0] {
		return h.browserClient.FetchHTMLDoc(ctx, url)
	}
	return h.httpClient.FetchHTMLDoc(ctx, url)
}

func (h *HTMLFetcher) IsUsingBrowserAPI() bool {
	return true
}

func (h *HTMLFetcher) String() string {
	return fmt.Sprintf("HTMLFetcher(browser_api=true, url=%s)", h.browserClient.apiURL)
}
