package system

import (
	"oplian/global"
	"time"
)

type SysOpRelations struct {
	global.ZC_MODEL
	RelationId    string    `gorm:"index;comment:关联id" json:"relation_id"`
	RelationType  int       `gorm:"index;comment:关联类型 1 脚本执行,2 文件分发,3 告警策略管理,4 扇区恢复" json:"relation_type"`
	GateWayId     string    `gorm:"comment:机房gateWayId" json:"gate_way_id"`
	OpId          string    `gorm:"index;comment:主机OpId" json:"op_id"`
	ServerName    string    `gorm:"comment:主机名称" json:"server_name"`
	AssetsNum     string    `gorm:"comment:资产编号" json:"assets_num"`
	DeviceSn      string    `gorm:"comment:设备SN" json:"device_sn"`
	OperateSystem string    `gorm:"comment:操作系统" json:"operate_system"`
	Ip            string    `gorm:"comment:IP" json:"ip"`
	Port          string    `gorm:"comment:端口号" json:"port"`
	OpFilePath    string    `gorm:"index;comment:主机OpId" json:"op_file_path"`
	RoomId        string    `gorm:"index;comment:机房ID" json:"room_id"`
	RoomName      string    `gorm:"comment:机房名称" json:"room_name"`
	Status        int       `gorm:"default:3;comment:执行状态 1成功 2失败 3执行中" json:"status"`
	ResMsg        string    `gorm:"type:text;comment:返回消息" json:"res_msg"`
	BeginTime     time.Time `gorm:"default:CURRENT_TIMESTAMP(3);comment:开始时间" json:"begin_time"`
	TimeLength    string    `gorm:"comment:时长" json:"time_length"`
	Remark        string    `gorm:"remark;comment:备注" json:"remark"`
}

func (SysOpRelations) TableName() string {
	return "sys_op_relations"
}
