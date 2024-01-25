package extractor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"ipmanlk/cnapi/common"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

func extractItemsFromFeed(feedURL string) ([]*gofeed.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(feedURL, ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %v", err)
	}

	return feed.Items, nil
}

func getGoQueryDocFromURL(url string) (*goquery.Document, error) {
	httpClient := &http.Client{Timeout: 30 * time.Second}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to extract HTML from %s: %v", url, err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(html)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from %s: %v", url, err)
	}

	return doc, nil
}

func cleanContent(article common.Article, replacements map[string]string) common.Article {
	for k, v := range replacements {
		article.Title = strings.ReplaceAll(article.Title, k, v)
		article.ContentText = strings.ReplaceAll(article.ContentText, k, v)
		article.ContentHTML = strings.ReplaceAll(article.ContentHTML, k, v)
	}
	return article
}

var fetchingBlockedPhrases = []string{
	"access to this site has been limited",
}

func isFetchingBlocked(article *common.Article) bool {
	for _, phrase := range fetchingBlockedPhrases {
		if strings.Contains(article.Title, phrase) || strings.Contains(article.ContentText, phrase) {
			return true
		}
	}
	return false
}

var splashURL = os.Getenv("SPLASH_URL")

func formatURL(url string, useSplash bool) string {
	if useSplash {
		url = fmt.Sprintf("%s/?url=%s", splashURL, url)
	}
	return url
}

func hasSkipQuery(article *common.Article, skipQueries []string) bool {
	for _, skipQuery := range skipQueries {
		if strings.Contains(article.Title, skipQuery) {
			return true
		}
	}
	return false
}


// removeAttributes removes specified attributes from all elements in the given goquery.Selection.
func removeAttributes(selection *goquery.Selection, attributes ...string) {
	for _, attribute := range attributes {
		selection.RemoveAttr(attribute)
	}
}

func cleanHTMLString(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	doc.Find("*").Each(func(_ int, selection *goquery.Selection) {
		removeAttributes(selection, "class", "id")
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
