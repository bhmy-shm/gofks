package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespField(k string, v interface{}) gin.H {
	return gin.H{
		k: v,
	}
}

func InternalResp(c *gin.Context, h gin.H) {
	c.JSON(http.StatusInsufficientStorage, h)
}

func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
