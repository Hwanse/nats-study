package main

import (
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

	subject := "sync.test"
	tick := time.Tick(time.Second * 3)

	for {
		select {
		case <-tick:
			err := conn.Publish(subject, []byte("hello!"))
			if err != nil {
				log.Fatalf("message publish error occurred : %+v", err)
				return
			}
		}
	}
}
