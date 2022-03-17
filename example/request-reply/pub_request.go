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
	tick := time.Tick(time.Second * 2)

	for {
		select {
		case <-tick:
			inbox := nats.NewInbox()
			subscription, err := conn.SubscribeSync(inbox)
			if err != nil {
				log.Fatalf("inbox subscribe error occurred : %+v", err)
				return
			}
			defer subscription.Unsubscribe()

			err = conn.PublishRequest(subject, inbox, []byte("Hello"))
			if err != nil {
				log.Fatalf("message publish error occurred : %+v", err)
				return
			}

			// expect respond message
			respMsg, err := subscription.NextMsg(time.Second * 1)
			if err != nil {
				log.Fatalf("receive response message error : %+v", err)
				return
			}
			fmt.Println("received respond message : ", string(respMsg.Data))
		}
	}

}
