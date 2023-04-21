package events

import "github.com/google/uuid"

type DocumentEvent string

func (e DocumentEvent) String() string {
	return string(e)
}

const (
	DocumentCreated = DocumentEvent("document.created")
	DocumentDeleted = DocumentEvent("document.deleted")
)

type DocumentCreatedEvent struct {
	DocumentID uuid.UUID
}

func (DocumentCreatedEvent) Event() string {
	return DocumentCreated.String()
}

type DocumentDeletedEvent struct {
	DocumentID uuid.UUID
}

func (DocumentDeletedEvent) Event() string {
	return DocumentDeleted.String()
}
