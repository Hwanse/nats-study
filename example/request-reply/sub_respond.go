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

	subject := "request.test.1"
	subscription, err := conn.SubscribeSync(subject)
	if err != nil {
		log.Fatalf("subscribe error occurred : %+v", err)
		return
	}
	defer subscription.Unsubscribe()

	tick := time.Tick(time.Second * 2)

	for {
		select {
		case <-tick:
			msg, err := subscription.NextMsg(time.Second * 2)
			if err != nil {
				log.Fatalf("receive message error occurred : %+v", err)
				return
			}

			msg.Respond([]byte(fmt.Sprintf("%s World!", string(msg.Data))))
		}
	}
}
