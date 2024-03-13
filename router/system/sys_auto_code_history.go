package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type AutoCodeHistoryRouter struct{}

func (s *AutoCodeRouter) InitAutoCodeHistoryRouter(Router *gin.RouterGroup) {
	autoCodeHistoryRouter := Router.Group("autoCode")
	autoCodeHistoryApi := v1.ApiGroupApp.SystemApiGroup.AutoCodeHistoryApi
	{
		autoCodeHistoryRouter.POST("getMeta", autoCodeHistoryApi.First)         // Get meta information based on id
		autoCodeHistoryRouter.POST("rollback", autoCodeHistoryApi.RollBack)     // rollback
		autoCodeHistoryRouter.POST("delSysHistory", autoCodeHistoryApi.Delete)  // Delete the rollback record
		autoCodeHistoryRouter.POST("getSysHistory", autoCodeHistoryApi.GetList) // Gets the rollback record page
	}
}
