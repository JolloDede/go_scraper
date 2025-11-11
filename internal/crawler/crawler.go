package crawler

import (
	"fmt"
	"sync"

	uri "net/url"

	"github.com/JolloDede/go-crawler/internal/fetcher"
)

func Crawl(url string, depth uint64) {
	fetchedUrls := make(map[string]int)
	var mu sync.Mutex
	baseUrl, err := uri.ParseRequestURI(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	var crawl func(string, uint64)
	crawl = func(currentUrl string, depth uint64) {
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
			newUrl, err := uri.ParseRequestURI(u)

			if err != nil {
				newUrl, err = baseUrl.Parse(u)

				if err != nil {
					fmt.Println("Couldnt call ", u)
					return
				}
			}

			crawl(newUrl.String(), depth-1)
		}
	}

	crawl(baseUrl.String(), depth)
}
