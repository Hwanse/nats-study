package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	done := make(chan bool)

	// Connect to NATS
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic("nats connect failed")
	}
	defer conn.Close()

	log.Println("nats connect success")
	subject := "pub.1.req"

	// subscribe subject and receive message
	conn.Subscribe(subject, func(m *nats.Msg) {
		log.Println(string(m.Data))
		close(done)
	})
	// message publish to subject
	conn.Publish(subject, []byte("hello world!"))

	<- done
	log.Println("program end")
}