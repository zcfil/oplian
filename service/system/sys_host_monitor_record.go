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


//@function: CreateSysHostMonitorRecord
//@description: Create host information data
//@param: sysHostMonitorRecord model.SysHostMonitorRecord
//@return: err error

type HostMonitorRecordService struct{}

func (hostMonitorRecordService *HostMonitorRecordService) CreateSysHostMonitorRecord(sysHostMonitorRecord system.SysHostMonitorRecord) (err error) {
	err = global.ZC_DB.Create(&sysHostMonitorRecord).Error
	return err
}


//@function: GetSysHostMonitorRecord
//@description: Get a single piece of data within 5 minutes of host monitoring information according to uuid
//@param: uuid uuid.UUID
//@return: sysHostRecord system.SysHostMonitorRecord, err error

func (hostMonitorRecordService *HostMonitorRecordService) GetSysHostMonitorRecord(hostUUID, beginTime, endTime string) (sysHostRecord system.SysHostMonitorRecord, err error) {
	err = global.ZC_DB.Where("host_uuid = ?", hostUUID).Where("created_at>? and created_at<=?", beginTime, endTime).
		Order("id desc").First(&sysHostRecord).Error
	return
}


//@function: GetSysHostMonitorRecord
//@description: Get the records of the host within half an hour according to the uuid
//@param: uuid uuid.UUID
//@return: sysHostMonitorRecord system.SysHostMonitorRecord, err error

func (hostMonitorRecordService *HostMonitorRecordService) GetSysHostMonitorRecordList(hostUUID, gpuId, keyword, beginTime, endTime string) (sysHostMonitorRecord []response.HostMonitorRate, err error) {
	db := global.ZC_DB.Model(&system.SysHostMonitorRecord{}).Where("host_uuid = ? AND gpu_id = ?", hostUUID, gpuId).
		Where("created_at>? and created_at<=?", beginTime, endTime)
	if keyword != "gpu" {
		db = db.Where("gpu_id = 0")
	}
	err = db.Order("id desc").Find(&sysHostMonitorRecord).Error
	return
}


// @function: GetList
// @description: Page to get the host information list
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostMonitorRecordService *HostMonitorRecordService) GetSysHostMonitorRecordInfoList(info systemReq.SysHostMonitorRecordSearch, beginTime, endTime string) (list []response.SysHostMonitorRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostMonitorRecord{})
	var sysHostMonitorRecords []response.SysHostMonitorRecord
	db = db.Where("created_at>? and created_at<=?", beginTime, endTime)

	if info.HostUUID != "" {
		db = db.Where("host_uuid = ?", info.HostUUID).Order("id desc")
	}
	if info.Keyword == "cpu" {
		if info.UseRate > 0 {
			if !info.IsLessThan {
				db = db.Where("cpu_use_rate >= ?", info.UseRate).Order("cpu_use_rate desc")
			} else {
				db = db.Where("cpu_use_rate < ?", info.UseRate).Order("cpu_use_rate asc")
			}
		}
	}
	if info.Keyword == "disk" {
		if info.UseRate > 0 {
			if !info.IsLessThan {
				db = db.Where("disk_use_rate >= ?", info.UseRate).Order("disk_use_rate desc")
			} else {
				db = db.Where("disk_use_rate < ?", info.UseRate).Order("disk_use_rate asc")
			}
		}
	}
	if info.Keyword == "memory" {
		if info.UseRate > 0 {
			if !info.IsLessThan {
				db = db.Where("memory_use_rate >= ?", info.UseRate).Order("memory_use_rate desc")
			} else {
				db = db.Where("memory_use_rate < ?", info.UseRate).Order("memory_use_rate asc")
			}
		}
	}
	if info.Keyword == "gpu" {
		if info.UseRate > 0 {
			if !info.IsLessThan {
				db = db.Where("gpu_use_rate >= ?", info.UseRate).Order("gpu_use_rate desc")
			} else {
				db = db.Where("gpu_use_rate < ?", info.UseRate).Order("gpu_use_rate asc")
			}
		}
	}

	if info.Keyword != "gpu" {
		db = db.Where("gpu_id = 0")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&sysHostMonitorRecords).Error
	return sysHostMonitorRecords, total, err
}


//@function: GetLastHostMonitorRecords
//@description: Get the most recently recorded inspection host information
//@param: uuid uuid.UUID
//@return: sysHostRecord system.SysHostMonitorRecord, err error

func (hostMonitorRecordService *HostMonitorRecordService) GetLastHostMonitorRecords() (monitorRecords []system.SysHostMonitorRecord, err error) {
	var monitorRecord system.SysHostMonitorRecord

	err = global.ZC_DB.Where("gpu_id = 0").Order("id desc").First(&monitorRecord).Error
	if err != nil {
		return
	}

	beginTime := monitorRecord.CreatedAt.Add(-time.Minute * 2)
	endTime := monitorRecord.CreatedAt.Add(time.Minute * 2)

	err = global.ZC_DB.Where("gpu_id = 0").Model(&system.SysHostMonitorRecord{}).Where("created_at>? and created_at<=?", beginTime, endTime).
		Find(&monitorRecords).Error
	return
}

//@function: GetLastHostMonitorLists
//@description: Get the most recently recorded inspection host information, paging
//@param: uuid uuid.UUID
//@return: sysHostRecord system.SysHostMonitorRecord, err error

func (hostMonitorRecordService *HostMonitorRecordService) GetLastHostMonitorLists(info systemReq.GetHostRunListReq) (monitorRecords []system.SysHostMonitorRecord, total int64, err error) {
	var sysHostRecord system.SysHostMonitorRecord

	err = global.ZC_DB.Order("id desc").First(&sysHostRecord).Error
	if err != nil {
		return
	}

	beginTime := sysHostRecord.CreatedAt.Add(-time.Minute * 4)
	endTime := sysHostRecord.CreatedAt.Add(time.Minute * 4)

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostMonitorRecord{})
	db = db.Where("created_at>? and created_at<=?", beginTime, endTime).Where("gpu_id = 0")

	if info.RoomId != "" {
		db1 := global.ZC_DB.Model(&system.SysHostRecord{})
		var dataId []string
		err = db1.Where("room_id = ?", info.RoomId).Pluck("uuid", &dataId).Error
		if err != nil {
			return
		}
		db = db.Where("host_uuid in ?", dataId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&monitorRecords).Error
	return
}


//@function: GetLastHostMonitorLists
//@description: Get the most recently recorded inspection host information, paging
//@param: uuid uuid.UUID
//@return: sysHostRecord system.SysHostMonitorRecord, err error

func (hostMonitorRecordService *HostMonitorRecordService) GetLastHostStorageInfoMonitorLists(info request.PageInfo, beginTime, endTime string) (monitorRecords []system.SysHostMonitorRecord, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostMonitorRecord{})
	db = db.Where("created_at>? and created_at<=?", beginTime, endTime).Where("gpu_id = 0")

	dbHost := global.ZC_DB.Model(&system.SysHostRecord{})
	var dataId []string
	err = dbHost.Where("host_classify = ?", config.HostStorageType).Pluck("uuid", &dataId).Error
	if err != nil {
		return
	}
	db = db.Where("host_uuid in ?", dataId)

	if info.Keyword != "" {
		db1 := global.ZC_DB.Model(&system.SysHostRecord{})
		err = db1.Where("intranet_ip like ?", "%"+info.Keyword+"%").Pluck("uuid", &dataId).Error
		if err != nil {
			return
		}
		db = db.Where("host_uuid in ?", dataId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&monitorRecords).Error
	return
}

// DeleteByTime Delete the corresponding information based on the time
func (hostMonitorRecordService *HostMonitorRecordService) DeleteByTime() (err error) {
	sql := `DELETE FROM sys_host_monitor_records WHERE created_at < ` + time.Now().Format("2006-01-02 00:00:00")
	err = global.ZC_DB.Exec(sql).Error
	return err
}
