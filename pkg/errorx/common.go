package errorx

import "errors"

var (
	FileNotExist    = errors.New("the file in the path does not exist")
	FileReadFail    = errors.New("failed to read the file content. Procedure")
	WatcherFileStop = errors.New("water stopped")
)
