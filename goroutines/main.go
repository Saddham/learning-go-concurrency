package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	// Go routine
	go hello(&wg)

	wg.Wait()

	goodbye()

	// Make main goroutine wait for a second before shutting down
	// time.Sleep(1 * time.Second) // Never use sleep in prod
}

func hello(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Hello, world!")
}

func goodbye() {
	fmt.Println("Goodbye, world!")
}
