package service

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Feature interface {
	AddFeature(ctx context.Context, feature *model.Feature) (*model.Feature, error)
	DeleteFeature(ctx context.Context, feature *model.Feature) error
	AddFeaturesToUser(ctx context.Context, user *model.User, features []model.Feature) error
	DeleteFeatureFromUser(ctx context.Context, feature model.Feature, user model.User) error
	DeleteFeaturesFromUser(ctx context.Context, features []model.Feature, user model.User) error
}
type User interface {
	AddUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
	GetUserWithFeatures(ctx context.Context, user *model.User) (*model.User, error)
}

type History interface {
	GetHistory(ctx context.Context, after time.Time, before time.Time) (model.History, error)
	Export(ctx context.Context, after time.Time, before time.Time) (string, error)
}

type Services struct {
	Feature
	User
	History
}

func NewServices(repos *repository.RepoContainer, saver Saver, logger *logrus.Logger) *Services {
	return &Services{
		User:    NewUserService(repos, logger),
		Feature: NewFeatureService(repos, logger),
		History: NewHistoryService(repos, saver, logger),
	}
}
