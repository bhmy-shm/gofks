package event

import (
	"encoding/xml"
)

type (
	XMLMsg struct {
		XMLName xml.Name   `xml:"Msg"`
		Head    *XMLHeader `xml:"Head"`
		Body    *XMLBody   `xml:"Body"`
	}
	XMLHeader struct {
		Id      uint64 `xml:"id"`      //当前消息ID
		Type    string `xml:"msgType"` //消息类型一级类
		SubType string `xml:"subType"` //消息子类型（二级分类）

		//执行的操作[ack,select,increment]
		//ack：当前消息为编辑响应消息
		//select：当前消息为查询响应消息
		//increment：当前消息为增量推送消息
		Operator string `xml:"operator"`
	}
	XMLBody struct {
		Items []struct {
			Data string `xml:",innerxml"`
		} `xml:"items"`
	}
)

func defaultXmlHeader() *XMLHeader {
	return &XMLHeader{}
}

func AbleXmlMsg() *XMLMsg {
	return &XMLMsg{}
}

func (m *XMLMsg) getId() uint64 {
	return m.Head.Id
}

func (m *XMLMsg) PackEncode() []byte {
	if m.Body == nil || m.Head.Type == "" {
		return nil
	}

	return m.Marshal()
}

func (m *XMLMsg) PackDecode(data []byte, value interface{}) error {
	return nil
	//msg, err := m.UnMarshal(data)
	//if err != nil {
	//	log.Println(err)
	//	return nil, err
	//}
	//
	////如果消息头和消息体有一个没有拿到，则返回发送异常
	//XmlMsg := msg.(*XMLMsg)
	//if XmlMsg.Body == nil || XmlMsg.Head == nil {
	//	return nil, fmt.Errorf("incorrect delivery content")
	//}
	//
	////封装成msg进行返回
	//return &Message[*XMLMsg]{
	//	Able: XmlMsg,
	//	Opts: &MsgOptions{
	//		Id: XmlMsg.Head.Id,
	//		JsonNotify: &JsonNotify{
	//			Operator: MessageOperator(XmlMsg.Head.Operator),
	//			Type:     MessageType(XmlMsg.Head.Type),
	//			SubType:  MessageSub(XmlMsg.Head.SubType),
	//		},
	//		XMLBody: XmlMsg.Body,
	//	},
	//}, nil
}

func (m *XMLMsg) Marshal() []byte {
	data, err := xml.Marshal(m)
	if err != nil {
		return nil
	}
	return data
}

func (m *XMLMsg) UnMarshal(buf []byte) (PackageMsg, error) {
	err := xml.Unmarshal(buf, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *XMLMsg) SetMessage(option *MsgOptions) {
	m.Head = &XMLHeader{
		Id:       option.Id,
		Type:     string(option.JsonNotify.Type),
		SubType:  string(option.JsonNotify.SubType),
		Operator: string(option.JsonNotify.Operator),
	}

	m.Body = option.XMLBody
}

func (m *XMLMsg) GetMessage() []byte {
	return m.PackEncode()
}
