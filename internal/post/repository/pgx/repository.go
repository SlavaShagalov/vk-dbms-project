package pgx

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"

	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgPost "github.com/SlavaShagalov/vk-dbms-project/internal/post"
)

type repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) pkgPost.Repository {
	return &repository{pool: pool, log: log}
}

const getPostById = `
SELECT  id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE id = $1;`

func (rep *repository) GetPost(id int) (models.Post, error) {
	tmp := models.Post{}
	row := rep.pool.QueryRow(context.Background(), getPostById, id)
	if err := row.Scan(&tmp.Id, &tmp.Parent, &tmp.Author, &tmp.Message, &tmp.IsEdited, &tmp.Forum, &tmp.Thread, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrPostNotFound
		}
		return tmp, pkgErrors.ErrInternal
	}

	return tmp, nil
}

const getPostAuthor = `
SELECT id, nickname, fullname, about, email
FROM users
WHERE nickname = $1;`

func (rep *repository) GetPostAuthor(post *models.Post) (models.User, error) {
	tmp := models.User{}
	row := rep.pool.QueryRow(context.Background(), getPostAuthor, post.Author)
	if err := row.Scan(&tmp.ID, &tmp.Nickname, &tmp.Fullname, &tmp.About, &tmp.Email); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrUserNotFound
		}
		return tmp, pkgErrors.ErrInternal
	}
	return tmp, nil
}

const getPostForum = `
SELECT id, title, user_nickname, slug, posts, threads
FROM forums
WHERE slug = $1;`

func (rep *repository) GetPostForum(post *models.Post) (models.Forum, error) {
	tmp := models.Forum{}
	row := rep.pool.QueryRow(context.Background(), getPostForum, post.Forum)
	if err := row.Scan(&tmp.ID, &tmp.Title, &tmp.User, &tmp.Slug, &tmp.Posts, &tmp.Threads); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrUserNotFound
		}
		return tmp, pkgErrors.ErrInternal
	}
	return tmp, nil
}

const getPostThread = `
SELECT  id, title, author, forum, message, slug, votes, created
FROM threads
WHERE id = $1;`

func (rep *repository) GetPostThread(post *models.Post) (models.Thread, error) {
	tmp := models.Thread{}
	row := rep.pool.QueryRow(context.Background(), getPostThread, post.Thread)
	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Votes, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrUserNotFound
		}
		return tmp, pkgErrors.ErrInternal
	}
	return tmp, nil
}

const updatePost = `
UPDATE posts
SET isEdited = case when (trim($2) = '') OR (trim($2) = trim(message)) then false else true end,
	message = case when trim($2) = '' then message else $2 end
WHERE id = $1
RETURNING id, parent, author, message, isEdited, forum, thread, created;`

func (rep *repository) UpdatePost(post *models.Post) (models.Post, error) {
	tmp := models.Post{}

	row := rep.pool.QueryRow(context.Background(), updatePost, post.Id, post.Message)
	if err := row.Scan(&tmp.Id, &tmp.Parent, &tmp.Author, &tmp.Message, &tmp.IsEdited, &tmp.Forum, &tmp.Thread, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrPostNotFound
		}
		rep.log.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.ErrInternal
	}

	return tmp, nil
}
