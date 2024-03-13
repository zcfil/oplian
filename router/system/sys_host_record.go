package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type HostRecordRouter struct{}

func (s *HostRecordRouter) InitSysHostRecordRouter(Router *gin.RouterGroup) {
	hostRecordRouter := Router.Group("sysHostRecord")
	hostRecordApi := v1.ApiGroupApp.SystemApiGroup.HostRecordApi
	{
		//hostRecordRouter.POST("createSysHostRecord", hostRecordApi.CreateSysHostRecord)                 // Create a new SysHostRecord
		hostRecordRouter.POST("updateSysHostRecord", hostRecordApi.UpdateSysHostRecord)                 // Update SysHostRecord
		hostRecordRouter.DELETE("deleteSysHostRecordByUUIDs", hostRecordApi.DeleteSysHostRecordByUUIDs) // Batch delete SysHostRecord
		hostRecordRouter.GET("getSysHostRecordList", hostRecordApi.GetSysHostRecordList)                // Gets the SysHostRecord list
		hostRecordRouter.GET("getSysHostListNormal", hostRecordApi.GetSysHostListNormal)                // Gets a list of running hosts
		hostRecordRouter.GET("getSysHostList", hostRecordApi.GetSysHostList)                            // Get the SysHost list
		hostRecordRouter.GET("hostTestRecordList", hostRecordApi.GetSysHostTestRecordList)              // Get the SysHost list (Classifiable query)
		hostRecordRouter.POST("opInfoList", hostRecordApi.OpInfoList)                                   // Get Op information

		hostRecordRouter.GET("getOpHardwareInfo", hostRecordApi.GetOpHardwareInfo) // Get Op hardware information
		hostRecordRouter.GET("netHostList", hostRecordApi.GetNetHostList)          // Gets the list of selected network hosts
		hostRecordRouter.GET("patrolHostList", hostRecordApi.GetPatrolHostList)    // Gets a list of patrol target hosts
	}
}
