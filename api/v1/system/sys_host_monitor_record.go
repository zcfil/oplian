package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"oplian/config"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	systemReq "oplian/model/system/request"
	response2 "oplian/model/system/response"
	"oplian/utils"
	"strconv"
	"time"
)

type HostMonitorRecordApi struct{}

// GetSysHostMonitorRecordList
// @Tags      SysHostMonitorRecord
// @Summary   Obtain a list of host specific hardware usage rates
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostMonitorRecord/getSysHostMonitorRecordList [get]
func (s *HostMonitorRecordApi) GetSysHostMonitorRecordList(c *gin.Context) {
	var pageInfo systemReq.SysHostMonitorRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var beginTime, endTime string
	beginTime = time.Now().Add(-time.Minute * 5).Format(config.TimeFormat)
	endTime = utils.GetNowStr()

	list, total, err := hostMonitorRecordService.GetSysHostMonitorRecordInfoList(pageInfo, beginTime, endTime)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	for i := 0; i < len(list); i++ {
		// 获取主机信息
		if list[i].HostUUID != "" {
			hostInfo, err := hostRecordService.GetSysHostRecord(list[i].HostUUID)
			if err != nil {
				global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
				response.FailWithMessage("Failed to obtain host information", c)
				return
			}
			list[i].HostName = hostInfo.HostName
			list[i].IntranetIP = hostInfo.IntranetIP
			list[i].InternetIP = hostInfo.InternetIP
		}
	}

	countResp := map[string]interface{}{}
	if pageInfo.Keyword == "" {
		allNum, err := hostRecordService.GetHostNum()
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain the total number of hosts in the list", zap.Error(err))
			response.FailWithMessage("Failed to obtain the total number of hosts in the list", c)
			return
		}
		countResp["allNum"] = allNum
		countResp["normalNum"] = total
	}

	response.OkWithDetailed(response.PageResult{
		Main:     countResp,
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetSysHostMonitorLineChart
// @Tags      SysHostMonitorRecord
// @Summary   Get a line chart of host usage
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostMonitorRecord/getSysHostMonitorLineChart [post]
func (s *HostMonitorRecordApi) GetSysHostMonitorLineChart(c *gin.Context) {
	var req systemReq.HostUUIDsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var beginTime, endTime string
	beginTime = time.Now().Add(-time.Minute * 30).Format(config.TimeFormat)
	endTime = utils.GetNowStr()

	var resp []response2.SysHostMonitorLineChart
	for _, val := range req.HostUUIDs {
		// Get host information
		hostInfo, err := hostRecordService.GetSysHostRecord(val.HostUUID)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
			response.FailWithMessage("Failed to obtain host information", c)
			return
		}
		// Obtain hardware monitoring information of the host within half an hour
		list, err := hostMonitorRecordService.GetSysHostMonitorRecordList(val.HostUUID, val.GPUID, req.Keyword, beginTime, endTime)
		if err != nil {
			global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
			response.FailWithMessage("Acquisition failed", c)
			return
		}
		var sysHostMonitorLineChart response2.SysHostMonitorLineChart
		sysHostMonitorLineChart.HostName = hostInfo.HostName
		sysHostMonitorLineChart.HostUUID = val.HostUUID
		sysHostMonitorLineChart.GPUID = val.GPUID
		sysHostMonitorLineChart.HostInfo = list
		resp = append(resp, sysHostMonitorLineChart)
	}

	response.OkWithDetailed(resp, "Successfully obtained", c)
}

// GetStorageInformationList
// @Tags      SysHostMonitorRecord
// @Summary   Storage Machine Storage Information List Information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostMonitorRecord/getStorageInfoList [get]
func (s *HostMonitorRecordApi) GetStorageInformationList(c *gin.Context) {
	var req request.PageInfo
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var beginTime, endTime string
	beginTime = time.Now().Add(-time.Minute * 5).Format(config.TimeFormat)
	endTime = utils.GetNowStr()

	list, total, err := hostMonitorRecordService.GetLastHostStorageInfoMonitorLists(req, beginTime, endTime)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	var resp []response2.StorageRateResp
	var diskAllSizeSum string
	var diskUseSizeSum string
	for _, val := range list {
		// 获取主机信息
		hostInfo, err := hostRecordService.GetSysHostRecord(val.HostUUID)
		if err != nil {
			global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
			response.FailWithMessage("Failed to obtain host information", c)
			return
		}
		// 处理返回信息
		allNum, _ := strconv.ParseFloat(utils.GetNumFromStr(val.DiskSize), 64)
		useNum, _ := strconv.ParseFloat(utils.GetNumFromStr(val.DiskUseSize), 64)
		storageInfo := response2.StorageRateResp{
			HostName:        hostInfo.HostName,
			IntranetIP:      hostInfo.IntranetIP,
			InternetIP:      hostInfo.InternetIP,
			DiskUseRate:     float32(math.Ceil(utils.PercentageOfTwoSize(val.DiskUseSize, val.DiskSize) * 100)),
			DiskAllSize:     allNum,
			DiskAllSizeUnit: utils.GetLetterFromStr(val.DiskSize),
			DiskUseSize:     useNum,
			DiskUseSizeUnit: utils.GetLetterFromStr(val.DiskUseSize),
			MinerID:         "--",
		}
		resp = append(resp, storageInfo)
		diskAllSizeSum = utils.AddInputTwoSizes(diskAllSizeSum, val.DiskSize)
		diskUseSizeSum = utils.AddInputTwoSizes(diskUseSizeSum, val.DiskUseSize)
	}
	allNum, _ := strconv.ParseFloat(utils.GetNumFromStr(diskAllSizeSum), 64)
	useNum, _ := strconv.ParseFloat(utils.GetNumFromStr(diskUseSizeSum), 64)
	diskInfo := map[string]interface{}{
		"diskAllSizeSum":  allNum,
		"diskAllSizeUnit": utils.GetLetterFromStr(diskAllSizeSum),
		"diskUseSizeSum":  useNum,
		"diskUseSizeUnit": utils.GetLetterFromStr(diskUseSizeSum),
		"diskUseRate":     math.Ceil(utils.PercentageOfTwoSize(diskUseSizeSum, diskAllSizeSum) * 100),
	}

	response.OkWithDetailed(response.PageResult{
		Main:     diskInfo,
		List:     resp,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "Successfully obtained", c)
}
