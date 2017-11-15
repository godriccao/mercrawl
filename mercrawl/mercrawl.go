package main

import (
	"fmt"
	"os"

	"github.com/godriccao/mercrawl"
)

func main() {
	search := os.Args[1]
	go mercrawl.Start(search)
	mercrawl.WaitInterupt()
	fmt.Printf("\n")
}
