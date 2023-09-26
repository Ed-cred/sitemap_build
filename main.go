package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	link "github.com/Ed-cred/html_link_parser"
)

/*
1. GET web page
2. Parse all links on page
3. Build proper urls from the links
4. Filter out links that lead to different domains
5. Find all pages contained in the site
*/

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the URL that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 9, "maximum number of links to traverse")
	flag.Parse()
	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		fmt.Println(page)
	}
}

type empty struct{}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]empty)
	var q map[string]empty
	nq := map[string]empty{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]empty)
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = empty{}
			for _, link := range get(url) {
				nq[link] = empty{}
			}

		}
	}
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	reqUrl := resp.Request.URL

	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)

		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(prfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prfx)
	}
}
