package repository

import (
	"context"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/google/uuid"
)

type Repository interface {
	StoreStructuralElement(
		ctx context.Context,
		bodyId uuid.UUID,
		element *models.StructuralElement,
	) error

	GetStructuralElements(
		ctx context.Context,
		bodyId uuid.UUID,
		cursor int,
		limit uint,
	) ([]models.StructuralElement, int, error)

	FindStructuralElementByID(
		ctx context.Context,
		seId int,
	) (models.StructuralElement, error)

	FindStructuralElementByIDAndBodyID(
		ctx context.Context,
		seId int,
		bodyId uuid.UUID,
	) (models.StructuralElement, error)

	UpdateStructuralElement(
		ctx context.Context,
		element models.StructuralElement,
	) error

	DeleteStructuralElement(
		ctx context.Context,
		element models.StructuralElement,
	) error

	StoreParagraphElement(
		ctx context.Context,
		paragraphId int,
		element *models.ParagraphElement,
	) error

	GetParagraphElement(
		ctx context.Context,
		paragraphId int,
		paragraphElementId int,
	) (models.ParagraphElement, error)

	DeleteParagraphElement(
		ctx context.Context,
		element models.ParagraphElement,
	) error

	UpdateParagraphElement(
		ctx context.Context,
		element models.ParagraphElement,
	) error

	GetParagraphElements(
		ctx context.Context,
		paragraphId int,
		cursor int,
		limit uint,
	) ([]models.ParagraphElement, int, error)
}
