package rest

import (
	"context"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	repository "github.com/0B1t322/Documents-Service/sessions/internal/repository/elements"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"net/http"
)

type Repository struct {
	client *documents.Client
}

func New(client *documents.Client) *Repository {
	return &Repository{client: client}
}

func (r Repository) GetParagraphElementByIndexes(
	ctx context.Context,
	documentId uuid.UUID,
	seId, peId int,
) (models.ParagraphElement, error) {
	res, err := r.client.GetParagraphElementByIndexes(
		ctx, documents.GetParagraphElementByIndexesParams{
			ID:                     documentId,
			StructuralElementIndex: seId,
			ParagraphElementIndex:  peId,
		},
	)
	if err != nil {
		return models.ParagraphElement{}, err
	}

	var element models.ParagraphElement
	switch res := res.(type) {
	case *documents.ParagraphElement:
		element = models.ParagraphElement{
			ID:    res.ID,
			Index: res.Index,
		}

		if textRun, ok := res.Element.GetTextRun(); ok {
			element.TextRun = &models.TextRun{
				ID:      textRun.ID,
				Content: textRun.Content,
				TextStyle: models.TextStyle{
					ID: textRun.TextStyleId,
				},
			}
		}
	case *documents.ErrorStatusCode:
		if res.StatusCode == http.StatusFound {
			return models.ParagraphElement{}, repository.ErrParagraphElementNotFound
		}
	}

	return element, nil
}

func (r Repository) UpdateParagraphElementByIndexes(
	ctx context.Context,
	documentId uuid.UUID,
	seId, peId int,
	element models.ParagraphElement,
) error {
	params := documents.UpdateParagraphElementByIndexesParams{
		ID:                     documentId,
		StructuralElementIndex: seId,
		ParagraphElementIndex:  peId,
	}

	req := &documents.UpdateParagraphElement{}
	{
		req.SetElement(
			documents.UpdateParagraphElementElement{
				Type: documents.CreateUpdateTextRunUpdateParagraphElementElement,
				CreateUpdateTextRun: documents.CreateUpdateTextRun{
					Content:     element.TextRun.Content,
					TextStyleId: element.TextRun.TextStyle.ID,
				},
			},
		)
	}

	_, err := r.client.UpdateParagraphElementByIndexes(
		ctx,
		req,
		params,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) IsNotFound(err error) bool {
	return errors.Is(err, repository.ErrParagraphElementNotFound)
}
