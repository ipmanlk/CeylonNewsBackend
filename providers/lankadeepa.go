package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type Lankadeepa struct {
}

func (p *Lankadeepa) SourceName() string {
	return "Lankadeepa"
}

func (p *Lankadeepa) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangSi: "https://www.lankadeepa.lk/rss/latest_news/1",
	}
}

func (p *Lankadeepa) FetchArticles(lang common.Lang) ([]common.Article, error) {
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
