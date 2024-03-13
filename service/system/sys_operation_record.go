package system

import (
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
)

//@function: CreateSysOperationRecord
//@description: Create record
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

type OperationRecordService struct{}

func (operationRecordService *OperationRecordService) CreateSysOperationRecord(sysOperationRecord system.SysOperationRecord) (err error) {
	err = global.ZC_DB.Create(&sysOperationRecord).Error
	return err
}

//@function: DeleteSysOperationRecordByIds
//@description: Batch delete record
//@param: ids request.IdsReq
//@return: err error

func (operationRecordService *OperationRecordService) DeleteSysOperationRecordByIds(ids request.IdsReq) (err error) {
	err = global.ZC_DB.Delete(&[]system.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	return err
}

//@function: DeleteSysOperationRecord
//@description: Delete operation record
//@param: sysOperationRecord model.SysOperationRecord
//@return: err error

func (operationRecordService *OperationRecordService) DeleteSysOperationRecord(sysOperationRecord system.SysOperationRecord) (err error) {
	err = global.ZC_DB.Delete(&sysOperationRecord).Error
	return err
}

//@function: DeleteSysOperationRecord
//@description: Obtain a single operation record based on the id
//@param: id uint
//@return: sysOperationRecord system.SysOperationRecord, err error

func (operationRecordService *OperationRecordService) GetSysOperationRecord(id uint) (sysOperationRecord system.SysOperationRecord, err error) {
	err = global.ZC_DB.Where("id = ?", id).First(&sysOperationRecord).Error
	return
}

//@function: GetSysOperationRecordInfoList
//@description: Paginate for a list of operation records
//@param: info systemReq.SysOperationRecordSearch
//@return: list interface{}, total int64, err error

func (operationRecordService *OperationRecordService) GetSysOperationRecordInfoList(info systemReq.SysOperationRecordSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysOperationRecord{})
	var sysOperationRecords []system.SysOperationRecord

	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	return sysOperationRecords, total, err
}
