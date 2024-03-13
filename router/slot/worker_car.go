package slot

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type WorkerCarRouter struct{}

// InitWorkerCarRouter 调度路由
func (u *SlotRouter) InitWorkerCarRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	workerCarRouter := Router.Group("workerCar")
	workerCarApi := v1.ApiGroupApp.SlotApiGroup.WorkerCarApi
	{
		workerCarRouter.GET("workerCarList", workerCarApi.WorkerCarList)             // 获取workerCar任务列表
		workerCarRouter.POST("modifyTaskStatus", workerCarApi.ModifyTaskStatus)      // 修改workerCar任务状态
		workerCarRouter.POST("modifyWorkerNum", workerCarApi.ModifyWorkerNum)        // 修改workerCar worker任务数
		workerCarRouter.GET("workerCarTaskDetail", workerCarApi.WorkerCarTaskDetail) // workerCar任务详情
		workerCarRouter.POST("addWorkerCarTask", workerCarApi.AddWorkerCarTask)      // 新增workerCar任务
		workerCarRouter.GET("dcQuotaWalletList", workerCarApi.DcQuotaWalletList)     // 获取DC额度钱包
	}

	return workerCarRouter
}
