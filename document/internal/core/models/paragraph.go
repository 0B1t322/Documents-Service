package models

type Paragraph struct {
	ID               int
	ParagraphStyleId int
	Elements         []ParagraphElement
}
