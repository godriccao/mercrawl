package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/godriccao/mercrawl"
)

func main() {
	search := os.Args[1]
	go mercrawl.Start(search)
	pend()
	fmt.Printf("\n")
}

// pend wait an interupt signal ^c to end the program
func pend() {
	var end sync.WaitGroup
	endSignal := make(chan os.Signal, 1)
	signal.Notify(endSignal, os.Interrupt)

	end.Add(1)
	go func() {
		<-endSignal
		end.Done()
	}()

	end.Wait()
}
