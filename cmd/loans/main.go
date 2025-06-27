// Package main ...
package main

import (
	"context"
	"os"
	"time"

	"github.com/cicingik/loans-service/app"
	"github.com/cicingik/loans-service/config"
	env "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var cfg *config.AppConfig

func init() {
	// load .env file when existed.
	_, err := os.Stat(".env")
	if err == nil {
		if err := env.Load(); err != nil {
			log.Error("error loading .env: ", err.Error())
		}
	}

	cfg = config.LoadConfig()
}

func main() {
	webApp, err := app.NewWebApplication(cfg)
	if err != nil {
		log.Error("found error: ", err)
		os.Exit(1)
	}

	if err := webApp.Start(context.Background()); err != nil {
		log.Error("found error: ", err)
		os.Exit(1)
	}

	log.Infof("app started")

	<-webApp.Done() // wait for signal

	log.Infof("app is shutting down")

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := webApp.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
