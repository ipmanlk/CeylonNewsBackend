package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type Derana struct {
}

func (p *Derana) SourceName() string {
	return "Ada Derana"
}

func (p *Derana) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://www.adaderana.lk/rss.php",
		common.LangSi: "https://sinhala.adaderana.lk/rsshotnews.php",
		common.LangTa: "https://tamil.adaderana.lk/rss.php",
	}
}

func (p *Derana) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	// selectors := map[common.Lang]string{
	// 	common.LangEn: ".story-text a",
	// 	common.LangSi: ".story-text a",
	// 	common.LangTa: ".story-text a",
	// }

	// prefixes := map[common.Lang]string{
	// 	common.LangEn: "",
	// 	common.LangSi: "https://sinhala.adaderana.lk/",
	// 	common.LangTa: "https://tamil.adaderana.lk/",
	// }

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:              supportedLangs[lang],
		ListStrategy:     common.ExtractorStrategyRSS,
		ArticleStrategy:  common.ExtractorStrategyHTML,
		Limit:            maxNewsItemsPerFetch,
		UseSplash:        true,
	})
}
