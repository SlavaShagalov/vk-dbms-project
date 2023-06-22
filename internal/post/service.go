package post

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
)

type Service interface {
	GetPost(id int, related []string) (models.FullPost, error)
	UpdatePost(post *models.Post) (models.Post, error)
}
