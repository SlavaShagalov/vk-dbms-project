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

	// Forum
	ErrForumNotFound      = errors.New("forum not found")
	ErrForumAlreadyExists = errors.New("forum already exists")

	// Thread
	ErrThreadAlreadyExists = errors.New("thread already exists")

	// Params
	ErrInvalidLimitParam = errors.New("invalid limit param")

	// HTTP
	ErrReadBody = errors.New("read request body error")

	// JSON
	ErrParseJson = errors.New("parse json error")
)
