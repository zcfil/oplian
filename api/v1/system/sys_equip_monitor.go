package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"oplian/api/v1/lotus"
	"oplian/config"
	"oplian/define"
	"oplian/global"
	"oplian/model/common/response"
	responseNodel "oplian/model/lotus/response"
	"oplian/model/system/request"
	"oplian/service/pb"
	"oplian/utils"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type EquipMonitorApi struct{}

// GetMinerList
// @Tags      EquipMonitorApi
// @Summary   Get a list of Miner machines
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getMinerList [get]
func (s *EquipMonitorApi) GetMinerList(c *gin.Context) {
	var pageInfo request.HostMonitorReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := lotus.DeployService.GetMinerMonitorList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	listInfoChan := make(chan *pb.MonitorInfo, len(list))

	// Query the corresponding hardware information for each machine
	for i := 0; i < len(list); i++ {
		go func(info responseNodel.MinerMonitorInfo) {
			resInfo := &pb.MonitorInfo{OpId: info.OpId}

			client := global.GateWayClinets.GetGateWayClinet(info.GateId)
			if client == nil {
				log.Println("Connection gateway failed: " + info.GateId)
				listInfoChan <- resInfo
				return
			}

			res, err := client.GetOpMonitorInfo(context.TODO(),
				&pb.OpHardwareInfo{
					HostUUID:     info.OpId,
					HostClassify: config.HostMinerType,
				})
			if err != nil {
				log.Println("client.OpInformationTest err: ", zap.Error(err))
				listInfoChan <- resInfo
				return
			}
			listInfoChan <- res
		}(list[i])
	}

	listInfo := make(map[string]*pb.MonitorInfo, len(list))

	for i := 0; i < len(list); i++ {
		info := <-listInfoChan
		listInfo[info.OpId] = info
	}

	for i := 0; i < len(list); i++ {
		res, isExit := listInfo[list[i].OpId]
		if !isExit {
			continue
		}
		list[i].LotusStatus = res.LotusStatus
		list[i].MinerStatus = res.MinerStatus
		list[i].BoostStatus = res.BoostStatus
		list[i].CPUUseRate = float64(res.CpuUseRate)
		list[i].SysDiskUseRate = res.SysDiskUseRate
		list[i].SysDiskLeave = res.SysDiskLeave
		list[i].MainDiskUseRate = res.MainDiskUseRate
		list[i].MainDiskLeave = res.MainDiskLeave
		list[i].DiskStatus = res.DiskStatus
		list[i].GPUStatus = res.GpuStatus
		list[i].MountStatus = res.MountStatus
		list[i].LotusHeightStatus = res.LotusHeightStatus
		// Obtain the storage machine corresponding to this node
		var mountList []responseNodel.StorageMountErrorList
		mountList, err = lotus.DeployService.GetStorageMountInfoByActor(list[i].Actor)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		var noMountNum int
		for _, val := range mountList {
			if !utils.IsInStrList(val.Ip, res.Ips) {
				noMountNum++
			}
		}
		if noMountNum == 0 {
			list[i].MountStatus = true
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetWorkerList
// @Tags      EquipMonitorApi
// @Summary   Get a list of Worker machines
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getWorkerList [get]
func (s *EquipMonitorApi) GetWorkerList(c *gin.Context) {
	var pageInfo request.HostMonitorReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := lotus.DeployService.GetWorkerMonitorList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	listInfoChan := make(chan *pb.MonitorInfo, len(list))

	// Query the corresponding hardware information for each machine
	for i := 0; i < len(list); i++ {
		go func(info responseNodel.WorkerMonitorInfo) {
			resInfo := &pb.MonitorInfo{OpId: info.OpId}

			client := global.GateWayClinets.GetGateWayClinet(info.GateId)
			if client == nil {
				log.Println("Connection gateway failed: " + info.GateId)
				listInfoChan <- resInfo
				return
			}

			res, err := client.GetOpMonitorInfo(context.TODO(),
				&pb.OpHardwareInfo{
					HostUUID:     info.OpId,
					HostClassify: config.HostWorkerType,
				})
			if err != nil {
				log.Println("client.OpInformationTest err: ", zap.Error(err))
				listInfoChan <- resInfo
				return
			}
			listInfoChan <- res
		}(list[i])
	}

	listInfo := make(map[string]*pb.MonitorInfo, len(list))

	for i := 0; i < len(list); i++ {
		info := <-listInfoChan
		listInfo[info.OpId] = info
	}

	for i := 0; i < len(list); i++ {
		res, isExit := listInfo[list[i].OpId]
		if !isExit {
			continue
		}
		list[i].WorkerStatus = res.WorkerStatus
		list[i].CPUUseRate = float64(res.CpuUseRate)
		list[i].SysDiskUseRate = res.SysDiskUseRate
		list[i].SysDiskLeave = res.SysDiskLeave
		list[i].MainDiskUseRate = res.MainDiskUseRate
		list[i].MainDiskLeave = res.MainDiskLeave
		list[i].DiskStatus = res.DiskStatus
		list[i].GPUStatus = res.GpuStatus
		list[i].MountStatus = res.MountStatus
		// Obtain the storage machine corresponding to this node
		var mountList []responseNodel.StorageMountErrorList
		mountList, err = lotus.DeployService.GetStorageMountInfoByActor(list[i].Actor)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		var noMountNum int
		for _, val := range mountList {
			if !utils.IsInStrList(val.Ip, res.Ips) {
				noMountNum++
			}
		}
		if noMountNum < 0 {
			list[i].MountStatus = true
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetStorageList
// @Tags      EquipMonitorApi
// @Summary   Get Storage Machine List
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getStorageList [get]
func (s *EquipMonitorApi) GetStorageList(c *gin.Context) {
	var pageInfo request.HostMonitorReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := lotus.DeployService.GetStorageMonitorList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	listInfoChan := make(chan *pb.MonitorInfo, len(list))

	// Query the corresponding hardware information for each machine
	for i := 0; i < len(list); i++ {
		go func(info responseNodel.StorageMonitorInfo) {
			resInfo := &pb.MonitorInfo{OpId: info.OpId}

			client := global.GateWayClinets.GetGateWayClinet(info.GateId)

			if client == nil {
				log.Println("Connection gateway failed: " + info.GateId)
				listInfoChan <- resInfo
				return
			}

			res, err := client.GetOpMonitorInfo(context.TODO(),
				&pb.OpHardwareInfo{
					HostUUID:     info.OpId,
					HostClassify: config.HostStorageType,
				})
			if err != nil {
				log.Println("client.OpInformationTest err: ", zap.Error(err))
				listInfoChan <- resInfo
				return
			}
			listInfoChan <- res
		}(list[i])

	}
	listInfo := make(map[string]*pb.MonitorInfo, len(list))

	for i := 0; i < len(list); i++ {
		info := <-listInfoChan
		listInfo[info.OpId] = info
	}

	for i := 0; i < len(list); i++ {
		res, isExit := listInfo[list[i].OpId]
		if !isExit {
			continue
		}
		list[i].Actor = list[i].ColonyName
		list[i].StorageType = list[i].ColonyType
		list[i].CPUUseRate = float64(res.CpuUseRate)
		list[i].DiskStatus = res.DiskStatus
		list[i].SysDiskUseRate = res.SysDiskUseRate
		list[i].SysDiskLeave = res.SysDiskLeave
		list[i].NfsStatus = res.NfsStatus
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetDCStorageList
// @Tags      EquipMonitorApi
// @Summary   Obtain the list of DC original value machines
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getDCStorageList [get]
func (s *EquipMonitorApi) GetDCStorageList(c *gin.Context) {
	var pageInfo request.HostMonitorReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := hostRecordService.GetDCStorageList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}

	listInfoChan := make(chan *pb.MonitorInfo, len(list))

	// Query the corresponding hardware information for each machine
	for i := 0; i < len(list); i++ {
		go func(info responseNodel.DCStorageMonitorInfo) {
			resInfo := &pb.MonitorInfo{OpId: info.UUID}

			client := global.GateWayClinets.GetGateWayClinet(info.GatewayId)
			if client == nil {
				listInfoChan <- resInfo
				log.Println("Connection gateway failed: " + info.GatewayId)
				return
			}

			res, err := client.GetOpMonitorInfo(context.TODO(),
				&pb.OpHardwareInfo{
					HostUUID:     info.UUID,
					HostClassify: config.HostDCStorageType,
				})
			if err != nil {
				log.Println("client.OpInformationTest err: ", zap.Error(err))
				listInfoChan <- resInfo
				return
			}
			listInfoChan <- res
		}(list[i])

	}
	listInfo := make(map[string]*pb.MonitorInfo, len(list))

	for i := 0; i < len(list); i++ {
		info := <-listInfoChan
		listInfo[info.OpId] = info
	}

	for i := 0; i < len(list); i++ {
		res, isExit := listInfo[list[i].UUID]
		if !isExit {
			continue
		}
		list[i].CPUUseRate = float64(res.CpuUseRate)
		list[i].DiskStatus = res.DiskStatus
		list[i].SysDiskUseRate = res.SysDiskUseRate
		list[i].SysDiskLeave = res.SysDiskLeave
		list[i].NfsStatus = res.NfsStatus
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}

// GetHostScriptResult
// @Tags      EquipMonitorApi
// @Summary   Node machine script information return
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getHostScriptResult [get]
func (s *EquipMonitorApi) GetHostScriptResult(c *gin.Context) {
	var scriptReq request.HostScriptReq
	err := c.ShouldBindQuery(&scriptReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(scriptReq.GatewayId)

	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + scriptReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	if !utils.IsInStrList(scriptReq.ScriptInfo, config.AllowScriptList) {
		response.FailWithMessage("The script does not comply with the specifications", c)
		return
	}

	switch scriptReq.ScriptInfo {
	case config.CarDirName:
		scriptReq.ScriptInfo = config.CarDir
	case config.SealDirName:
		scriptReq.ScriptInfo = config.SealDir
	}

	res, err := client.GetOpScriptInfo(context.TODO(),
		&pb.OpScriptInfo{
			HostUUID:     scriptReq.UUID,
			HostClassify: scriptReq.HostClassify,
			ScriptInfo:   scriptReq.ScriptInfo,
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	response.OkWithDetailed(res, "Successfully obtained script information", c)
}


// LocalDiskReMounting 
func (s *EquipMonitorApi) LocalDiskReMounting(c *gin.Context) {
	var scriptReq request.HostScriptReq
	err := c.ShouldBindQuery(&scriptReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	client := global.GateWayClinets.GetGateWayClinet(scriptReq.GatewayId)

	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + scriptReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	_, err = client.GetOpScriptInfo(context.TODO(),
		&pb.OpScriptInfo{
			HostUUID:   scriptReq.UUID,
			ScriptInfo: define.PathDiskInitialization,
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	response.OkWithMessage("Execution succeeded", c)
}

// GetDiskLetter
// @Tags      EquipMonitorApi
// @Summary   Query drive letter
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getDiskLetter [get]
func (s *EquipMonitorApi) GetDiskLetter(c *gin.Context) {
	var diskLetterReq request.DiskLetterReq
	err := c.ShouldBindQuery(&diskLetterReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	client := global.GateWayClinets.GetGateWayClinet(diskLetterReq.GatewayId)

	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + diskLetterReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	res, err := client.GetDiskLetter(context.TODO(),
		&pb.DiskLetterReq{
			HostUUID:   diskLetterReq.UUID,
			DiskLetter: diskLetterReq.DiskLetter,
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	response.OkWithDetailed(res, "Successfully obtained", c)
}

// GetAbnormalDiskInfo
// @Tags      EquipMonitorApi
// @Summary   Get abnormal hard drive logs
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getAbnormalDiskInfo [get]
func (s *EquipMonitorApi) GetAbnormalDiskInfo(c *gin.Context) {
	var scriptReq request.HostScriptReq
	err := c.ShouldBindQuery(&scriptReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	client := global.GateWayClinets.GetGateWayClinet(scriptReq.GatewayId)

	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + scriptReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	res, err := client.GetOpScriptInfo(context.TODO(),
		&pb.OpScriptInfo{
			HostUUID:   scriptReq.UUID,
			ScriptInfo: `dmesg | grep dev | grep -E "I/O"`,
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	response.OkWithDetailed(res, "Successfully obtained script information", c)
}

// DiskReMounting
// @Tags      EquipMonitorApi
// @Summary   Local disk re mounting
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/diskReMounting [get]
func (s *EquipMonitorApi) DiskReMounting(c *gin.Context) {
	var scriptReq request.DiskReMountReq
	err := c.ShouldBindQuery(&scriptReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	client := global.GateWayClinets.GetGateWayClinet(scriptReq.GatewayId)

	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + scriptReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	_, err = client.DiskReMounting(context.TODO(),
		&pb.DiskReMountReq{
			HostClassify: scriptReq.HostClassify,
			HostUUID:     scriptReq.UUID,
			Actor:        scriptReq.Actor,
			NodeIP:       scriptReq.NodeIP,
			MountOpId:    scriptReq.MountOpId,
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	response.OkWithMessage("Execution succeeded", c)
}

// RestartRelatedServices
// @Tags      EquipMonitorApi
// @Summary   Restart related services
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/restartServices [get]
func (s *EquipMonitorApi) RestartRelatedServices(c *gin.Context) {
	var scriptReq request.HostScriptReq
	err := c.ShouldBindQuery(&scriptReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(scriptReq.GatewayId)

	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + scriptReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	if scriptReq.ScriptInfo == config.LotusLogType {
		if scriptReq.HostClassify != config.HostWorkerType {
			response.FailWithMessage("Host type error", c)
			return
		}
	} else {
		if scriptReq.HostClassify != config.HostMinerType {
			response.FailWithMessage("Host type error", c)
			return
		}
	}

	switch scriptReq.ScriptInfo {
	case config.LotusLogType:
		scriptReq.ScriptInfo = config.LotusRestartCmd
	case config.MinerLogType:
		scriptReq.ScriptInfo = config.MinerRestartCmd
	case config.BoostLogType:
		scriptReq.ScriptInfo = config.BoostRestartCmd
	case config.WorkerLogType:
		scriptReq.ScriptInfo = config.WorkerRestartCmd
	}

	res, err := client.GetOpScriptInfo(context.TODO(),
		&pb.OpScriptInfo{
			HostUUID:     scriptReq.UUID,
			HostClassify: config.HostDCStorageType,
			ScriptInfo:   scriptReq.ScriptInfo,
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	response.OkWithDetailed(res, "Successfully restart services", c)
}

// GetNodeStorageInfo
// @Tags      EquipMonitorApi
// @Summary   Associated storage machine information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getNodeStorageInfo [get]
func (s *EquipMonitorApi) GetNodeStorageInfo(c *gin.Context) {
	var scriptReq request.GetNodeStorageInfoReq
	err := c.ShouldBindQuery(&scriptReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(scriptReq.GatewayId)
	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + scriptReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	// Obtain the storage machine corresponding to this node
	var list []responseNodel.StorageMountErrorList
	list, err = lotus.DeployService.GetStorageMountInfoByActor(scriptReq.Actor)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Obtain the IP information of the actual hard disk mounted on the node
	res, err := client.GetOpMountInfo(context.TODO(),
		&pb.DiskLetterReq{
			HostUUID: scriptReq.UUID,
		})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	var resp []responseNodel.StorageMountErrorList
	var isMountError = true
	for _, val := range list {
		if !utils.IsInStrList(val.Ip, res.Ips) {
			resp = append(resp, val)
		}
	}
	if len(resp) > 0 {
		isMountError = false
	}

	response.OkWithDetailed(map[string]interface{}{
		"isMountError":   isMountError,
		"mountErrorList": resp,
	}, "Successfully obtained information", c)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240,
	WriteBufferSize: 10240,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GetTestWSTest �����洢����Ϣ
func (s *EquipMonitorApi) GetTestWSTest(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		// ��ȡ��Ϣ
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// ������Ϣ
		log.Printf("�յ���Ϣ: %s", msg)

		beginTime := time.Now()
		for {
			time.Sleep(5 * time.Second)
			err = conn.WriteMessage(websocket.TextMessage, []byte("Received: "+time.Now().Format("15:04:05")))
			if err != nil {
				log.Println(err)
				break
			}
			timeSub := time.Now().Sub(beginTime).Minutes()
			if timeSub > 1 {
				break
			}
		}
	}
}

// GetHostLogsTest �����洢����Ϣ
func (s *EquipMonitorApi) GetHostLogsTest(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		beginTime := time.Now()
		if string(msg) == "lotus" {
			var info string
			logBeginNum := getOplianLogLen()
			for {
				time.Sleep(500 * time.Millisecond)

				logBeginNum, info = getOplianLog(logBeginNum)

				err = conn.WriteMessage(websocket.TextMessage, []byte(info))
				if err != nil {
					log.Println(err)
					continue
				}
				timeSub := time.Now().Sub(beginTime).Minutes()
				if timeSub > 30 {
					continue
				}
			}
		}
	}
}

func getOplianLogLen() int {
	var logRowNum int

	logLen, err := exec.Command("bash", "-c", "wc -l /mnt/md0/ipfs/logs/lotus.log | awk '{print $1}'").CombinedOutput()
	if err != nil {
		log.Println("cmd `wc -l /mnt/md0/ipfs/logs/lotus.log | awk '{print $1}'` Filed", err.Error())
		return logRowNum
	}
	logRowNum, _ = strconv.Atoi(string(logLen)[:len(string(logLen))-1])
	return logRowNum
}

func getOplianLog(logBeginNum int) (int, string) {
	//log.Println("op NodeAddShareDir begin")
	var logRowNum int
	var logEndNum = logBeginNum

	logLen, err := exec.Command("bash", "-c", "wc -l /mnt/md0/ipfs/logs/lotus.log | awk '{print $1}'").CombinedOutput()
	if err != nil {
		log.Println("cmd `wc -l /mnt/md0/ipfs/logs/lotus.log | awk '{print $1}'` Filed", err.Error())
		return logEndNum, ""
	}
	if len(string(logLen)) < 0 {
		log.Println("get log info failed", err.Error())
		return logEndNum, ""
	}

	logRowNum, _ = strconv.Atoi(string(logLen)[:len(string(logLen))-1])

	if logRowNum > logBeginNum {
		if logRowNum > logBeginNum+100 {
			logEndNum = logBeginNum + 100
		} else {
			logEndNum = logRowNum
		}
		logCmd := "cat /mnt/md0/ipfs/logs/lotus.log | head -n " + strconv.Itoa(logEndNum) + " | tail -n +" + strconv.Itoa(logBeginNum+1)

		out, err := exec.Command("bash", "-c", logCmd).CombinedOutput()
		if err != nil {
			log.Println("cmd `"+logCmd+"` Filed", err.Error())
			return logEndNum, ""
		}
		return logEndNum, string(out)
	}

	return logEndNum, ""
}

// GetHostLogs
func (s *EquipMonitorApi) GetHostLogs(c *gin.Context) {
	var nodeInfoReq request.GetNodeLogInfoReq
	err := c.ShouldBindQuery(&nodeInfoReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	// Link to corresponding node gateway
	client := global.GateWayClinets.GetGateWayClinet(nodeInfoReq.GatewayId)
	if client == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Connection gateway failed"))
		return
	}

	// Receive corresponding identification and retrieve different log information
	for {
		// Obtain the corresponding information identifier, and only return information for specific identifiers
		beginTime := time.Now()
		switch nodeInfoReq.HostClassify {
		case config.HostWorkerType:
			if nodeInfoReq.LogType != config.WorkerLogType {
				conn.WriteMessage(websocket.TextMessage, []byte("host type error"))
				return
			}
		case config.HostMinerType:
			if !utils.IsInStrList(nodeInfoReq.LogType, []string{config.MinerProveLogType, config.MinerWdpostLogType,
				config.MinerRealTimeLogType, config.LotusLogType, config.BoostLogType}) {
				conn.WriteMessage(websocket.TextMessage, []byte("host type error"))
				return
			}
		default:
			conn.WriteMessage(websocket.TextMessage, []byte("host type error"))
			return
		}

		// Get the latest number of log lines
		opLogLenResp, err := client.GetOpLogLen(context.TODO(),
			&pb.OpLogInfoReq{
				HostUUID:     nodeInfoReq.UUID,
				HostClassify: nodeInfoReq.HostClassify,
				LogType:      nodeInfoReq.LogType,
			})
		if err != nil {
			log.Println(err)
			conn.WriteMessage(websocket.TextMessage, []byte("failed to connect to machine to get log len"))
			return
		}
		if opLogLenResp.GetLogLenNum() == 0 {
			conn.WriteMessage(websocket.TextMessage, []byte("the log information is empty"))
			return
		}
		logBeginNum := opLogLenResp.LogLenNum

		for {
			time.Sleep(500 * time.Millisecond)

			opLogInfoResp, err := client.GetOpLogInfo(context.TODO(),
				&pb.OpLogInfoReq{
					HostUUID:     nodeInfoReq.UUID,
					HostClassify: nodeInfoReq.HostClassify,
					LogType:      nodeInfoReq.LogType,
					LogBeginNum:  logBeginNum,
					GetNum:       config.DefaultGetNum,
				})
			if err != nil {
				log.Println(err)
				conn.WriteMessage(websocket.TextMessage, []byte("failed to connect to machine to get log info"))
				return
			}

			// Get log information
			logBeginNum = opLogInfoResp.LogBeginNum

			// Return websocket information
			err = conn.WriteMessage(websocket.TextMessage, []byte(opLogInfoResp.LogResp))
			if err != nil {
				log.Println(err)
				return
			}
			timeSub := time.Now().Sub(beginTime).Minutes()
			if timeSub > config.LogTimeDuration {
				return
			}
		}
	}
}

type DeadLineInfo struct {
	DeadLineId          int    `json:"deadLineId"`
	Partitions          int    `json:"partitions"`
	LivePartitions      int    `json:"livePartitions"`
	Sectors             int    `json:"sectors"`
	Live                int    `json:"live"`
	Active              int    `json:"active"`
	Fault               int    `json:"fault"`
	Recovery            int    `json:"recovery"`
	Terminated          int    `json:"terminated"`
	Unproven            int    `json:"unproven"`
	ProvenPartitions    string `json:"provenPartitions"`
	Current             string `json:"current"`
	ProvenPartitionsStr string `json:"provenPartitionsStr"`
	PeriodStart         int    `json:"periodStart"`
	PeriodEnd           int    `json:"periodEnd"`
	CurrentEpoch        int    `json:"currentEpoch"`
	Open                int    `json:"open"`
	OpenTime            string `json:"openTime"`
}

type Resp struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []DeadLineInfo `json:"data"`
}

// GetNodeMinerInfo
// @Tags      EquipMonitorApi
// @Summary   Obtain node miner information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /equipMonitor/getNodeMinerInfo [get]
func (s *EquipMonitorApi) GetNodeMinerInfo(c *gin.Context) {
	var scriptReq request.GetNodeStorageInfoReq
	err := c.ShouldBindQuery(&scriptReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(scriptReq.GatewayId)
	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: " + scriptReq.GatewayId)
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	opMinerInfo, err := lotus.DeployService.GetMinerByOpId(scriptReq.UUID)
	if err != nil {
		global.ZC_LOG.Error("get miner info failed" + err.Error())
		response.FailWithMessage("get miner info failed", c)
		return
	}

	// MinerInfo
	nodeResp := responseNodel.NodeMinerInfoResp{}

	if opMinerInfo.IsWnpost {
		nodeResp.MinerAttribute = config.MinerWdpostLogType
	}
	if opMinerInfo.IsWdpost {
		if len(nodeResp.MinerAttribute) == 0 {
			nodeResp.MinerAttribute = config.MinerWnpostLogType
		} else {
			nodeResp.MinerAttribute = nodeResp.MinerAttribute + ", " + config.MinerWnpostLogType
		}
	}

	// Obtain corresponding node window information
	url := config.FilMinerUrl + scriptReq.Actor + "/deadline"
	resp, err := http.Get(url)
	if err != nil {
		global.ZC_LOG.Error("get miner info by url err: ", zap.Error(err))
		response.FailWithMessage("Failed to get miner info by url", c)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.FailWithMessage("Failed to get miner info by url", c)
		return
	}

	info := new(Resp)
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("JSON parsing information failed!", err.Error())
		response.FailWithMessage("Failed to get information", c)
		return
	}

	currentKey := 0
	for key, val := range info.Data {
		if strings.Contains(val.Current, "current") {
			currentKey = key
		}
	}
	var nextTime time.Time
	if currentKey < 47 {
		nextTime, _ = time.ParseInLocation(config.TimeFormat, info.Data[currentKey+1].OpenTime, time.Local)
		nodeResp.WinRemainingTime = nextTime.Sub(time.Now()).String()
	}

	nodeResp.ProofRequiredTime = "30m"

	// Obtain node information
	res, err := client.GetNodeMinerInfo(context.TODO(),
		&pb.OpHardwareInfo{HostUUID: scriptReq.UUID})
	if err != nil {
		global.ZC_LOG.Error("client.OpInformationTest err: ", zap.Error(err))
		response.FailWithMessage("Failed to connect to machine", c)
		return
	}

	nodeResp.MinerProcessStatus = res.MinerProcessStatus
	nodeResp.MessageOut = res.MessageOut
	nodeResp.RestartStatus = res.RestartStatus
	nodeResp.DailyExplosiveNum = int(res.DailyExplosiveNum)

	response.OkWithDetailed(nodeResp, "Successfully obtained information", c)
}

// GetHostLogsNum
func (s *EquipMonitorApi) GetHostLogsNum(c *gin.Context) {
	// ��ȡ��Ӧ�ڵ���Ϣ
	var nodeInfoReq request.GetNodeLogInfoReq
	err := c.ShouldBindQuery(&nodeInfoReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// ���Ӷ�Ӧ�ڵ�gateway
	client := global.GateWayClinets.GetGateWayClinet(nodeInfoReq.GatewayId)
	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: ", zap.Error(err))
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	switch nodeInfoReq.HostClassify {
	case config.HostWorkerType:
		if nodeInfoReq.LogType != config.WorkerLogType {
			response.FailWithMessage("Host type error", c)
			return
		}
	case config.HostMinerType:
		if !utils.IsInStrList(nodeInfoReq.LogType, []string{config.MinerProveLogType, config.MinerWdpostLogType,
			config.MinerRealTimeLogType, config.LotusLogType, config.BoostLogType}) {
			response.FailWithMessage("Host type error", c)
			return
		}
	default:
		response.FailWithMessage("Host type error", c)
		return
	}

	// ��ȡ������־����
	opLogLenResp, err := client.GetOpLogLen(context.TODO(),
		&pb.OpLogInfoReq{
			HostUUID:     nodeInfoReq.UUID,
			HostClassify: nodeInfoReq.HostClassify,
			LogType:      nodeInfoReq.LogType,
		})
	if err != nil {
		log.Println(err)
		response.FailWithMessage("Failed to connect to machine to get log len", c)
		return
	}
	if opLogLenResp.GetLogLenNum() == 0 {
		response.FailWithMessage("The log information is empty", c)
		return
	}
	logBeginNum := opLogLenResp.LogLenNum

	response.OkWithDetailed(map[string]int64{"beginNum": logBeginNum}, "Successfully obtained information", c)
}

// GetHostLogs
func (s *EquipMonitorApi) GetHostLogsInfo(c *gin.Context) {
	var nodeInfoReq request.GetNodeLogInfoInfoReq
	err := c.ShouldBindQuery(&nodeInfoReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	client := global.GateWayClinets.GetGateWayClinet(nodeInfoReq.GatewayId)
	if client == nil {
		global.ZC_LOG.Error("Connection gateway failed: ", zap.Error(err))
		response.FailWithMessage("Connection gateway failed", c)
		return
	}

	switch nodeInfoReq.HostClassify {
	case config.HostWorkerType:
		if nodeInfoReq.LogType != config.WorkerLogType {
			response.FailWithMessage("Host type error", c)
			return
		}
	case config.HostMinerType:
		if !utils.IsInStrList(nodeInfoReq.LogType, []string{config.MinerProveLogType, config.MinerWdpostLogType,
			config.MinerRealTimeLogType, config.LotusLogType, config.BoostLogType}) {
			response.FailWithMessage("Host type error", c)
			return
		}
	default:
		response.FailWithMessage("Host type error", c)
		return
	}

	opLogInfoResp, err := client.GetOpLogInfo(context.TODO(),
		&pb.OpLogInfoReq{
			HostUUID:     nodeInfoReq.UUID,
			HostClassify: nodeInfoReq.HostClassify,
			LogType:      nodeInfoReq.LogType,
			LogBeginNum:  nodeInfoReq.BeginNum,
			GetNum:       config.DefaultGetNum,
		})
	if err != nil {
		log.Println(err)
		response.FailWithMessage("Failed to connect to machine to get log info", c)
		return
	}

	logBeginNum := opLogInfoResp.LogBeginNum

	type logResp struct {
		LogResp  string `json:"logResp"`
		BeginNum int64  `json:"beginNum"`
	}
	response.OkWithDetailed(logResp{
		LogResp:  opLogInfoResp.LogResp,
		BeginNum: logBeginNum,
	}, "Successfully obtained information", c)
}
