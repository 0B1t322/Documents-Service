package events

type EventHandler func(event Event) error

type EventSubscriber interface {
	Subscribe(event Event, handler EventHandler)
}
