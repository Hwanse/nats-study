package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	subscriberNo := rand.Intn(10)
	log.Printf("sync subscriber[%d] running.. \n", subscriberNo)

	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	subject := "sync.test"

	// synchronous subscription
	subscription, err := conn.SubscribeSync(subject)
	if err != nil {
		log.Fatalf("%s < subscription failed : %+v", subject, err)
		return
	}
	defer subscription.Unsubscribe()

	for {
		// 인자로 넘겨준 시간만큼 메세지 수신을 기대, 시간이 지나면 error return
		msg, err := subscription.NextMsg(time.Second * 4)
		if err != nil {
			log.Fatalf("receive message error : %+v", err)
			return
		}

		message := string(msg.Data)
		fmt.Println("received message :", fmt.Sprintf("%s World!", message))
	}

}
