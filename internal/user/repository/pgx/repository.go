package pgx

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"

	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	"github.com/SlavaShagalov/vk-dbms-project/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgUser "github.com/SlavaShagalov/vk-dbms-project/internal/user"
)

type repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, logger *zap.Logger) pkgUser.Repository {
	return &repository{pool: pool, log: logger}
}

const createCmd = `
INSERT INTO users (nickname, fullname, about, email)
VALUES ($1, $2, $3, $4)
RETURNING id, nickname, fullname, about, email;`

func (rep *repository) Create(ctx context.Context, params *pkgUser.CreateParams) ([]models.User, error) {
	row := rep.pool.QueryRow(ctx, createCmd, params.Nickname, params.Fullname, params.About, params.Email)

	user := new(models.User)
	if err := row.Scan(&user.ID, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_email_key" || pgErr.ConstraintName == "users_nickname_key" {
				users := make([]models.User, 0, 2)

				if user, err = rep.GetByNickname(ctx, params.Nickname); err == nil {
					users = append(users, *user)
				}

				if user, err = rep.GetByEmail(ctx, params.Email); err == nil {
					if len(users) > 0 {
						if user.ID != users[0].ID {
							users = append(users, *user)
						}
					} else {
						users = append(users, *user)
					}
				}
				return users, pkgErrors.ErrUserAlreadyExists
			} else {
				rep.log.Error(constants.DBError, zap.Error(err), zap.String("cmd", createCmd),
					zap.Any("params", params))
				return nil, pkgErrors.ErrInternal
			}
		}
	}

	return []models.User{*user}, nil
}

const getByNicknameCmd = `
SELECT id, nickname, fullname, about, email
FROM users
WHERE nickname = $1;`

func (rep *repository) GetByNickname(ctx context.Context, nickname string) (*models.User, error) {
	row := rep.pool.QueryRow(ctx, getByNicknameCmd, nickname)

	user := new(models.User)
	if err := row.Scan(&user.ID, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, pkgErrors.ErrUserNotFound
		} else {
			rep.log.Error(constants.DBError, zap.Error(err), zap.String("cmd", getByNicknameCmd),
				zap.String("nickname", nickname))
			return user, pkgErrors.ErrInternal
		}
	}

	return user, nil
}

const getByEmailCmd = `
SELECT id, nickname, fullname, about, email
FROM users
WHERE email = $1;`

func (rep *repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := rep.pool.QueryRow(ctx, getByEmailCmd, email)

	user := new(models.User)
	if err := row.Scan(&user.ID, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, pkgErrors.ErrUserNotFound
		} else {
			rep.log.Error(constants.DBError, zap.Error(err), zap.String("cmd", getByEmailCmd),
				zap.String("email", email))
			return user, pkgErrors.ErrInternal
		}
	}

	return user, nil
}

const updateCmd = `
UPDATE users
SET fullname = $1,
    about    = $2,
    email    = $3
WHERE nickname = $4
RETURNING id, nickname, fullname, about, email;`

func (rep *repository) Update(ctx context.Context, params *pkgUser.UpdateParams) (*models.User, error) {
	row := rep.pool.QueryRow(ctx, updateCmd, params.Nickname, params.Fullname, params.About, params.Email)

	user := new(models.User)
	if err := row.Scan(&user.ID, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgErrors.ErrUserNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_email_key" {
				user, _ = rep.GetByNickname(ctx, params.Nickname)
				return user, pkgErrors.ErrUserAlreadyExists
			} else {
				rep.log.Error(constants.DBError, zap.Error(err), zap.String("cmd", updateCmd),
					zap.Any("params", params))
				return user, pkgErrors.ErrInternal
			}
		}
	}

	return user, nil
}
