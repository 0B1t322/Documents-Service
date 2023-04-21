package events

import (
	"context"
)

type EventPublisher interface {
	PublishEvent(ctx context.Context, event Event)
}
