package gormx

import (
	"context"
	"database/sql"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type SqlSession interface {
	RawDB() *gorm.DB
	QueryScanContext(context.Context, interface{}, string, ...string) error
	QueryCountContext(ctx context.Context, total *int64, sql string, args ...string) error
	ExecSqlContext(context.Context, string, ...string) (bool, error)
	TransactCtx(context.Context, func(context context.Context, tx *gorm.DB) error) error

	QueryFromDB(ctx context.Context, db *gorm.DB, result interface{}) error
	CountFromDB(ctx context.Context, db *gorm.DB, total *int64) error
}

type (
	TransFunc     func(context.Context, *gorm.DB) error
	QueryFunc     func(context.Context, SqlSession) error
	QueryTakeFunc func(context.Context, string, interface{}, SqlSession)

	QueryOpts func(context.Context, ...SqlOptions) error
)

func SqlInit(c *gofkConfs.DBConfig) (*gorm.DB, error) {

	var (
		gormDB *gorm.DB
		sqlDB  *sql.DB
		err    error
	)

	gConf := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	switch c.DB.Debug {
	case "info":
		gConf.Logger = logger.Default.LogMode(logger.Info)
	case "error":
		gConf.Logger = logger.Default.LogMode(logger.Error)
	}

	switch c.DB.Types {
	case "PgSQL":
		gormDB, err = gorm.Open(postgres.Open(c.DB.DataSourceName), gConf)
		sqlDB, err = gormDB.DB()
	default:
		gormDB, err = gorm.Open(mysql.Open(c.DB.DataSourceName), gConf)
		sqlDB, err = gormDB.DB()
	}

	sqlDB.SetMaxIdleConns(c.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.DB.MaxLifetime) * time.Second)
	return gormDB, err
}
