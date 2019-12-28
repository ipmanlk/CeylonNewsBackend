package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type HiruNews struct {
}

func (p *HiruNews) SourceName() string {
	return "Hiru News"
}

func (p *HiruNews) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://www.hirunews.lk/english/",
		common.LangSi: "https://www.hirunews.lk/",
		common.LangTa: "https://www.hirunews.lk/tamil/",
	}
}

func (p *HiruNews) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	selectors := map[common.Lang]string{
		common.LangEn: ".row .section-tittle a",
		common.LangSi: ".row .section-tittle a",
		common.LangTa: ".row .section-tittle a",
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:             supportedLangs[lang],
		ListStrategy:    common.ExtractorStrategyHTML,
		ListSelector:    selectors[lang],
		Limit:           maxNewsItemsPerFetch,
		ArticleStrategy: common.ExtractorStrategyHTML,
	})
}
