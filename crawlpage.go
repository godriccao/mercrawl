// Package mercrawl crawls pages from a begin point and the following pages in parallel
package mercrawl

import (
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

const base string = "https://www.mercari.com/jp/search/?"

var itemRegexp = regexp.MustCompile("^https://item.mercari.com/jp/m[0-9]+/")
var pageRegexp = regexp.MustCompile("^/jp/search/?page=[0-9]+")

// Start starts crawling all items of the search result page with search condition string
func Start(search string) {
	url := base + search
	go crawlPage(url)
}

func crawlPage(url string) {
	fmt.Println("Crawling " + url)
	res, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to crwal " + url)
		return
	}

	b := res.Body
	defer b.Close()

	tokens := html.NewTokenizer(b)

	for {
		tt := tokens.Next()

		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := tokens.Token()
			if t.Data == "a" {
				ok, href := getHref(t)
				if !ok {
					continue
				}
				switch {
				case itemRegexp.MatchString(href):
					go crawlItem(href)
				case pageRegexp.MatchString(href):
					go crawlPage(href)
				}
			}
		}
	}
}

func crawlItem(url string) {
	fmt.Println(url)
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	return
}
