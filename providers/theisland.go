package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type TheIsland struct {
}

func (p *TheIsland) SourceName() string {
	return "The Island"
}

func (p *TheIsland) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "http://island.lk/feed/",
	}
}

func (p *TheIsland) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:             supportedLangs[lang],
		ListStrategy:    common.ExtractorStrategyRSS,
		ArticleStrategy: common.ExtractorStrategyHTML,
		Limit:           maxNewsItemsPerFetch,
	})
}
