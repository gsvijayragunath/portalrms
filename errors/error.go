package errors

import "errors"

const (
	SystemError     = "system-error"
	RecordNotFound  = "record-not-found"
	DatabaseError   = "database-error"
	ClientError     = "client-error"
	Unauthorized    = "unauthorized-error"
	Conflict        = "conflict-error"
	StatusForbidden = "forbidden"
)

var (
	ErrInvalidRequest = errors.New(ClientError)
	ErrUnauthorized   = errors.New(Unauthorized)
	ErrNotFound       = errors.New(RecordNotFound)
	ErrDatabase       = errors.New(DatabaseError)
	ErrConflict       = errors.New(Conflict)
	ErrForbidden      = errors.New(StatusForbidden)
)
