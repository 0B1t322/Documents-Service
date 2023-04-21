package events

import "github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"

type ElementsEvent string

func (e ElementsEvent) String() string {
	return string(e)
}

const (
	StructuralElementCreated = ElementsEvent("structuralElement.created")
	StructuralElementDeleted = ElementsEvent("structuralElement.deleted")
	StructuralElementUpdated = ElementsEvent("structuralElement.updated")
)

type StructuralElementCreatedEvent struct {
	StructuralElement models.StructuralElement
}

func (StructuralElementCreatedEvent) Event() string {
	return StructuralElementCreated.String()
}

type StructuralElementUpdatedEvent struct {
	Was    models.StructuralElement
	Become models.StructuralElement
}

func (StructuralElementUpdatedEvent) Event() string {
	return StructuralElementUpdated.String()
}

type StructuralElementDeletedEvent struct {
	StructuralElement models.StructuralElement
}

func (StructuralElementDeletedEvent) Event() string {
	return StructuralElementDeleted.String()
}
