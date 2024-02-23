package gormx

import (
	"context"
	"fmt"
	"github.com/bhmy-shm/gofks/core/breaker"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/tracex"
	"gorm.io/gorm"
	"log"
	"reflect"
)

type (
	commonSql struct {
		connProv connProvider    //数据库连接
		brk      breaker.Breaker //熔断器
		//事务
		//数据库stat 统计
	}
	connProvider func() (*gorm.DB, error)
	SqlOptions   func(newTx *gorm.DB) *gorm.DB
)

func NewSql(c *gofkConfs.DBConfig) SqlSession {
	conn := &commonSql{
		connProv: func() (*gorm.DB, error) {
			return SqlInit(c)
		},
		brk: breaker.NewBreaker(),
	}
	return conn
}

func NewSqlFromDB(db *gorm.DB) SqlSession {
	conn := &commonSql{
		connProv: func() (*gorm.DB, error) {
			return db, nil
		},
		brk: breaker.NewBreaker(),
	}
	return conn
}

func (c *commonSql) RawDB() *gorm.DB {
	db, err := c.connProv()
	if err != nil {
		errorx.Fatal(err)
	}
	return db
}

func (c *commonSql) QueryScanContext(ctx context.Context, result interface{}, sql string, args ...string) error {

	var (
		err error
		db  = c.RawDB()
	)

	ctx, span := tracex.StartDbSpan(ctx, "QueryScan")
	defer func() {
		tracex.EndDbSpan(span, err)
	}()

	if len(args) > 0 {
		db = db.Raw(sql, args)
	} else {
		db = db.Raw(sql)
	}

	err = db.Scan(result).Error
	log.Println(result)
	return err
}

func (c *commonSql) QueryCountContext(ctx context.Context, total *int64, sql string, args ...string) error {

	var (
		err error
		db  = c.RawDB()
	)

	ctx, span := tracex.StartDbSpan(ctx, "QueryCount")
	defer func() {
		tracex.EndDbSpan(span, err)
	}()

	if len(args) > 0 {
		db = db.Raw(sql, args)
	} else {
		db = db.Raw(sql)
	}

	return db.Scan(total).Error
}

func (c *commonSql) ExecSqlContext(ctx context.Context, sql string, args ...string) (bool, error) {

	var (
		err  error
		rows int64
	)

	ctx, span := tracex.StartDbSpan(ctx, "ExecSQL")
	defer func() {
		tracex.EndDbSpan(span, err)
	}()

	if len(args) > 0 {
		rows = c.RawDB().Exec(sql, args).RowsAffected
	} else {
		rows = c.RawDB().Exec(sql).RowsAffected
	}

	return rows > 0, err
}

func (c *commonSql) TransactCtx(ctx context.Context, fn func(context context.Context, tx *gorm.DB) error) error {
	var err error

	ctx, span := tracex.StartDbSpan(ctx, "transaction")
	defer func() {
		tracex.EndDbSpan(span, err)
	}()

	return c.RawDB().Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (c *commonSql) QueryFromDB(ctx context.Context, db *gorm.DB, result interface{}) error {
	var err error

	ctx, span := tracex.StartDbSpan(ctx, "QueryFromDB")
	defer func() {
		tracex.EndDbSpan(span, err)
	}()

	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("【QueryFromDB】 错误传参")
	}

	return db.Find(&result).Error
}

func (c *commonSql) CountFromDB(ctx context.Context, db *gorm.DB, total *int64) error {
	var (
		err error
	)

	ctx, span := tracex.StartDbSpan(ctx, "CountFromDB")
	defer func() {
		tracex.EndDbSpan(span, err)
	}()

	return db.Count(total).Error
}
