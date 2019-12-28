package providers

import (
	"ipmanlk/cnapi/common"
)

const maxNewsItemsPerFetch = 10

var bbcProvider = &BBCProvider{}
var dailymirrorProvider = &DailyMirror{}
var deranaProvider = &Derana{}
var newsFirstProvider = &NewsFirst{}
var theislandProvider = &TheIsland{}
var nethnewsProvider = &NethNews{}
var newsLkProvider = &NewsLk{}
var lankaPuvathProvider = &LankaPuvath{}
var hiruNewsProvider = &HiruNews{}

// Disabled
// var gaganaProvider = &Gagana{}

var ActiveProviders = map[string]common.NewsProvider{
	bbcProvider.SourceName():         bbcProvider,
	dailymirrorProvider.SourceName(): dailymirrorProvider,
	newsFirstProvider.SourceName():   newsFirstProvider,
	theislandProvider.SourceName():   theislandProvider,
	nethnewsProvider.SourceName():    nethnewsProvider,
	newsLkProvider.SourceName():      newsLkProvider,
	lankaPuvathProvider.SourceName(): lankaPuvathProvider,
	hiruNewsProvider.SourceName():    hiruNewsProvider,
	deranaProvider.SourceName(): deranaProvider,
}

var ActiveProviderCount = 0
var ActiveProviderNames = make([]string, 0)
var ActiveLangs = []common.Lang{
	common.LangEn,
	common.LangSi,
	common.LangTa,
}
var LangSourceNames = map[common.Lang]map[string]struct{}{}

func init() {
	addedProviderNames := make(map[string]struct{})

	for _, provider := range ActiveProviders {
		providerName := provider.SourceName()

		for lang := range provider.SupportedLangs() {
			if _, ok := LangSourceNames[lang]; !ok {
				LangSourceNames[lang] = map[string]struct{}{}
			}
			LangSourceNames[lang][providerName] = struct{}{}
		}
		ActiveProviderCount++

		if _, ok := addedProviderNames[providerName]; ok {
			continue
		}

		addedProviderNames[providerName] = struct{}{}
		ActiveProviderNames = append(ActiveProviderNames, providerName)
	}
}
