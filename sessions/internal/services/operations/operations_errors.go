package operations

import "github.com/go-faster/errors"

var (
	ErrParagraphElementNotFound     = errors.New("Paragraph element not found")
	ErrNotFoundRevision             = errors.New("Not found revision for document")
	ErrNotFoundRevisionsForDocument = errors.New("Not found revision for document")
)

func (s Service) IsNotFound(err error) bool {
	return errors.Is(err, ErrParagraphElementNotFound) ||
		errors.Is(err, ErrNotFoundRevision) ||
		errors.Is(err, ErrNotFoundRevisionsForDocument)
}
