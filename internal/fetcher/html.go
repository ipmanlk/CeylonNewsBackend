package fetcher

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HTMLProcessor struct{}

func NewHTMLProcessor() *HTMLProcessor {
	return &HTMLProcessor{}
}

func (h *HTMLProcessor) GetFirstImageFromHTML(html []byte) *string {
	var thumbnailURL string
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err == nil {
		thumbnailURL = doc.Find("img").First().AttrOr("src", "")
	}
	if thumbnailURL != "" {
		return &thumbnailURL
	}
	return nil
}

func (h *HTMLProcessor) CleanHTMLString(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	// Remove unwanted attributes
	doc.Find("*").Each(func(_ int, selection *goquery.Selection) {
		for _, attribute := range []string{"style", "width", "height", "onclick", "onerror", "align", "border", "cellpadding", "cellspacing"} {
			selection.RemoveAttr(attribute)
		}
	})

	// Check if the body tag exists
	bodyExists := false
	doc.Find("body").Each(func(_ int, selection *goquery.Selection) {
		bodyExists = true
	})

	// If the body tag exists, return the body content
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

	// If the body tag does not exist, return the whole document
	str, err := doc.Html()
	return str, err
}

func (h *HTMLProcessor) GetTextFromHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(doc.Text()), nil
}

func (h *HTMLProcessor) ExtractLinks(doc *goquery.Document, selector, urlPattern string) []string {
	var links []string
	doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if exists && href != "" && strings.HasPrefix(href, urlPattern) {
			links = append(links, href)
		}
	})
	return links
}
