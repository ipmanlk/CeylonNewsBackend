package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type Mawrata struct {
}

func (p *Mawrata) SourceName() string {
	return "Mawrata"
}

func (p *Mawrata) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://mawratanews.lk/feed/",
		common.LangSi: "https://sinhala.mawratanews.lk/feed/",
	}
}

func (p *Mawrata) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:             supportedLangs[lang],
		ListStrategy:    common.ExtractorStrategyRSS,
		ArticleStrategy: common.ExtractorStrategyRSS,
		Limit:           maxNewsItemsPerFetch,
	})
}
