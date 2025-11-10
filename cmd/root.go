package cmd

import (
	"flag"

	"github.com/JolloDede/go-crawler/internal/crawler"
)

func Execute() {
	var filename string
	flag.StringVar(&filename, "o", "output.csv", "define output filename")
	var location string
	flag.StringVar(&location, "u", "http://localhost:80", "starting location for the crawl")
	var depth uint64
	flag.Uint64Var(&depth, "d", 10, "starting location for the crawl")

	flag.Parse()

	crawler.Crawl(location, depth)
}
