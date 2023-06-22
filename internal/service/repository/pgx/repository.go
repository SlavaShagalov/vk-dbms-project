package pgx

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"

	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgService "github.com/SlavaShagalov/vk-dbms-project/internal/service"
)

type repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) pkgService.Repository {
	return &repository{pool: pool, log: log}
}

var getStatusCmd = `
SELECT
(SELECT count(nickname) FROM users)  AS users,
(SELECT count(slug) FROM forums) AS forums,
(SELECT count(id) FROM threads) AS threads,
(SELECT count(id) FROM posts)  AS posts;`

func (rep *repository) GetStatus() (models.Status, error) {
	tmp := models.Status{}
	row := rep.pool.QueryRow(context.Background(), getStatusCmd)

	if err := row.Scan(&tmp.User, &tmp.Forum, &tmp.Thread, &tmp.Post); err != nil {
		rep.log.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.ErrInternal
	}
	return tmp, nil
}

var clearDbCmd = `
TRUNCATE TABLE
users,
forums,
threads,
posts,
votes
CASCADE;`

func (rep *repository) Clear() error {
	_, err := rep.pool.Exec(context.Background(), clearDbCmd)
	if err != nil {
		rep.log.Error("DB error", zap.Error(err))
		return pkgErrors.ErrInternal
	}
	return nil
}
