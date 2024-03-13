// 自动生成模板SysHostInfo
package system

import (
	"oplian/global"
)

type SysHostGroup struct {
	global.ZC_MODEL
	GroupName string `json:"groupName" form:"groupName" gorm:"column:group_name;comment:分组名称"`
}

func (SysHostGroup) TableName() string {
	return "sys_host_groups"
}
