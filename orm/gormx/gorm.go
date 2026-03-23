package gormx

import (
	"context"
	"database/sql"
	"fmt"
	
	"gorm.io/gorm"
)

type Init func(params *Params) (App, error)

func New(params *Params) (App, error) {
	connectDB, err := params.Setup()
	if err != nil {
		return nil, err
	}
	
	return connectDB, nil
}

func (params *Params) Get() *gorm.DB {
	return params.DB
}

func (params *Params) Ping(ctx context.Context) error {
	db, err := params.DB.DB()
	if err != nil {
		return fmt.Errorf("%w: %w", errGetDB, err)
	}
	
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("%w: %w", errPingDB, err)
	}
	
	return nil
}

func (params *Params) Health() error {
	return params.Ping(context.Background())
}

func (params *Params) Stats() sql.DBStats {
	db, err := params.DB.DB()
	if err != nil {
		return sql.DBStats{}
	}
	
	return db.Stats()
}

func (params *Params) Close() error {
	db, err := params.DB.DB()
	if err != nil {
		return fmt.Errorf("%w: %w", errGetDB, err)
	}
	
	if err = db.Close(); err != nil {
		return fmt.Errorf("%w: %w", errCloseDB, err)
	}
	
	return nil
}
