package lotus

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	request1 "oplian/model/lotus/request"
	response1 "oplian/model/lotus/response"
	system "oplian/model/system/request"
	"oplian/service/pb"
	"oplian/utils"
	"path"
	"strings"
	"sync"
)

type DeployApi struct{}

// AddLotus
// @Tags      AddLotus
// @Summary   Add lotus
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request1.AddLotusInfo
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/addLotus [post]
func (deploy *DeployApi) AddLotus(c *gin.Context) {
	var l request1.AddLotusInfo
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	gclient := global.GateWayClinets.GetGateWayClinet(l.GateId)
	if gclient == nil {
		response.FailWithMessage(l.GateId+" not exist！", c)
		return
	}
	wallets := make([]*pb.OpWallet, len(l.Wallets))
	for i, v := range l.Wallets {
		wallets[i] = &pb.OpWallet{OpId: v.OpId, Address: v.Address}
	}
	var info = &pb.LotusInfo{
		LotusId:       l.Id,
		GateId:        l.GateId,
		OpId:          l.OpId,
		Ip:            l.Ip,
		SecpCount:     l.SecpCount,
		BlsCount:      l.BlsCount,
		ImportMode:    l.ImportMode,
		FileName:      l.FileName,
		WalletNewMode: l.WalletNewMode,
		Wallets:       wallets,
	}
	if _, err = gclient.AddLotus(context.Background(), info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.Ok(c)
}

// DownloadSnapshot
// @Tags      DowloadSnapshot
// @Summary   Download Snapshot
// @accept    application/json
// @Produce   application/json
// @Param     data  body      gid:gatewayID,path
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/dowloadSnapshot [post]
func (deploy *DeployApi) DownloadSnapshot(c *gin.Context) {
	gid := c.Request.FormValue("gid")
	path := c.Request.FormValue("path")
	//todo
	//The cloud the snapshots path
	url := `https://snapshots.calibrationnet.filops.net/minimal/latest.zst`
	//url := `https://snapshots.mainnet.filops.net/minimal/latest.zst`
	//lotus.OpClinets.LockRW.RLock()
	client := global.GateWayClinets.GetGateWayClinet(gid)
	res, err := client.DownloadSnapshot(c, &pb.Downtown{Url: url, Path: path})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(res, c)
	//lotus.OpClinets.LockRW.RUnlock()
}

// GetLotusList
// @Tags      GetLotusList
// @Summary   Get lotus List
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/getLotusList [post]
func (deploy *DeployApi) GetLotusList(c *gin.Context) {
	var pageInfo request1.LotusInfoPage
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := DeployService.GetLotusList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetWalletList
// @Tags      GetWalletList
// @Summary   Gets a list of lotus wallets
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/getWalletList [post]
func (deploy *DeployApi) GetWalletList(c *gin.Context) {
	var req pb.RequestOp
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.RequesOpVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	l, err := DeployService.GetLotusByOpID(req.OpId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.Token = l.Token
	req.Ip = l.Ip

	gclient := global.GateWayClinets.GetGateWayClinet(req.GateId)
	if gclient == nil {
		response.FailWithMessage(req.GateId+",not exist！", c)
		return
	}
	res, err := gclient.GetWalletList(context.Background(), &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result := make([]request1.OpWallet, len(res.Wallets))
	for i, w := range res.Wallets {
		result[i].Balance = w.Balance
		result[i].OpId = w.OpId
		result[i].Address = w.Address
	}

	response.OkWithData(result, c)
}

// GetRoomWalletList
// @Tags      GetRoomWalletList
// @Summary   Get wallet list
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/getWalletList [post]
func (deploy *DeployApi) GetRoomWalletList(c *gin.Context) {
	var param request1.LotusInfoPage
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	ls, err := DeployService.GetRoomAllLotus(param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var wallet pb.WalletList

	for _, v := range ls {
		for _, w := range v.Wallets {
			wallet.Wallets = append(wallet.Wallets, &pb.Wallet{Address: w.Address, OpId: v.OpId, Balance: w.Balance})
		}
	}

	response.OkWithData(wallet, c)
}

// AddMiner
// @Tags      AddMiner
// @Summary   Add miner
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pb.MinerInfo
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/addMiner [put]
func (deploy *DeployApi) AddMiner(c *gin.Context) {
	var l request1.MinerInfo
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//_, err = service.ServiceGroupApp.SystemServiceGroup.GetFileName(&pb.FileNameInfo{GateWayId: l.GateId, FileType: define.ProveFile.Int64()})
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	if l.Actor != "" {
		var param = request1.MinerParam{
			l.Id,
			l.LotusId,
			l.Actor,
			l.Partitions,
			l.IsManage,
			l.IsWdpost,
			l.IsWnpost,
		}
		_, err = DeployService.CheckMinerRole(param)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	//Get lotus token
	info, err := DeployService.GetLotus(l.LotusId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if info.SyncStatus != define.SyncFinish.Int() {
		response.FailWithMessage("The lotus synchronous abnormal state！", c)
		return
	}
	if info.DeployStatus != define.DeployFinish.Int() {
		response.FailWithMessage("The lotus deployment status is abnormal！", c)
		return
	}
	//发起添加lotus RPC请求
	gclient := global.GateWayClinets.GetGateWayClinet(l.GateId)
	if gclient == nil {
		response.FailWithMessage(l.GateId+"not exist！", c)
		return
	}
	if l.Actor == "" {
		wallet, err := gclient.WalletBalance(context.Background(), &pb.FilParam{Token: info.Token, Ip: info.Ip, Param: l.Owner})
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("Failed to get the wallet balance：%s", err.Error()), c)
			return
		}
		if wallet.Balance < 0.01 {
			response.FailWithMessage("owner wallet balance insufficient ", c)
			return
		}

	}

	miner := &pb.MinerInfo{
		AddType:     int32(l.AddType),
		MinerId:     uint64(l.Id),
		LotusId:     l.LotusId,
		OpId:        l.OpId,
		Ip:          l.Ip,
		Actor:       l.Actor,
		LotusToken:  info.Token,
		Partitions:  l.Partitions,
		LotusIp:     info.Ip,
		Owner:       l.Owner,
		IsManage:    l.IsManage,
		IsWdpost:    l.IsWdpost,
		IsWnpost:    l.IsWnpost,
		SectorSize:  l.SectorSize,
		StorageType: int32(l.ColonyType),
	}
	res, err := gclient.AddMiner(context.Background(), miner)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if res.Code != 200 {
		response.FailWithMessage(res.Msg, c)
		return
	}
	response.Ok(c)
}

// CheckMiner
// @Tags      CheckMiner
// @Summary   Check miner role
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request1.MinerParam
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/checkMiner [post]
func (deploy *DeployApi) CheckMiner(c *gin.Context) {
	var param request1.MinerParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	msg := ""
	if param.Actor != "" {
		msg, err = DeployService.CheckMinerRole(param)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	response.OkWithMessage(msg, c)
}

// ModifyMinerRole
// @Tags      ModifyMinerRole
// @Summary   Modifying miner roles
// @accept    application/json
// @Produce   application/json
// @Param     data  body     request1.MinerParam
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/modifyMinerRole [post]
func (deploy *DeployApi) ModifyMinerRole(c *gin.Context) {
	var param request1.MinerParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	_, err = DeployService.CheckMinerRole(param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//Get miner
	miner, err := DeployService.GetMinerRun(uint64(param.Id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	gclent := global.GateWayClinets.GetGateWayClinet(miner.GateId)
	if gclent == nil {
		response.FailWithMessage(miner.GateId+" not exist！", c)
		return
	}

	//Restart service
	var run = &pb.MinerRun{
		Ip:         miner.Ip,
		OpId:       miner.OpId,
		Actor:      param.Actor,
		LotusToken: miner.LotusToken,
		Partitions: param.Partitions,
		IsManage:   param.IsManage,
		IsWdpost:   param.IsWdpost,
		IsWnpost:   param.IsWnpost,
		LotusIp:    miner.LotusIp,
	}
	if param.LotusId != 0 {
		//Change lotus connection
		l, err := DeployService.GetLotus(param.LotusId)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		run.LotusIp = l.Ip
		run.LotusToken = l.Token
	}
	if _, err = gclent.RestartMiner(context.Background(), run); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if param.LotusId != 0 {
		if err = DeployService.UpdateMinerStatusAndLink(param.Id, define.RunStatusRunning.Int(), param.LotusId); err != nil {
			response.FailWithMessage(err.Error(), c)
		}
	}
	if err = DeployService.ModifyMinerRole(param); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// GetMinerList
// @Tags      GetMinerList
// @Summary   Get miner list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/getMinerList [post]
func (deploy *DeployApi) GetMinerList(c *gin.Context) {
	var pageInfo request1.MinerTypePage
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := DeployService.GetMinerList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetMinerSelectList
// @Tags      GetMinerSelectList
// @Summary   获取miner下拉列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /deploy/getMinerList [post]
func (deploy *DeployApi) GetMinerSelectList(c *gin.Context) {
	var req request1.MinerListReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := DeployService.GateMinerSelectList(req.GateId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list,
	}, "获取成功", c)
}

// GetNodeList
// @Tags      GetNodeList
// @Summary   Get node list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/GetNodeList [post]
func (deploy *DeployApi) GetNodeList(c *gin.Context) {
	var pageInfo system.ColonyPageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := DeployService.GetNodeList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetNodesNum
// @Tags      GetNodesNum
// @Summary   Get the number of nodes
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/getNodeNumList [post]
func (deploy *DeployApi) GetNodesNum(c *gin.Context) {
	var gateId request.GetGateID
	err := c.ShouldBindJSON(&gateId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := DeployService.GetNodesNum(gateId.GateId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	gcli := global.GateWayClinets.GetGateWayClinet(gateId.GateId)
	if gcli != nil {
		for i := 0; i < len(list); i++ {
			fpath := path.Join(define.FileGateWayDir, fmt.Sprintf("%s_%s.zip", define.MinerName, list[i].Actor))
			if exit, err := gcli.GatewayFileExist(context.Background(), &pb.String{Value: fpath}); err == nil {
				list[i].MinerFile = exit.Value
			} else {
				global.ZC_LOG.Error(err.Error())
			}
		}
	}
	response.OkWithData(list, c)
}

// AddWorker
// @Tags      AddWorker
// @Summary   Add worker
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pb.AddWorker
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/addWorker [put]
func (deploy *DeployApi) AddWorker(c *gin.Context) {
	var bw pb.BatchWorker
	err := c.ShouldBindJSON(&bw)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	gclient := global.GateWayClinets.GetGateWayClinet(bw.GateId)
	if gclient == nil {
		response.FailWithMessage(bw.GateId+" not exist！", c)
		return
	}
	_, err = gclient.AddWorker(context.Background(), &bw)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// AddStorage
// @Tags      AddStorage
// @Summary   Add storage
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pb.AddWorker
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/addStorage [put]
func (deploy *DeployApi) AddStorage(c *gin.Context) {
	var bs pb.BatchStroage
	err := c.ShouldBindJSON(&bs)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	gclient := global.GateWayClinets.GetGateWayClinet(bs.GateId)
	if gclient == nil {
		response.FailWithMessage(bs.GateId+" not exist！", c)
		return
	}
	//新增存储
	_, err = gclient.AddStorage(context.Background(), &bs)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// AddBoost
// @Tags      AddBoost
// @Summary   Add boost
// @accept    application/json
// @Produce   application/json
// @Param     data  body      pb.boostInfo
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/addBoost [put]
func (deploy *DeployApi) AddBoost(c *gin.Context) {
	var l pb.BoostInfo
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	info, _ := DeployService.GetBoost(l.Id)
	miner, _ := DeployService.GetMiner(l.MinerId)
	if !miner.IsManage {
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("Non-scheduler, no need to set!"), c)
			return
		}
	}

	if l.DcQuotaWallet != "" {
		if strings.Contains(l.DcQuotaWallet, ",") {

			walletAr := strings.Split(l.DcQuotaWallet, ",")
			for _, v := range walletAr {
				vAr := strings.Split(v, "|")
				if len(vAr) != 2 {
					response.FailWithMessage(fmt.Sprintf("Abnormal DC quota wallet:%s", v), c)
					return
				}
				_, err = address.NewFromString(vAr[0])
				if err != nil {
					response.FailWithMessage(fmt.Sprintf("Abnormal DC quota wallet:%s", vAr[0]), c)
					return
				}
			}

		} else {

			vAr := strings.Split(l.DcQuotaWallet, "|")
			if len(vAr) != 2 {
				response.FailWithMessage(fmt.Sprintf("Abnormal DC quota wallet:%s", l.DcQuotaWallet), c)
				return
			}
			_, err = address.NewFromString(vAr[0])
			if err != nil {
				response.FailWithMessage(fmt.Sprintf("Abnormal DC quota wallet:%s", l.DcQuotaWallet), c)
				return
			}
		}
	}

	if l.DcQuotaWallet != info.DcQuotaWallet && info.ID != 0 {
		info.DcQuotaWallet = l.DcQuotaWallet
		if err = DeployService.AddBoost(&info); err != nil {
			log.Println("addBoost:", err.Error())
		}
	}

	if info.DeployStatus != define.DeployFail.Int() {
		if l.Id == uint64(info.ID) && l.Id != 0 && info.LanIp == l.LanIp && info.LanPort == l.LanPort && info.NetworkType == int(l.NetworkType) {
			if info.NetworkType == 0 {
				if info.InternetIp == l.InternetIp && info.InternetPort == l.InternetPort {
					response.Ok(c)
					return
				}
			} else {
				response.Ok(c)
				return
			}
		}
	}

	gclient := global.GateWayClinets.GetGateWayClinet(l.GateId)
	if gclient == nil {
		response.FailWithMessage(l.GateId+"not exist！", c)
		return
	}
	_, err = gclient.AddBoost(context.Background(), &l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// GetBoost
// @Tags      GetBoost
// @Summary   Get boost Info
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetOpID
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/getBoost [get]
func (deploy *DeployApi) GetBoost(c *gin.Context) {
	var l request.GetOpID
	err := c.ShouldBindQuery(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.OpIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := DeployService.GetBoostByOpId(l.OpId)
	if err != nil && err != gorm.ErrRecordNotFound {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// EditRunHost
// @Tags      EditRunHost
// @Summary    Modify the host service status
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request1.EditRun
// @Success   200   {object}  response.Ok{code=0,data="",msg=string}
// @Router    /deploy/editRunHost [post]
func (deploy *DeployApi) EditRunHost(c *gin.Context) {
	var ws []request1.EditRun
	err := c.ShouldBindJSON(&ws)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, v := range ws {
		gclient := global.GateWayClinets.GetGateWayClinet(v.GateId)
		if gclient == nil {
			response.FailWithMessage(v.GateId+" not exist！", c)
			return
		}
		_, err := gclient.RunStopService(context.Background(), &pb.RunStop{IsRun: v.IsRun, Id: v.Id, OpId: v.OpId, ServiceType: int32(v.ServiceType)})
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	response.Ok(c)
}

// ResetWorker
// @Tags      ResetWorker
// @Summary   Reset worker
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request1.ResetWorker
// @Success   200   {object}  response.Ok{code=0,data="",msg=string}
// @Router    /deploy/ResetWorker [post]
func (deploy *DeployApi) ResetWorker(c *gin.Context) {
	var ws []request1.ResetWorker
	err := c.ShouldBindJSON(&ws)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var wait sync.WaitGroup
	var errmsg string
	for _, vs := range ws {
		go func(v request1.ResetWorker) {
			wait.Add(1)
			defer wait.Done()
			w, _ := DeployService.GetWorker(v.Id)
			if w.MinerId == v.LinkId {
				errmsg = fmt.Sprintf("worker %d Connect the miner unchanged！", v.Id)
				return
			}
			gclient := global.GateWayClinets.GetGateWayClinet(v.GateId)
			if gclient == nil {
				log.Println(v.GateId + " not exist！")
				return
			}
			if !v.IsClear {
				//Stop P1
				if err = dispatchService.OnOffByWorkerID(v.Id, false); err != nil {
					log.Println(err.Error())
					return
				}
				//Take effect immediately
				_, err = gclient.SetWorkerTask(context.Background(), &pb.WorkerConfig{Id: v.Id, PreCount1: 0, PreCount2: -1, OpId: v.OpId})
				if err != nil {
					log.Println(err.Error())
					return
				}
				//abort AP P1
				_, err = gclient.SealingAbort(context.Background(), &pb.ResetWorker{Id: v.Id, Ip: w.Ip, LinkId: v.LinkId, OpId: v.OpId})
				if err != nil {
					log.Println(err.Error())
					return
				}
				return
			}
			//Stop of service
			go func() {
				_, err = gclient.RunStopService(context.Background(), &pb.RunStop{IsRun: false, Id: v.Id, OpId: v.OpId, ServiceType: int32(define.ServiceWorkerTask)})
				if err != nil {
					log.Println(err.Error())
					return
				}
			}()
			//Force clear cache
			_, err = gclient.ClearWorker(context.Background(), &pb.RequestOp{GateId: v.GateId, OpId: v.OpId})
			if err != nil {
				log.Println(err.Error())
				return
			}
			//Start to serve
			var hosts = []*pb.HostParam{
				{OpId: v.OpId, Id: v.Id, Ip: w.Ip},
			}
			_, err = gclient.AddWorker(context.Background(), &pb.BatchWorker{MinerId: v.LinkId, GateId: v.GateId, Host: hosts})
			if err != nil {
				log.Println(err.Error())
				return
			}
		}(vs)

	}
	if errmsg == "" {
		response.Ok(c)
		return
	}
	response.FailWithMessage(errmsg, c)
}

// GetWorkerList
// @Tags      GetWorkerList
// @Summary   Get worker list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/getWorkerList [post]
func (deploy *DeployApi) GetWorkerList(c *gin.Context) {
	var pageInfo request1.WorkerInfoPage
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := DeployService.GetWorkerList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetStorageList
// @Tags      GetStorageList
// @Summary   Get storage list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/getStorageList [get]
func (deploy *DeployApi) GetStorageList(c *gin.Context) {
	var pageInfo request1.WorkerInfoPage
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := DeployService.GetStorageList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// RoomAllLotus
// @Tags      RoomAllLotus
// @Summary   lotus machine Room
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/roomAllLotus [post]
func (deploy *DeployApi) RoomAllLotus(c *gin.Context) {
	var req request1.LotusInfoPage
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.GatewayVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := DeployService.GetRoomAllLotus(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// RelationMinerList
// @Tags      RelationMinerList
// @Summary   Get the associated miner list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request1.IDActor
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/relationMinerList [post]
func (deploy *DeployApi) RelationMinerList(c *gin.Context) {
	var id request1.IDActor
	err := c.ShouldBindJSON(&id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := DeployService.GetRelationMinerList(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// RelationWorkerList
// @Tags      RelationWorkerList
// @Summary   Gets the list of associated workers
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request1.IDActor
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/relationWorkerList [post]
func (deploy *DeployApi) RelationWorkerList(c *gin.Context) {
	var worker request1.IDActor
	err := c.ShouldBindJSON(&worker)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := DeployService.GetRelationWorkerList(worker)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// RelationStorageList
// @Tags      RelationStorageList
// @Summary   Get the associated storage list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request1.IDActor
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /deploy/relationStorageList [post]
func (deploy *DeployApi) RelationStorageList(c *gin.Context) {
	var worker request1.IDActor
	err := c.ShouldBindJSON(&worker)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := DeployService.GetRelationStorageList(worker)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// WorkerSectors
// @Tags      WorkerSectors
// @Summary   Query sector list on the worker
// @accept    application/json
// @Produce   application/json
// @Param     data  body      []request.GateOpID
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /deploy/WorkerSectors [get]
func (deploy *DeployApi) WorkerSectors(c *gin.Context) {
	var gateOps []request.GateOpID
	err := c.ShouldBindQuery(&gateOps)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var wait sync.WaitGroup
	var sectors = make([][]*pb.StatusCount, len(gateOps))
	for i, goid := range gateOps {
		go func(j int, id request.GateOpID) {
			wait.Add(1)
			defer wait.Done()
			gclient := global.GateWayClinets.GetGateWayClinet(id.GateId)
			if gclient == nil {
				response.FailWithMessage(id.GateId+" not exist！", c)
				return
			}
			sector, err := gclient.OpLocalSectors(context.Background(), &pb.OpMiner{OpId: id.OpId, Miner: id.Actor})
			if err != nil {
				response.FailWithMessage(id.GateId+" not exist！", c)
				return
			}
			sectors[j] = sector.Sectors
		}(i, goid)
	}
	wait.Wait()

	var statusCountMap = make(map[string]int32)
	var statusStoreMap = make(map[string]int32)
	var statusPreMsgMap = make(map[string]int32)
	var statusRemarkMap = make(map[string]string)

	for i := 0; i < len(sectors); i++ {
		for j := 0; j < len(sectors[i]); j++ {
			statCount := sectors[i][j]
			switch statCount.Status {
			case define.PreCommit1, define.AddPiece.String():
				statStr := fmt.Sprintf("%s,%s", define.PreCommit1, define.AddPiece.String())
				statusCountMap[statStr] += statCount.Count
				statusStoreMap[statStr] += statCount.Store
				statusPreMsgMap[statStr] += statCount.PreMsg
				statusRemarkMap[statStr] = fmt.Sprintf(`%d（已发过P2消息，需等待完成，否则终止会扣质押）
%d（正常P1）`, statusPreMsgMap[statStr], statusCountMap[statStr]-statusPreMsgMap[statStr])
			case define.PreCommit2, define.PreCommitting, define.SubmitPreCommitBatch:
				statStr := fmt.Sprintf("%s,%s,%s", define.PreCommit2, define.PreCommitting, define.SubmitPreCommitBatch)
				statusCountMap[statStr] += statCount.Count
				statusStoreMap[statStr] += statCount.Store
				statusPreMsgMap[statStr] += statCount.PreMsg
				statusRemarkMap[statStr] = fmt.Sprintf("P2等待发消息、批量推送消息")
			case define.WaitSeed.String(), define.PreCommitWait, define.PreCommitBatchWait:
				statStr := fmt.Sprintf("%s,%s,%s", define.WaitSeed.String(), define.PreCommitWait, define.PreCommitBatchWait)
				statusCountMap[statStr] += statCount.Count
				statusStoreMap[statStr] += statCount.Store
				statusPreMsgMap[statStr] += statCount.PreMsg
				statusRemarkMap[statStr] = fmt.Sprintf("等待150个高度，约1.5小时内，否则终止则会扣质押")
			case define.CommitFinalize, define.FinalizeSector.String():
				statStr := fmt.Sprintf("%s,%s", define.CommitFinalize, define.FinalizeSector.String())
				statusCountMap[statStr] += statCount.Count
				statusStoreMap[statStr] += statCount.Store
				statusPreMsgMap[statStr] += statCount.PreMsg
				statusRemarkMap[statStr] = fmt.Sprintf("%d（存储）%d（未存储）", statusStoreMap[statStr], statusCountMap[statStr]-statusStoreMap[statStr])
			case define.Committing, define.SubmitCommitAggregate:
				statStr := fmt.Sprintf("%s,%s", define.Committing, define.SubmitCommitAggregate)
				statusCountMap[statStr] += statCount.Count
				statusStoreMap[statStr] += statCount.Store
				statusPreMsgMap[statStr] += statCount.PreMsg
				statusRemarkMap[statStr] = fmt.Sprintf("C2等待发消息、批量推送消息")
			case define.SubmitCommit, define.CommitWait, define.CommitAggregateWait:
				statStr := fmt.Sprintf("%s,%s,%s", define.SubmitCommit, define.CommitWait, define.CommitAggregateWait)
				statusCountMap[statStr] += statCount.Count
				statusStoreMap[statStr] += statCount.Store
				statusPreMsgMap[statStr] += statCount.PreMsg
				statusRemarkMap[statStr] = fmt.Sprintf("消息等待C2消息推送")
			case define.Proving:
				statStr := fmt.Sprintf("%s", define.Proving)
				statusCountMap[statStr] += statCount.Count
				statusStoreMap[statStr] += statCount.Store
				statusPreMsgMap[statStr] += statCount.PreMsg
				statusRemarkMap[statStr] = fmt.Sprintf("%d（存储）%d（未存储）", statusStoreMap[statStr], statusCountMap[statStr]-statusStoreMap[statStr])
			default:
				statusCountMap[statCount.Status] += statCount.Count
				statusStoreMap[statCount.Status] += statCount.Store
				statusPreMsgMap[statCount.Status] += statCount.PreMsg
			}
		}
	}

	var res []response1.SectorStatusInfo
	for stat, v := range statusCountMap {
		statusInfo := response1.SectorStatusInfo{
			Status: stat,
			Count:  v,
			Store:  statusStoreMap[stat],
			PreMsg: statusPreMsgMap[stat],
			Remark: statusRemarkMap[stat],
		}
		res = append(res, statusInfo)
	}

	response.OkWithData(res, c)
}
