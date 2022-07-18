package gofk

import (
	"fmt"
	"github.com/bhmy-shm/gofks/Injector"
	"github.com/bhmy-shm/gofks/pkg/config"
	"reflect"
	"strings"
)

type Annotation interface {
	SetTag(tag reflect.StructTag)
}

var AnnotationList []Annotation

func IsAnnotation(t reflect.Type) bool {
	for _, item := range AnnotationList {
		if reflect.TypeOf(item) == t {
			return true
		}
	}
	return false
}

func init() {
	AnnotationList = make([]Annotation, 0)
	AnnotationList = append(AnnotationList, new(Value))
}

type Value struct {
	tag reflect.StructTag
}

func (this *Value) SetTag(tag reflect.StructTag) {
	this.tag = tag
}

func (this *Value) String() string {

	get_prefix := this.tag.Get("prefix")
	if get_prefix == "" {
		return ""
	}

	prefix := strings.Split(get_prefix, ".")
	if conf := Injector.BeanFactory.Get((*config.SysConfig)(nil)); conf != nil {
		get_value := config.GetConfigValue(conf.(*config.SysConfig).Config, prefix, 0)
		if get_value != nil {
			return fmt.Sprintf("%v", get_value)
		} else {
			return ""
		}
	} else {
		return ""
	}
}
