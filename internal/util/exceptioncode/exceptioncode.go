package exceptioncode

import "errors"

var (
	ErrEmptyResult    = errors.New("empty result")
	ErrInvalidRequest = errors.New("invalid request")

	// spesific postgre error
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrUniqueViolation     = errors.New("unique violation")
)

const (
	CodeDataNotFound        = "DATA_NOT_FOUND"
	CodeInvalidRequest      = "INVALID_REQUEST"
	CodeInvalidValidation   = "INVALID_VALIDATION"
	CodeBadRequest          = "BAD_REQUEST"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

type (
	errorType struct {
		ErrorMessage string
	}
	ErrorNotFound            errorType
	ErrorForeignKeyViolation errorType
)

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}
