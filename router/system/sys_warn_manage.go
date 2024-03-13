package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type WarnManageRouter struct {
}

func (w *WarnManageRouter) InitWarnManageRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	warnManageRouter := Router.Group("warnManage") //.UsePercent(middleware.OperationRecord())
	warnManageApi := v1.ApiGroupApp.SystemApiGroup.WarnManageApi
	{
		warnManageRouter.POST("warnTotal", warnManageApi.WarnTotal)                       //Statistics of the alarm center
		warnManageRouter.POST("warnTrend", warnManageApi.WarnTrend)                       //Status diagram of the alarm center
		warnManageRouter.POST("warnList", warnManageApi.WarnList)                         //Alarm center Indicates the alarm list
		warnManageRouter.POST("modifyWarnStatus", warnManageApi.ModifyWarnStatus)         //The alarm center updates the alarm status
		warnManageRouter.POST("strategyId", warnManageApi.StrategyId)                     //Get policy Id
		warnManageRouter.POST("strategyType", warnManageApi.StrategyType)                 //Acquisition policy type
		warnManageRouter.POST("saveStrategy", warnManageApi.SaveStrategy)                 //Save or modify the alarm policy
		warnManageRouter.POST("strategyList", warnManageApi.StrategyList)                 //Alarm policy list
		warnManageRouter.POST("strategyDetail", warnManageApi.StrategyDetail)             //Alarm Policy Details
		warnManageRouter.POST("delStrategy", warnManageApi.DelStrategy)                   //Deleting an Alarm policy
		warnManageRouter.POST("modifyStrategyStatus", warnManageApi.ModifyStrategyStatus) //Example Change the status of an alarm policy
	}

	return warnManageRouter
}
