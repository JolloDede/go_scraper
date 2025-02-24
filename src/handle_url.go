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
			fmt.Println("Check URL: ", x)
			scraper.checkUrl(x)
		}
	}

	return scraper.searched
}

type UrlScraper struct {
	scheme     string
	host       string
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
	return &UrlScraper{scheme: res.Scheme, host: res.Host, searched: make(map[string]int, 0), unsearched: []string{p}}
}

func (s *UrlScraper) checkUrl(u string) {
	res, err := http.Get(u)

	if err != nil {
		s.searched[u] = 404
		return
	}
	defer res.Body.Close()

	if sameHost(s.scheme+"://"+s.host, u) {
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
					s.addToUnsearched(n.Attr[i].Val)
				}
			}
		}
	}
}

func (s *UrlScraper) addToUnsearched(p string) {
	u, err := url.Parse(p)

	if err != nil {
		return
	}
	if !slices.Contains([]string{"http", "https", ""}, u.Scheme) {
		return
	}

	if u.Host == "" {
		// remove the // infront of the url
		u.Scheme = s.scheme
		u.Host = s.host
	}

	if !s.pathAlreadySearched(u.String()) {
		s.unsearched = append(s.unsearched, u.String())
	}
}

func (s *UrlScraper) pathAlreadySearched(p string) bool {
	if s.searched[p] == 0 {
		return false
	} else {
		return true
	}
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
