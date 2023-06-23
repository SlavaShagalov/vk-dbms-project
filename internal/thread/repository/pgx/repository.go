package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/SlavaShagalov/vk-dbms-project/internal/pkg/constants"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"

	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgThread "github.com/SlavaShagalov/vk-dbms-project/internal/thread"
)

type repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, log *zap.Logger) pkgThread.Repository {
	return &repository{pool: pool, log: log}
}

const createPostBeginCmd = `
INSERT INTO posts (parent, author, message, thread, forum, created) 
VALUES `

const checkPostAuthor = `
SELECT id
FROM users
WHERE nickname = $1;`

const checkPostParent = `
SELECT id
FROM posts
WHERE id = $1 AND thread = $2;`

func (rep *repository) CreatePosts(slugOrId string, posts []models.Post) ([]models.Post, error) {
	result := make([]models.Post, 0)

	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return result, err
	}

	if len(posts) == 0 {
		return result, nil
	}

	postTmp := models.Post{}
	created := time.Unix(0, time.Now().UnixNano()/1e6*1e6)
	cmd := createPostBeginCmd
	args := make([]interface{}, 0, 6*len(posts))

	postTmp.Created = created
	for ind, post := range posts {
		tmpId := 0

		row := rep.pool.QueryRow(context.Background(), checkPostAuthor, post.Author)
		if err := row.Scan(&tmpId); err != nil {
			return result, pkgErrors.ErrUserNotFound
		}

		if post.Parent != 0 {
			row = rep.pool.QueryRow(context.Background(), checkPostParent, post.Parent, thread.Id)
			if err := row.Scan(&tmpId); err != nil {
				return result, pkgErrors.ErrParentPostNotFound
			}
		}

		cmd += fmt.Sprintf(" ($%d, $%d, $%d, $%d, $%d, $%d)", 6*ind+1, 6*ind+2, 6*ind+3, 6*ind+4, 6*ind+5, 6*ind+6)
		args = append(args, post.Parent, post.Author, post.Message, thread.Id, thread.Forum, created)
		if ind != len(posts)-1 {
			cmd += ","
		}
	}
	cmd += " RETURNING id, parent, author, message, isEdited, forum, thread;"

	rows, err := rep.pool.Query(context.Background(), cmd, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			rep.log.Error("TEST", zap.Error(err))
			if pgErr.Message == "Invalid parent" {
				return []models.Post{}, pkgErrors.ErrParentPostNotFound
			}
			switch pgErr.ConstraintName {
			case "posts_forum_fkey":
				return []models.Post{}, pkgErrors.ErrForumNotFound
			case "posts_thread_fkey":
				return []models.Post{}, pkgErrors.ErrThreadNotFound
			case "thread_check":
				return []models.Post{}, pkgErrors.ErrThreadNotFound
			case "posts_author_fkey":
				return []models.Post{}, pkgErrors.ErrUserNotFound
			}
		}
		rep.log.Error(constants.DBError, zap.Error(err))
		return []models.Post{}, pkgErrors.ErrInternal
	}

	for rows.Next() {
		if rows.Scan(&postTmp.Id, &postTmp.Parent, &postTmp.Author, &postTmp.Message, &postTmp.IsEdited, &postTmp.Forum, &postTmp.Thread); err != nil {
			rep.log.Error(constants.DBError, zap.Error(err))
			return []models.Post{}, pkgErrors.ErrInternal
		}
		result = append(result, postTmp)
	}

	return result, nil
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

const updateThreadByIdCmd = `
UPDATE threads
SET message = case when trim($2) = '' then message else $2 end, 
	title = case when trim($3) = '' then title else $3 end
WHERE id = $1
RETURNING  id, title, author, forum, message, slug, votes, created;`

const updateThreadBySlugCmd = `
UPDATE threads
SET message = case when trim($2) = '' then message else $2 end, 
	title = case when trim($3) = '' then title else $3 end
WHERE slug = $1
RETURNING  id, title, author, forum, message, slug, votes, created;`

func (rep *repository) UpdateThread(slugOrId string, thread *models.Thread) (models.Thread, error) {
	tmp := models.Thread{}
	var row pgx.Row

	if id, err := strconv.Atoi(slugOrId); err == nil {
		row = rep.pool.QueryRow(context.Background(), updateThreadByIdCmd, id, thread.Message, thread.Title)
	} else {
		row = rep.pool.QueryRow(context.Background(), updateThreadBySlugCmd, slugOrId, thread.Message, thread.Title)
	}

	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Votes, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrThreadNotFound
		}
		rep.log.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.ErrInternal
	}

	return tmp, nil
}

const getPostsAscCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE thread = $1 AND id > $2
ORDER BY created, id
LIMIT $3;`

const getPostsDescWithSinceCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE thread = $1 AND id < $2
ORDER BY created DESC, id DESC 
LIMIT $3;`

const getPostsDescCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE thread = $1
ORDER BY created DESC, id DESC 
LIMIT $2;`

func (rep *repository) GetPostsFlat(slugOrId string, limit, since int, desc bool) (models.PostList, error) {
	var rows pgx.Rows
	var err error

	tmp := make([]models.Post, 0)
	post := models.Post{}
	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return []models.Post{}, err
	}

	if desc {
		if since != 0 {
			rows, err = rep.pool.Query(context.Background(), getPostsDescWithSinceCmd, thread.Id, since, limit)
			if err != nil {
				return tmp, pkgErrors.ErrInternal
			}
		} else {
			rows, err = rep.pool.Query(context.Background(), getPostsDescCmd, thread.Id, limit)
			if err != nil {
				return tmp, pkgErrors.ErrInternal
			}
		}
	} else {
		rows, err = rep.pool.Query(context.Background(), getPostsAscCmd, thread.Id, since, limit)
	}

	if err != nil {
		return tmp, pkgErrors.ErrInternal
	}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return tmp, pkgErrors.ErrInternal
		}
		tmp = append(tmp, post)
	}

	return tmp, nil
}

const getPostsTreeAscCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE thread = $1
ORDER BY path, id
LIMIT $2;`

const getPostsTreeWithSinceAscCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE thread = $1 AND path > (SELECT path FROM posts WHERE id = $2) 
ORDER BY path, id
LIMIT $3;`

const getPostsTreeDescCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE thread = $1
ORDER BY path DESC, id
LIMIT $2;`

const getPostsTreeWithSinceDescCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE thread = $1 AND path < (SELECT path FROM posts WHERE id = $2) 
ORDER BY path DESC, id
LIMIT $3;`

func (rep *repository) GetPostsTree(slugOrId string, limit, since int, desc bool) (models.PostList, error) {
	var rows pgx.Rows
	var err error

	tmp := make([]models.Post, 0)
	post := models.Post{}
	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return []models.Post{}, err
	}

	cmd := ""

	if desc {
		cmd = getPostsTreeDescCmd
	} else {
		cmd = getPostsTreeAscCmd
	}

	if since != 0 {
		switch cmd {
		case getPostsTreeAscCmd:
			cmd = getPostsTreeWithSinceAscCmd
		case getPostsTreeDescCmd:
			cmd = getPostsTreeWithSinceDescCmd
		}
	}

	if since != 0 {
		rows, err = rep.pool.Query(context.Background(), cmd, thread.Id, since, limit)
	} else {
		rows, err = rep.pool.Query(context.Background(), cmd, thread.Id, limit)
	}

	if err != nil {
		return tmp, pkgErrors.ErrInternal
	}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return tmp, pkgErrors.ErrInternal
		}
		tmp = append(tmp, post)
	}

	return tmp, nil
}

const getPostsParentTreeAscCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id ASC LIMIT $2)
ORDER BY path, id;`

const getPostsParentTreeWithSinceAscCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created
FROM posts
WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND path[1] >
(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id ASC LIMIT $3) 
ORDER BY path, id;`

const getPostsParentTreeDescCmd = `
SELECT id, parent, author, message, isEdited, forum, thread,created
FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2)
ORDER BY path[1] DESC, path, id;`

const getPostsParentTreeWithSinceDescCmd = `
SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts
WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND path[1] <
(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3)
ORDER BY path[1] DESC, path, id;`

func (rep *repository) GetPostsParentTree(slugOrId string, limit, since int, desc bool) (models.PostList, error) {
	var rows pgx.Rows
	var err error

	tmp := make([]models.Post, 0)
	post := models.Post{}
	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return []models.Post{}, err
	}

	cmd := ""

	if desc {
		cmd = getPostsParentTreeDescCmd
	} else {
		cmd = getPostsParentTreeAscCmd
	}

	if since != 0 {
		switch cmd {
		case getPostsParentTreeAscCmd:
			cmd = getPostsParentTreeWithSinceAscCmd
		case getPostsParentTreeDescCmd:
			cmd = getPostsParentTreeWithSinceDescCmd
		}
	}

	if since != 0 {
		rows, err = rep.pool.Query(context.Background(), cmd, thread.Id, since, limit)
	} else {
		rows, err = rep.pool.Query(context.Background(), cmd, thread.Id, limit)
	}

	if err != nil {
		rep.log.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.ErrInternal
	}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return tmp, pkgErrors.ErrInternal
		}
		tmp = append(tmp, post)
	}

	return tmp, nil
}

const getVoteCmd = `
SELECT voice
FROM votes
WHERE nickname = $1 AND thread = $2;`

func (rep *repository) GetVote(thread *models.Thread, vote *models.Vote) (models.Vote, error) {
	tmp := *vote
	row := rep.pool.QueryRow(context.Background(), getVoteCmd, vote.Nickname, thread.Id)
	if err := row.Scan(&tmp.Voice); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ErrVoiceNotFound
		}
		return tmp, pkgErrors.ErrInternal
	}
	return tmp, nil
}

const addVoteCmd = `
INSERT INTO votes
(nickname, voice, thread)
VALUES ($1, $2, $3)
RETURNING id;`

func (rep *repository) AddVote(thread *models.Thread, vote *models.Vote) (models.Thread, error) {
	row := rep.pool.QueryRow(context.Background(), addVoteCmd, vote.Nickname, vote.Voice, thread.Id)
	id := 0

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "votes_nickname_thread_key":
				return models.Thread{}, pkgErrors.ErrVoiceAlreadyExists
			case "votes_nickname_fkey":
				return models.Thread{}, pkgErrors.ErrUserNotFound
			default:
				rep.log.Error("DB error", zap.Error(err))
				return models.Thread{}, pkgErrors.ErrInternal
			}
		}
	}

	thread.Votes += vote.Voice
	return *thread, nil
}

const updateVoteCmd = `
UPDATE votes
SET voice = $1
WHERE nickname = $2 AND thread = $3 AND voice != $1
RETURNING id;`

func (rep *repository) UpdateVote(slugOrId string, thread *models.Thread, vote *models.Vote) (models.Thread, error) {
	row := rep.pool.QueryRow(context.Background(), updateVoteCmd, vote.Voice, vote.Nickname, thread.Id)
	id := 0
	if err := row.Scan(&id); err != nil {
		if !errors.Is(pgx.ErrNoRows, err) {
			rep.log.Error("DB error", zap.Error(err))
			return models.Thread{}, pkgErrors.ErrInternal
		}
	}

	return rep.GetThread(slugOrId)
}
