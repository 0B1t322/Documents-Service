package services

import "github.com/go-faster/errors"

var (
	ErrNotFoundStyleInDocument = errors.New("Not found styles in document")
	ErrStyleWithThisNameExist  = errors.New("Style with this name exist")
)

func (s StylesService) IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFoundStyleInDocument)
}

func (s StylesService) IsValidation(err error) bool {
	return errors.Is(err, ErrStyleWithThisNameExist)
}
