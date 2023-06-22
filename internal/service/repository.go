package post

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
)

type Repository interface {
	GetStatus() (models.Status, error)
	Clear() error
}
