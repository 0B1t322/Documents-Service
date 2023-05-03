package app

import (
	"github.com/0B1t322/Documents-Service/sessions/internal/app/operations"
	"github.com/0B1t322/Documents-Service/sessions/internal/repository/elements/rest"
	"github.com/0B1t322/Documents-Service/sessions/internal/repository/operations/influxdb"
)

type (
	Operations = operations.App
)

var (
	NewOperations = operations.New
)

var (
	NewInfluxOperationRepository = influxdb.New
	NewRestElementsRepository    = rest.New
)
