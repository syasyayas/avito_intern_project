package db

import (
	"avito_project/config"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	PgPoolSize     = 30
	PgConnAttempts = 5
	PgConnTimeout  = time.Second
)

func NewPgPool(ctx context.Context, log *logrus.Logger, cfg *config.Config) (*pgxpool.Pool, error) {
	var pgPool *pgxpool.Pool
	var err error

	pgCfg, err := pgxpool.ParseConfig(cfg.Pg.URL)
	pgCfg.MaxConns = PgPoolSize
	log.Infof("Connecting: %s", cfg.Pg.URL)

	for i := PgConnAttempts; i > 0; i-- {
		pgPool, err = pgxpool.New(ctx, cfg.Pg.URL)
		if err == nil {
			err = pgPool.Ping(ctx)
			if err != nil {
				continue
			}
			return pgPool, nil
		}
		log.Infof("Failed to establish posgres connection, attempts left: %d", i-1)
		time.Sleep(PgConnTimeout)
	}

	return nil, err
}
