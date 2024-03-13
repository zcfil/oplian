package system

import (
	"oplian/service"
)

type ApiGroup struct {
	DBApi
	BaseApi
	SystemApi
	CasbinApi
	AutoCodeApi
	SystemApiApi
	AuthorityApi
	DictionaryApi
	AuthorityMenuApi
	OperationRecordApi
	AutoCodeHistoryApi
	DictionaryDetailApi
	AuthorityBtnApi
	MachineRoomRecordApi
	JobPlatformApi
	HostGroupApi
	HostRecordApi
	WarnManageApi
	HostMonitorRecordApi
	MonitorCenterApi
	HostTestRecordApi
	HostPatrolRecordApi
	PatrolConfigApi
	HomeInterfaceApi
	EquipMonitorApi
	AppStoreApi
}

var (
	apiService               = service.ServiceGroupApp.SystemServiceGroup.ApiService
	menuService              = service.ServiceGroupApp.SystemServiceGroup.MenuService
	userService              = service.ServiceGroupApp.SystemServiceGroup.UserService
	initDBService            = service.ServiceGroupApp.SystemServiceGroup.InitDBService
	casbinService            = service.ServiceGroupApp.SystemServiceGroup.CasbinService
	autoCodeService          = service.ServiceGroupApp.SystemServiceGroup.AutoCodeService
	baseMenuService          = service.ServiceGroupApp.SystemServiceGroup.BaseMenuService
	authorityService         = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	dictionaryService        = service.ServiceGroupApp.SystemServiceGroup.DictionaryService
	systemConfigService      = service.ServiceGroupApp.SystemServiceGroup.SystemConfigService
	operationRecordService   = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService
	autoCodeHistoryService   = service.ServiceGroupApp.SystemServiceGroup.AutoCodeHistoryService
	dictionaryDetailService  = service.ServiceGroupApp.SystemServiceGroup.DictionaryDetailService
	authorityBtnService      = service.ServiceGroupApp.SystemServiceGroup.AuthorityBtnService
	machineRoomRecordService = service.ServiceGroupApp.SystemServiceGroup.MachineRoomRecordService
	jobPlatform              = service.ServiceGroupApp.SystemServiceGroup.JobPlatformService
	hostGroupService         = service.ServiceGroupApp.SystemServiceGroup.HostGroupService
	hostRecordService        = service.ServiceGroupApp.SystemServiceGroup.HostRecordService
	warnManageService        = service.ServiceGroupApp.SystemServiceGroup.WarnManageService
	hostMonitorRecordService = service.ServiceGroupApp.SystemServiceGroup.HostMonitorRecordService
	monitorCenterServer      = service.ServiceGroupApp.SystemServiceGroup.MonitorCenterServer
	hostTestRecordService    = service.ServiceGroupApp.SystemServiceGroup.HostTestRecordService
	hostPatrolRecordService  = service.ServiceGroupApp.SystemServiceGroup.HostPatrolRecordService
	patrolConfigService      = service.ServiceGroupApp.SystemServiceGroup.PatrolConfigService
)
