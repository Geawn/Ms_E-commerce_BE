package event

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQService(host, port, user, password string) (*RabbitMQService, error) {
	// Tạo connection string
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	// Kết nối đến RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Tạo channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Khai báo các exchange và queue
	exchanges := []string{"product_events", "order_events", "inventory_events"}
	queues := []string{"product_created", "product_updated", "product_deleted", "order_created", "inventory_updated"}

	// Khai báo exchanges
	for _, exchange := range exchanges {
		err = ch.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to declare exchange %s: %v", exchange, err)
		}
	}

	// Khai báo queues và bindings
	for _, queue := range queues {
		_, err = ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to declare queue %s: %v", queue, err)
		}

		// Bind queue to appropriate exchange based on queue name
		var exchange string
		var routingKey string

		switch queue {
		case "product_created", "product_updated", "product_deleted":
			exchange = "product_events"
			routingKey = queue
		case "order_created":
			exchange = "order_events"
			routingKey = "order.created"
		case "inventory_updated":
			exchange = "inventory_events"
			routingKey = "inventory.updated"
		}

		err = ch.QueueBind(
			queue,      // queue name
			routingKey, // routing key
			exchange,   // exchange
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind queue %s to exchange %s: %v", queue, exchange, err)
		}
	}

	return &RabbitMQService{
		conn:    conn,
		channel: ch,
	}, nil
}

func (s *RabbitMQService) PublishEvent(exchange, routingKey string, message []byte) error {
	return s.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
}

func (s *RabbitMQService) ConsumeEvents(queue string, handler func([]byte) error) error {
	msgs, err := s.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				// Reject the message and requeue it
				d.Reject(true)
			} else {
				// Acknowledge the message
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (s *RabbitMQService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
