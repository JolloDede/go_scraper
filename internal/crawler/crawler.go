package crawler

import (
	"fmt"
	"sync"

	uri "net/url"

	"github.com/JolloDede/go-crawler/internal/fetcher"
)

func Crawl(url string, depth uint64) map[string]int {
	fetchedUrls := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	if !urlIsAbsolut(url) {
		println("Url is not absolut")
		return nil
	}

	baseUrl, _ := uri.Parse(url)

	var crawl func(string, uint64)
	crawl = func(currentUrl string, depth uint64) {
		defer wg.Done()

		if depth == 0 {
			return
		}
		resp, urls, err := fetcher.Fetch(currentUrl)

		if err != nil {
			fmt.Println(err)
			return
		}

		mu.Lock()
		fetchedUrls[currentUrl] = resp.StatusCode
		mu.Unlock()

		for _, u := range urls {
			newUrl, _ := uri.Parse(u)
			if !urlIsAbsolut(u) {
				newUrl, _ = baseUrl.Parse(u)
			}

			mu.Lock()
			_, ok := fetchedUrls[newUrl.String()]
			mu.Unlock()
			if ok {
				return
			}

			newDepth := depth - 1
			if baseUrl.Host != newUrl.Host {
				newDepth = 1
			}
			wg.Add(1)
			go crawl(newUrl.String(), newDepth)
		}
	}
	wg.Add(1)
	crawl(url, depth)
	wg.Wait()

	return fetchedUrls
}

func urlIsAbsolut(u string) bool {
	newUrl, err := uri.ParseRequestURI(u)

	if err != nil {
		return false
	}

	return newUrl.IsAbs()
}
