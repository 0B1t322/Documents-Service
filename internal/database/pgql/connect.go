package pgql

import (
	"context"
	"github.com/go-kit/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type connectionPoolOptions struct {
	logger log.Logger
}

func (c connectionPoolOptions) apply(cfg *pgxpool.Config) {

}

type ConnectionPoolOptions func(options *connectionPoolOptions)

func NewConnectionPool(ctx context.Context, url string, opts ...ConnectionPoolOptions) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		if err := initTypes(conn); err != nil {
			return err
		}

		return nil
	}

	return pgxpool.NewWithConfig(ctx, config)
}
