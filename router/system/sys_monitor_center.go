package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type MonitorCenterRouter struct {
}

func (s *JobPlatformRouter) InitMonitorCenterRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	monitorCenterRouter := Router.Group("monitorCenter") //.UsePercent(middleware.OperationRecord())
	monitorCenterApi := v1.ApiGroupApp.SystemApiGroup.MonitorCenterApi
	{
		monitorCenterRouter.GET("businessMonitor", monitorCenterApi.BusinessMonitor) //Business monitoring
	}

	return monitorCenterRouter
}
