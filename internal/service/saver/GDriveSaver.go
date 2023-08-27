package saver

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GDriveSaver struct {
	drive *drive.Service

	log *logrus.Logger
}

func NewGDriveSaver(APIKey string, log *logrus.Logger) (*GDriveSaver, error) {
	if APIKey == "" {
		return nil, fmt.Errorf("no google drive credentials provided")
	}
	APIKey = "AIzaSyDnkbxg1YSOVJoMXxAZHsG4FbITboJ0FoM" // TODO FIX

	driveService, err := drive.NewService(context.Background(), option.WithAPIKey(APIKey))
	if err != nil {
		return nil, err
	}

	return &GDriveSaver{
		drive: driveService,
		log:   log,
	}, nil
}

func (s *GDriveSaver) Save(ctx context.Context, data [][]string) (string, error) {
	return "", nil
}
