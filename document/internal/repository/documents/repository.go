package repository

import (
	"context"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (models.Document, error)
	Store(ctx context.Context, document *models.Document) error
	Update(ctx context.Context, document models.Document) error
	Get(
		ctx context.Context,
		cursor uuid.UUID,
		limit uint64,
	) ([]models.Document, uuid.UUID, error)
}
