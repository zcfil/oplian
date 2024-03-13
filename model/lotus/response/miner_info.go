package response

import "time"

type MinerInfo struct {
	Id            uint      `json:"id"`
	Bid           uint      `json:"bid"`
	LotusId       uint      `json:"lotusId"`
	OpId          string    `json:"opId"`
	GateId        string    `json:"gateId"`
	RoomId        string    `json:"roomId"`
	RoomName      string    `json:"roomName"`
	HostName      string    `json:"hostName"`
	DeviceSN      string    `json:"deviceSN"`
	CPUModel      string    `json:"CPUModel"`
	MemorySize    int       `json:"memorySize"`
	Ip            string    `json:"ip"`
	InternetIp    string    `json:"internetIp"`
	Actor         string    `json:"actor"`
	SectorSize    uint64    `json:"sectorSize"`
	Port          string    `json:"port"`
	DeployStatus  int       `json:"deployStatus"`
	IsManage      bool      `json:"isManage"`
	IsWdpost      bool      `json:"isWdpost"`
	IsWnpost      bool      `json:"isWnpost"`
	Partitions    string    `json:"partitions"`
	StorageCount  int       `json:"storageCount"`
	WorkerCount   int       `json:"workerCount"`
	Online        bool      `json:"online" gorm:"-"`
	LotusIp       string    `json:"lotusIp"`
	LotusHostName string    `json:"lotusHostName"`
	LotusToken    string    `json:"lotusToken"`
	WalletCount   int       `json:"walletCount"`
	Wallets       []Wallet  `json:"wallets" gorm:"-"`
	StartAt       time.Time `json:"startAt"`
	FinishAt      time.Time `json:"finishAt"`
	ErrMsg        string    `json:"errMsg"`
	ColonyType    int       `json:"colonyType"`
	DcQuotaWallet string    `json:"dcQuotaWallet"`
	AddType       int       `json:"addType"`
}

type ActorInfo struct {
	Actor      string `json:"actor"`
	GateId     string `json:"gateId"`
	SectorSize int    `json:"sectorSize"`
	ColonyType int    `json:"colonyType"`
	MinerFile  bool   `json:"minerFile" gorm:"-"`
}

type ManageInfo struct {
	Actor         string  `json:"actor"`
	SectorSize    int     `json:"sectorSize"`
	GateId        string  `json:"gateId"`
	OpId          string  `json:"opId"`
	Ip            string  `json:"ip"`
	LotusId       uint64  `json:"lotusId"`
	WorkerWallet  string  `json:"workerWallet" gorm:"-"`
	WorkerBalance float64 `json:"workerBalance" gorm:"-"`
}

type MinerMonitorInfo struct {
	OpId              string  `json:"opId"`
	GateId            string  `json:"gateId"`
	Ip                string  `json:"ip"`
	Actor             string  `json:"actor"`
	LotusStatus       bool    `json:"lotusStatus" gorm:"-"`
	MinerStatus       bool    `json:"minerStatus" gorm:"-"`
	BoostStatus       bool    `json:"boostStatus" gorm:"-"`
	CPUUseRate        float64 `json:"cpuUseRate" gorm:"-"`
	SysDiskUseRate    string  `json:"sysDiskUseRate" gorm:"-"`
	SysDiskLeave      string  `json:"sysDiskLeave" gorm:"-"`
	MainDiskUseRate   string  `json:"mainDiskUseRate" gorm:"-"`
	MainDiskLeave     string  `json:"mainDiskLeave" gorm:"-"`
	DiskStatus        bool    `json:"diskStatus" gorm:"-"`
	GPUStatus         bool    `json:"gpuStatus" gorm:"-"`
	MountStatus       bool    `json:"mountStatus" gorm:"-"`
	LotusHeightStatus bool    `json:"lotusHeightStatus" gorm:"-"`
}

type NodeMinerInfoResp struct {
	MinerAttribute     string `json:"minerAttribute"`
	MinerProcessStatus bool   `json:"minerProcessStatus"`
	WinRemainingTime   string `json:"winRemainingTime"`
	MessageOut         bool   `json:"messageOut"`
	ProofRequiredTime  string `json:"proofRequiredTime"`
	RestartStatus      bool   `json:"restartStatus"`
	DailyExplosiveNum  int    `json:"dailyExplosiveNum"`
}

type MinerSelectInfoResp struct {
	Id     uint   `json:"id"`     
	OpId   string `json:"opId"`   
	GateId string `json:"gateId"` 
	Ip     string `json:"ip"`     
	Port   string `json:"port"`   
	Actor  string `json:"actor"`  
}
