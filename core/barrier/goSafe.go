package barrier

import (
	"github.com/bhmy-shm/gofks/core/logx"
)

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func()) {
	go RunSafe(fn)
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func()) {
	defer Recover()
	fn()
}

// Recover is used with defer to do cleanup on panics.
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.Error(p)
	}
}
