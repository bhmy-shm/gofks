package errorx

import (
	"fmt"
	"runtime"
)

type stack []uintptr

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func stackToString(s *stack) string {
	str := ""
	for _, pc := range *s {
		funcName := runtime.FuncForPC(pc).Name()
		str += fmt.Sprintf("%s\n", funcName)
	}
	return str
}
