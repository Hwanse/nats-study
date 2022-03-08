package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	subject := "request.test.2"
	tick := time.Tick(time.Second * 2)

	for {
		select {
		case <-tick:
			msg, err := conn.Request(subject, []byte("Hello"), time.Second*2)
			if err != nil {
				log.Fatalf("request message error occurred : %+v", err)
				return
			}

			reply := string(msg.Data)
			fmt.Println("receive response message :", reply)
		}
	}
}
