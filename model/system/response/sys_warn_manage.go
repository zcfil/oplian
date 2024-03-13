package response

import (
	"oplian/model/system"
	"time"
)

type WarnManageRes struct {
	ID               int64     `json:"id"`
	WarnId           string    `json:"warn_id"`
	WarnName         string    `json:"warn_name"`
	WarnType         int64     `json:"warn_type"`
	WarnStatus       int64     `json:"warn_status"`
	ComputerId       int64     `json:"computer_id"`
	Ip               string    `json:"ip"`
	Sn               string    `json:"sn"`
	WarnInfo         string    `json:"warn_info"`
	NotifyPerson     string    `json:"notify_person"`
	PolicyId         string    `json:"policy_id"`
	ComputerType     int64     `json:"computer_type"`
	ComputerRoomId   string    `json:"computer_room_id"`
	ComputerRoomName string    `json:"computer_room_name"`
	ComputerRoomNo   string    `json:"computer_room_no"`
	WarnTime         time.Time `json:"warn_time"`
	NotifyTime       time.Time `json:"notify_time"`
	FinishTime       time.Time `json:"finish_time"`
	Remark           string    `json:"remark"`
	DeletedAt        time.Time `json:"deleted_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type WarnStrategiesRes struct {
	WarnStrategies system.SysWarnStrategies
	IpInfo         []IpInfo
}

type WarnAllRes struct {
	PolicyTotal    int `json:"policy_total"`    //策略告警数量
	PatrolTotal    int `json:"patrol_total"`    //巡检告警数量
	BusinessTotal  int `json:"business_total"`  //业务告警数量
	ProcessedTotal int `json:"processed_total"` //已处理告警数量
	WarnTotal      int `json:"warn_total"`      //告警总数量
}

type WarnTrendRes struct {
	TimeStr string `json:"time_str"`
	Number  int    `json:"number"`
}

type IpInfo struct {
	Id         int    `json:"id"`
	Ip         string `json:"ip"`
	OpId       string `json:"op_id"`
	RoomId     string `json:"room_id"`
	GateWayId  string `json:"gate_way_id"`
	ServerName string `json:"server_name"`
	RelationId string `json:"relation_id"`
}

type SysWarnStrategies struct {
	Id               int    `json:"id"`
	StrategiesId     string `json:"strategies_id"`     // 策略ID
	StrategiesName   string `json:"strategies_name"`   // 策略名称
	StrategiesType   int    `json:"strategies_type"`   // 策略类型 1指标 2日志 3 关联 4事件
	StrategiesStatus int    `json:"strategies_status"` // 策略状态 1启用,0停用
	CycleType        int    `json:"cycle_type"`        // 周期 1周期,0实时
	CycleLength      int    `json:"cycle_length"`      // 时长,单位秒
	CycleUnit        int64  `json:"cycle_unit"`        // 单位 1秒 2分
	WarnLevel        int    `json:"warn_level"`        // 级别 1致命,2预警,3提醒
	KeyWord          string `json:"key_word"`          // 关键字
	MonitorItem      string `json:"monitor_item"`      // 监控项
	TriggerNum       int    `json:"trigger_num"`       // 触发条件,周期
	TriggerCheck     int    `json:"trigger_check"`     // 触发条件,检测
	RecoverNum       int    `json:"recover_num"`       // 恢复条件,周期
	LostEnable       int    `json:"lost_enable"`       // 无数据,1启用,0停用
	LostNum          int    `json:"lost_num"`          // 无数据,周期
	EffectivePeriod  string `json:"effective_period"`  // 生效时间
	NotificationSet  string `json:"notification_set"`  // 通知设置 1告警触发时,2告警恢复时,3告警关闭时
	IntervalLength   int    `json:"interval_length"`   // 通知间隔,单位分钟
	NotifyPerson     string `json:"notify_person"`     // 通知人
	IpList           string `json:"ip_list"`           // 关联IP
	CreatedAt        string `json:"created_at"`        // 创建时间
}

type ItemPrj struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
