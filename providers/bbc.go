package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type BBCProvider struct {
}

func (p *BBCProvider) SourceName() string {
	return "BBC"
}

func (p *BBCProvider) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://www.bbc.com/news/topics/cywd23g0gxgt",
		common.LangSi: "https://www.bbc.com/sinhala/topics/cg7267dz901t",
		common.LangTa: "https://www.bbc.com/tamil/topics/cz74k7p3qw7t",
	}
}

func (p *BBCProvider) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	selectors := map[common.Lang]string{
		common.LangEn: ".ssrcss-1mrs5ns-PromoLink.exn3ah91",
		common.LangSi: ".focusIndicatorDisplayBlock.bbc-uk8dsi.e1d658bg0",
		common.LangTa: ".focusIndicatorDisplayBlock.bbc-uk8dsi.e1d658bg0",
	}

	prefixes := map[common.Lang]string{
		common.LangEn: "https://www.bbc.com",
		common.LangSi: "",
		common.LangTa: "",
	}

	replacements := map[string]string{
		" - BBC News සිංහල": "",
		" - BBC News தமிழ்": "",
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:              supportedLangs[lang],
		ListStrategy:     common.ExtractorStrategyHTML,
		ListSelector:     selectors[lang],
		ArticleStrategy:  common.ExtractorStrategyHTML,
		ArticleURLPrefix: prefixes[lang],
		Replacements:     replacements,
		Limit:            maxNewsItemsPerFetch,
	})
}
