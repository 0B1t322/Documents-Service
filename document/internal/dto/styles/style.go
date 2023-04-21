package dto

type (
	StyleDto struct {
		Name           string
		ParagraphStyle ParagraphStyleDto
		TextStyle      TextStyleDto
	}

	ParagraphStyleDto struct {
		Alignment   Alignment
		LineSpacing int
	}

	TextStyleDto struct {
		FontFamily      string
		FontWeight      int
		Bold            bool
		Underline       bool
		Italic          bool
		BackgroundColor Color
		ForegroundColor Color
	}
)
