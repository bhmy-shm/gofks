package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Writer interface {
		Close() error
		Error(v interface{}, fields ...LogField)
		Info(v interface{}, fields ...LogField)
		Stat(v interface{}, fields ...LogField)
		Slow(v interface{}, fields ...LogField)
		Write(p []byte) (int, error)
	}

	atomicWriter struct {
		writer Writer
		lock   sync.RWMutex
	}

	concreteWriter struct {
		buf      bytes.Buffer
		infoLog  io.WriteCloser
		errorLog io.WriteCloser
	}
)

func NewWriter(w io.Writer) Writer {
	lw := newLogWriter(log.New(w, "", flags))

	return &concreteWriter{
		infoLog:  lw,
		errorLog: lw,
	}
}

func getWriter() Writer {
	w := writer.Load()
	if w == nil {
		w = newConsoleWriter()
		writer.Swap(w)
	}
	return w
}

func SetWriter(w Writer) {
	if writer.Load() == nil {
		writer.Store(w)
	}
}

// atomicWriter

func (w *atomicWriter) Load() Writer {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.writer
}

func (w *atomicWriter) Store(v Writer) {
	w.lock.Lock()
	w.writer = v
	w.lock.Unlock()
}

func (w *atomicWriter) Swap(v Writer) Writer {
	w.lock.Lock()
	old := w.writer
	w.writer = v
	w.lock.Unlock()
	return old
}

// concreteWriter

func newConsoleWriter() Writer {
	outLog := newLogWriter(log.New(os.Stdout, "", flags))
	errLog := newLogWriter(log.New(os.Stderr, "", flags))
	return &concreteWriter{
		infoLog:  outLog,
		errorLog: errLog,
	}
}

func newFileWriter(c gofkConfs.LogConfig) (Writer, error) {
	var (
		err      error
		fns      []LogOptionFunc
		infoLog  io.WriteCloser
		errorLog io.WriteCloser
	)

	if len(c.Path()) == 0 {
		return nil, errorx.Wrap(errorx.ErrCodeLogPathNotSet, "writer notFound Path")
	}

	fns = append(fns, WithCoolDownMillis(c.StackCoolDownMillis()))
	if c.Compress() {
		fns = append(fns, WithGzip())
	}
	if c.KeepDays() > 0 {
		fns = append(fns, WithKeepDays(c.KeepDays()))
	}

	infoFile := path.Join(c.Path(), infoFilename)
	errorFile := path.Join(c.Path(), errorFilename)

	for _, fn := range fns {
		fn(&options)
	}

	SetLevel(c)

	if infoLog, err = createOutput(infoFile); err != nil {
		return nil, err
	}

	if errorLog, err = createOutput(errorFile); err != nil {
		return nil, err
	}

	return &concreteWriter{
		infoLog:  infoLog,
		errorLog: errorLog,
	}, err
}

func (w *concreteWriter) Close() error {

	if err := w.infoLog.Close(); err != nil {
		return err
	}

	if err := w.errorLog.Close(); err != nil {
		return err
	}

	return nil
}

func (w *concreteWriter) Error(v interface{}, fields ...LogField) {
	output(w.errorLog, levelError, v, fields...)
}

func (w *concreteWriter) Info(v interface{}, fields ...LogField) {
	output(w.infoLog, levelInfo, v, fields...)
}

func (w *concreteWriter) Stat(v interface{}, fields ...LogField) {
	output(w.infoLog, levelStat, v, fields...)
}

func (w *concreteWriter) Slow(v interface{}, fields ...LogField) {
	output(w.errorLog, levelSlow, v, fields...)
}

func (w concreteWriter) Write(data []byte) (n int, err error) {
	return w.buf.Write(data)
}

func createOutput(path string) (io.WriteCloser, error) {
	if len(path) == 0 {
		return nil, errorx.ErrCodeLogPathNotSet
	}

	return NewLogger(path, DefaultRotateRule(path, backupFileDelimiter,
		options.keepDays, options.gzipEnabled), options.gzipEnabled)
}

func output(writer io.Writer, level string, val interface{}, fields ...LogField) {
	fields = append(fields, Field(callerKey, getCaller(callerDepth)))

	switch atomic.LoadUint32(&encoding) {
	case jsonEncodingType:
		entry := make(map[string]interface{})
		for _, field := range fields {
			entry[field.Key] = field.Value
		}
		entry[timestampKey] = getFormatTime()
		entry[levelKey] = level
		entry[contentKey] = val
		writeJson(writer, entry)
	}
}

func writeJson(writer io.Writer, info interface{}) {
	if content, err := json.MarshalIndent(info, "", "	"); err != nil {
		log.Printf("[logx-write] writeJson Marshal failed:%v", err.Error())
	} else if writer == nil {
		log.Println(content)
	} else {

		content = bytes.ReplaceAll(content, []byte("\\n"), []byte("\n"))
		content = bytes.ReplaceAll(content, []byte("\\t"), []byte("\t"))

		_, err = writer.Write(append(content, '\n'))
		if err != nil {
			log.Printf("[logx-write] writeJson write failed:%v", err)
		}
	}
}

func Field(key string, value interface{}) LogField {
	switch val := value.(type) {
	case error:
		return LogField{Key: key, Value: val.Error()}
	case []error:
		var errs []string
		for _, err := range val {
			errs = append(errs, "\n"+err.Error())
		}
		return LogField{Key: key, Value: errs}
	case time.Duration:
		return LogField{Key: key, Value: fmt.Sprint(val)}
	default:
		return LogField{Key: key, Value: val}
	}
}
