package system

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"oplian/config"
	"oplian/global"
	"oplian/model/common/response"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	"oplian/service/pb"
	"oplian/utils"
	"time"
)

type HostPatrolRecordApi struct{}

// GetSysHostPatrolRecordList
// @Tags      SysHostPatrolRecord
// @Summary   Paging to obtain the SysHostPatrolRecord list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/getSysHostPatrolRecordList [get]
func (s *HostPatrolRecordApi) GetSysHostPatrolRecordList(c *gin.Context) {
	var pageInfo systemReq.SysHostPatrolRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := hostPatrolRecordService.GetSysHostPatrolRecordInfoList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	for i := 0; i < len(list); i++ {
		list[i].BeginTime = time.Unix(list[i].PatrolBeginAt, 0).Format(config.TimeFormat)
		// Calculate time difference
		if list[i].PatrolEndAt > 0 {
			list[i].EndTime = time.Unix(list[i].PatrolEndAt, 0).Format(config.TimeFormat)
			beginTime, _ := time.Parse(config.TimeFormat, list[i].BeginTime)
			endTime, _ := time.Parse(config.TimeFormat, list[i].EndTime)
			list[i].PatrolTakeTime = endTime.Sub(beginTime).String()
		}

		// Get host information
		var hostInfo system.SysHostRecord
		if len(list[i].HostUUID) != 0 {
			hostInfo, err = hostRecordService.GetSysHostRecord(list[i].HostUUID)
			if err != nil {
				global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
				response.FailWithMessage("Failed to obtain host information", c)
				return
			}
			list[i].HostName = hostInfo.HostName
			list[i].InternetIP = hostInfo.InternetIP
			list[i].IntranetIP = hostInfo.IntranetIP
			list[i].DeviceSN = hostInfo.DeviceSN
			list[i].AssetNumber = hostInfo.AssetNumber
		}
		// Obtain the name of the computer room
		if hostInfo.RoomId != "" {
			list[i].RoomId = hostInfo.RoomId
			list[i].RoomName = hostInfo.RoomName
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetSysHostPatrolReport
// @Tags      SysHostPatrolRecord
// @Summary   Obtain inspection management single information inspection report
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/getSysHostPatrolReport [get]
func (s *HostPatrolRecordApi) GetSysHostPatrolReport(c *gin.Context) {
	var reportReq systemReq.GetHostPatrolReportReq
	err := c.ShouldBindQuery(&reportReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := hostPatrolRecordService.GetSysHostPatrolReport(reportReq)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	info.BeginTime = time.Unix(info.PatrolBeginAt, 0).Format(config.TimeFormat)
	// 计算时间差
	if info.PatrolEndAt > 0 {
		info.EndTime = time.Unix(info.PatrolEndAt, 0).Format(config.TimeFormat)
		beginTime, _ := time.Parse(config.TimeFormat, info.BeginTime)
		endTime, _ := time.Parse(config.TimeFormat, info.EndTime)
		info.PatrolTakeTime = endTime.Sub(beginTime).String()
	}

	// Get host information
	var hostInfo system.SysHostRecord
	if len(info.HostUUID) != 0 {
		hostInfo, err = hostRecordService.GetSysHostRecord(info.HostUUID)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
			response.FailWithMessage("Failed to obtain host information", c)
			return
		}
		info.HostName = hostInfo.HostName
		info.IntranetIP = hostInfo.IntranetIP
		info.DeviceSN = hostInfo.DeviceSN
		info.AssetNumber = hostInfo.AssetNumber
	}
	// Obtain the name of the computer room
	if hostInfo.RoomId != "" {
		info.RoomId = hostInfo.RoomId
		info.RoomName = hostInfo.RoomName
	}

	// Program version processing
	version := &utils.LotusPackageVersion{}
	if len(info.PackageVersion) > 0 {
		err = json.Unmarshal([]byte(info.PackageVersion), version)
		if err != nil {
			global.ZC_LOG.Error("Failed to parse host version information,hostUUID: "+hostInfo.UUID, zap.Error(err))
			response.FailWithMessage("Failed to obtain inspection information", c)
			return
		}
		info.HostPackageVersion = *version
	}

	response.OkWithDetailed(info, "Successfully obtained", c)
}

// AddHostPatrolByHand
// @Tags      SysHostPatrolRecord
// @Summary   Manual addition of host inspection
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/addHostPatrolByHand [post]
func (s *HostPatrolRecordApi) AddHostPatrolByHand(c *gin.Context) {
	var patrolReq systemReq.AddHostPatrolByHandReq
	err := c.ShouldBindJSON(&patrolReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Get host information
	info, err := hostRecordService.GetSysHostRecord(patrolReq.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	if info.HostClassify != config.HostC2WorkerType {
		if info.HostClassify != int(patrolReq.PatrolType) {
			response.FailWithMessage("The selection of inspection type does not match the host type", c)
			return
		}
	} else {
		info.HostClassify = config.HostWorkerType
	}

	// Check if the host is undergoing inspection
	patrolInfo, err := hostPatrolRecordService.GetSysHostPatrolByResult(patrolReq.HostUUID, config.HostUnderTest)
	if err != nil && err != gorm.ErrRecordNotFound {
		global.ZC_LOG.Error("Failed to obtain inspection information!", zap.Error(err))
		response.FailWithMessage("Failed to obtain inspection information", c)
		return
	} else if patrolInfo.ID != 0 {
		response.FailWithMessage("The host already has type information under inspection. Please try again later", c)
		return
	}

	var otherInfo system.SysHostRecord
	if info.HostClassify == config.HostStorageType {
		// Obtain an additional machine IP for ping testing
		otherInfo, err = hostRecordService.GetSysOtherHostRecord(info.UUID, info.GatewayId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				response.FailWithMessage("There is only one host under the current gateway, and network status cannot be tested", c)
				return
			}
			global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
			response.FailWithMessage("Acquisition failed", c)
		}
	}

	client := global.GateWayClinets.GetGateWayClinet(info.GatewayId)

	go func(info system.SysHostRecord) {
		if client != nil {
			_, err := client.OpInformationPatrol(context.TODO(),
				&pb.HostPatrolInfo{
					HostUUID:     info.UUID,
					HostClassify: int64(info.HostClassify),
					PatrolMode:   int64(config.ManualTrigger),
					PatrolHostIP: otherInfo.IntranetIP,
				})
			if err != nil {
				global.ZC_LOG.Error("client.OpInformationPatrol err: ", zap.Error(err))
				return
			}
		}
	}(info)

	response.OkWithMessage("New successfully added", c)
}
