package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type HomeInterfaceRouter struct{}

func (s *HomeInterfaceRouter) InitHomeInterfaceRouter(Router *gin.RouterGroup) {
	homeInterfaceRouter := Router.Group("sysHomeInterface")
	homeInterfaceApi := v1.ApiGroupApp.SystemApiGroup.HomeInterfaceApi
	{
		homeInterfaceRouter.GET("getDataScreening", homeInterfaceApi.GetHostDataScreening) // Get host data overview
		homeInterfaceRouter.GET("getHostUseData", homeInterfaceApi.GetHostUseData)         // Gets the host usage parameters
		homeInterfaceRouter.GET("getHostRunList", homeInterfaceApi.GetHostRunList)         // Gets the host operation
	}
}
