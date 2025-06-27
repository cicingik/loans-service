package postgre

import (
	"errors"
	"fmt"
	"time"

	"github.com/cicingik/loans-service/config"
	"github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	DbEngine struct {
		G   *gorm.DB
		cfg config.AppConfig
	}
)

func NewDbService(cfg config.AppConfig) *DbEngine {
	db, err := initDb(cfg)
	if err != nil {
		fmt.Printf("Cannot connect to %s database. Details %v", cfg.DBConfig.DbDriver, err)
	}

	return &DbEngine{
		G:   db,
		cfg: cfg,
	}
}

func initDb(cfg config.AppConfig) (connection *gorm.DB, err error) {
	if len(cfg.DBConfig.DbDriver) < 1 {
		log.Errorf("Cannot specify database driver")
		return nil, err
	}

	if cfg.DBConfig.DbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			cfg.DBConfig.DbHost, cfg.DBConfig.DbPort, cfg.DBConfig.DbUser, cfg.DBConfig.DbName, cfg.DBConfig.DbPassword,
		)

		// Register augments the provided driver with tracing, enabling it to be loaded by gormtrace.Open.
		sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithServiceName("loan-svc"))

		sqlDb, err := sqltrace.Open("pgx", DBURL)
		if err != nil {
			log.Errorf("Cannot connect to %s database. Details %v", cfg.DBConfig.DbDriver, err)
			return nil, err
		} else {
			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			sqlDb.SetMaxIdleConns(10)

			// SetMaxOpenConns sets the maximum number of open connections to the database.
			sqlDb.SetMaxOpenConns(100)

			// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
			sqlDb.SetConnMaxLifetime(time.Hour)

			log.Infof("Connected to the %s database", cfg.DBConfig.DbDriver)
			// sqlDb.LogMode(cfg.DBConfig.DbDebug == 1)

			sqlDb.SetMaxIdleConns(cfg.DBConfig.MaxIdleConns)
			sqlDb.SetMaxOpenConns(cfg.DBConfig.MaxOpenConns)
			sqlDb.SetConnMaxLifetime(time.Duration(cfg.DBConfig.MaxConnLifetimeSeconds) * time.Second)

			// defer connection.Close()
			// return connection, err
		}

		connection, err = gormtrace.Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{}, gormtrace.WithServiceName(`loan-svc`))
		if err != nil {
			log.Errorf("Cannot trace database. Details %v", err)
			return nil, err
		}

		return connection, nil
	}

	return connection, errors.New("cannot connect to %s database. details missing some parameter")
}
