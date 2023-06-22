package models

//go:generate easyjson -all -snake_case post.go

import "time"

//easyjson:json
type PostList []Post

type Post struct {
	Id       int       `json:"id"`
	Parent   int       `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created"`
}

type FullPost struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}
