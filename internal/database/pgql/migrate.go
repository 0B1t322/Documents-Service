package pgql

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
)

func Migrate(_ context.Context, url, source string) error {
	m, err := migrate.New(
		source,
		url,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
