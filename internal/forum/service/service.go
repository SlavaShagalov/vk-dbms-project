package service

import (
	"context"
	"github.com/SlavaShagalov/vk-dbms-project/internal/forum"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	"go.uber.org/zap"
)

type service struct {
	rep forum.Repository
	log *zap.Logger
}

func NewService(rep forum.Repository, log *zap.Logger) forum.Service {
	return &service{rep: rep, log: log}
}

func (serv *service) Create(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	return serv.rep.Create(ctx, forum)
}

func (serv *service) CreateThread(thread *models.Thread) (models.Thread, error) {
	return serv.rep.CreateThread(thread)
}

func (serv *service) Get(ctx context.Context, slug string) (*models.Forum, error) {
	return serv.rep.Get(ctx, slug)
}

func (serv *service) GetForumUsers(ctx context.Context, slug string, limit int, since string,
	desc bool) (models.UserList, error) {
	return serv.rep.GetForumUsers(ctx, slug, limit, since, desc)
}

func (serv *service) GetForumThreads(ctx context.Context, slug string, limit int, since string,
	desc bool) (models.ThreadList, error) {
	return serv.rep.GetForumThreads(ctx, slug, limit, since, desc)
}
