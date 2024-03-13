package system

import (
	"oplian/global"
	"time"
)

type SysWarnManage struct {
	global.ZC_MODEL
	WarnId           string    `gorm:"comment:告警ID" json:"warn_id"`
	WarnName         string    `gorm:"comment:告警名称" json:"warn_name"`
	WarnType         int       `gorm:"index;comment:告警类型 1 指标告警 2日志告警 3 关联告警 4事件告警(策略告警) 5巡检告警 6业务告警" json:"warn_type"`
	WarnStatus       int       `gorm:"index;default:2;comment:告警状态 1已恢复,2告警中,3已关闭" json:"warn_status"`
	ComputerId       string    `gorm:"index;comment:电脑ID" json:"computer_id"`
	Ip               string    `gorm:"index;comment:IP" json:"ip"`
	AssetsNum        string    `gorm:"comment:资产编号" json:"assets_num"`
	Sn               string    `gorm:"comment:SN" json:"sn"`
	WarnInfo         string    `gorm:"type:text;comment:告警消息" json:"warn_info"`
	NotifyPerson     string    `gorm:"comment:通知人" json:"notify_person"`
	StrategiesId     string    `gorm:"strategies_id;comment:策略ID" json:"strategies_id"`
	ComputerType     int       `gorm:"comment:主机类型" json:"computer_type"`
	ComputerRoomId   string    `gorm:"comment:机房ID" json:"computer_room_id"`
	ComputerRoomName string    `gorm:"comment:机房名称" json:"computer_room_name"`
	ComputerRoomNo   string    `gorm:"comment:机房编号" json:"computer_room_no"`
	WarnTime         time.Time `gorm:"default:null;comment:告警时间" json:"warn_time"`
	NotifyTime       time.Time `gorm:"default:null;comment:通知时间" json:"notify_time"`
	FinishTime       time.Time `gorm:"default:null;comment:完成时间" json:"finish_time"`
	Remark           string    `gorm:"remark;comment:备注" json:"remark"`
}

func (SysWarnManage) TableName() string {
	return "sys_warn_manage"
}
