package event

type MessageType string

const (
	MonitorHeartBeat   MessageType = "HeartBeat" //监控心跳
	WebSocketKeepalive MessageType = "webSocketKeepalive"
	DataSync           MessageType = "DataSync"     //数据同步（增量递增，单位：天）
	DataResponse       MessageType = "DataResponse" //数据响应
)

type MessageOperator string

const (
	AddIncrement  MessageOperator = "addIncrement"
	EditIncrement MessageOperator = "editIncrement"
	QueryOperator MessageOperator = "queryOperator"
)

type MessageSub string

const (
	GlobalMidIncrementKey = "Auto-Increment-Message-Mid-Global"
)
