package lotus

import (
	"context"
	"github.com/gin-gonic/gin"
	"oplian/define"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	request1 "oplian/model/lotus/request"
	"oplian/service"
	"oplian/service/pb"
	"oplian/utils"
	"time"
)

// QueryAsk
// @Tags      QueryAsk
// @Summary   QueryAsk
// @Summary   inquiry
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetMinerID
// @Success   200   {object}  pb.AskInfo
// @Router    /deploy/queryAsk [get]
func (deploy *DeployApi) QueryAsk(c *gin.Context) {
	var l request1.GetMinerID
	err := c.ShouldBindQuery(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.ActorVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// find available lotus
	lotusList, _, err := service.ServiceGroupApp.LotusServiceGroup.GetLotusList(request1.LotusInfoPage{SyncStatus: define.SyncFinish.Int(), DeployStatus: define.DeployFinish.Int()})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(lotusList) == 0 {
		response.OkWithMessage("No available lotus", c)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15)*time.Second)
	defer cancel()
	var ask *pb.AskInfo
	online := false
	for _, info := range lotusList {
		if info.Online {
			online = true
			gclient := global.GateWayClinets.GetGateWayClinet(info.GateId)
			if gclient == nil {
				global.ZC_LOG.Error(info.GateId + " not exist！")
				continue
			}
			param := &pb.QueryParam{
				LotusIp:    info.Ip,
				LotusToken: info.Token,
				Param:      l.Actor,
			}
			ask, err = gclient.QueryAsk(ctx, param)
			if err != nil {
				//Inquiry failed
				if ctx.Err() == context.DeadlineExceeded {
					response.OkWithMessage(err.Error(), c)
					return
				}
				global.ZC_LOG.Error("Inquiry failed！" + err.Error())
				continue
			}
			break
		}
	}
	if ask == nil {
		response.FailWithMessage("Inquiry failed！", c)
		return
	}
	if !online {
		response.FailWithMessage("No online lotus！", c)
		return
	}
	response.OkWithData(ask, c)
}

// QueryDataCap
// @Tags      QueryDataCap
// @Summary   Query DC Quota
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.WalletAddr
// @Success   200   {object}  string
// @Router    /deploy/queryDataCap [get]
func (deploy *DeployApi) QueryDataCap(c *gin.Context) {
	var l request.WalletAddr
	err := c.ShouldBindQuery(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.AddrVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Find available lotus
	lotusList, _, err := service.ServiceGroupApp.LotusServiceGroup.GetLotusList(request1.LotusInfoPage{SyncStatus: define.SyncFinish.Int(), DeployStatus: define.DeployFinish.Int()})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(lotusList) == 0 {
		response.OkWithMessage("No available lotus", c)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	var dc *pb.String
	online := false
	for _, info := range lotusList {
		//过滤掉测试网测试
		if info.Ip == "10.0.5.48" || info.Ip == "10.0.1.77" {
			continue
		}
		if info.Online {
			online = true
			gclient := global.GateWayClinets.GetGateWayClinet(info.GateId)
			if gclient == nil {
				global.ZC_LOG.Error(info.GateId + " not exist！")
				continue
			}
			param := &pb.QueryParam{
				LotusIp:    info.Ip,
				LotusToken: info.Token,
				Param:      l.Addr,
			}
			dc, err = gclient.QueryDataCap(ctx, param)
			if err != nil {
				global.ZC_LOG.Warn("Query failed！" + err.Error())
				continue
			}
			if dc == nil {
				global.ZC_LOG.Warn("Query failed！" + err.Error())
				continue
			}
			break
		}
	}
	if !online {
		response.OkWithMessage("No online lotus！", c)
		return
	}
	if dc == nil {
		response.OkWithData("0 GiB", c)
		return
	}
	if dc.Value == "" {
		dc.Value = "0 GiB"
	}
	response.OkWithData(dc.Value, c)
}
