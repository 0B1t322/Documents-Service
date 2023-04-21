package dto

type (
	CreateParagraphElementDto struct {
		StructuralElementID int
		Index               int
		ElementType         ParagraphElementType

		TextRun CreateTextRun
	}

	CreateTextRun struct {
		Content     string
		TextStyleID int
	}
)
