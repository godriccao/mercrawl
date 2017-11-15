package main

import (
	"fmt"
	"os"

	"github.com/godriccao/mercrawl"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage:\n\tmermail \"mail_to_address\"")
		os.Exit(1)
	}
	to := os.Args[1]
	go mercrawl.Mail(to)
	mercrawl.WaitInterupt()
	fmt.Printf("\n")
}
