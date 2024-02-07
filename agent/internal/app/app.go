package app

import (
	"context"
	"github.com/VadimGossip/calculator/agent/internal/config"
	"github.com/VadimGossip/calculator/agent/internal/domain"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	*Factory
	name         string
	configDir    string
	appStartedAt time.Time
	cfg          *domain.Config
}

func NewApp(name, configDir string, appStartedAt time.Time) *App {
	return &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.Init(app.configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	app.cfg = cfg

	dbAdapter := NewDBAdapter()
	if err := dbAdapter.Connect(); err != nil {
		logrus.Fatalf("Fail to connect db %s", err)
	}
	app.Factory = newFactory(app.cfg, dbAdapter)

	if err := app.rabbitService.Run(ctx); err != nil {
		logrus.Fatalf("Fail to run RabbitMQ service %s", err)
	}

	app.rabbitConsumer.Subscribe(ctx)

	logrus.Infof("[%s] started", app.name)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	var running = true
	for running {
		s := <-c
		logrus.Infof("[%s] got signal: [%s]", app.name, s)
		switch s {
		case syscall.SIGINT,
			syscall.SIGTERM:
			running = false
		}
	}
	cancel()
	if err := app.rabbitService.Shutdown(); err != nil {
		logrus.Fatalf("Fail to shutdown RabbitMQ service %s", err)
	}

	if err := dbAdapter.Close(); err != nil {
		logrus.Fatalf("Fail to close db %s", err)
	}

	logrus.Infof("[%s] stopped", app.name)
}
