package lotus

import (
	"oplian/global"
)

type LotusSectorAbnormal struct {
	global.ZC_MODEL
	MinerId      string `gorm:"comment:矿工号" json:"minerId"`
	SectorId     int    `gorm:"index;comment:扇区ID" json:"sectorId"`
	AbnormalTime string `gorm:"default:null;comment:异常时间" json:"abnormalTime"`
	Count        int    `gorm:"default:0;comment:异常次数" json:"count"`
}

func (LotusSectorAbnormal) TableName() string {
	return "lotus_sector_abnormal"
}
