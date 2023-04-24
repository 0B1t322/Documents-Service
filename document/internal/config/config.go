package config

import (
	"encoding/json"
	"github.com/caarlos0/env/v8"
	"os"
)

type Config struct {
	DatabaseURL string `json:"databaseUrl" env:"DOCUMENTS_DATABASE_URL"`
	AppPort     string `json:"appPort" env:"DOCUMENTS_APP_PORT"`
	Development bool   `json:"development" env:"DOCUMENTS_DEVELOPMENT"`
}

var (
	GlobalConfig Config
)

func init() {
	GlobalConfig = Config{
		DatabaseURL: "postgres://postgres:password@localhost:5432/Documents?sslmode=disable",
		AppPort:     "8080",
		Development: true,
	}
}

func FromJSON(src string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&GlobalConfig); err != nil {
		return err
	}

	return nil
}

func FromEnv() error {
	if err := env.Parse(&GlobalConfig); err != nil {
		return err
	}

	return nil
}
