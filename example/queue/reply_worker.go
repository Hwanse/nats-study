package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {

	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	rand.Seed(time.Now().UnixNano())
	subject := "queue.req.reply.work"

	subscription, err := conn.QueueSubscribe(subject, "queue.worker.group.2", func(msg *nats.Msg) {
		fmt.Println("receive job :", string(msg.Data))
		msg.Respond([]byte(fmt.Sprintf("job end - %s", string(msg.Data))))
	})
	if err != nil {
		log.Fatalf("subscribe error occurred : %+v", err)
		return
	}
	defer subscription.Unsubscribe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}