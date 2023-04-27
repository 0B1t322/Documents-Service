package amqp

import (
	"fmt"
	"github.com/0B1t322/Documents-Service/internal/core/events"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/ogen-go/ogen/json"
	"github.com/streadway/amqp"
	"reflect"
)

type eventUnmarshaller func(delivery amqp.Delivery) (events.Event, error)

type EventsSubscriber struct {
	ch             *amqp.Channel
	logger         log.Logger
	queue          amqp.Queue
	exchangeName   string
	queueName      string
	serviceSubName string
	consumer       <-chan amqp.Delivery
	handlers       map[string][]events.EventHandler
	unmarshals     map[string]eventUnmarshaller
	close          chan int
}

func NewEventsSubscriber(
	ch *amqp.Channel, logger log.Logger, exchangeName string,
	serviceSubName string,
) (*EventsSubscriber, error) {
	queueName := fmt.Sprintf("%s.subscribers", serviceSubName)
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	consumer, err := ch.Consume(queueName, serviceSubName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	e := &EventsSubscriber{
		ch:             ch,
		logger:         logger,
		queue:          q,
		exchangeName:   exchangeName,
		queueName:      queueName,
		serviceSubName: serviceSubName,
		consumer:       consumer,
		handlers:       map[string][]events.EventHandler{},
		unmarshals:     map[string]eventUnmarshaller{},
		close:          make(chan int),
	}

	go e.Consume()

	return e, nil
}

func (e EventsSubscriber) Close() error {
	e.close <- 1
	close(e.close)
	return e.ch.Close()
}

func (e EventsSubscriber) Subscribe(event events.Event, handler events.EventHandler) {
	if err := e.ch.QueueBind(e.queueName, event.Event(), e.exchangeName, false, nil); err != nil {
		level.Error(e.logger).Log("description", "failed to bind queue to exchange", "err", err)
	}

	e.addHandler(event, handler)
	e.addUnmarshaller(event)
}

func (e EventsSubscriber) addHandler(event events.Event, handler events.EventHandler) {
	handlers, find := e.handlers[event.Event()]
	if !find {
		handlers = make([]events.EventHandler, 0)
	}

	handlers = append(handlers, handler)
	e.handlers[event.Event()] = handlers
}

func (e EventsSubscriber) addUnmarshaller(event events.Event) {
	_, find := e.unmarshals[event.Event()]
	if find {
		return
	}

	unmarshaller := func(delivery amqp.Delivery) (events.Event, error) {
		typ := reflect.TypeOf(event)

		v := reflect.New(typ)
		value := v.Interface()

		if err := json.Unmarshal(delivery.Body, &value); err != nil {
			return nil, err
		}

		return value.(events.Event), nil
	}

	e.unmarshals[event.Event()] = unmarshaller
}

func (e EventsSubscriber) HandleMessage(msg amqp.Delivery) {
	//defer msg.Reject(true)
	key := msg.RoutingKey
	unmarshaller, find := e.unmarshals[key]
	if !find {
		return
	}

	event, err := unmarshaller(msg)
	if err != nil {
		level.Error(e.logger).Log("description", "Failed to unmarshall event", "err", err, "msg", msg)
		return
	}

	handlers := e.handlers[key]
	if len(handlers) == 0 {
		return
	}

	for _, handler := range handlers {
		if err := handler(event); err != nil {
			return
		}
	}

	msg.Ack(true)
}

func (e EventsSubscriber) Consume() {
	for {
		select {
		case msg := <-e.consumer:
			go e.HandleMessage(msg)
		case <-e.close:
			return
		}
	}
}
