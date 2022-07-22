package user

import (
	"github.com/bhmy-shm/gofks/gofk"
	"github.com/bhmy-shm/gofks/pkg/cache/noSql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type User struct {
	Db    *gorm.DB           `inject:"-"`
	Cache *noSql.SimpleCache `inject:"-"`
}

func UserController() *User {
	return &User{}
}

func (this *User) UserDetail(ctx *gin.Context) {
	param := &UserModel{}
	err := ctx.ShouldBindJSON(param)
	if err != nil {
		log.Fatal("should bind json is failed", err)
	}

	var name string
	if this.Db == nil {
		log.Println("db 依赖注入初始化失败了")
	} else {
		err = this.Db.Table("users").Where("id = ?", 1).Find(&name).Error
		if err != nil {
			gofk.InternalResp(ctx, err)
			return
		}
	}

	//this.Cache = noSql.GetCache().
	//	SetOperation(noSql.StringType()).
	//	Tactics(noSql.RegexpInterLength)
	//
	//this.Cache.DBGetter = func() pkg.Value {
	//
	//}
	//
	//this.Cache.GetCache(fmt.Sprintf("%d",101))

	Data := &UserModel{
		Id:      101,
		Name:    name,
		Address: param.Address,
	}
	gofk.Successful(ctx, Data)
}

func (this *User) Build(gofk *gofk.Gofk) {
	user := gofk.Group("Sys")
	user.GET("/userInfo", this.UserDetail)
}

func (this *User) Injector() string {
	return "dbConfig"
}
