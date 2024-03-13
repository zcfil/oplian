package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"oplian/config"
	"oplian/global"
	"oplian/model/common/response"
	systemReq "oplian/model/system/request"
	response2 "oplian/model/system/response"
	"oplian/utils"
	"strconv"
)

type HomeInterfaceApi struct{}

// GetHostDataScreening
// @Tags      SysHostMonitorRecord
// @Summary   Overview of obtaining host data
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHomeInterface/getSysPatrolConfigList [get]
func (s *HomeInterfaceApi) GetHostDataScreening(c *gin.Context) {
	hostTestInfo := response2.DataScreeningReport{}

	// Obtain the total number of data centers
	machineRoomNum, err := machineRoomRecordService.GetSysMachineRoomNum()
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain the number of data centers!", zap.Error(err))
		response.FailWithMessage("Failed to obtain the number of data centers", c)
		return
	}
	hostTestInfo.MachineRoomNum = machineRoomNum

	// Get idle hosts
	freeHostNum, err := hostRecordService.GetFreeHostNum(false)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain the number of idle hosts", zap.Error(err))
		response.FailWithMessage("Failed to obtain the number of idle hosts", c)
		return
	}
	hostTestInfo.FreeHostNum = freeHostNum

	// Get the total host
	allHostNum, err := hostRecordService.GetFreeHostNum(true)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain the total number of hosts", zap.Error(err))
		response.FailWithMessage("Failed to obtain the total number of hosts", c)
		return
	}
	hostTestInfo.AllHostNum = allHostNum

	// Obtain the total number of alarms
	warnNum, err := warnManageService.GetWarnNumTotal()
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain the total number of alarms", zap.Error(err))
		response.FailWithMessage("Failed to obtain the total number of alarms", c)
		return
	}
	hostTestInfo.WarnNum = warnNum

	response.OkWithDetailed(hostTestInfo, "Successfully obtained", c)
}

// GetHostUseData
// @Tags      SysHostMonitorRecord
// @Summary   Get Host Usage Parameters
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHomeInterface/getHostUseData [get]
func (s *HomeInterfaceApi) GetHostUseData(c *gin.Context) {
	records, err := hostMonitorRecordService.GetLastHostMonitorRecords()
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain list information", zap.Error(err))
		response.FailWithMessage("Failed to obtain list information", c)
		return
	}
	// Memory Information
	var memoryAllSize int64
	var memoryUseSize int64

	var resp response2.HostUseDataResp
	for _, val := range records {
		resp.CPUUseRate += val.CPUUseRate
		memoryAllSize += val.MemorySize
		memoryUseSize += val.MemoryUseSize
		if val.DiskSize != "" {
			resp.DiskAllSize = utils.AddInputTwoSizes(resp.DiskAllSize, val.DiskSize)
		}
		if val.DiskUseSize != "" {
			resp.DiskUseSize = utils.AddInputTwoSizes(resp.DiskUseSize, val.DiskUseSize)
		}
	}

	// Unit conversion, size processing
	memoryAllSizeFloat := utils.DealSizeUnit(float64(memoryAllSize), "M")
	MemoryUseSizeFloat := utils.DealSizeUnit(float64(memoryUseSize), "M")
	if memoryAllSizeFloat > 1024 {
		memoryAllSizeFloat /= 1024
		resp.MemoryAllSize = strconv.FormatFloat(memoryAllSizeFloat, 'f', 2, 64) + "T"
	} else {
		resp.MemoryAllSize = strconv.FormatFloat(memoryAllSizeFloat, 'f', 2, 64) + "G"
	}
	if MemoryUseSizeFloat > 1024 {
		MemoryUseSizeFloat /= 1024
		resp.MemoryUseSize = strconv.FormatFloat(MemoryUseSizeFloat, 'f', 2, 64) + "T"
	} else {
		resp.MemoryUseSize = strconv.FormatFloat(MemoryUseSizeFloat, 'f', 2, 64) + "G"
	}

	// Get the number of list entries
	recordsNum := len(records)

	resp.CPUUseRate = resp.CPUUseRate / float32(recordsNum)
	resp.CPUUseNum = recordsNum
	resp.MemoryUseRate = math.Ceil((float64(memoryUseSize) / float64(memoryAllSize)) * 100)
	resp.DiskUseRate = float32(math.Ceil(utils.PercentageOfTwoSize(resp.DiskUseSize, resp.DiskAllSize) * 100))

	resp.CPUAllNum, err = hostRecordService.GetHostNum()
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain the total number of hosts in the list", zap.Error(err))
		response.FailWithMessage("Failed to obtain the total number of hosts in the list", c)
		return
	}

	response.OkWithDetailed(resp, "Successfully obtained", c)
}

// GetHostRunList
// @Tags      SysHostMonitorRecord
// @Summary   Get the running status of the host
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHomeInterface/getHostRunList [get]
func (s *HomeInterfaceApi) GetHostRunList(c *gin.Context) {
	var param systemReq.GetHostRunListReq
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Obtain the latest recorded patrol host information
	records, total, err := hostMonitorRecordService.GetLastHostMonitorLists(param)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain inspection information list information", zap.Error(err))
		response.FailWithMessage("Failed to obtain inspection information list information", c)
		return
	}
	// Query grouping information
	groupMap, err := hostGroupService.GetHostGroupListMap()
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain host group list information", zap.Error(err))
		response.FailWithMessage("Failed to obtain host group list information", c)
		return
	}
	// Processing return information
	var resp []response2.HostRunDataResp
	for _, val := range records {
		hostRunData := response2.HostRunDataResp{
			ID:             val.ID,
			HostUUID:       val.HostUUID,
			CPUUseRate:     val.CPUUseRate,
			CPUTemperature: val.CPUTemperature,
			MemoryUseRate:  val.MemoryUseRate,
			DiskUseRate:    val.DiskUseRate,
			CreatedAt:      val.CreatedAt.Format(config.TimeFormat),
		}
		// Get host information
		hostInfo, err := hostRecordService.GetSysHostRecord(val.HostUUID)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
			response.FailWithMessage("Failed to obtain host information", c)
			return
		}
		// Data supplementation
		hostRunData.HostName = hostInfo.HostName
		hostRunData.RoomId = hostInfo.RoomId
		hostRunData.RoomName = hostInfo.RoomName
		hostRunData.IntranetIP = hostInfo.IntranetIP
		hostRunData.InternetIP = hostInfo.InternetIP
		hostRunData.GroupName = groupMap[hostInfo.HostGroupId]

		resp = append(resp, hostRunData)
	}

	response.OkWithDetailed(response.PageResult{
		List:     resp,
		Total:    total,
		Page:     param.Page,
		PageSize: param.PageSize,
	}, "Successfully obtained", c)
}
