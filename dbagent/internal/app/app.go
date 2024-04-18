package app

import (
	"context"
	"github.com/VadimGossip/calculator/dbagent/internal/config"
	"github.com/VadimGossip/calculator/dbagent/internal/domain"
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
	grpcServer   *GrpcServer
}

func NewApp(name, configDir string, appStartedAt time.Time) *App {
	return &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}
}

func (app *App) Run() {
	_, cancel := context.WithCancel(context.Background())
	cfg, err := config.Init(app.configDir)
	if err != nil {
		logrus.Fatalf("Config initialization error %s", err)
	}
	app.cfg = cfg
	logrus.Infof("[%s] got config: [%+v]", app.name, app.cfg)

	dbAdapter := NewDBAdapter(app.cfg.Db)
	if err = dbAdapter.Connect(); err != nil {
		logrus.Fatalf("Fail to connect db %s", err)
	}
	app.Factory = newFactory(dbAdapter)

	go func() {
		app.grpcServer = NewGrpcServer(cfg.AppGrpcServer.Port)
		grpcRouter := initGrpcRouter(app)
		if err := app.grpcServer.Listen(grpcRouter); err != nil {
			logrus.Fatalf("Failed to start GRPC server %s", err)
		}
	}()

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
	if err := dbAdapter.Close(); err != nil {
		logrus.Fatalf("Fail to close db %s", err)
	}

	logrus.Infof("[%s] stopped", app.name)
}
