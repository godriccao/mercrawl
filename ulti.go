package mercrawl

import (
	"golang.org/x/net/html"
)

// GetAttr get href value from html token
func GetAttr(t html.Token, attr string) (ok bool, val string) {
	for _, a := range t.Attr {
		if a.Key == attr {
			val = a.Val
			ok = true
		}
	}

	return
}

// ParsePrice parse a JPY price like ¥ 168,800 to float32 168800.0
func ParsePrice(s string) (price int) {
	weight := 1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			price += int(s[i]-'0') * weight
			weight *= 10
		}
	}
	return
}
