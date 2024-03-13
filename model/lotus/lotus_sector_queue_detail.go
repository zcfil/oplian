package lotus

import (
	"oplian/global"
)

type LotusSectorQueueDetail struct {
	global.ZC_MODEL
	SectorId  uint64 `json:"sectorId"  gorm:"index;comment:扇区ID"`
	Actor     string `json:"actor"  gorm:"comment:节点ID"`
	QueueId   uint64 `json:"queueId" gorm:"index;comment:任务队列ID"`
	RunIndex  int    `json:"runCount"  gorm:"comment:任务队列序号"`
	JobStatus int    `json:"jobStatus"  gorm:"default:1;comment:任务状态：1.待创建，2.创建中,3.创建失败,4.已完成"`
}

func (LotusSectorQueueDetail) TableName() string {
	return "lotus_sector_queue_detail"
}
