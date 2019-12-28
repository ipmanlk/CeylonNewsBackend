package api

import "ipmanlk/cnapi/common"

var supportedLangs = map[common.Lang]struct{}{
	common.LangEn: {},
	common.LangSi: {},
	common.LangTa: {},
}

type getNewsData struct {
	Langs    []common.Lang
	Sources  []string
	Query    string
	Page     int
	PageSize int
}

type newsSource struct {
	Name string `json:"name"`
}
