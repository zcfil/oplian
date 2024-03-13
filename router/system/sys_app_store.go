package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type AppStoreRouter struct{}

func (s *AppStoreRouter) InitAppStoreRouter(Router *gin.RouterGroup) {
	appStoreRouter := Router.Group("slot")
	appStoreApi := v1.ApiGroupApp.SystemApiGroup.AppStoreApi
	{
		appStoreRouter.POST("getSlotList", appStoreApi.GetSlotList)         // Get a list of plug-ins
		appStoreRouter.POST("getSlotFileList", appStoreApi.GetSlotFileList) // Obtain the plug-in file list

		appStoreRouter.POST("replaceSlotFile", appStoreApi.ReplaceSlotFile) // Replacement file
	}
}
