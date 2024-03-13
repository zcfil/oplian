package request

// PageInfo Paging common input parameter structure
type MinerInfoPage struct {
	Page         int    `json:"page" form:"page"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
	Keyword      string `json:"keyword" form:"keyword"`
	SyncStatus   int    `json:"syncStatus" form:"syncStatus"`
	DeployStatus int    `json:"deployStatus" form:"deployStatus"`
}
type MinerTypePage struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	Keyword    string `json:"keyword" form:"keyword"`
	GateId     string `json:"gateId" form:"gateId"`
	IsManage   bool   `json:"isManage" form:"isManage"`
	ColonyType int    `json:"colonyType" form:"colonyType"`
}

type MinerListReq struct {
	GateId string `json:"gateId" form:"gateId"` 
}

// 新增miner
type MinerInfo struct {
	Id         uint
	AddType    int
	GateId     string
	OpId       string
	LotusId    uint64
	Ip         string
	Actor      string
	Partitions string
	Owner      string
	IsManage   bool
	IsWdpost   bool
	IsWnpost   bool
	SectorSize uint64
	ColonyType int //1 NFS
}

// miner角色
type MinerParam struct {
	Id         uint
	LotusId    uint64
	Actor      string
	Partitions string
	IsManage   bool
	IsWdpost   bool
	IsWnpost   bool
}
