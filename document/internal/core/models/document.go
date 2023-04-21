package model

import "github.com/google/uuid"

type Document struct {
	ID    uuid.UUID
	Title string
	Style DocumentStyle
	Body  Body
}
