package documents

import (
	"github.com/0B1t322/Documents-Service/document/internal/core/events"
	repository "github.com/0B1t322/Documents-Service/document/internal/repository/documents"
	"github.com/0B1t322/Documents-Service/document/internal/services"
	"github.com/go-kit/log"
)

type (
	Repository = repository.Repository
)

type App struct {
	*services.DocumentService
}

func New(repository Repository, logger log.Logger, publisher events.EventPublisher) *App {
	app := &App{
		DocumentService: services.NewDocumentService(repository, publisher, logger),
	}

	return app
}

//func (a *App) HTTPHandler(basePath string) http.Handler {
//	h, err := documents.NewServer(a.httpController, documents.WithPathPrefix(basePath))
//	if err != nil {
//		panic(err)
//	}
//
//	return h
//}
