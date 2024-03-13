package response

import "oplian/utils"

type SysHostPatrolRecord struct {
	ID            uint   `json:"ID" gorm:"column:id"`
	HostUUID      string `json:"hostUUID" gorm:"column:host_uuid"`
	PatrolType    int    `json:"patrolType" gorm:"column:patrol_type"`
	PatrolBeginAt int64  `json:"-" gorm:"column:patrol_begin_at"`
	PatrolEndAt   int64  `json:"-" gorm:"column:patrol_end_at"`
	BeginTime     string `json:"beginTime" gorm:"-"`
	EndTime       string `json:"endTime" gorm:"-"`
	PatrolResult  int    `json:"patrolResult" gorm:"column:patrol_result"`
	PatrolMode    int    `json:"patrolMode" gorm:"column:patrol_mode"`
	// 需要查询其余表获取的内容
	HostName       string `json:"hostName" gorm:"-"`
	InternetIP     string `json:"internetIp" gorm:"-"`
	IntranetIP     string `json:"intranetIp" gorm:"-"`
	DeviceSN       string `json:"deviceSN" gorm:"-"`
	AssetNumber    string `json:"assetNumber" gorm:"-"`
	RoomId         string `json:"roomId" gorm:"-"`
	RoomName       string `json:"roomName" gorm:"-"`
	PatrolTakeTime string `json:"patrolTakeTime" gorm:"-"`
}

type SysHostPatrolInfo struct {
	ID            uint   `json:"ID" gorm:"column:id"`
	HostUUID      string `json:"hostUUID" gorm:"column:host_uuid"`
	PatrolType    int    `json:"patrolType" gorm:"column:patrol_type"`
	PatrolBeginAt int64  `json:"-" gorm:"column:patrol_begin_at"`
	PatrolEndAt   int64  `json:"-" gorm:"column:patrol_end_at"`
	BeginTime     string `json:"beginTime" gorm:"-"`
	EndTime       string `json:"endTime" gorm:"-"`
	PatrolResult  int    `json:"patrolResult" gorm:"column:patrol_result"`
	PatrolMode    int    `json:"patrolMode" gorm:"column:patrol_mode"`
	// 需要查询其余表获取的内容
	HostName       string `json:"hostName" gorm:"-"`
	IntranetIP     string `json:"intranetIp" gorm:"-"`
	DeviceSN       string `json:"deviceSN" gorm:"-"`
	AssetNumber    string `json:"assetNumber" gorm:"-"`
	RoomId         string `json:"roomId" gorm:"-"`
	RoomName       string `json:"roomName" gorm:"-"`
	PatrolTakeTime string `json:"patrolTakeTime" gorm:"-"`
	// 详细信息
	DiskIO                 bool                      `json:"diskIO" gorm:"column:disk_io"`
	DiskIODuration         string                    `json:"diskIODuration" gorm:"column:disk_io_duration"`
	HostIsDown             bool                      `json:"hostIsDown" gorm:"column:host_is_down"`
	HostIsDownDuration     string                    `json:"hostIsDownDuration" gorm:"column:host_is_down_duration"`
	HostNetStatus          bool                      `json:"hostNetStatus" gorm:"column:host_net_status"`
	HostNetDuration        string                    `json:"hostNetDuration" gorm:"column:host_net_duration"`
	LogInfoStatus          bool                      `json:"logInfoStatus" gorm:"column:log_info_status"`
	LogInfoDuration        string                    `json:"logInfoDuration" gorm:"column:log_info_duration"`
	LogOvertimeStatus      bool                      `json:"logOvertimeStatus" gorm:"column:log_overtime_status"`
	LogOvertimeDuration    string                    `json:"logOvertimeDuration" gorm:"column:log_overtime_duration"`
	WalletBalanceStatus    bool                      `json:"walletBalanceStatus" gorm:"column:wallet_balance_status"`
	WalletBalance          float64                   `json:"walletBalance" gorm:"column:wallet_balance"`
	WalletBalanceDuration  string                    `json:"walletBalanceDuration" gorm:"column:wallet_balance_duration"`
	LotusSyncStatus        bool                      `json:"lotusSyncStatus" gorm:"column:lotus_sync_status"`
	LotusSyncDuration      string                    `json:"lotusSyncDuration" gorm:"column:lotus_sync_duration"`
	GPUDriveStatus         bool                      `json:"gpuDriveStatus" gorm:"column:gpu_drive_status"`
	GPUDriveDuration       string                    `json:"gpuDriveDuration" gorm:"column:gpu_drive_duration"`
	PackageVersionStatus   bool                      `json:"packageVersionStatus" gorm:"column:package_version_status"`
	PackageVersion         string                    `json:"-" gorm:"column:package_version"`
	HostPackageVersion     utils.LotusPackageVersion `json:"packageVersion" gorm:"-"`
	PackageVersionDuration string                    `json:"packageVersionDuration" gorm:"column:package_version_duration"`
	DataCatalogStatus      bool                      `json:"dataCatalogStatus" gorm:"column:data_catalog_status"`
	DataCatalogDuration    string                    `json:"dataCatalogDuration" gorm:"column:data_catalog_duration"`
	BlockLogStatus         bool                      `json:"blockLogStatus" gorm:"column:block_log_status"`
	BlockLogDuration       string                    `json:"blockLogDuration" gorm:"column:block_log_duration"`
	TimeSyncStatus         bool                      `json:"timeSyncStatus" gorm:"column:time_sync_status"`
	TimeSyncDuration       string                    `json:"timeSyncDuration" gorm:"column:time_sync_duration"`
	PingNetStatus          bool                      `json:"pingNetStatus" gorm:"column:ping_net_status"`
	PingNetDuration        string                    `json:"pingNetDuration" gorm:"column:ping_net_duration"`
}

type SysHostPatrolReport struct {
	ID       uint   `json:"ID" gorm:"column:id"`
	HostUUID string `json:"hostUUID" gorm:"column:host_uuid"`
}
