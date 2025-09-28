package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mibrgmv/document-service/internal/config"
)

func NewPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := "postgres://" + cfg.Postgres.User + ":" + cfg.Postgres.Password +
		"@" + cfg.Postgres.Host + ":" + cfg.Postgres.Port + "/" + cfg.Postgres.DBName +
		"?sslmode=" + cfg.Postgres.SSLMode

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("connected to postgres")
	return pool, nil
}
