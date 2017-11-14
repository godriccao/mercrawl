// Package mercrawl crawls pages from a begin point and the following pages in parallel
package mercrawl

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"golang.org/x/net/html"
)

const base string = "https://www.mercari.com/jp/search/?"
const pageWorkers int = 5
const itemWorkers int = 10

var itemRegexp = regexp.MustCompile("^https://item\\.mercari\\.com/jp/m[0-9]+/")
var pageRegexp = regexp.MustCompile("^/jp/search/\\?page=([0-9]+)")

// Start starts crawling all items of the search result page with search condition string
func Start(search string) {
	pageSem := make(chan bool, pageWorkers)
	itemSem := make(chan bool, itemWorkers)
	pageState := PageState{make(map[string]bool), &sync.Mutex{}}

	url := base + search
	go crawlPage(url, pageSem, itemSem, pageState)
}

func crawlPage(url string, pageSem chan bool, itemSem chan bool, pageState PageState) {
	pageSem <- true
	defer func() { <-pageSem }()

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
					go crawlItem(href, itemSem)
				case pageRegexp.MatchString(href):
					num := pageRegexp.FindStringSubmatch(href)[1]
					fmt.Println("found ", num)
					if pageState.Get(num) {
						return
					}
					pageState.Set(num)
					go crawlPage(base+href, pageSem, itemSem, pageState)
				}
			}
		}
	}
}

func crawlItem(url string, itemSem chan bool) {
	itemSem <- true
	defer func() { <-itemSem }()
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
