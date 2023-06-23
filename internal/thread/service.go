package thread

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
)

type Service interface {
	CreatePosts(slugOrId string, posts []models.Post) (models.PostList, error)
	GetThread(slugOrId string) (models.Thread, error)
	UpdateThread(slugOrId string, thread *models.Thread) (models.Thread, error)
	GetPosts(slugOrId string, limit, since int, sort string, desc bool) (models.PostList, error)
	AddVote(slugOrId string, vote *models.Vote) (models.Thread, error)
}
