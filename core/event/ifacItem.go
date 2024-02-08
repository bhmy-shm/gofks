package event

type Item interface {
	MsgType() string    //消息一级类型
	MsgSubType() string //消息二级类型
	PubBus()            //发送数据到消息总线

	Set(k, v string) //设置
	Updates() error  //更新
}

type Subject interface {
	NotifyUp() error
	EventUp() error
}

type EventMonitorItem interface {
	ProgramID() string          //返回实例的名字
	Type() string               //返回实例的消息类型
	Send(topic ...string) error //发送消息到消息队列
}

type IMessage interface {
	GetId() uint64     //获取消息ID
	GetMethod() string //获取消息协议
	GetBytes() []byte  //获取消息内容
	GetData() string   //获取消息内容的字符串
}
