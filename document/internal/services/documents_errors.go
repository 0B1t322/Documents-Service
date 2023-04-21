package services

import (
	"github.com/0B1t322/Online-Document-Redactor/internal/core/utils/cursor"
	"github.com/go-faster/errors"
)

var (
	ErrDocumentNotFound = errors.New("Document not found")
)

func (DocumentService) IsNotFound(err error) bool {
	return err == ErrDocumentNotFound
}

func (DocumentService) IsParametersNotValid(err error) bool {
	if _, ok := errors.Into[*cursor.CursorDecodeError](err); ok {
		return ok
	}

	return false
}
