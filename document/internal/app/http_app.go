package app

import (
	"github.com/0B1t322/Online-Document-Redactor/document/internal/config"
	"net/http"
)

type closer interface {
	Close()
}

type OASApp struct {
	app     *App
	handler http.Handler
}

func NewOASAppFromConfig(cfg config.Config) (*OASApp, error) {
	app, err := NewAppFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &OASApp{
		app: app,
	}, nil
}

func (a OASApp) Shutdown() {
	a.app.dbCloser.Close()
}
