package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vladjong/user_balance/config"
	postgressql "github.com/vladjong/user_balance/internal/adapters/db/postgres_sql"
	"github.com/vladjong/user_balance/internal/controller/handler"
	"github.com/vladjong/user_balance/internal/usecase"
	"github.com/vladjong/user_balance/pkg/fileworker"
	"github.com/vladjong/user_balance/pkg/postgres"
	"github.com/vladjong/user_balance/pkg/server"
)

type Service struct {
	cfg            *config.Config
	postgresClient *sqlx.DB
}

func NewService(cfg *config.Config) (service Service, err error) {
	postgresClient, err := postgres.NewClient(
		postgres.PostgresConfig{
			Host:     cfg.PostgresSQL.Host,
			Port:     cfg.PostgresSQL.Port,
			Username: cfg.PostgresSQL.Username,
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   cfg.PostgresSQL.DBName,
			SSLMode:  cfg.PostgresSQL.SSLMode,
		})
	if err != nil {
		return service, err
	}
	return Service{
		cfg:            cfg,
		postgresClient: postgresClient,
	}, nil
}

func (s *Service) Run() error {
	logrus.Info("initializing openWeatherApi service storage interface")
	s.startHTTP()
	return nil
}

func (s *Service) startHTTP() {
	logrus.Info("HTTP Server initializing")
	server := new(server.Server)
	fileworker := fileworker.New()
	userBalancePostgres := postgressql.New(s.postgresClient)
	userBalanceUseCase := usecase.New(userBalancePostgres, fileworker)
	handlers := handler.New(userBalanceUseCase)
	go func() {
		if err := server.Run(s.cfg.Listen.Port, handlers.NewRouter()); err != nil {
			logrus.Fatalf("error: occured while running HTTP Server: %s", err.Error())
		}
	}()
	logrus.Info("HTTP Server start")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Info("HTTP Server Shutdown")
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error: occured on server shutdown: %s", err.Error())
	}
	if err := s.postgresClient.Close(); err != nil {
		logrus.Errorf("error: occured on db connection close: %s", err.Error())
	}
}
