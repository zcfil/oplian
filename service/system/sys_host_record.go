package system

import (
	"fmt"
	"oplian/config"
	"oplian/global"
	"oplian/model/common/request"
	request1 "oplian/model/lotus/request"
	response1 "oplian/model/lotus/response"
	"oplian/model/system"
	systemReq "oplian/model/system/request"
	"oplian/model/system/response"
	"oplian/utils"
	"strconv"
	"time"
)

//@function: CreateSysHostRecord
//@description: Create host information data
//@param: sysHostRecord model.SysHostRecord
//@return: err error

type HostRecordService struct{}

func (hostRecordService *HostRecordService) CreateSysHostRecord(sysHostRecord system.SysHostRecord) (err error) {
	err = global.ZC_DB.Create(&sysHostRecord).Error
	return err
}

//@function: DeleteSysHostRecord
//@description: Delete host information data
//@param: sysHostRecord model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) DeleteSysHostRecord(sysHostRecord system.SysHostRecord) (err error) {
	err = global.ZC_DB.Delete(&sysHostRecord).Error
	return err
}


//@function: DeleteSysHostRecordByIds
//@description: Batch delete record
//@param: ids request.IdsReq
//@return: err error

func (hostRecordService *HostRecordService) DeleteSysHostRecordByUUIDs(ids request.UUIDsReq) (err error) {
	err = global.ZC_DB.Delete(&[]system.SysHostRecord{}, "uuid in (?)", ids.UUIDs).Error
	return err
}


//@function: UpdateSysHostRecordAuto
//@description: The system automatically updates the host information
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) UpdateSysHostRecordAuto(sysHostRecord *system.SysHostRecord) (err error) {
	var dict system.SysHostRecord
	sysHostRecordMap := map[string]interface{}{
		"IntranetIP":       sysHostRecord.IntranetIP,
		"InternetIP":       sysHostRecord.InternetIP,
		"DeviceSN":         sysHostRecord.DeviceSN,
		"HostManufacturer": sysHostRecord.HostManufacturer,
		"HostModel":        sysHostRecord.HostModel,
		"OperatingSystem":  sysHostRecord.OperatingSystem,
		"CPUCoreNum":       sysHostRecord.CPUCoreNum,
		"CPUModel":         sysHostRecord.CPUModel,
		"MemorySize":       sysHostRecord.MemorySize,
		"DiskNum":          sysHostRecord.DiskNum,
		"DiskSize":         sysHostRecord.DiskSize,
		"ServerDNS":        sysHostRecord.ServerDNS,
		"SubnetMask":       sysHostRecord.SubnetMask,
		"Gateway":          sysHostRecord.Gateway,
		"GatewayId":        sysHostRecord.GatewayId,
		"GPUNum":           sysHostRecord.GPUNum,
		"SystemVersion":    sysHostRecord.SystemVersion,
		"SystemBits":       sysHostRecord.SystemBits,
	}
	db := global.ZC_DB.Where("uuid = ?", sysHostRecord.UUID).First(&dict)
	err = db.Updates(sysHostRecordMap).Error
	return err
}


//@function: UpdateSysHostRecordClassify
//@description: The system automatically updates the host status
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) UpdateSysHostRecordClassify(sysHostRecord *system.SysHostRecord) (err error) {
	var dict system.SysHostRecord
	sysHostRecordMap := map[string]interface{}{
		"HostClassify": sysHostRecord.HostClassify,
	}
	if sysHostRecord.IsGroupArray {
		sysHostRecordMap["IsGroupArray"] = true
	}
	db := global.ZC_DB.Where("uuid = ?", sysHostRecord.UUID).First(&dict)
	err = db.Updates(sysHostRecordMap).Error
	return err
}


//@function: UpdateSysHostRecordAuto
//@description: Update Host information Update time of host monitoring information
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) UpdateSysHostRecordMonitorTime(sysHostRecord *system.SysHostRecord) (err error) {
	var dict system.SysHostRecord
	sysHostRecordMap := map[string]interface{}{
		"MonitorTime": sysHostRecord.MonitorTime,
	}
	db := global.ZC_DB.Where("uuid = ?", sysHostRecord.UUID).First(&dict)
	err = db.Updates(sysHostRecordMap).Error
	return err
}


//@function: UpdateSysHostRecord
//@description: Update the host information
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) UpdateSysHostRecord(sysHostRecord *system.SysHostRecord) (err error) {
	db := global.ZC_DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	var dict system.SysHostRecord
	sysHostRecordMap := map[string]interface{}{
		"RoomId":       sysHostRecord.RoomId,
		"HostName":     sysHostRecord.HostName,
		"HostType":     sysHostRecord.HostType,
		"HostClassify": sysHostRecord.HostClassify,
		"HostGroupId":  sysHostRecord.HostGroupId,
		"AssetNumber":  sysHostRecord.AssetNumber,
	}
	dbHost := db.Where("uuid = ?", sysHostRecord.UUID).First(&dict)
	originRoomId := dict.RoomId
	if len(sysHostRecord.RoomId) > 0 {
		var machineRoom system.SysMachineRoomRecord
		roomDB := db.Where("room_id = ?", sysHostRecord.RoomId).First(&machineRoom)
		if len(machineRoom.GatewayId) != 0 && machineRoom.GatewayId != dict.GatewayId {
			err = global.GatewayIdMismatchError
			return
		}
		sysMachineRoomRecordMap := map[string]interface{}{
			"GatewayId": dict.GatewayId,
		}
		err = roomDB.Updates(sysMachineRoomRecordMap).Error
		if err != nil {
			return
		}

		sysHostRecordMap["RoomName"] = machineRoom.RoomName
	}

	err = dbHost.Updates(sysHostRecordMap).Error
	if err != nil {
		return
	}

	if len(originRoomId) != 0 && originRoomId != sysHostRecord.RoomId {
		var total int64
		err = db.Model(&system.SysHostRecord{}).Where("room_id = ?", originRoomId).Count(&total).Error
		if err != nil {
			return
		}

		if total == 0 {

			var machineRoomRecord system.SysMachineRoomRecord
			hostDB := db.Where("room_id = ?", originRoomId).First(&machineRoomRecord)
			sysMachineRoomMap := map[string]interface{}{
				"GatewayId": "",
			}
			err = hostDB.Updates(sysMachineRoomMap).Error
			if err != nil {
				return
			}
		}
	}
	return
}


//@function: GetSysHostRecord
//@description: Obtain single host information based on the uuid
//@param: uuid uuid.UUID
//@return: sysHostRecord system.SysHostRecord, err error

func (hostRecordService *HostRecordService) GetSysHostRecord(hostUUID string) (sysHostRecord system.SysHostRecord, err error) {
	err = global.ZC_DB.Where("uuid = ?", hostUUID).Find(&sysHostRecord).Error
	return
}


//@function: GetSysHostRecord
//@description: Obtain another host information under the current gateway
//@param: uuid uuid.UUID
//@return: sysHostRecord system.SysHostRecord, err error

func (hostRecordService *HostRecordService) GetSysOtherHostRecord(hostUUID, gatewayId string) (sysHostRecord system.SysHostRecord, err error) {
	err = global.ZC_DB.Where("uuid != ? and gateway_id = ?", hostUUID, gatewayId).First(&sysHostRecord).Error
	return
}

// GetSysHostRecordByIPAndGatewayId 获取当前网关下面的另外一条主机信息单条数据
func (hostRecordService *HostRecordService) GetSysHostRecordByIPAndGatewayId(intranetIp, gatewayId string) (sysHostRecord system.SysHostRecord, err error) {
	err = global.ZC_DB.Where("intranet_ip = ? and gateway_id = ?", intranetIp, gatewayId).Limit(1).Find(&sysHostRecord).Error
	return
}


//@function: GetSysHostRecord
//@description: Obtain the number of bound hosts based on the roomId
//@param: roomId uuid.UUID
//@return: sysHostRecord system.SysHostRecord, err error

func (hostRecordService *HostRecordService) GetSysHostRecordCountByRoomId(roomId string) (data []string, total int64, err error) {
	db := global.ZC_DB.Model(&system.SysHostRecord{})
	err = db.Where("room_id = ?", roomId).Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Select("uuid").Find(&data).Error
	return
}


//@function: GetSysHostRecord
//@description: Obtain the number of bound nodes based on roomId
//@param: roomId uuid.UUID
//@return: sysHostRecord system.SysHostRecord, err error

func (hostRecordService *HostRecordService) GetNodeHostNumByRoomId(roomId string) (total int64, err error) {
	db := global.ZC_DB.Model(&system.SysHostRecord{})
	err = db.Where("room_id = ? and host_classify in ?", roomId, []int64{config.HostMinerType, config.HostLotusType}).Count(&total).Error
	return
}


// @function: GetList
// @description: Page to get the host information list
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetSysHostRecordInfoList(info systemReq.SysHostRecordSearch) (list []response.SysHostRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostRecord{})
	var sysHostRecords []response.SysHostRecord

	if info.UUID != "" {
		db = db.Where("uuid = ?", info.UUID)
	}

	if info.RoomKeyword != "" {
		db = db.Where(`room_id != ""`).Where("room_id LIKE ? OR room_name LIKE ?", "%"+info.RoomKeyword+"%", "%"+info.RoomKeyword+"%")
	}

	if info.Keyword != "" {
		db = db.Where("host_name LIKE ? OR room_name LIKE ? OR asset_number LIKE ? OR device_sn LIKE ? OR internet_ip LIKE ? OR intranet_ip LIKE ?",
			"%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.HostKeyword != "" {
		db = db.Where("host_name LIKE ? OR intranet_ip LIKE ?", "%"+info.HostKeyword+"%", "%"+info.HostKeyword+"%")
	}
	if info.RoomId != "" {
		db = db.Where("room_id = ?", info.RoomId)
	}
	if info.DeviceSN != "" {
		db = db.Where("device_sn LIKE ?", "%"+info.DeviceSN+"%")
	}
	if info.HostType != nil {
		db = db.Where("host_type = ?", *info.HostType)
	}
	if info.HostClassify != nil {
		if *info.HostClassify == config.HostMinerType {
			db = db.Where("host_classify in ?", []int64{config.HostMinerType, config.HostLotusType})
		} else {
			db = db.Where("host_classify = ?", *info.HostClassify)
		}
	}
	if info.HostGroupId != nil {
		db = db.Where("host_group_id = ?", *info.HostGroupId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&sysHostRecords).Error
	return sysHostRecords, total, err
}


//@function: UpdateSysHostRecord
//@description: Update the host information
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) UpdateSysHostRecordRoomId(sysHostRecordReq *systemReq.BindSysHostRecordsReq) (err error) {
	tx := global.ZC_DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var machineRoomRecord system.SysMachineRoomRecord
	err = tx.Where("room_id = ?", sysHostRecordReq.RoomId).First(&machineRoomRecord).Error
	if err != nil {
		return
	}

	for _, val := range sysHostRecordReq.HostUUIDs {
		var hostRecord system.SysHostRecord
		hostDB := tx.Where("uuid = ?", val).First(&hostRecord)

		if len(machineRoomRecord.GatewayId) != 0 {
			if len(hostRecord.GatewayId) != 0 && machineRoomRecord.GatewayId != hostRecord.GatewayId {
				err = global.GatewayIdMismatchError
				return
			}
		} else {
			machineRoomRecord.GatewayId = hostRecord.GatewayId
			err = tx.Where("room_id = ?", machineRoomRecord.RoomId).Updates(machineRoomRecord).Error
			if err != nil {
				return
			}
		}
		sysHostRecordMap := map[string]interface{}{
			"RoomId":   sysHostRecordReq.RoomId,
			"RoomName": machineRoomRecord.RoomName,
		}
		err = hostDB.Updates(sysHostRecordMap).Error
		if err != nil {
			return
		}
	}
	return
}

// @author: Nathan
// @function: GetList
// @description: Paginate to get a brief list of hosts
// @param: info request.SysHost
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetSysHostRecordListNormal(page request.PageInfo) (list []response.SysHost, total, totalNum int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostRecord{})
	var sysHost []response.SysHost

	timeStamp := time.Now().Add(-time.Minute * 5).Unix()
	db = db.Where("monitor_time >= ?", timeStamp)

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&sysHost).Error
	if err != nil {
		return
	}
	err = global.ZC_DB.Model(&system.SysHostRecord{}).Count(&totalNum).Error
	if err != nil {
		return
	}
	return sysHost, total, totalNum, err
}

// @author: Nathan
// @function: GetList
// @description: Paginate to get a brief list of hosts
// @param: info request.SysHost
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetSysHostRecordList(page request.PageInfo) (list []response.SysHostInfo, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)

	db := global.ZC_DB.Model(&system.SysHostRecord{})
	var sysHost []response.SysHostInfo

	if page.Keyword != "" {
		db = db.Where("host_name LIKE ? or intranet_ip LIKE ? ", "%"+page.Keyword+"%", "%"+page.Keyword+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&sysHost).Error
	return sysHost, total, err
}

// @author: Nathan
// @function: GetList
// @description: Paginate to get a brief list of hosts
// @param: info request1.HostPage
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetSysHostTestRecordList(info request1.HostPage) (list []response.SysHostInfo, total int64, err error) {
	limit := ""
	sql := `SELECT h.gateway_id,h.room_id,h.room_name,h.uuid,h.host_name,h.intranet_ip,h.memory_size,h.disk_size,
                   h.host_classify,h.asset_number,h.device_sn,h.gpu_num
					FROM sys_host_test_records ht
					INNER JOIN sys_host_records h ON h.uuid = ht.host_uuid
					WHERE test_type = ? AND test_result = 2`
	if info.Keyword != "" {
		sql += fmt.Sprintf(` and (host_name LIKE '%%%s%%' or intranet_ip LIKE '%%%s%%') `, info.Keyword, info.Keyword)
	}
	if info.Classify == config.HostMinerTest {
		limit = utils.LimitAndOrder("IFNULL(h.host_classify,-1)", "DESC", info.Page, info.PageSize)
	} else {
		limit = utils.LimitAndOrder("IFNULL(h.host_classify,-1)", "ASC", info.Page, info.PageSize)
	}

	if err = global.ZC_DB.Raw("SELECT COUNT(1) FROM sys_host_test_records ht INNER JOIN sys_host_records h ON h.uuid = ht.host_uuid WHERE test_type = ? AND test_result = 2", info.Classify).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var sysHost []response.SysHostInfo
	err = global.ZC_DB.Raw(sql+limit, info.Classify).Scan(&sysHost).Error
	return sysHost, total, err
}

func (hostRecordService *HostRecordService) GetSysHostTestRecordList1(info request1.HostPage) (list []response.SysHostInfo, total int64, err error) {
	limit := ""
	sql := `SELECT h.id,h.gateway_id,h.room_id,h.room_name,h.uuid,h.host_name,h.intranet_ip,h.memory_size,h.disk_size,
                   h.host_classify,h.asset_number,h.device_sn,h.gpu_num
					FROM sys_host_records h
					WHERE (if(host_classify=6 or host_classify=3 or host_classify=2 or host_classify=1 or host_classify=5,host_classify,-1)=? or host_classify=0)`
	if info.Keyword != "" {
		sql += fmt.Sprintf(` and (host_name LIKE '%%%s%%' or intranet_ip LIKE '%%%s%%') `, info.Keyword, info.Keyword)
	}
	if info.GateId != "" {
		sql += fmt.Sprintf(" and gateway_id='%s'", info.GateId)
	}
	if info.Classify == config.HostMinerTest {
		limit = utils.LimitAndOrder("IFNULL(h.host_classify,-1)", "DESC", info.Page, info.PageSize)
	} else {
		limit = utils.LimitAndOrder("IFNULL(h.host_classify,-1)", "ASC", info.Page, info.PageSize)
	}

	if err = global.ZC_DB.Raw("SELECT COUNT(1) FROM sys_host_test_records ht INNER JOIN sys_host_records h ON h.uuid = ht.host_uuid WHERE test_type = ? AND test_result = 2", info.Classify).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var sysHost []response.SysHostInfo

	err = global.ZC_DB.Raw(sql+limit, info.Classify).Scan(&sysHost).Error
	return sysHost, total, err
}

//@function: GetSysHostRecordCountByHostGroupId
//@description: Obtain the number of hosts to be bound based on the groupId
//@param: roomId uuid.UUID
//@return: sysHostRecord system.SysHostRecord, err error

func (hostRecordService *HostRecordService) GetSysHostRecordCountByHostGroupId(groupId int) (total int64, err error) {
	db := global.ZC_DB.Model(&system.SysHostRecord{})
	err = db.Where("host_group_id = ?", groupId).Count(&total).Error
	return
}

// @author: Nathan
// @function: GetList
// @description: Obtain the host binding information list
// @param: info request.SysHost
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetSysHostRecordBindList(req systemReq.BindSysHostRecordListReq) (list []response.SysHostBindList, err error) {

	db := global.ZC_DB.Model(&system.SysHostRecord{})
	if req.Bind {
		db = db.Where("room_id = ?", req.RoomId)
	} else {
		db = db.Where(`room_id = ""`)
	}
	if req.Keyword != "" {
		db = db.Where("host_name LIKE ? or asset_number LIKE ? or internet_ip LIKE ? or device_sn LIKE ? ",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	err = db.Order("id desc").Find(&list).Error
	return
}


//@function: UpdateSysHostRecordRoomIdUnbind
//@description: Update the host information
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) UpdateSysHostRecordRoomIdUnbind(sysHostRecordReq *systemReq.UnbindSysHostRecordsReq) (err error) {
	tx := global.ZC_DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	for _, val := range sysHostRecordReq.HostUUIDs {
		var hostRecord system.SysHostRecord
		hostDB := tx.Where("uuid = ?", val).First(&hostRecord)
		sysHostRecordMap := map[string]interface{}{
			"RoomId":   "",
			"RoomName": "",
		}
		err = hostDB.Updates(sysHostRecordMap).Error
		if err != nil {
			return
		}
	}

	var total int64
	err = tx.Model(&system.SysHostRecord{}).Where("room_id = ?", sysHostRecordReq.RoomId).Count(&total).Error
	if err != nil {
		return
	}

	if total == 0 {

		var machineRoomRecord system.SysMachineRoomRecord
		hostDB := tx.Where("room_id = ?", sysHostRecordReq.RoomId).First(&machineRoomRecord)
		sysMachineRoomMap := map[string]interface{}{
			"GatewayId": "",
		}
		err = hostDB.Updates(sysMachineRoomMap).Error
		if err != nil {
			return
		}
	}
	return
}


// @function: GetOpInfoList
// @description: The host information list can be queried under the following conditions: Host category, host name, and host IP address combined fuzzy query
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetOpInfoList(info systemReq.SysHostRecordListByClassifyReq) (list []response.OpInfo, err error) {

	var opInfo []response.OpInfo
	var resInfo []response.OpInfo
	var sysHostRecord []system.SysHostRecord
	var sysMachineRoomRecord []system.SysMachineRoomRecord

	db := global.ZC_DB.Model(&system.SysHostRecord{})
	sql := " 1=1 "
	var param []interface{}
	if info.GateWayId != "" {
		sql += "and gateway_id = ? "
		param = append(param, info.GateWayId)
	}
	if info.HostClassify != 0 {
		if info.HostClassify == config.HostMinerType {
			sql += "and host_classify in ? "
			param = append(param, []int64{config.HostMinerType, config.HostLotusType})
		} else {
			sql += "and host_classify = ? "
			param = append(param, info.HostClassify)
		}
	}
	if info.KeyWord != "" {
		info.KeyWord = "%" + info.KeyWord + "%"
		sql += "and (host_name LIKE ? OR internet_ip LIKE ?) "
		param = append(param, info.KeyWord, info.KeyWord)
	}
	err = db.Where(sql, param...).Find(&sysHostRecord).Error
	if err != nil {
		return nil, err
	}

	roomIdStr := ""
	for _, v := range sysHostRecord {
		if v.RoomId != "" {
			roomIdStr += `'` + v.RoomId + `',`
		}

		t := &response.OpInfo{
			HostName:        v.HostName,
			IntranetIP:      v.IntranetIP,
			InternetIP:      v.InternetIP,
			UUID:            v.UUID,
			DeviceSN:        v.DeviceSN,
			HostModel:       v.HostModel,
			OperatingSystem: v.OperatingSystem,
			CPUCoreNum:      v.CPUCoreNum,
			CPUModel:        v.CPUModel,
			MemorySize:      v.MemorySize,
			DiskNum:         v.DiskNum,
			DiskSize:        v.DiskSize,
			HostShelfLife:   v.HostShelfLife,
			HostType:        v.HostType,
			HostClassify:    v.HostClassify,
			ServerDNS:       v.ServerDNS,
			SubnetMask:      v.SubnetMask,
			Gateway:         v.Gateway,
			HostGroupId:     v.HostGroupId,
			RoomId:          v.RoomId,
			GatewayId:       v.GatewayId,
			GPUNum:          v.GPUNum,
			AssetNumber:     v.AssetNumber,
			SystemVersion:   v.SystemVersion,
			SystemBits:      v.SystemBits,
		}
		opInfo = append(opInfo, *t)
	}

	sql = ""
	param = nil
	db = global.ZC_DB.Model(&system.SysMachineRoomRecord{})
	if roomIdStr != "" {
		roomIdStr = utils.SubStr(roomIdStr, 0, len(roomIdStr)-1)
		sql += " room_id in(" + roomIdStr + ")"
	}

	err = db.Where(sql, param...).Find(&sysMachineRoomRecord).Error
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, v := range sysMachineRoomRecord {
		m[v.RoomId] = v.RoomName
	}

	for _, v := range opInfo {
		if v1, ok := m[v.RoomId]; ok {
			v.RoomName = v1
		}
		resInfo = append(resInfo, v)
	}

	return resInfo, err
}


//@function: UpdateSysHostRecord
//@description: Perform the binding operation on the host table
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) HostBindRoomByGatewayId(sysRoom *system.SysMachineRoomRecord) (err error) {

	sysHostRecordMap := map[string]interface{}{
		"RoomId":   sysRoom.RoomId,
		"RoomName": sysRoom.RoomName,
	}
	err = global.ZC_DB.Model(&system.SysHostRecord{}).Where("gateway_id = ?", sysRoom.GatewayId).
		Updates(sysHostRecordMap).Error
	if err != nil {
		return
	}
	return
}

// @author: Lex
// @function: GetList
// @description: Gets a list of hosts of the type that you have obtained
// @param: info request.SysHost
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetSysHostRecordListForPatrol(classify []int64, gatewayId string) (list []response.SysHostRecordPatrol, err error) {

	db := global.ZC_DB.Model(&system.SysHostRecord{}).Where("host_classify <> ? and host_classify in ? and gateway_id = ?", config.HostDCStorageType, classify, gatewayId)
	err = db.Order("id desc").Find(&list).Error
	return
}

// GetSysHostRecordListForReplace Get host information except DC machines
func (hostRecordService *HostRecordService) GetSysHostRecordListForReplace(gatewayId string) (list []response.SysHostRecordPatrol, err error) {

	db := global.ZC_DB.Model(&system.SysHostRecord{}).Where("host_classify <> ? and gateway_id = ?", config.HostDCStorageType, gatewayId)
	err = db.Order("id desc").Find(&list).Error
	return
}


// @function: GetNetHostList
// @description: Obtain the list of selected network hosts. The search criteria can be joint fuzzy query of host name and host IP address
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetNetHostList(info systemReq.GetNetHostListReq, gatewayId string) (list []response.SysHostBindList, err error) {
	// 创建db
	db := global.ZC_DB.Model(&system.SysHostRecord{}).Where(`uuid != ?`, info.UUID).
		Where("net_occupy_time < ?", time.Now().Unix()).Where("gateway_id = ?", gatewayId)
	if info.Keyword != "" {
		db = db.Where("host_name LIKE ? or intranet_ip LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	err = db.Order("updated_at desc").Find(&list).Error
	return
}


//@function: UpdateHostNetOccupyTime
//@description: Updates when this host is used to test the network of the rest of the hosts
//@param: sysHostRecord *model.SysHostRecord
//@return: err error

func (hostRecordService *HostRecordService) UpdateHostNetOccupyTime(sysHostRecord *system.SysHostRecord) (err error) {
	var dict system.SysHostRecord
	sysHostTestRecordMap := map[string]interface{}{
		"NetOccupyTime": sysHostRecord.NetOccupyTime,
	}
	db := global.ZC_DB.Where("uuid = ?", sysHostRecord.UUID).First(&dict)
	err = db.Updates(sysHostTestRecordMap).Error
	return err
}

// @author: Nathan
// @function: GetList
// @description: Obtain the number of free hosts
// @param: info request.SysHost
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetFreeHostNum(isAll bool) (total int64, err error) {
	// 创建db
	db := global.ZC_DB.Model(&system.SysHostRecord{})
	if !isAll {
		db = db.Where("host_classify = 0")
	}
	err = db.Count(&total).Error
	return total, err
}


// @function: GetNetHostList
// @description: Obtain the list of selected inspection hosts. The search criteria can be combined fuzzy search by host type, host name, and host IP address
// @param: info request.SysHostInfoSearch
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetPatrolHostList(info systemReq.GetPatrolHostListReq) (list []response.SysHostBindList, err error) {

	db := global.ZC_DB.Model(&system.SysHostRecord{}).Where("host_type <> ?", config.HostDCStorageType)
	if info.PatrolType != 0 {
		if info.PatrolType == config.HostWorkerType {
			db = db.Where("host_classify in ?", []int64{config.HostWorkerType, config.HostC2WorkerType})
		} else if info.PatrolType == config.HostMinerType {
			db = db.Where("host_classify in ?", []int64{config.HostMinerType, config.HostLotusType})
		} else {
			db = db.Where("host_classify = ?", info.PatrolType)
		}
	}
	if info.Keyword != "" {
		db = db.Where("host_name LIKE ? or intranet_ip LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	err = db.Order("updated_at desc").Find(&list).Error
	return
}

// @author: Nathan
// @function: GetHostNum
// @description: Get host count
// @param: info request.SysHost
// @return: list interface{}, total int64, err error

func (hostRecordService *HostRecordService) GetHostNum() (total int64, err error) {
	err = global.ZC_DB.Model(&system.SysHostRecord{}).Count(&total).Error
	return total, err
}

func (hostRecordService *HostRecordService) PushHostTestRecord(testType int, opId string) error {

	var total int64
	if err := global.ZC_DB.Raw("SELECT COUNT(1) FROM sys_host_test_records WHERE test_type = ? AND test_result = 2 AND host_uuid = ?", testType, opId).Count(&total).Error; err != nil {
		return err
	}
	if total == 0 {
		sql := `INSERT INTO sys_host_test_records (test_type,test_result,host_uuid)VALUES(?,2,?)`
		return global.ZC_DB.Exec(sql, testType, opId).Error
	}
	return nil
}

// GetDCStorageList  Obtain the monitoring DC check-in list
func (hostRecordService *HostRecordService) GetDCStorageList(info systemReq.HostMonitorReq) (list []response1.DCStorageMonitorInfo, total int64, err error) {
	sqlparam := " WHERE sh.host_classify = " + strconv.Itoa(config.HostDCStorageType)
	if info.GateId != "" {
		sqlparam = fmt.Sprintf(` and sh.gateway_id = '%s' `, info.GateId)
	}
	if info.Keyword != "" {
		sqlparam += ` and ( sh.intranet_ip like '%` + info.Keyword + `%')`
	}

	sqltotal := `SELECT count(1) FROM sys_host_records sh` + sqlparam
	if err = global.ZC_DB.Raw(sqltotal).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var storageList []response1.DCStorageMonitorInfo
	sql := `SELECT sh.uuid,sh.gateway_id,sh.intranet_ip FROM sys_host_records sh`

	sql += sqlparam + utils.LimitAndOrder("id", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&storageList).Error; err != nil {
		return nil, 0, err
	}

	return storageList, total, err
}
