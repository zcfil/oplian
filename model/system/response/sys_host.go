package response

import (
	"time"
)

type SysHostRecord struct {
	ID               uint      `json:"ID" gorm:"column:id"`
	CreatedAt        time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
	UUID             string    `json:"uuid" gorm:"column:uuid"`
	HostName         string    `json:"hostName" gorm:"column:host_name"`
	HostType         int       `json:"hostType" gorm:"column:host_type"`
	HostClassify     int       `json:"hostClassify" gorm:"column:host_classify"`
	HostGroupId      int       `json:"hostGroupId" gorm:"column:host_group_id"`
	RoomId           string    `json:"roomId" gorm:"column:room_id"`
	RoomName         string    `json:"roomName" gorm:"room_name"`
	HostGroupName    string    `json:"hostGroupName" gorm:"-"`
	IntranetIP       string    `json:"intranetIp" gorm:"column:intranet_ip"`
	InternetIP       string    `json:"internetIp" gorm:"column:internet_ip"`
	DeviceSN         string    `json:"deviceSN" gorm:"column:device_sn"`
	HostManufacturer string    `json:"hostManufacturer" gorm:"column:host_manufacturer"`
	HostModel        string    `json:"hostModel" gorm:"column:host_model"`
	OperatingSystem  string    `json:"operatingSystem" gorm:"column:operating_system"`
	CPUCoreNum       int       `json:"cpuCoreNum" gorm:"column:cpu_core_num"`
	CPUModel         string    `json:"cpuModel" gorm:"column:cpu_model"`
	MemorySize       int       `json:"memorySize" gorm:"column:memory_size"`
	DiskNum          int       `json:"diskNum" gorm:"column:disk_num"`
	DiskSize         float64   `json:"diskSize" gorm:"column:disk_size"`
	HostShelfLife    int       `json:"hostShelfLife" gorm:"column:host_shelf_life"`
	ServerDNS        string    `json:"serverDns" gorm:"column:server_dns"`
	SubnetMask       string    `json:"subnetMask" gorm:"column:subnet_mask"`
	Gateway          string    `json:"gateway" gorm:"column:gateway"`
	GatewayId        string    `json:"gatewayId" gorm:"column:gateway_id"`
	AssetNumber      string    `json:"assetNumber" gorm:"column:asset_number"`
	SystemVersion    string    `json:"systemVersion" gorm:"column:system_version"`
	SystemBits       int       `json:"systemBits" gorm:"column:system_bits"`
}

type SysHostInfo struct {
	ID           uint    `json:"ID"`
	UUID         string  `json:"uuid"`
	AssetNumber  string  `json:"assetNumber"`
	DeviceSn     string  `json:"deviceSn"`
	HostName     string  `json:"hostName"`
	RoomId       string  `json:"roomId"`
	RoomName     string  `json:"roomName"`
	GateId       string  `json:"gateId" gorm:"column:gateway_id"`
	IntranetIP   string  `json:"intranetIp"`
	MemorySize   int     `json:"memorySize"`
	DiskSize     float64 `json:"diskSize"`
	HostClassify int     `json:"hostClassify"`
	GpuNum       int     `json:"gpuNum"`
}

type SysHost struct {
	UUID       string `json:"uuid" gorm:"column:uuid"`
	HostName   string `json:"hostName" gorm:"column:host_name"`
	InternetIP string `json:"internetIp" gorm:"column:internet_ip"`
	IntranetIP string `json:"intranetIp" gorm:"column:intranet_ip"`
}

type SysHostPageResult struct {
	List     interface{} `json:"list"`
	TotalNum int64       `json:"totalNum"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type SysHostBindList struct {
	UUID        string `json:"uuid" gorm:"column:uuid"`
	HostName    string `json:"hostName" gorm:"column:host_name"`
	InternetIP  string `json:"internetIp" gorm:"column:internet_ip"`
	IntranetIP  string `json:"intranetIp" gorm:"column:intranet_ip"`
	DeviceSN    string `json:"deviceSN" gorm:"column:device_sn"`
	AssetNumber string `json:"assetNumber" gorm:"column:asset_number"`
}

type SysHostRecordPatrol struct {
	UUID         string `json:"uuid" gorm:"column:uuid"`
	HostClassify int64  `json:"hostClassify" gorm:"column:host_classify"`
	IntranetIP   string `json:"intranetIp" gorm:"column:intranet_ip"`
	GatewayId    string `json:"gatewayId" gorm:"column:gateway_id"`
}
