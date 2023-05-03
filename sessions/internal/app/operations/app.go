package operations

import (
	elemRepository "github.com/0B1t322/Documents-Service/sessions/internal/repository/elements"
	opRepository "github.com/0B1t322/Documents-Service/sessions/internal/repository/operations"
	"github.com/0B1t322/Documents-Service/sessions/internal/services/operations"
)

type (
	OperationsRepository = opRepository.Repository
	ElementsRepository   = elemRepository.Repository
)

type App struct {
	*operations.Service
}

func New(opRepository OperationsRepository, elemRepository ElementsRepository) *App {
	app := &App{
		Service: operations.New(elemRepository, opRepository),
	}

	return app
}
