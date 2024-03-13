package system

import (
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	"oplian/model/system/response"
)

//@function: CreateSysMachineRoomRecords
//@description: Create room information data
//@param: sysMachineRoomRecords model.SysMachineRoomRecord
//@return: err error

type MachineRoomRecordService struct{}

func (machineRoomRecordService *MachineRoomRecordService) CreateSysMachineRoomRecord(sysMachineRoomRecords system.SysMachineRoomRecord) (err error) {
	err = global.ZC_DB.Create(&sysMachineRoomRecords).Error
	return err
}

//@function: DeleteSysMachineRoomRecords
//@description: Delete room information data
//@param: sysMachineRoomRecords model.SysMachineRoomRecord
//@return: err error

func (machineRoomRecordService *MachineRoomRecordService) DeleteSysMachineRoomRecord(sysMachineRoomRecord system.SysMachineRoomRecord) (err error) {
	err = global.ZC_DB.Delete(&sysMachineRoomRecord, "room_id = ?", sysMachineRoomRecord.RoomId).Error
	return err
}

//@function: UpdateSysMachineRoomRecords
//@description: Update the equipment room information
//@param: sysMachineRoomRecords *model.SysMachineRoomRecord
//@return: err error

func (machineRoomRecordService *MachineRoomRecordService) UpdateSysMachineRoomRecord(sysMachineRoomRecords *system.SysMachineRoomRecord) (err error) {
	tx := global.ZC_DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var dict system.SysMachineRoomRecord
	sysMachineRoomRecordsMap := map[string]interface{}{
		"RoomName": sysMachineRoomRecords.RoomName,
		"RoomType": sysMachineRoomRecords.RoomType,
		//"CabinetsNum":     sysMachineRoomRecords.CabinetsNum,
		"RoomTemp":        sysMachineRoomRecords.RoomTemp,
		"RoomLeader":      sysMachineRoomRecords.RoomLeader,
		"RoomLeaderPhone": sysMachineRoomRecords.RoomLeaderPhone,
		"RoomSupplier":    sysMachineRoomRecords.RoomSupplier,
		"SupplierContact": sysMachineRoomRecords.SupplierContact,
		"SupplierPhone":   sysMachineRoomRecords.SupplierPhone,
		"RoomAdmin":       sysMachineRoomRecords.RoomAdmin,
		"RoomOwner":       sysMachineRoomRecords.RoomOwner,
		"PhysicalAddress": sysMachineRoomRecords.PhysicalAddress,
		"RoomArea":        sysMachineRoomRecords.RoomArea,
	}
	err = tx.Where("room_id = ?", sysMachineRoomRecords.RoomId).First(&dict).Updates(sysMachineRoomRecordsMap).Error

	sysHostRecordMap := map[string]interface{}{
		"RoomName": sysMachineRoomRecords.RoomName,
	}

	err = tx.Model(&system.SysHostRecord{}).Where("room_id = ?", sysMachineRoomRecords.RoomId).Updates(sysHostRecordMap).Error
	if err != nil {
		return
	}
	return
}

//@function: GetSysMachineRoomRecords
//@description: Obtain equipment room information based on the roomId
//@param: uuid uuid.UUID
//@return: sysMachineRoomRecords system.SysMachineRoomRecord, err error

func (machineRoomRecordService *MachineRoomRecordService) GetSysMachineRoomRecord(roomId string) (sysMachineRoomRecords system.SysMachineRoomRecord, err error) {
	err = global.ZC_DB.Where("room_id = ?", roomId).First(&sysMachineRoomRecords).Error
	return
}

//@function: GetSysMachineRoomRecords
//@description: Obtain equipment room information based on the roomId
//@param: uuid uuid.UUID
//@return: sysMachineRoomRecords system.SysMachineRoomRecord, err error

func (machineRoomRecordService *MachineRoomRecordService) GetRoomByGatewayId(gatewayId string) (sysMachineRoomRecords system.SysMachineRoomRecord, err error) {
	err = global.ZC_DB.Where("gateway_id = ?", gatewayId).First(&sysMachineRoomRecords).Error
	return
}

//@function: DeleteSysMachineRoomRecordByRoomIds
//@description: Batch delete record
//@param: ids request.IdsReq
//@return: err error

func (machineRoomRecordService *MachineRoomRecordService) DeleteSysMachineRoomRecordByIds(roomIds request.RoomIdsReq) (err error) {
	err = global.ZC_DB.Delete(&[]system.SysMachineRoomRecord{}, "room_id in (?)", roomIds.RoomIds).Error
	return err
}

// @function: GetSysMachineRoomRecordsInfoList
// @description: Obtain the equipment room information list by paging
// @param: info request.SysMachineRoomRecordsSearch
// @return: list interface{}, total int64, err error

func (machineRoomRecordService *MachineRoomRecordService) GetSysMachineRoomRecordInfoList(info systemReq.SysMachineRoomRecordSearch) (list []response.SysMachineRoomRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysMachineRoomRecord{})
	var sysMachineRoomRecords []response.SysMachineRoomRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Keyword != "" {
		db = db.Where("room_id LIKE ? or room_name LIKE ? ", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&sysMachineRoomRecords).Error
	return sysMachineRoomRecords, total, err
}

// @function: GetSysMachineRoomRecordsInfoList
// @description: Get a list of room information
// @param: info request.SysMachineRoomRecordsSearch
// @return: list interface{}, total int64, err error

func (machineRoomRecordService *MachineRoomRecordService) GetRoomRecordList() (list []response.RoomRecordList, err error) {
	db := global.ZC_DB.Model(&system.SysMachineRoomRecord{})
	var sysMachineRoomRecords []response.RoomRecordList
	err = db.Order("updated_at desc").Find(&sysMachineRoomRecords).Error
	return sysMachineRoomRecords, err
}

// @function: GetSysMachineRoomNum
// @description: Get the number of computer rooms
// @param: info request.SysMachineRoomRecordsSearch
// @return: list interface{}, total int64, err error

func (machineRoomRecordService *MachineRoomRecordService) GetSysMachineRoomNum() (total int64, err error) {
	db := global.ZC_DB.Model(&system.SysMachineRoomRecord{})
	err = db.Count(&total).Error
	return total, err
}
