package repository

import (
	"context"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	"github.com/google/uuid"
)

type Repository interface {
	SaveRevision(ctx context.Context, documentId uuid.UUID, operations []models.RevisionOperation) (int, error)
	GetLastRevision(ctx context.Context, documentId uuid.UUID) (int, error)
	GetRevisions(ctx context.Context, documentId uuid.UUID) ([]models.DocumentRevision, error)

	GetRevisionsAfter(
		ctx context.Context,
		documentId uuid.UUID,
		revId int,
	) ([]models.DocumentRevision, error)

	IsNotFound(err error) bool
}
