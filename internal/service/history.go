package service

import (
	"avito_project/internal/model"
	"avito_project/internal/repository"
	"avito_project/internal/repository/repoerr"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Saver interface {
	Save(ctx context.Context, data [][]string) (string, error)
}

type HistoryService struct {
	repos *repository.RepoContainer
	saver Saver

	log *logrus.Logger
}

func NewHistoryService(repos *repository.RepoContainer, saver Saver, log *logrus.Logger) *HistoryService {
	return &HistoryService{repos: repos, saver: saver, log: log}
}

func (s *HistoryService) GetHistory(ctx context.Context, after time.Time, before time.Time) (model.History, error) {
	hist, err := s.repos.GetHistory(ctx, after, before)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			return nil, ErrHistoryEmpty
		}
		return nil, err
	}
	return hist, nil

}

func (s *HistoryService) Export(ctx context.Context, after time.Time, before time.Time) (string, error) {
	data, err := s.GetHistory(ctx, after, before)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	return s.saver.Save(ctx, data.ParseToCSV())
}
