package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
)

type NewsLk struct {
}

func (p *NewsLk) SourceName() string {
	return "News.lk"
}

func (p *NewsLk) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://news.lk/news?format=feed",
		common.LangSi: "https://sinhala.news.lk/news?format=feed",
		common.LangTa: "https://tamil.news.lk/news?format=feed",
	}
}

func (p *NewsLk) FetchArticles(lang common.Lang) ([]common.Article, error) {
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
