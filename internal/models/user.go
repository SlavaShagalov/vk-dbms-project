package models

//go:generate easyjson -all -snake_case user.go

type User struct {
	ID       int
	Nickname string
	Fullname string
	About    string
	Email    string
}
