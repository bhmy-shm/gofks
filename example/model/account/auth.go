package account

import (
	"context"
	"time"
)

type (
	UserAuthModel interface {
		Insert(ctx context.Context, data *UserAuth) error
		FindOne(ctx context.Context, id int64) (*UserAuth, error)
		FindOneByAuthTypeAuthKey(ctx context.Context, authType string, authKey string) (*UserAuth, error)
		FindOneByUserIdAuthType(ctx context.Context, userId int64, authType string) (*UserAuth, error)
		Update(ctx context.Context, data *UserAuth) error
		UpdateWithVersion(ctx context.Context, data *UserAuth) error
		Delete(ctx context.Context, id int64) error
	}

	UserAuth struct {
		Id         int64     `gorm:"id"`
		CreateTime time.Time `gorm:"create_time"`
		UpdateTime time.Time `gorm:"update_time"`
		DeleteTime time.Time `gorm:"delete_time"`
		DelState   int64     `gorm:"del_state"`
		Version    int64     `gorm:"version"` // 版本号
		UserId     int64     `gorm:"user_id"`
		AuthKey    string    `gorm:"auth_key"`  // 平台唯一id
		AuthType   string    `gorm:"auth_type"` // 平台类型
	}
)

func (u *UserAuth) Insert(ctx context.Context, data *UserAuth) error {
	return nil
}

func (u *UserAuth) FindOne(ctx context.Context, id int64) (*UserAuth, error) {
	return nil, nil
}

func (u *UserAuth) FindOneByAuthTypeAuthKey(ctx context.Context, authType string, authKey string) (*UserAuth, error) {
	return nil, nil
}

func (u *UserAuth) FindOneByUserIdAuthType(ctx context.Context, userId int64, authType string) (*UserAuth, error) {
	return nil, nil
}

func (u *UserAuth) Update(ctx context.Context, data *UserAuth) error {
	return nil
}

func (u *UserAuth) UpdateWithVersion(ctx context.Context, data *UserAuth) error {
	return nil
}

func (u *UserAuth) Delete(ctx context.Context, id int64) error {
	return nil
}
