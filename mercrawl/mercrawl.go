package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/godriccao/mercrawl"
)

func main() {
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
