package mercrawl

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

var interval int
var dialer *gomail.Dialer

// Mail scans database and send new item info to the mail address
func Mail(to string) {
	s := os.Getenv("INTERVAL")

	interval, err := strconv.Atoi(s)

	// default scan inteval is 5 seconds
	if err != nil || interval <= 0 {
		interval = 30
	}

	interval = 10
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	go func() {
		for i := 0; i < 10; i++ {
			<-ticker.C
			go scanAndSend(to)
		}
	}()
}

func scanAndSend(to string) {
	println("send to " + to)
	db = GetDb()

	sql := "SELECT count(*) as total FROM items WHERE sent = false"
	total := 0
	db := GetDb()

	err := db.QueryRow(sql).Scan(&total)
	if err != nil {
		log.Fatal(err)
	}
	if total > 0 {
		items := GetUnsentItems()
		send(items, to)
	}
}

func send(items []Item, to string) {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "[Mercrawl] Found "+strconv.Itoa(len(items))+" items that match your search condition!")
	var body string
	body = "Hello <b>Bob</b> and <i>Cora</i>!"
	m.SetBody("text/html", body)

	d := getDialer()
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(items), "items sent.")
}

func getDialer() *gomail.Dialer {
	if dialer == nil {
		smtp := os.Getenv("SMTP_SERVER")
		username := os.Getenv("SMTP_USER")
		password := os.Getenv("SMTP_PWD")
		port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
		if err != nil {
			log.Fatal(err)
		}

		dialer = gomail.NewDialer(smtp, port, username, password)
	}
	return dialer
}
