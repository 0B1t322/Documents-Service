package dto

import "github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"

type CreateDocumentDto struct {
	Title string
	Size  models.Size
}
