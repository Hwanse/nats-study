package main

import (
	"bufio"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

func main() {
	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	subject := "subject.1.test"

	requireInputMessage()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		inputStr := scanner.Text()
		conn.Publish(subject, []byte(inputStr))
		requireInputMessage()
	}
}

func requireInputMessage() {
	fmt.Printf("input message :\t")
}
