package rabbitmq

import (
	"fmt"
	"log"
	"time"
)

func ExampleRabbitMQ_Push() {
	name := "job_queue"
	addr := "amqp://guest:guest@localhost:5672/"
	queue := NewRabbitMQ(name, addr)
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
	queue := NewRabbitMQ(name, addr)
	for {
		msgs, err := queue.Stream()
		if !queue.GetIsReady() {
			msgs, err = queue.Stream()
		}
		if err != nil {
			log.Println("Failed to register a consumer, ", err)
			time.Sleep(time.Second * 5)
			continue
		}

		for d := range msgs {
			fmt.Printf("Received message: %s\n", d.Body)
		}
	}
}
