package src

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"golang.org/x/net/html"
)

func HandleUrl(url string, recursiv bool) map[string][]Searched {
	scraper := NewUrlScraper(url)

	scraper.checkUrl(Link{Url: url, link: url})
	if recursiv {
		for len(scraper.unsearched) != 0 {
			var x Link
			x, scraper.unsearched = scraper.unsearched[len(scraper.unsearched)-1], scraper.unsearched[:len(scraper.unsearched)-1]
			fmt.Println("Check URL: ", x.link)
			scraper.checkUrl(x)
		}
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

type Link struct {
	Url     string
	link    string
	Content string
}

type Searched struct {
	Link       Link
	StatusCode int
}

type UrlScraper struct {
	scheme      string
	host        string
	searchedMap map[string]Searched
	unsearched  []Link
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
	return &UrlScraper{scheme: res.Scheme, host: res.Host, searchedMap: make(map[string]Searched, 0), unsearched: []Link{{Url: p, link: p}}}
}

func (s *UrlScraper) checkUrl(u Link) {
	res, err := http.Get(u.link)

	if err != nil {
		s.searchedMap[u.link] = Searched{Link: u, StatusCode: 404}
		return
	}
	defer res.Body.Close()

	if sameHost(s.scheme+"://"+s.host, u.link) {
		s.parseBody(u.link, res.Body)
	}

	if slices.Contains([]int{300, 303}, res.StatusCode) {
		locations := res.Header["Location"]
		for i := 0; i < len(locations); i++ {
			s.checkUrl(Link{Url: u.link, link: locations[i]})
		}
	}

	if slices.Contains([]int{301, 302, 305, 307, 308}, res.StatusCode) {
		location := res.Request.Response.Header["Location"]
		s.checkUrl(Link{Url: u.link, link: location[0]})
	}

	if s.searchedMap[u.link].StatusCode == 0 {
		s.searchedMap[u.link] = Searched{Link: u, StatusCode: res.StatusCode}
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
					content := ""
					if n.FirstChild != nil {
						content = n.FirstChild.Data
					}
					s.addToUnsearched(Link{Url: cur, link: n.Attr[i].Val, Content: content})
				}
			}
		}
	}
}

func (s *UrlScraper) addToUnsearched(link Link) {
	u, err := url.Parse(link.link)

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
		s.unsearched = append(s.unsearched, Link{Url: link.Url, link: u.String(), Content: link.Content})
	}
}

func (s *UrlScraper) pathAlreadySearched(p string) bool {
	if s.searchedMap[p].StatusCode == 0 {
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
