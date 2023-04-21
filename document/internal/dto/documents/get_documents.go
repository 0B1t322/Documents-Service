package dto

import "github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"

type (
	GetDocumentsDto struct {
		// Can be empty
		Cursor string
		// Default is 10
		Limit uint64
	}

	GetDocumentsResponse struct {
		Document []models.Document
		Cursor   string
	}
)
