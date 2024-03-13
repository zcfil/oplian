package system

import "oplian/global"

type SysColony struct {
	global.ZC_MODEL
	ColonyName string `json:"colonyName" form:"colonyName" gorm:"unique;column:colony_name;comment:集群名称"`
	ColonyType int    `json:"colonyType" form:"colonyType" gorm:"column:colony_type;comment:集群类型：1.NFS"`
}

func (SysColony) TableName() string {
	return "sys_colony"
}
