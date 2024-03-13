package deploy

import (
	"fmt"
	"oplian/global"
	model "oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	systemReq "oplian/model/system/request"
	"oplian/utils"
)

// AddStorage
// @author: nathan
// @function: AddStorage
// @description: Add storage
// @Param     model.LotusStorageInfo
// @return:  error
func (deploy *DeployService) AddStorage(storage *model.LotusStorageInfo) error {
	if storage.ID == 0 {
		return global.ZC_DB.Save(storage).Error
	}
	return global.ZC_DB.Updates(storage).Error
}

// UpdateStorage
// @author: nathan
// @function: UpdateStorage
// @description: Update store
// @Param     model.LotusStorageInfo
// @return:  error
func (deploy *DeployService) UpdateStorage(storage *model.LotusStorageInfo) error {
	return global.ZC_DB.Updates(storage).Error
}

// GetStorage
// @author: nathan
// @function: GetStorage
// @description: Fetch store
// @Param     id uint64
// @return:  error
func (deploy *DeployService) GetStorage(id uint64) (model.LotusStorageInfo, error) {
	var storage model.LotusStorageInfo
	return storage, global.ZC_DB.Model(model.LotusStorageInfo{}).Where("id = ?", id).First(&storage).Error
}

// GetStorageList
// @author: nathan
// @function: GetStorageList
// @description: Get storage list
// @Param     request.WorkerInfoPage
// @return:  error
func (deploy *DeployService) GetStorageList(info request.WorkerInfoPage) (list interface{}, total int64, err error) {
	sqlparam := ""
	if info.GateId != "" {
		sqlparam = fmt.Sprintf(` and s.gate_id = '%s' `, info.GateId)
	}
	if info.Actor != "" {
		sqlparam += fmt.Sprintf(` and colony_name = '%s' `, info.Actor)
	}
	if info.Keyword != "" {
		sqlparam += ` and ( s.ip like '%` + info.Keyword + `%' or host_name like '%` + info.Keyword + `%' or colony_name like '%` + info.Keyword + `%')`
	}

	sqltotal := `SELECT count(1) FROM lotus_storage_info s LEFT JOIN sys_host_records r ON s.op_id = r.uuid WHERE 1=1 ` + sqlparam
	if err = global.ZC_DB.Raw(sqltotal).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var storageList []*response.StorageInfo
	sql := `SELECT s.id,s.op_id,s.gate_id,r.room_id,room_name,host_name,device_sn,s.ip,s.colony_type,
       		s.deploy_status,s.start_at,s.finish_at,miner_id,s.colony_name,s.err_msg
			FROM lotus_storage_info s
			LEFT JOIN sys_host_records r ON s.op_id = r.uuid
			WHERE 1=1 `

	sql += sqlparam + utils.LimitAndOrder("colony_name", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&storageList).Error; err != nil {
		return nil, 0, err
	}

	return storageList, total, err
}

func (deploy *DeployService) GetRelationStorageList(id request.IDActor) (list interface{}, err error) {
	var workerList []response.RelationWorker
	sqlparam := fmt.Sprintf(" WHERE colony_name = '%s'", id.Actor)
	if id.Actor == "" {
		sqlparam = fmt.Sprintf(" WHERE miner_id = %d", id.ID)
	}
	sql := `SELECT w.id,w.op_id,w.gate_id,r.room_id,room_name,host_name,device_sn,w.ip
			FROM lotus_storage_info w 
			LEFT JOIN sys_host_records r ON w.op_id = r.uuid
			` + sqlparam
	return workerList, global.ZC_DB.Raw(sql).Scan(&workerList).Error
}

func (deploy *DeployService) GetStorageByOpID(opID string) (model.LotusStorageInfo, error) {
	var storage model.LotusStorageInfo
	return storage, global.ZC_DB.Model(model.LotusStorageInfo{}).Where("op_id = ?", opID).First(&storage).Error
}

func (deploy *DeployService) GetStorageInfoByActor(colonyName string) ([]model.LotusStorageInfo, error) {
	var storage []model.LotusStorageInfo
	return storage, global.ZC_DB.Model(model.LotusStorageInfo{}).Where("colony_name = ?", colonyName).Find(&storage).Error
}

// GetStorageMonitorList 获取监控Storage列表
func (deploy *DeployService) GetStorageMonitorList(info systemReq.HostMonitorReq) (list []response.StorageMonitorInfo, total int64, err error) {
	sqlparam := " WHERE 1=1"
	if info.GateId != "" {
		sqlparam = fmt.Sprintf(` and ls.gate_id = '%s' `, info.GateId)
	}
	if info.Keyword != "" {
		sqlparam += ` and ( ls.ip like '%` + info.Keyword + `%')`
	}
	//求总数
	sqltotal := `SELECT count(1) FROM lotus_storage_info ls ` + sqlparam
	if err = global.ZC_DB.Raw(sqltotal).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var storageList []response.StorageMonitorInfo
	sql := `SELECT ls.op_id,ls.gate_id,ls.ip,ls.colony_name,ls.colony_type FROM lotus_storage_info ls`

	sql += sqlparam + utils.LimitAndOrder("colony_name", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&storageList).Error; err != nil {
		return nil, 0, err
	}

	return storageList, total, err
}

func (deploy *DeployService) GetStorageMountInfoByActor(colonyName string) ([]response.StorageMountErrorList, error) {
	var storage []response.StorageMountErrorList
	return storage, global.ZC_DB.Model(model.LotusStorageInfo{}).Where("colony_name = ?", colonyName).Find(&storage).Error
}
