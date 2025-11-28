package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/todoist/backend/pkg/events"
)

// RabbitMQPublisher implements event publishing using RabbitMQ
type RabbitMQPublisher struct {
	conn         *amqp091.Connection
	channel      *amqp091.Channel
	exchangeName string
}

// NewRabbitMQPublisher creates a new RabbitMQ event publisher
func NewRabbitMQPublisher(url, exchangeName string) (*RabbitMQPublisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare topic exchange
	err = channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQPublisher{
		conn:         conn,
		channel:      channel,
		exchangeName: exchangeName,
	}, nil
}

// Publish publishes an event to RabbitMQ
func (p *RabbitMQPublisher) Publish(ctx context.Context, event interface{}) error {
	routingKey := p.getRoutingKey(event)

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchangeName, // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Close closes the RabbitMQ connection
func (p *RabbitMQPublisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}

func (p *RabbitMQPublisher) getRoutingKey(event interface{}) string {
	switch e := event.(type) {
	case events.UserRegistered:
		return e.EventType
	case events.UserLoggedIn:
		return e.EventType
	default:
		return "unknown"
	}
}
