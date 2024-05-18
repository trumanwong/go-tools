package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"testing"
	"time"
)

func TestNewRabbitMQ(t *testing.T) {
	arguments := make(amqp.Table)
	arguments["x-max-priority"] = int64(9)
	NewRabbitMQ(&Options{
		Name:          os.Getenv("RABBITMQ_QUEUE_NAME"),
		Addr:          os.Getenv("RABBITMQ_ADDR"),
		PrefetchCount: 1,
		PrefetchSize:  0,
		Global:        false,
		Consume: func(msgs <-chan amqp.Delivery) {
			for d := range msgs {
				fmt.Println("receive data: ", string(d.Body))
				time.Sleep(10 * time.Second)
				_ = d.Ack(false)
			}
		},
		Arguments: arguments,
	})
	select {}
}

func TestRabbitMQ_PushV2(t *testing.T) {
	arguments := make(amqp.Table)
	arguments["x-max-priority"] = int64(9)
	queue := NewRabbitMQ(&Options{
		Name:          os.Getenv("RABBITMQ_QUEUE_NAME"),
		Addr:          os.Getenv("RABBITMQ_ADDR"),
		PrefetchCount: 0,
		PrefetchSize:  0,
		Global:        false,
		Consume:       nil,
		Arguments:     arguments,
	})
	messages := []string{"Hello, World!", "Hello, RabbitMQ!", "Hello, Go!"}
	// Attempt to push a message every 2 seconds
	for _, message := range messages {
		for {
			if err := queue.PushV2(amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
				Priority:    0,
			}); err != nil {
				fmt.Printf("Push failed: %s\n", err)
				time.Sleep(time.Second * 3)
			} else {
				fmt.Println("Push succeeded!")
				break
			}
		}
	}
	for {
		if err := queue.PushV2(amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Test Priority"),
			Priority:    9,
		}); err != nil {
			time.Sleep(time.Second * 3)
			fmt.Printf("Push failed: %s\n", err)
		} else {
			fmt.Println("Push succeeded!")
			break
		}
	}
}
