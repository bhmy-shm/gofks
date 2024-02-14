package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"os"
	"reflect"
	"sort"
	"strings"
)

type (
	ErrorInter interface {
		Is(err error) bool
	}
	Error struct {
		*Status `json:"status,omitempty"`
		Cause   error  `json:"cause,omitempty"`
		Stack   string `json:"-"`
	}
)

// Stack 打印全部错误得堆栈信息
func Stack(err error) string {
	return err.(*Error).Stack
}

// Wrap	携带携带上层error,赋值给Cause 根因后返回
func Wrap(err error, format string, a ...interface{}) error {

	var (
		wrapErr error
		cause   = fmt.Errorf(format, a...)
	)

	if ok := isErrCode(err); ok {
		wrapErr = New(err.(ErrCode)).setCause(err).setStack()
	} else {
		wrapErr = New(err, WithReason(cause.Error())).setStack()
	}

	return wrapErr
}

// WrapErr 原生err 携带 errCode 错误提示
func WrapErr(err error, code ErrCode) error {
	return New(code).setCause(err).setStack()
}

// Cause 返回调用链最底层错误的根因
func Cause(err error) error {
	type causer interface {
		causePrint() error
	}

	for err != nil {
		Err, ok := err.(causer)
		if !ok {
			break
		}
		err = Err.causePrint()
	}
	return err
}

// Fatal panic退出程序
func Fatal(err error, msg ...string) {
	if err == nil {
		return
	}
	if len(msg) == 0 {
		panic(err)
	}
	panic(fmt.Errorf("err=%v,msg=%s \n", err, msg[0]))
}

// ErrStatus api端将rpc的Status转换成 *Error
func ErrStatus(status interface{}) *Error {
	v := reflect.ValueOf(status)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return &Error{
		Status: &Status{
			Code:     v.FieldByName("Code").Uint(),
			Reason:   v.FieldByName("Reason").String(),
			Message:  v.FieldByName("Message").String(),
			Metadata: v.FieldByName("Metadata").Interface().(map[string]string),
		},
	}
}

func New(err error, opts ...StatusFunc) *Error {

	var Err = new(Error)

	if ss, ok := isStatus(err, opts...); ok {
		Err.Status = ss
	}

	Err.Cause = err
	return Err
}

func (e Error) Error() string {
	str := e.Cause.Error()
	e.Cause = nil
	body, _ := json.Marshal(e)
	return string(body) + "\n cause: " + str
}

// Is 传入error 判断是否是自定的 Error
// 通过As 会将传入的 err 赋值给 new(Error) 的 se。
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
	}
	return false
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

// Unwrap 忽略错误直接panic
func (e *Error) Unwrap() {
	if e.Cause != nil {
		panic(e.Cause)
	}
}

// UnwrapFunc 执行额外错误实现
func (e *Error) UnwrapFunc(f func() interface{}) {
	if e.Status.Code != uint64(ErrCodeOK) {
		f()
	}
}

// SetCause 写入错误根因
func (e *Error) setCause(err error) *Error {
	if err != nil {
		e.Cause = err
		e.Status.Reason = err.Error()
	}
	return e
}

// SetStack 写入错误根因
func (e *Error) setStack() *Error {
	e.Stack = stackToString(callers())
	return e
}

// SetMetadata 写入错误的元数据信息
func (e *Error) setMetadata(md map[string]string) *Error {
	if md != nil {
		e.Metadata = md
	}
	return e
}

// 输出根因
func (e *Error) causePrint() error {
	return e.Cause
}

// ===================================================
// ================ 输出errCode 到文件 ================
// ===================================================

type ErrorModels []*Status

func CreateBaseErrCode() {

	var errIndex ErrorModels

	for _, n := range errMap {
		errIndex = append(errIndex, n)
	}

	sort.Slice(errIndex, func(i, j int) bool {
		return errIndex[i].Code < errIndex[j].Code
	})

	var (
		content10 strings.Builder
		content16 strings.Builder
	)

	for _, n := range errIndex {
		content10.WriteString(n.String())
		content16.WriteString(n.StringTo16())
	}

	err := os.WriteFile("ErrCode-10.txt", []byte(content10.String()), 0644)
	if err != nil {
		fmt.Println("write ErrCode-10 failed:", err.Error())
	}

	err = os.WriteFile("ErrCode-16.txt", []byte(content16.String()), 0644)
	if err != nil {
		fmt.Println("write ErrCode-16 failed:", err.Error())
	}
}
