package services

import (
	"context"
	"fmt"
	"github.com/0B1t322/Documents-Service/document/internal/core/events"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	dto "github.com/0B1t322/Documents-Service/document/internal/dto/elements"
	repository "github.com/0B1t322/Documents-Service/document/internal/repository/elements"
	"github.com/0B1t322/Documents-Service/internal/core/utils/cursor"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
)

type (
	elementsRepository interface {
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

		FindParagraphElementByIndexes(
			ctx context.Context,
			bodyId uuid.UUID,
			seId,
			peId int,
		) (models.ParagraphElement, error)
	}
)

type ElementsService struct {
	repository elementsRepository
	logger     log.Logger
	publisher  events.EventPublisher
}

func NewElementsService(
	elementsRepository elementsRepository,
	publisher events.EventPublisher,
	logger log.Logger,
) *ElementsService {
	return &ElementsService{
		repository: elementsRepository,
		logger:     log.WithPrefix(logger, "service", "ElementsService"),
		publisher:  publisher,
	}
}

func (s ElementsService) createStructuralElement(from dto.CreateElementDto) (models.StructuralElement, error) {
	element := models.StructuralElement{
		Index: from.Index,
	}

	switch from.ElementType {
	case dto.Paragraph:
		element.Paragraph = &models.Paragraph{
			ParagraphStyleId: from.ElementStyleID,
		}
	case dto.SectionBreak:
		element.SectionBreak = &models.SectionBreak{
			SectionBreakStyleId: from.ElementStyleID,
		}
	default:
		return models.StructuralElement{}, fmt.Errorf("Unknown element type")
	}

	return element, nil
}

func (s ElementsService) CreateStructuralElement(
	ctx context.Context,
	req dto.CreateElementDto,
) (models.StructuralElement, error) {
	element, err := s.createStructuralElement(req)
	if err != nil {
		return element, err
	}

	if err := s.repository.StoreStructuralElement(ctx, req.BodyID, &element); err == repository.ErrBadIndex {
		return element, ErrStructuralElementBadIndex
	} else if err != nil {
		level.Error(s.logger).Log("Failed to create structural element", err)
		return element, err
	}

	s.publisher.PublishEvent(ctx, events.StructuralElementCreatedEvent{StructuralElement: element})

	return element, nil
}

func (ElementsService) elementTypeToDtoType(elementType models.SEType) dto.ElementType {
	switch elementType {
	case models.SEParagraph:
		return dto.Paragraph
	case models.SESectionBreak:
		return dto.SectionBreak
	}

	panic("Unknown type")
}

func (s ElementsService) typesIsEqual(reqType dto.ElementType, elementType models.SEType) error {
	dtoType := s.elementTypeToDtoType(elementType)

	if dtoType != reqType {
		return ErrSEChildElementBadType
	}

	return nil
}

func (s ElementsService) updateStructuralElement(
	ctx context.Context,
	element models.StructuralElement,
	req dto.UpdateElementDto,
) (models.StructuralElement, error) {
	// check type of child element
	if err := s.typesIsEqual(req.ElementType, element.GetType()); err != nil {
		return models.StructuralElement{}, err
	}

	event := events.StructuralElementUpdatedEvent{
		Was: element,
	}

	switch req.ElementType {
	case dto.Paragraph:
		element.Paragraph.ParagraphStyleId = req.ElementStyleID
	case dto.SectionBreak:
		element.SectionBreak.SectionBreakStyleId = req.ElementStyleID
	}
	event.Become = element

	if err := s.repository.UpdateStructuralElement(ctx, element); err != nil {
		return models.StructuralElement{}, err
	}

	s.publisher.PublishEvent(ctx, event)

	return element, nil
}

func (s ElementsService) UpdateStructuralElement(
	ctx context.Context,
	req dto.UpdateElementDto,
) (models.StructuralElement, error) {
	element, err := s.GetStructuralElement(ctx, req.ID)
	if err != nil {
		return models.StructuralElement{}, err
	}

	return s.updateStructuralElement(ctx, element, req)
}

func (s ElementsService) UpdateStructuralElementWithBodyID(
	ctx context.Context,
	bodyID uuid.UUID,
	req dto.UpdateElementDto,
) (models.StructuralElement, error) {
	element, err := s.GetStructuralElementByBodyAndID(ctx, bodyID, req.ID)
	if err != nil {
		return models.StructuralElement{}, err
	}

	return s.updateStructuralElement(ctx, element, req)
}

func (s ElementsService) GetStructuralElements(
	ctx context.Context,
	req dto.GetElementsDto,
) (dto.GetElementsResponse, error) {
	var cur int
	{
		if req.Cursor == "" {
			cur = -1
		} else if decoded, err := cursor.CursorToInt(req.Cursor); err != nil {
			return dto.GetElementsResponse{}, err
		} else {
			cur = decoded
		}

	}

	elements, nextCursor, err := s.repository.GetStructuralElements(ctx, req.BodyID, cur, req.Limit)
	if err != nil {
		level.Error(s.logger).Log("Failed to get structural elements", err)
		return dto.GetElementsResponse{}, err
	}

	return dto.GetElementsResponse{
		Elements:   elements,
		NextCursor: cursor.IntToCursor(nextCursor),
	}, nil
}

func (s ElementsService) DeleteStructuralElement(ctx context.Context, id int) error {
	element, err := s.GetStructuralElement(ctx, id)
	if err != nil {
		return err
	}

	if err := s.deleteStructuralElement(ctx, element); err != nil {
		return err
	}

	s.publisher.PublishEvent(
		ctx, events.StructuralElementDeletedEvent{
			StructuralElement: element,
		},
	)

	return nil
}

func (s ElementsService) deleteStructuralElement(ctx context.Context, element models.StructuralElement) error {
	if err := s.repository.DeleteStructuralElement(ctx, element); err != nil {
		level.Error(s.logger).Log("Failed to delete structural element", err)
		return err
	}

	return nil
}

func (s ElementsService) DeleteStructuralElementWithBodyID(ctx context.Context, id int, bodyId uuid.UUID) error {
	element, err := s.GetStructuralElementByBodyAndID(ctx, bodyId, id)
	if err != nil {
		return err
	}

	if err := s.deleteStructuralElement(ctx, element); err != nil {
		return err
	}

	s.publisher.PublishEvent(
		ctx, events.StructuralElementDeletedEvent{
			StructuralElement: element,
		},
	)

	return nil
}

func (s ElementsService) GetStructuralElement(
	ctx context.Context,
	id int,
) (models.StructuralElement, error) {
	element, err := s.repository.FindStructuralElementByID(ctx, id)
	if err == repository.ErrStructuralElementNotFound {
		return models.StructuralElement{}, ErrStructuralElementNotFound
	} else if err != nil {
		level.Error(s.logger).Log("Failed to get structural element", err)
		return models.StructuralElement{}, err
	}

	return element, nil
}

func (s ElementsService) GetStructuralElementByBodyAndID(
	ctx context.Context,
	bodyId uuid.UUID,
	id int,
) (models.StructuralElement, error) {
	element, err := s.repository.FindStructuralElementByIDAndBodyID(ctx, id, bodyId)
	if err == repository.ErrStructuralElementNotFound {
		return models.StructuralElement{}, ErrStructuralElementNotFound
	} else if err != nil {
		return element, err
	}

	return element, nil
}

func createParagraphElement(from dto.CreateParagraphElementDto) (models.ParagraphElement, error) {
	element := models.ParagraphElement{
		Index: from.Index,
	}

	switch from.ElementType {
	case dto.TextRun:
		element.TextRune = &models.TextRune{
			Content:     from.TextRun.Content,
			TextStyleID: from.TextRun.TextStyleID,
		}
	default:
		return models.ParagraphElement{}, fmt.Errorf("Unknown element type")
	}

	return element, nil
}

func (s ElementsService) CreateParagraphElementWithBodyID(
	ctx context.Context,
	bodyId uuid.UUID,
	req dto.CreateParagraphElementDto,
) (models.ParagraphElement, error) {
	structuralElement, err := s.repository.FindStructuralElementByIDAndBodyID(ctx, req.StructuralElementID, bodyId)
	if err == repository.ErrStructuralElementNotFound {
		return models.ParagraphElement{}, ErrStructuralElementNotFound
	} else if err != nil {
		return models.ParagraphElement{}, err
	}

	if structuralElement.GetType() != models.SEParagraph {
		return models.ParagraphElement{}, ErrSEChildIsNotParagraph
	}

	element, err := createParagraphElement(req)
	if err != nil {
		return models.ParagraphElement{}, err
	}

	if err := s.repository.StoreParagraphElement(
		ctx,
		structuralElement.Paragraph.ID,
		&element,
	); err == repository.ErrBadIndex {
		return models.ParagraphElement{}, ErrParagraphElementBadIndex
	} else if err != nil {
		level.Error(s.logger).Log("err", err)
		return models.ParagraphElement{}, err
	}

	return element, nil
}

func (s ElementsService) GetParagraphElements(
	ctx context.Context,
	req dto.GetParagraphsElementDto,
) (dto.GetParagraphsElementResponse, error) {
	var cur int
	{
		if req.Cursor == "" {
			cur = -1
		} else if decoded, err := cursor.CursorToInt(req.Cursor); err != nil {
			cur = decoded
		} else {
			return dto.GetParagraphsElementResponse{}, err
		}
	}

	structuralElement, err := s.repository.FindStructuralElementByIDAndBodyID(ctx, req.StructuralElementID, req.BodyID)
	if err == repository.ErrStructuralElementNotFound {
		return dto.GetParagraphsElementResponse{}, ErrStructuralElementNotFound
	} else if err != nil {
		level.Error(s.logger).Log("Can't find structural element", err)
		return dto.GetParagraphsElementResponse{}, err
	}

	if structuralElement.GetType() != models.SEParagraph {
		return dto.GetParagraphsElementResponse{}, ErrSEChildIsNotParagraph
	}

	elements, nextCursor, err := s.repository.GetParagraphElements(ctx, structuralElement.Paragraph.ID, cur, req.Limit)
	if err != nil {
		level.Error(s.logger).Log("Can't get paragraph elements", err)
		return dto.GetParagraphsElementResponse{}, err
	}

	return dto.GetParagraphsElementResponse{
		Elements:   elements,
		NextCursor: cursor.IntToCursor(nextCursor),
	}, nil
}

func (s ElementsService) DeleteParagraphElementWithBodyID(
	ctx context.Context,
	seId int,
	peId int,
	bodyId uuid.UUID,
) error {
	structuralElement, err := s.repository.FindStructuralElementByIDAndBodyID(ctx, seId, bodyId)
	if err == repository.ErrStructuralElementNotFound {
		return ErrStructuralElementNotFound
	} else if err != nil {
		return err
	}

	if structuralElement.GetType() != models.SEParagraph {
		return ErrSEChildIsNotParagraph
	}

	element, err := s.repository.GetParagraphElement(ctx, structuralElement.Paragraph.ID, peId)
	if err == repository.ErrParagraphElementNotFound {
		return ErrParagraphElementNotFound
	} else if err != nil {
		return err
	}

	if err := s.repository.DeleteParagraphElement(ctx, element); err != nil {
		level.Error(s.logger).Log("err", err)
		return err
	}

	return nil
}

func paragraphElementTypeToDtoType(peType models.PEType) dto.ParagraphElementType {
	switch peType {
	case models.PETextRune:
		return dto.TextRun
	}
	panic("Unknown type")
}

func paragraphElementTypesIsEqual(reqType dto.ParagraphElementType, elementType models.PEType) error {
	dtoType := paragraphElementTypeToDtoType(elementType)

	if dtoType != reqType {
		return ErrPEChildElementBadType
	}

	return nil
}

func (s ElementsService) UpdateParagraphElementWithBodyID(
	ctx context.Context,
	seId int,
	bodyId uuid.UUID,
	id int,
	req dto.UpdateParagraphElementDto,
) (models.ParagraphElement, error) {
	structuralElement, err := s.repository.FindStructuralElementByIDAndBodyID(ctx, seId, bodyId)
	if err == repository.ErrStructuralElementNotFound {
		return models.ParagraphElement{}, ErrStructuralElementNotFound
	} else if err != nil {
		return models.ParagraphElement{}, err
	}

	if structuralElement.GetType() != models.SEParagraph {
		return models.ParagraphElement{}, ErrSEChildIsNotParagraph
	}

	element, err := s.repository.GetParagraphElement(ctx, structuralElement.Paragraph.ID, id)
	if err == repository.ErrParagraphElementNotFound {
		return models.ParagraphElement{}, ErrParagraphElementNotFound
	} else if err != nil {
		return models.ParagraphElement{}, err
	}

	element, err = s.updateParagraphElement(ctx, req, element)
	if err != nil {
		return models.ParagraphElement{}, err
	}

	return element, nil
}

func (s ElementsService) updateParagraphElement(
	ctx context.Context,
	req dto.UpdateParagraphElementDto,
	element models.ParagraphElement,
) (models.ParagraphElement, error) {
	if err := paragraphElementTypesIsEqual(req.ElementType, element.GetType()); err != nil {
		return models.ParagraphElement{}, err
	}

	switch element.GetType() {
	case models.PETextRune:
		element.TextRune.TextStyleID = req.TextRun.TextStyleID
		element.TextRune.Content = req.TextRun.Content
	}

	if err := s.repository.UpdateParagraphElement(ctx, element); err == repository.ErrBadIndex {
		return models.ParagraphElement{}, ErrParagraphElementBadIndex
	} else if err != nil {
		level.Error(s.logger).Log("err", err)
		return models.ParagraphElement{}, err
	}

	return element, nil
}

func (s ElementsService) GetParagraphElementByIndexes(
	ctx context.Context,
	bodyId uuid.UUID,
	structuralElementIndex,
	paragraphElementIndex int,
) (models.ParagraphElement, error) {
	element, err := s.repository.FindParagraphElementByIndexes(
		ctx, bodyId, structuralElementIndex,
		paragraphElementIndex,
	)
	if err == repository.ErrParagraphElementNotFound {
		return element, ErrParagraphElementNotFound
	} else if err != nil {
		level.Error(s.logger).Log(
			"description", "Failed to get paragraph element by indexes",
			"err", err,
		)
		return element, nil
	}

	return element, nil
}

func (s ElementsService) UpdateParagraphElementByIndexes(
	ctx context.Context,
	bodyId uuid.UUID,
	structuralElementIndex,
	paragraphElementIndex int,
	req dto.UpdateParagraphElementDto,
) (models.ParagraphElement, error) {
	element, err := s.GetParagraphElementByIndexes(ctx, bodyId, structuralElementIndex, paragraphElementIndex)
	if err != nil {
		return models.ParagraphElement{}, err
	}

	element, err = s.updateParagraphElement(ctx, req, element)
	if err != nil {
		return models.ParagraphElement{}, err
	}

	return element, nil
}
