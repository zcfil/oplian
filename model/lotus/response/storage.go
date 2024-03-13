package response

import "time"

type StorageInfo struct {
	Id           string    `json:"id"`
	OpId         string    `json:"opId"`
	GateId       string    `json:"gateId"`
	RoomId       string    `json:"roomId"`
	RoomName     string    `json:"roomName"`
	HostName     string    `json:"hostName"`
	DeviceSN     string    `json:"deviceSN"`
	Ip           string    `json:"ip"`
	MinerId      uint64    `json:"minerId"`
	DeployStatus int       `json:"deployStatus"`
	ColonyName   string    `json:"colonyName"`
	ColonyType   int       `json:"colonyType"`
	StartAt      time.Time `json:"startAt"`
	FinishAt     time.Time `json:"finishAt"`
	ErrMsg       string    `json:"errMsg"`
}

type StorageMonitorInfo struct {
	OpId           string  `json:"opId"`
	GateId         string  `json:"gateId"`
	Ip             string  `json:"ip"`
	ColonyName     string  `json:"-" gorm:"colony_name"`
	ColonyType     int     `json:"-" gorm:"colony_type"`
	Actor          string  `json:"actor" gorm:"-"`
	StorageType    int     `json:"storageType" gorm:"-"`
	CPUUseRate     float64 `json:"cpuUseRate" gorm:"-"`
	SysDiskUseRate string  `json:"sysDiskUseRate" gorm:"-"`
	SysDiskLeave   string  `json:"sysDiskLeave" gorm:"-"`
	DiskStatus     bool    `json:"diskStatus" gorm:"-"`
	NfsStatus      bool    `json:"nfsStatus" gorm:"-"`
}

type DCStorageMonitorInfo struct {
	UUID           string  `json:"opId" gorm:"uuid"`
	GatewayId      string  `json:"gateId" gorm:"gateway_id"`
	IntranetIp     string  `json:"ip" gorm:"intranet_ip"`
	CPUUseRate     float64 `json:"cpuUseRate" gorm:"-"`
	SysDiskUseRate string  `json:"sysDiskUseRate" gorm:"-"`
	SysDiskLeave   string  `json:"sysDiskLeave" gorm:"-"`
	DiskStatus     bool    `json:"diskStatus" gorm:"-"`
	NfsStatus      bool    `json:"nfsStatus" gorm:"-"`
}

type StorageMountErrorList struct {
	OpId       string `json:"opId" gorm:"op_id"`
	GateId     string `json:"gateId" gorm:"gate_id"`
	Ip         string `json:"ip" gorm:"ip"`
	ColonyName string `json:"colonyName" gorm:"colony_name"`
}
