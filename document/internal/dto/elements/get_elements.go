package dto

import (
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/google/uuid"
)

type (
	GetElementsDto struct {
		Cursor string
		BodyID uuid.UUID
		Limit  uint
	}

	GetElementsResponse struct {
		Elements   []models.StructuralElement
		NextCursor string
	}
)
