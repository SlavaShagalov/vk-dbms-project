package http

import "github.com/SlavaShagalov/vk-dbms-project/internal/models"

//go:generate easyjson -all -snake_case api_models.go

// API requests
type createRequest struct {
	Fullname string
	About    string
	Email    string
}

type updateRequest struct {
	Fullname string
	About    string
	Email    string
}

// API responses
type createResponse struct {
	ID       int
	Nickname string
	Fullname string
	About    string
	Email    string
}

func newCreateResponse(user *models.User) *createResponse {
	return &createResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		About:    user.About,
		Email:    user.Email,
	}
}

//easyjson:json
type createAlreadyExistsResponse []models.User

func newCreateAlreadyExistsResponse(users []models.User) createAlreadyExistsResponse {
	return users
}

type getResponse struct {
	ID       int
	Nickname string
	Fullname string
	About    string
	Email    string
}

func newGetResponse(user *models.User) *getResponse {
	return &getResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		About:    user.About,
		Email:    user.Email,
	}
}

type updateResponse struct {
	ID       int
	Nickname string
	Fullname string
	About    string
	Email    string
}

func newUpdateResponse(user *models.User) *updateResponse {
	return &updateResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		About:    user.About,
		Email:    user.Email,
	}
}
