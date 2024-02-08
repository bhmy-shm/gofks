package logx

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

/*
getCaller 通过调用runtime.Caller(callDepth)获取调用者的信息，包括文件路径、行号等。
@request callDepth int，表示调用栈的深度。

如果调用成功，则将获取的文件路径和行号传递给prettyCaller()函数进行格式化，并返回格式化后的字符串。
如果调用失败，则返回空字符串。
*/
func getCaller(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return ""
	}

	return prettyCaller(file, line)
}

// 返回指定时间的格式化时间
func getFormatTime() string {
	return time.Now().Format(timeFormat)
}

func getNowDate() string {
	return time.Now().Format(dateFormat)
}

// prettyCaller
// 该函数接受文件路径和行号作为参数，对文件路径进行格式化，以便更易读。
// 它通过查找最后一个'/'字符来确定文件名的起始位置，并返回格式化后的字符串，格式为: "文件名:行号"。
func prettyCaller(file string, line int) string {
	idx := strings.LastIndexByte(file, '/')
	if idx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}

	idx = strings.LastIndexByte(file[:idx], '/')
	if idx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return fmt.Sprintf("%s:%d", file[idx+1:], line)
}
