package errors

import "net/http"

var httpCodes = map[error]int{
	// Common repository
	ErrInternal: http.StatusInternalServerError,
}

func GetHTTPCodeByError(err error) (int, bool) {
	httpCode, exist := httpCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
