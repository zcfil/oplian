package system

import (
	"oplian/global"
	"time"
)

type SysJobExecuteRecords struct {
	global.ZC_MODEL
	TaskName    string    `gorm:"comment:任务名称" json:"task_name"`
	TaskNum     string    `gorm:"comment:任务编码" json:"task_num"`
	TaskType    int       `gorm:"index;comment:任务类型  1脚本执行,2文件分发" json:"task_type"`
	OperateUser string    `gorm:"comment:操作用户" json:"operate_user"`
	OpNumber    int       `gorm:"comment:主机数量" json:"op_number"`
	Status      int       `gorm:"index;default:3;comment:执行状态 1成功 2失败 3执行中" json:"status"`
	BeginTime   time.Time `gorm:"default:CURRENT_TIMESTAMP(3);comment:开始时间" json:"begin_time"`
	TimeLength  string    `gorm:"comment:时长" json:"time_length"`
	ScriptName  string    `gorm:"comment:脚本名称" json:"script_name"`
	Remark      string    `gorm:"comment:备注" json:"remark"`
}

func (SysJobExecuteRecords) TableName() string {
	return "sys_job_execute_records"
}
