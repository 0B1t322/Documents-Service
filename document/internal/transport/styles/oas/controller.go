package oas

import (
	"context"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
)

type StylesController struct {
}

func NewStylesController() *StylesController {
	return &StylesController{}
}

func (s StylesController) CreateDocumentStyle(
	ctx context.Context,
	req *documents.CreateUpdateStyle,
	params documents.CreateDocumentStyleParams,
) (documents.CreateDocumentStyleRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s StylesController) DeleteStyleById(
	ctx context.Context,
	params documents.DeleteStyleByIdParams,
) (documents.DeleteStyleByIdRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s StylesController) DocumentsIDStylesGet(
	ctx context.Context,
	params documents.DocumentsIDStylesGetParams,
) (documents.DocumentsIDStylesGetRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s StylesController) UpdateStyleById(
	ctx context.Context,
	req *documents.CreateUpdateStyle,
	params documents.UpdateStyleByIdParams,
) (documents.UpdateStyleByIdRes, error) {
	//TODO implement me
	panic("implement me")
}
