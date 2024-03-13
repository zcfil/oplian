package system

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	"oplian/utils"
)

type MachineRoomRecordApi struct{}

// CreateSysMachineRoomRecord
// @Tags      SysMachineRoomRecord
// @Summary   create SysMachineRoomRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysMachineRoomRecord/createSysMachineRoomRecord [post]
func (s *MachineRoomRecordApi) CreateSysMachineRoomRecord(c *gin.Context) {
	var sysMachineRoomRecord system.SysMachineRoomRecord
	err := c.ShouldBindJSON(&sysMachineRoomRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = machineRoomRecordService.CreateSysMachineRoomRecord(sysMachineRoomRecord)
	if err != nil {
		global.ZC_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// UpdateSysMachineRoomRecord
// @Tags      SysMachineRoomRecord
// @Summary   update SysMachineRoomRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysMachineRoomRecord/updateSysMachineRoomRecord [post]
func (s *MachineRoomRecordApi) UpdateSysMachineRoomRecord(c *gin.Context) {
	var sysMachineRoomRecord system.SysMachineRoomRecord
	err := c.ShouldBindJSON(&sysMachineRoomRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(sysMachineRoomRecord, utils.MachineRoomVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = machineRoomRecordService.UpdateSysMachineRoomRecord(&sysMachineRoomRecord)
	if err != nil {
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// DeleteSysMachineRoomRecordByRoomIds
// @Tags      SysMachineRoomRecord
// @Summary   Batch Delete SysMachineRoomRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysMachineRoomRecord/deleteSysMachineRoomRecordByIds [delete]
func (s *MachineRoomRecordApi) DeleteSysMachineRoomRecordByRoomIds(c *gin.Context) {
	var roomIds request.RoomIdsReq
	err := c.ShouldBindJSON(&roomIds)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Query the data center to be deleted. Below is the bound host
	for _, val := range roomIds.RoomIds {
		_, hostCount, err := hostRecordService.GetSysHostRecordCountByRoomId(val)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain the number of bound host entries!", zap.Error(err))
			response.FailWithMessage("Failed to obtain the number of bound host entries!", c)
			return
		}
		if hostCount > 0 {
			response.FailWithMessage("Select a host that is associated with the data center. Please remove the associated host before deleting it!", c)
			return
		}
	}

	err = machineRoomRecordService.DeleteSysMachineRoomRecordByIds(roomIds)
	if err != nil {
		global.ZC_LOG.Error("Batch deletion failed!", zap.Error(err))
		response.FailWithMessage("Batch deletion failed", c)
		return
	}
	response.OkWithMessage("Batch deletion successful", c)
}

// GetSysMachineRoomRecordList
// @Tags      SysMachineRoomRecord
// @Summary   Paging to retrieve the SysMachineRoomRecord list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysMachineRoomRecord/getSysMachineRoomRecordList [get]
func (s *MachineRoomRecordApi) GetSysMachineRoomRecordList(c *gin.Context) {
	var pageInfo systemReq.SysMachineRoomRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := machineRoomRecordService.GetSysMachineRoomRecordInfoList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	for i := 0; i < len(list); i++ {
		// 获取绑定主机信息
		data, hostCount, err := hostRecordService.GetSysHostRecordCountByRoomId(list[i].RoomId)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain the number of bound host entries", zap.Error(err))
			response.FailWithMessage("Acquisition failed", c)
			return
		}
		list[i].RoomBindHost.BindHostNum = hostCount
		list[i].RoomBindHost.BindHostUUIDs = data

		// 获取绑定节点机数量
		nodeHostNum, err := hostRecordService.GetNodeHostNumByRoomId(list[i].RoomId)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain the number of bound host entries", zap.Error(err))
			response.FailWithMessage("Failed to obtain the number of bound host entries", c)
			return
		}
		list[i].BindNodeHostNum = nodeHostNum

		// 获取机房负责人name
		if len(list[i].RoomLeader) != 0 {
			leaderUser, err := userService.GetUserInfoTidy(list[i].RoomLeader)
			if err != nil {
				global.ZC_LOG.Error("Failed to obtain the information of the data center manager!", zap.Error(err))
				response.FailWithMessage("Failed to obtain the information of the data center manager", c)
				return
			}
			list[i].RoomLeaderName = leaderUser.NickName
		}

		// 获取机房管理员name
		if len(list[i].RoomAdmin) != 0 {
			adminUser, err := userService.GetUserInfoTidy(list[i].RoomAdmin)
			if err != nil {
				global.ZC_LOG.Error("Failed to obtain server room administrator information!", zap.Error(err))
				response.FailWithMessage("Failed to obtain server room administrator information", c)
				return
			}
			list[i].RoomAdminName = adminUser.NickName
		}

		// 获取机房owner name
		if len(list[i].RoomOwner) != 0 {
			ownerUser, err := userService.GetUserInfoTidy(list[i].RoomOwner)
			if err != nil {
				global.ZC_LOG.Error("Failed to obtain data center owner information!", zap.Error(err))
				response.FailWithMessage("Failed to obtain data center owner information", c)
				return
			}
			list[i].RoomOwnerName = ownerUser.NickName
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// BindSysHostRecords
// @Tags      SysHostRecord
// @Summary   update SysHostRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysMachineRoomRecord/updateSysMachineRoomRecord [post]
func (s *MachineRoomRecordApi) BindSysHostRecords(c *gin.Context) {
	var bindSysHostRecordsReq systemReq.BindSysHostRecordsReq
	err := c.ShouldBindJSON(&bindSysHostRecordsReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(bindSysHostRecordsReq, utils.BindSysHostRecordsVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = hostRecordService.UpdateSysHostRecordRoomId(&bindSysHostRecordsReq)
	if err != nil {
		if err == global.GatewayIdMismatchError {
			global.ZC_LOG.Error("Modification failed!", zap.Error(err))
			response.FailWithMessage("The gateway ID corresponding to the host does not match the gateway ID of the data center", c)
			return
		}
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// GetMachineRoomId
// @Tags      SysMachineRoomRecord
// @Summary   Add a new data center to obtain the generated roomId
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=map[string]string{},msg=string}
// @Router    /sysMachineRoomRecord/getSysMachineRoomRecordList [get]
func (s *MachineRoomRecordApi) GetMachineRoomId(c *gin.Context) {
	resp := map[string]string{}
	resp["roomId"] = uuid.NewV4().String()
	response.OkWithDetailed(resp, "Successfully obtained", c)
}

// BindSysHostRecordsList
// @Tags      SysMachineRoomRecord
// @Summary   Get a list of bound and unbound hosts
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysMachineRoomRecord/bindSysHostRecordList [get]
func (s *MachineRoomRecordApi) BindSysHostRecordsList(c *gin.Context) {
	var bindReq systemReq.BindSysHostRecordListReq
	err := c.ShouldBindQuery(&bindReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := hostRecordService.GetSysHostRecordBindList(bindReq)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	response.OkWithDetailed(list, "Successfully obtained", c)
}

// UnbindSysHostRecords
// @Tags      SysHostRecord
// @Summary   Unbind host information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /sysMachineRoomRecord/updateSysMachineRoomRecord [post]
func (s *MachineRoomRecordApi) UnbindSysHostRecords(c *gin.Context) {
	var unBindSysHostRecordsReq systemReq.UnbindSysHostRecordsReq
	err := c.ShouldBindJSON(&unBindSysHostRecordsReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(unBindSysHostRecordsReq, utils.BindSysHostRecordsVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = hostRecordService.UpdateSysHostRecordRoomIdUnbind(&unBindSysHostRecordsReq)
	if err != nil {
		global.ZC_LOG.Error("Modification failed!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// GetRoomRecordList
// @Tags      SysMachineRoomRecord
// @Summary   Paging to retrieve the SysMachineRoomRecord list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysMachineRoomRecord/getSysMachineRoomRecordList [get]
func (s *MachineRoomRecordApi) GetRoomRecordList(c *gin.Context) {
	list, err := machineRoomRecordService.GetRoomRecordList()
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(list, "Successfully obtained", c)
}
