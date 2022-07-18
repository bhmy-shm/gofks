package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type TokenCheck struct {
}

func NewTokenCheck() *TokenCheck {
	return &TokenCheck{}
}

func (t *TokenCheck) OnRequest(ctx *gin.Context) error {
	if ctx.Query("token") == "" {
		return fmt.Errorf("token requred")
	}
	return nil
}

func (t *TokenCheck) OnResponse(result interface{}) (interface{}, error) {
	return result, nil
}
