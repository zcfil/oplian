package lotus

import (
	"oplian/global"
	"time"
)

type LotusWorkerRelations struct {
	global.ZC_MODEL
	RelationId   string    `gorm:"index;comment:关联id" json:"relationId"`
	RelationType int       `gorm:"comment:关联类型 1worker集群" json:"relationType"`
	GateId       string    `gorm:"comment:网关Id" json:"gateId"`
	OpId         string    `gorm:"index;comment:主机OpId" json:"opId"`
	TaskId       string    `gorm:"comment:任务ID" json:"taskId"`
	TaskAgent    string    `gorm:"comment:任务代理商" json:"taskAgent"`
	TaskAgentNo  string    `gorm:"comment:任务代理商编号" json:"taskAgentNo"`
	Ip           string    `gorm:"comment:IP" json:"ip"`
	ServerName   string    `gorm:"comment:主机名称" json:"serverName"`
	RoomId       string    `gorm:"index;comment:机房ID" json:"roomId"`
	RoomName     string    `gorm:"comment:机房名称" json:"roomName"`
	Miner        string    `gorm:"comment:节点号" json:"miner"`
	Number       int       `gorm:"comment:扇区号" json:"number"`
	SectorSize   int       `gorm:"comment:扇区大小" json:"sectorSize"`
	IsRemove     int       `gorm:"comment:扇区文件是否移除 1是,0否" json:"isRemove"`
	SyncStatus   int       `gorm:"comment:文件同步状态 1完成 ,0未完成" json:"syncStatus"`
	TaskStatus   int       `gorm:"index;default:4;comment:任务状态 1已完成 2失败 3进行中,4排队中" json:"taskStatus"`
	ResMsg       string    `gorm:"comment:返回消息" json:"resMsg"`
	BeginTime    time.Time `gorm:"default:null;comment:开始时间" json:"beginTime"`
	TimeLength   string    `gorm:"comment:时长" json:"timeLength"`
	EndTime      time.Time `gorm:"default:null;comment:开始时间" json:"endTime"`
}

func (LotusWorkerRelations) TableName() string {
	return "lotus_worker_relations"
}
