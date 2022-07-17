package gofk

import (
	"github.com/gin-gonic/gin"
	"github.com/bhmy-shm/gofks/ifac"
	"reflect"
	"sync"
)

var responderList []Responder
var onec_resp_list sync.Once

func getResponderList() []Responder {
	onec_resp_list.Do(func() {
		responderList = []Responder{
			(StringResponder)(nil),
			(ModelResponder)(nil),
			(ModelsResponder)(nil),
		}
	})
	return responderList
}


type Responder interface {
	RespondTo() gin.HandlerFunc
}

func Convert(handle interface{}) gin.HandlerFunc {
	h_ref := reflect.ValueOf(handle)

	//循环遍历Responder接口中实现的几个对象
	for _, responder := range getResponderList() {
		//拿到每一个对象的类型
		r_ref := reflect.TypeOf(responder)

		//如果传入的对象类型 与 ResponderList 中注册的类型相匹配，ConvertibleTo可以进行转换。
		if h_ref.Type().ConvertibleTo(r_ref) {
			return h_ref.Convert(r_ref).Interface().(Responder).RespondTo()
		}
	}
	return nil
}


type StringResponder func(*gin.Context) string

func (s StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.String(200, s(context))
	}
}

type ModelResponder func(*gin.Context) ifac.Model

func (m ModelResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, m(context))
	}
}

type ModelsResponder func(*gin.Context) ifac.Models

func (m ModelsResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Content-type", "application/json")
		context.Writer.WriteString(string(m(context)))
	}
}
