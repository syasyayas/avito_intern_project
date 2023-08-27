package postgres

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserRepo struct {
	db *pgxpool.Pool

	log *logrus.Logger
}

func NewUserRepo(db *pgxpool.Pool, log *logrus.Logger) *UserRepo {
	return &UserRepo{db, log}
}

func (r *UserRepo) AddUser(ctx context.Context, id string) error {
	r.log.Debugf("Ading user %s", id)

	_, err := r.db.Exec(ctx, "INSERT INTO avito_features.users (id) VALUES ($1)", id)

	return repository.PgErrorWrapper(err)
}

func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	r.log.Debugf("Deleting user %s", id)

	_, err := r.db.Exec(ctx, "UPDATE users SET deleted_at = NOW() WHERE id = ? RETURNING id", id)

	return repository.PgErrorWrapper(err)
}

// TODO delete?
func (r *UserRepo) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user = &model.User{}
	row := r.db.QueryRow(ctx, "SELECT id FROM avito_features.users WHERE id = ? AND deleted_at IS NOT NULL", id)

	err := row.Scan(user)
	if err != nil {
		r.log.Errorf("couldn't get user %s. Error: %v", id, err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}
