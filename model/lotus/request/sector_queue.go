package request

type SectorQueue struct {
	TaskName    string `json:"taskName" gorm:"comment:任务名称"`
	SectorType  int    `json:"sectorType"  gorm:"comment:CC：1，DC：2"`
	SectorSize  int    `json:"sectorSize"  gorm:"comment:32,64"`
	Actor       string `json:"actor" gorm:"index;comment:节点号"`
	GateId      string `json:"gateId" gorm:"comment:服务器gateID"`
	OpId        string `json:"opId" gorm:"comment:服务器ID"`
	SectorTotal int    `json:"sectorTotal" gorm:"comment:扇区总数"`
}

type DealInfo struct {
	DealUuid  string `json:"dealUuid"`
	PieceCid  string `json:"pieceCid"`
	EndEpoch  int64  `json:"endEpoch"`
	CarPath   string `json:"carPath"`
	JobStatus int    `json:"jobStatus"`
	FileOpId  string `json:"fileOpId"`
}

type DealPage struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	Actor      string `json:"actor" form:"actor"`
	QueueId    uint   `json:"queueId" form:"queueId"`
	Status     int    `json:"status" form:"status"`
	SectorType uint   `json:"sectorType" form:"sectorType"`
}

type DealCarInfo struct {
	PieceCid string   `json:"pieceCid" form:"pieceCid"`
	GateId   string   `json:"gateId" form:"gateId"`
	OpIds    []string `json:"opIds" form:"opIds"`
	CarDir   string   `json:"carDir" form:"carDir"`
	ID       uint     `json:"id" form:"id"`
	Actor    string   `json:"actor" form:"actor"`
}

type SectorRecoverDetail struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Id       int `json:"id" form:"id"`
}
