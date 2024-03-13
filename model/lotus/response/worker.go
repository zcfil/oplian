package response

import "time"

type WorkerInfo struct {
	Id           string    `json:"id"`
	OpId         string    `json:"opId"`
	GateId       string    `json:"gateId"`
	RoomId       string    `json:"roomId"`
	RoomName     string    `json:"roomName"`
	HostName     string    `json:"hostName"`
	DeviceSN     string    `json:"deviceSN"`
	CPUModel     string    `json:"CPUModel"`
	MemorySize   int       `json:"memorySize"`
	Online       bool      `json:"online" gorm:"-"`
	Ip           string    `json:"ip"`
	MinerIp      string    `json:"minerIp"`
	MinerId      uint64    `json:"minerId"`
	Actor        string    `json:"actor"`
	SectorSize   uint64    `json:"sectorSize"`
	Port         string    `json:"port"`
	DeployStatus int       `json:"deployStatus"`
	ErrMsg       string    `json:"errMsg"`
	RunStatus    int       `json:"runStatus"`
	RunTaskCount int       `json:"runTaskCount"`
	StartAt      time.Time `json:"startAt"`
	FinishAt     time.Time `json:"finishAt"`
}

type RelationWorker struct {
	Id       string `json:"id"`
	OpId     string `json:"opId"`
	RoomId   string `json:"roomId"`
	GateId   string `json:"gateId"`
	RoomName string `json:"roomName"`
	HostName string `json:"hostName"`
	DeviceSN string `json:"deviceSN"`
	Ip       string `json:"ip"`
}

type WorkerConfig struct {
	Id           string  `json:"id"`
	OpId         string  `json:"opId"`
	GateId       string  `json:"gateId"`
	RoomId       string  `json:"roomId"`
	RoomName     string  `json:"roomName"`
	HostName     string  `json:"hostName"`
	AssetNumber  string  `json:"assetNumber"`
	DeviceSN     string  `json:"deviceSN"`
	CPUCoreNum   int     `json:"cpuCoreNum"`
	CPUModel     string  `json:"cpuModel"`
	MemorySize   int     `json:"memorySize"`
	DiskNum      int     `json:"diskNum"`
	DiskSize     float64 `json:"diskSize"`
	Ip           string  `json:"ip"`
	PreCount1    int     `json:"preCount1"`
	PreCount2    int     `json:"preCount2"`
	Actor        string  `json:"actor"`
	OnOff1       bool    `json:"onOff1"`
	DeployStatus int     `json:"deployStatus"`
}

type WorkerMonitorInfo struct {
	OpId            string  `json:"opId"`
	GateId          string  `json:"gateId"`
	Ip              string  `json:"ip"`
	Actor           string  `json:"actor"`
	WorkerStatus    bool    `json:"workerStatus" gorm:"-"`
	CPUUseRate      float64 `json:"cpuUseRate" gorm:"-"`
	SysDiskUseRate  string  `json:"sysDiskUseRate" gorm:"-"`
	SysDiskLeave    string  `json:"sysDiskLeave" gorm:"-"`
	MainDiskUseRate string  `json:"mainDiskUseRate" gorm:"-"`
	MainDiskLeave   string  `json:"mainDiskLeave" gorm:"-"`
	DiskStatus      bool    `json:"diskStatus" gorm:"-"`
	GPUStatus       bool    `json:"gpuStatus" gorm:"-"`
	MountStatus     bool    `json:"mountStatus" gorm:"-"`
}
