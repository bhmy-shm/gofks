package errorx

import (
	"errors"
)

const (
	HTTP_STATUS = "GOFT_STATUS"
)

var (
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
	ErrClosed               = errors.New("pool is closed")

	FileNotExist    = errors.New("the file in the path does not exist")
	FileReadFail    = errors.New("failed to read the file content. Procedure")
	WatcherFileStop = errors.New("water stopped")
)

//func (e *Error) Error() string {
//	return ""
//}
//
//// BadRequest new BadRequest error that is mapped to a 400 response.
//func BadRequest(reason, message string) *Error {
//	return New(400, reason, message)
//}
//
//// New returns an error object for the code, message.
//func New(code int, reason, message string) *Error {
//	return &Error{
//		Code:    int32(code),
//		Message: message,
//		Reason:  reason,
//	}
//}

func Error(err error, format ...string) {
	if err == nil {
		return
	} else {
		errMsg := err.Error()
		if len(format) > 0 {
			errMsg += format[0]
		}
		panic(errMsg)
	}
}
