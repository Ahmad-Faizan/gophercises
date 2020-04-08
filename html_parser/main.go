package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	filename := flag.String("file", "ex1.html", "the HTML file to parse the links")
	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	rootNode, err := html.Parse(file)
	if err != nil {
		fmt.Print(err)
	}

	links := generateLinks(parseLinkNodes(rootNode))

	fmt.Printf("%v\n", links)
}

func parseLinkNodes(root *html.Node) []*html.Node {
	if root.Type == html.ElementNode && root.Data == "a" {
		return []*html.Node{root}
	}

	var linkSlice []*html.Node
	for n := root.FirstChild; n != nil; n = n.NextSibling {
		linkSlice = append(linkSlice, parseLinkNodes(n)...)
	}

	return linkSlice
}

func generateLinks(links []*html.Node) []Link {
	var list []Link
	for _, l := range links {
		list = append(list, buildLink(l))
	}
	return list
}

func buildLink(n *html.Node) Link {
	var link Link
	for _, attrib := range n.Attr {
		if attrib.Key == "href" {
			link.Href = attrib.Val
		}
	}
	link.Text = strings.Join(strings.Fields(buildText(n)), " ")
	return link
}

func buildText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += buildText(c)
	}
	return text
}
