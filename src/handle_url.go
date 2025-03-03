package src

import (
	"fmt"
	"net/url"
)

func HandleUrl(url string, recursiv bool) map[string][]Searched {
	scraper := NewUrlScraper(url)

	scraper.Wg.Add(1)
	scraper.checkUrl(Link{Url: url, link: url})
	if recursiv {
		for len(scraper.unsearched) != 0 {
			var x Link
			x, scraper.unsearched = scraper.unsearched[len(scraper.unsearched)-1], scraper.unsearched[:len(scraper.unsearched)-1]
			fmt.Println("Check URL: ", x.link)

			scraper.Wg.Add(1)

			go scraper.checkUrl(x)

			if len(scraper.unsearched) == 0 {
				scraper.Wg.Wait()
			}
		}
		scraper.Wg.Wait()
	}

	r := make(map[string][]Searched, 0)

	for key, val := range scraper.searchedMap {
		if len(r[key]) == 0 {
			r[key] = make([]Searched, 0)
		}
		r[key] = append(r[key], val)
	}

	return r
}

func sameHost(u1 string, u2 string) bool {
	r1, e1 := url.Parse(u1)
	r2, e2 := url.Parse(u2)

	if e1 != nil {
		panic(e1)
	}
	if e2 != nil {
		panic(e2)
	}

	if r1.Host == r2.Host {
		return true
	} else {
		return false
	}
}
