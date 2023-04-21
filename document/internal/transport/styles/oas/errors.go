package oas

import (
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"net/http"
)

func NotFound(err error) (*documents.ErrorStatusCode, error) {
	return Status(http.StatusFound, err.Error()), nil
}

func BadRequest(err error) (*documents.ErrorStatusCode, error) {
	return Status(http.StatusBadRequest, err.Error()), nil
}

func FailedToCreateStyleInDocument() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to create style in document"), nil
}

func FailedToDeleteStyleInDocument() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to delete style in document"), nil
}

func FailedToGetStylesInDocument() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to get styles in document"), nil
}

func FailedToUpdateInDocument() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to update style in document"), nil
}
