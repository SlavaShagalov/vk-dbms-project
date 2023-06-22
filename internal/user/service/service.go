package service

import (
	"context"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	"github.com/SlavaShagalov/vk-dbms-project/internal/user"
	"go.uber.org/zap"

	pkgUser "github.com/SlavaShagalov/vk-dbms-project/internal/user"
)

type service struct {
	rep user.Repository
	log *zap.Logger
}

func NewService(rep user.Repository, log *zap.Logger) user.Service {
	return &service{rep: rep, log: log}
}

func (serv *service) Create(ctx context.Context, params *pkgUser.CreateParams) ([]models.User, error) {
	return serv.rep.Create(ctx, params)
}

func (serv *service) GetByNickname(ctx context.Context, nickname string) (*models.User, error) {
	return serv.rep.GetByNickname(ctx, nickname)
}

func (serv *service) Update(ctx context.Context, params *pkgUser.UpdateParams) (*models.User, error) {
	return serv.rep.Update(ctx, params)
}
