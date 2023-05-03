package elements

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations"
)

type DeleteStructuralElement struct {
	StructuralElementIndex int
	DeleteParagraphElement
}

func NewDeleteStructuralElement(
	structuralElementIndex, paragraphElementIndex, localIndex int,
	deleteText string,
) DeleteStructuralElement {
	return DeleteStructuralElement{
		StructuralElementIndex: structuralElementIndex,
		DeleteParagraphElement: NewDeleteParagraphElement(paragraphElementIndex, localIndex, deleteText),
	}
}

func (d DeleteStructuralElement) GetIndex() int {
	return d.StructuralElementIndex
}

func (d DeleteStructuralElement) SetIndex(i int) operations.Delete {
	d.StructuralElementIndex = i
	return d
}

func (d DeleteStructuralElement) DeleteLength() int {
	return d.DeleteLength()
}

func (d DeleteStructuralElement) IncludeInsert(insert operations.Insert) operations.Delete {
	switch op := insert.(type) {
	case InsertStructuralElement:
		if d.StructuralElementIndex != op.StructuralElementIndex {
			return d
		}
		d.DeleteParagraphElement = d.DeleteParagraphElement.IncludeInsertParagraphElement(op.InsertParagraphElement)

		return d
	default:
		panic("Unsupported insert operation")
	}
}

func (d DeleteStructuralElement) IncludeDelete(delete operations.Delete) operations.Delete {
	switch op := delete.(type) {
	case DeleteStructuralElement:
		if d.StructuralElementIndex != op.StructuralElementIndex {
			return d
		}
		d.DeleteParagraphElement = d.DeleteParagraphElement.IncludeDeleteParagraphElement(op.DeleteParagraphElement)

		return d
	default:
		panic("Unsupported delete operation")
	}
}

func (i DeleteStructuralElement) Apply(to StructuralElement) {
	i.DeleteParagraphElement.Apply(to.GetParagraphElement(i.StructuralElementIndex))
}
