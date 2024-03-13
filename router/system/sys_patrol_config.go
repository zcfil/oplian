package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type PatrolConfigRouter struct{}

func (s *PatrolConfigRouter) InitSysPatrolConfigRouter(Router *gin.RouterGroup) {
	patrolConfigRouter := Router.Group("sysPatrolConfig")
	patrolConfigApi := v1.ApiGroupApp.SystemApiGroup.PatrolConfigApi
	{
		patrolConfigRouter.GET("getSysPatrolConfigList", patrolConfigApi.GetSysPatrolConfigList) // Gets a list of patrol Settings
		patrolConfigRouter.POST("updateSysPatrolConfig", patrolConfigApi.UpdateSysPatrolConfig)  // Update host inspection Settings
	}
}
