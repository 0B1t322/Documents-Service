package app

import (
	"context"
	"errors"
	"github.com/0B1t322/Documents-Service/document/internal/config"
	"github.com/0B1t322/Documents-Service/document/internal/core/events"
	documentsPgql "github.com/0B1t322/Documents-Service/document/internal/repository/documents/postgresql"
	elementsPgql "github.com/0B1t322/Documents-Service/document/internal/repository/elements/pgql"
	stylesPgql "github.com/0B1t322/Documents-Service/document/internal/repository/styles/pgql"
	"github.com/0B1t322/Documents-Service/internal/database/pgql"
	zapLog "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/ogen-go/ogen/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type closer interface {
	Close()
}

type closerFunc func() error

func (c closerFunc) Close() error {
	return c()
}

type App struct {
	closers         []io.Closer
	cfg             config.Config
	logger          log.Logger
	eventsPublisher events.EventPublisher

	documentsRepository DocumentsRepository
	elementsRepository  ElementsRepository
	stylesRepository    StylesRepository

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

	if err := app.initEventPublisher(); err != nil {
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
	a.closers = append(a.closers, closerFunc(func() error { pool.Close(); return nil }))

	a.documentsRepository = documentsPgql.New(pool)
	a.elementsRepository = elementsPgql.New(pool)
	a.stylesRepository = stylesPgql.New(pool)

	return nil
}

func (a *App) initEventPublisher() error {
	conn, err := amqp.Dial(a.cfg.AMQPUrl)
	if err != nil {
		return err
	}
	a.closers = append(a.closers, conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	a.closers = append(a.closers, ch)

	pub, err := NewAMQPEventPublisher(ch, a.logger, a.cfg.AMQPExchangeName)
	if err != nil {
		return err
	}

	a.eventsPublisher = pub

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
	a.Documents = NewDocumentApp(a.documentsRepository, a.logger, a.eventsPublisher)
	a.Elements = NewElementsApp(a.elementsRepository, a.logger, a.eventsPublisher)
	a.Styles = NewStylesApp(a.stylesRepository, a.logger, a.eventsPublisher)

	return nil
}

func (a *App) Close() error {
	var errs []error
	for _, c := range a.closers {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
