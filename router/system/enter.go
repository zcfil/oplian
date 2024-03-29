package system

type RouterGroup struct {
	ApiRouter
	SysRouter
	BaseRouter
	InitRouter
	MenuRouter
	UserRouter
	CasbinRouter
	AuthorityRouter
	DictionaryRouter
	OperationRecordRouter
	DictionaryDetailRouter
	AuthorityBtnRouter
	MachineRoomRecordRouter
	JobPlatformRouter
	HostGroupRouter
	HostRecordRouter
	WarnManageRouter
	HostMonitorRecordRouter
	HostTestRecordRouter
	HostPatrolRecordRouter
	PatrolConfigRouter
	HomeInterfaceRouter
	EquipMonitorRouter
	AppStoreRouter
}
