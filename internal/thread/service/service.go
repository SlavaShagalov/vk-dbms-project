package service

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	"github.com/SlavaShagalov/vk-dbms-project/internal/thread"
	"go.uber.org/zap"
)

type service struct {
	rep thread.Repository
	log *zap.Logger
}

func NewService(rep thread.Repository, log *zap.Logger) thread.Service {
	return &service{rep: rep, log: log}
}

func (serv *service) CreatePosts(slugOrId string, posts []models.Post) (models.PostList, error) {
	return serv.rep.CreatePosts(slugOrId, posts)
}

func (serv *service) GetThread(slugOrId string) (models.Thread, error) {
	return serv.rep.GetThread(slugOrId)
}

func (serv *service) UpdateThread(slugOrId string, thread *models.Thread) (models.Thread, error) {
	return serv.rep.UpdateThread(slugOrId, thread)
}

func (serv *service) GetPosts(slugOrId string, limit, since int, sort string, desc bool) (models.PostList, error) {
	switch sort {
	case "tree":
		return serv.rep.GetPostsTree(slugOrId, limit, since, desc)
	case "parent_tree":
		return serv.rep.GetPostsParentTree(slugOrId, limit, since, desc)
	default:
		return serv.rep.GetPostsFlat(slugOrId, limit, since, desc)
	}
}

func (serv *service) AddVote(slugOrId string, vote *models.Vote) (models.Thread, error) {
	thread, err := serv.rep.GetThread(slugOrId)
	if err != nil {
		return thread, err
	}

	if _, err := serv.rep.GetVote(&thread, vote); err == nil {
		return serv.rep.UpdateVote(slugOrId, &thread, vote)
	} else {
		return serv.rep.AddVote(&thread, vote)
	}
}
