package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type HostTestRecordRouter struct{}

func (s *HostRecordRouter) InitSysHostTestRecordRouter(Router *gin.RouterGroup) {
	hostTestRecordRouter := Router.Group("sysHostTestRecord")
	hostTestRecordApi := v1.ApiGroupApp.SystemApiGroup.HostTestRecordApi
	{
		hostTestRecordRouter.GET("getSysHostTestRecordList", hostTestRecordApi.GetSysHostTestRecordList) // Get the test management list
		hostTestRecordRouter.GET("getSysHostTestReport", hostTestRecordApi.GetSysHostTestReport)         // Get the test management single message test report
		hostTestRecordRouter.POST("addHostTestByHand", hostTestRecordApi.AddHostTestByHand)              // Manually add a host test
		hostTestRecordRouter.POST("closeHostTest", hostTestRecordApi.CloseHostTest)                      // Shut down host testing
		hostTestRecordRouter.POST("restartAddHostTest", hostTestRecordApi.RestartAddHostTest)            // Re-run the host test
		hostTestRecordRouter.GET("defaultHostTestInfo", hostTestRecordApi.DefaultHostTestInfo)           // Default test information
	}
}
