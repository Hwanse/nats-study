package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	subject := "request.test.2"
	subscription, err := conn.Subscribe(subject, func(msg *nats.Msg) {
		log.Println("receive message : ", string(msg.Data))
		msg.Respond([]byte(fmt.Sprintf("%s World!", string(msg.Data))))
	})
	if err != nil {
		log.Fatalf("subscribe error occurred : %+v", err)
		return
	}
	defer subscription.Unsubscribe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	<-c

	log.Println("subscriber exit")
}
