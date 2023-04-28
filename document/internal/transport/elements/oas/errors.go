package oas

import (
	"fmt"
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"
	"github.com/google/uuid"
	"net/http"
)

func NotFound(err error) (*documents.ErrorStatusCode, error) {
	return Status(http.StatusNotFound, err.Error()), nil
}

func DocumentNotFound(documentId uuid.UUID) (*documents.ErrorStatusCode, error) {
	return Status(http.StatusNotFound, fmt.Sprintf("Not found document with id=%s", documentId)), nil
}

func StructuralElementNotFound(id int) (*documents.ErrorStatusCode, error) {
	return Status(http.StatusNotFound, fmt.Sprintf("Not found structural element with id=%s", id)), nil
}

func StructuralElementNotFoundInDocument(id int, documentId uuid.UUID) (*documents.ErrorStatusCode, error) {
	return Status(
		http.StatusNotFound,
		fmt.Sprintf("Not found structural element with id=%v in document with id=%s", id, documentId),
	), nil
}

func ParagraphElementNotFoundInDocument(id int, seId int, documentId uuid.UUID) (*documents.ErrorStatusCode, error) {
	return Status(
		http.StatusNotFound,
		fmt.Sprintf(
			"Paragraph element with id=%s not found in document with id=%s and is structural element with id"+
				"=%s", id, documentId, seId,
		),
	), nil
}

func FailedToGetElements() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to get elements"), nil
}

func BadRequest(err error) (*documents.ErrorStatusCode, error) {
	return Status(http.StatusBadRequest, err.Error()), nil
}

func FailedDeleteStructuralElementByID() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to delete structural element by id"), nil
}

func FailedToCreateElement() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to create element"), nil
}

func FailedToUpdateElement() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to update element"), nil
}

func FailedToCreateParagraphElement() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to create paragraph element"), nil
}

func FailedToDeleteParagraphElement() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to delete paragraph element"), nil
}

func FailedToUpdateParagraphElement() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to update paragraph element"), nil
}

func FailedToGetParagraphElements() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to get paragraph elements"), nil
}

func FailedToGetParagraphElementByIndexes() (*documents.ErrorStatusCode, error) {
	return Status(http.StatusInternalServerError, "Failed to get paragraph element by indexes"), nil
}
