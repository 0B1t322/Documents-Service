package oas

import (
	"context"
	dto "github.com/0B1t322/Online-Document-Redactor/document/internal/dto/documents"

	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/0B1t322/Online-Document-Redactor/pkg/gen/open-api/documents"
	"github.com/google/uuid"
)

type documentsService interface {
	CreateDocument(
		ctx context.Context,
		req dto.CreateDocumentDto,
	) (models.Document, error)

	GetDocument(ctx context.Context, documentId uuid.UUID) (models.Document, error)

	UpdateDocument(ctx context.Context, req dto.UpdateDocumentDto) (models.Document, error)

	GetDocuments(ctx context.Context, req dto.GetDocumentsDto) (dto.GetDocumentsResponse, error)

	IsNotFound(err error) bool

	IsParametersNotValid(err error) bool
}

type DocumentsController struct {
	service documentsService
	mapper  mapper
}

func NewDocumentController(service documentsService) *DocumentsController {
	return &DocumentsController{
		service: service,
		mapper:  mapper{},
	}
}

func (d DocumentsController) CreateDocument(
	ctx context.Context,
	req *documents.CreateUpdateDocumentView,
) (documents.CreateDocumentRes, error) {
	document, err := d.service.CreateDocument(ctx, d.mapper.CreateDocumentReq(req))
	if err != nil {
		return FailedToCreateDocument()
	}

	return d.mapper.Document(document), nil
}

func (d DocumentsController) DocumentsGet(
	ctx context.Context,
	params documents.DocumentsGetParams,
) (documents.DocumentsGetRes, error) {
	res, err := d.service.GetDocuments(ctx, d.mapper.GetDocumentsReq(params))
	if d.service.IsParametersNotValid(err) {
		return BadRequest(err)
	} else if err != nil {
		return FailedToGetDocuments()
	}

	return d.mapper.PaginatedDocuments(res), nil
}

func (d DocumentsController) GetDocumentById(
	ctx context.Context,
	params documents.GetDocumentByIdParams,
) (documents.GetDocumentByIdRes, error) {
	document, err := d.service.GetDocument(ctx, params.ID)
	if d.service.IsNotFound(err) {
		return DocumentNotFound(params.ID)
	} else if err != nil {
		return FailedToGetDocument()
	}

	return d.mapper.Document(document), nil
}

func (d DocumentsController) UpdateDocumentById(
	ctx context.Context,
	req *documents.CreateUpdateDocumentView,
	params documents.UpdateDocumentByIdParams,
) (documents.UpdateDocumentByIdRes, error) {
	document, err := d.service.UpdateDocument(ctx, d.mapper.UpdateDocumentReq(req, params))
	if d.service.IsNotFound(err) {
		return DocumentNotFound(params.ID)
	} else if err != nil {
		return FailedToUpdateDocument()
	}

	return d.mapper.Document(document), nil
}
