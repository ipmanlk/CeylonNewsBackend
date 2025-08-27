package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/markusmobius/go-trafilatura"
)

func GetFirstImageFromHTML(html []byte) *string {
	var thumbnailURL string
	var doc *goquery.Document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err == nil {
		thumbnailURL = doc.Find("img").First().AttrOr("src", "")
	}
	if thumbnailURL != "" {
		return &thumbnailURL
	}
	return nil
}

func CleanHTMLString(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	doc.Find("*").Each(func(_ int, selection *goquery.Selection) {
		for _, attribute := range []string{"style", "width", "height", "onclick", "onerror", "align", "border", "cellpadding", "cellspacing"} {
			selection.RemoveAttr(attribute)
		}
	})

	// check if the body tag exists
	bodyExists := false
	doc.Find("body").Each(func(_ int, selection *goquery.Selection) {
		bodyExists = true
	})

	// if the body tag exists, return the body content
	if bodyExists {
		var bodyContent string
		doc.Find("body").Each(func(_ int, selection *goquery.Selection) {
			bodyContent, err = selection.Html()
		})
		if err != nil {
			return "", err
		}
		return bodyContent, nil
	}

	// if the body tag does not exist, return the whole document
	str, err := doc.Html()
	return str, err
}

func GetTextFromHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(doc.Text()), nil
}

func GetGoQueryDocFromURL(ctx context.Context, url string) (*goquery.Document, error) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code (%d) for %s", resp.StatusCode, url)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from %s: %w", url, err)
	}

	return doc, nil
}

func ScrapeArticleFromURL(ctx context.Context, url string) (*trafilatura.ExtractResult, error) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code (%d) for %s", resp.StatusCode, url)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if len(html) == 0 {
		return nil, fmt.Errorf("empty response body for %s", url)
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
