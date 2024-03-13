package response

import (
	"time"
)

type SysHostMonitorRecord struct {
	ID            uint      `json:"ID" gorm:"column:id"`
	CreatedAt     time.Time `json:"CreatedAt" gorm:"column:created_at"`
	HostName      string    `json:"hostName" gorm:"-"`
	IntranetIP    string    `json:"intranetIp" gorm:"-"`
	InternetIP    string    `json:"internetIp" gorm:"-"`
	HostUUID      string    `json:"hostUUID" gorm:"column:host_uuid"`
	CPUUseRate    float32   `json:"cpuUseRate" gorm:"column:cpu_use_rate"`
	DiskUseRate   float32   `json:"diskUseRate" gorm:"column:disk_use_rate"`
	MemoryUseRate float32   `json:"memoryUseRate" gorm:"column:memory_use_rate"`
	GPUUseRate    float32   `json:"gpuUseRate" gorm:"column:gpu_use_rate"`
	GPUID         string    `json:"gpuId" gorm:"column:gpu_id"`
}

type SysHostMonitorLineChart struct {
	HostName string            `json:"hostName"`
	HostUUID string            `json:"hostUUID"`
	GPUID    string            `json:"gpuId"`
	HostInfo []HostMonitorRate `json:"hostInfo"`
}

type HostMonitorRate struct {
	CPUUseRate    float32   `json:"cpuUseRate" gorm:"column:cpu_use_rate"`
	DiskUseRate   float32   `json:"diskUseRate" gorm:"column:disk_use_rate"`
	MemoryUseRate float32   `json:"memoryUseRate" gorm:"column:memory_use_rate"`
	GPUUseRate    float32   `json:"gpuUseRate" gorm:"column:gpu_use_rate"`
	GPUID         string    `json:"gpuId" gorm:"column:gpu_id"`
	CreatedAt     time.Time `json:"CreatedAt" gorm:"column:created_at"`
}

type StorageRateResp struct {
	HostName        string  `json:"hostName"`
	IntranetIP      string  `json:"intranetIp"`
	InternetIP      string  `json:"internetIp"`
	DiskUseRate     float32 `json:"diskUseRate"`
	DiskAllSize     float64 `json:"diskAllSize"`
	DiskAllSizeUnit string  `json:"diskAllSizeUnit"`
	DiskUseSize     float64 `json:"diskUseSize"`
	DiskUseSizeUnit string  `json:"diskUseSizeUnit"`
	MinerID         string  `json:"minerId"`
}
