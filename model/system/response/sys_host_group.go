package response

type SysHostGroup struct {
	ID        uint   `json:"ID" gorm:"column:id"`
	GroupName string `json:"groupName" gorm:"column:group_name"`
}

type OpInfo struct {
	HostName         string  `json:"hostName" `
	IntranetIP       string  `json:"intranetIp"`
	InternetIP       string  `json:"internetIp" `
	UUID             string  `json:"uuid" `
	DeviceSN         string  `json:"deviceSN" `
	HostManufacturer string  `json:"hostManufacturer"`
	HostModel        string  `json:"hostModel" `
	OperatingSystem  string  `json:"operatingSystem"`
	CPUCoreNum       int     `json:"cpuCoreNum" `
	CPUModel         string  `json:"cpuModel"`
	MemorySize       int     `json:"memorySize"`
	DiskNum          int     `json:"diskNum" `
	DiskSize         float64 `json:"diskSize" `
	HostShelfLife    int     `json:"hostShelfLife"`
	HostType         int     `json:"hostType" `
	HostClassify     int     `json:"hostClassify"`
	ServerDNS        string  `json:"serverDns" `
	SubnetMask       string  `json:"subnetMask" `
	Gateway          string  `json:"gateway" `
	HostGroupId      int     `json:"hostGroupId" `
	RoomId           string  `json:"roomId" `
	RoomName         string  `json:"roomName" `
	MonitorTime      int     `json:"monitorTime" `
	GatewayId        string  `json:"gatewayId"`
	GPUNum           int     `json:"gpuNum"`
	AssetNumber      string  `json:"assetNumber" `
	SystemVersion    string  `json:"systemVersion" `
	SystemBits       int     `json:"systemBits"`
}
