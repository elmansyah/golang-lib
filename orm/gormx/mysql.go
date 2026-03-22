package gormx

import (
	"fmt"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openMySQL(params *Params, logs logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=%t&charset=utf8mb4&loc=Local&timeout=5s&readTimeout=5s&writeTimeout=5s",
		params.User,
		params.Password,
		params.Host,
		params.Port,
		params.Name,
		params.ParseTime,
	)
	
	openDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logs,
		SkipDefaultTransaction: params.SkipDefaultTransaction,
		TranslateError:         params.TranslateError,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenMySQL, err)
	}
	
	sqlDB, err := openDB.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenMySQL, err)
	}
	
	sqlDB.SetMaxOpenConns(params.MaxOpenConns)
	sqlDB.SetMaxIdleConns(params.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(params.ConnMaxLifetime)
	
	return openDB, nil
}
