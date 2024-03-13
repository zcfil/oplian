package request

import (
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/lotus"
)

type WorkerCluster struct {
	request.PageInfo
	OpKeyWord   string `json:"opKeyWord"`
	RoomKeyWord string `json:"roomKeyWord"`
	SectorSize  int    `json:"sectorSize"`
}

type AddWorkerCluster struct {
	WorkerCluster []lotus.LotusWorkerCluster `json:"workerCluster"`
}

type WorkerOp struct {
	ID         int    `json:"id" form:"id"`
	YearMonth  string `json:"yearMonth" form:"yearMonth"`
	OpKeyWord  string `json:"opKeyWord" form:"opKeyWord"`
	RoomId     string `json:"roomId" form:"roomId"`
	SectorSize string `json:"sectorSize" form:"sectorSize"`
	TaskStatus string `json:"taskStatus" form:"taskStatus"`
	request.PageInfo
}

type SealerParam struct {
	GateWayId string `json:"gateWayId"`
	OpId      string `json:"opId"`
	OpC2Id    string `json:"opC2Id"`
	Number    uint64 `json:"number"`
	Miner     string `json:"miner"`
}

type Commit2In struct {
	SectorNum  int64
	Phase1Out  []byte
	SectorSize uint64
}

type SectorsRecover struct {
	request.PageInfo
	SectorId      string `json:"sectorId"`
	SectorSize    int    `json:"sectorSize"`
	SectorType    int    `json:"sectorType"`
	SectorStatus  int    `json:"sectorStatus"`
	BelongingNode string `json:"belongingNode"`
}

type Ids struct {
	Id string `json:"id"`
}

type SectorsRecoverTask struct {
	global.ZC_MODEL
	TaskName          string `json:"taskName"`
	Actor             string `json:"actor"`
	Ids               []Ids  `json:"ids"`
	SectorType        int    `json:"sectorType"`
	SectorSize        int    `json:"sectorSize"`
	SectorTotal       int    `json:"sectorTotal"`
	OriginalValueOpId string `json:"originalValueOpId"`
	OriginalValueDir  string `json:"originalValueDir"`
	StorageOpId       string `json:"storageOpId"`
	StorageOpIp       string `json:"storageOpIp"`
	StorageOpName     string `json:"StorageOpName"`
	WorkerOp          []Ids  `json:"workerOp"`
}

type LotusSectorTask struct {
	global.ZC_MODEL
	TaskName          string `json:"taskName"`
	Actor             string `json:"actor"`
	SectorType        int    `json:"sectorType"`
	SectorTotal       int    `json:"sectorTotal"`
	SectorSize        int    `json:"sectorSize"`
	OriginalValueOpId string `json:"originalValueOpId"`
	OriginalValueDir  string `json:"originalValueDir"`
	StorageOpId       string `json:"storageOpId"`
	StorageOpIp       string `json:"storageOpIp"`
	StorageOpName     string `json:"StorageOpName"`
	TaskStatus        int    `json:"taskStatus"`
}

type StoragePaths struct {
	StoragePaths []Path `json:"StoragePaths"`
}

type Path struct {
	Path string `json:"Path"`
}

type DirFile struct {
	GateWayId string `json:"gateWayId"`
	Ip        string `json:"ip"`
	OpId      string `json:"opId"`
	Path      string `json:"path"`
}

type SectorStatus struct {
	Id     int `json:"id"`
	Status int `json:"status"` //1.进行中，2.暂停中,3.已完成,4.已终止
}

type WorkerInfo struct {
	GateWayId  string `json:"gateWayId"`
	WorkerType string `json:"workerType"`
	KeyWord    string `json:"keyWord"`
}

type C2TaskInfo struct {
	Miner      string `json:"miner"`
	Number     int    `json:"number"`
	DelType    int    `json:"delType"`
	SectorSize int    `json:"sectorSize"`
}

func (LotusSectorTask) TableName() string {
	return "lotus_sector_task"
}
