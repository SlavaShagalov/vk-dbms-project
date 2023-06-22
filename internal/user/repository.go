package user

import (
	"context"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
)

type CreateParams struct {
	Nickname string
	Fullname string
	About    string
	Email    string
}

type UpdateParams struct {
	Nickname string
	Fullname string
	About    string
	Email    string
}

type Repository interface {
	Create(ctx context.Context, params *CreateParams) ([]models.User, error)
	GetByNickname(ctx context.Context, nickname string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, params *UpdateParams) (*models.User, error)
}
