package mercrawl

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/lib/pq"
)

// Item represents a mercari item
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Photos      [4]string `json:"photos,omitempty"`
	Status      string    `json:"status,omitempty"`
	Price       int       `json:"price,omitempty"`
	ShippingFee string    `json:"shippingFee,omitempty"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url,omitempty"`
	Sent        bool      `json:"sent,omitempty"`
}

func crawlItem(url string, itemSem chan bool) {
	itemSem <- true
	defer func() { <-itemSem }()

	item := Item{ID: itemRegexp.FindStringSubmatch(url)[1], URL: url}
	if item.Exists() {
		fmt.Println("Item " + item.ID + " already exists. Skip.")
		return
	}

	fmt.Println("[Crawling Item]" + url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("ERROR: Failed to crwal " + url)
		return
	}

	// item.name
	doc.Find("h2.item-name").Each(func(i int, s *goquery.Selection) {
		item.Name = s.Text()
	})
	doc.Find(".item-photo img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("data-src")
		item.Photos[i] = src
	})
	doc.Find("table.item-detail-table tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("th").Text() == "商品の状態" {
			item.Status = s.Find("td").Text()
		}
	})
	doc.Find("span.item-price").Each(func(i int, s *goquery.Selection) {
		item.Price = ParsePrice(s.Text())
	})
	doc.Find("span.item-shipping-fee").Each(func(i int, s *goquery.Selection) {
		item.ShippingFee = s.Text()
	})
	doc.Find("div.item-description").Each(func(i int, s *goquery.Selection) {
		item.Description = s.Text()
	})

	item.Save()
}

// Exists checks if an item is already in database
func (item *Item) Exists() bool {
	sql := "SELECT count(*) as total FROM items WHERE id = $1"
	total := 0
	db := GetDb()
	err := db.QueryRow(sql, item.ID).Scan(&total)
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
		item.ID,
		item.Name,
		item.Photos[0],
		item.Photos[1],
		item.Photos[2],
		item.Photos[3],
		item.Status,
		item.Price,
		item.ShippingFee,
		item.Description,
		item.URL).Scan()

}

// GetUnsentItems gets all unsent items from database
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
			&item.ID,
			&item.Name,
			&item.Photos[0],
			&item.Photos[1],
			&item.Photos[2],
			&item.Photos[3],
			&item.Status,
			&item.Price,
			&item.ShippingFee,
			&item.Description,
			&item.URL)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return
}

// GetAllItems gets all items from database
func GetAllItems() (items []Item) {
	sql := "SELECT id, name, photo1, photo2, photo3, photo4, status, price, shippingFee, description, url, sent FROM items"
	db = GetDb()

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		item := Item{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Photos[0],
			&item.Photos[1],
			&item.Photos[2],
			&item.Photos[3],
			&item.Status,
			&item.Price,
			&item.ShippingFee,
			&item.Description,
			&item.URL,
			&item.Sent)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return
}

// GetItem gets one item from database
func GetItem(id string) (item Item) {
	sql := "SELECT id, name, photo1, photo2, photo3, photo4, status, price, shippingFee, description, url, sent FROM items WHERE id = $1"
	db = GetDb()

	err := db.QueryRow(sql, id).Scan(
		&item.ID,
		&item.Name,
		&item.Photos[0],
		&item.Photos[1],
		&item.Photos[2],
		&item.Photos[3],
		&item.Status,
		&item.Price,
		&item.ShippingFee,
		&item.Description,
		&item.URL,
		&item.Sent)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// MarkAsSent marks items as sent
func MarkAsSent(items []Item) {
	db := GetDb()
	sql := "UPDATE items SET sent = true WHERE id = ANY ($1)"

	var ids []string
	for _, item := range items {
		ids = append(ids, item.ID)
	}

	rows, err := db.Query(sql, pq.Array([]string(ids)))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
