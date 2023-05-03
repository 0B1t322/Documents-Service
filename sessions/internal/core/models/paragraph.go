package models

type Paragraph struct {
	ID             int
	ParagraphStyle ParagraphStyle
	Elements       []ParagraphElement
}
