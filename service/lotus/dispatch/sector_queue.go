package dispatch

import (
	"errors"
	"fmt"
	"oplian/define"
	"oplian/global"
	"oplian/model/common/request"
	model "oplian/model/lotus"
	"oplian/service/pb"
	utils "oplian/utils"
	"strconv"
	"time"
)

// GetSectorTaskList
// @author: nathan
// @function: GetSectorTaskList
// @description: Gets the sector task list
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error
func (w *DispatchService) GetSectorTaskList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	condition := ""
	var param []interface{}
	if info.Keyword != "" {
		condition = "where actor = ?"
		param = append(param, info.Keyword, info.Keyword)
	}
	sql := ` select id, task_name,created_at,finish_at,sector_size,sector_type,actor,job_total,run_count,complete_count,concurrent_import,task_status
			from lotus_sector_queue ` + condition + `
			union all
			select id, task_name,begin_time,end_time,sector_size,sector_type,actor,job_total,run_count,finish_count,1,task_status
			from lotus_sector_task  ` + condition

	sqlQuery := fmt.Sprintf("select x.* from(%s)x order by x.created_at desc limit %d,%d ", sql, offset, limit)
	sqlTotal := fmt.Sprintf("select count(1) as total from(%s)x", sql)
	var res []model.LotusSectorQueue
	err = global.ZC_DB.Raw(sqlQuery, param...).Scan(&res).Error
	if err != nil {
		return
	}

	err = global.ZC_DB.Raw(sqlTotal, param...).Count(&total).Error
	if err != nil {
		return
	}

	return res, total, err
}

// AddSectorTaskQueue
// @author: nathan
// @function: AddSectorTaskQueue
// @description: Sector task queue
// @param: info *model.LotusSectorQueue
// @return: err error
func (w *DispatchService) AddSectorTaskQueue(info *model.LotusSectorQueue) (err error) {
	return global.ZC_DB.Save(info).Error
}

// UpdateSectorTaskQueue
// @author: nathan
// @function: UpdateSectorTaskQueue
// @description: Sector task queue
// @param: info *model.LotusSectorQueue
// @return: err error
func (w *DispatchService) UpdateSectorTaskQueue(info *model.LotusSectorQueue) (err error) {
	return global.ZC_DB.Updates(info).Error
}

// CheckDealCar
// @author: nathan
// @function: CheckDealCar
// @description: Sector task queue
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error
func (w *DispatchService) CheckDealCar(deal map[string][]string, actor string, queueId uint64) (err error) {
	if !w.IsExistNode(actor) {
		return define.ErrorNodeNotFound
	}

	table := fmt.Sprintf("%s_%s", model.LotusSectorPiece{}.TableName(), actor)
	var i int
	sql := `INSERT INTO ` + table + ` (actor,queue_id,run_index,job_status,expiration_time,deal_uuid,piece_cid)VALUES`
	for uuid, cid_at := range deal {
		if len(cid_at) < 2 {
			continue
		}
		i++
		epoch, _ := strconv.ParseInt(cid_at[1], 10, 64)
		end := utils.BlockHeightToTime(epoch).Format("2006-01-02 15:04:05")
		if i < len(deal)-1 {
			sql += fmt.Sprintf("('%s',%d,%d,%d,'%s','%s','%s'),", actor, queueId, i, define.QueueSectorStatusWait, end, uuid, cid_at[0])
		} else {
			sql += fmt.Sprintf("('%s',%d,%d,%d,'%s','%s','%s')", actor, queueId, i, define.QueueSectorStatusWait, end, uuid, cid_at[0])
			return global.ZC_DB.Exec(sql).Error
		}

	}
	return errors.New("Import order data is 0")
}

// EditTaskQueueStatus
// @author: nathan
// @function: EditTaskQueueStatus
// @description: Modify task status
// @param: id uint64, status int
// @return: err error
func (w *DispatchService) EditTaskQueueStatus(id int, status int) (err error) {
	return global.ZC_DB.Model(model.LotusSectorQueue{}).Where("id = ?", id).Update("task_status", status).Error
}

// EditTaskQueueDetailStatus
// @author: nathan
// @function: EditTaskQueueDetailStatus
// @description: Modify the order task status
// @param: info request.ActorIdStatus
// @return: err error
func (w *DispatchService) EditTaskQueueDetailStatus(info request.ActorIdStatus, opId string) (err error) {
	if !w.IsExistNode(info.Actor) {
		return define.ErrorNodeNotFound
	}
	sqlparam := ""
	if info.Value != "" && info.Status == define.QueueSectorStatusWait {
		sqlparam = fmt.Sprintf(",car_path = '%s'", info.Value)
	}
	if opId != "" {
		sqlparam += fmt.Sprintf(",car_op_id = '%s'", opId)
	}
	sql := fmt.Sprintf(`UPDATE %s_%s SET job_status = ? %s WHERE id = ?`, model.LotusSectorPiece{}.TableName(), info.Actor, sqlparam)
	return global.ZC_DB.Exec(sql, info.Status, info.ID).Error
}

// EditTaskQueueDetailStatusBatch
// @author: nathan
// @function: EditTaskQueueDetailStatusBatch
// @description: Batch modify order task status
// @param: info request.ActorIdStatus
// @return: err error
func (w *DispatchService) EditTaskQueueDetailStatusBatch(infos []request.ActorIdStatus) (err error) {
	DB := global.ZC_DB.Begin()
	defer func() {
		if err != nil {
			DB.Rollback()
			return
		}
		DB.Commit()
	}()
	for _, info := range infos {
		if !w.IsExistNode(info.Actor) {
			return define.ErrorNodeNotFound
		}
		sqlparam := ""
		if info.Value != "" && info.Status == define.QueueSectorStatusWait {
			sqlparam = fmt.Sprintf(",car_path = '%s'", info.Value)
		}
		sql := fmt.Sprintf(`UPDATE %s_%s SET job_status = ? %s WHERE id = ?`, model.LotusSectorPiece{}.TableName(), info.Actor, sqlparam)
		err = DB.Exec(sql, info.Status, info.ID).Error
	}
	return err
}

// EditTaskQueueDetailConcurrent
// @author: nathan
// @function: EditTaskQueueDetailConcurrent
// @description: Modify the order task status
// @param: info request.ActorIdStatus
// @return: err error
func (w *DispatchService) EditTaskQueueDetailConcurrent(info request.IdCount) (err error) {
	return global.ZC_DB.Model(model.LotusSectorQueue{}).Where("id = ?", info.ID).Update("concurrent_import", info.Count).Error
}

// GetRunningTaskQueue
// @author: nathan
// @function: GetRunningTaskQueue
// @description: Gets a list of sector tasks in progress
// @param: actor string
// @return: res []*pb.TaskQueue, err error
func (w *DispatchService) GetRunningTaskQueue(actor string) (res []*pb.TaskQueue, err error) {
	sql := `select q.*,m.ip miner_ip,m.token miner_token from lotus_sector_queue q
				LEFT JOIN lotus_miner_info m ON q.actor = m.actor AND is_manage = 1
				where q.actor = ? and task_status=?`
	return res, global.ZC_DB.Raw(sql, actor, define.QueueStatusRun).Scan(&res).Error
}

// AddCompleteCountByID
// @author: nathan
// @function: AddCompleteCountByID
// @description: Record the completion of a sector
// @param: id uint64
// @return:  err error
func (w *DispatchService) AddCompleteCountByID(id uint64) (err error) {

	var date model.LotusSectorQueue
	if err = global.ZC_DB.Where("id = ?", id).Find(&date).Error; err != nil {
		return err
	}
	if (date != model.LotusSectorQueue{}) {
		date.CompleteCount++
		date.RunCount--
		if date.RunCount < 0 {
			date.RunCount = 0
		}
		if date.CompleteCount >= date.JobTotal {
			date.FinishAt = time.Now()
			date.CompleteCount = date.JobTotal
			date.TaskStatus = define.QueueStatusFinish
		}
		return global.ZC_DB.Updates(&date).Error
	}

	return nil
}

// AddRunCountByID
// @author: nathan
// @function: AddRunCountByID
// @description: Record the start of a sector
// @param: id uint64
// @return: err error
func (w *DispatchService) AddRunCountByID(id uint64) (err error) {
	var date model.LotusSectorQueue
	if err = global.ZC_DB.Where("id = ?", id).First(&date).Error; err != nil {
		return err
	}
	date.RunCount++
	return global.ZC_DB.Updates(&date).Error
}
