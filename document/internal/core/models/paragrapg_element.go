package models

type ParagraphElement struct {
	ID         int
	Index      int
	StartIndex int
	EndIndex   int

	TextRune *TextRune
}

func (p ParagraphElement) GetType() PEType {
	if p.TextRune != nil {
		return PETextRune
	}

	return PEUnknown
}

func (p ParagraphElement) GetChildElementID() int {
	switch p.GetType() {
	case PETextRune:
		return p.TextRune.ID
	}

	panic("Uknown element type")
}

type PEType int

const (
	PEUnknown PEType = iota
	PETextRune
)
