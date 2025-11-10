package crawler

import (
	"fmt"
	"sync"

	"github.com/JolloDede/go-crawler/internal/fetcher"
)

func Crawl(uri string, depth uint64) {
	print("test")

	fetchedUrls := make(map[string]int)
	var mu sync.Mutex

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
			crawl(u, depth-1)
		}
	}

	crawl(uri, depth)
}
