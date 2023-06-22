package models

//go:generate easyjson -all -snake_case vote.go

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}
