package request

import (
	"time"
)

type SysExecuteRecords struct {
	TaskName    string    `json:"task_name"`
	TaskNum     string    `json:"task_num"`
	TaskType    int       `json:"task_type"`
	OperateUser string    `json:"operate_user"`
	Status      int       `json:"status"`
	BeginTime   time.Time `json:"begin_time"`
	TimeLength  int       `json:"time_length"`
	Remark      string    `json:"remark"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}
