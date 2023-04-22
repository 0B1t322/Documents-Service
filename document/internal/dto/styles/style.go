package dto

import "github.com/0B1t322/Documents-Service/document/internal/core/models"

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
		FontSize        models.Dimension
		Bold            bool
		Underline       bool
		Italic          bool
		BackgroundColor Color
		ForegroundColor Color
	}
)
