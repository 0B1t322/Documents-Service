package oas

import (
	"fmt"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/google/uuid"
	"net/http"
)

func DocumentNotFound(givenId uuid.UUID) (*documents.ErrorStatusCode, error) {
	return &documents.ErrorStatusCode{
		StatusCode: http.StatusNotFound,
		Response: documents.Error{
			Status:  http.StatusNotFound,
			Message: fmt.Sprintf("Not found document with id=%s", givenId),
		},
	}, nil
}

func BadRequest(err error) (*documents.ErrorStatusCode, error) {
	return &documents.ErrorStatusCode{
		StatusCode: http.StatusBadRequest,
		Response: documents.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		},
	}, nil
}

func FailedToGetDocument() (*documents.ErrorStatusCode, error) {
	return &documents.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: documents.Error{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get document",
		},
	}, nil
}

func FailedToCreateDocument() (*documents.ErrorStatusCode, error) {
	return &documents.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: documents.Error{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create document",
		},
	}, nil
}

func FailedToUpdateDocument() (*documents.ErrorStatusCode, error) {
	return &documents.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: documents.Error{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update document",
		},
	}, nil
}

func FailedToGetDocuments() (*documents.ErrorStatusCode, error) {
	return &documents.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: documents.Error{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get documents",
		},
	}, nil
}
