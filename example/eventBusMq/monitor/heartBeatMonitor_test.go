package monitor

import (
	"github.com/bhmy-shm/gofks/core/event"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"log"
	"time"
)

const (
	DataAccessPrefix    string = "DataAccess-"
	SystemGroundMonitor string = "Ground"
	SystemAutoMonitor   string = "Auto"
	SystemFlowMonitor   string = "AirportFlow"
)

func MonitorHeartbeat() {

	//创建一个心跳类型
	beat := NewHearBeat[*event.JsonMsg]()

	//注册到全局中心
	err := RegisterMonitor(beat.MqRedis, beat)
	if err != nil {
		log.Println(err)
	}

	//发送到消息队列
	timer := timex.NewTicker(time.Second * 10)
	defer timer.Stop()
	for {
		select {
		case <-timer.Chan():
			err = beat.Send(DataAccessPrefix + SystemGroundMonitor) //没间隔10秒钟发送一次消息信息
		}
	}
}
