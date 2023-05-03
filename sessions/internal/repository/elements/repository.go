package repository

import (
	"context"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetParagraphElementByIndexes(
		ctx context.Context, documentId uuid.UUID, seId,
		peId int,
	) (models.ParagraphElement, error)

	UpdateParagraphElementByIndexes(
		ctx context.Context, documentId uuid.UUID,
		seId, peId int,
		element models.ParagraphElement,
	) error

	IsNotFound(err error) bool
}
