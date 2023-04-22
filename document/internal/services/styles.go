package services

import (
	"context"
	"github.com/0B1t322/Documents-Service/document/internal/core/events"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	dto "github.com/0B1t322/Documents-Service/document/internal/dto/styles"
	repository "github.com/0B1t322/Documents-Service/document/internal/repository/styles"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
)

type (
	stylesRepository interface {
		StoreStyleInDocument(ctx context.Context, documentId uuid.UUID, style *models.Style) error
		FindStyleInDocumentByID(ctx context.Context, documentId uuid.UUID, styleId uuid.UUID) (models.Style, error)
		UpdateStyle(ctx context.Context, style models.Style) error
		GetAllStylesInDocument(ctx context.Context, documentId uuid.UUID) ([]models.Style, error)
		DeleteStyleInDocument(ctx context.Context, documentId uuid.UUID, style models.Style) error
	}
)

type StylesService struct {
	repository stylesRepository
	logger     log.Logger
	publisher  events.EventPublisher
}

func NewStylesService(repository stylesRepository, logger log.Logger, publisher events.EventPublisher) *StylesService {
	return &StylesService{
		repository: repository,
		logger:     log.WithPrefix(logger, "service", "StylesService"),
		publisher:  publisher,
	}
}

func alignmentFromDto(alignment dto.Alignment) models.Alignment {
	switch alignment {
	case dto.AlignmentCenter:
		return models.ALIGNMENT_CENTER
	case dto.AlignmentJustified:
		return models.ALIGNMENT_JUSTIFIED
	case dto.AlignmentEnd:
		return models.ALIGNMENT_END
	case dto.AlignmentStart:
		return models.ALIGNMENT_START
	}

	panic("Unknown alignment type")
}

func colorFromDto(color dto.Color) models.Color {
	return models.Color{
		Red:   color.Red,
		Blue:  color.Blue,
		Green: color.Green,
	}
}

func styleFromDto(dto dto.StyleDto) models.Style {
	return models.Style{
		Name: dto.Name,
		ParagraphStyle: models.ParagraphStyle{
			Alignment:   alignmentFromDto(dto.ParagraphStyle.Alignment),
			LineSpacing: dto.ParagraphStyle.LineSpacing,
		},
		TextStyle: models.TextStyle{
			FontSize:        dto.TextStyle.FontSize,
			FontFamily:      dto.TextStyle.FontFamily,
			FontWeight:      dto.TextStyle.FontWeight,
			Bold:            dto.TextStyle.Bold,
			Underline:       dto.TextStyle.Underline,
			Italic:          dto.TextStyle.Italic,
			BackgroundColor: colorFromDto(dto.TextStyle.BackgroundColor),
			ForegroundColor: colorFromDto(dto.TextStyle.ForegroundColor),
		},
	}
}

func updateFromDto(style models.Style, dto dto.StyleDto) models.Style {
	fromDto := styleFromDto(dto)
	fromDto.ID = style.ID
	fromDto.ParagraphStyle.ID = style.ParagraphStyle.ID
	fromDto.TextStyle.ID = style.TextStyle.ID

	return fromDto
}

func (s StylesService) CreateStyleInDocument(
	ctx context.Context, documentId uuid.UUID,
	dto dto.StyleDto,
) (models.Style, error) {
	style := styleFromDto(dto)

	if err := s.repository.StoreStyleInDocument(ctx, documentId, &style); err == repository.ErrStyleExist {
		return models.Style{}, ErrStyleWithThisNameExist
	} else if err != nil {
		level.Error(s.logger).Log("description", "can't Store Style In Document", "err", err)
		return models.Style{}, err
	}

	return style, nil
}

func (s StylesService) UpdateStyleInDocument(
	ctx context.Context, documentId uuid.UUID, styleId uuid.UUID,
	styleDto dto.StyleDto,
) (models.Style, error) {
	style, err := s.repository.FindStyleInDocumentByID(ctx, documentId, styleId)
	if err == repository.ErrStyleNotFound {
		return models.Style{}, ErrNotFoundStyleInDocument
	} else if err != nil {
		return models.Style{}, err
	}

	updated := updateFromDto(style, styleDto)

	if err := s.repository.UpdateStyle(ctx, updated); err == repository.ErrStyleExist {
		return models.Style{}, ErrStyleWithThisNameExist
	} else if err != nil {
		return models.Style{}, err
	}

	return updated, nil
}

func (s StylesService) GetStylesInDocument(ctx context.Context, documentId uuid.UUID) ([]models.Style, error) {
	styles, err := s.repository.GetAllStylesInDocument(ctx, documentId)
	if err != nil {
		level.Error(s.logger).Log(
			"description", "Can't get styles in document",
			"err", err,
		)
		return nil, err
	}

	return styles, nil
}

func (s StylesService) DeleteStyleInDocument(ctx context.Context, documentId uuid.UUID, styleId uuid.UUID) error {
	style, err := s.repository.FindStyleInDocumentByID(ctx, documentId, styleId)
	if err == repository.ErrStyleNotFound {
		return ErrNotFoundStyleInDocument
	} else if err != nil {
		return err
	}

	if err := s.repository.DeleteStyleInDocument(ctx, documentId, style); err != nil {
		return err
	}

	return nil
}
