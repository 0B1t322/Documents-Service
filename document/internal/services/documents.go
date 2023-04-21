package services

import (
	"context"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/events"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	dto "github.com/0B1t322/Online-Document-Redactor/document/internal/dto/documents"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/repository/documents"
	"github.com/0B1t322/Online-Document-Redactor/internal/core/utils/cursor"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
)

type (
	documentRepository interface {
		FindByID(ctx context.Context, id uuid.UUID) (models.Document, error)
		Store(ctx context.Context, document *models.Document) error
		Update(ctx context.Context, document models.Document) error
		Get(
			ctx context.Context,
			cursor uuid.UUID,
			limit uint64,
		) ([]models.Document, uuid.UUID, error)
	}
)

type DocumentService struct {
	documentRepository documentRepository
	eventPublisher     events.EventPublisher
	logger             log.Logger
}

func NewDocumentService(
	documentRepository documentRepository,
	eventPublisher events.EventPublisher,
	logger log.Logger,
) *DocumentService {
	return &DocumentService{
		documentRepository: documentRepository,
		eventPublisher:     eventPublisher,
		logger:             log.WithPrefix(logger, "service", "DocumentService"),
	}
}

func (d DocumentService) CreateDocument(
	ctx context.Context,
	req dto.CreateDocumentDto,
) (models.Document, error) {
	document := models.Document{
		Title: req.Title,
		Style: models.DocumentStyle{
			PageSize: req.Size,
		},
	}

	if err := d.documentRepository.Store(ctx, &document); err != nil {
		level.Error(d.logger).Log("Failed to save document", err)

		return models.Document{}, err
	}

	d.eventPublisher.PublishEvent(ctx, events.DocumentCreatedEvent{DocumentID: document.ID})

	return document, nil
}

func (d DocumentService) GetDocument(ctx context.Context, documentId uuid.UUID) (models.Document, error) {
	document, err := d.documentRepository.FindByID(ctx, documentId)
	if err == repository.ErrNotFound {
		return models.Document{}, ErrDocumentNotFound
	} else if err != nil {
		level.Error(d.logger).Log("Failed to get document", err)
		return models.Document{}, err
	}

	return document, nil
}

func (d DocumentService) UpdateDocument(ctx context.Context, req dto.UpdateDocumentDto) (models.Document, error) {
	document, err := d.GetDocument(ctx, req.ID)
	if err != nil {
		return models.Document{}, err
	}

	document.Title = req.Title
	document.Style.PageSize = req.Style.PageSize

	if err := d.documentRepository.Update(ctx, document); err != nil {
		level.Error(d.logger).Log("Failed to update document", err)
		return models.Document{}, err
	}

	return document, nil
}

func (d DocumentService) GetDocuments(
	ctx context.Context,
	req dto.GetDocumentsDto,
) (dto.GetDocumentsResponse, error) {
	var cur uuid.UUID = uuid.Nil
	{
		if req.Cursor != "" {
			if decoded, err := cursor.CursorToUUID(req.Cursor); err != nil {
				return dto.GetDocumentsResponse{}, err
			} else {
				cur = decoded
			}
		}
	}

	var limit uint64 = req.Limit

	documents, nextCursor, err := d.documentRepository.Get(ctx, cur, limit)
	if err != nil {
		return dto.GetDocumentsResponse{}, err
	}

	return dto.GetDocumentsResponse{
		Document: documents,
		Cursor:   cursor.UUIDToCursor(nextCursor),
	}, nil
}
