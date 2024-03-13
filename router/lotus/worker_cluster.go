package lotus

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type WorkerClusterRouter struct{}

// InitWorkerClusterRouter 调度路由
func (s *DispatchRouter) InitWorkerClusterRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	workerClusterRouter := Router.Group("workerCluster")
	workerClusterApi := v1.ApiGroupApp.LotusApiGroup.WorkerClusterApi
	{
		workerClusterRouter.POST("workerClusterList", workerClusterApi.WorkerClusterList) // Gets the list of worker clusters
		workerClusterRouter.POST("addWorkerOp", workerClusterApi.AddWorkerOp)             // Added worker op
		workerClusterRouter.POST("delWorkerOp", workerClusterApi.DelWorkerOp)             // Delete worker op
		workerClusterRouter.POST("workerTaskDetail", workerClusterApi.WorkerTaskDetail)   // Worker task details
		workerClusterRouter.GET("exportWorkerTask", workerClusterApi.ExportWorkerTask)    // Export the list of worker tasks
	}

	return workerClusterRouter
}
