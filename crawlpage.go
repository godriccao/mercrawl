// Package mercrawl crawls pages from a begin point and the following pages in parallel
package mercrawl

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/net/html"
)

const domain string = "https://www.mercari.com"
const base string = "https://www.mercari.com/jp/search/?"

var itemRegexp = regexp.MustCompile("^https://item\\.mercari\\.com/jp/(m[0-9]+)/")
var pageRegexp = regexp.MustCompile("/jp/search/\\?page=([0-9]+)")

// Start starts crawling all items of the search result page with search condition string
func Start(search string) {
	var pageWorkers int
	var itemWorkers int
	var err error
	pageWorkers, err = strconv.Atoi(os.Getenv("PAGE_WORKERS"))
	if err != nil || pageWorkers <= 0 {
		pageWorkers = 5
	}
	itemWorkers, err = strconv.Atoi(os.Getenv("ITEM_WORKERS"))
	if err != nil || itemWorkers <= 0 {
		itemWorkers = 20
	}

	url := base + search
	pageSem := make(chan bool, pageWorkers)
	itemSem := make(chan bool, itemWorkers)
	pageState := PageState{make(map[string]bool), &sync.Mutex{}}

	// Will not crawl ?page=1 since it is same with a page who does not have `page` parameter.
	if !pageRegexp.MatchString(url) {
		pageState.Set("1")
	}

	go crawlPage(url, pageSem, itemSem, &pageState)
}

func crawlPage(url string, pageSem chan bool, itemSem chan bool, pageState *PageState) {
	pageSem <- true
	defer func() { <-pageSem }()

	fmt.Println("[Crawling Page] " + url)

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
				ok, href := GetAttr(t, "href")
				if !ok {
					continue
				}
				switch {
				case itemRegexp.MatchString(href):
					go crawlItem(href, itemSem)
				case pageRegexp.MatchString(href):
					num := pageRegexp.FindStringSubmatch(href)[1]
					if pageState.Get(num) {
						return
					}
					pageState.Set(num)
					go crawlPage(domain+href, pageSem, itemSem, pageState)
				}
			}
		}
	}
}
