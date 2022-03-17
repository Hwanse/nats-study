package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"math/rand"
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
	tick := time.Tick(time.Second * 1)

	for {
		select {
		case <-tick:
			random := rand.Int()
			jobID := fmt.Sprintf("[id=%d]", random)
			fmt.Println("send job -", jobID)
			msg, err := conn.Request(subject, []byte(jobID), time.Second*1)
			if err != nil {
				log.Fatalf("request message error occurred : %+v", err)
				return
			}
			reply := string(msg.Data)
			fmt.Println("receive worker-response message :", reply)
		}
	}
}
