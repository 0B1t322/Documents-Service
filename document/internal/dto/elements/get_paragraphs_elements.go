package dto

import (
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/google/uuid"
)

type (
	GetParagraphsElementDto struct {
		Cursor              string
		BodyID              uuid.UUID
		StructuralElementID int
		Limit               uint
	}

	GetParagraphsElementResponse struct {
		Elements   []models.ParagraphElement
		NextCursor string
	}
)
