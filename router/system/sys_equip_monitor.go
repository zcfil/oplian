package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type EquipMonitorRouter struct{}

func (s *EquipMonitorRouter) InitEquipMonitorRouter(Router *gin.RouterGroup) {
	equipMonitorRouter := Router.Group("equipMonitor")
	equipMonitorApi := v1.ApiGroupApp.SystemApiGroup.EquipMonitorApi
	{
		equipMonitorRouter.GET("getMinerList", equipMonitorApi.GetMinerList)         // Get a list of Miner machines
		equipMonitorRouter.GET("getWorkerList", equipMonitorApi.GetWorkerList)       // Get a list of Worker machines
		equipMonitorRouter.GET("getStorageList", equipMonitorApi.GetStorageList)     // Gets a list of Storage machines
		equipMonitorRouter.GET("getDCStorageList", equipMonitorApi.GetDCStorageList) // Get the DC original check-in list

		equipMonitorRouter.GET("getHostScriptResult", equipMonitorApi.GetHostScriptResult) // The node script information is returned
		equipMonitorRouter.GET("getDiskLetter", equipMonitorApi.GetDiskLetter)             // Query drive
		equipMonitorRouter.GET("getAbnormalDiskInfo", equipMonitorApi.GetAbnormalDiskInfo) // Query the abnormal hard disk log
		equipMonitorRouter.GET("localDiskReMounting", equipMonitorApi.LocalDiskReMounting) // The local disk is mounted again
		equipMonitorRouter.GET("diskReMounting", equipMonitorApi.DiskReMounting)           // The disk is remounted
		equipMonitorRouter.GET("restartServices", equipMonitorApi.RestartRelatedServices)  // Restart related services

		equipMonitorRouter.GET("getNodeStorageInfo", equipMonitorApi.GetNodeStorageInfo) // Associate storage device information
		equipMonitorRouter.GET("getNodeMinerInfo", equipMonitorApi.GetNodeMinerInfo)     // Get node miner information

		equipMonitorRouter.GET("getHostLogsNum", equipMonitorApi.GetHostLogsNum)   // 
		equipMonitorRouter.GET("getHostLogsInfo", equipMonitorApi.GetHostLogsInfo) // 

	}
}

func (s *EquipMonitorRouter) InitEquipWSRouter(Router *gin.RouterGroup) {
	equipMonitorRouter := Router.Group("equipMonitor")
	equipMonitorApi := v1.ApiGroupApp.SystemApiGroup.EquipMonitorApi
	{
		// websocket
		equipMonitorRouter.GET("getTestWSTest", equipMonitorApi.GetTestWSTest)     // test webSocket
		equipMonitorRouter.GET("getHostLogsTest", equipMonitorApi.GetHostLogsTest) // 
		equipMonitorRouter.GET("getHostLogs", equipMonitorApi.GetHostLogs)         // 
	}
}
