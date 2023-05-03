package oas

import (
	"context"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/sessions"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	"github.com/google/uuid"
)

type (
	operationsService interface {
		SaveOperations(
			ctx context.Context, documentId uuid.UUID,
			operations []models.RevisionOperation,
		) error

		GetLastRevisionID(
			ctx context.Context,
			documentId uuid.UUID,
		) (int, error)

		GetDocumentHistory(
			ctx context.Context,
			documentId uuid.UUID,
		) ([]models.DocumentRevision, error)

		GetDocumentOperationsAfter(
			ctx context.Context,
			documentId uuid.UUID,
			afterRevisionId int,
		) ([]models.RevisionOperation, error)

		IsNotFound(err error) bool
	}
)

type OperationsController struct {
	operationsService operationsService
	mapper            mapper
}

func New(operationsService operationsService) *OperationsController {
	return &OperationsController{
		operationsService: operationsService,
		mapper:            mapper{},
	}
}

func (o OperationsController) SyncDocumentsById(
	ctx context.Context,
	params sessions.SyncDocumentsByIdParams,
) (sessions.SyncDocumentsByIdRes, error) {
	ops, err := o.operationsService.GetDocumentOperationsAfter(ctx, params.ID, params.RevId)
	if err != nil {
		return FailedToSyncDocument()
	}

	return o.mapper.DocumentsOperations(ops), nil
}

func (o OperationsController) GetDocumentRevisionId(
	ctx context.Context,
	params sessions.GetDocumentRevisionIdParams,
) (sessions.GetDocumentRevisionIdRes, error) {
	id, err := o.operationsService.GetLastRevisionID(ctx, params.ID)
	switch {
	case o.operationsService.IsNotFound(err):
		return NotFound(err)
	case err != nil:
		return FailedGetDocumentRevision()
	}

	return &sessions.Document{
		ID:         params.ID,
		RevisionId: id,
	}, nil
}

func (o OperationsController) PushOperationToDocument(
	ctx context.Context,
	req *sessions.SaveDocumentRequest,
	params sessions.PushOperationToDocumentParams,
) (sessions.PushOperationToDocumentRes, error) {
	err := o.operationsService.SaveOperations(ctx, params.ID, o.mapper.SaveOperations(req))
	switch {
	case o.operationsService.IsNotFound(err):
		return NotFound(err)
	case err != nil:
		return FailedToSaveOperations()
	}

	return &sessions.PushOperationToDocumentNoContent{}, nil
}

func (o OperationsController) GetDocumentHistory(
	ctx context.Context,
	params sessions.GetDocumentHistoryParams,
) (sessions.GetDocumentHistoryRes, error) {
	documentRevisions, err := o.operationsService.GetDocumentHistory(ctx, params.ID)
	if o.operationsService.IsNotFound(err) {
		return NotFound(err)
	} else if err != nil {
		return FailedToGetDocumentHistory()
	}

	return o.mapper.DocumentHistory(documentRevisions), nil
}
