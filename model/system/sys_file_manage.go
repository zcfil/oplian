package system

import "oplian/global"

type SysFileManage struct {
	global.ZC_MODEL
	FileName       string `gorm:"file_name;comment:文件名称" json:"file_name"`
	FileSize       int    `gorm:"file_size;comment:文件大小" json:"file_size"`
	Percentage     int    `gorm:"percentage;comment:文件进度" json:"percentage"`
	FileType       int    `gorm:"file_type;comment:文件类型 1一般文件,2证明文件,3高度文件,4快照文件,5miner文件" json:"file_type"`
	FileUrl        string `gorm:"file_url;comment:文件URL" json:"file_url"`
	OpId           string `gorm:"op_id;comment:主机OpId" json:"op_id"`
	GateWayId      string `gorm:"gate_way_id;comment:gateWayId" json:"gate_way_id"`
	ServerName     string `gorm:"server_name;comment:主机名称" json:"server_name"`
	Ip             string `gorm:"ip;comment:IP" json:"ip"`
	RoomName       string `gorm:"room_name;comment:机房名称" json:"room_name"`
	RoomId         string `gorm:"room_id;comment:机房编号" json:"room_id"`
	FileStatus     int    `gorm:"file_status;comment:文件状态 1获取中,2完成,3异常" json:"file_status"`
	ComputerSystem string `gorm:"computer_system;comment:操作系统" json:"computer_system"`
	Version        string `gorm:"version;comment:版本号" json:"version"`
	Remark         string `gorm:"remark;comment:备注" json:"remark"`
}

func (SysFileManage) TableName() string {
	return "sys_file_manage"
}
