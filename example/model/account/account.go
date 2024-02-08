package account

import (
	"context"
	"github.com/bhmy-shm/gofks/core/cache/nosql/redisx"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/gormx"
	"gorm.io/gorm"
)

var _ UserModel = (*defaultUserModel)(nil)

type (
	UserModel interface {
		userQueryModel
		userEditModel
	}

	userQueryModel interface {
		Trans(context.Context, gormx.TransFunc) error
		Query(context.Context, interface{}, ...gormx.SqlOptions) error
		QueryCount(context.Context, *int64, ...gormx.SqlOptions) error
		QueryTakeScan(ctx context.Context, key string, value interface{}, queryFunc gormx.QueryFunc) error
		QueryTakeTotal(ctx context.Context, key string, value *int64, queryCount gormx.QueryFunc) error
	}

	userEditModel interface {
		AutoMigrates()
		Insert(context.Context, *gorm.DB, *User) error
		FindOne(context.Context, int64) (*User, error)
		FindOneByMobile(context.Context, string) (*User, error)
		Update(context.Context, *User) error
		UpdateWithVersion(context.Context, *User) error
		Delete(context.Context, int64) error
	}

	User struct {
		gorm.Model
		Account string `gorm:"del_state"`
		Name    string `gorm:"version"`
		Pass    string `gorm:"mobile"`
		Gender  string `gorm:"password"`
		Mobile  string `gorm:"nickname"`
		Phone   string `gorm:"sex"`
		Avatar  string `gorm:"avatar"`
		Info    string `gorm:"info"`
	}

	defaultUserModel struct {
		session   gormx.SqlSession
		cache     redisx.CacheSession
		tableName string
	}
)

func NewUserModel(c *gofkConf.Config) UserModel {
	if c.GetDB().IsLoad() && c.GetRedis().IsLoad() {
		model := newUserModel(c)
		model.AutoMigrates()
		return model
	}
	return nil
}

func newUserModel(c *gofkConf.Config) *defaultUserModel {
	return &defaultUserModel{
		session: gormx.NewSql(c.GetDB()),
		cache:   redisx.NewSqlCache(c.GetRedis()),
	}
}

func (m *defaultUserModel) Trans(ctx context.Context, transFunc gormx.TransFunc) error {
	return m.session.TransactCtx(ctx, transFunc)
}

func (m *defaultUserModel) Query(ctx context.Context, result interface{}, opts ...gormx.SqlOptions) error {
	db := m.session.RawDB().Model(&User{})
	for _, opt := range opts {
		opt(db)
	}
	return m.session.QueryFromDB(ctx, db, result)
}

func (m *defaultUserModel) QueryCount(ctx context.Context, total *int64, opts ...gormx.SqlOptions) error {
	db := m.session.RawDB().Model(&User{})
	for _, opt := range opts {
		opt(db)
	}
	return m.session.CountFromDB(ctx, db, total)
}

func (m *defaultUserModel) QueryTakeScan(ctx context.Context, key string, value interface{}, queryFunc gormx.QueryFunc) error {
	return m.cache.TakeCtx(ctx, value, key,
		func(val interface{}) error {
			return queryFunc(ctx, m.session)
		},
	)
}

func (m *defaultUserModel) QueryTakeTotal(ctx context.Context, key string, value *int64, queryCount gormx.QueryFunc) error {
	return m.cache.TakeCtx(ctx, value, key,
		func(val interface{}) error {
			return queryCount(ctx, m.session)
		},
	)
}

// --- ------ -------- ------------- ----------------

func (m *defaultUserModel) Insert(ctx context.Context, tx *gorm.DB, user *User) error {
	return tx.Model(&User{}).Create(user).Error
}

func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}

func (m *defaultUserModel) FindOneByMobile(ctx context.Context, mobile string) (*User, error) {
	return nil, nil
}

func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	return nil
}

func (m *defaultUserModel) UpdateWithVersion(ctx context.Context, data *User) error {
	return nil
}

func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *defaultUserModel) AutoMigrates() {
	err := m.session.RawDB().AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}

func (m *User) TableName() string {
	return "shm_test_user"
}

func WithUserMobile(mobile string) gormx.SqlOptions {
	return func(newTx *gorm.DB) *gorm.DB {
		if len(mobile) > 0 {
			newTx = newTx.Where("mobile = ?", mobile)
		}
		return newTx
	}
}
