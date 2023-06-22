package service

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgService "github.com/SlavaShagalov/vk-dbms-project/internal/service"
	"go.uber.org/zap"
)

type service struct {
	rep pkgService.Repository
	log *zap.Logger
}

func NewService(rep pkgService.Repository, log *zap.Logger) pkgService.Service {
	return &service{rep: rep, log: log}
}
func (serv *service) GetStatus() (models.Status, error) {
	return serv.rep.GetStatus()
}

func (serv *service) Clear() error {
	return serv.rep.Clear()
}
