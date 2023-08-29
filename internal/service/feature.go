package service

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"avito_project/internal/repository/repoerr"
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

func (s *FeatureService) AddFeature(ctx context.Context, feature *model.Feature) error {
	if feature == nil || feature.Slug == "" {
		return ErrFeatureEmpty
	}
	err := s.repo.AddFeature(ctx, feature.Slug)
	if err != nil {
		if errors.Is(err, repoerr.ErrAlreadyExists) {
			return ErrFeatureAlreadyExists
		}
		return err
	}
	return nil
}

func (s *FeatureService) AddFeaturePercent(ctx context.Context, feature *model.Feature, percent float64) error {
	if feature == nil || feature.Slug == "" {
		return ErrFeatureEmpty
	}
	return nil
}

func (s *FeatureService) DeleteFeature(ctx context.Context, feature *model.Feature) error {
	if feature == nil || feature.Slug == "" {
		return ErrFeatureEmpty
	}
	if err := s.repo.DeleteFeature(ctx, feature.Slug); err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			return ErrFeatureNotFound
		}
		return err
	}
	return nil
}

func (s *FeatureService) AddFeaturesToUser(ctx context.Context, user *model.User, features []model.Feature) error {
	_, err := s.repo.GetUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	err = s.repo.AddFeaturesToUser(ctx, user.ID, features)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			return ErrFeatureNotFound
		}
		if errors.Is(err, repoerr.ErrInvalidExpiresAt) {
			return ErrFeatureInvalid
		}
	}
	return nil
}

func (s *FeatureService) DeleteFeaturesFromUser(ctx context.Context, features []model.Feature, user *model.User) error {
	_, err := s.repo.GetUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	for _, feature := range features {
		err = s.repo.DeleteFeatureFromUser(ctx, user.ID, feature.Slug)
		if err != nil {
			if errors.Is(err, repoerr.ErrNotFound) {
				continue
			}
			return err
		}
	}
	return nil
}
