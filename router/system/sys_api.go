package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
	"oplian/middleware"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("api").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := Router.Group("api")
	apiRouterApi := v1.ApiGroupApp.SystemApiGroup.SystemApiApi
	{
		apiRouter.POST("createApi", apiRouterApi.CreateApi)               // Create Api
		apiRouter.POST("deleteApi", apiRouterApi.DeleteApi)               // Delete Api
		apiRouter.POST("getApiById", apiRouterApi.GetApiById)             // Gets a single Api message
		apiRouter.POST("updateApi", apiRouterApi.UpdateApi)               // Update api
		apiRouter.DELETE("deleteApisByIds", apiRouterApi.DeleteApisByIds) // Delete selected api
	}
	{
		apiRouterWithoutRecord.POST("getAllApis", apiRouterApi.GetAllApis) // Get all apis
		apiRouterWithoutRecord.POST("getApiList", apiRouterApi.GetApiList) // Get a list of apis
	}
}
