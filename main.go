package main

import (
	"flag"
	"fmt"
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
	flag.Parse()
	fmt.Println(urlFlag)
	resp, err := http.Get(*urlFlag)
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

	links, _ := link.Parse(resp.Body)
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)

		}
	}
	for _, href := range hrefs {
		fmt.Println(href)
	}
}
