package system

import (
	"github.com/gin-gonic/gin"
	v1 "oplian/api/v1"
)

type MachineRoomRecordRouter struct{}

func (s *MachineRoomRecordRouter) InitSysMachineRoomRecordRouter(Router *gin.RouterGroup) {
	machineRoomRecordRouter := Router.Group("sysMachineRoomRecord")
	machineRoomRecordApi := v1.ApiGroupApp.SystemApiGroup.MachineRoomRecordApi
	{
		machineRoomRecordRouter.POST("createSysMachineRoomRecord", machineRoomRecordApi.CreateSysMachineRoomRecord)                     // Create a new SysMachineRoomRecord
		machineRoomRecordRouter.POST("updateSysMachineRoomRecord", machineRoomRecordApi.UpdateSysMachineRoomRecord)                     // Update SysMachineRoomRecord
		machineRoomRecordRouter.DELETE("deleteSysMachineRoomRecordByRoomIds", machineRoomRecordApi.DeleteSysMachineRoomRecordByRoomIds) // Batch delete SysMachineRoomRecord
		machineRoomRecordRouter.GET("getSysMachineRoomRecordList", machineRoomRecordApi.GetSysMachineRoomRecordList)                    // Gets the SysMachineRoomRecord list
		machineRoomRecordRouter.GET("getMachineRoomId", machineRoomRecordApi.GetMachineRoomId)                                          // New Room Gets the generated roomId

		machineRoomRecordRouter.POST("bindSysHostRecord", machineRoomRecordApi.BindSysHostRecords)        // Bind host Information
		machineRoomRecordRouter.GET("bindSysHostRecordList", machineRoomRecordApi.BindSysHostRecordsList) // Gets a list of bound and unbound hosts
		machineRoomRecordRouter.POST("unBindSysHostRecord", machineRoomRecordApi.UnbindSysHostRecords)    // Unbind the host information
		machineRoomRecordRouter.GET("getRoomRecordList", machineRoomRecordApi.GetRoomRecordList)          //Get the list of computer rooms (drop-down)
	}
}
