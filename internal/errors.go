package internal

import "errors"

var (
	ErrUnauthorized = errors.New("monime: unauthorized")
	ErrBadRequest   = errors.New("monime: bad request")
	ErrServerError  = errors.New("monime: internal server error")
)