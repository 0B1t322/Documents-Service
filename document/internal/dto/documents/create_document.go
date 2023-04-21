package dto

import "github.com/0B1t322/Documents-Service/document/internal/core/models"

type CreateDocumentDto struct {
	Title string
	Size  models.Size
}
