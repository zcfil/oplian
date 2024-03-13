package lotus

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"oplian/define"
	"oplian/global"
	request1 "oplian/model/common/request"
	"oplian/model/common/response"
	model "oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/service"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DispatchApi struct{}

// WorkerConfigList
// @Tags      WorkerConfigList
// @Summary   Gets the worker configuration list
// @accept    application/json
// @Produce   application/json
// @Param     data  body     request.RoomPageInfo
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/workerConfigList [post]
func (dispatch *DispatchApi) WorkerConfigList(c *gin.Context) {
	var pageInfo request.RoomPageInfo
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
	list, total, err := service.ServiceGroupApp.LotusServiceGroup.DispatchService.GetWorkerConfigList(pageInfo)
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

// SetWorkerConfig
// @Tags      SetWorkerConfig
// @Summary   Set worker Settings
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.WorkerPre
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/setWorkerConfig [post]
func (dispatch *DispatchApi) SetWorkerConfig(c *gin.Context) {
	var pres []request.WorkerPre
	err := c.ShouldBindJSON(&pres)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, pre := range pres {
		err = utils.Verify(pre, utils.IdVerify)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	for _, pre := range pres {
		if err = service.ServiceGroupApp.LotusServiceGroup.DispatchService.SetConfig(request.PreConfig{ID: pre.ID, PreCount1: pre.PreCount1, PreCount2: pre.PreCount2}); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		gclient := global.GateWayClinets.GetGateWayClinet(pre.GateId)
		if gclient == nil {
			continue
		}
		_, err = gclient.SetWorkerTask(context.Background(), &pb.WorkerConfig{Id: pre.ID, PreCount1: pre.PreCount1, PreCount2: pre.PreCount2, OpId: pre.OpId})
		if err != nil {
			continue
		}
	}
	response.Ok(c)
}

// OnOffPre1
// @Tags      OnOffPre1
// @Summary    Start stop task
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.WorkerOnOff
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/onOffPre1 [post]
func (dispatch *DispatchApi) OnOffPre1(c *gin.Context) {
	var pres []request.WorkerOnOff
	err := c.ShouldBindJSON(&pres)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, pre := range pres {
		err = utils.Verify(pre, utils.IdVerify)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	for _, pre := range pres {
		if err = service.ServiceGroupApp.LotusServiceGroup.DispatchService.OnOff(pre.ID, pre.OnOff); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		gclient := global.GateWayClinets.GetGateWayClinet(pre.GateId)
		if gclient == nil {
			continue
		}

		_, err = gclient.SetWorkerTask(context.Background(), &pb.WorkerConfig{Id: pre.ID, PreCount1: 0, PreCount2: -1, OpId: pre.OpId})
		if err != nil {
			continue
		}
	}
	response.Ok(c)
}

// GetSectorsList
// @Tags      GetSectorsList
// @Summary   Gets a list of sectors
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SectorPage
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/getSectorsList [get]
func (dispatch *DispatchApi) GetSectorsList(c *gin.Context) {
	var param request.SectorPage
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.ActorVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, total, err := service.ServiceGroupApp.LotusServiceGroup.GetSectorsList(param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     res,
		Total:    total,
		Page:     param.Page,
		PageSize: param.PageSize,
	}, "Successfully obtained", c)
}

// GetSectorDetails
// @Tags      GetSectorDetails
// @Summary   Get the sector details
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SectorPage
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/getSectorDetails [get]
func (dispatch *DispatchApi) GetSectorDetails(c *gin.Context) {
	var param request.SectorPage
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.ActorSectorIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, total, err := service.ServiceGroupApp.LotusServiceGroup.GetSectorDetails(param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	miner, err := service.ServiceGroupApp.SystemServiceGroup.GetMinerManage(param.Actor)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	gclinet := global.GateWayClinets.GetGateWayClinet(miner.GateId)
	var paths []string
	if gclinet != nil {
		actorID, _ := strconv.ParseUint(param.Actor[2:], 10, 64)
		spath, err := gclinet.SectorStorage(context.Background(), &pb.SectorActorID{Token: miner.Token, Ip: miner.Ip, Miner: actorID, Number: param.Number})
		if err == nil && spath != nil {
			for _, v := range spath.Strs {
				paths = append(paths, v.Value)
			}
		} else {
			global.ZC_LOG.Error(err.Error())
		}
	}
	result := make(map[string]interface{})
	result["piece"] = res.Piece
	result["sectorLog"] = res.SectorLog
	result["sectorInfo"] = res.SectorInfo

	result["sectorPaths"] = paths
	result["total"] = total
	result["page"] = param.Page
	result["pageSize"] = param.PageSize

	response.OkWithDetailed(result, "Successfully obtained", c)
}

// GetSectorTaskList
// @Tags      GetSectorTaskList
// @Summary   Get the sector task list
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SectorPage
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/getSectorTaskList [get]
func (dispatch *DispatchApi) GetSectorTaskList(c *gin.Context) {
	var param request1.PageInfo
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, total, err := service.ServiceGroupApp.LotusServiceGroup.GetSectorTaskList(param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     res,
		Total:    total,
		Page:     param.Page,
		PageSize: param.PageSize,
	}, "Successfully obtained", c)
}

// AddSectorTask
// @Tags      AddSectorTask
// @Summary   Add Task Queue
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SectorPage
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/addSectorTask [put]
func (dispatch *DispatchApi) AddSectorTask(c *gin.Context) {
	var param request.SectorQueue
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.SectorTaskVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	queue := &model.LotusSectorQueue{
		TaskName:   param.TaskName,
		SectorSize: param.SectorSize,
		JobTotal:   param.SectorTotal,
		SectorType: param.SectorType,
		Actor:      param.Actor,
	}
	if err = service.ServiceGroupApp.LotusServiceGroup.AddSectorTaskQueue(queue); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.Ok(c)
}

// AddSectorDcTask
// @Tags      AddSectorDcTask
// @Summary   Add DC task queue
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/addSectorDcTask [put]
func (dispatch *DispatchApi) AddSectorDcTask(c *gin.Context) {
	s := c.Request.FormValue("sectorSize")
	sectorSize, err := strconv.Atoi(s)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	concurrent, _ := strconv.Atoi(c.Request.FormValue("concurrent"))
	dir := c.Request.FormValue("carDir")
	if dir == "" {
		dir = c.Request.FormValue("path")
	}
	actor := c.Request.FormValue("actor")
	opIds := utils.FormDataStrToArray(c.Request.FormValue("carIds"))
	gid := c.Request.FormValue("gateId")
	file, fileHead, err := c.Request.FormFile("dcInfo")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer file.Close()

	buf := make([]byte, fileHead.Size)
	file.Read(buf)

	str := utils.TrimBlankSpace(string(buf))
	deals := strings.Split(str, ";")
	if len(deals) == 0 {
		response.FailWithMessage("Incorrect format1，Please follow：uuid,pieceCid,Height of maturity;", c)
		return
	}
	format := strings.Split(deals[0], ",")
	if len(format) < 3 {
		response.FailWithMessage("Incorrect format2，Please follow：uuid,pieceCid,Height of maturity;", c)
		return
	}
	if _, err = uuid.FromString(format[0]); err != nil {
		response.FailWithMessage("Incorrect format3，Please follow：uuid,pieceCid,Height of maturity;", c)
		return
	}
	if !strings.HasPrefix(format[1], "baga6ea4seaq") {
		response.FailWithMessage("Incorrect format4，Please follow：uuid,pieceCid,Height of maturity;", c)
		return
	}
	exp, err := strconv.ParseUint(format[2], 10, 64)
	if err != nil {
		response.FailWithMessage("Incorrect format5，Please follow：uuid,pieceCid,Height of maturity;", c)
		return
	}
	if exp+120 <= utils.BlockHeight() {
		response.FailWithMessage("Order has expired.！", c)
		return
	}
	//保存任务队列数据
	queue := &model.LotusSectorQueue{
		TaskName:         c.Request.FormValue("taskName"),
		SectorSize:       sectorSize,
		Actor:            actor,
		SectorType:       define.SectorTypeDC,
		JobTotal:         0,
		TaskStatus:       define.QueueStatusAnalyzing,
		ConcurrentImport: concurrent,
	}
	queue.CreatedAt = time.Now()
	if err = service.ServiceGroupApp.LotusServiceGroup.AddSectorTaskQueue(queue); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	go func() {
		dealMap := make(map[string]*request.DealInfo) //dealuuid -- [0]pieceCid [1]endTtime [2]carPath
		type carPathUuid struct {
			car *pb.CarFile
			uid string
		}
		var carFiles = make([]*carPathUuid, len(deals))
		var waitDeal sync.WaitGroup
		for i, v := range deals {

			dealstr := strings.Split(v, ",")
			if len(dealstr) < 3 {
				continue
			}

			exp, err := strconv.ParseUint(format[2], 10, 64)
			if err != nil {
				global.ZC_LOG.Error(fmt.Sprintf("Incorrect format：%s,%s,%s", format[0], format[1], format[2]))
				continue
			}
			if exp+120 <= utils.BlockHeight() {
				global.ZC_LOG.Error(fmt.Sprintf("Order has expired：%s,%s,%s", format[0], format[1], format[2]))
				continue
			}
			epoch, _ := strconv.ParseInt(dealstr[2], 10, 64)
			deal := &request.DealInfo{
				DealUuid:  dealstr[0],
				PieceCid:  dealstr[1],
				EndEpoch:  epoch,
				JobStatus: define.QueueSectorStatusFitFail,
			}
			dealMap[dealstr[0]] = deal
			// match car files
			waitDeal.Add(1)
			go func(j int, uuid string) {
				defer waitDeal.Done()
				var waitCar sync.WaitGroup
				for _, id := range opIds {
					gclient := global.GateWayClinets.GetGateWayClinet(gid)
					if gclient == nil {
						continue
					}
					waitCar.Add(1)
					go func(opId string) {
						defer waitCar.Done()

						res, err := gclient.CarFilePath(context.Background(), &pb.CarFile{OpId: opId, Path: dir, FileName: dealstr[1] + ".car"})
						if err == nil {
							ucar := &carPathUuid{
								car: res,
								uid: uuid,
							}
							carFiles[j] = ucar
						}

					}(id)
					waitCar.Wait()
				}
			}(i, dealstr[0])
		}
		waitDeal.Wait()
		defer func() {

			queue.JobTotal = len(dealMap)
			queue.TaskStatus = define.QueueStatusRun
			if err != nil {
				queue.TaskStatus = define.QueueStatusAnalyzingFail
			}
			if err = service.ServiceGroupApp.LotusServiceGroup.UpdateSectorTaskQueue(queue); err != nil {
				global.ZC_LOG.Error(err.Error())
			}
		}()

		for _, ucar := range carFiles {
			if ucar != nil {
				dealMap[ucar.uid].CarPath = ucar.car.Path
				dealMap[ucar.uid].FileOpId = ucar.car.OpId
				dealMap[ucar.uid].JobStatus = define.QueueSectorStatusWait
			}
		}

		if err = service.ServiceGroupApp.LotusServiceGroup.AddDCSectorQueueDetail(dealMap, actor, queue.ID); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}()

	response.Ok(c)
}

// EditTaskQueueStatus
// @Tags      EditTaskQueueStatus
// @Summary   Modify the task queue
// @accept    application/json
// @Produce   application/json
// @Param     data  body      param request1.IdStatus
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/editTaskQueueStatus [put]
func (dispatch *DispatchApi) EditTaskQueueStatus(c *gin.Context) {
	var param request1.IdStatus
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.IdStatusVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = service.ServiceGroupApp.LotusServiceGroup.EditTaskQueueStatus(param.ID, param.Status); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// EditTaskQueueDetailStatus
// @Tags      EditTaskQueueDetailStatus
// @Summary   Modify the task detail state
// @accept    application/json
// @Produce   application/json
// @Param     data  body      param request1.IdStatus
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/editTaskQueueDetailStatus [put]
func (dispatch *DispatchApi) EditTaskQueueDetailStatus(c *gin.Context) {
	var params []request1.ActorIdStatus
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, param := range params {
		err = utils.Verify(param, utils.ActorIdStatusVerify)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	if err = service.ServiceGroupApp.LotusServiceGroup.EditTaskQueueDetailStatusBatch(params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// EditTaskQueueDetailConcurrent
// @Tags      EditTaskQueueDetailConcurrent
// @Summary   Modify the task queue number state
// @accept    application/json
// @Produce   application/json
// @Param     data  body      param request1.IdCount
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/editTaskQueueDetailConcurrent [put]
func (dispatch *DispatchApi) EditTaskQueueDetailConcurrent(c *gin.Context) {
	var param request1.IdCount
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.IdCountVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = service.ServiceGroupApp.LotusServiceGroup.EditTaskQueueDetailConcurrent(param); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// GetManageMiners
// @Tags      GetManageMiners
// @Summary   Get the schedule list
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SectorPage
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/getManageMiners [get]
func (dispatch *DispatchApi) GetManageMiners(c *gin.Context) {
	var param request1.PageInfo
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := service.ServiceGroupApp.LotusServiceGroup.GetManageMiners(param.Keyword)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// GetSectorTaskDetailList
// @Tags      GetSectorTaskDetailList
// @Summary   Get the sector mission details list
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.DealPage
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/getSectorTaskDetailList [get]
func (dispatch *DispatchApi) GetSectorTaskDetailList(c *gin.Context) {
	var param request.DealPage
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.DealPageVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, total, err := service.ServiceGroupApp.LotusServiceGroup.GetSectorQueueDetailList(param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     res,
		Total:    total,
		Page:     param.Page,
		PageSize: param.PageSize,
	}, "Successfully obtained", c)
}

// SectorRecoverDetail
// @Tags      SectorRecoverDetail
// @Summary   Sector recovery task details
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SectorRecoverDetail
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/sectorRecoverDetail [get]
func (dispatch *DispatchApi) SectorRecoverDetail(c *gin.Context) {
	var param request.SectorRecoverDetail
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.RecoverDetailVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	task, list, total, err := service.ServiceGroupApp.LotusServiceGroup.GetSectorRecoverDetail(param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		Main:     task,
		List:     list,
		Total:    total,
		Page:     param.Page,
		PageSize: param.PageSize,
	}, "Successfully obtained", c)
}

// CheckDealCar
// @Tags      CheckDealCar
// @Summary   Matching order
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SectorPage
// @Success   200   {object}  response.Response{data=model.LotusInfo,msg=string}
// @Router    /dispatch/checkDealCar [get]
func (dispatch *DispatchApi) CheckDealCar(c *gin.Context) {
	var deals []request.DealCarInfo
	err := c.ShouldBindJSON(&deals)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	for _, v := range deals {
		err = utils.Verify(v, utils.DealCarVerify)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	go func() {
		var carFiles = make([]*pb.CarFile, len(deals))
		var waitDeal sync.WaitGroup
		for i, deal := range deals {
			// match car files
			waitDeal.Add(1)
			go func(j int) {
				defer waitDeal.Done()
				var waitCar sync.WaitGroup
				for _, id := range deal.OpIds {
					gclient := global.GateWayClinets.GetGateWayClinet(deal.GateId)
					if gclient == nil {
						continue
					}
					waitCar.Add(1)
					go func(opId string) {
						defer waitCar.Done()
						res, err := gclient.CarFilePath(context.Background(), &pb.CarFile{OpId: opId, Path: deal.CarDir, FileName: deal.PieceCid + ".car"})
						if err == nil {
							carFiles[j] = res
						}
					}(id)
					waitCar.Wait()
				}
			}(i)
		}
		waitDeal.Wait()

		for i, car := range carFiles {
			if car != nil {
				param := request1.ActorIdStatus{
					ID:     deals[i].ID,
					Status: define.QueueSectorStatusWait,
					Actor:  deals[i].Actor,
					Value:  car.Path,
				}
				if err = service.ServiceGroupApp.LotusServiceGroup.EditTaskQueueDetailStatus(param, car.OpId); err != nil {
					global.ZC_LOG.Error(err.Error())
				}
			}
		}
	}()

	response.Ok(c)
}
