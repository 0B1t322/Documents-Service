package oas

import (
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/sessions"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/transport/oas"
	"net/http"
)

var (
	Status = oas.Status
)

func NotFound(err error) (*sessions.ErrorStatusCode, error) {
	return Status(http.StatusNotFound, err.Error()), nil
}

func FailedToSaveOperations() (*sessions.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to save operations"), nil
}

func FailedGetDocumentRevision() (*sessions.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to get last revision id"), nil
}

func FailedToGetDocumentHistory() (*sessions.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to get document history"), nil
}

func FailedToSyncDocument() (*sessions.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to get sync documents"), nil
}
