package lotus

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/lotus/request"
	"oplian/service"
	"oplian/utils"
)

type WorkerClusterApi struct {
}

// WorkerClusterList
// @Tags      WorkerClusterList
// @Summary   Gets the list of worker clusters
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.WorkerCluster
// @Success   200   {object}  response.Response{Device: device}
// @Router    /workerCluster/workerClusterList [post]
func (w *WorkerClusterApi) WorkerClusterList(c *gin.Context) {

	var req request.WorkerCluster
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := service.ServiceGroupApp.LotusServiceGroup.GetWorkerClusterList(req)
	if err != nil {
		global.ZC_LOG.Error("workerPlatform.GetWorkerClusterList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "Successfully obtained", c)
}

// AddWorkerOp
// @Tags      AddWorkerOp
// @Summary   Added worker op
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      lotus.LotusWorkerCluster
// @Success   200   {object}  response.Response{}
// @Router    /workerCluster/addWorkerOp [post]
func (w *WorkerClusterApi) AddWorkerOp(c *gin.Context) {

	var req request.AddWorkerCluster
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.ServiceGroupApp.LotusServiceGroup.AddWorkerOp(req)
	if err != nil {
		global.ZC_LOG.Error("AddWorkerOp:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "Successfully obtained", c)
}

// DelWorkerOp
// @Tags      DelWorkerOp
// @Summary   Delete worker op
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body     lotus.LotusWorkerCluster "删除worker ID"
// @Success   200   {object}  response.Response{} "成功或是失败"
// @Router    /workerCluster/addWorkerOp [post]
func (w *WorkerClusterApi) DelWorkerOp(c *gin.Context) {

	var req request.WorkerOp
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.ServiceGroupApp.LotusServiceGroup.DelWorkerOp(req)
	if err != nil {
		global.ZC_LOG.Error("DelWorkerOp:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "Successfully obtained", c)
}

// WorkerTaskDetail
// @Tags      WorkerTaskDetail
// @Summary   Worker task details
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      lotus.LotusWorkerCluster
// @Success   200   {object}  response.Response{Device: device}
// @Router    /workerCluster/workerTaskDetail [post]
func (w *WorkerClusterApi) WorkerTaskDetail(c *gin.Context) {

	var req request.WorkerOp
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := service.ServiceGroupApp.LotusServiceGroup.WorkerTaskDetail(req)
	if err != nil {
		global.ZC_LOG.Error("WorkerTaskDetail:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "Successfully obtained", c)
}

// ExportWorkerTask
// @Tags      ExportWorkerTask
// @Summary   Export the list of worker tasks
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      lotus.LotusWorkerCluster
// @Success   200   {object}  response.Response{Device: device}
// @Router    /workerCluster/workerTaskDetail [post]
func (w *WorkerClusterApi) ExportWorkerTask(c *gin.Context) {

	var req request.WorkerOp
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := service.ServiceGroupApp.LotusServiceGroup.ExportWorkerTask(req)
	if err != nil {
		global.ZC_LOG.Error("WorkerTaskDetail:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}

	title := []string{"taskId", "taskAgent", "sectorSize", "serverName", "taskStatus", "beginTime", "timeLength", "endTime"}
	col := []string{"taskId", "taskAgent", "sectorSize", "serverName", "taskStatus", "beginTime", "timeLength", "endTime"}
	utils.ExportExcelFile(c, title, col, result, "C2集群任务详情")
}
