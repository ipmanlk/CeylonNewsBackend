package extractor

import (
	"bytes"
	"fmt"
	"io"
	"ipmanlk/cnapi/common"
	"log"
	"net/http"
	nurl "net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	rb "github.com/go-shiori/go-readability"
	traf "github.com/markusmobius/go-trafilatura"
	"github.com/mmcdole/gofeed"
)

// At least this is better than the previous one
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⢀⣠⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣷⣦⣄⡀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⣴⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⢀⣾⠟⠛⢉⣁⣀⠉⠉⢛⣿⣿⡛⠉⠁⣀⣉⡉⠛⢿⣷⡀⠀⠀⠀⠀
// ⠀⠀⠀⠀⣾⣿⣷⣾⠟⠋⣡⣴⣾⣿⣿⣿⣿⣷⣦⣈⠙⠻⣷⣾⣿⣷⠀⠀⠀⠀
// ⠀⠀⠀⢀⣿⣿⣿⣿⣶⣿⣿⣿⡿⠿⠿⠿⢿⣿⣿⣿⣿⣾⣿⣿⣿⣿⡄⠀⠀⠀
// ⠀⠀⠀⢸⣿⣿⣿⡿⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠿⣿⣿⣿⣿⡇⠀⠀⠀
// ⠀⠀⠀⠘⣿⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣿⣿⠃⠀⠀⠀
// ⠀⠀⠀⠀⢿⡇⠀⠀⢀⣠⣤⣶⣶⣶⣶⣶⣶⣶⣶⣦⣄⡀⠀⠀⢸⡿⠀⠀⠀⠀
// ⠀⠀⠀⠀⠈⢿⣦⣄⣈⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⠟⢉⣁⣤⣴⡿⠁⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠻⣿⣿⣷⡄⢹⣿⣿⣿⣿⣿⣿⡏⢰⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠈⠙⠿⡇⢸⣿⣿⣿⣿⣿⣿⡇⢸⠿⠋⠁⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠛⠛⠛⠛⠛⠛⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// still shit though.

func ExtractArticles(c common.ExtractorConfig) ([]common.Article, error) {
	switch c.ListStrategy {
	case common.ExtractorStrategyRSS:
		return extractArticlesFromFeed(c)
	case common.ExtractorStrategyHTML:
		return extractArticlesFromHTML(c)
	default:
		return nil, fmt.Errorf("unsupported list strategy: %s", c.ListStrategy)
	}
}

func extractArticlesFromFeed(c common.ExtractorConfig) ([]common.Article, error) {
	if c.ListStrategy == common.ExtractorStrategyHTML {
		return nil, fmt.Errorf("HTML Link strategy is not supported for feed extractors")
	}

	items, err := extractItemsFromFeed(formatURL(c.URL, c.UseSplash))
	if err != nil {
		return nil, err
	}

	time.Sleep(c.Delay)

	var articleURLs = make(map[string]struct{})
	var articles []common.Article

	for _, item := range items {
		if _, ok := articleURLs[item.Link]; ok {
			continue
		}

		if item.Link == "" {
			log.Printf("skipping [EMPTY ARTICLE URL] for %s\n", c.URL)
			continue
		}

		articleURL := fmt.Sprintf("%s%s", c.ArticleURLPrefix, item.Link)
		articleURLs[articleURL] = struct{}{}

		var err error
		var article *common.Article

		switch c.ArticleStrategy {
		case common.ExtractorStrategyHTML:
			article, err = extractArticleFromURL(articleURL, c.UseSplash)
		case common.ExtractorStrategyRSS:
			article, err = extractArticleFromFeedItem(item)
		}

		if err != nil {
			log.Printf("failed to extract article from %s: %v\n", item.Link, err)
			continue
		}

		if article == nil {
			log.Printf("failed to extract article from %s: article is nil\n", item.Link)
			continue
		}

		if isFetchingBlocked(article) {
			return nil, fmt.Errorf("fetching blocked for %s", article.URL)
		}

		if hasSkipQuery(article, c.SkipQueries) {
			continue
		}

		cleanedArticle := cleanContent(*article, c.Replacements)
		articles = append(articles, cleanedArticle)

		if len(articles) >= c.Limit {
			break
		}

		time.Sleep(c.Delay)
	}

	if len(articles) == 0 {
		return nil, fmt.Errorf("no articles extracted from %s", c.URL)
	}

	return articles, nil
}

func extractArticlesFromHTML(c common.ExtractorConfig) ([]common.Article, error) {
	if c.ArticleStrategy == common.ExtractorStrategyRSS ||
		c.ListStrategy == common.ExtractorStrategyRSS {
		return nil, fmt.Errorf("RSS strategy is not supported for URL extractors")
	}

	doc, err := getGoQueryDocFromURL(c.URL)
	if err != nil {
		return nil, err
	}

	var articleURLs = make(map[string]struct{})

	doc.Find(c.ListSelector).Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("href", "")
		if url == "" {
			log.Printf("skipping [EMPTY ARTICLE URL] for %s\n", c.URL)
			return
		}

		articleURL := fmt.Sprintf("%s%s", c.ArticleURLPrefix, url)

		if _, ok := articleURLs[articleURL]; ok {
			return
		}

		articleURLs[articleURL] = struct{}{}
	})

	time.Sleep(c.Delay)

	var articles []common.Article

	for articleURL := range articleURLs {
		article, err := extractArticleFromURL(articleURL, c.UseSplash)
		if err != nil {
			log.Printf("failed to extract article from %s: %v\n", articleURL, err)
			continue
		}

		if article == nil {
			log.Printf("failed to extract article from %s: article is nil\n", articleURL)
			continue
		}

		if isFetchingBlocked(article) {
			return nil, fmt.Errorf("fetching blocked for %s", article.URL)
		}

		if hasSkipQuery(article, c.SkipQueries) {
			continue
		}

		cleanedArticle := cleanContent(*article, c.Replacements)
		articles = append(articles, cleanedArticle)

		if len(articles) >= c.Limit {
			break
		}

		time.Sleep(c.Delay)
	}

	if len(articles) == 0 {
		return nil, fmt.Errorf("no articles extracted from %s", c.URL)
	}

	return articles, nil
}

func extractArticleFromURL(url string, useSplash bool) (*common.Article, error) {
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", url, err)
	}

	// Fetch article
	httpClient := &http.Client{Timeout: 30 * time.Second}

	resp, err := httpClient.Get(formatURL(url, useSplash))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return extractArticleFromBytes(body, parsedURL)
}

func extractArticleFromBytes(body []byte, pageURL *nurl.URL) (*common.Article, error) {
	// Extract data using readability
	readabilityRes, err := rb.FromReader(bytes.NewReader(body), pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article: %v", err)
	}

	// Extract data using trafilatura
	opts := traf.Options{
		IncludeImages:   true,
		ExcludeComments: true,
		Deduplicate:     true,
	}

	trafRes, err := traf.Extract(bytes.NewReader(body), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article: %v", err)
	}

	thumbnailURL := trafRes.Metadata.Image
	if thumbnailURL == "" {
		thumbnailURL = readabilityRes.Image
	}

	createdAt := trafRes.Metadata.Date
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	cleanContentHTML, err := cleanHTMLString(readabilityRes.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to clean HTML: %v", err)
	}

	return &common.Article{
		Title:        trafRes.Metadata.Title,
		ContentText:  trafRes.ContentText,
		ContentHTML:  cleanContentHTML,
		URL:          pageURL.String(),
		ThumbnailURL: &thumbnailURL,
		CreatedAt:    createdAt,
	}, nil
}

func extractArticleFromFeedItem(item *gofeed.Item) (*common.Article, error) {
	createdAt := *item.PublishedParsed
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	thumbnailURL := ""

	if item.Image != nil {
		thumbnailURL = item.Image.URL
	}

	if thumbnailURL == "" {
		var doc *goquery.Document
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(item.Description)))
		if err == nil {
			thumbnailURL = doc.Find("img").First().AttrOr("src", "")
		}
	}

	thumbnailP := &thumbnailURL
	if thumbnailURL == "" {
		thumbnailP = nil
	}

	return &common.Article{
		Title:        item.Title,
		ContentText:  item.Description,
		ContentHTML:  item.Content,
		URL:          item.Link,
		ThumbnailURL: thumbnailP,
		CreatedAt:    createdAt,
	}, nil
}
