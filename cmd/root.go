package cmd

import (
	"flag"
	"fmt"

	"github.com/JolloDede/go-crawler/internal/crawler"
)

func Execute() {
	var filename string
	flag.StringVar(&filename, "o", "output.csv", "define output filename")
	var location string
	flag.StringVar(&location, "u", "http://localhost:80", "starting location for the crawl")
	var depth uint64
	flag.Uint64Var(&depth, "d", 10, "starting location for the crawl")
	var code int
	flag.IntVar(&code, "c", 0, "select a specific status code to filter")

	flag.Parse()

	fetchedUrls := crawler.Crawl(location, depth)

	for u, statusCode := range fetchedUrls {
		if code == 0 || statusCode == code {
			fmt.Printf("url: %s\tStatusCode: %d\n", u, statusCode)
		}
	}
}
