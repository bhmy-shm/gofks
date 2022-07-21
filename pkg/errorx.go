package pkg

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
