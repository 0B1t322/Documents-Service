package amqp

import (
	"context"
	"github.com/0B1t322/Documents-Service/internal/core/events"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/ogen-go/ogen/json"
	"github.com/streadway/amqp"
	"time"
)

type EventsPublisher struct {
	ch           *amqp.Channel
	logger       log.Logger
	exchangeName string
}

func NewEventsPublisher(ch *amqp.Channel, logger log.Logger, exchangeName string) (*EventsPublisher, error) {
	err := ch.ExchangeDeclare(
		exchangeName, "topic", true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	return &EventsPublisher{
		ch:           ch,
		logger:       log.With(logger, "Service", "AMQPEventsPublisher"),
		exchangeName: exchangeName,
	}, nil
}

func (a EventsPublisher) PublishEvent(ctx context.Context, event events.Event) {
	body, err := json.Marshal(event)
	if err != nil {
		level.Error(a.logger).Log("description", "Failed to encode event", "err", err)
	}

	err = a.ch.Publish(
		a.exchangeName, event.Event(), false, false, amqp.Publishing{
			ContentType: "application/json",
			Timestamp:   time.Now(),
			Body:        body,
		},
	)
	if err != nil {
		level.Error(a.logger).Log("description", "Failed to publish event", "err", err)
	}
}
