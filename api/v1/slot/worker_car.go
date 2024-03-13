package slot

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/define"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/lotus"
	"oplian/model/slot/request"
	"oplian/service"
	"oplian/utils"
)

type WorkerCarApi struct {
}

// WorkerCarList peter
// @Tags      WorkerCarList
// @Summary   workerCar任务列表
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.QueryWorkerCarReq "WorkerCar任务参数"
// @Success   200   {object}  response.Response{Device: device}  "WorkerCar任务列表"
// @Router    /workerCar/workerClusterList [GET]
func (w *WorkerCarApi) WorkerCarList(c *gin.Context) {

	var req request.QueryWorkerCarReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := service.ServiceGroupApp.SlotServiceGroup.GetWorkerCarTaskList(req)
	if err != nil {
		global.ZC_LOG.Error("SlotServiceGroup.GetWorkerCarTaskList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// ModifyTaskStatus peter
// @Tags      ModifyTaskStatus
// @Summary   修改workerCar任务状态
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.ModifyWorkerCarReq "修改WorkerCar任务状态参数"
// @Success   200   {object}  err  "修改WorkerCar任务状态结果"
// @Router    /workerCar/ModifyTaskStatus [POST]
func (w *WorkerCarApi) ModifyTaskStatus(c *gin.Context) {

	var req request.ModifyWorkerCarReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.WorkerCarStatusVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.ServiceGroupApp.SlotServiceGroup.ModifyTaskStatus(req)
	if err != nil {
		global.ZC_LOG.Error("SlotServiceGroup.ModifyTaskStatus:"+err.Error(), zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "获取成功", c)
}

// ModifyWorkerNum peter
// @Tags      ModifyWorkerNum
// @Summary   修改workerCar worker任务数
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.ModifyWorkerCarReq "修改workerCar worker任务数参数"
// @Success   200   {object}  err  "修改workerCar worker任务数结果"
// @Router    /workerCar/modifyWorkerCarNum [POST]
func (w *WorkerCarApi) ModifyWorkerNum(c *gin.Context) {

	var req request.ModifyWorkerCarReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.WorkerNumVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.ServiceGroupApp.SlotServiceGroup.ModifyWorkerNum(req)
	if err != nil {
		global.ZC_LOG.Error("SlotServiceGroup.ModifyWorkerCarNum:"+err.Error(), zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "获取成功", c)
}

// WorkerCarTaskDetail peter
// @Tags      WorkerCarTaskDetail
// @Summary   workerCar任务详情
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.QueryWorkerCarDetailReq "workerCar任务详情参数"
// @Success   200   {object}  map[string]interface{}, error  "workerCar任务详情结果"
// @Router    /workerCar/workerCarTaskDetail [GET]
func (w *WorkerCarApi) WorkerCarTaskDetail(c *gin.Context) {

	var req request.QueryWorkerCarDetailReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	res, err := service.ServiceGroupApp.SlotServiceGroup.GetWorkerCarTaskDetail(req)
	if err != nil {
		global.ZC_LOG.Error("SlotServiceGroup.ModifyWorkerCarNum:"+err.Error(), zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "获取成功", c)
}

// AddWorkerCarTask peter
// @Tags      AddWorkerCarTask
// @Summary   新增workerCar任务
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      lotus.WorkerCarTask "新增workerCar任务参数"
// @Success   200   {object}  error  "新增workerCar任务结果"
// @Router    /workerCar/addWorkerCarTask [POST]
func (w *WorkerCarApi) AddWorkerCarTask(c *gin.Context) {

	taskMain, err := service.ServiceGroupApp.SlotServiceGroup.CheckWorkerCarParam(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if taskMain.TaskType == define.CarTaskTypeAuto {
		_, err := service.ServiceGroupApp.SlotServiceGroup.AddWorkerCarTask(taskMain)
		if err != nil {
			global.ZC_LOG.Error("SlotServiceGroup.AddWorkerCarTask:"+err.Error(), zap.Error(err))
			response.FailWithMessage("获取失败:"+err.Error(), c)
			return
		}
	} else {

		taskDetail, err := service.ServiceGroupApp.SlotServiceGroup.CheckImportFile(c, taskMain.MinerId)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}

		id, err := service.ServiceGroupApp.SlotServiceGroup.AddWorkerCarTask(taskMain)
		if err != nil {
			global.ZC_LOG.Error("SlotServiceGroup.AddWorkerCarTask:"+err.Error(), zap.Error(err))
			response.FailWithMessage("获取失败:"+err.Error(), c)
			return
		}

		var paramMain []lotus.LotusSectorPiece
		for _, v := range taskDetail {
			v.QueueId = uint64(id)
			v.Actor = taskMain.MinerId
			paramMain = append(paramMain, v)
		}

		err = service.ServiceGroupApp.SlotServiceGroup.AddWorkerCarTaskDetail(taskMain.MinerId, paramMain)
		if err != nil {
			global.ZC_LOG.Error("SlotServiceGroup.AddWorkerCarTaskDetail:"+err.Error(), zap.Error(err))
			response.FailWithMessage("获取失败:"+err.Error(), c)
			return
		}
	}

	response.OkWithDetailed(nil, "获取成功", c)
}

// DcQuotaWalletList peter
// @Tags      DcQuotaWalletList
// @Summary   获取DC钱包
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      lotus.WorkerCarTask "获取DC钱包参数"
// @Success   200   {object}  error  "获取DC钱包结果"
// @Router    /workerCar/dcQuotaWalletList [GET]
func (w *WorkerCarApi) DcQuotaWalletList(c *gin.Context) {

	var req request.QueryWorkerCarReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	res, err := service.ServiceGroupApp.SlotServiceGroup.GetDcQuotaWalletList(req.MinerId)
	if err != nil {
		global.ZC_LOG.Error("SlotServiceGroup.GetDcQuotaWalletList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "获取成功", c)
}
