package svc

import (
	pkgConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/model/account"
)

type ServiceContext struct {
	config        *pkgConf.Config
	userModel     account.UserModel     `inject:"-"`
	userAuthModel account.UserAuthModel `inject:"-"`
}

func NewServiceContext(c *pkgConf.Config) *ServiceContext {
	return &ServiceContext{
		config: c,
		//userModel: model.NewUserModel(c),
	}
}

func (s *ServiceContext) UserModel() account.UserModel {
	return s.userModel
}

func (s *ServiceContext) UserAuthModel() account.UserAuthModel {
	return s.userAuthModel
}
