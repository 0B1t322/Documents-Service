package models

import "github.com/google/uuid"

type Body struct {
	ID       uuid.UUID
	Elements []StructuralElement
}
