// 自动生成模板SysHostInfo
package system

import (
	"oplian/global"
)

// 如果含有time.Time 请自行import time包
type SysHostRecord struct {
	global.ZC_MODEL
	HostName         string  `json:"hostName" form:"hostName" gorm:"column:host_name;comment:主机名称"`
	IntranetIP       string  `json:"intranetIp" form:"intranetIp" gorm:"column:intranet_ip;comment:内网IP"`
	InternetIP       string  `json:"internetIp" form:"internetIp" gorm:"column:internet_ip;comment:外网IP"`
	UUID             string  `json:"uuid" gorm:"index;unique;comment:资产编号UUID"`
	DeviceSN         string  `json:"deviceSN" form:"deviceSN" gorm:"column:device_sn;comment:设备SN"`
	HostManufacturer string  `json:"hostManufacturer" form:"hostManufacturer" gorm:"column:host_manufacturer;comment:主机厂商"`
	HostModel        string  `json:"hostModel" form:"hostModel" gorm:"column:host_model;comment:主机型号"`
	OperatingSystem  string  `json:"operatingSystem" form:"operatingSystem" gorm:"column:operating_system;comment:操作系统"`
	CPUCoreNum       int     `json:"cpuCoreNum" form:"cpuCoreNum" gorm:"column:cpu_core_num;comment:CPU逻辑核心数"`
	CPUModel         string  `json:"cpuModel" form:"cpuModel" gorm:"column:cpu_model;comment:CPU型号"`
	MemorySize       int     `json:"memorySize" form:"memorySize" gorm:"column:memory_size;comment:内存大小(GB)"`
	DiskNum          int     `json:"diskNum" form:"diskNum" gorm:"column:disk_num;comment:磁盘数量"`
	DiskSize         float64 `json:"diskSize" form:"diskSize" gorm:"column:disk_size;comment:磁盘容量(TB)"`
	HostShelfLife    int     `json:"hostShelfLife" form:"hostShelfLife" gorm:"column:host_shelf_life;comment:主机保质年限"`
	HostType         int     `json:"hostType" form:"hostType" gorm:"column:host_type;comment:主机类型, 1 物理机, 2 虚拟机"`
	HostClassify     int     `json:"hostClassify" form:"hostClassify" gorm:"column:host_classify;comment:主机分类,0 默认值未赋值类型, 1 miner机, 2 worker机, 3 存储机, 4 DC主机, 5 C2 worker机, 6 单lotus机(中间状态)"`
	ServerDNS        string  `json:"serverDns" form:"serverDns" gorm:"column:server_dns;comment:DNS服务器"`
	SubnetMask       string  `json:"subnetMask" form:"subnetMask" gorm:"column:subnet_mask;comment:子网掩码"`
	Gateway          string  `json:"gateway" form:"gateway" gorm:"column:gateway;comment:网关"`
	HostGroupId      int     `json:"hostGroupId" form:"hostGroupId" gorm:"column:host_group_id;comment:主机所属分组ID"`
	RoomId           string  `json:"roomId" form:"roomId" gorm:"column:room_id;comment:机房编号"`
	MonitorTime      int     `json:"monitorTime" form:"monitorTime" gorm:"column:monitor_time;comment:主机监控信息更新时间"`
	GatewayId        string  `json:"gatewayId" form:"gatewayId" gorm:"column:gateway_id;comment:oplian-gateway对应的网关ID"`
	GPUNum           int     `json:"gpuNum" form:"gpuNum" gorm:"column:gpu_num;comment:主机GPU数量"`
	AssetNumber      string  `json:"assetNumber" form:"assetNumber" gorm:"column:asset_number;comment:资产编号"`
	SystemVersion    string  `json:"systemVersion" form:"systemVersion" gorm:"column:system_version;comment:系统版本号"`
	SystemBits       int     `json:"systemBits" form:"systemBits" gorm:"column:system_bits;comment:系统位数"`
	RoomName         string  `json:"roomName" form:"roomName" gorm:"column:room_name;comment:机房名称"`
	IsGroupArray     bool    `json:"isGroupArray" form:"isGroupArray" gorm:"column:is_group_array;comment:是否组阵列,true为组阵列, false未组阵列"`
	NetOccupyTime    int64   `json:"netOccupyTime" form:"netOccupyTime" gorm:"column:net_occupy_time;comment:该主机被使用于测试其余主机网络的时间"`
}

func (SysHostRecord) TableName() string {
	return "sys_host_records"
}
