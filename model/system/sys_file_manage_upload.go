package system

import "oplian/global"

type SysFileManageUpload struct {
	global.ZC_MODEL
	RelationId    string `gorm:"index;comment:关联id" json:"relation_id"`
	FileType      int    `gorm:"comment:文件类型 1本地文件,2主机文件,3文件管理文件" json:"file_type"`
	FileName      string `gorm:"comment:文件名称" json:"file_name"`
	FileSize      int    `gorm:"default:0;comment:文件大小" json:"file_size"`
	OpId          string `gorm:"default:null;comment:opId" json:"op_id"`
	Ip            string `gorm:"default:null;comment:IP" json:"ip"`
	ServerName    string `gorm:"default:null;comment:主机名称" json:"server_name"`
	FilePath      string `gorm:"comment:主机文件路径" json:"file_path"`
	RoomId        string `gorm:"default:null;comment:机房ID" json:"room_id"`
	RoomName      string `gorm:"default:null;comment:机房名称" json:"room_name"`
	OperateSystem string `gorm:"default:null;comment:操作系统" json:"operate_system"`
	Remark        string `gorm:"remark;comment:备注" json:"remark"`
}

func (SysFileManageUpload) TableName() string {
	return "sys_file_manage_upload"
}
