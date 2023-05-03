package elements

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations"
)

type InsertStructuralElement struct {
	StructuralElementIndex int
	InsertParagraphElement
}

func NewInsertStructuralElement(
	structuralElementIndex, paragraphElementIndex, localIndex int,
	insertText string,
) InsertStructuralElement {
	return InsertStructuralElement{
		StructuralElementIndex: structuralElementIndex,
		InsertParagraphElement: NewInsertParagraphElement(paragraphElementIndex, localIndex, insertText),
	}
}

func (i InsertStructuralElement) GetIndex() int {
	return i.StructuralElementIndex
}

func (i InsertStructuralElement) SetIndex(index int) operations.Insert {
	i.Index = index
	return i
}

func (i InsertStructuralElement) InsertLength() int {
	return i.InsertParagraphElement.InsertLength()
}

func (i InsertStructuralElement) IncludeInsert(insert operations.Insert) operations.Insert {
	switch op := insert.(type) {
	case InsertStructuralElement:
		if i.StructuralElementIndex != op.StructuralElementIndex {
			return i
		}
		i.InsertParagraphElement = i.InsertParagraphElement.IncludeInsertParagraphElement(op.InsertParagraphElement)
		return i
	default:
		panic("Unsupported insert operation")
	}
}

func (i InsertStructuralElement) IncludeDelete(delete operations.Delete) operations.Insert {
	switch op := delete.(type) {
	case DeleteStructuralElement:
		if i.StructuralElementIndex != op.StructuralElementIndex {
			return i
		}

		i.InsertParagraphElement = i.InsertParagraphElement.IncludeDeleteParagraphElement(op.DeleteParagraphElement)
		return i
	default:
		panic("Unsupported delete operation")
	}
}

func (i InsertStructuralElement) Apply(to StructuralElement) {
	i.InsertParagraphElement.Apply(to.GetParagraphElement(i.StructuralElementIndex))
}
