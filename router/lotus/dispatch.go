package lotus

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type DispatchRouter struct{}

// InitDispatchRouter Scheduling Routing
func (s *DispatchRouter) InitDispatchRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	dispatchRouter := Router.Group("dispatch")
	dispatchApi := v1.ApiGroupApp.LotusApiGroup.DispatchApi
	{
		dispatchRouter.GET("workerConfigList", dispatchApi.WorkerConfigList)                      // Gets the worker configuration list
		dispatchRouter.POST("setWorkerConfig", dispatchApi.SetWorkerConfig)                       // Set worker Settings
		dispatchRouter.GET("getSectorsList", dispatchApi.GetSectorsList)                          // Gets a list of sectors
		dispatchRouter.GET("getSectorDetails", dispatchApi.GetSectorDetails)                      // Get the sector details
		dispatchRouter.GET("getSectorTaskList", dispatchApi.GetSectorTaskList)                    // Get the sector task list
		dispatchRouter.GET("getSectorTaskDetailList", dispatchApi.GetSectorTaskDetailList)        // Get the sector mission details list
		dispatchRouter.GET("sectorRecoverDetail", dispatchApi.SectorRecoverDetail)                // Sector recovery task details
		dispatchRouter.PUT("addSectorTask", dispatchApi.AddSectorTask)                            // Add Task Queue
		dispatchRouter.POST("addSectorDcTask", dispatchApi.AddSectorDcTask)                       // Add DC task queue
		dispatchRouter.POST("editTaskQueueStatus", dispatchApi.EditTaskQueueStatus)               // Modify the task queue
		dispatchRouter.POST("editTaskQueueDetailStatus", dispatchApi.EditTaskQueueDetailStatus)   // Modify the task detail state
		dispatchRouter.POST("editTaskQueueConcurrent", dispatchApi.EditTaskQueueDetailConcurrent) // Modify the task queue number state
		dispatchRouter.GET("getManageMiners", dispatchApi.GetManageMiners)                        // Get the schedule list
		dispatchRouter.POST("checkDealCar", dispatchApi.CheckDealCar)                             // Matching order
		dispatchRouter.POST("onOffPre1", dispatchApi.OnOffPre1)                                   // Start stop task
	}

	return dispatchRouter
}
