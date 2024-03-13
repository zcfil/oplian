package system

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	request1 "oplian/model/lotus/request"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	systemResp "oplian/model/system/response"
	"oplian/service/pb"
	"oplian/utils"
)

type HostRecordApi struct{}

// CreateSysHostRecord
// @Tags      SysHostRecord
// @Summary   create SysHostRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysHostRecord/createSysHostRecord [post]
func (s *HostRecordApi) CreateSysHostRecord(c *gin.Context) {
	var sysHostRecord system.SysHostRecord
	err := c.ShouldBindJSON(&sysHostRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = hostRecordService.CreateSysHostRecord(sysHostRecord)
	if err != nil {
		global.ZC_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// UpdateSysHostRecord
// @Tags      SysHostRecord
// @Summary   update SysHostRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysHostRecord/updateSysHostRecord [post]
func (s *HostRecordApi) UpdateSysHostRecord(c *gin.Context) {
	var sysHostRecord system.SysHostRecord
	err := c.ShouldBindJSON(&sysHostRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(sysHostRecord, utils.HostRecordVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = hostRecordService.UpdateSysHostRecord(&sysHostRecord)
	if err != nil {
		if err == global.GatewayIdMismatchError {
			global.ZC_LOG.Error("Modification failed!", zap.Error(err))
			response.FailWithMessage("The gateway ID corresponding to the host does not match the gateway ID of the data center,", c)
			return
		}
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// DeleteSysHostRecordByUUIDs
// @Tags      SysHostRecord
// @Summary   Batch deletion of SysHostRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysHostRecord/deleteSysHostRecordByUUIDs [delete]
func (s *HostRecordApi) DeleteSysHostRecordByUUIDs(c *gin.Context) {
	var IDS request.UUIDsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = hostRecordService.DeleteSysHostRecordByUUIDs(IDS)
	if err != nil {
		global.ZC_LOG.Error("Batch deletion failed!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed", c)
		return
	}
	response.OkWithMessage("Batch deletion successful", c)
}

// GetSysHostRecordList
// @Tags      SysHostRecord
// @Summary   Paging to retrieve the SysHostRecord list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/getSysHostRecordList [get]
func (s *HostRecordApi) GetSysHostRecordList(c *gin.Context) {
	var hostSearchReq systemReq.SysHostRecordSearch
	err := c.ShouldBindQuery(&hostSearchReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := hostRecordService.GetSysHostRecordInfoList(hostSearchReq)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	for i := 0; i < len(list); i++ {
		// 获取分组名称
		if list[i].HostGroupId != 0 {
			hostGroupInfo, err := hostGroupService.GetSysHostGroup(list[i].HostGroupId)
			if err != nil {
				global.ZC_LOG.Error("Failed to obtain host grouping information!", zap.Error(err))
				response.FailWithMessage("Failed to obtain host grouping information", c)
				return
			}
			list[i].HostGroupName = hostGroupInfo.GroupName
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     hostSearchReq.Page,
		PageSize: hostSearchReq.PageSize,
	}, "Successfully obtained", c)
}

// GetSysHostList
// @Tags      SysHost
// @Summary   Paging to retrieve the SysHost list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/getSysHostist [get]
func (s *HostRecordApi) GetSysHostList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := hostRecordService.GetSysHostRecordList(pageInfo)
	if err != nil {
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetSysHostTestRecordList
// @Tags      SysHost
// @Summary   Get a list of hosts
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     r
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/hostTestRecordList [get]
func (s *HostRecordApi) GetSysHostTestRecordList(c *gin.Context) {
	var pageInfo request1.HostPage
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.ClassifyVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := hostRecordService.GetSysHostTestRecordList1(pageInfo)
	if err != nil {
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetSysHostListNormal
// @Tags      SysHost
// @Summary   Get a list of running hosts
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/getSysHostListNormal [get]
func (s *HostRecordApi) GetSysHostListNormal(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, totalNum, err := hostRecordService.GetSysHostRecordListNormal(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	response.OkWithDetailed(systemResp.SysHostPageResult{
		List:     list,
		TotalNum: totalNum,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// OpInfoList
// @Tags      SysHost
// @Summary   Get host information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{res}
// @Router    /sysHostRecord/opInfo [post]
func (s *HostRecordApi) OpInfoList(c *gin.Context) {
	var param systemReq.SysHostRecordListByClassifyReq
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := hostRecordService.GetOpInfoList(param)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	response.OkWithDetailed(res, "Successfully obtained", c)
}

// GetOpHardwareInfo
// @Tags      SysHost
// @Summary   Obtain Op hardware information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{res}
// @Router    /sysHostRecord/getOpHardwareInfo [get]
func (s *HostRecordApi) GetOpHardwareInfo(c *gin.Context) {
	var param systemReq.GetOpHardwareInfoReq
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := hostRecordService.GetSysHostRecord(param.UUID)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(info.GatewayId)

	if client == nil {
		response.FailWithMessage("Connection gateway failed", c)
	}

	res, err := client.GetOpHardwareInfo(context.TODO(),
		&pb.OpHardwareInfo{
			HostUUID:     info.UUID,
			HostClassify: int64(info.HostClassify),
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	// 解析返回的参数
	opStorageInfo := &utils.OpStorageInfo{}
	err = json.Unmarshal([]byte(res.GetValue()), opStorageInfo)
	if err != nil {
		global.ZC_LOG.Error("Failed to parse host information,hostUUID: "+info.UUID, zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	response.OkWithDetailed(opStorageInfo, "Successfully obtained", c)
}

// GetNetHostList
// @Tags      SysHost
// @Summary   Get a list of selected network hosts
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{res}
// @Router    /sysHostRecord/netHostList [get]
func (s *HostRecordApi) GetNetHostList(c *gin.Context) {
	var param systemReq.GetNetHostListReq
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := hostRecordService.GetSysHostRecord(param.UUID)
	if err != nil {
		global.ZC_LOG.Error("Failed to parse host information!", zap.Error(err))
		response.FailWithMessage("Failed to parse host information", c)
		return
	}
	res, err := hostRecordService.GetNetHostList(param, info.GatewayId)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	response.OkWithDetailed(res, "Successfully obtained", c)
}

// GetPatrolHostList
// @Tags      SysHost
// @Summary   Obtain the list of inspection target hosts
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{res}
// @Router    /sysHostRecord/patrolHostList [get]
func (s *HostRecordApi) GetPatrolHostList(c *gin.Context) {
	var param systemReq.GetPatrolHostListReq
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := hostRecordService.GetPatrolHostList(param)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	response.OkWithDetailed(res, "Successfully obtained", c)
}
