package lotus

import (
	"oplian/global"
)

type LotusDeals struct {
	global.ZC_MODEL
	SectorId uint64 `json:"sectorId" gorm:"index;comment:扇区ID"`
	MinerId  string `json:"minerId"  gorm:"index;comment:节点ID"`
	PieceCid string `json:"pieceCid"  gorm:"comment:重做DC扇区ID"`
	DealId   uint64 `json:"dealId" gorm:"comment:订单ID"`
}

func (LotusDeals) TableName() string {
	return "lotus_deals"
}
