package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func run() {
	for _, url := range os.Args[1:] {
		links, err := findlinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "findlinks:%v\n", err)
			continue
		}
		for _, link := range links {
			fmt.Println(link)
		}
	}

}
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

// 返回链接列表和错误信息
func findlinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//传播错误：
	//当对html.Parse的调用失败时，findLinks不会直接返回html.Parse的错误，因为缺少两条重要信息：
	// 1、发生错误时的解析器（html parser）；（哪个组件出错）
	// 2、发生错误的url。（哪个位置出错）
	// 因此，findLinks构造了一个新的错误信息，既包含了这两项，也包括了底层的解析出错的信息。
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	return visit(nil, doc), nil

}
