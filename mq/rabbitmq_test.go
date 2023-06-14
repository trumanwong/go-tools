package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func ExampleRabbitMQ_Push() {
	name := "job_queue"
	addr := "amqp://guest:guest@localhost:5672/"
	queue := NewRabbitMQ(name, addr, nil)
	message := []byte("message")
	// Attempt to push a message every 2 seconds
	for {
		time.Sleep(time.Second * 3)
		if err := queue.Push(message); err != nil {
			fmt.Printf("Push failed: %s\n", err)
		} else {
			fmt.Println("Push succeeded!")
		}
	}
}

func ExampleRabbitMQ_Stream() {
	name := "job_queue"
	addr := "amqp://guest:guest@localhost:5672/"
	NewRabbitMQ(name, addr, func(deliveries <-chan amqp.Delivery) {
		// Consume messages
	})
}
