package system

import (
	"oplian/global"
	"oplian/model/system"
	"oplian/model/system/response"
)

//@function: CreateSysPatrolConfig
//@description: Create host patrol setup data
//@param: sysPatrolConfig model.SysPatrolConfig
//@return: err error

type PatrolConfigService struct{}

func (hostPatrolRecordService *PatrolConfigService) CreateSysPatrolConfig(sysPatrolConfig system.SysPatrolConfig) (err error) {
	err = global.ZC_DB.Create(&sysPatrolConfig).Error
	return err
}

//@function: UpdateSysPatrolConfig
//@description: The system automatically updates the host inspection Settings
//@param: sysPatrolConfig *model.SysPatrolConfig
//@return: err error

func (hostPatrolRecordService *PatrolConfigService) UpdateSysPatrolConfig(sysPatrolConfig *system.SysPatrolConfig) (err error) {
	var dict system.SysPatrolConfig
	sysPatrolConfigMap := map[string]interface{}{
		"IntervalHours":   sysPatrolConfig.IntervalHours,
		"IntervalMinutes": sysPatrolConfig.IntervalMinutes,
		"IntervalSeconds": sysPatrolConfig.IntervalSeconds,
		"IntervalTime":    sysPatrolConfig.IntervalTime,
	}
	db := global.ZC_DB.Where("patrol_type = ?", sysPatrolConfig.PatrolType).First(&dict)
	err = db.Updates(sysPatrolConfigMap).Error
	return err
}

// @function: GetSysPatrolConfigInfoList
// @description: Gets a list of host inspection Settings
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostPatrolRecordService *PatrolConfigService) GetSysPatrolConfigInfoList() (list []response.SysPatrolConfig, err error) {

	var sysPatrolConfigs []response.SysPatrolConfig
	err = global.ZC_DB.Model(&system.SysPatrolConfig{}).Order("patrol_type asc").Find(&sysPatrolConfigs).Error
	return sysPatrolConfigs, err
}

// @function: GetSysPatrolConfigInfoByType
// @description: Obtain the host inspection Settings list based on different types
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostPatrolRecordService *PatrolConfigService) GetSysPatrolConfigInfoByType(patrolType int) (info system.SysPatrolConfig, err error) {

	var sysPatrolConfig system.SysPatrolConfig
	err = global.ZC_DB.Model(&system.SysPatrolConfig{}).Where("patrol_type = ?", patrolType).First(&sysPatrolConfig).Error
	return sysPatrolConfig, err
}
