package system

import (
	"oplian/global"
)

type SysApi struct {
	global.ZC_MODEL
	Path        string `json:"path" gorm:"comment:api路径"`
	Description string `json:"description" gorm:"comment:api中文描述"`
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`
	Method      string `json:"method" gorm:"default:POST;comment:方法"`
}

func (SysApi) TableName() string {
	return "sys_apis"
}
