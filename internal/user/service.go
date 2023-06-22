package user

import (
	"context"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
)

type Service interface {
	Create(ctx context.Context, params *CreateParams) ([]models.User, error)
	GetByNickname(ctx context.Context, nickname string) (*models.User, error)
	Update(ctx context.Context, params *UpdateParams) (*models.User, error)
}
