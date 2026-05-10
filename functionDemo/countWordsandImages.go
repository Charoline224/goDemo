package main

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func CountWordsandImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		resp.Body.Close()
		return
	}
	words, images = countWordsandImages(doc)
	return

}

func countWordsandImages(n *html.Node) (words, images int) {
	if n.Type == html.TextNode {
		//strings.Fields("hello world foo") → ["hello", "world", "foo"]
		words += len(strings.Fields(n.Data))
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		images++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w, i := countWordsandImages(c)
		words += w
		images += i
	}
	return words, images
}
