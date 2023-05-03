package models

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations/elements"
	"github.com/samber/lo"
)

type StructuralElement struct {
	ID           int
	Index        int
	Paragraph    *Paragraph
	SectionBreak *SectionBreak
}

func (s StructuralElement) GetParagraphElement(i int) elements.ParagraphElement {
	// TODO Handle not find
	element, _ := lo.Find(
		s.Paragraph.Elements, func(item ParagraphElement) bool {
			return item.Index == i
		},
	)

	return &element
}
