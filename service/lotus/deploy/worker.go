package deploy

import (
	"context"
	"fmt"
	"oplian/define"
	"oplian/global"
	model "oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	systemReq "oplian/model/system/request"
	"oplian/service/pb"
	"oplian/utils"
	"time"
)

//@author: nathan
//@function: UpdateMinerStatusAndLink
//@description: Update Worker running status Connect miner
//@param: ID uint, status int
//@return: error

func (deploy *DeployService) UpdateWorkerStatusAndLink(ID uint, runStatus, deployStatus int, linkId uint64) error {
	var worker model.LotusWorkerInfo
	db := global.ZC_DB.Model(model.LotusWorkerInfo{})
	if err := db.Where("id = ?", ID).First(&worker).Error; err != nil {
		return err
	}
	if worker.RunStatus > -1 {
		worker.RunStatus = runStatus
	}
	if linkId != 0 {
		worker.MinerId = linkId
	}
	if deployStatus > 0 {
		worker.DeployStatus = deployStatus
	}
	return db.Save(&worker).Error
}

//@author: nathan
//@function: GetWorkerList
//@description: Paging for data
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (deploy *DeployService) GetWorkerList(info request.WorkerInfoPage) (list interface{}, total int64, err error) {
	sqlparam := ""
	if info.GateId != "" {
		sqlparam = fmt.Sprintf(` and w.gate_id = '%s' `, info.GateId)
	}
	//if info.Actor != "" {
	//	sqlparam += fmt.Sprintf(" AND miner.actor like '%%%s%%'", info.Keyword, info.Actor)
	//}
	if info.Keyword != "" {
		sqlparam += ` and (w.ip like '%` + info.Keyword + `%' or host_name like '%` + info.Keyword + `%' or miner.actor like '%` + info.Keyword + `%') `
	}
	if info.DeployStatus != 0 {
		sqlparam += fmt.Sprintf(` and w.deploy_status = %d `, info.DeployStatus)
	}

	sqltotal := `SELECT count(1) FROM lotus_worker_info w LEFT JOIN lotus_miner_info miner ON w.miner_id = miner.id LEFT JOIN sys_host_records r ON w.op_id = r.uuid WHERE 1=1 ` + sqlparam
	if err = global.ZC_DB.Raw(sqltotal).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var workerList []*response.WorkerInfo
	sql := `SELECT w.id,w.op_id,w.gate_id,r.room_id,room_name,host_name,device_sn,cpu_model,memory_size,w.ip,miner.actor,
       		sector_size,w.deploy_status,w.start_at,w.finish_at,w.run_status,miner.ip miner_ip,miner_id,w.err_msg,wc.pre_count1
			FROM lotus_worker_info w
			LEFT JOIN lotus_miner_info miner ON w.miner_id = miner.id
			LEFT JOIN sys_host_records r ON w.op_id = r.uuid
			LEFT JOIN lotus_worker_config wc ON w.id = wc.worker_id
			WHERE 1=1  `
	sql += sqlparam + utils.LimitAndOrder("miner.actor", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&workerList).Error; err != nil {
		return nil, 0, err
	}

	for _, v := range workerList {
		gcli := global.GateWayClinets.GetGateWayClinet(v.GateId)
		if gcli == nil {
			global.ZC_LOG.Error(v.GateId + "not exist！")
			continue
		}
		//是否在线
		if online, err := gcli.OpOnline(context.Background(), &pb.String{Value: v.OpId}); err == nil {
			v.Online = online.Value
		} else {
			global.ZC_LOG.Error(v.GateId + ",OpOnline:" + err.Error())
		}

		if v.Online {
			count, err := gcli.GetRunningCount(context.Background(), &pb.OpTask{OpId: v.OpId, TType: define.SealPreCommit1.String()})
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				continue
			}
			v.RunTaskCount = int(count.TCount)
		}

	}
	return workerList, total, err
}

//@author: nathan
//@function: GateWorkerList
//@description: Get machine room worker information
//@param: gateId string
//@return: list []model.LotusWorkerInfo, err error

func (deploy *DeployService) GateWorkerList(gateId string) (list []model.LotusWorkerInfo, err error) {
	sql := `SELECT w.*
			FROM lotus_worker_info w
			WHERE w.deploy_status = ? AND w.gate_id = ?`

	if err = global.ZC_DB.Raw(sql, define.DeployFinish, gateId).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, err
}

//@author: nathan
//@function: GetStorageByActor
//@description: Get node storage information
//@param: actor string
//@return: param []*pb.FilParam, err error

func (deploy *DeployService) GetStorageByActor(actor string) (param []*pb.FilParam, err error) {
	sql := `SELECT w.ip,actor param
			FROM lotus_storage_info w
			LEFT JOIN lotus_miner_info miner ON w.miner_id = miner.id 
			WHERE actor = ? and w.deploy_status = 2`

	if err = global.ZC_DB.Raw(sql, actor).Scan(&param).Error; err != nil {
		return nil, err
	}
	return param, err
}

func (deploy *DeployService) GetWorkerByActor(actor string) (param []*pb.FilParam, err error) {
	sql := `SELECT w.ip,actor param,w.op_id opid
			FROM lotus_worker_info w
			LEFT JOIN lotus_miner_info miner ON w.miner_id = miner.id 
			WHERE actor = ? and w.deploy_status = 2`

	if err = global.ZC_DB.Raw(sql, actor).Scan(&param).Error; err != nil {
		return nil, err
	}
	return param, err
}

// AddWorker
// @author: nathan
// @function: AddWorker
// @description: Add worker
// @Param     model.LotusWorkerInfo
// @return:  error
func (deploy *DeployService) AddWorker(worker *model.LotusWorkerInfo) error {
	if worker.ID == 0 {
		return global.ZC_DB.Save(worker).Error
	}
	return global.ZC_DB.Updates(worker).Error
}

// UpdateWorker
// @author: nathan
// @function: UpdateWorker
// @description: Update worker
// @Param     model.LotusWorkerInfo
// @return:  error
func (deploy *DeployService) UpdateWorker(worker *model.LotusWorkerInfo) error {
	return global.ZC_DB.Updates(worker).Error
}

// GetWorker
// @author: nathan
// @function: GetWorker
// @description: Get worker
// @Param     model.LotusWorkerInfo
// @return:  error
func (deploy *DeployService) GetWorker(id uint64) (model.LotusWorkerInfo, error) {
	var worker model.LotusWorkerInfo
	return worker, global.ZC_DB.Model(model.LotusWorkerInfo{}).Where("id = ?", id).First(&worker).Error
}

//GetRelationWorkerList
//@author: nathan
//@function: GetRelationWorkerList
//@description: Get the associated worker
//@param: info request.PageInfo
//@return: err error, list interface{}

func (deploy *DeployService) GetRelationWorkerList(id request.IDActor) (list interface{}, err error) {
	var workerList []response.RelationWorker
	sqlparam := fmt.Sprintf(" WHERE actor = '%s'", id.Actor)
	if id.Actor == "" {
		sqlparam = fmt.Sprintf(" WHERE w.miner_id = %d", id.ID)
	}
	sql := `SELECT w.id,w.op_id,w.gate_id,r.room_id,room_name,host_name,device_sn,w.ip
			FROM lotus_worker_info w 
			LEFT JOIN lotus_miner_info m ON w.miner_id = m.id
			LEFT JOIN sys_host_records r ON w.op_id = r.uuid
			` + sqlparam
	return workerList, global.ZC_DB.Raw(sql).Scan(&workerList).Error
}

// GetWorkerMiner
// @author: nathan
// @function: GetWorkerMiner
// @description: Gets the miner of the worker connection
// @Param     model.LotusWorkerInfo
// @return:  error
func (deploy *DeployService) GetWorkerMiner(id uint64) (response.ServerInfo, error) {
	var server response.ServerInfo
	sql := `SELECT m.id,m.ip,m.token,m.actor
			FROM lotus_worker_info w 
			LEFT JOIN lotus_miner_info m ON w.miner_id = m.id
			WHERE w.id = ? `
	return server, global.ZC_DB.Raw(sql, id).Scan(&server).Error
}

// GetWorkerByOPId Get work information with op_id
func (deploy *DeployService) GetWorkerByOPId(opId string) (model.LotusWorkerInfo, error) {
	var worker model.LotusWorkerInfo
	return worker, global.ZC_DB.Model(model.LotusWorkerInfo{}).Where("op_id = ?", opId).First(&worker).Error
}

// ModifyWorkerStatus Change worker status
func (deploy *DeployService) ModifyWorkerStatus(worker map[string]int) error {

	ipStr := ""
	for k, _ := range worker {
		if utils.IsNull(ipStr) {
			ipStr = fmt.Sprintf("'%s'", k)
		} else {
			ipStr += fmt.Sprintf(",'%s'", k)
		}
	}

	if ipStr != "" {

		modifyMap := make(map[uint]int)
		var workerList []model.LotusWorkerInfo
		err := global.ZC_DB.Model(&model.LotusWorkerInfo{}).Where("ip in(" + ipStr + ")").Find(&workerList).Error
		if err != nil {
			return err
		}

		if len(workerList) > 0 {

			for _, v := range workerList {
				if val, ok := worker[v.Ip]; ok {
					if val != v.RunStatus {
						modifyMap[v.ID] = val
					}
				}
			}

			if len(modifyMap) > 0 {
				for k, v := range modifyMap {

					err = global.ZC_DB.Model(&model.LotusWorkerInfo{}).Where("id", k).Update("run_status", v).Error
					if err != nil {
						return err
					}
					time.Sleep(time.Second)
				}
			}
		}
	}

	return nil
}

// GetWorkerMonitorList Get a list of monitoring workers
func (deploy *DeployService) GetWorkerMonitorList(info systemReq.HostMonitorReq) (list []response.WorkerMonitorInfo, total int64, err error) {
	sqlparam := " WHERE 1=1 "
	if info.Keyword != "" {
		sqlparam += fmt.Sprintf(" AND (lw.ip like '%%%s%%') ", info.Keyword)
	}
	if info.GateId != "" {
		sqlparam += fmt.Sprintf(" AND lw.gate_id = '%s' ", info.GateId)
	}
	sql := `SELECT lw.op_id,lw.gate_id,lw.ip,lm.actor FROM lotus_worker_info lw LEFT JOIN lotus_miner_info lm ON lm.id = lw.miner_id `
	//求总数
	sqlTotal := `SELECT COUNT(1) FROM lotus_worker_info lw `
	err = global.ZC_DB.Model(&model.LotusMinerInfo{}).Raw(sqlTotal + sqlparam).Count(&total).Error
	sql += sqlparam + utils.LimitAndOrder("actor", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return
}
