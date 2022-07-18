package gofk

import "github.com/gin-gonic/gin"

type Responder interface {
	RespondTo() gin.HandlerFunc
}
