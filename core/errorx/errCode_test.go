package errorx

import (
	"log"
	"os"
	"testing"
)

func f1Err() error {
	err := f2Err()

	return Wrap(err, "f1 调用 f2 出错")
}

func f2Err() error {
	err := f3Err()
	return Wrap(err, "asdsa")
}

func f3Err() error {
	return New(ErrCodeParamsErr)
}

func FileOpen(path string) error {
	_, err := os.Open(path)
	if err != nil {
		log.Println("open failed:", err)
		return Wrap(err, "os.Open is failed")
	}
	return nil
}

func TestErrWrapCode(t *testing.T) {
	err1 := FileOpen("jhah")
	log.Println("err1:", Stack(err1))

}

func TestErrCauseCode(t *testing.T) {

	err2 := f2Err()
	if err2 != nil {
		aa := Stack(err2)

		log.Println("err2 Cause", aa)
	}
}

func TestErrStackCode(t *testing.T) {
	err := f1Err()
	if err != nil {
		log.Println("stack:", Stack(err))
	}
}

func TestWrapErr(t *testing.T) {

	var newErr error

	err := FileOpen("aaxx")
	if err != nil {
		newErr = WrapErr(err, ErrCodeBusy)
	}

	log.Println(newErr)
}
