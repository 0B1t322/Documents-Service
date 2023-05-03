package elements

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations/text"
)

type InsertParagraphElement struct {
	ParagraphIndex int
	text.InsertText
}

func NewInsertParagraphElement(paragraphElementIndex int, localIndex int, insertText string) InsertParagraphElement {
	return InsertParagraphElement{
		ParagraphIndex: paragraphElementIndex,
		InsertText:     text.NewInsertText(localIndex, insertText),
	}
}

func (i InsertParagraphElement) GetIndex() int {
	return i.ParagraphIndex
}

func (i InsertParagraphElement) SetIndex(index int) operations.Insert {
	i.ParagraphIndex = index
	return i
}

func (i InsertParagraphElement) InsertLength() int {
	return i.InsertText.InsertLength()
}

func (i InsertParagraphElement) IncludeInsert(insert operations.Insert) operations.Insert {
	switch op := insert.(type) {
	case text.InsertText:
		i.InsertText = i.InsertText.IncludeInsertText(op)
		return i
	case InsertParagraphElement:
		return i.IncludeInsertParagraphElement(op)
	default:
		panic("Unsupported insert operation")
	}
}

func (i InsertParagraphElement) IncludeInsertParagraphElement(op InsertParagraphElement) InsertParagraphElement {
	if i.GetIndex() != op.GetIndex() {
		return i
	}
	i.InsertText = i.InsertText.IncludeInsertText(op.InsertText)
	return i
}

func (i InsertParagraphElement) IncludeDelete(delete operations.Delete) operations.Insert {
	switch op := delete.(type) {
	case text.DeleteText:
		i.InsertText = i.InsertText.IncludeDeleteText(op)
		return i
	case DeleteParagraphElement:
		return i.IncludeDeleteParagraphElement(op)
	default:
		panic("Unsupported delete operation")
	}
}

func (i InsertParagraphElement) Apply(to ParagraphElement) {
	to.SetContent(i.InsertText.Apply(to.GetContent()))
}

func (i InsertParagraphElement) IncludeDeleteParagraphElement(op DeleteParagraphElement) InsertParagraphElement {
	if i.GetIndex() != op.GetIndex() {
		return i
	}

	i.InsertText = i.InsertText.IncludeDeleteText(op.DeleteText)
	return i
}
