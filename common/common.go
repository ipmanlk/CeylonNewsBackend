package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	ThumbnailURL *string   `json:"thumbnailURL,omitempty"`
	Language     Lang      `gorm:"type:VARCHAR(2);index:idx_lang" json:"language"`
	SourceName   string    `gorm:"type:VARCHAR(255);index:idx_source_name" json:"sourceName"`
	CreatedAt    time.Time `json:"createdAt"`
	ContentText  string    `gorm:"type:TEXT" json:"-"`
	ContentHTML  string    `gorm:"type:TEXT" json:"contentHTML"`
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

type PaginationDirection string

var (
	PaginationDirectionNext PaginationDirection = "next"
	PaginationDirectionPrev PaginationDirection = "prev"
)

type PaginationPaging struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type PaginationResponse struct {
	Data   []NewsItem       `json:"data"`
	Paging PaginationPaging `json:"paging"`
}

type CursorData struct {
	ItemID    uint
	Direction PaginationDirection
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

func StringToLangs(langsStr string) ([]Lang, error) {
	langs := make([]Lang, 0)

	for _, langStr := range []rune(langsStr) {
		switch string(langStr) {
		case "en":
			langs = append(langs, LangEn)
		case "si":
			langs = append(langs, LangSi)
		case "ta":
			langs = append(langs, LangTa)
		default:
			return nil, ErrUnsupportedLang
		}
	}

	return langs, nil
}


func CreateCursor(itemID uint, direction PaginationDirection) string {
	cursorStr := fmt.Sprintf("%d/%s", itemID, direction)
	return base64.StdEncoding.EncodeToString([]byte(cursorStr))
}

func DecodeCursor(cursor string) (*CursorData, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, errors.New("failed to decode cursor")
	}

	cursorParts := strings.Split(string(decoded), "/")

	if len(cursorParts) != 2 {
		return nil, errors.New("invalid cursor format")
	}

	itemID, err := strconv.Atoi(cursorParts[0])
	if err != nil {
		return nil, errors.New("failed to parse itemID")
	}

	if itemID < 0 {
		return nil, errors.New("invalid itemID")
	}

	directionStr := cursorParts[1]

	var direction PaginationDirection
	switch PaginationDirection(directionStr) {
	case PaginationDirectionNext, PaginationDirectionPrev:
		direction = PaginationDirection(directionStr)
	default:
		return nil, errors.New("invalid pagination direction")
	}

	return &CursorData{
		ItemID:    uint(itemID),
		Direction: direction,
	}, nil
}
