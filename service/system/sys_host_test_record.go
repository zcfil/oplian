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


//@function: CreateSysHostTestRecord
//@description: Create host test information data
//@param: sysHostTestRecord model.SysHostTestRecord
//@return: err error

type HostTestRecordService struct{}

func (hostTestRecordService *HostTestRecordService) CreateSysHostTestRecord(sysHostTestRecord system.SysHostTestRecord) (err error) {
	err = global.ZC_DB.Create(&sysHostTestRecord).Error
	return err
}


//@function: DeleteSysHostTestRecordByIds
//@description: Batch delete record
//@param: ids request.IdsReq
//@return: err error

func (hostTestRecordService *HostTestRecordService) DeleteSysHostTestRecordByUUIDs(ids request.UUIDsReq) (err error) {
	err = global.ZC_DB.Delete(&[]system.SysHostTestRecord{}, "host_uuid in (?)", ids.UUIDs).Error
	return err
}


//@function: UpdateSysHostTestRecordAuto
//@description: The system automatically updates the host test data
//@param: sysHostTestRecord *model.SysHostTestRecord
//@return: err error

func (hostTestRecordService *HostTestRecordService) UpdateSysHostTestRecord(sysHostTestRecord *system.SysHostTestRecord) (err error) {
	var dict system.SysHostTestRecord
	sysHostTestRecordMap := map[string]interface{}{
		"TestBeginAt":      sysHostTestRecord.TestBeginAt,
		"IsAddPower":       sysHostTestRecord.IsAddPower,
		"SelectHostUUIDs":  sysHostTestRecord.SelectHostUUIDs,
		"SelectHostIPs":    sysHostTestRecord.SelectHostIPs,
		"TestEndAt":        sysHostTestRecord.TestEndAt,
		"TestResult":       sysHostTestRecord.TestResult,
		"CPUHardInfo":      sysHostTestRecord.CPUHardInfo,
		"CPUHardScore":     sysHostTestRecord.CPUHardScore,
		"GPUHardInfo":      sysHostTestRecord.GPUHardInfo,
		"GPUHardScore":     sysHostTestRecord.GPUHardScore,
		"MemoryHardInfo":   sysHostTestRecord.MemoryHardInfo,
		"MemoryHardScore":  sysHostTestRecord.MemoryHardScore,
		"DiskHardInfo":     sysHostTestRecord.DiskHardInfo,
		"DiskHardScore":    sysHostTestRecord.DiskHardScore,
		"NetTestInfo":      sysHostTestRecord.NetTestInfo,
		"NetTestScore":     sysHostTestRecord.NetTestScore,
		"GPUTestInfo":      sysHostTestRecord.GPUTestInfo,
		"GPUTestScore":     sysHostTestRecord.GPUTestScore,
		"DiskIO":           sysHostTestRecord.DiskIO,
		"DiskAllRate":      sysHostTestRecord.DiskAllRate,
		"DiskAllRateScore": sysHostTestRecord.DiskAllRateScore,
		"DiskSSDRate":      sysHostTestRecord.DiskSSDRate,
		"DiskSSDRateScore": sysHostTestRecord.DiskSSDRateScore,
	}
	db := global.ZC_DB.Where("host_uuid = ? AND test_type = ?", sysHostTestRecord.HostUUID, sysHostTestRecord.TestType).First(&dict)
	err = db.Updates(sysHostTestRecordMap).Error
	return err
}


//@function: UpdateSysHostTestRecordAuto
//@description: Update the host test result
//@param: sysHostTestRecord *model.SysHostTestRecord
//@return: err error

func (hostTestRecordService *HostTestRecordService) UpdateSysHostTestResult(sysHostTestRecord *system.SysHostTestRecord) (err error) {
	var dict system.SysHostTestRecord
	sysHostTestRecordMap := map[string]interface{}{
		"TestResult": sysHostTestRecord.TestResult,
		"TestEndAt":  sysHostTestRecord.TestEndAt,
	}
	db := global.ZC_DB.Where("host_uuid = ? AND id = ?", sysHostTestRecord.HostUUID, sysHostTestRecord.ID).First(&dict)
	err = db.Updates(sysHostTestRecordMap).Error
	return err
}


//@function: GetSysHostTestInfo
//@description: Obtain the single piece of host test information based on host_uuid
//@param: uuid uuid.UUID
//@return: sysHostTestRecord system.SysHostTestRecord, err error

func (hostTestRecordService *HostTestRecordService) GetSysHostTestInfo(id int, hostUUID string) (sysHostTestRecord system.SysHostTestRecord, err error) {
	err = global.ZC_DB.Where("id = ? and host_uuid = ?", id, hostUUID).First(&sysHostTestRecord).Error
	return
}


//@function: GetSysHostTestInfo
//@description: Obtain a single piece of host test information based on the test type test_type and host_uuid
//@param: uuid uuid.UUID
//@return: sysHostTestRecord system.SysHostTestRecord, err error

func (hostTestRecordService *HostTestRecordService) GetSysHostTestInfoByType(hostUUID string, testType int64) (sysHostTestRecord system.SysHostTestRecord, err error) {
	err = global.ZC_DB.Where("host_uuid = ? and test_type = ?", hostUUID, testType).First(&sysHostTestRecord).Error
	return
}


//@function: GetSysHostTestInfo
//@description: Obtain a single piece of host test information based on the test result test_result and host_uuid
//@param: uuid uuid.UUID
//@return: sysHostTestRecord system.SysHostTestRecord, err error

func (hostTestRecordService *HostTestRecordService) GetSysHostTestInfoByTestResult(hostUUID string, testResult int64) (sysHostTestRecord system.SysHostTestRecord, err error) {
	err = global.ZC_DB.Where("host_uuid = ? and test_result = ?", hostUUID, testResult).First(&sysHostTestRecord).Error
	return
}


//@function: GetSysHostTestInfo
//@description: Obtain the single piece of host test information based on host_uuid
//@param: uuid uuid.UUID
//@return: sysHostTestRecord system.SysHostTestRecord, err error

func (hostTestRecordService *HostTestRecordService) GetSysHostTestReport(info systemReq.GetHostTestReportReq) (sysHostTestRecord system.SysHostTestRecord, err error) {
	err = global.ZC_DB.Where("id = ? and host_uuid = ?", info.ID, info.HostUUID).First(&sysHostTestRecord).Error
	return
}


// @function: GetList
// @description: Paginate for a list of host test information
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostTestRecordService *HostTestRecordService) GetSysHostTestRecordInfoList(info systemReq.SysHostTestRecordSearch) (list []response.SysHostTestRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostTestRecord{})
	var sysHostTestRecords []response.SysHostTestRecord

	if info.TestType != 0 {
		db = db.Where("test_type = ?", info.TestType)
	}
	if info.TestResult != 0 {
		db = db.Where("test_result = ?", info.TestResult)
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
	err = db.Order("updated_at desc").Limit(limit).Offset(offset).Find(&sysHostTestRecords).Error
	return sysHostTestRecords, total, err
}


//@function: UpdateSysHostTestRecordAuto
//@description: Shut down the system to update host test data
//@param: sysHostTestRecord *model.SysHostTestRecord
//@return: err error

func (hostTestRecordService *HostTestRecordService) UpdateSysHostTestRecordClose(sysHostTestRecord *system.SysHostTestRecord) (err error) {
	var dict system.SysHostTestRecord
	sysHostTestRecordMap := map[string]interface{}{
		"TestResult": sysHostTestRecord.TestResult,
		"TestEndAt":  sysHostTestRecord.TestEndAt,
	}
	db := global.ZC_DB.Where("host_uuid = ? AND id = ?", sysHostTestRecord.HostUUID, sysHostTestRecord.ID).First(&dict)
	err = db.Updates(sysHostTestRecordMap).Error
	return err
}


//@function: UpdateSysHostTestRecordAuto
//@description: The system starts to update the host test data again
//@param: sysHostTestRecord *model.SysHostTestRecord
//@return: err error

func (hostTestRecordService *HostTestRecordService) RestartUpdateSysHostTestRecord(sysHostTestRecord *system.SysHostTestRecord) (err error) {
	sysHostTestRecordMap := map[string]interface{}{
		"TestBeginAt": sysHostTestRecord.TestBeginAt,
		"TestEndAt":   sysHostTestRecord.TestEndAt,
		"TestResult":  sysHostTestRecord.TestResult,
	}
	err = global.ZC_DB.Where("host_uuid = ? AND id = ?", sysHostTestRecord.HostUUID, sysHostTestRecord.ID).
		First(&system.SysHostTestRecord{}).Updates(sysHostTestRecordMap).Error
	return err
}


// @function: GetList
// @description: Gets a list of execution timeouts
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostTestRecordService *HostTestRecordService) GetTimeoutTestList() (list []system.SysHostTestRecord, err error) {

	db := global.ZC_DB.Model(&system.SysHostTestRecord{}).Where("test_result = ?", config.HostUnderTest)

	checkTime := 1*time.Hour + 30*time.Minute
	db = db.Where("test_begin_at < ?", time.Now().Add(-checkTime).Unix())

	err = db.Order("id desc").Find(&list).Error
	return
}
