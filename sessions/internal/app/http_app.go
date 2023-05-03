package app

import (
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/sessions"
	"github.com/0B1t322/Documents-Service/sessions/internal/config"
	operationsOas "github.com/0B1t322/Documents-Service/sessions/internal/transport/operations/oas"
	"net/http"
)

type (
	OperationsController = operationsOas.OperationsController
)

var (
	NewOperationsController = operationsOas.New
)

type HTTPApp struct {
	app *App

	HTTPController
}

type HTTPController struct {
	*OperationsController
}

func NewHTTPFromConfig(cfg config.Config) (*HTTPApp, error) {
	app, err := NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &HTTPApp{
		app: app,
		HTTPController: HTTPController{
			OperationsController: NewOperationsController(app),
		},
	}, nil
}

func (a *HTTPApp) ToHandler(basePath string) (http.Handler, error) {
	s, err := sessions.NewServer(a, sessions.WithPathPrefix(basePath))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (a HTTPApp) Shutdown() {

}
