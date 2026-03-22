package gormx

import (
	"context"
	"database/sql"
	"errors"
	"time"
	
	"github.com/elmansyah/golang-lib/connection/sshx"
	"gorm.io/gorm"
)

var (
	errGetDB               = errors.New("failed to get database connection")
	errPingDB              = errors.New("failed to ping database connection")
	errCloseDB             = errors.New("failed to close database connection")
	errUnsupportedDBDriver = errors.New("unsupported database driver")
	errOpenPostgres        = errors.New("failed to open postgres connection")
	errOpenMySQL           = errors.New("failed to open mysql connection")
	errOpenSQLServer       = errors.New("failed to open sqlserver connection")
	errOpenSQLite          = errors.New("failed to open sqlite connection")
	errSSHTunnel           = errors.New("failed to establish SSH tunnel")
	errFailedToOpenDB      = errors.New("failed to open database connection")
	errOpenDBNil           = errors.New("database connection is nil")
	errUnknownDBLocation   = errors.New("unknown database location")
)

type Params struct {
	Closed                 func() error
	DB                     *gorm.DB
	Host                   string
	User                   string
	Password               string
	Name                   string
	Timezone               string
	AppMode                string
	DBDriver               string
	DBLocation             string
	DBTunnel               string
	SSHParams              sshx.Params
	Port                   int
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetime        time.Duration
	SSLMode                bool
	SkipDefaultTransaction bool
	TranslateError         bool
	Encrypt                bool
	ParseTime              bool
}

type App interface {
	Get() *gorm.DB
	Ping(context.Context) error
	Health() error
	Stats() sql.DBStats
	Close() error
}
