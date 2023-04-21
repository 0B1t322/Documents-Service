package events

type EventHandler func(event Event)

type EventSubscriber interface {
	Subscribe(event Event, handler EventHandler)
}
