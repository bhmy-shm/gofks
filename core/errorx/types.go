package errorx

import "fmt"

type StatusFunc func(status *Status)

func WithReason(reason string) StatusFunc {
	return func(status *Status) {
		status.Reason = reason
	}
}

func WithField(key, value string) StatusFunc {
	return func(status *Status) {
		status.Metadata[key] = value
	}
}

func (s *Status) StringTo16() string {
	return fmt.Sprintf("0x%0*x  %s\n", 8, s.Code, s.Message)
}

// ==========================================================

// BadRequest new 请求错误 that is mapped to a 400 response.
func BadRequest() *Error {
	return New(ErrCodeParamsErr)
}

// Unauthorized new 未授权错误 error that is mapped to a 401 response.
func Unauthorized(reason, message string) *Error {
	return New(ErrCodeNotAuthorized)
}

//// Forbidden new Forbidden error that is mapped to a 403 response.
//func Forbidden(reason, message string) *Error {
//	return New(403, reason, message)
//}

// NotFound new NotFound error that is mapped to a 404 response.
func NotFound(reason, message string) *Error {
	return New(ErrCodeNotFound)
}

//// Conflict new Conflict error that is mapped to a 409 response.
//func Conflict(reason, message string) *Error {
//	return New(409, reason, message)
//}

//// InternalServer new InternalServer error that is mapped to a 500 response.
//func InternalServer(reason, message string) *Error {
//	return New(500, reason, message)
//}
//
//// ServiceUnavailable new ServiceUnavailable error that is mapped to an HTTP 503 response.
//func ServiceUnavailable(reason, message string) *Error {
//	return New(503, reason, message)
//}
//
//// GatewayTimeout new GatewayTimeout error that is mapped to an HTTP 504 response.
//func GatewayTimeout(reason, message string) *Error {
//	return New(504, reason, message)
//}
