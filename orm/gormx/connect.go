package gormx

import (
	"fmt"
	"log"
	"os"
	"time"
	
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func open(params *Params) (*gorm.DB, error) {
	logLevel := logger.Silent
	
	if params.AppMode == "dev" {
		logLevel = logger.Info
	}
	
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logLevel,
		},
	)
	
	db, err := openByDriver(params, newLogger)
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

func openByDriver(params *Params, logs logger.Interface) (*gorm.DB, error) {
	switch params.DBDriver {
	case "postgresql":
		return openPostgres(params, logs)
	case "mysql":
		return openMySQL(params, logs)
	case "sqlserver":
		return openSQLServer(params, logs)
	case "sqlite":
		return openSQLite(params, logs)
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedDBDriver, params.DBDriver)
	}
}
