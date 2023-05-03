package oas

import (
	"context"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	dto "github.com/0B1t322/Documents-Service/document/internal/dto/elements"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type (
	documentService interface {
		GetDocument(ctx context.Context, documentId uuid.UUID) (models.Document, error)
		IsNotFound(err error) bool
	}

	elementsService interface {
		CreateStructuralElement(
			ctx context.Context,
			req dto.CreateElementDto,
		) (models.StructuralElement, error)

		UpdateStructuralElementWithBodyID(
			ctx context.Context,
			bodyID uuid.UUID,
			req dto.UpdateElementDto,
		) (models.StructuralElement, error)

		GetStructuralElements(
			ctx context.Context,
			req dto.GetElementsDto,
		) (dto.GetElementsResponse, error)

		DeleteStructuralElementWithBodyID(ctx context.Context, id int, bodyId uuid.UUID) error

		CreateParagraphElementWithBodyID(
			ctx context.Context,
			bodyId uuid.UUID,
			req dto.CreateParagraphElementDto,
		) (models.ParagraphElement, error)

		DeleteParagraphElementWithBodyID(
			ctx context.Context,
			seId int,
			peId int,
			bodyId uuid.UUID,
		) error

		UpdateParagraphElementWithBodyID(
			ctx context.Context,
			seId int,
			bodyId uuid.UUID,
			id int,
			req dto.UpdateParagraphElementDto,
		) (models.ParagraphElement, error)

		GetParagraphElements(
			ctx context.Context,
			req dto.GetParagraphsElementDto,
		) (dto.GetParagraphsElementResponse, error)

		GetParagraphElementByIndexes(
			ctx context.Context,
			bodyId uuid.UUID,
			structuralElementIndex,
			paragraphElementIndex int,
		) (models.ParagraphElement, error)

		UpdateParagraphElementByIndexes(
			ctx context.Context,
			bodyId uuid.UUID,
			structuralElementIndex,
			paragraphElementIndex int,
			req dto.UpdateParagraphElementDto,
		) (models.ParagraphElement, error)

		IsParametersNotValid(err error) bool
		IsNotFound(err error) bool
		IsValidation(err error) bool
	}
)

type ElementsController struct {
	service         elementsService
	documentService documentService
	mapper          mapper
}

func NewElementsController(service elementsService, documentsService documentService) *ElementsController {
	return &ElementsController{
		service:         service,
		documentService: documentsService,
		mapper:          mapper{},
	}
}

func (e ElementsController) CreateElement(
	ctx context.Context,
	req *documents.CreateUpdateStructuralElement,
	params documents.CreateElementParams,
) (documents.CreateElementRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	switch {
	case e.service.IsNotFound(err):
		return DocumentNotFound(params.ID)
	case err != nil:
		return FailedToCreateElement()
	}

	element, err := e.service.CreateStructuralElement(
		ctx,
		e.mapper.CreateElementReq(req, document.Body.ID),
	)
	switch {
	case e.service.IsValidation(err):
		return BadRequest(err)
	case err != nil:
		return FailedToCreateElement()
	}

	return lo.ToPtr(e.mapper.StructuralElement(element)), nil
}

func (e ElementsController) CreateParagraphElement(
	ctx context.Context,
	req *documents.CreateUpdateParagraphElement,
	params documents.CreateParagraphElementParams,
) (documents.CreateParagraphElementRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	switch {
	case e.service.IsNotFound(err):
		return DocumentNotFound(params.ID)
	case err != nil:
		return FailedToCreateParagraphElement()
	}

	element, err := e.service.CreateParagraphElementWithBodyID(
		ctx,
		document.Body.ID,
		e.mapper.CreateParagraphElementReq(params.SeId, req),
	)
	switch {
	case e.service.IsNotFound(err):
		return StructuralElementNotFoundInDocument(params.SeId, params.ID)
	case e.service.IsValidation(err):
		return BadRequest(err)
	case err != nil:
		return FailedToCreateParagraphElement()
	}

	return lo.ToPtr(e.mapper.ParagraphElement(element)), nil
}

func (e ElementsController) DeleteParagraphElement(
	ctx context.Context,
	params documents.DeleteParagraphElementParams,
) (documents.DeleteParagraphElementRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	switch {
	case e.documentService.IsNotFound(err):
		return DocumentNotFound(params.ID)
	case err != nil:
		return FailedToDeleteParagraphElement()
	}

	err = e.service.DeleteParagraphElementWithBodyID(ctx, params.SeId, params.ElementId, document.Body.ID)
	switch {
	case e.service.IsNotFound(err):
		return ParagraphElementNotFoundInDocument(params.ElementId, params.SeId, params.ID)
	case err != nil:
		return FailedToDeleteParagraphElement()
	}

	return &documents.DeleteParagraphElementNoContent{}, nil
}

func (e ElementsController) DeleteStructuralElementByID(
	ctx context.Context,
	params documents.DeleteStructuralElementByIDParams,
) (documents.DeleteStructuralElementByIDRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	if e.documentService.IsNotFound(err) {
		return DocumentNotFound(params.ID)
	} else if err != nil {
		return FailedDeleteStructuralElementByID()
	}

	if err := e.service.DeleteStructuralElementWithBodyID(
		ctx, params.SeId,
		document.Body.ID,
	); e.service.IsNotFound(err) {
		return StructuralElementNotFoundInDocument(params.SeId, params.ID)
	} else if err != nil {
		return FailedDeleteStructuralElementByID()
	}

	return &documents.DeleteStructuralElementByIDNoContent{}, nil
}

func (e ElementsController) DocumentsIDElementsSeIdGet(
	ctx context.Context,
	params documents.DocumentsIDElementsSeIdGetParams,
) (documents.DocumentsIDElementsSeIdGetRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	if e.documentService.IsNotFound(err) {
		return DocumentNotFound(params.ID)
	} else if err != nil {
		return FailedToGetParagraphElements()
	}

	resp, err := e.service.GetParagraphElements(
		ctx, e.mapper.GetParagraphElementsReq(
			params.Cursor.Value,
			document.Body.ID, params.SeId, params.Limit.Or(0),
		),
	)
	if err != nil {
		return FailedToGetParagraphElements()
	}

	return e.mapper.PaginatedParagraphsElements(resp), nil
}

func (e ElementsController) GetElements(
	ctx context.Context,
	params documents.GetElementsParams,
) (documents.GetElementsRes, error) {
	var (
		bodyId uuid.UUID
		cursor string
		limit  uint
	)
	{
		document, err := e.documentService.GetDocument(ctx, params.ID)
		if e.documentService.IsNotFound(err) {
			return DocumentNotFound(params.ID)
		} else if err != nil {
			return FailedToGetElements()
		}

		bodyId = document.Body.ID
		cursor = params.Cursor.Value
		limit = params.Limit.Or(10)
	}

	resp, err := e.service.GetStructuralElements(ctx, e.mapper.GetElementsReq(bodyId, cursor, limit))
	if e.service.IsParametersNotValid(err) {
		return BadRequest(err)
	} else if err != nil {
		return FailedToGetElements()
	}

	return e.mapper.PaginatedElements(resp), nil
}

func (e ElementsController) UpdateParagraphElement(
	ctx context.Context,
	req *documents.UpdateParagraphElement,
	params documents.UpdateParagraphElementParams,
) (documents.UpdateParagraphElementRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	if e.documentService.IsNotFound(err) {
		return DocumentNotFound(params.ID)
	} else if err != nil {
		return FailedToUpdateParagraphElement()
	}

	element, err := e.service.UpdateParagraphElementWithBodyID(
		ctx, params.SeId, document.Body.ID, params.ElementId,
		e.mapper.UpdateParagraphElementReq(req),
	)
	switch {
	case e.service.IsNotFound(err):
		return NotFound(err)
	case err != nil:
		return FailedToUpdateParagraphElement()
	}

	return lo.ToPtr(e.mapper.ParagraphElement(element)), nil
}

func (e ElementsController) UpdateStructuralElement(
	ctx context.Context,
	req *documents.UpdateStyleOfStructuralElement,
	params documents.UpdateStructuralElementParams,
) (documents.UpdateStructuralElementRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	switch {
	case e.service.IsNotFound(err):
		return DocumentNotFound(params.ID)
	case err != nil:
		return FailedToUpdateElement()
	}

	element, err := e.service.UpdateStructuralElementWithBodyID(
		ctx, document.Body.ID,
		e.mapper.UpdateElementReq(params.SeId, req),
	)
	switch {
	case e.service.IsNotFound(err):
		return StructuralElementNotFoundInDocument(params.SeId, params.ID)
	case err != nil:
		return FailedToUpdateElement()
	}

	return lo.ToPtr(e.mapper.StructuralElement(element)), nil
}

func (e ElementsController) GetParagraphElementByIndexes(
	ctx context.Context,
	params documents.GetParagraphElementByIndexesParams,
) (documents.GetParagraphElementByIndexesRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	switch {
	case e.service.IsNotFound(err):
		return DocumentNotFound(params.ID)
	case err != nil:
		return FailedToGetParagraphElementByIndexes()
	}

	element, err := e.service.GetParagraphElementByIndexes(
		ctx,
		document.Body.ID,
		params.StructuralElementIndex,
		params.ParagraphElementIndex,
	)
	switch {
	case e.service.IsNotFound(err):
		return NotFound(err)
	case err != nil:
		return FailedToGetParagraphElementByIndexes()
	}

	return lo.ToPtr(e.mapper.ParagraphElement(element)), nil
}

func (e ElementsController) UpdateParagraphElementByIndexes(
	ctx context.Context, req *documents.UpdateParagraphElement,
	params documents.UpdateParagraphElementByIndexesParams,
) (documents.UpdateParagraphElementByIndexesRes, error) {
	document, err := e.documentService.GetDocument(ctx, params.ID)
	switch {
	case e.service.IsNotFound(err):
		return DocumentNotFound(params.ID)
	case err != nil:
		return FailedToUpdateParagraphElementByIndexes()
	}

	element, err := e.service.UpdateParagraphElementByIndexes(
		ctx, document.Body.ID, params.StructuralElementIndex,
		params.ParagraphElementIndex, e.mapper.UpdateParagraphElementReq(req),
	)
	switch {
	case e.service.IsNotFound(err):
		return NotFound(err)
	case e.service.IsValidation(err):
		return BadRequest(err)
	case err != nil:
		return FailedToUpdateParagraphElementByIndexes()
	}

	return lo.ToPtr(e.mapper.ParagraphElement(element)), nil
}
