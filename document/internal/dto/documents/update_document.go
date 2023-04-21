package dto

import (
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/google/uuid"
)

type (
	UpdateDocumentDto struct {
		ID    uuid.UUID
		Title string
		Style UpdateDocumentStyleDto
	}

	UpdateDocumentStyleDto struct {
		PageSize models.Size
	}
)
