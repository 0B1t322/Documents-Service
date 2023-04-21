package models

import "github.com/google/uuid"

type Style struct {
	ID             uuid.UUID
	Name           string
	ParagraphStyle ParagraphStyle
	TextStyle      TextStyle
}
