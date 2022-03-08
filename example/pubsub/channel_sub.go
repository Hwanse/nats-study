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
	log.Printf("subscriber[%d] running.. \n", subscriberNo)

	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	subject := "async.test"

	// asynchronous subscription, receive message via channel
	// buffered 채널을 통해 message를 수신할 경우, channel buffer 사이즈에 주의해야한다.
	// 서버에서 pub 메세지 속도에 비해 subscriber의 메세지 처리속도가 느릴 경우 해당 채널에 메세지가 쌓이게 되며
	// 서버에서는 해당 subscriber client는 slow consumer가 되어 message drop 현상을 유발하거나, server와 connection이 강제로 끊긴다
	msgChannel := make(chan *nats.Msg, 10)
	subscription, err := conn.ChanSubscribe(subject, msgChannel)
	if err != nil {
		log.Fatalf("%s < subscription failed : %+v", subject, err)
		return
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case msg, _ := <-msgChannel:
			message := string(msg.Data)
			fmt.Println("received message :", message)
			if message == "exit" {
				return
			}
		}
	}

	log.Println("subscriber end..")
}
