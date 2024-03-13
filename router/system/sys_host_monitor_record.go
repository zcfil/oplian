package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type HostMonitorRecordRouter struct{}

func (s *HostMonitorRecordRouter) InitSysHostMonitorRecordRouter(Router *gin.RouterGroup) {
	hostMonitorRecordRouter := Router.Group("sysHostMonitorRecord")
	hostMonitorRecordApi := v1.ApiGroupApp.SystemApiGroup.HostMonitorRecordApi
	{
		hostMonitorRecordRouter.GET("getSysHostMonitorRecordList", hostMonitorRecordApi.GetSysHostMonitorRecordList) // Gets a list of host-specific hardware usage
		hostMonitorRecordRouter.POST("getSysHostMonitorLineChart", hostMonitorRecordApi.GetSysHostMonitorLineChart)  // Gets a line chart of host usage
		hostMonitorRecordRouter.GET("getStorageInfoList", hostMonitorRecordApi.GetStorageInformationList)            // Storage machine stores information list information
	}
}
