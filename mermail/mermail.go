package main

import (
	"fmt"
	"os"

	"github.com/godriccao/mercrawl"
)

func main() {
	to := os.Args[1]
	go mercrawl.Mail(to)
	mercrawl.WaitInterupt()
	fmt.Printf("\n")
}
