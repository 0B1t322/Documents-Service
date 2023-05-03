package operations

import (
	"context"
	"fmt"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	"github.com/google/uuid"
)

type paragraphOperation interface {
	ApplyToParagraph(to *models.ParagraphElement)
	StructuralElementIndex() int
	ParagraphElementIndex() int
}

type operationsAndElement struct {
	Element    models.ParagraphElement
	Operations []paragraphOperation
}

func createOperationsAndElement(operations []models.RevisionOperation) (slice operationsAndElements) {
	m := map[string]operationsAndElement{}
	{
		for _, op := range operations {
			key := fmt.Sprintf("%v%v", op.StructuralElementIndex(), op.ParagraphElementIndex())
			if _, find := m[key]; !find {
				m[key] = operationsAndElement{
					Element:    models.ParagraphElement{},
					Operations: []paragraphOperation{op},
				}
				continue
			}
			value := m[key]
			value.Operations = append(value.Operations, op)
			m[key] = value
		}
	}

	for _, v := range m {
		slice = append(slice, v)
	}

	return
}

type elementsGetter interface {
	GetParagraphElementByIndexes(
		ctx context.Context, documentId uuid.UUID, seId,
		peId int,
	) (models.ParagraphElement, error)
}

func (o *operationsAndElement) getElement(ctx context.Context, documentId uuid.UUID, getter elementsGetter) error {
	var (
		seId int
		peId int
	)
	seId = o.Operations[0].StructuralElementIndex()
	peId = o.Operations[0].ParagraphElementIndex()

	element, err := getter.GetParagraphElementByIndexes(ctx, documentId, seId, peId)
	if err != nil {
		return err
	}

	o.Element = element

	return nil
}

func (o *operationsAndElement) applyOperations() {
	for _, op := range o.Operations {
		op.ApplyToParagraph(&o.Element)
	}
}

type elementsUpdater interface {
	UpdateParagraphElementByIndexes(
		ctx context.Context, documentId uuid.UUID,
		seId, peId int,
		element models.ParagraphElement,
	) error
}

func (o *operationsAndElement) updateElement(ctx context.Context, documentId uuid.UUID, updater elementsUpdater) error {
	var (
		seId int
		peId int
	)
	seId = o.Operations[0].StructuralElementIndex()
	peId = o.Operations[0].ParagraphElementIndex()

	if err := updater.UpdateParagraphElementByIndexes(ctx, documentId, seId, peId, o.Element); err != nil {
		return err
	}

	return nil
}

type operationsAndElements []operationsAndElement

func (o operationsAndElements) getElements(ctx context.Context, documentId uuid.UUID, getter elementsGetter) error {
	for i := range o {
		if err := o[i].getElement(ctx, documentId, getter); err != nil {
			return err
		}
	}

	return nil
}

func (o operationsAndElements) applyOperations() {
	for i := range o {
		o[i].applyOperations()
	}
}

func (o operationsAndElements) updateElements(
	ctx context.Context, documentId uuid.UUID,
	updater elementsUpdater,
) error {
	for i := range o {
		if err := o[i].updateElement(ctx, documentId, updater); err != nil {
			return err
		}
	}

	return nil
}
