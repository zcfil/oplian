package lotus

import (
	"oplian/global"
	"time"
)

type LotusSectorTask struct {
	global.ZC_MODEL
	TaskName          string    `gorm:"comment:gateWayId" json:"taskName"`
	Actor             string    `gorm:"comment:节点" json:"actor"`
	SectorType        int       `gorm:"comment:扇区类型 1 cc,2 dc" json:"sectorType"`
	SectorTotal       int       `gorm:"comment:扇区数量" json:"sectorTotal"`
	JobTotal          int       `gorm:"comment:工作扇区数量" json:"jobTotal"`
	RunCount          int       `gorm:"comment:正在跑数量" json:"runCount"`
	FinishCount       int       `gorm:"comment:完成数量" json:"finishCount"`
	SectorSize        int       `gorm:"comment:扇区大小" json:"sectorSize"`
	OriginalValueOpId string    `gorm:"type:text;comment:原值数据主机" json:"originalValueOpId"`
	OriginalValueDir  string    `gorm:"comment:原值数据主机目录" json:"originalValueDir"`
	StorageOpId       string    `gorm:"comment:存储机" json:"storageOpId"`
	StorageOpIp       string    `gorm:"comment:存储机IP" json:"storageOpIp"`
	StorageOpName     string    `gorm:"comment:存储机名称" json:"storageOpName"`
	TaskStatus        int       `gorm:"default 0;comment:0.未开始 1.进行中 2.暂停中,3.已完成,4.已终止" json:"taskStatus"`
	BeginTime         time.Time `gorm:"default null;comment:开始时间" json:"beginTime"`
	EndTime           time.Time `gorm:"default null;comment:结束时间" json:"endTime"`
}

func (LotusSectorTask) TableName() string {
	return "lotus_sector_task"
}
