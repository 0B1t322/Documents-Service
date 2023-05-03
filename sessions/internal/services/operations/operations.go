package operations

import (
	"context"
	"fmt"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	"github.com/google/uuid"
)

type (
	elementsRepository interface {
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

	operationRepository interface {
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
)

type Service struct {
	elementsRepository  elementsRepository
	operationRepository operationRepository
}

func New(elementsRepository elementsRepository, operationRepository operationRepository) *Service {
	return &Service{
		elementsRepository:  elementsRepository,
		operationRepository: operationRepository,
	}
}

func (s Service) SaveOperations(
	ctx context.Context, documentId uuid.UUID,
	operations []models.RevisionOperation,
) error {
	elements := createOperationsAndElement(operations)
	err := elements.getElements(ctx, documentId, s.elementsRepository)
	if s.elementsRepository.IsNotFound(err) {
		return ErrParagraphElementNotFound
	} else if err != nil {
		fmt.Println("failed createOperationsAndElement", err)
		return err
	}

	elements.applyOperations()

	if err := elements.updateElements(ctx, documentId, s.elementsRepository); err != nil {
		fmt.Println("failed updateElements", err)
		return err
	}

	if _, err := s.operationRepository.SaveRevision(ctx, documentId, operations); err != nil {
		fmt.Println("failed SaveRevision", err)
		return err
	}

	return nil
}

func (s Service) GetLastRevisionID(
	ctx context.Context,
	documentId uuid.UUID,
) (int, error) {
	revisionID, err := s.operationRepository.GetLastRevision(ctx, documentId)
	if err != nil {
		return 0, err
	} else if revisionID == -1 {
		return 0, ErrNotFoundRevision
	}

	return revisionID, nil
}

func (s Service) GetDocumentHistory(
	ctx context.Context,
	documentId uuid.UUID,
) ([]models.DocumentRevision, error) {
	revisions, err := s.operationRepository.GetRevisions(ctx, documentId)
	if s.operationRepository.IsNotFound(err) {
		return revisions, ErrNotFoundRevisionsForDocument
	} else if err != nil {
		return revisions, err
	}

	return revisions, nil
}

func (s Service) GetDocumentOperationsAfter(
	ctx context.Context,
	documentId uuid.UUID,
	afterRevisionId int,
) ([]models.RevisionOperation, error) {
	revs, err := s.operationRepository.GetRevisionsAfter(ctx, documentId, afterRevisionId)
	if s.operationRepository.IsNotFound(err) {
		return []models.RevisionOperation{}, nil
	} else if err != nil {
		return nil, err
	}

	var revOps []models.RevisionOperation
	{
		for _, rev := range revs {
			revOps = append(revOps, rev.Operations...)
		}
	}

	return revOps, nil
}
