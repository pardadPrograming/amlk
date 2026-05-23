package events

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(ctx context.Context, event Event) error
	Close() error
}

type Event struct {
	Type        string      `json:"type"`
	AggregateID string      `json:"aggregateId"`
	OccurredAt  time.Time   `json:"occurredAt"`
	Payload     interface{} `json:"payload,omitempty"`
}

const SearchInvalidateEvent = "search.invalidate"

type SearchInvalidatePayload struct {
	BusinessID string `json:"businessId"`
	UserID     string `json:"userId"`
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(ctx context.Context, event Event) error { return nil }
func (NoopPublisher) Close() error                                   { return nil }

type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	logger   *slog.Logger
}

func NewRabbitMQPublisher(url, exchange string, logger *slog.Logger) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if exchange == "" {
		exchange = "amlak.events"
	}
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	return &RabbitMQPublisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		logger:   logger,
	}, nil
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, event Event) error {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.channel.PublishWithContext(
		ctx,
		p.exchange,
		event.Type,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    event.OccurredAt,
			Type:         event.Type,
			Body:         body,
		},
	)
}

func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func SafePublish(ctx context.Context, logger *slog.Logger, publisher Publisher, event Event) {
	if publisher == nil {
		return
	}
	if err := publisher.Publish(ctx, event); err != nil && logger != nil {
		logger.Warn("event publish failed", "type", event.Type, "aggregate_id", event.AggregateID, "error", err)
	}
}
