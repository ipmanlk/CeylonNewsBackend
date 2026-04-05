package scraper

type ListingStrategy string

const (
	StrategyRSS                ListingStrategy = "rss"
	StrategyRSSBrowser         ListingStrategy = "rss_browser"
	StrategyRSSContentSelector ListingStrategy = "rss_content_selector"
	StrategyRSSArticleBrowser  ListingStrategy = "rss_article_browser"
	StrategyHTML               ListingStrategy = "html"
	StrategyHTMLBrowser        ListingStrategy = "html_browser"
)

type LangConfig struct {
	Language        string          `toml:"language"`
	Strategy        ListingStrategy `toml:"strategy"`
	FeedURL         string          `toml:"feed_url"`
	PageURL         string          `toml:"page_url"`
	LinkSelectors   []string        `toml:"link_selectors"`
	URLPrefix       string          `toml:"url_prefix"`
	BaseURL         string          `toml:"base_url"`
	MaxItems        int             `toml:"max_items"`
	ContentSelector string          `toml:"content_selector"`
}

type SkipRule struct {
	Contains      string `toml:"contains"`
	CaseSensitive bool   `toml:"case_sensitive"`
}

type ReplaceRule struct {
	Pattern       string `toml:"pattern"`
	With          string `toml:"with"`
	CaseSensitive bool   `toml:"case_sensitive"`
}

type TitleTransform struct {
	Skip    []SkipRule    `toml:"skip"`
	Replace []ReplaceRule `toml:"replace"`
}

type SourceConfig struct {
	Name           string         `toml:"name"`
	Languages      []LangConfig   `toml:"languages"`
	TitleTransform TitleTransform `toml:"title_transform"`
}
