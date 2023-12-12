package main

import (
	"fmt"
	"strconv"
	"time"
)

var hellos = []string{"Hello!", "Ciao!", "Hola!", "Hej!", "Salut!"}
var goodbyes = []string{"Goodbye!", "Arrivederci!", "Adios!", "Hej Hej!", "La revedere!"}

func main() {

	// Buffered channel
	ch := make(chan string, 1)
	ch2 := make(chan string, 1)
	done := make(chan struct{})
	done2 := make(chan struct{})

	go greetWithDone(hellos, ch, done)
	go greetWithDone(goodbyes, ch2, done2)

	time.Sleep(5 * time.Second)

	fmt.Println("Main ready!")

	receiveGreetingUsingDoneChannel(1, ch, done)
	receiveGreetingUsingDoneChannel(2, ch2, done2)
}

// "Signalling work is done" concurrency pattern
func receiveGreetingUsingDoneChannel(chanNumber int, ch <-chan string, done <-chan struct{}) {
	for {
		select {
		case gr := <-ch:
			printGreeting(gr)
		case <-done:
			fmt.Printf("Done receiving greetings from channel %s!\n", strconv.Itoa(chanNumber))
			return
		}
	}
}

// "Signalling work is done" concurrency pattern
func greetWithDone(greetings []string, ch chan<- string, done chan<- struct{}) {
	fmt.Printf("Greeter ready!\nGreeter waiting to send greeting...\n")

	for _, greeting := range greetings {
		ch <- greeting
	}

	close(done)

	fmt.Println("Greeter completed!")
}

func printGreeting(greeting string) {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Greeting received!", greeting)
}
