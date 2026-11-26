package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation error")
	ErrBadRequest   = errors.New("bad request")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

type ValidationErrors struct {
	List []*ValidationError
}

func (e *ValidationErrors) Error() string {
	if len(e.List) == 0 {
		return "validation error"
	}
	out := ""
	for i, it := range e.List {
		if i > 0 {
			out += "; "
		}
		out += it.Error()
	}
	return out
}
