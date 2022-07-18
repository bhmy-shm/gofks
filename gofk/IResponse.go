package gofk

import "github.com/gin-gonic/gin"

type Fairing interface {
	OnRequest(ctx *gin.Context) error
}

type Responder interface {
	RespondTo() gin.HandlerFunc
}
