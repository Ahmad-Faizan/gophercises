package main

import (
	"flag"
	"fmt"

	"./sitemap"
)

func main() {
	siteURL := flag.String("url", "http://calhoun.io", "the URL of the website to build the sitemap")
	// depth := flag.Int("depth", 10, "the max depth of the sitemap tree")

	flag.Parse()

	links := sitemap.MakeRequest(*siteURL)
	generateXML(links)
}

func generateXML(links []sitemap.Link) {
	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, link := range links {
		fmt.Printf("  <url>\n    <loc>%s</loc>\n  </url>\n", link.Href)
	}
	fmt.Println(`</urlset>`)
}
