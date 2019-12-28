package providers

import (
	"ipmanlk/cnapi/common"
	"ipmanlk/cnapi/extractor"
	"time"
)

type LankaPuvath struct {
}

func (p *LankaPuvath) SourceName() string {
	return "Lanka Puvath"
}

func (p *LankaPuvath) SupportedLangs() map[common.Lang]string {
	return map[common.Lang]string{
		common.LangEn: "https://english.lankapuvath.lk/feed/",
		common.LangSi: "https://sinhala.lankapuvath.lk/feed/",
	}
}

func (p *LankaPuvath) FetchArticles(lang common.Lang) ([]common.Article, error) {
	supportedLangs := p.SupportedLangs()

	if _, ok := supportedLangs[lang]; !ok {
		return nil, common.ErrUnsupportedLang
	}

	return extractor.ExtractArticles(common.ExtractorConfig{
		URL:             supportedLangs[lang],
		ListStrategy:    common.ExtractorStrategyRSS,
		ArticleStrategy: common.ExtractorStrategyHTML,
		Limit:           maxNewsItemsPerFetch,
		Delay:           time.Second * 10,
	})
}
