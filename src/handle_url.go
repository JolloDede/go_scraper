package src

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"golang.org/x/net/html"
)

func HandleUrl(url string, recursiv bool) map[string]int {
	scraper := NewUrlScraper(url)

	scraper.checkUrl(url)
	if recursiv {
		for len(scraper.unsearched) != 0 {
			var x string
			x, scraper.unsearched = scraper.unsearched[len(scraper.unsearched)-1], scraper.unsearched[:len(scraper.unsearched)-1]
			scraper.checkUrl(x)
		}
	}

	return scraper.searched
}

type UrlScraper struct {
	baseUrl    string
	searched   map[string]int
	unsearched []string
}

func NewUrlScraper(u string) *UrlScraper {
	res, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	p := res.Path
	res.Path = ""
	res.RawQuery = ""
	res.RawFragment = ""
	return &UrlScraper{baseUrl: res.String(), searched: make(map[string]int, 0), unsearched: []string{p}}
}

func (s *UrlScraper) checkUrl(u string) {
	if u[0] == '/' {
		u = s.baseUrl + u
	}
	res, err := http.Get(u)

	fmt.Println(u)

	if err != nil {
		return
	}
	defer res.Body.Close()

	if sameBaseUrl(s.baseUrl, u) {
		s.parseBody(res.Body)
	}

	if slices.Contains([]int{300, 303}, res.StatusCode) {
		locations := res.Header["Location"]
		for i := 0; i < len(locations); i++ {
			s.checkUrl(locations[i])
		}
	}

	if slices.Contains([]int{301, 302, 305, 307, 308}, res.StatusCode) {
		location := res.Request.Response.Header["Location"]
		s.checkUrl(location[0])
	}

	if s.searched[u] == 0 {
		s.searched[u] = res.StatusCode
	}
}

func (s *UrlScraper) parseBody(body io.ReadCloser) {
	doc, err := html.Parse(body)

	if err != nil {
		return
	}

	s.traverseRec(doc)
}

func (s *UrlScraper) traverseRec(node *html.Node) {
	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i := 0; i < len(n.Attr); i++ {
				if n.Attr[i].Key == "href" {
					if n.Attr[i].Val[0] == '/' || (len(n.Attr[i].Val) > 4 && n.Attr[i].Val[0:4] == "http") {
						s.unsearched = append(s.unsearched, n.Attr[i].Val)
					}
				}
			}
		}
	}
}

func sameBaseUrl(u1 string, u2 string) bool {
	r1, e1 := url.Parse(u1)
	r2, e2 := url.Parse(u2)

	if e1 != nil {
		panic(e1)
	}
	if e2 != nil {
		panic(e2)
	}

	r1.Path = ""
	r1.RawQuery = ""
	r1.RawFragment = ""

	r2.Path = ""
	r2.RawQuery = ""
	r2.RawFragment = ""

	if r1.String() == r2.String() {
		return true
	} else {
		return false
	}
}
