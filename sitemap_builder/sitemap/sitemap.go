package sitemap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func MakeRequest(url string) []Link {
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	rawHTML := readBody(response)
	links := parseBody(rawHTML)
	return links
}

func readBody(response *http.Response) []byte {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return body
}

func parseBody(body []byte) []Link {
	rootNode, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	linkNodes := parseLinkNodes(rootNode)
	return generateLinks(linkNodes)
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
