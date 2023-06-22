package errors

import "net/http"

var httpCodes = map[error]int{
	// Common
	ErrInternal: http.StatusInternalServerError,

	// User
	ErrUserNotFound:      http.StatusNotFound,
	ErrUserAlreadyExists: http.StatusConflict,

	// Forum
	ErrForumNotFound:      http.StatusNotFound,
	ErrForumAlreadyExists: http.StatusConflict,

	// Thread
	ErrThreadNotFound:      http.StatusNotFound,
	ErrThreadAlreadyExists: http.StatusConflict,

	// Voice
	ErrVoiceNotFound:      http.StatusNotFound,
	ErrVoiceAlreadyExists: http.StatusConflict,

	// Post
	ErrPostNotFound:       http.StatusNotFound,
	ErrParentPostNotFound: http.StatusNotFound,

	// Params
	ErrInvalidIDParam:    http.StatusBadRequest,
	ErrInvalidLimitParam: http.StatusBadRequest,

	// HTTP
	ErrReadBody: http.StatusBadRequest,

	// JSON
	ErrParseJSON: http.StatusBadRequest,
}

func GetHTTPCodeByError(err error) (int, bool) {
	httpCode, exist := httpCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
