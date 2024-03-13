package system

import (
	"oplian/global"
)

type SysWarnStrategies struct {
	global.ZC_MODEL
	StrategiesId     string `gorm:"strategies_id;comment:策略ID" json:"strategies_id"`
	StrategiesName   string `gorm:"strategies_name;comment:策略名称" json:"strategies_name"`
	StrategiesType   int    `gorm:"strategies_type;comment:策略类型 1指标 2日志 3 关联 4事件" json:"strategies_type"`
	StrategiesStatus int    `gorm:"strategies_status;default:1;comment:策略状态 1启用,0停用" json:"strategies_status"`
	CycleType        int    `gorm:"cycle_type;comment:周期 1周期,0实时" json:"cycle_type"`
	CycleLength      int    `gorm:"cycle_length;comment:时长,单位秒" json:"cycle_length"`
	CycleUnit        int64  `gorm:"cycle_length;comment:时长,单位秒" json:"cycle_unit"`
	WarnLevel        int    `gorm:"warn_level;comment:级别 1致命,2预警,3提醒" json:"warn_level"`
	KeyWord          string `gorm:"key_word;comment:关键字" json:"key_word"`
	MonitorItem      int    `gorm:"monitor_item;comment:监控项" json:"monitor_item"`
	TriggerNum       int    `gorm:"trigger_num;comment:触发条件,周期" json:"trigger_num"`
	TriggerCheck     int    `gorm:"trigger_check;comment:触发条件,检测" json:"trigger_check"`
	RecoverNum       int    `gorm:"recover_num;comment:恢复条件,周期" json:"recover_num"`
	LostEnable       int    `gorm:"lost_enable;comment:无数据,1启用,0停用" json:"lost_enable"`
	LostNum          int    `gorm:"lost_num;comment:无数据,周期" json:"lost_num"`
	EffectivePeriod  string `gorm:"effective_period;comment:生效时间" json:"effective_period"`
	NotificationSet  string `gorm:"notification_set;comment:通知设置 1告警触发时,2告警恢复时,3告警关闭时" json:"notification_set"`
	IntervalLength   int    `gorm:"interval_length;comment:通知间隔,单位分钟" json:"interval_length"`
	NotifyPerson     string `gorm:"notify_person;comment:通知人" json:"notify_person"`
	IpList           string `gorm:"ip_list;comment:通知人" json:"ip_list"`
	Remark           string `gorm:"remark;comment:备注" json:"remark"`
}

func (SysWarnStrategies) TableName() string {
	return "sys_warn_strategies"
}
