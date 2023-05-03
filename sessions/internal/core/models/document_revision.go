package models

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations/elements"
	"github.com/google/uuid"
)

type (
	DocumentRevision struct {
		DocumentID uuid.UUID
		RevisionID int
		Operations []RevisionOperation
	}

	RevisionOperation struct {
		Insert *elements.InsertStructuralElement `json:",omitempty"`
		Delete *elements.DeleteStructuralElement `json:",omitempty"`
	}
)

func (r RevisionOperation) ApplyToParagraph(to *ParagraphElement) {
	if r.Insert != nil {
		r.Insert.InsertParagraphElement.Apply(to)
		return
	} else if r.Delete != nil {
		r.Delete.DeleteParagraphElement.Apply(to)
		return
	}
}

func (r RevisionOperation) StructuralElementIndex() int {
	switch r.Type() {
	case RevisionOperationTypeInsert:
		return r.Insert.StructuralElementIndex
	case RevisionOperationTypeDelete:
		return r.Delete.StructuralElementIndex
	}

	panic("Unknown type")
}

func (r RevisionOperation) ParagraphElementIndex() int {
	switch r.Type() {
	case RevisionOperationTypeInsert:
		return r.Insert.ParagraphIndex
	case RevisionOperationTypeDelete:
		return r.Delete.ParagraphIndex
	}

	panic("Unknown type")
}

func (r RevisionOperation) LocalIndex() int {
	switch r.Type() {
	case RevisionOperationTypeInsert:
		return r.Insert.Index
	case RevisionOperationTypeDelete:
		return r.Delete.Index
	}

	panic("Unknown type")
}

type RevisionOperationType string

const (
	RevisionOperationTypeInsert = "insert"
	RevisionOperationTypeDelete = "delete"
)

func (r RevisionOperation) Type() RevisionOperationType {
	if r.Insert != nil {
		return RevisionOperationTypeInsert
	} else if r.Delete != nil {
		return RevisionOperationTypeDelete
	}

	panic("Unknown type")
}
