package event

import (
	"encoding/json"
)

type (
	JsonMsg struct {
		*JsonRequest
		*JsonParams
	}

	JsonRequest struct {
		ID uint64 `json:"id"`
		*JsonBase
	}
)

func AbleJsonMsg() *JsonMsg {
	return &JsonMsg{}
}

func (m *JsonMsg) GetId() uint64 {
	return m.ID
}

func (m *JsonMsg) GetRequest() *JsonRequest {
	return m.JsonRequest
}

func (m *JsonMsg) GetParams() interface{} {
	return m.JsonParams.Params
}

// ------

func (m *JsonMsg) PackEncode() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return data
}

func (m *JsonMsg) PackDecode(buf []byte, value interface{}) error {
	return json.Unmarshal(buf, value)
}

// ----------

func (m *JsonMsg) SetMessage(option *MsgOptions) {

	m.JsonRequest = &JsonRequest{
		ID:       option.Id,
		JsonBase: option.JsonBase,
	}
	m.JsonParams = option.JsonParams
}

func (m *JsonMsg) GetMessage() []byte {
	return m.PackEncode()
}
