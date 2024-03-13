package lotus

import (
	"oplian/global"
	"time"
)

type LotusSectorPiece struct {
	global.ZC_MODEL
	SectorId       uint64    `json:"sectorId" gorm:"index;comment:扇区ID"`
	Actor          string    `json:"actor"  gorm:"comment:节点ID"`
	WorkerIp       string    `json:"workerIp" gorm:"comment:worker主机IP"`
	PieceCid       string    `json:"pieceCid" gorm:"comment:piece id"`
	PieceSize      uint64    `json:"pieceSize"  gorm:"comment:订单大小"`
	CarSize        int       `gorm:"comment:car文件大小" json:"carSize"`
	DataCid        string    `gorm:"comment:数据CID" json:"dataCid"`
	DealId         uint64    `json:"dealId" gorm:"comment:订单ID"`
	DealUuid       string    `json:"dealUuid" gorm:"comment:订单UUID"`
	QueueId        uint64    `json:"queueId" gorm:"index;comment:任务队列ID"`
	RunIndex       int       `json:"runCount"  gorm:"comment:任务队列序号"`
	CarPath        string    `json:"carPath" gorm:"comment:car文件路径（浮动非永久有效）"`
	CarOpId        string    `json:"carOpId" gorm:"comment:car文件OpID（浮动非永久有效）"`
	JobStatus      int       `json:"jobStatus"  gorm:"default:1;comment:任务状态：1.待创建，2.创建中,3.已完成,4.创建失败,5.匹配失败,6.已过期"`
	ExpirationTime time.Time `json:"expirationTime" gorm:"comment:到期时间"`
}

func (LotusSectorPiece) TableName() string {
	return "lotus_sector_piece"
}
