package dto

type (
	UpdateParagraphElementDto struct {
		ElementType ParagraphElementType

		TextRun UpdateTextRun
	}

	UpdateTextRun struct {
		Content     string
		TextStyleID int
	}
)
