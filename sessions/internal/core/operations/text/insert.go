package text

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations"
	"strings"
)

type InsertText struct {
	Index int
	Text  string
}

func NewInsertText(index int, text string) InsertText {
	return InsertText{
		Index: index,
		Text:  text,
	}
}

func (op InsertText) Apply(to string) string {
	b := strings.Builder{}

	b.WriteString(to[:op.Index])
	b.WriteString(op.Text)
	b.WriteString(to[op.Index:])

	return b.String()
}

func (op InsertText) GetIndex() int {
	return op.Index
}

func (op InsertText) SetIndex(i int) operations.Insert {
	op.Index = i
	return op
}

func (op1 InsertText) IncludeInsert(insert operations.Insert) operations.Insert {
	op2, ok := insert.(InsertText)
	if !ok {
		panic("Unsupported insert operation")
	}

	return op1.IncludeInsertText(op2)
}

func (op1 InsertText) IncludeInsertText(op2 InsertText) InsertText {
	if op1.GetIndex() < op2.GetIndex() {
		return op1
	}

	// Вставка была до другой вставки
	newIndex := op1.GetIndex() + op2.InsertLength()
	op1.Index = newIndex
	return op1
}

func (op1 InsertText) IncludeDelete(delete operations.Delete) operations.Insert {
	op2, ok := delete.(DeleteText)
	if !ok {
		panic("Unsupported insert operation")
	}

	return op1.IncludeDeleteText(op2)
}

func (op1 InsertText) IncludeDeleteText(op2 DeleteText) InsertText {
	// Удаление была после индексу до вставки в текст
	if op1.GetIndex() < op2.GetIndex() {
		return op1
	}

	// Удаление произошло до вставки по индексу
	newIndex := op1.GetIndex() - op2.DeleteLength()
	if newIndex < 0 {
		newIndex = 0
	}
	op1.Index = newIndex
	return op1
}

func (op InsertText) InsertLength() int {
	return len(op.Text)
}
