package fetcher

import (
	"bytes"
	"context"
	"fmt"

	"github.com/markusmobius/go-trafilatura"
)

// ContentExtractor provides advanced content extraction using Trafilatura
type ContentExtractor struct {
	httpClient *HTTPClient
}

// NewContentExtractor creates a new content extractor
func NewContentExtractor(httpClient *HTTPClient) *ContentExtractor {
	return &ContentExtractor{
		httpClient: httpClient,
	}
}

// ExtractArticleFromURL extracts article content from a URL using Trafilatura
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

// ExtractArticleFromHTML extracts article content from HTML bytes using Trafilatura
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
