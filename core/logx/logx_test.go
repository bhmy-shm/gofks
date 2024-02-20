package logx

import (
	"fmt"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/errorx"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type test struct {
	name string
	f    LogField
	want map[string]interface{}
}

func TestLogInfo(t *testing.T) {

	test1 := &test{
		name: "error",
		f:    Field("foo", fmt.Errorf("bar")),
		want: map[string]interface{}{
			"foo": "bar",
		},
	}

	conf := gofkConf.Load(gofkConf.WithPath("application_log.yaml"))

	if err := SetUp(conf.GetLog()); err != nil {
		log.Println("err:", err)
	}

	Infow("foo 1 console xxx", WithField("first", test1.f))
	Infow("foo 2 console xxx", WithField("second", test1.f))
	Infow("foo 3 console xxx", WithField("thread", test1.f))

	Infof("foo 1 console:%s xxx", test1.f)
	Infof("foo 2 console:%s xxx", test1.f)
	Infof("foo 3 console:%s xxx", test1.f)

	fmt.Println("successful")
	time.Sleep(time.Second * 1)
}

func TestRotateLoggerWrite(t *testing.T) {
	//filename, err := fs.TempFilenameWithText("foo")
	filename := "haha"

	err := errorx.New(errorx.ErrCodeDBQueryFailed)
	if err != nil {
		log.Println("fileName failed:", err)
	}
	log.Println("filename:", filename)

	rule := new(DailyRotateRule)
	logger, err2 := NewLogger(filename, rule, true)
	if err2 != nil {
		log.Println("newLogger fileName failed:", err)
	}

	if len(filename) > 0 {
		defer func() {
			os.Remove(logger.getBackupFilename())
			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
		}()
	}

	// the following write calls cannot be changed to Write, because of DATA RACE.
	logger.write([]byte(`foo`))
	//rule.rotatedTime = time.Now().Add(-time.Hour * 24).Format(dateFormat)
	logger.write([]byte(`bar`))
	logger.Close()
	logger.write([]byte(`baz`))
}
