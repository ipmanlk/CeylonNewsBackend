package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type Divaina struct {
}

func (p *Divaina) SourceName() string {
	return "Divaina"
}

func (p *Divaina) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangSi: "https://divaina.lk/feed",
	}
}

func (p *Divaina) FetchArticles(lang common.Lang) ([]common.Article, error) {
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
