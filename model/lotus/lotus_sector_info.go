package lotus

import (
	"oplian/global"
	"time"
)

type LotusSectorInfo struct {
	global.ZC_MODEL
	SectorId     uint64    `json:"sectorId" gorm:"index;unique;comment:扇区ID"`
	Actor        string    `json:"actor"  gorm:"comment:节点ID"`
	SectorStatus string    `json:"sectorStatus"  gorm:"comment:扇区状态"`
	SectorType   int       `json:"sectorType"  gorm:"comment:扇区类型：1：CC，2：DC"`
	SectorSize   uint64    `json:"sectorSize"  gorm:"comment:扇区大小"`
	CidCommD     string    `json:"cidCommD" gorm:"comment:Unsealed"`
	CidCommR     string    `json:"cidCommR" gorm:"comment:Sealed"`
	Ticket       string    `json:"ticket" gorm:"comment:重做扇区字段"`
	TicketH      uint64    `json:"ticketH" gorm:"comment:ticket获取高度"`
	Seed         string    `json:"seed" gorm:"comment:WaitSeed"`
	SeedH        uint64    `json:"seedH" gorm:"comment:WaitSeed高度"`
	PreCid       string    `json:"preCid" gorm:"comment:P2 消息ID"`
	CommitCid    string    `json:"commitCid" gorm:"comment:C2 消息ID"`
	Proof        string    `json:"proof" gorm:"type:text;comment:C2证明结果"`
	FinishAt     time.Time `json:"finishAt" gorm:"default:null;comment:扇区完成时间"`
}

func (LotusSectorInfo) TableName() string {
	return "lotus_sector_info"
}
