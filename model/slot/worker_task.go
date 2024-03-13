package slot

import "oplian/global"

type LotusWorkerTask struct {
	global.ZC_MODEL
	RelationId string `gorm:"relation_id;comment:关联id" json:"relation_id"`                 // 关联id
	TaskId     string `gorm:"task_id;comment:任务ID" json:"task_id"`                         // 任务ID
	Agent      string `gorm:"agent;comment:代理商" json:"agent"`                              // 代理商
	SectorSize int    `gorm:"sector_size;comment:扇区大小" json:"sector_size"`                 // 扇区大小
	ServerName string `gorm:"server_name;comment:主机名称" json:"server_name"`                 // 主机名称
	Ip         string `gorm:"ip;comment:IP" json:"ip"`                                     // IP
	TaskType   int    `gorm:"task_type;comment:任务类型 1ap,2p1,3p2,4c1,5c2" json:"task_type"` // 任务类型 1ap,2p1,3p2,4c1,5c2
	TaskStatus int    `gorm:"task_status;comment:任务状态 0进行中,1完成,2失败" json:"task_status"`    // 任务状态 0进行中,1完成,2失败
	BeginTime  string `gorm:"begin_time;comment:开始时间" json:"begin_time"`                   // 开始时间
	EndTime    string `gorm:"end_time;comment:开始时间" json:"end_time"`                       // 开始时间
	TimeLength string `gorm:"time_length;comment:耗时" json:"time_length"`                   // 耗时
	Remark     string `gorm:"remark;comment:备注" json:"remark"`                             // 备注
}

func (LotusWorkerTask) TableName() string {
	return "lotus_worker_task"
}
