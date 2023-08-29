package saver

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
)

type LocalSaver struct {
	appPort int
	log     *logrus.Logger
}

func NewLocalSaver(log *logrus.Logger, port int) *LocalSaver {
	return &LocalSaver{log: log, appPort: port}
}

func (s *LocalSaver) Save(ctx context.Context, data [][]string) (string, error) {
	fileNameUUID := uuid.New().String()
	f, err := os.Create(fileNameUUID + ".csv")

	defer func(f *os.File) {
		_ = f.Close()
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
	return fmt.Sprintf("http://localhost:%d/static/%s", s.appPort, f.Name()), nil
}
