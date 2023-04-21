package elements

import (
	"github.com/0B1t322/Documents-Service/document/internal/core/events"
	repository "github.com/0B1t322/Documents-Service/document/internal/repository/elements"
	"github.com/0B1t322/Documents-Service/document/internal/services"
	"github.com/go-kit/log"
)

type (
	Repository = repository.Repository
)

type App struct {
	*services.ElementsService
}

func New(repository Repository, logger log.Logger, publisher events.EventPublisher) *App {
	app := &App{
		ElementsService: services.NewElementsService(repository, publisher, logger),
	}

	return app
}
