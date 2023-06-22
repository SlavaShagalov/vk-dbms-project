package models

//go:generate easyjson -all -snake_case forum.go

type Forum struct {
	ID      int64
	Title   string
	User    string
	Slug    string
	Posts   int64
	Threads int64
}
