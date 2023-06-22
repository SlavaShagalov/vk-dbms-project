package models

//go:generate easyjson -all -snake_case user.go

//easyjson:json
type UserList []User

type User struct {
	ID       int
	Nickname string
	Fullname string
	About    string
	Email    string
}
