package lotus

import (
	"time"
)

type LotusSectorLog struct {
	ID           string    `gorm:"primarykey"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP(3)"`
	SectorId     uint64    `json:"sectorId" gorm:"index;comment:扇区ID"`
	Actor        string    `json:"actor"  gorm:"comment:节点ID"`
	SectorStatus string    `json:"sectorStatus" gorm:"comment:扇区状态"`
	ErrorMsg     string    `json:"errorMsg" gorm:"type:text;comment:错误信息"`
	WorkerId     string    `json:"workerId"  gorm:"comment:worker UUID"`
	WorkerIp     string    `json:"WorkerIp" gorm:"comment:worker IP"`
	FinishAt     time.Time `json:"finishAt" gorm:"default:null;comment:结束时间"`
}

func (LotusSectorLog) TableName() string {
	return "lotus_sector_log"
}
