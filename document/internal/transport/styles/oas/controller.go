package oas

import (
	"context"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	dto "github.com/0B1t322/Documents-Service/document/internal/dto/styles"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type styleService interface {
	CreateStyleInDocument(ctx context.Context, documentId uuid.UUID, dto dto.StyleDto) (models.Style, error)
	UpdateStyleInDocument(
		ctx context.Context, documentId uuid.UUID, styleId uuid.UUID,
		styleDto dto.StyleDto,
	) (models.Style, error)
	GetStylesInDocument(ctx context.Context, documentId uuid.UUID) ([]models.Style, error)
	DeleteStyleInDocument(ctx context.Context, documentId uuid.UUID, styleId uuid.UUID) error

	IsNotFound(err error) bool
	IsValidation(err error) bool
}

type StylesController struct {
	service styleService
	mapper  mapper
}

func NewStylesController(service styleService) *StylesController {
	return &StylesController{
		service: service,
		mapper:  mapper{},
	}
}

func (s StylesController) CreateDocumentStyle(
	ctx context.Context,
	req *documents.CreateUpdateStyle,
	params documents.CreateDocumentStyleParams,
) (documents.CreateDocumentStyleRes, error) {
	style, err := s.service.CreateStyleInDocument(ctx, params.ID, s.mapper.CreateUpdateStyleInDocument(req))
	if s.service.IsNotFound(err) {
		return NotFound(err)
	} else if s.service.IsValidation(err) {
		return BadRequest(err)
	} else if err != nil {
		return FailedToCreateStyleInDocument()
	}

	return lo.ToPtr(s.mapper.Style(style)), nil
}

func (s StylesController) DeleteStyleById(
	ctx context.Context,
	params documents.DeleteStyleByIdParams,
) (documents.DeleteStyleByIdRes, error) {
	if err := s.service.DeleteStyleInDocument(ctx, params.ID, params.StyleId); s.service.IsNotFound(err) {
		return NotFound(err)
	} else if err != nil {
		return FailedToDeleteStyleInDocument()
	}

	return &documents.DeleteStyleByIdNoContent{}, nil
}

func (s StylesController) DocumentsIDStylesGet(
	ctx context.Context,
	params documents.DocumentsIDStylesGetParams,
) (documents.DocumentsIDStylesGetRes, error) {
	styles, err := s.service.GetStylesInDocument(ctx, params.ID)
	if s.service.IsNotFound(err) {
		return NotFound(err)
	} else if err != nil {
		return FailedToGetStylesInDocument()
	}

	return s.mapper.Styles(styles), nil
}

func (s StylesController) UpdateStyleById(
	ctx context.Context,
	req *documents.CreateUpdateStyle,
	params documents.UpdateStyleByIdParams,
) (documents.UpdateStyleByIdRes, error) {
	style, err := s.service.UpdateStyleInDocument(
		ctx, params.ID, params.StyleId,
		s.mapper.CreateUpdateStyleInDocument(req),
	)
	if s.service.IsNotFound(err) {
		return NotFound(err)
	} else if s.service.IsValidation(err) {
		return BadRequest(err)
	} else if err != nil {
		return FailedToUpdateInDocument()
	}

	return lo.ToPtr(s.mapper.Style(style)), nil
}
