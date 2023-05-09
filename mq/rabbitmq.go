package rabbitmq

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

const (
	// QueueLog 日志
	QueueLog = "api4.0_queue_log"
	// ExchangeXDelayedMessage 延迟队列名
	ExchangeXDelayedMessage = "x-delayed-message"
)

type RabbitMQ struct {
	name            string
	logger          *log.Logger
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("rabbitMQ is shutting down")
)

// NewRabbitMQ creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewRabbitMQ(name string, addr string) *RabbitMQ {
	rabbitMQ := RabbitMQ{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		name:   name,
		done:   make(chan bool),
	}
	go rabbitMQ.handleReconnect(addr)
	return &rabbitMQ
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (rabbitMQ *RabbitMQ) handleReconnect(addr string) {
	for {
		rabbitMQ.isReady = false
		log.Println("Attempting to connect")

		conn, err := rabbitMQ.connect(addr)

		if err != nil {
			log.Println("Failed to connect. Retrying..., err:", err)

			select {
			case <-rabbitMQ.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := rabbitMQ.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (rabbitMQ *RabbitMQ) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	rabbitMQ.changeConnection(conn)
	log.Println("Connected!")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (rabbitMQ *RabbitMQ) handleReInit(conn *amqp.Connection) bool {
	for {
		rabbitMQ.isReady = false

		err := rabbitMQ.init(conn)

		if err != nil {
			log.Println("Failed to initialize channel. Retrying..., err:", err)

			select {
			case <-rabbitMQ.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-rabbitMQ.done:
			return true
		case <-rabbitMQ.notifyConnClose:
			log.Println("Connection closed. Reconnecting...")
			return false
		case <-rabbitMQ.notifyChanClose:
			log.Println("Channel closed. Re-running init...")
		}
	}
}

// init will initialize channel & declare queue
func (rabbitMQ *RabbitMQ) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	err = ch.Confirm(false)

	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		rabbitMQ.name,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		return err
	}

	rabbitMQ.changeChannel(ch)
	rabbitMQ.isReady = true
	log.Println("Setup!")

	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (rabbitMQ *RabbitMQ) changeConnection(connection *amqp.Connection) {
	rabbitMQ.connection = connection
	rabbitMQ.notifyConnClose = make(chan *amqp.Error)
	rabbitMQ.connection.NotifyClose(rabbitMQ.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (rabbitMQ *RabbitMQ) changeChannel(channel *amqp.Channel) {
	rabbitMQ.channel = channel
	rabbitMQ.notifyChanClose = make(chan *amqp.Error)
	rabbitMQ.notifyConfirm = make(chan amqp.Confirmation, 1)
	rabbitMQ.channel.NotifyClose(rabbitMQ.notifyChanClose)
	rabbitMQ.channel.NotifyPublish(rabbitMQ.notifyConfirm)
}

// Push will push data onto the queue, and wait for a confirm.
// If no confirms are received until within the resendTimeout,
// it continuously re-sends messages until a confirm is received.
// This will block until the server sends a confirm. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (rabbitMQ *RabbitMQ) Push(data []byte) error {
	if !rabbitMQ.isReady {
		return errors.New("failed to push push: not connected")
	}
	for {
		err := rabbitMQ.UnsafePush(data)
		if err != nil {
			rabbitMQ.logger.Println("Push failed. Retrying...")
			select {
			case <-rabbitMQ.done:
				return errShutdown
			case <-time.After(resendDelay):
			}
			continue
		}
		select {
		case confirm := <-rabbitMQ.notifyConfirm:
			if confirm.Ack {
				rabbitMQ.logger.Println("Push confirmed!")
				return nil
			}
		case <-time.After(resendDelay):
		}
		rabbitMQ.logger.Println("Push didn't confirm. Retrying...")
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// recieve the message.
func (rabbitMQ *RabbitMQ) UnsafePush(data []byte) error {
	if !rabbitMQ.isReady {
		return errNotConnected
	}
	return rabbitMQ.channel.PublishWithContext(
		context.Background(),
		"",            // Exchange
		rabbitMQ.name, // Routing key
		false,         // Mandatory
		false,         // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

// Stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (rabbitMQ *RabbitMQ) Stream() (<-chan amqp.Delivery, error) {
	if !rabbitMQ.isReady {
		return nil, errNotConnected
	}
	return rabbitMQ.channel.Consume(
		rabbitMQ.name,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
}

// Close will cleanly shutdown the channel and connection.
func (rabbitMQ *RabbitMQ) Close() error {
	if !rabbitMQ.isReady {
		return errAlreadyClosed
	}
	err := rabbitMQ.channel.Close()
	if err != nil {
		return err
	}
	err = rabbitMQ.connection.Close()
	if err != nil {
		return err
	}
	close(rabbitMQ.done)
	rabbitMQ.isReady = false
	return nil
}

// GetIsReady returns whether the queue is ready to be used.
func (rabbitMQ *RabbitMQ) GetIsReady() bool {
	return rabbitMQ.isReady
}
