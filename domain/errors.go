package domain

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Your requested item is not found")
	ErrConflict            = errors.New("Your item already exist")
	ErrBadParamInput       = errors.New("Given param is not valid")
	ErrLinkIsExists        = errors.New("Link is exists")
)
