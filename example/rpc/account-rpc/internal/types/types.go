package types

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"google.golang.org/grpc/codes"
	"net/http"
)

func SuccessStatus() *user.Status {
	return &user.Status{
		Code:    http.StatusOK,
		Message: codes.OK.String(),
	}
}

func InterStatus(err error) *user.Status {
	ex := errorx.New(err)
	return &user.Status{
		Code:     ex.Code,
		Reason:   ex.Reason,
		Message:  ex.Message,
		Metadata: make(map[string]string),
	}
}
