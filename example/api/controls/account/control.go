package account

import (
	"github.com/bhmy-shm/gofks"
	"github.com/bhmy-shm/gofks/example/api/wire"
)

type AccountCase struct {
	*wire.ServiceWire `inject:"-"`
}

func UserController() *AccountCase {
	return &AccountCase{}
}

func (s *AccountCase) Build(gofk *gofks.Gofk) {

	account := gofk.Group("account")

	user := account.Group("user")
	user.Handle("POST", "/userDetail", s.UserDetail)
	user.Handle("POST", "/userList", s.UserList)
	user.Handle("POST", "/userAdd", s.UserAdd)

	org := account.Group("org")
	org.Handle("POST", "/orgDetail", s.OrgDetail)
}

func (s *AccountCase) Name() string {
	return "userCase"
}

func (s *AccountCase) Wire() *wire.ServiceWire {
	return s.ServiceWire
}

// ===========
