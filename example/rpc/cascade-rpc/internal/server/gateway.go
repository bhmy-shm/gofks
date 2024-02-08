package server

import (
	"github.com/bhmy-shm/gofks/core/cache/nosql/redisx"
	"github.com/bhmy-shm/gofks/core/gormx"
	"github.com/bhmy-shm/gofks/example/model/account"
	cascadeModel "github.com/bhmy-shm/gofks/example/model/cascade"
)

type (
	GateWayConfig struct {
		Enabled           map[string]struct{}
		SyncTimeEnable    int
		OutboundAddr      string
		NetroutetoolPort  int
		SipRegsvrAddr     string
		CacheExpireSecs   int
		HeartBeatInterval int
	}

	Gateway struct {
		stopCh chan struct{}  //网关结束运行开关
		conf   *GateWayConfig //配置参数
		cache  redisx.Redis
		db     gormx.SqlSession
	}
)

type CascadeLower struct {
	data          cascadeModel.Cascade
	jwtUserClaims *account.JwtUserClaims
}

type HttBaseAuthRequest struct {
	token     *string
	cascadeId string
}

//func (g *Gateway) Start() {
//
//	//下级平台心跳
//	go g.lowerHeartBeat()
//
//	//上级平台心跳
//	go g.superiorHeartBeat()
//
//	//下级平台校时
//
//	done := false
//	for !done {
//		select {
//		case <-g.stopCh: // 接收到停止信号时退出循环
//			done = true
//		}
//	}
//}

//func (g *Gateway) lowerHeartBeat() {
//
//	ticker := timex.NewTicker(time.Second * 30)
//	defer ticker.Stop()
//	done := false
//
//	lm := redsync.NewLockerManager(g.cache.Client())
//
//	for !done {
//		select {
//		case <-ticker.Chan():
//
//			locker, ok, err := lm.TryLock("gw_heartbeat_lock", time.Second*time.Duration(g.conf.HeartBeatInterval*3/2))
//
//			if !ok {
//				logx.Error("GateWay StartData busy waite", err)
//				break
//			}
//			logx.Info("cascadegw StartLower HeartBeat is start")
//
//			defer locker.UnLock()
//			sn, snlist, lowers, err := g.StartData()
//			if err != nil {
//				logx.Error("GateWay StartData failed", err)
//				break
//			}
//			cascadePro.HeartBeat
//			hb := proCascadegw.ReqHeartbeat{}
//
//			hb.Method = "cascade.heartBeat"
//
//			hb.Params.CbToken = ""
//
//			hb.Params.Sn = sn
//			if len(snlist) > 0 {
//				hb.Params.SnList = snlist
//			}
//
//			ch := make(chan string)
//			for _, v := range lowers {
//
//				go funHeartBeat(hb, v, ch)
//			}
//
//			for range lowers {
//				logx.Info("GateWay hb [done]", <-ch)
//			}
//
//		case <-g.stopCh: // 接收到停止信号时退出循环
//			done = true
//		}
//	}
//}
//
//func (g *Gateway) superiorHeartBeat() {}
//
//func (g *Gateway) StartData() (string, string, []model.CascadeData, error) {
//	reqSearch := &proCascade.ReqLowerSearch{}
//
//	reqSearch.Params.Gid = "root"
//	reqSearch.Params.IgnoreChild = false
//
//	reqSearch.Params.PageNum = 0
//	reqSearch.Params.PageSize = 1000
//
//	respSearch := &proCascade.RespLowerSearch{}
//
//	if error := g.sql.LowerSearch(reqSearch, respSearch); error != nil {
//		log.Error("LowerSearch [db err]", error)
//		return "", "", nil, error
//	}
//
//	req := rpcFlag.ReqSuperiorGet{}
//	resp := rpcFlag.RespSuperiorGet{}
//
//	if error := g.sql.FectchCascadeSuperior(&req, &resp); error != nil {
//		log.Error("FectchCascadeSuperior [db err]", error)
//		return "", "", nil, error
//	}
//
//	reqSN := rpcFlag.ReqLocalSN{}
//	respSN := rpcFlag.RespLocalSN{}
//
//	if error := g.sql.FectchCascadeLocalSN(&reqSN, &respSN); error != nil {
//		log.Error("FectchCascadeLocalSN failed [db err]", error)
//		return "", "", nil, error
//	}
//	sn := respSN.Sn
//	var snlist string
//	if resp.Superior != nil {
//		snlist = respSN.Sn + "." + resp.Superior.SnList
//	}
//	return sn, snlist, respSearch.Lowers, nil
//}
