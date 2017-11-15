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
	photos      [4]string
	status      string
	price       int
	shippingFee string
	description string
	url         string
}

// Exists checks if an item is already in database
func (item *Item) Exists() bool {
	sql := "SELECT count(*) as total FROM items WHERE id = $1"
	total := 0
	db := GetDb()
	err := db.QueryRow(sql, item.id).Scan(&total)
	if err != nil {
		log.Fatal(err)
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
	db := GetDb()

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

// GetUnsentItems get all unsent items from database
func GetUnsentItems() (items []Item) {
	sql := `SELECT id, name, photo1, photo2, photo3, photo4, status, price, shippingFee, description, url 
			FROM items 
			WHERE sent = false`
	db = GetDb()

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		item := Item{}
		err := rows.Scan(
			&item.id,
			&item.name,
			&item.photos[0],
			&item.photos[1],
			&item.photos[2],
			&item.photos[3],
			&item.status,
			&item.price,
			&item.shippingFee,
			&item.description,
			&item.url)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return
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
	if item.Exists() {
		fmt.Println("Item " + item.id + " already exists. Skip.")
	}

	// item.name
	doc.Find("h2.item-name").Each(func(i int, s *goquery.Selection) {
		item.name = s.Text()
	})
	doc.Find(".item-photo img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("data-src")
		item.photos[i] = src
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
