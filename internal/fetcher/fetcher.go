package fetcher

import (
	"net/http"

	"golang.org/x/net/html"
)

func Fetch(uri string) (*http.Response, []string, error) {
	urls := make([]string, 0)
	resp, err := http.Get(uri)

	if err != nil {
		return nil, nil, err
	}

	node, err := html.Parse(resp.Body)

	if err != nil {
		return nil, nil, err
	}

	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i := 0; i < len(n.Attr); i++ {
				if n.Attr[i].Key == "href" {
					// convert relative url to full

					urls = append(urls, n.Attr[i].Val)
				}
			}
		}
	}

	return resp, urls, nil
}
