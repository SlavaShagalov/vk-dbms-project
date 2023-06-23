package models

//go:generate easyjson -all -snake_case thread.go

import "time"

//easyjson:json
type ThreadList []Thread

type Thread struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}
