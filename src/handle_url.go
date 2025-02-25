package src

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"golang.org/x/net/html"
)

func HandleUrl(url string, recursiv bool) map[string]Searched {
	scraper := NewUrlScraper(url)

	scraper.checkUrl(ToSearch{currentPath: url, linkUrl: url})
	if recursiv {
		for len(scraper.unsearched) != 0 {
			var x ToSearch
			x, scraper.unsearched = scraper.unsearched[len(scraper.unsearched)-1], scraper.unsearched[:len(scraper.unsearched)-1]
			fmt.Println("Check URL: ", x.linkUrl)
			scraper.checkUrl(x)
		}
	}

	return scraper.searched
}

type ToSearch struct {
	currentPath string
	linkUrl     string
}

type Searched struct {
	path       string
	statusCode int
}

type UrlScraper struct {
	scheme     string
	host       string
	searched   map[string]Searched
	unsearched []ToSearch
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
	return &UrlScraper{scheme: res.Scheme, host: res.Host, searched: make(map[string]Searched, 0), unsearched: []ToSearch{{currentPath: p, linkUrl: p}}}
}

func (s *UrlScraper) checkUrl(u ToSearch) {
	res, err := http.Get(u.linkUrl)

	if err != nil {
		s.searched[u.linkUrl] = Searched{path: u.currentPath, statusCode: 404}
		return
	}
	defer res.Body.Close()

	if sameHost(s.scheme+"://"+s.host, u.linkUrl) {
		s.parseBody(u.linkUrl, res.Body)
	}

	if slices.Contains([]int{300, 303}, res.StatusCode) {
		locations := res.Header["Location"]
		for i := 0; i < len(locations); i++ {
			s.checkUrl(ToSearch{currentPath: u.linkUrl, linkUrl: locations[i]})
		}
	}

	if slices.Contains([]int{301, 302, 305, 307, 308}, res.StatusCode) {
		location := res.Request.Response.Header["Location"]
		s.checkUrl(ToSearch{currentPath: u.linkUrl, linkUrl: location[0]})
	}

	if s.searched[u.currentPath].statusCode == 0 {
		s.searched[u.currentPath] = Searched{path: u.currentPath, statusCode: res.StatusCode}
	}
}

func (s *UrlScraper) parseBody(cur string, body io.ReadCloser) {
	doc, err := html.Parse(body)

	if err != nil {
		return
	}

	s.traverseRec(cur, doc)
}

func (s *UrlScraper) traverseRec(cur string, node *html.Node) {
	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i := 0; i < len(n.Attr); i++ {
				if n.Attr[i].Key == "href" {
					s.addToUnsearched(cur, n.Attr[i].Val)
				}
			}
		}
	}
}

func (s *UrlScraper) addToUnsearched(cur string, p string) {
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
		s.unsearched = append(s.unsearched, ToSearch{currentPath: cur, linkUrl: u.String()})
	}
}

func (s *UrlScraper) pathAlreadySearched(p string) bool {
	if s.searched[p].statusCode == 0 {
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
