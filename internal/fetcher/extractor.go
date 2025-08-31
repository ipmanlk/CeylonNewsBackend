package fetcher

import (
	"bytes"
	"context"
	"fmt"

	"github.com/markusmobius/go-trafilatura"
)

type ContentExtractor struct {
	httpClient *HTTPClient
}

func NewContentExtractor(httpClient *HTTPClient) *ContentExtractor {
	return &ContentExtractor{
		httpClient: httpClient,
	}
}

func (e *ContentExtractor) ExtractArticleFromURL(ctx context.Context, url string) (*trafilatura.ExtractResult, error) {
	html, err := e.httpClient.FetchHTML(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HTML from %s: %w", url, err)
	}

	opts := trafilatura.Options{
		IncludeLinks:    true,
		IncludeImages:   true,
		ExcludeComments: true,
		EnableFallback:  true,
		Deduplicate:     true,
	}

	result, err := trafilatura.Extract(bytes.NewReader(html), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content from %s: %w", url, err)
	}

	return result, nil
}

func (e *ContentExtractor) ExtractArticleFromHTML(html []byte) (*trafilatura.ExtractResult, error) {
	opts := trafilatura.Options{
		IncludeLinks:    true,
		IncludeImages:   true,
		ExcludeComments: true,
		EnableFallback:  true,
		Deduplicate:     true,
	}

	result, err := trafilatura.Extract(bytes.NewReader(html), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content from HTML: %w", err)
	}

	return result, nil
}
