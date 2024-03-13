package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type HostGroupRouter struct{}

func (s *HostGroupRouter) InitSysHostGroupRouter(Router *gin.RouterGroup) {
	hostGroupRouter := Router.Group("sysHostGroup")
	hostGroupApi := v1.ApiGroupApp.SystemApiGroup.HostGroupApi
	{
		hostGroupRouter.POST("createSysHostGroup", hostGroupApi.CreateSysHostGroup)             // Create a new SysHostGroup
		hostGroupRouter.DELETE("deleteSysHostGroupByIds", hostGroupApi.DeleteSysHostGroupByIds) // Batch delete SysHostGroup
		hostGroupRouter.GET("getSysHostGroupList", hostGroupApi.GetSysHostGroupList)            // Get the SysHostGroup list
		hostGroupRouter.POST("dealSysHostGroupInfo", hostGroupApi.DealSysHostGroupInfo)         // Process host information in batches
	}
}
