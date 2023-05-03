package config

import (
	"encoding/json"
	"github.com/caarlos0/env/v8"
	"os"
)

type Config struct {
	AppPort                   string `json:"appPort" env:"SESSIONS_APP_PORT"`
	Development               bool   `json:"development" env:"SESSIONS_DEVELOPMENT"`
	AMQPUrl                   string `json:"development" env:"SESSIONS_AMQP_URL"`
	DocumentsAMQPExchangeName string `json:"amqp_exchange_name" env:"SESSIONS_DOCUMENT_AMQP_EXCHANGE_NAME"`
	DocumentsRestBaseURL      string `json:"documentsRestBaseURL" env:"SESSIONS_DOCUMENT_REST_BASE_URL"`

	InfluxDBUrl string `json:"influxDBUrl" env:"SESSIONS_INFLUXDB_URL"`
	InfluxToken string `json:"influxToken" env:"SESSIONS_INFLUXDB_TOKEN"`
}

var (
	GlobalConfig Config
)

func init() {
	GlobalConfig = Config{
		AppPort:                   "8082",
		Development:               true,
		AMQPUrl:                   "amqp://user:password@localhost:5672/",
		DocumentsAMQPExchangeName: "documents-service.events",
		DocumentsRestBaseURL:      "http://localhost:8080/api/documents/v1",
		InfluxDBUrl:               "http://localhost:8086",
		InfluxToken:               "some_token",
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
