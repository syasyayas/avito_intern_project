package saver

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"os"
)

// решил что оно того не стоит
// DEPRECATED
type GDriveSaver struct {
	drive *drive.Service

	log *logrus.Logger
}

func NewGDriveSaver(APIKey string, log *logrus.Logger) (*GDriveSaver, error) {

	if APIKey == "" {
		return nil, fmt.Errorf("no google drive credentials provided")
	}

	driveService, err := drive.NewService(context.Background(), option.WithAPIKey(APIKey), option.WithScopes(drive.DriveScope))
	if err != nil {
		return nil, err
	}

	return &GDriveSaver{
		drive: driveService,
		log:   log,
	}, nil
}

func (s *GDriveSaver) Save(ctx context.Context, data [][]string) (string, error) {
	fileNameUUID := uuid.New().String()
	f, err := os.Create(fileNameUUID + ".csv")

	defer func(f *os.File) {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}(f)

	if err != nil {
		return "", fmt.Errorf("couldn't create tmp file %w", err)
	}
	w := csv.NewWriter(f)
	if err := w.WriteAll(data); err != nil {
		_ = f.Close()
		return "", fmt.Errorf("couldn't write data to file: %w", err)
	}
	w.Flush()

	file := &drive.File{Name: fileNameUUID, MimeType: "text/csv"}
	permissions := &drive.Permission{Type: "anyone", Role: "reader"}

	savedFile, err := s.drive.Files.Create(file).Media(f, googleapi.ContentType(file.MimeType)).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("couldn't save file to gdrive: %w", err)
	}
	_, err = s.drive.Permissions.Create(savedFile.Id, permissions).Do()
	if err != nil {
		return "", fmt.Errorf("couldn't grant permissions to file: %w", err)
	}

	s.log.Infof("Created file on gdrive: id: %s name: %s", savedFile.Id, savedFile.Name)

	return fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", savedFile.Id), nil
}
