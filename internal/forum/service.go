package forum

import (
	"context"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
)

type Service interface {
	Create(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	CreateThread(thread *models.Thread) (models.Thread, error)
	GetForumUsers(ctx context.Context, slug string, limit int, since string, desc bool) (models.UserList, error)
	Get(ctx context.Context, slug string) (*models.Forum, error)
	GetForumThreads(ctx context.Context, slug string, limit int, since string, desc bool) (models.ThreadList, error)
}
