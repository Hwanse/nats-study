package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
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
	subject := "queue.pubsub.work"
	tick := time.Tick(time.Second * 1)

	for {
		select {
		case <-tick:
			random := rand.Int()
			jobID := fmt.Sprintf("[id=%d]", random)
			fmt.Println("send job -", jobID)
			conn.Publish(subject, []byte(jobID))

		}
	}
}
