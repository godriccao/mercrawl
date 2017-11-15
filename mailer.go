package mercrawl

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/leekchan/accounting"
	"gopkg.in/gomail.v2"
)

var interval int
var dialer *gomail.Dialer

// Mail scans database and send new item info to the mail address
func Mail(to string) {
	s := os.Getenv("INTERVAL")
	interval, err := strconv.Atoi(s)
	if err != nil || interval <= 0 {
		interval = 30 // default scan inteval is 30 seconds
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for ; true; <-ticker.C { // Ensure scanAndSend() run instantly. Otherwise it will run after a tick
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
	m.SetBody("text/html", mailBody(items))

	d := getDialer()
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(items), "items sent.")
}

func mailBody(items []Item) (body string) {
	ac := accounting.Accounting{Symbol: "￥", Precision: 0}
	body = "Mercrawl Found " + strconv.Itoa(len(items)) + " new items that match your search condition!<br>"

	body += "<table>\n"
	for _, item := range items {
		body += "<tr>\n"
		body += "<td><a href=\"" + item.url + "\"><b>" + item.name + "</b></a></td>"
		body += "<td>" + ac.FormatMoney(item.price) + "</td>"
		body += "<td>" + item.status + "</td>"
		body += "<td>" + item.shippingFee + "</td>"
		body += "</tr>\n"
	}
	body += "</table>\n"

	return
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
