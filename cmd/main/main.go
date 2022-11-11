package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vladjong/user_balance/config"
	"github.com/vladjong/user_balance/internal/service"
)

func main() {
	logrus.Info("config initializing")
	cfg := config.GetConfig()
	logrus.Info("env variables initializing")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	logrus.Info("running service")
	service, err := service.NewService(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := service.Run(); err != nil {
		logrus.Fatal(err)
	}
}
