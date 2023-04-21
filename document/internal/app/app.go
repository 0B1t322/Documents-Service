package app

import (
	"context"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/config"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/events"
	documentsPgql "github.com/0B1t322/Online-Document-Redactor/document/internal/repository/documents/postgresql"
	elementsPgql "github.com/0B1t322/Online-Document-Redactor/document/internal/repository/elements/pgql"
	"github.com/0B1t322/Online-Document-Redactor/internal/database/pgql"
	zapLog "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/ogen-go/ogen/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type closer interface {
	Close()
}

type App struct {
	dbCloser        closer
	cfg             config.Config
	logger          log.Logger
	eventsPublisher events.EventPublisher

	documentsRepository DocumentsRepository
	elementsRepository  ElementsRepository

	*Documents
	*Elements
	*Styles
}

func NewAppFromConfig(cfg config.Config) (*App, error) {
	app := &App{
		cfg: cfg,
	}

	if err := app.initLogger(); err != nil {
		return nil, err
	}

	if err := app.initRepository(); err != nil {
		return nil, err
	}

	if err := app.initApps(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initLogger() error {
	var (
		zapLogger *zap.Logger
		level     zapcore.Level
	)
	{
		if a.cfg.Development {
			if logger, err := zap.NewDevelopment(); err != nil {
				return err
			} else {
				zapLogger = logger
			}

			level = zap.DebugLevel
		} else {
			if logger, err := zap.NewProduction(); err != nil {
				return err
			} else {
				zapLogger = logger
			}

			level = zap.InfoLevel
		}
	}
	a.logger = zapLog.NewZapSugarLogger(zapLogger, level)

	return nil
}

func (a *App) initRepository() error {
	pool, err := pgql.NewConnectionPool(context.Background(), a.cfg.DatabaseURL)
	if err != nil {
		return err
	}

	a.dbCloser = pool

	a.documentsRepository = documentsPgql.New(pool)
	a.elementsRepository = elementsPgql.New(pool)

	return nil
}

type defaultEventPublisher struct {
	logger log.Logger
}

func (d defaultEventPublisher) PublishEvent(_ context.Context, event events.Event) {
	//	Pass
	data, _ := json.Marshal(event)
	level.Debug(d.logger).Log("Event", event.Event(), "EventPayload", string(data))
}

func (a *App) initApps() error {
	a.eventsPublisher = defaultEventPublisher{logger: a.logger}

	a.Documents = NewDocumentApp(a.documentsRepository, a.logger, a.eventsPublisher)
	a.Elements = NewElementsApp(a.elementsRepository, a.logger, a.eventsPublisher)
	a.Styles = NewStylesApp()

	return nil
}
