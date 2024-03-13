package system

import (
	"oplian/config"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	"oplian/model/system/response"
	"time"
)


//@function: CreateSysHostPatrolRecord
//@description: Create host inspection information data
//@param: sysHostPatrolRecord model.SysHostPatrolRecord
//@return: err error

type HostPatrolRecordService struct{}

func (hostPatrolRecordService *HostPatrolRecordService) CreateSysHostPatrolRecord(sysHostPatrolRecord system.SysHostPatrolRecord) (err error) {
	err = global.ZC_DB.Create(&sysHostPatrolRecord).Error
	return err
}


//@function: DeleteSysHostPatrolRecordByIds
//@description: Batch delete record
//@param: ids request.IdsReq
//@return: err error

func (hostPatrolRecordService *HostPatrolRecordService) DeleteSysHostPatrolRecordByUUIDs(ids request.UUIDsReq) (err error) {
	err = global.ZC_DB.Delete(&[]system.SysHostPatrolRecord{}, "host_uuid in (?)", ids.UUIDs).Error
	return err
}


//@function: UpdateSysHostPatrolRecordAuto
//@description: The system cannot connect to the host. Update the inspection status
//@param: sysHostPatrolRecord *model.SysHostPatrolRecord
//@return: err error

func (hostPatrolRecordService *HostPatrolRecordService) UpdateSysHostPatrolRecordConnectError(sysHostPatrolRecord *system.SysHostPatrolRecord) (err error) {
	var dict system.SysHostPatrolRecord
	sysHostPatrolRecordMap := map[string]interface{}{
		"PatrolEndAt":  sysHostPatrolRecord.PatrolEndAt,
		"PatrolResult": sysHostPatrolRecord.PatrolResult,
	}
	db := global.ZC_DB.Where("host_uuid = ? AND patrol_begin_at = ?", sysHostPatrolRecord.HostUUID, sysHostPatrolRecord.PatrolBeginAt).First(&dict)
	err = db.Updates(sysHostPatrolRecordMap).Error
	return err
}


//@function: UpdateSysHostPatrolRecordAuto
//@description: The system automatically updates the host inspection data
//@param: sysHostPatrolRecord *model.SysHostPatrolRecord
//@return: err error

func (hostPatrolRecordService *HostPatrolRecordService) UpdateSysHostPatrolRecord(sysHostPatrolRecord *system.SysHostPatrolRecord) (err error) {
	var dict system.SysHostPatrolRecord
	sysHostPatrolRecordMap := map[string]interface{}{
		"PatrolEndAt":                 sysHostPatrolRecord.PatrolEndAt,
		"PatrolResult":                sysHostPatrolRecord.PatrolResult,
		"DiskIO":                      sysHostPatrolRecord.DiskIO,
		"DiskIODuration":              sysHostPatrolRecord.DiskIODuration,
		"HostIsDown":                  sysHostPatrolRecord.HostIsDown,
		"HostIsDownDuration":          sysHostPatrolRecord.HostIsDownDuration,
		"HostNetStatus":               sysHostPatrolRecord.HostNetStatus,
		"HostNetDuration":             sysHostPatrolRecord.HostNetDuration,
		"LogInfoStatus":               sysHostPatrolRecord.LogInfoStatus,
		"LogInfoDuration":             sysHostPatrolRecord.LogInfoDuration,
		"LogOvertimeStatus":           sysHostPatrolRecord.LogOvertimeStatus,
		"LogOvertimeDuration":         sysHostPatrolRecord.LogOvertimeDuration,
		"WalletBalanceStatus":         sysHostPatrolRecord.WalletBalanceStatus,
		"WalletBalance":               sysHostPatrolRecord.WalletBalance,
		"WalletBalanceDuration":       sysHostPatrolRecord.WalletBalanceDuration,
		"LotusSyncStatus":             sysHostPatrolRecord.LotusSyncStatus,
		"LotusSyncDuration":           sysHostPatrolRecord.LotusSyncDuration,
		"GPUDriveStatus":              sysHostPatrolRecord.GPUDriveStatus,
		"GPUDriveDuration":            sysHostPatrolRecord.GPUDriveDuration,
		"PackageVersionStatus":        sysHostPatrolRecord.PackageVersionStatus,
		"PackageVersion":              sysHostPatrolRecord.PackageVersion,
		"PackageVersionDuration":      sysHostPatrolRecord.PackageVersionDuration,
		"DataCatalogStatus":           sysHostPatrolRecord.DataCatalogStatus,
		"DataCatalogDuration":         sysHostPatrolRecord.DataCatalogDuration,
		"EnvironmentVariableStatus":   sysHostPatrolRecord.EnvironmentVariableStatus,
		"EnvironmentVariableDuration": sysHostPatrolRecord.EnvironmentVariableDuration,
		"BlockLogStatus":              sysHostPatrolRecord.BlockLogStatus,
		"BlockLogDuration":            sysHostPatrolRecord.BlockLogDuration,
		"TimeSyncStatus":              sysHostPatrolRecord.TimeSyncStatus,
		"TimeSyncDuration":            sysHostPatrolRecord.TimeSyncDuration,
		"PingNetStatus":               sysHostPatrolRecord.PingNetStatus,
		"PingNetDuration":             sysHostPatrolRecord.PingNetDuration,
	}
	db := global.ZC_DB.Where("host_uuid = ? AND patrol_begin_at = ?", sysHostPatrolRecord.HostUUID, sysHostPatrolRecord.PatrolBeginAt).First(&dict)
	err = db.Updates(sysHostPatrolRecordMap).Error
	return err
}


//@function: GetSysHostPatrolInfo
//@description: Obtain the host inspection information based on host_uuid
//@param: uuid uuid.UUID
//@return: sysHostPatrolRecord system.SysHostPatrolRecord, err error

func (hostPatrolRecordService *HostPatrolRecordService) GetSysHostPatrolInfo(info systemReq.GetHostPatrolReportReq) (sysHostPatrolRecord system.SysHostPatrolRecord, err error) {
	err = global.ZC_DB.Where("id = ? and host_uuid = ?", info.ID, info.HostUUID).First(&sysHostPatrolRecord).Error
	return
}


//@function: GetSysHostPatrolInfo
//@description: Obtain host inspection information based on the test result patrol_result and host_uuid
//@param: uuid uuid.UUID
//@return: sysHostPatrolRecord system.SysHostPatrolRecord, err error

func (hostPatrolRecordService *HostPatrolRecordService) GetSysHostPatrolByResult(hostUUID string, patrolResult int64) (sysHostPatrolRecord system.SysHostPatrolRecord, err error) {
	err = global.ZC_DB.Where("host_uuid = ? and patrol_result = ?", hostUUID, patrolResult).First(&sysHostPatrolRecord).Error
	return
}


//@function: GetSysHostPatrolInfo
//@description: Obtain the host inspection information based on host_uuid
//@param: uuid uuid.UUID
//@return: sysHostPatrolRecord system.SysHostPatrolRecord, err error

func (hostPatrolRecordService *HostPatrolRecordService) GetSysHostPatrolReport(info systemReq.GetHostPatrolReportReq) (sysHostPatrolRecord response.SysHostPatrolInfo, err error) {
	err = global.ZC_DB.Model(&system.SysHostPatrolRecord{}).Where("id = ? and host_uuid = ?", info.ID, info.HostUUID).First(&sysHostPatrolRecord).Error
	return
}


// @function: GetList
// @description: Obtain the host inspection information list on a page
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostPatrolRecordService *HostPatrolRecordService) GetSysHostPatrolRecordInfoList(info systemReq.SysHostPatrolRecordSearch) (list []response.SysHostPatrolRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostPatrolRecord{})
	var sysHostPatrolRecords []response.SysHostPatrolRecord

	if info.PatrolType != 0 {
		db = db.Where("patrol_type = ?", info.PatrolType)
	}
	if info.PatrolResult != 0 {
		db = db.Where("patrol_result = ?", info.PatrolResult)
	}
	if info.RoomId != "" {
		db1 := global.ZC_DB.Model(&system.SysHostRecord{})
		var dataId []string
		err = db1.Where("room_id = ?", info.RoomId).Pluck("uuid", &dataId).Error
		if err != nil {
			return
		}
		db = db.Where("host_uuid in ?", dataId)
	}
	if info.HostNameIPKeyword != "" {
		db1 := global.ZC_DB.Model(&system.SysHostRecord{})
		var dataId []string
		err = db1.Where("host_name like ? or intranet_ip like ?", "%"+info.HostNameIPKeyword+"%", "%"+info.HostNameIPKeyword+"%").Pluck("uuid", &dataId).Error
		if err != nil {
			return
		}
		db = db.Where("host_uuid in ?", dataId)
	}
	if info.HostAssetSNKeyword != "" {
		db1 := global.ZC_DB.Model(&system.SysHostRecord{})
		var dataId []string
		err = db1.Where("asset_number like ? or device_sn like ?", "%"+info.HostAssetSNKeyword+"%", "%"+info.HostAssetSNKeyword+"%").Pluck("uuid", &dataId).Error
		if err != nil {
			return
		}
		db = db.Where("host_uuid in ?", dataId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&sysHostPatrolRecords).Error
	return sysHostPatrolRecords, total, err
}


// @function: GetList
// @description: Obtain the inspection list that has been running for a long time
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostPatrolRecordService *HostPatrolRecordService) GetTimeoutPatrolList() (list []system.SysHostPatrolRecord, err error) {
	var sysHostPatrolRecords []system.SysHostPatrolRecord

	db := global.ZC_DB.Model(&system.SysHostPatrolRecord{}).Where("patrol_result = ?", config.HostUnderTest)

	checkTime := config.HostNetTimeout
	db = db.Where("patrol_begin_at < ?", time.Now().Add(-checkTime).Unix())

	err = db.Order("id desc").Find(&sysHostPatrolRecords).Error
	return sysHostPatrolRecords, err
}
