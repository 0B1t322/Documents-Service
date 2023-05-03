package text

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations"
	"strings"
)

type DeleteText struct {
	Index int
	Text  string
}

func (d DeleteText) GetIndex() int {
	return d.Index
}

func (d DeleteText) SetIndex(i int) operations.Delete {
	d.Index = i
	return d
}

func (op1 DeleteText) IncludeDelete(delete operations.Delete) operations.Delete {
	op, ok := delete.(DeleteText)
	if !ok {
		panic("Unsupported insert operation")
	}

	return op1.IncludeDeleteText(op)
}

func (op1 DeleteText) IncludeDeleteText(op2 DeleteText) DeleteText {
	if op1.GetIndex() < op2.GetIndex() {
		return op1
	}

	newIndex := op1.GetIndex() - op2.DeleteLength()
	if newIndex < 0 {
		newIndex = 0
	}
	op1.Index = newIndex
	return op1
}

func (op1 DeleteText) IncludeInsert(insert operations.Insert) operations.Delete {
	op, ok := insert.(InsertText)
	if !ok {
		panic("Unsupported insert operation")
	}

	return op1.IncludeInsertText(op)
}

func (op1 DeleteText) IncludeInsertText(op2 operations.Insert) DeleteText {
	if op1.GetIndex() < op2.GetIndex() {
		return op1
	}

	newIndex := op1.GetIndex() + op2.InsertLength()
	op1.Index = newIndex
	return op1
}

func (d DeleteText) DeleteLength() int {
	return len(d.Text)
}

func NewDeleteText(index int, text string) DeleteText {
	return DeleteText{
		Index: index,
		Text:  text,
	}
}

func (op DeleteText) Apply(to string) string {
	b := strings.Builder{}

	b.WriteString(to[:op.Index])

	textToDelete := to[op.Index:]

	for i := range textToDelete {
		if i == len(op.Text) {
			b.WriteString(textToDelete[i:])
			break
		}
		if textToDelete[i] == op.Text[i] {

		} else {
			b.WriteByte(textToDelete[i])
		}
	}

	return b.String()
}
