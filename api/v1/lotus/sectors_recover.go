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

type SectorsRecoverApi struct {
}

// SectorsRecoverList
// @Tags      SectorsRecoverList
// @Summary   Gets the sector list
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.SectorsRecover
// @Success   200   {object}  response.Response{Device: device}
// @Router    /sectorsRecover/sectorsRecoverList [post]
func (s *SectorsRecoverApi) SectorsRecoverList(c *gin.Context) {

	var req request.SectorsRecover
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

	list, total, err := service.ServiceGroupApp.LotusServiceGroup.GetSectorsRecoverList(req)
	if err != nil {
		global.ZC_LOG.Error("sectorsRecover.GetSectorsRecoverList:"+err.Error(), zap.Error(err))
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

// AddSectorsRecoverTask
// @Tags      AddSectorsRecoverTask
// @Summary   Add the sector recovery task
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.SectorsRecover
// @Success   200   {object}  response.Response{Device: device}
// @Router    /sectorsRecover/addSectorsRecoverTask [post]
func (s *SectorsRecoverApi) AddSectorsRecoverTask(c *gin.Context) {

	var req request.SectorsRecoverTask
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

	res, err := service.ServiceGroupApp.LotusServiceGroup.AddSectorsRecoverTask(req)
	if err != nil {
		global.ZC_LOG.Error("sectorsRecover.AddSectorsRecoverTask:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// WorkerOpList
// @Tags      WorkerOpList
// @Summary   Get a list of worker hosts
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{res: res}  "worker机列表"
// @Router    /sectorsRecover/workerOpList [post]
func (s *SectorsRecoverApi) WorkerOpList(c *gin.Context) {

	var req request.WorkerInfo
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

	res, err := service.ServiceGroupApp.LotusServiceGroup.GetWorkerOpList(req)
	if err != nil {
		global.ZC_LOG.Error("sectorsRecover.GetWorkerOpList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// ModifySectorTaskStatus
// @Tags      ModifySectorTaskStatus
// @Summary   Changes the sector state
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      request.SectorStatus "扇区状态参数"
// @Success   200   {object}  response.Response{nil: }  ""
// @Router    /sectorsRecover/modifySectorTaskStatus [post]
func (s *SectorsRecoverApi) ModifySectorTaskStatus(c *gin.Context) {

	var req request.SectorStatus
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.SectorTaskStatusVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.ServiceGroupApp.LotusServiceGroup.ModifySectorTaskStatus(req)
	if err != nil {
		global.ZC_LOG.Error("sectorsRecover.ModifySectorTaskStatus:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "Successfully obtained", c)
}
