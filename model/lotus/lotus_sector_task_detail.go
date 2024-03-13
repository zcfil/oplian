package lotus

import "oplian/global"

type LotusSectorTaskDetail struct {
	global.ZC_MODEL
	RelationId      string `gorm:"comment:关联id" json:"relationId"`
	SectorRecoverId int    `gorm:"comment:扇区恢复ID" json:"sectorRecoverId"`
	TaskOrder       int    `gorm:"comment:任务顺序" json:"taskOrder"`
	MinerId         string `gorm:"comment:矿工号" json:"minerId"`
	SectorId        int    `gorm:"comment:扇区ID" json:"sectorId"`
	Ticket          string `gorm:"comment:授权码" json:"ticket"`
	SectorAddress   string `gorm:"comment:扇区地址" json:"sectorAddress"`
	CarFile         string `gorm:"comment:dc扇区car文件" json:"carFile"`
	SectorStatus    int    `gorm:"default:0;comment:扇区状态 0等待,1开始,2ap,3p1,4p2,5完成,6恢复失败" json:"sectorStatus"`
	SectorSize      int    `gorm:"comment:扇区大小 32GiB,64GiB" json:"sectorSize"`
	WaitTime        string `gorm:"comment:等待时间" json:"waitTime"`
}

func (LotusSectorTaskDetail) TableName() string {
	return "lotus_sector_task_detail"
}
