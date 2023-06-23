package middleware

//go:generate easyjson -all -snake_case error_response.go

type ErrorResponse struct {
	Message string `json:"message"`
}
