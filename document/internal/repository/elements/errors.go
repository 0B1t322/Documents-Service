package repository

import "github.com/go-faster/errors"

var (
	ErrStructuralElementNotFound = errors.New("Structural element not found")
	ErrBadIndex                  = errors.New("Can't create element with this index")
	ErrSEBadType                 = errors.New("Structural element contains different type")

	ErrParagraphElementNotFound = errors.New("Paragraph element not found")
)
