package service

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"avito_project/internal/repository/repoerr"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	Repo *repository.RepoContainer

	log *logrus.Logger
}

func NewUserService(repo *repository.RepoContainer, log *logrus.Logger) *UserService {
	return &UserService{
		Repo: repo,
		log:  log,
	}
}

func (s *UserService) AddUser(ctx context.Context, user *model.User) error {
	err := s.Repo.AddUser(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to add user %v: %v", user, err)
		if errors.Is(err, repoerr.ErrAlreadyExists) {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, user *model.User) error {
	err := s.Repo.DeleteUser(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to delete user %s: %v", user.ID, err)
		if errors.Is(err, repoerr.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *UserService) GetUserWithFeatures(ctx context.Context, user *model.User) (*model.User, error) {

	if _, err := s.Repo.GetUser(ctx, user.ID); err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			return nil, ErrUserNotFound
		}
	}
	res, err := s.Repo.GetUserWithFeatures(ctx, user.ID)
	if err != nil {
		s.log.Errorf("Failed to retrieve user %s with features: %v", user.ID, err)
		if errors.Is(err, repoerr.ErrNotFound) {
			return nil, ErrUserNotFound
		}
	}
	return res, nil
}
