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
	ErrThreadNotFound      = errors.New("thread not found")
	ErrThreadAlreadyExists = errors.New("thread already exists")

	// Voice
	ErrVoiceNotFound      = errors.New("voice not found")
	ErrVoiceAlreadyExists = errors.New("voice already exists")

	// Post
	ErrPostNotFound       = errors.New("post not found")
	ErrParentPostNotFound = errors.New("parent post not found")

	// Params
	ErrInvalidIDParam    = errors.New("invalid id param")
	ErrInvalidLimitParam = errors.New("invalid limit param")

	// HTTP
	ErrReadBody = errors.New("read request body error")

	// JSON
	ErrParseJSON = errors.New("parse json error")
)
