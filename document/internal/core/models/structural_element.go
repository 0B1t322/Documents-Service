package models

type StructuralElement struct {
	ID         int
	Index      int
	StartIndex int
	EndIndex   int

	Paragraph    *Paragraph
	SectionBreak *SectionBreak
}

// GetType return SEType of holding element
func (se StructuralElement) GetType() SEType {
	if se.Paragraph != nil {
		return SEParagraph
	} else if se.SectionBreak != nil {
		return SESectionBreak
	}

	return SEUnknown
}

// GetElementID return id of holding element
func (se StructuralElement) GetElementID() int {
	switch se.GetType() {
	case SEParagraph:
		return se.Paragraph.ID
	case SESectionBreak:
		return se.SectionBreak.ID
	default:
		return -1
	}
}

// SEType – Structural element type
type SEType int

const (
	SEUnknown SEType = iota
	SEParagraph
	SESectionBreak
)
