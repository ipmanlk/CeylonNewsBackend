package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type NewsFirst struct {
}

func (p *NewsFirst) SourceName() string {
	return "News First"
}

func (p *NewsFirst) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://english.newsfirst.lk/latest",
		common.LangSi: "https://sinhala.newsfirst.lk/latest",
		common.LangTa: "https://tamil.newsfirst.lk/latest",
	}
}

func (p *NewsFirst) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	selectors := map[common.Lang]string{
		common.LangEn: ".local_news .ng-star-inserted a",
		common.LangSi: ".local_news .ng-star-inserted a",
		common.LangTa: ".local_news .ng-star-inserted a",
	}

	prefixes := map[common.Lang]string{
		common.LangEn: "https://english.newsfirst.lk",
		common.LangSi: "https://sinhala.newsfirst.lk",
		common.LangTa: "https://tamil.newsfirst.lk",
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:              supportedLangs[lang],
		ListStrategy:     common.ExtractorStrategyHTML,
		ListSelector:     selectors[lang],
		ArticleStrategy:  common.ExtractorStrategyHTML,
		ArticleURLPrefix: prefixes[lang],
		Limit:            maxNewsItemsPerFetch,
		SkipQueries:      []string{"Newsfirst.lk - Sri Lanka"},
	})
}
