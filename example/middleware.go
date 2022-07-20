package user

import "github.com/gin-gonic/gin"

//只服务于 userCtl 的中间件

func (this *User) RequestCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("web_name", "hatMajor")
		c.Set("web_date", "2021-07-01")
		c.Set("web_version", "0.1")
		c.Set("user_name", "孙海铭")
		c.Next()
	}
}
