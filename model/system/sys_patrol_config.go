// 自动生成模板SysPatrolConfig
package system

import (
	"oplian/global"
)

// 如果含有time.Time 请自行import time包
type SysPatrolConfig struct {
	global.ZC_MODEL
	PatrolType      int64 `json:"patrolType" form:"patrolType" gorm:"column:patrol_type;comment:巡检类型, 1 节点机巡检, 2 存储机巡检, 3 worker机巡检"`
	IntervalHours   int64 `json:"intervalHours" form:"intervalHours" gorm:"column:interval_hours;comment:间隔时间(时)"`
	IntervalMinutes int64 `json:"intervalMinutes" form:"intervalMinutes" gorm:"column:interval_minutes;comment:间隔时间(分)"`
	IntervalSeconds int64 `json:"intervalSeconds" form:"intervalSeconds" gorm:"column:interval_seconds;comment:间隔时间(秒)"`
	IntervalTime    int64 `json:"intervalTime" form:"intervalTime" gorm:"column:interval_time;comment:间隔时间(数字类型,单位为秒)"`
}

func (SysPatrolConfig) TableName() string {
	return "sys_patrol_config"
}
