package gormx

import (
	"fmt"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openPostgres(params *Params, logs logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		params.Host,
		params.User,
		params.Password,
		params.Name,
		params.Port,
		params.SSLMode,
		params.Timezone,
	)
	
	openDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logs,
		SkipDefaultTransaction: params.SkipDefaultTransaction,
		TranslateError:         params.TranslateError,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenPostgres, err)
	}
	
	sqlDB, err := openDB.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errOpenPostgres, err)
	}
	
	sqlDB.SetMaxOpenConns(params.MaxOpenConns)
	sqlDB.SetMaxIdleConns(params.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(params.ConnMaxLifetime)
	
	return openDB, nil
}
