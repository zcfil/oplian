package dispatch

import (
	"fmt"
	"oplian/global"
	"oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/utils"
	"time"
)

var DispatchServiceApi = new(DispatchService)

type DispatchService struct{}

//GetWorkerConfigList
//@author: nathan
//@function: GetWrokerConfigList
//@description: Get a list of worker task configurations
//@param: param request.RoomPageInfo
//@return: []lotus.LoutsWorkerConfig, error

func (w *DispatchService) GetWorkerConfigList(param request.RoomPageInfo) (list interface{}, total int64, err error) {
	var workerList []*response.WorkerConfig
	sqlparam := " WHERE 1=1 "
	sql := `SELECT l.id,l.op_id,l.gate_id,r.asset_number,r.room_name,r.host_name,device_sn,l.ip,r.cpu_core_num,
                   r.cpu_model,r.disk_num,r.disk_size,r.memory_size,pre_count1,pre_count2,on_off1,m.actor,w.deploy_status
			FROM lotus_worker_config l
			LEFT JOIN lotus_worker_info w ON l.worker_id = w.id
			LEFT JOIN lotus_miner_info m ON w.miner_id = m.id
			LEFT JOIN sys_host_records r ON l.op_id = r.uuid`
	if param.Keyword != "" {
		sqlparam += ` and (l.ip like '%` + param.Keyword + `%' or host_name like '%` + param.Keyword + `%')`
	}
	if param.GateId != "" {
		sqlparam += ` and l.gate_id = '` + param.GateId + `'`
	}
	if param.OnOff != 0 {
		if param.OnOff != 1 {
			param.OnOff = 0
		}
		sqlparam += fmt.Sprintf(" and l.on_off1 = %d", param.OnOff)
	}
	if param.Actor != "" {
		sqlparam += fmt.Sprintf(" and m.actor = '%s'", param.Actor)
	}
	sqlTotal := `SELECT COUNT(1) FROM lotus_worker_config l
			LEFT JOIN lotus_worker_info w ON l.worker_id = w.id
			LEFT JOIN lotus_miner_info m ON w.miner_id = m.id
			LEFT JOIN sys_host_records r ON l.op_id = r.uuid`
	//求总数
	err = global.ZC_DB.Model(&lotus.LoutsWorkerConfig{}).Raw(sqlTotal + sqlparam).Count(&total).Error
	//获取数据
	sql += sqlparam
	sql += utils.LimitAndOrder("l.created_at", "desc", param.Page, param.PageSize)

	return workerList, total, global.ZC_DB.Raw(sql).Scan(&workerList).Error
}

// GetWorkerConfig
// @author: nathan
// @function: GetWorkerConfig
// @description: Get a single worker task configuration
// @param: OpId string
// @return: lotus.LoutsWorkerConfig, error
func (w *DispatchService) GetWorkerConfig(OpId string) (lotus.LoutsWorkerConfig, error) {
	var worker lotus.LoutsWorkerConfig
	if err := global.ZC_DB.Model(lotus.LoutsWorkerConfig{}).Where("op_id = ?", OpId).Scan(&worker).Error; err != nil {
		return worker, err
	}
	if !worker.OnOff1 {
		worker.PreCount1 = 0
	}
	return worker, nil
}

//SetWrokerConfig
//@author: nathan
//@function: SetWrokerConfig
//@description: Set the worker task configuration
//@param: pre request.WorkerPre
//@return: error

func (w *DispatchService) SetConfig(pre request.PreConfig) error {
	var worker lotus.LoutsWorkerConfig
	db := global.ZC_DB.Model(lotus.LoutsWorkerConfig{})
	db.Where("id = ?", pre.ID).First(&worker)
	worker.UpdatedAt = time.Now()
	if pre.PreCount1 >= 0 {
		worker.PreCount1 = int(pre.PreCount1)
	}
	if pre.PreCount2 >= 0 {
		worker.PreCount2 = int(pre.PreCount2)
	}
	return db.Save(&worker).Error
}

//SetWrokerConfig
//@author: nathan
//@function: SetWrokerConfig
//@description: Set the worker task configuration
//@param: pre request.WorkerPre
//@return: error

func (w *DispatchService) SetWrokerConfig(pre request.PreConfig) error {
	var worker lotus.LoutsWorkerConfig
	db := global.ZC_DB.Model(lotus.LoutsWorkerConfig{})
	db.Where("worker_id = ?", pre.ID).First(&worker)
	worker.UpdatedAt = time.Now()
	if pre.PreCount1 >= 0 {
		worker.PreCount1 = int(pre.PreCount1)
	}
	if pre.PreCount2 >= 0 {
		worker.PreCount2 = int(pre.PreCount2)
	}
	return db.Save(&worker).Error
}

//AddWrokerConfig
//@author: nathan
//@function: AddWrokerConfig
//@description: Add worker task configuration
//@param: worker lotus.LoutsWorkerConfig
//@return: error

func (w *DispatchService) AddWrokerConfig(worker *lotus.LoutsWorkerConfig) error {
	//log.Println("worker:", worker)
	var config lotus.LoutsWorkerConfig
	if err := global.ZC_DB.First(&config, "op_id = ?", worker.OpId).Error; err == nil && config.ID != 0 {
		worker.PreCount1 = config.PreCount1
		worker.PreCount2 = config.PreCount2
	}
	if global.ZC_DB.Model(lotus.LoutsWorkerConfig{}).Where("op_id = ?", worker.OpId).Updates(worker).RowsAffected == 0 {
		return global.ZC_DB.Save(worker).Error
	}
	return nil
}

//OnOff
//@author: nathan
//@function: OnOff
//@description: Start-stop task
//@param: id uint, onOff boo
//@return: error

func (w *DispatchService) OnOff(id uint64, onOff bool) error {
	return global.ZC_DB.Model(lotus.LoutsWorkerConfig{}).Where("id = ?", id).Update("on_off1", onOff).Error
}

//OnOffByWorkerID
//@author: nathan
//@function: OnOffByWorkerID
//@description: Start or stop a task using the worker ID
//@param: id uint, onOff boo
//@return: error

func (w *DispatchService) OnOffByWorkerID(workerId uint64, onOff bool) error {
	return global.ZC_DB.Model(lotus.LoutsWorkerConfig{}).Where("worker_id = ?", workerId).Update("on_off1", onOff).Error
}
