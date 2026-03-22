package gormx

import (
	"fmt"
	
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openSQLServer(params *Params, logs logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&encrypt=%t",
		params.User,
		params.Password,
		params.Host,
		params.Port,
		params.Name,
		params.Encrypt,
	)
	
	openDB, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger:                 logs,
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenSQLServer, err)
	}
	
	sqlDB, err := openDB.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenSQLServer, err)
	}
	
	sqlDB.SetMaxOpenConns(params.MaxOpenConns)
	sqlDB.SetMaxIdleConns(params.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(params.ConnMaxLifetime)
	
	return openDB, nil
}
