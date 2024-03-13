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
	response2 "oplian/model/system/response"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"strings"
	"time"
)

type HostTestRecordApi struct{}

// GetSysHostTestRecordList
// @Tags      SysHostTestRecord
// @Summary   Get Test Management List
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     re
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/getSysHostTestRecordList [get]
func (s *HostTestRecordApi) GetSysHostTestRecordList(c *gin.Context) {
	var pageInfo systemReq.SysHostTestRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := hostTestRecordService.GetSysHostTestRecordInfoList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	for i := 0; i < len(list); i++ {
		list[i].BeginTime = time.Unix(list[i].TestBeginAt, 0).Format(config.TimeFormat)
		// Calculate time difference
		if list[i].TestEndAt > 0 {
			list[i].EndTime = time.Unix(list[i].TestEndAt, 0).Format(config.TimeFormat)
			beginTime, _ := time.Parse(config.TimeFormat, list[i].BeginTime)
			endTime, _ := time.Parse(config.TimeFormat, list[i].EndTime)
			list[i].TestTakeTime = endTime.Sub(beginTime).String()
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
			list[i].IntranetIP = hostInfo.IntranetIP
			list[i].DeviceSN = hostInfo.DeviceSN
			list[i].AssetNumber = hostInfo.AssetNumber
		}
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

// GetSysHostTestReport
// @Tags      SysHostTestRecord
// @Summary   Obtain test management single piece information test report
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/getSysHostTestReport [get]
func (s *HostTestRecordApi) GetSysHostTestReport(c *gin.Context) {
	var hostTestReq systemReq.GetHostTestReportReq
	err := c.ShouldBindQuery(&hostTestReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	testInfo, err := hostTestRecordService.GetSysHostTestReport(hostTestReq)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	if testInfo.TestResult == config.HostUnderTest || testInfo.TestResult == config.HostInClose {
		response.FailWithMessage("Test information is running or closed", c)
		return
	}

	// Get host information
	hostInfo, err := hostRecordService.GetSysHostRecord(hostTestReq.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	hostTestInfo := response2.SysHostTestReport{
		ID:               testInfo.ID,
		HostUUID:         hostInfo.UUID,
		AssetNumber:      hostInfo.AssetNumber,
		DeviceSN:         hostInfo.DeviceSN,
		IntranetIP:       hostInfo.IntranetIP,
		HostName:         hostInfo.HostName,
		RoomName:         hostInfo.RoomName,
		RoomId:           hostInfo.RoomId,
		TestBeginAt:      time.Unix(testInfo.TestBeginAt, 0).Format(config.TimeFormat),
		TestEndAt:        time.Unix(testInfo.TestEndAt, 0).Format(config.TimeFormat),
		TestResult:       testInfo.TestResult,
		TestType:         testInfo.TestType,
		CPUHardScore:     testInfo.CPUHardScore,
		GPUHardScore:     testInfo.GPUHardScore,
		MemoryHardScore:  testInfo.MemoryHardScore,
		DiskHardScore:    testInfo.DiskHardScore,
		NetTestInfo:      testInfo.NetTestInfo,
		NetTestScore:     testInfo.NetTestScore,
		GPUTestInfo:      testInfo.GPUTestInfo,
		GPUTestScore:     testInfo.GPUTestScore,
		DiskIO:           testInfo.DiskIO,
		DiskAllRate:      testInfo.DiskAllRate,
		DiskAllRateScore: testInfo.DiskAllRateScore,
		DiskNFSRate:      testInfo.DiskNFSRate,
		DiskNFSRateScore: testInfo.DiskNFSRateScore,
		DiskSSDRate:      testInfo.DiskSSDRate,
		DiskSSDRateScore: testInfo.DiskSSDRateScore,
		IsAddPower:       testInfo.IsAddPower,
		SelectHostUUIDs:  testInfo.SelectHostUUIDs,
	}

	if testInfo.TestEndAt > 0 {
		testBeginAt := time.Unix(testInfo.TestBeginAt, 0).Format(config.TimeFormat)
		testEndAt := time.Unix(testInfo.TestEndAt, 0).Format(config.TimeFormat)
		beginTime, _ := time.Parse(config.TimeFormat, testBeginAt)
		endTime, _ := time.Parse(config.TimeFormat, testEndAt)
		hostTestInfo.TestTakeTime = endTime.Sub(beginTime).String()
	}

	// CPU information processing
	opCPUInfo := &utils.OpCPUInfo{}
	_ = json.Unmarshal([]byte(testInfo.CPUHardInfo), opCPUInfo)
	hostTestInfo.CPUHardInfo = *opCPUInfo
	// GPU information processing
	opGPUInfo := &[]utils.OpGPUInfo{}
	_ = json.Unmarshal([]byte(testInfo.GPUHardInfo), opGPUInfo)
	hostTestInfo.GPUHardInfo = *opGPUInfo
	// RAM information processing
	opRamInfo := &utils.OpRamInfo{}
	_ = json.Unmarshal([]byte(testInfo.MemoryHardInfo), opRamInfo)
	hostTestInfo.MemoryHardInfo = *opRamInfo
	// disk information processing
	hostDisk := &utils.HostDisk{}
	_ = json.Unmarshal([]byte(testInfo.DiskHardInfo), hostDisk)
	hostTestInfo.DiskHardInfo = *hostDisk

	response.OkWithDetailed(hostTestInfo, "Successfully obtained", c)
}

// AddHostTestByHand
// @Tags      SysHostTestRecord
// @Summary   Manually adding a new host for testing
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/addHostTestByHand [post]
func (s *HostTestRecordApi) AddHostTestByHand(c *gin.Context) {
	var hostTestReq systemReq.AddHostTestByHandReq
	err := c.ShouldBindJSON(&hostTestReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Determine if the selected number of hosts is qualified
	if hostTestReq.IsAddPower {
		if len(hostTestReq.HostUUIDs) != config.NetTestHostNumAddPower {
			response.FailWithMessage("The required number of auxiliary testing hosts for the newly added computing power type is"+strconv.Itoa(config.NetTestHostNumAddPower), c)
			return
		}
	} else {
		if len(hostTestReq.HostUUIDs) != config.NetTestHostNum {
			response.FailWithMessage("The required number of auxiliary testing hosts for non newly added computing power types is"+strconv.Itoa(config.NetTestHostNum), c)
			return
		}
	}

	// Get host information
	info, err := hostRecordService.GetSysHostRecord(hostTestReq.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	// Connect Gateway
	client := global.GateWayClinets.GetGateWayClinet(info.GatewayId)
	// Determine if the gateway connection is correct
	if client == nil {
		response.FailWithMessage("Failed to connect to the gateway corresponding to this host", c)
		return
	}
	// Determine whether the testing host is online
	_, err = client.CheckOpIsOnline(context.TODO(), &pb.OpHostUUID{HostUUID: info.UUID})
	if err != nil {
		response.FailWithMessage("The tested host is not online, please check its information", c)
		return
	}

	// Handling host uuids and IPs
	hostUUIDs := ""
	hostIntranetIPs := ""
	for _, val := range hostTestReq.HostUUIDs {
		hostUUIDs += val.HostUUID + ","
		hostIntranetIPs += val.IntranetIP + ","

		// Add 10 minutes of lock time to the corresponding selected auxiliary testing host
		sysHostTestRecord := system.SysHostRecord{
			UUID:          val.HostUUID,
			NetOccupyTime: time.Now().Add(10 * time.Minute).Unix(),
		}
		err = hostRecordService.UpdateHostNetOccupyTime(&sysHostTestRecord)
		if err != nil {
			global.ZC_LOG.Error("Failed to modify host information!", zap.Error(err))
			response.FailWithMessage("Modification failed", c)
			return
		}
	}

	// Check if the host is undergoing host testing
	_, err = hostTestRecordService.GetSysHostTestInfoByTestResult(hostTestReq.HostUUID, config.HostUnderTest)
	if err != nil && err != gorm.ErrRecordNotFound {
		global.ZC_LOG.Error("Obtain failed test information!", zap.Error(err))
		response.FailWithMessage("Obtain failed test information", c)
		return
	} else if err == nil {
		response.FailWithMessage("The host already has type information under testing. Please try again later", c)
		return
	}

	// Check if the host has already undergone this type of test
	testInfo, err := hostTestRecordService.GetSysHostTestInfoByType(hostTestReq.HostUUID, hostTestReq.TestType)
	if err != nil && err != gorm.ErrRecordNotFound {
		global.ZC_LOG.Error("Obtain failed test information!", zap.Error(err))
		response.FailWithMessage("Obtain failed test information", c)
		return
	}

	// time-on
	testBeginAt := time.Now().Unix()

	// Processing host information
	sysHostTestRecord := system.SysHostTestRecord{
		ZC_MODEL:        global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		HostUUID:        hostTestReq.HostUUID,
		TestType:        hostTestReq.TestType,
		TestBeginAt:     testBeginAt,
		TestMode:        config.ManualTrigger,
		TestEndAt:       0,
		TestResult:      config.HostUnderTest,
		IsAddPower:      hostTestReq.IsAddPower,
		SelectHostUUIDs: hostUUIDs[:len(hostUUIDs)-1],
		SelectHostIPs:   hostIntranetIPs[:len(hostIntranetIPs)-1],
	}
	if testInfo.ID == 0 {
		// Insert Host Test Information
		err = hostTestRecordService.CreateSysHostTestRecord(sysHostTestRecord)
		if err != nil {
			global.ZC_LOG.Error("Creation failed!", zap.Error(err))
			response.FailWithMessage("Creation failed", c)
			return
		}
	} else {
		// Modify the original host testing information
		err = hostTestRecordService.UpdateSysHostTestRecord(&sysHostTestRecord)
		if err != nil {
			global.ZC_LOG.Error("Creation failed!", zap.Error(err))
			response.FailWithMessage("Creation failed", c)
			return
		}
	}

	go func(info system.SysHostRecord) {
		if client != nil {
			_, err := client.OpInformationTest(context.TODO(),
				&pb.HostTestInfo{
					TestType:    hostTestReq.TestType,
					HostUUID:    info.UUID,
					TestMode:    int64(config.ManualTrigger),
					TestBeginAt: testBeginAt,
					IsAddPower:  hostTestReq.IsAddPower,
					HostUUIDs:   hostUUIDs[:len(hostUUIDs)-1],
					HostIPs:     hostIntranetIPs[:len(hostIntranetIPs)-1],
				})
			if err != nil {
				global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
				return
			}
		}
	}(info)

	response.OkWithMessage("New successfully added", c)
}

// CloseHostTest
// @Tags      SysHostTestRecord
// @Summary   Manually shut down host testing
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/closeHostTest [post]
func (s *HostTestRecordApi) CloseHostTest(c *gin.Context) {
	var hostTestReq systemReq.CloseHostTestReq
	err := c.ShouldBindJSON(&hostTestReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Get host information
	info, err := hostRecordService.GetSysHostRecord(hostTestReq.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	// Obtain host testing information
	testInfo, err := hostTestRecordService.GetSysHostTestInfo(hostTestReq.ID, hostTestReq.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	// Processing host testing information
	sysHostTestRecord := system.SysHostTestRecord{
		ZC_MODEL:   global.ZC_MODEL{ID: uint(hostTestReq.ID)},
		HostUUID:   hostTestReq.HostUUID,
		TestResult: config.HostInClose,
		TestEndAt:  time.Now().Unix(),
	}
	// Modify test information status
	err = hostTestRecordService.UpdateSysHostTestResult(&sysHostTestRecord)
	if err != nil {
		global.ZC_LOG.Error("Failed to modify test information!", zap.Error(err))
		response.FailWithMessage("Failed to shut down host testing", c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(info.GatewayId)

	go func(info system.SysHostRecord) {
		if client != nil {
			_, err := client.CloseOpInformationTest(context.TODO(),
				&pb.CloseHostTest{
					ID:              int64(hostTestReq.ID),
					HostUUID:        info.UUID,
					TestType:        testInfo.TestType,
					SelectHostUUIDs: testInfo.SelectHostUUIDs,
				})
			if err != nil {
				global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
				return
			}
		}
	}(info)

	response.OkWithMessage("Operation successful", c)
}

// RestartAddHostTest
// @Tags      SysHostTestRecord
// @Summary   Perform host testing again
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/restartAddHostTest [POST]
func (s *HostTestRecordApi) RestartAddHostTest(c *gin.Context) {
	var hostTestReq systemReq.AddHostTestRepeatReq
	err := c.ShouldBindJSON(&hostTestReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Get host information
	info, err := hostRecordService.GetSysHostRecord(hostTestReq.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	// Check if the host is undergoing host testing
	_, err = hostTestRecordService.GetSysHostTestInfoByTestResult(hostTestReq.HostUUID, config.HostUnderTest)
	if err != nil && err != gorm.ErrRecordNotFound {
		global.ZC_LOG.Error("Obtain failed test information!", zap.Error(err))
		response.FailWithMessage("Obtain failed test information", c)
		return
	} else if err == nil {
		response.FailWithMessage("The host already has type information under testing. Please try again later", c)
		return
	}

	// Obtain host testing information
	testInfo, err := hostTestRecordService.GetSysHostTestInfo(hostTestReq.ID, hostTestReq.HostUUID)
	if err != nil {
		global.ZC_LOG.Error("Failed to obtain host information!", zap.Error(err))
		response.FailWithMessage("Failed to obtain host information", c)
		return
	}

	hostUUIDs := strings.Split(testInfo.SelectHostUUIDs, ",")
	for _, val := range hostUUIDs {
		// Add 10 minutes of lock time to the corresponding selected auxiliary testing host
		sysHostTestRecord := system.SysHostRecord{
			UUID:          val,
			NetOccupyTime: time.Now().Add(10 * time.Minute).Unix(),
		}
		// 修改原有的主机信息
		err = hostRecordService.UpdateHostNetOccupyTime(&sysHostTestRecord)
		if err != nil {
			global.ZC_LOG.Error("Failed to modify host information!", zap.Error(err))
			response.FailWithMessage("Modification failed", c)
			return
		}
	}

	// Start Time
	testBeginAt := time.Now().Unix()

	// Processing host information
	sysHostTestRecord := system.SysHostTestRecord{
		ZC_MODEL:        global.ZC_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		HostUUID:        hostTestReq.HostUUID,
		TestType:        testInfo.TestType,
		TestBeginAt:     testBeginAt,
		TestMode:        config.ManualTrigger,
		TestEndAt:       0,
		TestResult:      config.HostUnderTest,
		IsAddPower:      testInfo.IsAddPower,
		SelectHostUUIDs: testInfo.SelectHostUUIDs,
		SelectHostIPs:   testInfo.SelectHostIPs,
	}

	// Modify the original host testing information
	err = hostTestRecordService.UpdateSysHostTestRecord(&sysHostTestRecord)
	if err != nil {
		global.ZC_LOG.Error("Creation failed!", zap.Error(err))
		response.FailWithMessage("Creation failed", c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(info.GatewayId)

	go func(info system.SysHostRecord) {
		if client != nil {
			_, err := client.OpInformationTest(context.TODO(),
				&pb.HostTestInfo{
					TestType:    testInfo.TestType,
					HostUUID:    info.UUID,
					TestMode:    int64(config.ManualTrigger),
					TestBeginAt: testBeginAt,
					IsAddPower:  testInfo.IsAddPower,
					HostUUIDs:   testInfo.SelectHostUUIDs,
					HostIPs:     testInfo.SelectHostIPs,
				})
			if err != nil {
				global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
				return
			}
		}
	}(info)

	response.OkWithMessage("Operation successful", c)
}

// DefaultHostTestInfo
// @Tags      SysHostTestRecord
// @Summary   Default test information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /sysHostRecord/defaultHostTestInfo [POST]
func (s *HostTestRecordApi) DefaultHostTestInfo(c *gin.Context) {
	var hostTestReq systemReq.GetDefaultHostTestInfoReq
	err := c.ShouldBindQuery(&hostTestReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	hostTestInfo := response2.DefaultHostTestReport{}

	if hostTestReq.TestType == config.HostMinerTest {
		hostTestInfo = response2.DefaultHostTestReport{
			CPUThreads:    config.CPUThreads,
			GPUModel:      config.GPUModelStandard,
			RamTotalMB:    strconv.Itoa(config.NodeRamSize) + "G",
			RamTotalMBAdd: strconv.Itoa(config.NodeRamSizeAdditional) + "G",
			DiskSize:      config.NodeSSDDiskSize,
			NetSpeed:      config.NetworkSpeed,
			NetSpeedAdd:   config.NetworkSpeedAdd,
			GPURunTime:    config.GPURunTime,
			DiskIO:        false,
			SSDDiskRate:   config.NodeSSDDiskSpeed,
			AllDiskRate:   "",
		}
	} else if hostTestReq.TestType == config.HostWorkerTest {
		hostTestInfo = response2.DefaultHostTestReport{
			CPUThreads:    config.CPUThreads,
			GPUModel:      config.GPUModelStandard,
			RamTotalMB:    strconv.Itoa(config.WorkerRamSize) + "G",
			RamTotalMBAdd: strconv.Itoa(config.WorkerRamSize) + "G",
			DiskSize:      config.WorkerSSDDiskSize,
			NetSpeed:      config.NetworkSpeed,
			NetSpeedAdd:   config.NetworkSpeed,
			GPURunTime:    config.GPURunTime,
			DiskIO:        false,
			SSDDiskRate:   config.WorkerSSDDiskSpeed,
			AllDiskRate:   "",
		}
	} else if hostTestReq.TestType == config.HostStorageTest {
		hostTestInfo = response2.DefaultHostTestReport{
			CPUThreads:    config.CPUThreads,
			RamTotalMB:    strconv.Itoa(config.StorageRamSize) + "G",
			RamTotalMBAdd: strconv.Itoa(config.StorageRamSize) + "G",
			DiskSize:      config.StorageDiskSize,
			NetSpeed:      config.NetworkSpeed,
			NetSpeedAdd:   config.NetworkSpeedAdd,
			DiskIO:        false,
			AllDiskRate:   config.StorageDiskOverallSpeed,
		}
	} else if hostTestReq.TestType == config.HostC2WorkerTest {
		hostTestInfo = response2.DefaultHostTestReport{
			CPUThreads: config.CPUThreads,
			GPUModel:   config.GPUModelStandard,
			RamTotalMB: strconv.Itoa(config.C2WorkerRamSize) + "G",
			NetSpeed:   config.NetworkC2WorkerSpeed,
			GPURunTime: config.GPURunTime,
		}
	}

	response.OkWithDetailed(hostTestInfo, "Successfully obtained", c)
}
