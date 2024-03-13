package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type JobPlatformRouter struct {
}

func (s *JobPlatformRouter) InitJobPlatformRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	jobPlatformRouter := Router.Group("jobPlatform") //.UsePercent(middleware.OperationRecord())
	jobPlatformApi := v1.ApiGroupApp.SystemApiGroup.JobPlatformApi
	{
		jobPlatformRouter.POST("executeRecordsDetail", jobPlatformApi.ExecuteRecordsDetail) //Script record details
		jobPlatformRouter.POST("executeScript", jobPlatformApi.ExecuteScript)               //Execute script
		jobPlatformRouter.POST("fileDistribution", jobPlatformApi.FileDistribution)         //Document distribution
		jobPlatformRouter.POST("executeRecordsList", jobPlatformApi.ExecuteRecordsList)     //Execute the history list
		jobPlatformRouter.POST("batchUploadFile", jobPlatformApi.BatchUploadFile)           //Batch file upload
		jobPlatformRouter.POST("opFileToGateWay", jobPlatformApi.OpFileToGateWay)           //Synchronize files from the host to the gateway
		jobPlatformRouter.POST("saveFileHost", jobPlatformApi.SaveFileHost)                 //Setup file host
		jobPlatformRouter.POST("fileHostList", jobPlatformApi.FileHostList)                 //File host list
		jobPlatformRouter.POST("fileManageList", jobPlatformApi.FileManageList)             //File management list
		jobPlatformRouter.POST("delFileManage", jobPlatformApi.DelFileManage)               //Delete file management
		jobPlatformRouter.POST("modifyFileManage", jobPlatformApi.ModifyFileManage)         //Update file management
		jobPlatformRouter.POST("addFile", jobPlatformApi.AddFile)                           //New file
		jobPlatformRouter.POST("sysFilePoint", jobPlatformApi.SysFilePoint)                 //Op Point-to-point synchronization
		jobPlatformRouter.POST("fileListByType", jobPlatformApi.FileListByType)             //Get the file based on the file type
		jobPlatformRouter.POST("downLoadFile", jobPlatformApi.DownLoadFile)                 //gateWay Server file download
		jobPlatformRouter.POST("lotusHeightList", jobPlatformApi.LotusHeightList)           //List of lotus height files
		jobPlatformRouter.POST("minerList", jobPlatformApi.MinerList)                       //miner list
		jobPlatformRouter.POST("fileForcedToStop", jobPlatformApi.FileForcedToStop)         //Forcible termination of file distribution
		jobPlatformRouter.POST("executeScriptStop", jobPlatformApi.ExecuteScriptStop)       //Script execution is forcibly terminated
	}

	return jobPlatformRouter
}
