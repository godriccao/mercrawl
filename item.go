package mercrawl

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

// Item represents a mercari item
type Item struct {
	id          string
	name        string
	photos      []string
	status      string
	price       int
	shippingFee string
	description string
	url         string
}

// Exists checks if an item is already in database
func (item *Item) Exists() bool {
	sql := "SELECT count(*) as total FROM items where id = $1"
	total := 0
	err := db.QueryRow(sql, item.id).Scan(&total)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if total > 0 {
		return true
	}
	return false
}

// Save persists the item to database
func (item *Item) Save() {
	insertSQL := `
	INSERT INTO items (id, name, photo1, photo2, photo3, photo4, status, price, shippingFee, description, url) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	// Scan() here is critical. Without calling it will cause leaking connections.
	// Refer to https://www.vividcortex.com/blog/2015/09/22/common-pitfalls-go/
	db.QueryRow(insertSQL,
		item.id,
		item.name,
		item.photos[0], item.photos[1], item.photos[2], item.photos[3],
		item.status,
		item.price,
		item.shippingFee,
		item.description,
		item.url).Scan()

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

	item := Item{id: itemRegexp.FindStringSubmatch(url)[1], url: url, photos: make([]string, 4)}
	if item.Exists() {
		fmt.Println("Item " + item.id + " already exists. Skip.")
	}

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

	item.Save()
}
