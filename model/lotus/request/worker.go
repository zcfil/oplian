package request

//type RelationWorker struct {
//	MinerId    uint64
//	WorkerType int //worker类型：0任务，1存储
//}

type ResetWorker struct {
	IsClear bool
	LinkId  uint64
	Id      uint64
	GateId  string
	OpId    string
}

type WorkerInfoPage struct {
	Page         int    `json:"page" form:"page"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
	Keyword      string `json:"keyword" form:"keyword"`
	Actor        string `json:"actor" form:"actor"`
	GateId       string `json:"gateId" form:"gateId"`
	DeployStatus int    `json:"deployStatus" form:"deployStatus"`
}

type WorkerConfigPage struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
	GateId   string `json:"gateId"`
}

type WorkerPre struct {
	ID        uint64 `json:"id"`
	GateId    string `json:"gateId"`
	OpId      string `json:"opId"`
	PreCount1 int32  `json:"preCount1"`
	PreCount2 int32  `json:"PreCount2"`
}

type PreConfig struct {
	ID        uint64 `json:"id"`
	PreCount1 int32  `json:"preCount1"`
	PreCount2 int32  `json:"PreCount2"`
}
type WorkerOnOff struct {
	ID     uint64 `json:"id"`
	GateId string `json:"gateId"`
	OpId   string `json:"opId"`
	OnOff  bool   `json:"onOff"`
}
