package pgx

import (
	"context"
	"errors"
	"github.com/SlavaShagalov/vk-dbms-project/internal/pkg/constants"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"

	pkgForum "github.com/SlavaShagalov/vk-dbms-project/internal/forum"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
)

type repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) pkgForum.Repository {
	return &repository{pool: pool, log: log}
}

const getForumUserCmd = `
SELECT nickname
FROM users
WHERE nickname = $1;`

const createCmd = `
INSERT INTO forums (title, user_nickname, slug)
VALUES($1, $2, $3)
RETURNING id, title, user_nickname, slug, posts, threads;`

func (rep *repository) Create(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	userRow := rep.pool.QueryRow(ctx, getForumUserCmd, forum.User)

	if err := userRow.Scan(&forum.User); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, pkgErrors.ErrInternal
	}

	row := rep.pool.QueryRow(ctx, createCmd, forum.Title, forum.User, forum.Slug)
	if err := row.Scan(&forum.ID, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "forums_pkey" {
				tmp, _ := rep.Get(ctx, forum.Slug)
				return tmp, pkgErrors.ErrForumAlreadyExists
			}

			if pgErr.ConstraintName == "forums_user_nickname_fkey" {
				return forum, pkgErrors.ErrUserNotFound
			}

			rep.log.Error(constants.DBError, zap.Error(err))
			return forum, pkgErrors.ErrInternal
		}
	}
	return forum, nil
}

const getCmd = `
SELECT id, title, user_nickname, slug, posts, threads
FROM forums
WHERE slug = $1;`

func (rep *repository) Get(ctx context.Context, slug string) (*models.Forum, error) {
	row := rep.pool.QueryRow(ctx, getCmd, slug)
	forum := new(models.Forum)
	if err := row.Scan(&forum.ID, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, pkgErrors.ErrForumNotFound
		} else {
			rep.log.Error(constants.DBError, zap.Error(err))
			return nil, pkgErrors.ErrInternal
		}
	}
	return forum, nil
}

const createThreadCmd = `
INSERT INTO threads (title, author, forum, message, slug, created)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING  id, title, author, (SELECT slug from forums WHERE slug = $3), message, slug, created;`

func (rep *repository) CreateThread(thread *models.Thread) (models.Thread, error) {
	if thread.Slug != "" {
		tmp, err := rep.GetThread(thread.Slug)
		if err == nil {
			return tmp, pkgErrors.ErrThreadAlreadyExists
		}
	}

	row := rep.pool.QueryRow(
		context.Background(),
		createThreadCmd,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Created,
	)
	tmp := models.Thread{}
	err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Created)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "threads_forum_fkey":
				return models.Thread{}, pkgErrors.ErrForumNotFound
			case "threads_author_fkey":
				return models.Thread{}, pkgErrors.ErrUserNotFound
			default:
				return tmp, pkgErrors.ErrInternal
			}
		}
	}

	return tmp, nil
}

const getThreadBySlugCmd = `
SELECT  id, title, author, forum, message, slug, votes, created
FROM threads
WHERE slug = $1;`

const getThreadByIdCmd = `
SELECT  id, title, author, forum, message, slug, votes, created
FROM threads
WHERE id = $1;`

func (rep *repository) GetThread(slugOrId string) (models.Thread, error) {
	tmp := models.Thread{}
	var row pgx.Row

	if id, err := strconv.Atoi(slugOrId); err == nil {
		row = rep.pool.QueryRow(context.Background(), getThreadByIdCmd, id)
	} else {
		row = rep.pool.QueryRow(context.Background(), getThreadBySlugCmd, slugOrId)
	}

	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Votes, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrThreadNotFound
		}
		rep.log.Error(constants.DBError, zap.Error(err))
		return tmp, pkgErrors.ErrInternal
	}

	return tmp, nil
}

const getForumUsersAsc = `
SELECT nickname, fullname, about, email
FROM forum_users
WHERE forum = $1
ORDER BY nickname
LIMIT $2;`

const getForumUsersWithSinceAsc = `
SELECT nickname, fullname, about, email
FROM forum_users
WHERE forum = $1 AND nickname > $2
ORDER BY nickname
LIMIT $3;`

const getForumUsersDesc = `
SELECT nickname, fullname, about, email
FROM forum_users
WHERE forum = $1
ORDER BY nickname DESC
LIMIT $2;`

const getForumUsersWithSinceDesc = `
SELECT nickname, fullname, about, email
FROM forum_users
WHERE forum = $1 AND nickname < $2
ORDER BY nickname DESC
LIMIT $3;`

func (rep *repository) GetForumUsers(ctx context.Context, slug string, limit int, since string, desc bool) ([]models.User, error) {
	var rows pgx.Rows
	var err error
	users := make([]models.User, 0)
	if _, err := rep.Get(ctx, slug); err != nil {
		return []models.User{}, err
	}

	if desc {
		if since != "" {
			rows, err = rep.pool.Query(ctx, getForumUsersWithSinceDesc, slug, since, limit)
		} else {
			rows, err = rep.pool.Query(ctx, getForumUsersDesc, slug, limit)
		}
	} else {
		if since != "" {
			rows, err = rep.pool.Query(ctx, getForumUsersWithSinceAsc, slug, since, limit)
		} else {
			rows, err = rep.pool.Query(ctx, getForumUsersAsc, slug, limit)
		}
	}
	defer rows.Close()

	if err != nil {
		rep.log.Error(constants.DBError, zap.Error(err))
		return []models.User{}, pkgErrors.ErrInternal
	}

	tmp := models.User{}
	for rows.Next() {
		if err := rows.Scan(&tmp.Nickname, &tmp.Fullname, &tmp.About, &tmp.Email); err != nil {
			rep.log.Error(constants.DBError, zap.Error(err))
			return []models.User{}, pkgErrors.ErrInternal
		}
		users = append(users, tmp)
	}
	return users, nil
}

const getThreadsDescCmd = `
SELECT id, title, author, forum, message, slug, votes, created
FROM threads
WHERE forum = $1
ORDER BY created DESC
LIMIT $2;`

const getThreadsDescWithFilterCmd = `
SELECT id, title, author, forum, message, slug, votes, created
FROM threads
WHERE forum = $1 AND created <= $3
ORDER BY created DESC
LIMIT $2;`

const getThreadsAscCmd = `
SELECT id, title, author, forum, message, slug, votes, created
FROM threads
WHERE forum = $1
ORDER BY created
LIMIT $2;`

const getThreadsAscWithFilterCmd = `
SELECT id, title, author, forum, message, slug, votes, created
FROM threads
WHERE forum = $1 AND created >= $3
ORDER BY created
LIMIT $2;`

func (rep *repository) GetForumThreads(ctx context.Context, slug string, limit int, since string,
	desc bool) (models.ThreadList, error) {
	getCmd := ""
	var rows pgx.Rows
	var err error

	if _, err := rep.Get(ctx, slug); err != nil {
		if errors.Is(pkgErrors.ErrForumNotFound, err) {
			return nil, pkgErrors.ErrForumNotFound
		} else {
			return nil, pkgErrors.ErrInternal
		}
	}

	if desc {
		if since == "" {
			getCmd = getThreadsDescCmd
		} else {
			getCmd = getThreadsDescWithFilterCmd
		}
	} else {
		if since == "" {
			getCmd = getThreadsAscCmd
		} else {
			getCmd = getThreadsAscWithFilterCmd
		}
	}

	if since == "" {
		rows, err = rep.pool.Query(ctx, getCmd, slug, limit)
	} else {
		rows, err = rep.pool.Query(ctx, getCmd, slug, limit, since)
	}

	if err != nil {
		rep.log.Error(constants.DBError, zap.Error(err))
	}

	threads := make([]models.Thread, 0)
	tmp := models.Thread{}

	for rows.Next() {

		if err := rows.Scan(
			&tmp.Id,
			&tmp.Title,
			&tmp.Author,
			&tmp.Forum,
			&tmp.Message,
			&tmp.Slug,
			&tmp.Votes,
			&tmp.Created,
		); err != nil {
			rep.log.Error(constants.DBError, zap.Error(err))
			return threads, pkgErrors.ErrInternal
		}

		threads = append(threads, tmp)
	}

	return threads, nil
}
