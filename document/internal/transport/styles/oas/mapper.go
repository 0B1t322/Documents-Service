package oas

import (
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	dto "github.com/0B1t322/Documents-Service/document/internal/dto/styles"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/samber/lo"
)

type mapper struct {
}

func (m mapper) CreateUpdateStyleInDocument(req *documents.CreateUpdateStyle) dto.StyleDto {
	return dto.StyleDto{
		Name: req.Name,
		ParagraphStyle: dto.ParagraphStyleDto{
			Alignment:   m.DtoAlignment(req.ParagraphStyle.Alignment),
			LineSpacing: req.ParagraphStyle.LineSpacing,
		},
		TextStyle: dto.TextStyleDto{
			FontFamily: req.TextStyle.FontFamily,
			FontWeight: req.TextStyle.FontWeight,
			Bold:       req.TextStyle.Bold,
			Underline:  req.TextStyle.Underline,
			Italic:     req.TextStyle.Italic,
			FontSize: models.Dimension{
				Magnitude: req.TextStyle.FontSize.Magnitude,
				Unit:      models.Unit(req.TextStyle.FontSize.Unit),
			},
			BackgroundColor: dto.Color{
				Red:   req.TextStyle.BackgroundColor.Red,
				Blue:  req.TextStyle.BackgroundColor.Blue,
				Green: req.TextStyle.BackgroundColor.Green,
			},
			ForegroundColor: dto.Color{
				Red:   req.TextStyle.ForegroundColor.Red,
				Blue:  req.TextStyle.ForegroundColor.Blue,
				Green: req.TextStyle.ForegroundColor.Green,
			},
		},
	}
}

func (mapper) DtoAlignment(alignment documents.Alignment) dto.Alignment {
	switch alignment {
	case documents.AlignmentStart:
		return dto.AlignmentStart
	case documents.AlignmentEnd:
		return dto.AlignmentEnd
	case documents.AlignmentCenter:
		return dto.AlignmentCenter
	case documents.AlignmentJustified:
		return dto.AlignmentJustified
	default:
		return dto.AlignmentUnknown
	}
}

func (m mapper) Styles(styles []models.Style) *documents.DocumentsIDStylesGetOKApplicationJSON {
	s := lo.Map(
		styles, func(item models.Style, index int) documents.Style {
			return m.Style(item)
		},
	)

	mapped := documents.DocumentsIDStylesGetOKApplicationJSON(s)

	return lo.ToPtr(mapped)
}

func (m mapper) Style(style models.Style) documents.Style {
	return documents.Style{
		ID:             style.ID,
		Name:           style.Name,
		ParagraphStyle: m.ParagraphStyle(style.ParagraphStyle),
		TextStyle:      m.TextStyle(style.TextStyle),
	}
}

func (m mapper) ParagraphStyle(style models.ParagraphStyle) documents.ParagraphStyle {
	return documents.ParagraphStyle{
		ID:          style.ID,
		Alignment:   m.Alignment(style.Alignment),
		LineSpacing: style.LineSpacing,
	}
}

func (m mapper) Alignment(alignment models.Alignment) documents.Alignment {
	switch alignment {
	case models.ALIGNMENT_START:
		return documents.AlignmentStart
	case models.ALIGNMENT_END:
		return documents.AlignmentEnd
	case models.ALIGNMENT_JUSTIFIED:
		return documents.AlignmentJustified
	case models.ALIGNMENT_CENTER:
		return documents.AlignmentCenter
	default:
		panic("Unknown alignment type")
	}
}

func (m mapper) TextStyle(style models.TextStyle) documents.TextStyle {
	return documents.TextStyle{
		ID:         style.ID,
		FontFamily: style.FontFamily,
		FontWeight: style.FontWeight,
		FontSize: documents.Dimension{
			Magnitude: style.FontSize.Magnitude,
			Unit:      documents.Unit(style.FontSize.Unit),
		},
		Bold:            style.Bold,
		Underline:       style.Underline,
		Italic:          style.Italic,
		BackgroundColor: m.Color(style.BackgroundColor),
		ForegroundColor: m.Color(style.ForegroundColor),
	}
}

func (m mapper) Color(color models.Color) documents.Color {
	return documents.Color{
		Red:   color.Red,
		Blue:  color.Blue,
		Green: color.Green,
	}
}
