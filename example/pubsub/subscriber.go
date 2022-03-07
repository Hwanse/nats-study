package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"math/rand"
	"sync"
	"time"
)

const EXIT = "exit"

func main() {
	rand.Seed(time.Now().UnixNano())
	subscriberNo := rand.Intn(10)
	log.Printf("subscriber[%d] running.. \n", subscriberNo)

	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	subject := "subject.1.test"
	wg := sync.WaitGroup{}

	wg.Add(1)

	// asynchronous subscription
	conn.Subscribe(subject, func(msg *nats.Msg) {
		message := string(msg.Data)
		printReceivedMessage(message)
		if message == EXIT {
			wg.Done()
		}
	})

	wg.Wait()
	log.Println("subscriber end..")
}

func printReceivedMessage(message string) {
	fmt.Println("received message :", message)
}
