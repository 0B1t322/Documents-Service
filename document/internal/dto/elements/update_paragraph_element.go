package dto

type (
	UpdateParagraphElementDto struct {
		ID          int
		ElementType ParagraphElementType

		TextRun UpdateTextRun
	}

	UpdateTextRun struct {
		Content     string
		TextStyleID int
	}
)
