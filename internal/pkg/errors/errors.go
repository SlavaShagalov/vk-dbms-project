package errors

import (
	"errors"
)

var (
	// Common
	ErrInternal = errors.New("internal error")

	// User
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// HTTP
	ErrReadBody = errors.New("read request body error")

	// JSON
	ErrParseJson = errors.New("parse json error")
)
