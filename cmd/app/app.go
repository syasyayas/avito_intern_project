package app

import (
	"avito_project/config"
	"avito_project/internal/repository"
	"avito_project/internal/repository/postgres/postgres/db"
	"avito_project/internal/service"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

const logDefaultLevel = "debug"

func Run(cfgPath string) error {
	cfg, err := config.New(cfgPath)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %v", err)
	}

	log, err := newLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to init logger: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Info("Setting up postgres connection")

	pool, err := db.NewPgPool(ctx, log, cfg)
	if err != nil {
		return fmt.Errorf("Failed to establish postgres connection: %v", err)
	}

	log.Info("Setting up repositories")
	repos := repository.NewPgRepos(pool, log)

	log.Info("Setting up services")

	services := service.NewServices(repos, log)

	log.Info("Setting up handlers")

	wg := &sync.WaitGroup{}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		s := <-sigChan
		log.Infof("Recived signal %v", s)
		log.Info("Initializing graceful shutdown")
		cancel()

		wg.Done()
	}()
	wg.Wait()
	return nil
}

func newLogger(cfg *config.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level, _ = logrus.ParseLevel(logDefaultLevel)
		logger.Errorf("Invalid log level %v, using default level: %v", err, logDefaultLevel)
	}
	logger.SetLevel(level)
	return logger, nil
}
