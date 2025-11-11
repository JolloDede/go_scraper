package crawler

import (
	uri "net/url"
	"strings"
)

func CreateUrlList(baseUrl string) UrlList {
	return UrlList{baseUrl: baseUrl, Urls: make([]string, 0)}
}

type UrlList struct {
	baseUrl *uri.URL
	Urls    []string
}

func (u *UrlList) Add(url string) {
	if strings.HasPrefix(url, u.baseUrl) {
		u.Urls = append(u.Urls, url)
	} else if true {

	}
}
