package elements

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations/text"
)

type DeleteParagraphElement struct {
	ParagraphIndex int
	text.DeleteText
}

func NewDeleteParagraphElement(paragraphElementIndex int, localIndex int, deleteText string) DeleteParagraphElement {
	return DeleteParagraphElement{
		ParagraphIndex: paragraphElementIndex,
		DeleteText:     text.NewDeleteText(localIndex, deleteText),
	}
}

func (d DeleteParagraphElement) GetIndex() int {
	return d.ParagraphIndex
}

func (d DeleteParagraphElement) SetIndex(i int) operations.Delete {
	d.ParagraphIndex = i
	return d
}

func (d DeleteParagraphElement) DeleteLength() int {
	return d.DeleteText.DeleteLength()
}

func (d DeleteParagraphElement) IncludeInsert(insert operations.Insert) operations.Delete {
	switch op := insert.(type) {
	case text.InsertText:
		d.DeleteText = d.DeleteText.IncludeInsertText(op)
		return d
	case InsertParagraphElement:
		return d.IncludeInsertParagraphElement(op)
	default:
		panic("Unsupported insert operation")
	}
}

func (d DeleteParagraphElement) IncludeInsertParagraphElement(op InsertParagraphElement) DeleteParagraphElement {
	if d.GetIndex() != op.GetIndex() {
		return d
	}

	d.DeleteText = d.DeleteText.IncludeInsertText(op.InsertText)

	return d
}

func (d DeleteParagraphElement) IncludeDelete(delete operations.Delete) operations.Delete {
	switch op := delete.(type) {
	case text.DeleteText:
		d.DeleteText = d.DeleteText.IncludeDeleteText(op)
		return d
	case DeleteParagraphElement:
		return d.IncludeDeleteParagraphElement(op)
	default:
		panic("Unsupported delete operation")
	}
}

func (d DeleteParagraphElement) IncludeDeleteParagraphElement(op DeleteParagraphElement) DeleteParagraphElement {
	if d.GetIndex() != op.GetIndex() {
		return d
	}

	d.DeleteText = d.DeleteText.IncludeDeleteText(op.DeleteText)
	return d
}

func (d DeleteParagraphElement) Apply(to ParagraphElement) {
	to.SetContent(d.DeleteText.Apply(to.GetContent()))
}
