package dto

import "github.com/google/uuid"

type CreateElementDto struct {
	Index          int
	ElementType    ElementType
	ElementStyleID int
	BodyID         uuid.UUID
}
