package slot

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type SlotRouter struct{}

func (u *SlotRouter) InitSlotRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	slotRouter := Router.Group("slot")
	slotApi := v1.ApiGroupApp.SlotApiGroup.SlotApi
	{
		slotRouter.POST("installSlot", slotApi.InstallSlot) // 安装插件
	}
	return slotRouter
}
