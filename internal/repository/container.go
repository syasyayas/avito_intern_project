package repository

import (
	"avito_project/internal/model"
	"avito_project/internal/repository/postgres/postgres"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type FeatureRepo interface {
	AddFeature(ctx context.Context, slug string) (int, error)
	DeleteFeature(ctx context.Context, slug string) error
	AddFeaturesToUser(ctx context.Context, userId string, features []model.Feature) error
	DeleteFeaturesFromUser(ctx context.Context, userId string, features []model.Feature) error
	GetUserWithFeatures(ctx context.Context, id string) (*model.User, error)
	DeleteFeatureFromUser(ctx context.Context, userID string, featureID int) error
}

type UserRepo interface {
	AddUser(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error
	GetUser(ctx context.Context, id string) (*model.User, error)
}

type HistoryRepo interface {
	GetHistory(ctx context.Context, after time.Time, before time.Time) (model.History, error)
}

type RepoContainer struct {
	FeatureRepo
	UserRepo
	HistoryRepo
}

func NewPgRepos(db *pgxpool.Pool, log *logrus.Logger) *RepoContainer {
	return &RepoContainer{
		FeatureRepo: postgres.NewFeatureRepo(db, log),
		UserRepo:    postgres.NewUserRepo(db, log),
		HistoryRepo: postgres.NewHistoryRepo(db, log),
	}
}
