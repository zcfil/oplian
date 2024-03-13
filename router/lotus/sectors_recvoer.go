package lotus

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type SectorsRecoverRouter struct{}

// InitSectorsRecoverRouter Scheduling Routing
func (s *DispatchRouter) InitSectorsRecoverRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	sectorsRecoverRouter := Router.Group("sectorsRecover")
	sectorsRecoverApi := v1.ApiGroupApp.LotusApiGroup.SectorsRecoverApi
	{
		sectorsRecoverRouter.POST("sectorsRecoverList", sectorsRecoverApi.SectorsRecoverList)         // Gets the sector list
		sectorsRecoverRouter.POST("addSectorsRecoverTask", sectorsRecoverApi.AddSectorsRecoverTask)   // Add the sector recovery task
		sectorsRecoverRouter.POST("workerOpList", sectorsRecoverApi.WorkerOpList)                     // Get a list of worker hosts
		sectorsRecoverRouter.POST("modifySectorTaskStatus", sectorsRecoverApi.ModifySectorTaskStatus) // Changes the sector state
	}

	return sectorsRecoverRouter
}
