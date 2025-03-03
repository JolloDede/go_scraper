package src

import (
	"io"
	"net/http"
	"net/url"
	"slices"
	"sync"

	"golang.org/x/net/html"
)

type UrlScraper struct {
	scheme      string
	host        string
	searchedMap map[string]Searched
	unsearched  []Link
	Wg          sync.WaitGroup
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
	defer s.Wg.Done()
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
