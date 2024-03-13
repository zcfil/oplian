package system

import "oplian/global"

type SysDictionaryConfig struct {
	global.ZC_MODEL
	RelationId      string `gorm:"relation_id;comment:关联id" json:"relation_id"`
	DictionaryKey   string `gorm:"dictionary_key;comment:字典Key" json:"dictionary_key"`
	DefaultValue    string `gorm:"default_value;comment:默认字典值" json:"default_value"`
	DictionaryValue string `gorm:"dictionary_value;comment:字典值" json:"dictionary_value"`
	DictionaryType  int    `gorm:"dictionary_type;comment:字典类型" json:"dictionary_type"`
	Remark          string `gorm:"remark;comment:备注" json:"remark"`
}

func (SysDictionaryConfig) TableName() string {
	return "sys_dictionary_config"
}
