package main

import (
	"fmt"
	"time"
)

var hellos = []string{"Hello!", "Ciao!", "Hola!", "Hej!", "Salut!"}
var goodbyes = []string{"Goodbye!", "Arrivederci!", "Adios!", "Hej Hej!", "La revedere!"}

func main() {
	// Unbuffered channel
	// ch := make(chan string)

	// Buffered channel
	ch := make(chan string, 1)
	ch2 := make(chan string, 1)

	go greet(hellos, ch)
	go greet(goodbyes, ch2)

	time.Sleep(5 * time.Second)

	fmt.Println("Main ready!")

	//receiveGreetingUsingForLoop(ch)
	//receiveGreetingUsingForRange(ch)
	receiveGreetingUsingSelect(ch, ch2)
}

func receiveGreetingUsingForLoop(ch <-chan string) {
	// Using for loop, we need to close the channel manually
	for {
		greeting, ok := <-ch
		if !ok {
			return
		}

		printGreeting(greeting)
	}
}

func receiveGreetingUsingForRange(ch <-chan string) {
	// Range automatically exits once channel is closed
	for greeting := range ch {
		printGreeting(greeting)
	}
}

func receiveGreetingUsingSelect(ch <-chan string, ch2 <-chan string) {
	for {
		select {
		case gr, ok := <-ch:
			if !ok {
				ch = nil
				break
			}

			printGreeting(gr)
		case gr2, ok := <-ch2:
			if !ok {
				ch2 = nil
				break
			}

			printGreeting(gr2)
		default:
			return
		}
	}
}

func greet(greetings []string, ch chan<- string) {
	fmt.Printf("Greeter ready!\nGreeter waiting to send greeting...\n")

	for _, greeting := range greetings {
		ch <- greeting
	}

	close(ch)

	fmt.Println("Greeter completed!")
}

func printGreeting(greeting string) {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Greeting received!", greeting)
}
