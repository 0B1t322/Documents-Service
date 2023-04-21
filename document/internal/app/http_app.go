package app

import (
	"github.com/0B1t322/Documents-Service/document/internal/config"
	documentsOas "github.com/0B1t322/Documents-Service/document/internal/transport/documents/oas"
	elementsOas "github.com/0B1t322/Documents-Service/document/internal/transport/elements/oas"
	stylesOas "github.com/0B1t322/Documents-Service/document/internal/transport/styles/oas"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"net/http"
)

// Alias HTTP Controller
type (
	DocumentsHttpController = documentsOas.DocumentsController
	ElementsHttpController  = elementsOas.ElementsController
	StylesHttpController    = stylesOas.StylesController
)

var (
	NewDocumentHttpController = documentsOas.NewDocumentController
	NewElementsHttpController = elementsOas.NewElementsController
	NewStylesHttpController   = stylesOas.NewStylesController
)

type HTTPApp struct {
	app *App

	HTTPController
}

type HTTPController struct {
	*DocumentsHttpController
	*ElementsHttpController
	*StylesHttpController
}

func NewHTTPAppFromConfig(cfg config.Config) (*HTTPApp, error) {
	app, err := NewAppFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &HTTPApp{
		app: app,
		HTTPController: HTTPController{
			DocumentsHttpController: NewDocumentHttpController(app.Documents),
			ElementsHttpController:  NewElementsHttpController(app.Elements, app.Documents),
			StylesHttpController:    NewStylesHttpController(),
		},
	}, nil
}

func (a *HTTPApp) ToHandler(basePath string) (http.Handler, error) {
	s, err := documents.NewServer(a, documents.WithPathPrefix(basePath))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (a HTTPApp) Shutdown() {
	a.app.dbCloser.Close()
}
