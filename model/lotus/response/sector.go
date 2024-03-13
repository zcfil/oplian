package response

import (
	"oplian/global"
	model "oplian/model/lotus"
	"time"
)

type SectorInfo struct {
	global.ZC_MODEL
	SectorId     uint64    `json:"sectorId" gorm:"index;unique;comment:扇区ID"`
	Actor        string    `json:"actor"  gorm:"comment:节点ID"`
	SectorStatus string    `json:"sectorStatus"  gorm:"comment:扇区状态"`
	SectorType   int       `json:"sectorType"  gorm:"comment:扇区类型：1：CC，2：DC"`
	SectorSize   uint64    `json:"sectorSize"  gorm:"comment:扇区大小"`
	PreCid       string    `json:"preCid"  gorm:"comment:P2消息"`
	FinishAt     time.Time `json:"finishAt" gorm:"default:null;comment:扇区完成时间"`
}

type SectorDetails struct {
	SectorInfo SectorInfo             `json:"sectorInfo"`
	Piece      []string               `json:"piece"`
	SectorLog  []model.LotusSectorLog `json:"sectorLog"`
}
