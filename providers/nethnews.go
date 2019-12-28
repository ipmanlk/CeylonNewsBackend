package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type NethNews struct {
}

func (p *NethNews) SourceName() string {
	return "Neth News"
}

func (p *NethNews) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangSi: "https://nethnews.lk/category/feed/5",
	}
}

func (p *NethNews) FetchArticles(lang common.Lang) ([]common.Article, error) {
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
