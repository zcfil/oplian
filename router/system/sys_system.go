package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
	"oplian/middleware"
)

type SysRouter struct{}

func (s *SysRouter) InitSystemRouter(Router *gin.RouterGroup) {
	sysRouter := Router.Group("system").Use(middleware.OperationRecord())
	systemApi := v1.ApiGroupApp.SystemApiGroup.SystemApi
	{
		sysRouter.POST("getSystemConfig", systemApi.GetSystemConfig) // Get the contents of the configuration file
		sysRouter.POST("setSystemConfig", systemApi.SetSystemConfig) // Set the content of the configuration file
		sysRouter.POST("getServerInfo", systemApi.GetServerInfo)     // Get server information
		sysRouter.POST("reloadSystem", systemApi.ReloadSystem)       // Restart service
	}
}
