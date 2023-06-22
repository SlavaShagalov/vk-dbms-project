package service

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgPost "github.com/SlavaShagalov/vk-dbms-project/internal/post"
	"go.uber.org/zap"
)

type service struct {
	rep pkgPost.Repository
	log *zap.Logger
}

func NewService(rep pkgPost.Repository, log *zap.Logger) pkgPost.Service {
	return &service{rep: rep, log: log}
}

func (serv *service) GetPost(id int, related []string) (models.FullPost, error) {
	post, err := serv.rep.GetPost(id)
	if err != nil {
		return models.FullPost{}, err
	}
	tmp := models.FullPost{
		Post:   &post,
		Author: nil,
		Thread: nil,
		Forum:  nil,
	}

	for _, tag := range related {
		switch tag {
		case "user":
			user, err := serv.rep.GetPostAuthor(tmp.Post)
			if err != nil {
				return tmp, err
			}
			tmp.Author = &user

		case "thread":
			thread, err := serv.rep.GetPostThread(tmp.Post)
			if err != nil {
				return tmp, err
			}
			tmp.Thread = &thread

		case "forum":
			forum, err := serv.rep.GetPostForum(tmp.Post)
			if err != nil {
				return tmp, err
			}
			tmp.Forum = &forum
		}
	}
	return tmp, nil
}

func (serv *service) UpdatePost(post *models.Post) (models.Post, error) {
	return serv.rep.UpdatePost(post)
}
