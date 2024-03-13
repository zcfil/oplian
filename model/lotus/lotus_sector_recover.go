package lotus

import (
	"oplian/global"
	"time"
)

type LotusSectorRecover struct {
	global.ZC_MODEL
	MinerId       string    `gorm:"comment:矿工号" json:"minerId"`
	SectorId      int       `gorm:"index;comment:扇区ID" json:"sectorId"`
	SectorSize    int       `gorm:"index;comment:扇区大小" json:"sectorSize"`
	Ticket        string    `gorm:"comment:授权码" json:"ticket"`
	SectorType    int       `gorm:"index;comment:扇区类型 1 cc,2 dc" json:"sectorType"`
	BelongingNode string    `gorm:"comment:所属节点ID" json:"belongingNode"`
	AbnormalTime  time.Time `gorm:"default:null;comment:异常时间" json:"abnormalTime"`
	SectorAddress string    `gorm:"comment:扇区地址" json:"sectorAddress"`
	RecoverTime   time.Time `gorm:"default:null;comment:恢复时间" json:"recoverTime"`
	SectorStatus  int       `gorm:"index;default:3;comment:扇区状态 1已恢复,2恢复失败,3异常,4恢复中" json:"sectorStatus"`
}

func (LotusSectorRecover) TableName() string {
	return "lotus_sector_recover"
}
