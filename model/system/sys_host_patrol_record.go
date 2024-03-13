// 自动生成模板SysHostTPatrolRecord
package system

import (
	"oplian/global"
)

// 如果含有time.Time 请自行import time包
type SysHostPatrolRecord struct {
	global.ZC_MODEL
	HostUUID                    string  `json:"hostUUID" form:"hostUUID" gorm:"column:host_uuid;comment:资产编号UUID"`
	PatrolType                  int64   `json:"patrolType" form:"patrolType" gorm:"column:patrol_type;comment:巡检类型, 1 节点机巡检, 2 存储机巡检, 3 worker机巡检"`
	PatrolBeginAt               int64   `json:"patrolBeginAt" form:"patrolBeginAt" gorm:"column:patrol_begin_at;comment:巡检开始时间 (时间戳)"`
	PatrolEndAt                 int64   `json:"patrolEndAt" form:"patrolEndAt" gorm:"column:patrol_end_at;comment:巡检结束时间 (时间戳)"`
	PatrolResult                int64   `json:"patrolResult" form:"patrolResult" gorm:"column:patrol_result;comment:巡检状态, 1 巡检中, 2 正常, 3 异常, 4 巡检失败"`
	PatrolMode                  int64   `json:"patrolMode" form:"patrolMode" gorm:"column:patrol_mode;comment:巡检方式, 1 自动触发, 2 手动触发"`
	DiskIO                      bool    `json:"diskIO" form:"diskIO" gorm:"column:disk_io;comment:磁盘有无IO, 0 无, 1 有"`
	DiskIODuration              string  `json:"diskIODuration" form:"diskIODuration" gorm:"column:disk_io_duration;comment:磁盘IO测试运行耗时"`
	HostIsDown                  bool    `json:"hostIsDown" form:"hostIsDown" gorm:"column:host_is_down;comment:主机是否宕机,true 已宕机 false 未宕机"`
	HostIsDownDuration          string  `json:"hostIsDownDuration" form:"hostIsDownDuration" gorm:"column:host_is_down_duration;comment:主机是否宕机测试耗时"`
	HostNetStatus               bool    `json:"hostNetStatus" form:"hostNetStatus" gorm:"column:host_net_status;comment:主机网络是否正常,true 正常 false 异常"`
	HostNetDuration             string  `json:"hostNetDuration" form:"hostNetDuration" gorm:"column:host_net_duration;comment:主机网络测试耗时"`
	LogInfoStatus               bool    `json:"logInfoStatus" form:"logInfoStatus" gorm:"column:log_info_status;comment:主机日志消息是否正常,true 正常 false 异常"`
	LogInfoDuration             string  `json:"logInfoDuration" form:"logInfoDuration" gorm:"column:log_info_duration;comment:主机日志消息测试耗时"`
	LogOvertimeStatus           bool    `json:"logOvertimeStatus" form:"logOvertimeStatus" gorm:"column:log_overtime_status;comment:主机日志是否超时,true 正常 false 异常"`
	LogOvertimeDuration         string  `json:"logOvertimeDuration" form:"logOvertimeDuration" gorm:"column:log_overtime_duration;comment:主机日志超时测试耗时"`
	WalletBalanceStatus         bool    `json:"walletBalanceStatus" form:"walletBalanceStatus" gorm:"column:wallet_balance_status;comment:wdpost钱包余额是否正常,true 正常 false 异常"`
	WalletBalance               float64 `json:"walletBalance" form:"walletBalance" gorm:"column:wallet_balance;comment:wdpost钱包余额"`
	WalletBalanceDuration       string  `json:"walletBalanceDuration" form:"walletBalanceDuration" gorm:"column:wallet_balance_duration;comment:wdpost钱包余额测试耗时"`
	LotusSyncStatus             bool    `json:"lotusSyncStatus" form:"lotusSyncStatus" gorm:"column:lotus_sync_status;comment:lotus高度同步是否正常,true 正常 false 异常"`
	LotusSyncDuration           string  `json:"lotusSyncDuration" form:"lotusSyncDuration" gorm:"column:lotus_sync_duration;comment:lotus高度同步测试耗时"`
	GPUDriveStatus              bool    `json:"gpuDriveStatus" form:"gpuDriveStatus" gorm:"column:gpu_drive_status;comment:GPU驱动是否正常,true 正常 false 异常"`
	GPUDriveDuration            string  `json:"gpuDriveDuration" form:"gpuDriveDuration" gorm:"column:gpu_drive_duration;comment:GPU驱动测试耗时"`
	PackageVersionStatus        bool    `json:"packageVersionStatus" form:"packageVersionStatus" gorm:"column:package_version_status;comment:程序包版本是否正常,true 正常 false 异常"`
	PackageVersion              string  `json:"packageVersion" form:"packageVersion" gorm:"type:text;column:package_version;comment:程序包版本号"`
	PackageVersionDuration      string  `json:"packageVersionDuration" form:"packageVersionDuration" gorm:"type:text;column:package_version_duration;comment:程序包版本测试耗时"`
	DataCatalogStatus           bool    `json:"dataCatalogStatus" form:"dataCatalogStatus" gorm:"column:data_catalog_status;comment:lotus与lotusminer数据目录是否正常,true 正常 false 异常"`
	DataCatalogDuration         string  `json:"dataCatalogDuration" form:"dataCatalogDuration" gorm:"column:data_catalog_duration;comment:lotus与lotusminer数据目录测试耗时"`
	EnvironmentVariableStatus   bool    `json:"environmentVariableStatus" form:"environmentVariableStatus" gorm:"column:environment_variable_status;comment:环境变量是否正常,true 正常 false 异常"`
	EnvironmentVariableDuration string  `json:"environmentVariableDuration" form:"environmentVariableDuration" gorm:"column:environment_variable_duration;comment:环境变量测试耗时"`
	BlockLogStatus              bool    `json:"blockLogStatus" form:"blockLogStatus" gorm:"column:block_log_status;comment:出块日志是否正常,true 正常 false 异常"`
	BlockLogDuration            string  `json:"blockLogDuration" form:"blockLogDuration" gorm:"column:block_log_duration;comment:出块日志测试耗时"`
	TimeSyncStatus              bool    `json:"timeSyncStatus" form:"timeSyncStatus" gorm:"column:time_sync_status;comment:时间同步是否正常,true 正常 false 异常"`
	TimeSyncDuration            string  `json:"timeSyncDuration" form:"timeSyncDuration" gorm:"column:time_sync_duration;comment:时间同步测试耗时"`
	PingNetStatus               bool    `json:"pingNetStatus" form:"pingNetStatus" gorm:"column:ping_net_status;comment:ping网络是否正常,true 正常 false 异常"`
	PingNetDuration             string  `json:"pingNetDuration" form:"pingNetDuration" gorm:"column:ping_net_duration;comment:ping网络测试耗时"`
}

func (SysHostPatrolRecord) TableName() string {
	return "sys_host_patrol_records"
}
