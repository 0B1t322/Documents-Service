package elements

type (
	ParagraphElement interface {
		GetContent() string
		SetContent(string)
	}

	StructuralElement interface {
		GetParagraphElement(int) ParagraphElement
	}
)
