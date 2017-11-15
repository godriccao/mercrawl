package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/godriccao/mercrawl"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage:\n\tmercrawl \"search_condition\"")
		os.Exit(1)
	}
	search := os.Args[1]

	s := os.Getenv("RECRAWL_INTERVAL")
	interval, err := strconv.Atoi(s)
	if err != nil || interval <= 0 {
		interval = 30 // default scan inteval is 30 seconds
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for ; true; <-ticker.C { // Ensure scanAndSend() run instantly. Otherwise it will run after a tick
			go mercrawl.Start(search)
		}
	}()

	mercrawl.WaitInterupt()
	fmt.Printf("\n")
}
