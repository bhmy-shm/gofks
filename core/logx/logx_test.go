package logx

import (
	"fmt"
	"log"
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

	conf := LogConf{
		ServiceName: "test-logx",
		Mode:        fileMode,
		Level:       levelInfo,
		Path:        "logs", //文件写入的必要操作，日志文件目录
		KeepDays:    1,
	}

	if err := SetUp(conf); err != nil {
		log.Println("err:", err)
	}

	Infow("foo 1 console:%s xxx", test1.f)
	Infow("foo 2 console:%s xxx", test1.f)
	Infow("foo 3 console:%s xxx", test1.f)

	fmt.Println("successful")
	time.Sleep(time.Second * 1)
}

//func TestRotateLoggerWrite(t *testing.T) {
//	//filename, err := fs.TempFilenameWithText("foo")
//	filename := "haha"
//	err := errorx.New(errorx.ErrCodeDbErr)
//	if err != nil {
//		log.Println("fileName failed:", err)
//	}
//	log.Println("filename:", filename)
//
//	rule := new(DailyRotateRule)
//	logger, err := NewLogger(filename, rule, true)
//
//	if len(filename) > 0 {
//		defer func() {
//			os.Remove(logger.getBackupFilename())
//			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
//		}()
//	}
//
//	// the following write calls cannot be changed to Write, because of DATA RACE.
//	logger.write([]byte(`foo`))
//	//rule.rotatedTime = time.Now().Add(-time.Hour * 24).Format(dateFormat)
//	logger.write([]byte(`bar`))
//	logger.Close()
//	logger.write([]byte(`baz`))
//}
