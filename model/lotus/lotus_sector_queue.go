package lotus

import (
	"oplian/global"
	"time"
)

type LotusSectorQueue struct {
	global.ZC_MODEL
	TaskName         string    `json:"taskName" gorm:"index;comment:任务名称"`
	FinishAt         time.Time `json:"finishAt" gorm:"default:null;comment:完成时间"`
	SectorSize       int       `json:"sectorSize"  gorm:"comment:扇区大小"`
	SectorType       int       `json:"sectorType"  gorm:"comment:CC：1，DC：2"`
	Actor            string    `json:"actor" gorm:"comment:节点号"`
	JobTotal         int       `json:"jobTotal" gorm:"comment:扇区总数"`
	RunCount         int       `json:"runCount"  gorm:"comment:正在跑数量"`
	CompleteCount    int       `json:"completeCount"  gorm:"comment:完成数量"`
	TaskStatus       int       `json:"taskStatus"  gorm:"default:1;comment:任务状态：1.进行中，2.暂停中,3.已完成,4.已终止,5.已过期,6,订单解析中,7.解析失败"`
	ConcurrentImport int       `json:"concurrentImport"  gorm:"default:1;comment:并发导入数量"`
}

func (LotusSectorQueue) TableName() string {
	return "lotus_sector_queue"
}
