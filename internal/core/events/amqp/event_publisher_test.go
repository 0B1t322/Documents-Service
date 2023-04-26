package amqp

import (
	"context"
	zapLog "github.com/go-kit/kit/log/zap"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

type Event struct {
	Payload struct {
		ID   int
		Name string
	}
}

func (Event) Event() string {
	return "test.event"
}

func TestFunc_EventsPublisher(t *testing.T) {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	require.NoError(t, err)

	ch, err := conn.Channel()
	require.NoError(t, err)
	logger := zapLog.NewZapSugarLogger(zap.L(), zap.DebugLevel)
	p, err := NewEventsPublisher(ch, logger, "documents-service.events")
	require.NoError(t, err)

	ch.QueueDeclare("documents.events.for.sessions", true, false, false, false, nil)

	require.NoError(t, ch.QueueBind("documents.events.for.sessions", "#", "documents-service.events", false, nil))

	p.PublishEvent(
		context.Background(), Event{
			Payload: struct {
				ID   int
				Name string
			}{ID: 1, Name: "s"},
		},
	)
}
