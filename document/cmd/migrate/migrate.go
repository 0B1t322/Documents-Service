package main

import (
	"github.com/0B1t322/Documents-Service/document/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func main() {
	log.Println("Start database migrations")

	if err := config.FromEnv(); err != nil {
		log.Fatalf("Failed to parse configs: %s", err)
	}

	source, ok := os.LookupEnv("DOCUMENTS_APP_MIGRATIONS_DIR")
	if !ok {
		source = "document/internal/migrations/pgql"
	}

	m, err := migrate.New(
		"file://"+source,
		config.GlobalConfig.DatabaseURL,
	)
	if err != nil {
		log.Fatalf("Failed create migrations: %s", err)
	}
	defer m.Close()

	if err := m.Up(); err == migrate.ErrNoChange {
		return
	} else {
		log.Fatalf("Failed to run migrations: %s", err)
	}

}
