package styles

import (
	"github.com/0B1t322/Documents-Service/document/internal/core/events"
	repository "github.com/0B1t322/Documents-Service/document/internal/repository/styles"
	"github.com/0B1t322/Documents-Service/document/internal/services"
	"github.com/go-kit/log"
)

type (
	Repository = repository.Repository
)

type App struct {
	*services.StylesService
}

func New(repository repository.Repository, logger log.Logger, publisher events.EventPublisher) *App {
	app := &App{
		StylesService: services.NewStylesService(repository, logger, publisher),
	}

	return app
}
