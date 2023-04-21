package models

import "github.com/google/uuid"

type Document struct {
	ID    uuid.UUID
	Title string
	Body  Body
	Style DocumentStyle
}
