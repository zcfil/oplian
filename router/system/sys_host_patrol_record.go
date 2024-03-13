package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type HostPatrolRecordRouter struct{}

func (s *HostPatrolRecordRouter) InitSysHostPatrolRecordRouter(Router *gin.RouterGroup) {
	hostPatrolRecordRouter := Router.Group("sysHostPatrolRecord")
	hostPatrolRecordApi := v1.ApiGroupApp.SystemApiGroup.HostPatrolRecordApi
	{
		hostPatrolRecordRouter.GET("getSysHostPatrolRecordList", hostPatrolRecordApi.GetSysHostPatrolRecordList) // Get the patrol management list
		hostPatrolRecordRouter.GET("getSysHostPatrolReport", hostPatrolRecordApi.GetSysHostPatrolReport)         // Get the inspection management single information inspection report
		hostPatrolRecordRouter.POST("addHostPatrolByHand", hostPatrolRecordApi.AddHostPatrolByHand)              // Manually initiate host inspection
	}
}
