package gofk

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

const (
	VarPattern       = `[0-9a-zA-Z_]+`
	CompareSign      = ">|>=|<=|<|==|!="
	CompareSignToken = "gt|ge|le|lt|eq|ne"
	ComparePattern   = `^(` + VarPattern + `)\s*(` + CompareSign + `)\s*(` + VarPattern + `)\s*$`
)

//执行表达式，临时方法后期需要修改
func ExecExpr(expr string, data map[string]interface{}) (string, error) {
	tpl := template.New("expr").Funcs(map[string]interface{}{
		"echo": func(params ...interface{}) interface{} {
			return fmt.Sprintf("echo:%v", params[0])
		},
	})

	t, err := tpl.Parse(fmt.Sprintf("{{%s}}", expr))
	if err != nil {
		log.Println("111", err)
		return "", err
	}
	var buf = &bytes.Buffer{}
	err = t.Execute(buf, data)
	if err != nil {
		log.Println("222", err)
		return "", err
	}
	return buf.String(), nil
}
