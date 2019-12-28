package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type Gagana struct {
}

func (p *Gagana) SourceName() string {
	return "Gagana"
}

func (p *Gagana) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangSi: "https://gagana.lk/news/srilanka",
	}
}

func (p *Gagana) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	selectors := map[common.Lang]string{
		common.LangSi: ".a-item a",
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:             supportedLangs[lang],
		ListStrategy:    common.ExtractorStrategyHTML,
		ListSelector:    selectors[lang],
		ArticleStrategy: common.ExtractorStrategyHTML,
		Limit:           maxNewsItemsPerFetch,
		UseSplash:       true,
	})
}
