package lotus

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type DeployRouter struct{}

// InitDeployRouter Gateway Connection Routing
func (s *DeployRouter) InitDeployRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	deployRouter := Router.Group("deploy")
	deployApi := v1.ApiGroupApp.LotusApiGroup.DeployApi
	{
		deployRouter.POST("editRunHost", deployApi.EditRunHost) // Modify the host service status

		//lotus Management
		deployRouter.PUT("addLotus", deployApi.AddLotus)                    // Add lotus
		deployRouter.POST("getLotusList", deployApi.GetLotusList)           // Get lotus List
		deployRouter.POST("getWalletList", deployApi.GetWalletList)         // Gets a list of lotus wallets
		deployRouter.POST("getRoomWalletList", deployApi.GetRoomWalletList) // Get wallet list
		deployRouter.POST("relationMinerList", deployApi.RelationMinerList) // Get the associated miner list
		//miner Management
		deployRouter.PUT("addMiner", deployApi.AddMiner)                        // Add miner
		deployRouter.GET("getNodeList", deployApi.GetNodeList)                  // Get node list
		deployRouter.POST("getMinerList", deployApi.GetMinerList)               // Get miner list
		deployRouter.POST("getNodeNumList", deployApi.GetNodesNum)              // Get the number of nodes
		deployRouter.POST("roomAllLotus", deployApi.RoomAllLotus)               // lotus machine Room
		deployRouter.POST("relationWorkerList", deployApi.RelationWorkerList)   // Gets the list of associated workers
		deployRouter.POST("relationStorageList", deployApi.RelationStorageList) // Get the associated storage list
		deployRouter.POST("checkMiner", deployApi.CheckMiner)                   // Check miner role
		deployRouter.POST("modifyMinerRole", deployApi.ModifyMinerRole)         // Modifying miner roles
		//boost Management
		deployRouter.PUT("addBoost", deployApi.AddBoost)         // Add boost
		deployRouter.GET("boostInfo", deployApi.GetBoost)        // Get boost Info
		deployRouter.GET("queryAsk", deployApi.QueryAsk)         // inquiry
		deployRouter.GET("queryDataCap", deployApi.QueryDataCap) // Query DC Quota

		//worker Management
		deployRouter.PUT("addWorker", deployApi.AddWorker)          // Add worker
		deployRouter.PUT("addStorage", deployApi.AddStorage)        // Add storage
		deployRouter.POST("getWorkerList", deployApi.GetWorkerList) // Get worker list

		deployRouter.POST("resetWorker", deployApi.ResetWorker)    // Reset worker
		deployRouter.GET("workerSectors", deployApi.WorkerSectors) // Query sector list on the worker
		deployRouter.GET("getMinerSelectList", deployApi.GetMinerSelectList)
		//Storage Management
		deployRouter.GET("getStorageList", deployApi.GetStorageList) // Get storage list

		deployRouter.POST("downloadSnapshot", deployApi.DownloadSnapshot) // Download Snapshot
	}

	return deployRouter
}
