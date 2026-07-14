package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ manages the AMQP connection and provides publish/consume.
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
	mu      sync.Mutex
}

// Event represents a domain event published to RabbitMQ.
type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
}

// NewRabbitMQ creates a new RabbitMQ connection.
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}
	mq := &RabbitMQ{conn: conn, channel: ch, url: url}
	for _, ex := range []string{"crm.events", "pbx.events", "recording.events"} {
		if err := ch.ExchangeDeclare(ex, "topic", true, false, false, false, nil); err != nil {
			return nil, fmt.Errorf("declare exchange %s: %w", ex, err)
		}
	}
	slog.Info("RabbitMQ exchanges declared")
	return mq, nil
}

func (mq *RabbitMQ) Close() error {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if mq.channel != nil { mq.channel.Close() }
	if mq.conn != nil { return mq.conn.Close() }
	return nil
}

func (mq *RabbitMQ) Publish(ctx context.Context, exchange, routingKey string, event Event) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if event.Timestamp.IsZero() { event.Timestamp = time.Now().UTC() }
	body, err := json.Marshal(event)
	if err != nil { return fmt.Errorf("marshal event: %w", err) }
	return mq.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json", DeliveryMode: amqp.Persistent, Timestamp: event.Timestamp, Body: body,
	})
}

func (mq *RabbitMQ) Consume(queueName, exchange, routingKey string) (<-chan amqp.Delivery, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	q, err := mq.channel.QueueDeclare(queueName, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange": queueName + ".dlx",
	})
	if err != nil { return nil, fmt.Errorf("declare queue: %w", err) }
	if err := mq.channel.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
		return nil, fmt.Errorf("bind queue: %w", err)
	}
	mq.channel.Qos(10, 0, false)
	return mq.channel.Consume(q.Name, "", false, false, false, false, nil)
}
