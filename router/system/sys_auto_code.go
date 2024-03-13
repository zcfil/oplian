package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type AutoCodeRouter struct{}

func (s *AutoCodeRouter) InitAutoCodeRouter(Router *gin.RouterGroup) {
	autoCodeRouter := Router.Group("autoCode")
	autoCodeApi := v1.ApiGroupApp.SystemApiGroup.AutoCodeApi
	{
		autoCodeRouter.GET("getDB", autoCodeApi.GetDB)                  // Get database
		autoCodeRouter.GET("getTables", autoCodeApi.GetTables)          // Gets the tables for the corresponding database
		autoCodeRouter.GET("getColumn", autoCodeApi.GetColumn)          // Gets information about all fields of a specified table
		autoCodeRouter.POST("preview", autoCodeApi.PreviewTemp)         // Get a preview of the automatically created code
		autoCodeRouter.POST("createTemp", autoCodeApi.CreateTemp)       // Create automated code
		autoCodeRouter.POST("createPackage", autoCodeApi.CreatePackage) // Create a package
		autoCodeRouter.POST("getPackage", autoCodeApi.GetPackage)       // Get package
		autoCodeRouter.POST("delPackage", autoCodeApi.DelPackage)       // Delete package
		autoCodeRouter.POST("createPlug", autoCodeApi.AutoPlug)         // Automatic plugin package template
		autoCodeRouter.POST("installPlugin", autoCodeApi.InstallPlugin) // Automatic installation plug-in
	}
}
