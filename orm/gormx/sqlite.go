package gormx

import (
	"fmt"
	
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openSQLite(params *Params, logs logger.Interface) (*gorm.DB, error) {
	openDB, err := gorm.Open(sqlite.Open(params.Name), &gorm.Config{
		Logger:                 logs,
		SkipDefaultTransaction: params.SkipDefaultTransaction,
		TranslateError:         params.TranslateError,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenSQLite, err)
	}
	
	sqlDB, err := openDB.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenSQLite, err)
	}
	
	if params.MaxOpenConns == 0 {
		sqlDB.SetMaxOpenConns(1)
	} else {
		sqlDB.SetMaxOpenConns(params.MaxOpenConns)
	}
	
	sqlDB.SetMaxIdleConns(params.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(params.ConnMaxLifetime)
	
	return openDB, nil
}
