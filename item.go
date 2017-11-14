package mercrawl

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Item represents a mercari item
type Item struct {
	id          string
	name        string
	photos      []string
	status      string
	price       float32
	shippingFee string
	description string
	url         string
}

func crawlItem(url string, itemSem chan bool) {
	itemSem <- true
	defer func() { <-itemSem }()

	doc, err := goquery.NewDocument(url)
	fmt.Println(url)
	if err != nil {
		fmt.Println("ERROR: Failed to crwal " + url)
		return
	}

	item := Item{id: itemRegexp.FindStringSubmatch(url)[1], url: url}

	// item.name
	doc.Find("h2.item-name").Each(func(i int, s *goquery.Selection) {
		item.name = s.Text()
	})
	doc.Find(".item-photo img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("data-src")
		item.photos = append(item.photos, src)
	})
	doc.Find("table.item-detail-table tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("th").Text() == "商品の状態" {
			item.status = s.Find("td").Text()
		}
	})
	doc.Find("span.item-price").Each(func(i int, s *goquery.Selection) {
		item.price = ParsePrice(s.Text())
	})
	doc.Find("span.item-shipping-fee").Each(func(i int, s *goquery.Selection) {
		item.shippingFee = s.Text()
	})
	doc.Find("div.item-description").Each(func(i int, s *goquery.Selection) {
		item.description = s.Text()
	})
	fmt.Printf("%+v\n", item)
}
