package event

import "time"

type (
	JsonParams struct {
		Params interface{} `json:"params,omitempty"`
	}

	JsonBody struct {
		Data string `json:"data,omitempty"`
	}

	JsonNotify struct {
		ProgramID     string          `json:"programId,omitempty"`     //系统程序id
		Id            uint64          `json:"id,omitempty"`            //消息id
		Type          MessageType     `json:"type,omitempty"`          //消息类型
		SubType       MessageSub      `json:"subType,omitempty"`       //一级类型(二级类型)
		Operator      MessageOperator `json:"operator,omitempty"`      //仅Data类型消息才有操作权限
		Timestamp     int64           `json:"timestamp,omitempty"`     //消息推送的时间戳
		KeepaliveTime time.Duration   `json:"keepaliveTime,omitempty"` //定时消息同步间隔时间
	}

	JsonBase struct {
		Method     string      `json:"method"`
		From       string      `json:"from,omitempty"`       //调用者（登录人）
		To         string      `json:"to,omitempty"`         //被调者
		Origin     string      `json:"origin,omitempty"`     //来源域名（或Ip）
		PlatFormId string      `json:"platformId,omitempty"` //平台id
		Zone       int         `json:"zone,omitempty"`       //时区
		RemoteIP   string      `json:"remoteIP,omitempty"`   //被调者Ip
		JWT        interface{} `json:"jwt,omitempty"`        //来源jwt验证信息
	}

	JsonError struct {
		Code    uint64      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
)

func defaultJsonBase() *JsonBase {
	return &JsonBase{}
}

func defaultJsonNotify() *JsonNotify {
	return &JsonNotify{}
}

func WithId(id uint64) MsgOption {
	return func(options *MsgOptions) {
		options.Id = id
	}
}

func WithMethod(method string) MsgOption {
	return func(options *MsgOptions) {
		options.Method = method
	}
}

func WithRequest(request *JsonRequest) MsgOption {
	return func(options *MsgOptions) {
		options.Id = request.ID
		options.JsonBase = request.JsonBase
	}
}

func WithParams(params interface{}) MsgOption {
	return func(options *MsgOptions) {
		options.JsonParams.Params = params
	}
}

func WithJsonNotify(notify *JsonNotify) MsgOption {
	return func(options *MsgOptions) {
		options.JsonNotify = notify
	}
}

//func WithJsonBase(base JsonBase) MsgOption {
//	return func(options *MsgOptions) {
//
//	}
//}
//
//func WithType(msgType string) MsgOption {
//	return func(options *MsgOptions) {
//		options.JsonNotify.Type = MessageType(msgType)
//	}
//}
//
//func WithOperator(operator string) MsgOption {
//	return func(options *MsgOptions) {
//		options.JsonNotify.Operator = MessageOperator(operator)
//	}
//}
//
//func WithSubType(subType string) MsgOption {
//	return func(options *MsgOptions) {
//		options.JsonNotify.SubType = MessageSub(subType)
//	}
//}
//
//func WithTimestamp(timeStamp int64) MsgOption {
//	return func(options *MsgOptions) {
//		options.Timestamp = timeStamp
//	}
//}
//
//func WithProgramID(programId string) MsgOption {
//	return func(options *MsgOptions) {
//		options.ProgramID = programId
//	}
//}
