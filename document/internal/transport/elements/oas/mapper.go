package oas

import (
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	dto "github.com/0B1t322/Documents-Service/document/internal/dto/elements"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type mapper struct{}

func (mapper) GetElementsReq(bodyId uuid.UUID, cursor string, limit uint) dto.GetElementsDto {
	req := dto.GetElementsDto{
		Cursor: cursor,
		BodyID: bodyId,
		Limit:  limit,
	}

	return req
}

func (mapper) CreateElementReq(req *documents.CreateUpdateStructuralElement, bodyId uuid.UUID) dto.CreateElementDto {
	var (
		elementType    dto.ElementType
		elementStyleID int
	)
	{
		switch req.Element.Type {
		case documents.CreateUpdateParagraphCreateUpdateStructuralElementElement:
			elementType = dto.Paragraph
			elementStyleID = req.GetElement().CreateUpdateParagraph.ParagraphStyleId.Value
		case documents.CreateUpdateSectionBreakCreateUpdateStructuralElementElement:
			elementType = dto.SectionBreak
			elementStyleID = req.GetElement().CreateUpdateSectionBreak.SectionBreakStyleId.Value
		}
	}

	return dto.CreateElementDto{
		Index:          req.Index,
		BodyID:         bodyId,
		ElementType:    elementType,
		ElementStyleID: elementStyleID,
	}
}

func (mapper) UpdateElementReq(seID int, req *documents.UpdateStyleOfStructuralElement) dto.UpdateElementDto {
	var (
		elementType    dto.ElementType
		elementStyleID int
	)
	{
		switch req.Element.Type {
		case documents.CreateUpdateParagraphUpdateStyleOfStructuralElementElement:
			elementType = dto.Paragraph
			elementStyleID = req.GetElement().CreateUpdateParagraph.ParagraphStyleId.Value
		case documents.CreateUpdateSectionBreakUpdateStyleOfStructuralElementElement:
			elementType = dto.SectionBreak
			elementStyleID = req.GetElement().CreateUpdateSectionBreak.SectionBreakStyleId.Value
		}
	}
	return dto.UpdateElementDto{
		ID:             seID,
		ElementType:    elementType,
		ElementStyleID: elementStyleID,
	}
}

func (m mapper) CreateParagraphElementReq(
	seID int,
	req *documents.CreateUpdateParagraphElement,
) dto.CreateParagraphElementDto {
	r := dto.CreateParagraphElementDto{
		StructuralElementID: seID,
		Index:               req.Index,
	}

	switch req.Element.Type {
	case documents.CreateUpdateTextRunCreateUpdateParagraphElementElement:
		r.ElementType = dto.TextRun
		r.TextRun = dto.CreateTextRun{
			Content:     req.Element.CreateUpdateTextRun.Content,
			TextStyleID: req.Element.CreateUpdateTextRun.TextStyleId,
		}
	}

	return r
}

func (m mapper) UpdateParagraphElementReq(
	req *documents.UpdateParagraphElement,
) dto.UpdateParagraphElementDto {
	r := dto.UpdateParagraphElementDto{}

	switch req.Element.Type {
	case documents.CreateUpdateTextRunUpdateParagraphElementElement:
		r.ElementType = dto.TextRun
		r.TextRun = dto.UpdateTextRun{
			Content:     req.Element.CreateUpdateTextRun.Content,
			TextStyleID: req.Element.CreateUpdateTextRun.TextStyleId,
		}
	}

	return r
}

func (mapper) GetParagraphElementsReq(
	cursor string, bodyId uuid.UUID, seId int,
	limit uint,
) dto.GetParagraphsElementDto {
	return dto.GetParagraphsElementDto{
		Cursor:              cursor,
		BodyID:              bodyId,
		StructuralElementID: seId,
		Limit:               limit,
	}
}

func (m mapper) PaginatedElements(resp dto.GetElementsResponse) *documents.PaginatedStructuralElements {
	return &documents.PaginatedStructuralElements{
		Elements: lo.Map(
			resp.Elements, func(item models.StructuralElement, index int) documents.StructuralElement {
				return m.StructuralElement(item)
			},
		),
		Cursor: resp.NextCursor,
	}
}

func (m mapper) PaginatedParagraphsElements(resp dto.GetParagraphsElementResponse) *documents.PaginatedParagrahElements {
	return &documents.PaginatedParagrahElements{
		Items: lo.Map(
			resp.Elements,
			func(item models.ParagraphElement, index int) documents.ParagraphElement {
				return m.ParagraphElement(item)
			},
		),
		Cursor: resp.NextCursor,
	}
}

func (m mapper) StructuralElement(element models.StructuralElement) documents.StructuralElement {
	return documents.StructuralElement{
		ID:      element.ID,
		Index:   element.Index,
		Element: m.StructuralElementElement(element),
	}
}

func (m mapper) StructuralElementElement(element models.StructuralElement) documents.StructuralElementElement {
	see := documents.StructuralElementElement{}

	switch element.GetType() {
	case models.SEParagraph:
		see.SetParagraph(m.Paragraph(element.Paragraph))
	case models.SESectionBreak:
		see.SetSectionBreak(m.SectionBreak(element.SectionBreak))
	}

	return see
}

func (mapper) Paragraph(paragraph *models.Paragraph) documents.Paragraph {
	return documents.Paragraph{
		ID: paragraph.ID,
		ParagraphStyleId: documents.OptInt{
			Value: paragraph.ID,
			Set:   true,
		},
	}
}

func (mapper) SectionBreak(sectionBreak *models.SectionBreak) documents.SectionBreak {
	return documents.SectionBreak{
		ID: sectionBreak.ID,
		SectionBreakStyleId: documents.OptInt{
			Value: sectionBreak.ID,
			Set:   true,
		},
	}
}

func (m mapper) ParagraphElement(element models.ParagraphElement) documents.ParagraphElement {
	return documents.ParagraphElement{
		ID:      element.ID,
		Index:   element.Index,
		Element: m.ParagraphElementElement(element),
	}
}

func (m mapper) ParagraphElementElement(element models.ParagraphElement) documents.ParagraphElementElement {
	r := documents.ParagraphElementElement{}

	switch element.GetType() {
	case models.PETextRune:
		r.SetTextRun(m.TextRun(element.TextRune))
	}

	return r
}

func (mapper) TextRun(textRune *models.TextRune) documents.TextRun {
	return documents.TextRun{
		ID:          textRune.ID,
		Content:     textRune.Content,
		TextStyleId: textRune.TextStyleID,
	}
}
