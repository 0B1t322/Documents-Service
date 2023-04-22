package repository

import "github.com/go-faster/errors"

var (
	ErrStyleNotFound = errors.New("Style not found")
	ErrStyleExist    = errors.New("Style with this name exist")
)
