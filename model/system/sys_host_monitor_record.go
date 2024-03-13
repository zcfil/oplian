// 自动生成模板SysHostInfo
package system

import (
	"oplian/global"
)

// 如果含有time.Time 请自行import time包
type SysHostMonitorRecord struct {
	global.ZC_MODEL
	HostUUID       string  `json:"hostUUID" form:"hostUUID" gorm:"column:host_uuid;comment:资产编号UUID"`
	CPUUseRate     float32 `json:"cpuUseRate" form:"cpuUseRate" gorm:"column:cpu_use_rate;comment:CPU使用率"`
	DiskUseRate    float32 `json:"diskUseRate" form:"diskUseRate" gorm:"column:disk_use_rate;comment:硬盘使用率"`
	MemoryUseRate  float32 `json:"memoryUseRate" form:"memoryUseRate" gorm:"column:memory_use_rate;comment:内存使用率"`
	GPUUseRate     float32 `json:"gpuUseRate" form:"gpuUseRate" gorm:"column:gpu_use_rate;comment:GPU使用率"`
	GPUID          string  `json:"gpuId" form:"gpuId" gorm:"column:gpu_id;comment:GPU的编号id"`
	CPUTemperature string  `json:"cpuTemperature" form:"cpuTemperature" gorm:"column:cpu_temperature;comment:CPU温度"`
	DiskSize       string  `json:"diskSize" form:"diskSize" gorm:"column:disk_size;comment:磁盘大小"`
	DiskUseSize    string  `json:"diskUseSize" form:"diskUseSize" gorm:"column:disk_use_size;comment:磁盘使用大小"`
	MemorySize     int64   `json:"memorySize" form:"memorySize" gorm:"column:memory_size;comment:内存大小"`
	MemoryUseSize  int64   `json:"memoryUseSize" form:"memoryUseSize" gorm:"column:memory_use_size;comment:内存使用大小"`
}

func (SysHostMonitorRecord) TableName() string {
	return "sys_host_monitor_records"
}
