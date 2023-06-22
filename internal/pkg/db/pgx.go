package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
)

func NewPgxPool(log *zap.Logger) (*pgxpool.Pool, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	conf, _ := pgxpool.ParseConfig("postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?" + "pool_max_conns=100")
	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		log.Error("Failed to connect to db ", zap.Error(err))
		return nil, err
	}

	return pool, nil
}
