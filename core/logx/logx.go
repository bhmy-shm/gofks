package logx

import (
	"fmt"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"io"
	"log"
	"sync/atomic"
)

// 全局变量
var (
	logLevel uint32
	encoding uint32 = jsonEncodingType
	writer          = new(atomicWriter)
	options  logOptions
)

// SetUp sets up the logx.
func SetUp(c *gofkConfs.LogConfig) error {
	SetLevel(*c)

	// encoding 编码方式
	switch c.Encoding() {
	case plainEncoding:
		atomic.StoreUint32(&encoding, plainEncodingType)
	default:
		atomic.StoreUint32(&encoding, jsonEncodingType)
	}

	// mode 日志类型
	switch c.Mode() {
	case fileMode:
		return setupWithFiles(*c)
	default:
		setupWithConsole()
		return nil
	}
}

// SetLevel sets the logging level. It can be used to suppress some logs.
func SetLevel(c gofkConfs.LogConfig) {
	switch c.Level() {
	case levelInfo:
		atomic.StoreUint32(&logLevel, InfoLevel)
	case levelError:
		atomic.StoreUint32(&logLevel, ErrorLevel)
	case levelStat:
		atomic.StoreUint32(&logLevel, StatLevel)
	default:
		//TODO
	}
}

func setupWithConsole() {
	SetWriter(newConsoleWriter())
}

func setupWithFiles(c gofkConfs.LogConfig) error {
	w, err := newFileWriter(c)
	if err != nil {
		return err
	}

	SetWriter(w)
	return nil
}

// =======================================================

type (
	logOptions struct {
		gzipEnabled           bool
		keepDays              int
		logStackCoolDownMills int //调用栈帧间隔
	}

	LogOptionFunc func(options *logOptions)

	logEntry struct {
		Timestamp string      `json:"@timestamp"`
		Level     string      `json:"level"`
		Duration  string      `json:"duration,omitempty"`
		Caller    string      `json:"caller,omitempty"`
		Content   interface{} `json:"content"`
	}

	LogField struct {
		Key   string
		Value interface{}
	}
)

// WithField k-v
func WithField(k string, v interface{}) LogField {
	return LogField{
		Key:   k,
		Value: v,
	}
}

// WithKeepDays customizes logging to keep logs with days.
func WithKeepDays(days int) LogOptionFunc {
	return func(opts *logOptions) {
		opts.keepDays = days
	}
}

// WithGzip customizes logging to automatically gzip the log files.
func WithGzip() LogOptionFunc {
	return func(opts *logOptions) {
		opts.gzipEnabled = true
	}
}

// WithCoolDownMillis customizes logging on writing call stack interval.
func WithCoolDownMillis(millis int) LogOptionFunc {
	return func(opts *logOptions) {
		opts.logStackCoolDownMills = millis
	}
}

// shallLog 日志等级过滤
func shallLog(level uint32) bool {
	return atomic.LoadUint32(&logLevel) <= level
}

func errorTextSync(msg string) {
	log.Println(msg)
	if shallLog(ErrorLevel) {
		getWriter().Error(msg)
	}
}

func errorFieldsSync(content string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Error(content, fields...)
	}
}

func infoTextSync(msg string) {
	if shallLog(InfoLevel) {
		getWriter().Info(msg)
	}
}

func infoFieldsSync(content string, fields ...LogField) {
	if shallLog(InfoLevel) {
		getWriter().Info(content, fields...)
	}
}

func statTextSync(msg string) {
	if shallLog(StatLevel) {
		getWriter().Info(msg)
	}
}

func statFieldSync(content string, fields ...LogField) {
	if shallLog(StatLevel) {
		getWriter().Stat(content, fields...)
	}
}

func slowFieldsSync(content string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(content, fields...)
	}
}

func slowTextSync(msg string) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(msg)
	}
}

// =================== 对外提供调用 =======================

func Close() error {
	if w := writer.Swap(nil); w != nil {
		return w.(io.Closer).Close()
	}
	return nil
}

// Error writes v into error log.
func Error(v ...interface{}) {
	errorTextSync(fmt.Sprint(v...))
}

// Errorf writes v with format into error log.
func Errorf(format string, v ...interface{}) {
	err := fmt.Errorf(format, v...)
	errorTextSync(err.Error())
}

// Errorw writes msg along with fields into error log.
func Errorw(msg string, fields ...LogField) {
	errorFieldsSync(msg, fields...)
}

// Info writes v into error log.
func Info(v ...interface{}) {
	infoTextSync(fmt.Sprint(v...))
}

// Infof writes v with format into error log.
func Infof(format string, v ...interface{}) {
	infoTextSync(fmt.Sprintf(format, v...))
}

// Infow writes msg along with fields into error log.
func Infow(msg string, fields ...LogField) {
	infoFieldsSync(msg, fields...)
}

// Stat writes v into error log.
func Stat(v ...interface{}) {
	statTextSync(fmt.Sprint(v...))
}

// Statf writes v with format into error log.
func Statf(format string, v ...interface{}) {
	statTextSync(fmt.Sprintf(format, v...))
}

// Statw writes msg along with fields into error log.
func Statw(msg string, fields ...LogField) {
	statFieldSync(msg, fields...)
}

// Slow writes v into slow log.
func Slow(v ...interface{}) {
	slowTextSync(fmt.Sprint(v...))
}

// Slowf writes v with format into slow log.
func Slowf(format string, v ...interface{}) {
	slowTextSync(fmt.Sprintf(format, v...))
}

// Sloww writes msg along with fields into slow log.
func Sloww(msg string, fields ...LogField) {
	slowFieldsSync(msg, fields...)
}
