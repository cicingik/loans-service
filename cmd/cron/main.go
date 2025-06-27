// Package main ...
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cicingik/loans-service/config"
	funding "github.com/cicingik/loans-service/repository/loan_funding"
	"github.com/cicingik/loans-service/repository/loans"
	"github.com/cicingik/loans-service/repository/postgre"
	"github.com/cicingik/loans-service/services"
	env "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var cfg *config.AppConfig

// Cron ...
type Cron struct {
	cfg *config.AppConfig
	db  *postgre.DbEngine
}

// NewCron ...
func NewCron(cfg *config.AppConfig, db *postgre.DbEngine) *Cron {
	return &Cron{
		cfg: cfg,
		db:  db,
	}
}

// Shutdown the service
func (c *Cron) Shutdown(_ context.Context, cancel context.CancelFunc) {
	log.Info("begin shutting down service...")

	defer func() {
		cancel()
	}()

	log.Info("shutdown procedure completed!")
}

// WatchSignal ...
func (c *Cron) WatchSignal(killMe chan<- bool) {
	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// block here, wait for signal
	osCall := <-sigint
	log.Infof("system call:%+v", osCall)

	// in case got the signal, notify parent signal
	killMe <- true
	log.Debugf("watchsignal routine done")
}

func loadConfig() {
	// load .env file when existed.
	_, err := os.Stat(".env")
	if err == nil {
		if err := env.Load(); err != nil {
			log.Error("error loading .env: ", err.Error())
		}
	}

	cfg = config.LoadConfig()
}

// GenerateLenderAgreeement ...
func (c *Cron) GenerateLenderAgreeement() error {
	lrepo, _ := loans.NewLoanRepository(c.cfg, c.db)
	frepo, _ := funding.NewLoanFundingRepository(c.cfg, c.db)

	_, err := lrepo.Invested()
	if err != nil {
		log.Errorf("failed update invested loan")
	}

	funding, _ := frepo.NoLenderAgreement()
	if err != nil {
		log.Errorf("failed loan funding")
	}

	if funding != nil {
		lenderAgreement, _ := services.CreateLenderAgreement(*funding)
		frepo.UpdateLenderAgreemnt(lenderAgreement, funding.ID)
	}

	return nil
}

func main() {
	loadConfig()

	sigShutdown := make(chan bool, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	db := postgre.NewDbService(*cfg)

	app := NewCron(cfg, db)

	go func() {
		app.WatchSignal(sigShutdown)
	}()

	go func() {
		if err := app.GenerateLenderAgreeement(); err != nil {
			log.Errorf("caught error on CronJob: %s", err)
		}

		sigShutdown <- true
	}()

	// watch for signal to perform shutdown
	<-sigShutdown

	// upon receive the signal, perform clean shutdown
	log.Infof("kthxbye!")
	app.Shutdown(ctx, cancel)

	os.Exit(0)
}
