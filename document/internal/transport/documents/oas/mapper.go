package oas

import (
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/dto/documents"
	"github.com/0B1t322/Online-Document-Redactor/pkg/gen/open-api/documents"
	"github.com/samber/lo"
)

type mapper struct{}

func (mapper) Document(document models.Document) *documents.Document {
	return &documents.Document{
		ID:    document.ID,
		Title: document.Title,
		Body: documents.Body{
			ID: document.Body.ID,
		},
		Style: documents.DocumentStyle{
			ID: document.Style.ID,
			Size: documents.Size{
				Height: documents.Dimension{
					Magnitude: document.Style.PageSize.Height.Magnitude,
					Unit:      documents.Unit(document.Style.PageSize.Height.Unit),
				},
				Width: documents.Dimension{
					Magnitude: document.Style.PageSize.Width.Magnitude,
					Unit:      documents.Unit(document.Style.PageSize.Width.Unit),
				},
			},
		},
	}
}

func (mapper) PaginatedDocuments(res dto.GetDocumentsResponse) *documents.PaginatedDocuments {
	return &documents.PaginatedDocuments{
		Items: lo.Map(
			res.Document, func(item models.Document, index int) documents.CompactDocument {
				return documents.CompactDocument{
					ID:    item.ID,
					Title: item.Title,
				}
			},
		),
		Cursor: res.Cursor,
	}
}

func (mapper) CreateDocumentReq(req *documents.CreateUpdateDocumentView) dto.CreateDocumentDto {
	return dto.CreateDocumentDto{
		Title: req.Title,
		Size: models.Size{
			Height: models.Dimension{
				Magnitude: req.Size.Height.Magnitude,
				Unit:      models.Unit(req.Size.Height.Unit),
			},
			Width: models.Dimension{
				Magnitude: req.Size.Width.Magnitude,
				Unit:      models.Unit(req.Size.Width.Unit),
			},
		},
	}
}

func (mapper) UpdateDocumentReq(
	req *documents.CreateUpdateDocumentView,
	params documents.UpdateDocumentByIdParams,
) dto.UpdateDocumentDto {
	return dto.UpdateDocumentDto{
		ID:    params.ID,
		Title: req.Title,
		Style: dto.UpdateDocumentStyleDto{
			PageSize: models.Size{
				Height: models.Dimension{
					Magnitude: req.Size.Height.Magnitude,
					Unit:      models.Unit(req.Size.Height.Unit),
				},
				Width: models.Dimension{
					Magnitude: req.Size.Width.Magnitude,
					Unit:      models.Unit(req.Size.Width.Unit),
				},
			},
		},
	}
}

func (mapper) GetDocumentsReq(params documents.DocumentsGetParams) dto.GetDocumentsDto {
	return dto.GetDocumentsDto{
		Cursor: params.Cursor.Value,
		Limit:  uint64(params.Limit.Value),
	}
}
