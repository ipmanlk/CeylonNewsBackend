package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type DailyMirror struct {
}

func (p *DailyMirror) SourceName() string {
	return "Daily Mirror"
}

func (p *DailyMirror) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://www.dailymirror.lk/rss/todays_headlines/419",
	}
}

func (p *DailyMirror) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	replacements := map[string]string{
		" - Breaking News | Daily Mirror": "",
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:             supportedLangs[lang],
		ListStrategy:    common.ExtractorStrategyRSS,
		ArticleStrategy: common.ExtractorStrategyHTML,
		Replacements:    replacements,
		Limit:           maxNewsItemsPerFetch,
		SkipQueries: 	 []string{"An Error Was"},
	})
}
