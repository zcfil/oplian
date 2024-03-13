package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system/request"
	"oplian/utils"
)

type JobPlatformApi struct {
	IsStop bool // Is the program terminated
}

// ExecuteScript
// @Tags      ScriptInfo
// @Summary   Execute script
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Device: device}
// @Router    /jobPlatform/executeScript [post]
func (a *JobPlatformApi) ExecuteScript(c *gin.Context) {

	var jobReq request.JobReq
	err := c.ShouldBindJSON(&jobReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(jobReq, utils.ExecuteScriptVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	jobReq.UserName = utils.GetUserName(c)
	err = jobPlatform.ExecuteScript(jobReq)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.ExecuteScript:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "Successfully obtained", c)
}

// ExecuteRecordsDetail
// @Tags      ExecuteRecordsDetail
// @Summary   Execution script recording details
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Device: device}
// @Router    /jobPlatform/executeRecordsDetail [post]
func (a *JobPlatformApi) ExecuteRecordsDetail(c *gin.Context) {

	var param request.CommonParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(param, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := jobPlatform.ExecuteRecordsDetail(param.ID)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.ExecuteRecordsDetail"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// ExecuteRecordsList
// @Tags      ExaCustomer
// @Summary   Execution History List
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /jobPlatform/executeRecordsList [post]
func (a *JobPlatformApi) ExecuteRecordsList(c *gin.Context) {

	var req request.ExecuteRecordsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.PateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := jobPlatform.GetExecuteRecordsList(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.GetExecuteRecordsList"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "Successfully obtained", c)
}

// BatchUploadFile
// @Tags      BatchUploadFile
// @Summary   Batch upload of files
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/batchUploadFile [post]
func (a *JobPlatformApi) BatchUploadFile(c *gin.Context) {

	res, err := jobPlatform.BatchUploadFile(c)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.BatchUploadFile:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// OpFileToGateWay
// @Tags      OpFileToGateWay
// @Summary   Host file synchronization to gateway
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/opFileToGateWay [post]
func (a *JobPlatformApi) OpFileToGateWay(c *gin.Context) {

	var req request.OpFileSync
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
	res, err := jobPlatform.OpFileToGateWay(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.OpFileToGateWay:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// FileDistribution
// @Tags      FileDistribution
// @Summary   File Distribution
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/fileDistribution [post]
func (a *JobPlatformApi) FileDistribution(c *gin.Context) {

	var distribute request.DistributeReq
	err := c.ShouldBindJSON(&distribute)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(distribute, utils.DistributionVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	distribute.UserName = utils.GetUserName(c)
	err = jobPlatform.FileDistribution(distribute)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.FileDistribution:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "Successfully obtained", c)
}

// FileForcedToStop
// @Tags      FileForcedToStop
// @Summary   Forced termination of file distribution
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/forcedToStop [post]
func (a *JobPlatformApi) FileForcedToStop(c *gin.Context) {

	var distribute request.DistributeReq
	err := c.ShouldBindJSON(&distribute)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(distribute, utils.ForcedToStopVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = jobPlatform.ForcedToStop(distribute)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.ForcedToStop:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "Successfully obtained", c)
}

// ExecuteScriptStop
// @Tags      ForcedToStop
// @Summary   Script execution forcibly terminates
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Data: fb}
// @Router    /jobPlatform/executeScriptStop [post]
func (a *JobPlatformApi) ExecuteScriptStop(c *gin.Context) {

	var req request.JobReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(req, utils.ScriptStopVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = jobPlatform.ExecuteScriptStop(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.ExecuteScriptStop:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(nil, "Successfully obtained", c)
}

// SaveFileHost
// @Tags      SaveFileHost
// @Summary   Set up file host
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Device: device}
// @Router    /jobPlatform/delFileManage [post]
func (a *JobPlatformApi) SaveFileHost(c *gin.Context) {

	var req request.AddFileReq
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

	err = jobPlatform.SaveFileHost(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.SaveFileHost:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}

// FileHostList
// @Tags      FileHostList
// @Summary   List of file hosts
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{res: res}
// @Router    /jobPlatform/fileHostList [post]
func (a *JobPlatformApi) FileHostList(c *gin.Context) {

	res, err := jobPlatform.GetFileHostList()
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.GetFileHostList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// FileManageList
// @Tags      FileManageList
// @Summary   File Management List
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /jobPlatform/fileManageList [post]
func (a *JobPlatformApi) FileManageList(c *gin.Context) {

	var req request.FileManageReq
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
	list, total, err := jobPlatform.GetFileManageList(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.GetFileManageList"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "Successfully obtained", c)
}

// DelFileManage
// @Tags      DelFileManage
// @Summary   Delete file management
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{Device: device}
// @Router    /jobPlatform/delFileManage [post]
func (a *JobPlatformApi) DelFileManage(c *gin.Context) {

	var req request.FileManageReq
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

	err = jobPlatform.DelFileManage(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.DelFileManage:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}

// ModifyFileManage
// @Tags      ModifyFileManage
// @Summary   Update file management
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{ok}
// @Router    /jobPlatform/modifyFileManage [post]
func (a *JobPlatformApi) ModifyFileManage(c *gin.Context) {

	var req request.FileManageReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.ModifyFileVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = jobPlatform.ModifyFileManage(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.ModifyFileManage:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}

// AddFile
// @Tags      AddFile
// @Summary   New file
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{ok}
// @Router    /jobPlatform/addFile [post]
func (a *JobPlatformApi) AddFile(c *gin.Context) {

	var req request.AddFileReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.AddFileVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = jobPlatform.AddFileGrpc(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.AddFile:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}

// SysFilePoint
// @Tags      SysFilePoint
// @Summary   OP point-to-point replication (supports one or multiple files)
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{ok}
// @Router    /jobPlatform/sysFilePoint [post]
func (a *JobPlatformApi) SysFilePoint(c *gin.Context) {

	var req request.SynFileReq
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

	err = jobPlatform.SysFilePoint(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.SysFilePoint:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}

// FileListByType
// @Tags      FileListByType
// @Summary   File Management List
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /jobPlatform/fileListByType [post]
func (a *JobPlatformApi) FileListByType(c *gin.Context) {

	var req request.FileTypeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.RequesFileTypeReqVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := jobPlatform.GetFileListByType(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.GetFileListByType"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithData(list, c)
}

// LotusHeightList
// @Tags      LotusHeightList
// @Summary   Lotus height file list
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{res: res}
// @Router    /jobPlatform/lotusHeightList [post]
func (a *JobPlatformApi) LotusHeightList(c *gin.Context) {

	var req request.GateWayReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.RequesFileTypeReqVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := jobPlatform.GetLotusHeightList(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.GetLotusHeightList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// MinerList
// @Tags      MinerList
// @Summary   List of Miner Hosts
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{res: res}
// @Router    /jobPlatform/lotusHeightList [post]
func (a *JobPlatformApi) MinerList(c *gin.Context) {

	var req request.GateWayReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.RequesFileTypeReqVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := jobPlatform.GetMinerList(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.GetLotusHeightList:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "Successfully obtained", c)
}

// DownLoadFile
// @Tags      DownLoadFile
// @Summary   File download (supports one or multiple file downloads)
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{ok}
// @Router    /jobPlatform/downLoadFile [post]
func (a *JobPlatformApi) DownLoadFile(c *gin.Context) {

	var req request.DownLoadReq
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

	err = jobPlatform.DownLoadFile(req)
	if err != nil {
		global.ZC_LOG.Error("jobPlatform.DownLoadFile:"+err.Error(), zap.Error(err))
		response.FailWithMessage("Acquisition failed:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("ok", "Successfully obtained", c)
}
