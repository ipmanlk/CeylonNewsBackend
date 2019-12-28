package common

import (
	"bytes"
	"encoding/json"
	"time"
)

type Lang string

const (
	LangEn Lang = "en"
	LangSi Lang = "si"
	LangTa Lang = "ta"
)

type Article struct {
	Title        string    `json:"title"`
	ContentText  string    `json:"content_text"`
	ContentHTML  string    `json:"content_html"`
	URL          string    `json:"url"`
	ThumbnailURL *string   `json:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at"`
}

type NewsItem struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	URL          string    `gorm:"type:VARCHAR(255);unique" json:"url"`
	ThumbnailURL *string   `json:"thumbnailUrl,omitempty"`
	Language     Lang      `gorm:"type:VARCHAR(2);index:idx_lang" json:"language"`
	SourceName   string    `gorm:"type:VARCHAR(255);index:idx_source_name" json:"sourceName"`
	CreatedAt    time.Time `json:"createdDate"`
	ContentText  string    `gorm:"type:TEXT" json:"-"`
	ContentHTML  string    `gorm:"type:TEXT" json:"content"`
}

type NewsProvider interface {
	SourceName() string
	FetchArticles(lang Lang) ([]Article, error)
	SupportedLangs() map[Lang]string
}

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrUnsupportedLang = &Err{
		Code:    400,
		Message: "unsupported language",
	}
	ErrUnexpectedResponse = &Err{
		Code:    500,
		Message: "unexpected server response",
	}
	ErrInvalidRequestMethod = &Err{
		Code:    405,
		Message: "invalid request method",
	}
	ErrBadRequest = &Err{
		Code:    400,
		Message: "bad request",
	}
	ErrInternalServer = &Err{
		Code:    500,
		Message: "internal server error",
	}
)

func (e *Err) Error() string {
	return e.Message
}

type ExtractorStrategy string

var (
	ExtractorStrategyRSS  ExtractorStrategy = "rss"
	ExtractorStrategyHTML ExtractorStrategy = "html"
)

type ExtractorConfig struct {
	URL              string
	ListStrategy     ExtractorStrategy
	ListSelector     string
	ArticleStrategy  ExtractorStrategy
	Replacements     map[string]string
	ArticleURLPrefix string
	Limit            int
	UseSplash        bool
	Delay            time.Duration
	SkipQueries      []string
}

// json.Marshal doesn't support SetEscapeHTML.
// This will marshal the given object without escaping HTML.
func JSONMarshal(t interface{}) ([]byte, error) {
    buffer := &bytes.Buffer{}
    encoder := json.NewEncoder(buffer)
    encoder.SetEscapeHTML(false)
    err := encoder.Encode(t)
    return buffer.Bytes(), err
}
