package service

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

type FeatureService struct {
	repo *repository.RepoContainer

	log *logrus.Logger
}

func NewFeatureService(r *repository.RepoContainer, l *logrus.Logger) *FeatureService {
	return &FeatureService{
		repo: r,
		log:  l,
	}
}

func (s *FeatureService) AddFeature(ctx context.Context, feature *model.Feature) (*model.Feature, error) {
	if feature == nil || feature.Slug == "" {
		return nil, ErrFeatureEmpty
	}
	var err error
	feature.ID, err = s.repo.AddFeature(ctx, feature.Slug)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, ErrFeatureAlreadyExists
		}
		return nil, err
	}
	return feature, nil
}

func (s *FeatureService) DeleteFeature(ctx context.Context, feature *model.Feature) error {
	if feature == nil || feature.Slug == "" {
		return ErrFeatureEmpty
	}
	if err := s.repo.DeleteFeature(ctx, feature.Slug); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrFeatureNotFound
		}
		return err
	}
	return nil
}

func (s *FeatureService) AddFeaturesToUser(ctx context.Context, user *model.User, features []model.Feature) error {
	_, err := s.repo.GetUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	err = s.repo.AddFeaturesToUser(ctx, user.ID, features)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrFeatureNotFound
		}
		if errors.Is(err, repository.ErrInvalidExpiresAt) {
			return ErrFeatureInvalid
		}
	}
	return nil
}

func (s *FeatureService) DeleteFeatureFromUser(ctx context.Context, feature model.Feature, user model.User) error {
	return nil
}

func (s *FeatureService) DeleteFeaturesFromUser(ctx context.Context, features []model.Feature, user model.User) error {
	_, err := s.repo.GetUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	err = s.repo.DeleteFeaturesFromUser(ctx, user.ID, features)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrFeatureNotFound
		}
	}
	return nil
}
