package dispatch

import (
	"errors"
	"fmt"
	"oplian/define"
	"oplian/global"
	model "oplian/model/lotus"
	request1 "oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/service/pb"
	"oplian/utils"
)

// GetSectorQueueDetail
// @author: nathan
// @function: GetSectorQueueDetail
// @description: Get sector queue details
// @param: miner string,sectorId uint64
// @return: in uint64, err error
func (w *DispatchService) GetSectorQueueDetail(miner string, sectorId uint64) (res model.LotusSectorQueueDetail, err error) {
	table := fmt.Sprintf("%s_%s", model.LotusSectorQueueDetail{}.TableName(), miner)
	//sql := `SELECT ifnull(MAX(run_index),0)run_index FROM ` + table + ` WHERE queue_id = ? `
	//global.ZC_DB.Raw(sql, queueId).Scan(&runIndex)
	//return runIndex + 1
	return res, global.ZC_DB.Table(table).Where("sector_id = ?", sectorId).Find(&res).Error
}

// GetSectorQueueDetailList
// @author: nathan
// @function: GetSectorQueueDetailList
// @description: Gets a list of sector queue details
// @param: sectorId uint64
// @return: sector model.LotusSectorQueueDetail, err error
func (w *DispatchService) GetSectorQueueDetailList(info request1.DealPage) (list interface{}, total int64, err error) {
	sqlparam := ""
	if info.Status == define.QueueSectorStatusExceed {
		info.Status = 0
		sqlparam += fmt.Sprintf(" AND expiration_time < NOW() AND job_status<> %d", define.QueueSectorStatusFinish)
	}
	if info.Status != 0 {
		sqlparam += fmt.Sprintf(" AND job_status = %d", info.Status)
	}
	var sql, sqltotal string
	//dc or cc
	switch info.SectorType {
	case define.SectorTypeDC:
		dcTable := fmt.Sprintf("%s_%s", model.LotusSectorPiece{}.TableName(), info.Actor)
		sql = `SELECT qd.id,qd.sector_id,qd.actor,IF(expiration_time<NOW(),IF(qd.job_status=3,3,6),qd.job_status) job_status,qd.run_index,qd.expiration_time,CONCAT(qd.car_path,'/',qd.piece_cid,'.car') car_path,qd.piece_cid,qd.deal_uuid,qd.created_at
							,q.sector_size,q.job_total,q.sector_type,q.task_name,q.created_at task_create_at
			                FROM ` + dcTable + ` qd
			                LEFT JOIN lotus_sector_queue q ON qd.queue_id = q.id  
							WHERE queue_id = ? ` + sqlparam
		sqltotal = `SELECT COUNT(1) FROM ` + dcTable + ` WHERE queue_id = ?` + sqlparam
	case define.SectorTypeCC:
		ccTable := fmt.Sprintf("%s_%s", model.LotusSectorQueueDetail{}.TableName(), info.Actor)
		sql = `SELECT qd.sector_id,qd.actor,qd.job_status,qd.run_index,qd.created_at
							,q.sector_size,q.job_total,q.sector_type,q.task_name,q.created_at task_create_at
			                FROM ` + ccTable + ` qd
			                LEFT JOIN lotus_sector_queue q ON qd.queue_id = q.id
							WHERE queue_id = ? ` + sqlparam
		sqltotal = `SELECT COUNT(1) FROM ` + ccTable + ` WHERE queue_id = ? ` + sqlparam
	}

	if err = global.ZC_DB.Raw(sqltotal, info.QueueId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sql += utils.LimitAndOrder("qd.created_at", "desc", info.Page, info.PageSize)
	var detail []response.SectorQueueDetail
	return detail, total, global.ZC_DB.Raw(sql, info.QueueId).Scan(&detail).Error
}

// GetSectorRecoverDetail
// @author: GetSectorRecoverDetail
// @function: GetSectorRecoverDetail
// @description: Get the sector recovery details
// @param: sectorId uint64
// @return: sector model.LotusSectorQueueDetail, err error
func (w *DispatchService) GetSectorRecoverDetail(info request1.SectorRecoverDetail) (task, list interface{}, total int64, err error) {

	var lotusSectorTask model.LotusSectorTask
	err = global.ZC_DB.Model(&model.LotusSectorTask{}).Where("id", info.Id).Find(&lotusSectorTask).Error
	if err != nil {
		return nil, nil, 0, err
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.ZC_DB.Model(&model.LotusSectorTaskDetail{}).Where("relation_id", info.Id)
	var res []model.LotusSectorTaskDetail
	err = db.Count(&total).Error
	if err != nil {
		return nil, nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&res).Error
	if err != nil {
		return nil, nil, 0, err
	}

	return lotusSectorTask, res, total, nil
}

// AddSectorQueueDetail
// @author: nathan
// @function: AddSectorQueueDetail
// @description: Added sector queue details
// @param: sectorId uint64
// @return: err error
func (w *DispatchService) AddSectorQueueDetail(args *pb.SectorQueueDetail) (err error) {
	if !w.IsExistNode(args.Sid.Miner) {
		return define.ErrorNodeNotFound
	}
	table := fmt.Sprintf("%s_%s", model.LotusSectorQueueDetail{}.TableName(), args.Sid.Miner)
	sql := `INSERT INTO ` + table + ` (sector_id,actor,queue_id,run_index,job_status)VALUE(?,?,?,?,?)`
	return global.ZC_DB.Exec(sql, args.Sid.Number, args.Sid.Miner, args.QueueId, w.GetNextRunIndex(args.Sid.Miner, args.QueueId), args.SectorStatus).Error
}

// AddDCSectorQueueDetail
// @author: nathan
// @function: AddDCSectorQueueDetail
// @description: Added DC sector queue details
// @param: sectorId uint64
// @return: err error
func (w *DispatchService) AddDCSectorQueueDetail(dealMap map[string]*request1.DealInfo, actor string, queueId uint) (err error) {
	if !w.IsExistNode(actor) {
		return define.ErrorNodeNotFound
	}

	table := fmt.Sprintf("%s_%s", model.LotusSectorPiece{}.TableName(), actor)
	var i int
	sql := `INSERT INTO ` + table + ` (actor,queue_id,run_index,job_status,expiration_time,deal_uuid,piece_cid,car_path,car_op_id)VALUES`
	for uuid, deal := range dealMap {
		i++
		end := utils.BlockHeightToTime(deal.EndEpoch).Format("2006-01-02 15:04:05")
		if i < len(dealMap) {
			sql += fmt.Sprintf("('%s',%d,%d,%d,'%s','%s','%s','%s','%s'),", actor, queueId, i, deal.JobStatus, end, uuid, deal.PieceCid, deal.CarPath, deal.FileOpId)
		} else {
			sql += fmt.Sprintf("('%s',%d,%d,%d,'%s','%s','%s','%s','%s')", actor, queueId, i, deal.JobStatus, end, uuid, deal.PieceCid, deal.CarPath, deal.FileOpId)
			return global.ZC_DB.Exec(sql).Error
		}
	}
	return errors.New("Import order data is 0")
}

// AddSectorPiece
// @author: nathan
// @function: AddSectorPiece
// @description: Sectors add order information
// @param: sector *pb.SectorPiece
// @return: error
func (w *DispatchService) AddSectorPiece(sector *pb.SectorPiece) error {

	if !w.IsExistNode(sector.Sector.Miner) {
		return define.ErrorNodeNotFound
	}
	table := fmt.Sprintf("%s_%s", model.LotusSectorPiece{}.TableName(), sector.Sector.Miner)

	sql := `UPDATE ` + table + ` SET deal_id = ?,piece_size = ?,sector_id = ?,job_status = 2 WHERE piece_cid = ? and actor = ?`
	return global.ZC_DB.Exec(sql, sector.DealId, sector.PieceSize, sector.Sector.Number, sector.PieceCid, sector.Sector.Miner).Error
}

// UpdateSectorQueueDetailStatus
// @author: nathan
// @function: UpdateSectorQueueDetailStatus
// @description: Modify the sector queue detail state
// @param: sectorId uint64
// @return: err error
func (w *DispatchService) UpdateSectorQueueDetailStatus(miner string, sectorId uint64, status int) (err error) {
	table := fmt.Sprintf("%s_%s", model.LotusSectorQueueDetail{}.TableName(), miner)
	return global.ZC_DB.Table(table).Where("sector_id = ?", sectorId).Update("job_status", status).Error
}

// UpdateSectorPieceStatus
// @author: nathan
// @function: UpdateSectorPieceStatus
// @description: Modify the DC sector queue detail state
// @param: sectorId uint64
// @return: err error
func (w *DispatchService) UpdateSectorPieceStatus(miner string, sectorId uint64, status int) (err error) {
	table := fmt.Sprintf("%s_%s", model.LotusSectorPiece{}.TableName(), miner)
	return global.ZC_DB.Table(table).Where("sector_id = ?", sectorId).Update("job_status", status).Error
}

// GetNextRunIndex
// @author: nathan
// @function: GetNextRunIndex
// @description: Gets the task queue execution index
// @param: queueId uint64
// @return: index int
func (w *DispatchService) GetNextRunIndex(miner string, queueId uint64) (runIndex int) {
	table := fmt.Sprintf("%s_%s", model.LotusSectorQueueDetail{}.TableName(), miner)
	sql := `SELECT ifnull(MAX(run_index),0)run_index FROM ` + table + ` WHERE queue_id = ? `
	global.ZC_DB.Raw(sql, queueId).Scan(&runIndex)
	return runIndex + 1
}

// GetSectorPiece
// @author: nathan
// @function: GetSectorPiece
// @description: Get sector order information
// @param: actor string, sectorId uint64
// @return: sector *model.GetSectorPiece, err error
func (w *DispatchService) GetSectorPiece(actor string, sectorId uint64) (sector []model.LotusSectorPiece, err error) {
	if !w.IsExistNode(actor) {
		return nil, define.ErrorNodeNotFound
	}
	sql := fmt.Sprintf("SELECT * FROM `%s_%s` WHERE sector_id = %d", model.LotusSectorPiece{}.TableName(), actor, sectorId)
	return sector, global.ZC_DB.Raw(sql).Scan(&sector).Error
}

// GetWaitImportDeal
// @author: nathan
// @function: GetWaitImportDeal
// @description: Get sector order information
// @param: info *pb.DealInfo
// @return: deal []model.LotusSectorPiece, err error
func (w *DispatchService) GetWaitImportDeal(info *pb.DealParam) (deals []*pb.DealInfo, err error) {
	if !w.IsExistNode(info.Actor) {
		return nil, define.ErrorNodeNotFound
	}
	sql := fmt.Sprintf("SELECT * FROM `%s_%s` WHERE job_status = ? AND queue_id = ? AND expiration_time>NOW() ORDER BY run_index LIMIT %d", model.LotusSectorPiece{}.TableName(), info.Actor, info.Count)
	return deals, global.ZC_DB.Raw(sql, define.QueueSectorStatusWait, info.QueueId).Scan(&deals).Error
}
