package services

import (
	"github.com/0B1t322/Documents-Service/internal/core/utils/cursor"
	"github.com/go-faster/errors"
)

var (
	ErrStructuralElementNotFound = errors.New("Structural element not found")
	ErrSEChildElementBadType     = errors.New("Structural element contains a child element with another type")
	ErrSEChildIsNotParagraph     = errors.New("Structural element is not contains paragraph")
	ErrStructuralElementBadIndex = errors.New("Structural element can have this index")
	ErrParagraphElementBadIndex  = errors.New("Paragraph element can have this index")
	ErrParagraphElementNotFound  = errors.New("Paragraph element not found")
	ErrPEChildElementBadType     = errors.New("Paragraph element contains a child element with another type")
)

func (ElementsService) IsNotFound(err error) bool {
	return errors.Is(err, ErrStructuralElementNotFound) ||
		errors.Is(err, ErrParagraphElementNotFound)
}

func (s ElementsService) IsValidation(err error) bool {
	return errors.Is(err, ErrSEChildElementBadType) ||
		errors.Is(err, ErrStructuralElementBadIndex) ||
		errors.Is(err, ErrSEChildIsNotParagraph) ||
		errors.Is(err, ErrParagraphElementBadIndex) ||
		errors.Is(err, ErrPEChildElementBadType)
}

func (ElementsService) IsParametersNotValid(err error) bool {
	if _, ok := errors.Into[*cursor.CursorDecodeError](err); ok {
		return ok
	}

	return false
}
