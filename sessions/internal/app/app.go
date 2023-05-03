package app

import (
	"github.com/0B1t322/Documents-Service/internal/database/influxdb"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/0B1t322/Documents-Service/sessions/internal/app/operations"
	"github.com/0B1t322/Documents-Service/sessions/internal/config"
)

type App struct {
	cfg config.Config

	elementsRepository   operations.ElementsRepository
	operationsRepository operations.OperationsRepository

	*Operations
}

func NewFromConfig(cfg config.Config) (*App, error) {
	app := &App{
		cfg: cfg,
	}

	if err := app.initRepository(); err != nil {
		return nil, err
	}

	if err := app.initApps(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initRepository() error {
	influxClient, err := influxdb.Connect(a.cfg.InfluxDBUrl, a.cfg.InfluxToken)
	if err != nil {
		return err
	}

	a.operationsRepository = NewInfluxOperationRepository(influxClient, "document_redactor", "documents_operations")

	restClient, err := documents.NewClient(a.cfg.DocumentsRestBaseURL)
	if err != nil {
		return err
	}

	a.elementsRepository = NewRestElementsRepository(restClient)

	return nil
}

func (a *App) initApps() error {
	a.Operations = NewOperations(a.operationsRepository, a.elementsRepository)
	return nil
}
