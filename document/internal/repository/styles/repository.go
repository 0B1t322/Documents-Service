package repository

import (
	"context"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	"github.com/google/uuid"
)

type Repository interface {
	StoreStyleInDocument(ctx context.Context, documentId uuid.UUID, style *models.Style) error
	FindStyleInDocumentByID(ctx context.Context, documentId uuid.UUID, styleId uuid.UUID) (models.Style, error)
	UpdateStyle(ctx context.Context, style models.Style) error
	GetAllStylesInDocument(ctx context.Context, documentId uuid.UUID) ([]models.Style, error)
	DeleteStyleInDocument(ctx context.Context, documentId uuid.UUID, style models.Style) error
}
