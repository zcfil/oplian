package dispatch

import (
	"fmt"
	"oplian/define"
	"oplian/global"
	model "oplian/model/lotus"
	request1 "oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/service/pb"
	"oplian/utils"
	"time"
)

// UpdateSectorStatus
// @author: nathan
// @function: UpdateSectorStatus
// @description: Modify the sector state
// @param: sector *pb.SectorStatus
// @return: error
func (w *DispatchService) UpdateSectorStatus(sector *pb.SectorStatus) (err error) {

	if !w.IsExistNode(sector.Sector.Miner) {
		return define.ErrorNodeNotFound
	}
	if sector.Status == define.Proving {
		var queueId uint64

		p, _ := w.GetSectorPiece(sector.Sector.Miner, sector.Sector.Number)
		if len(p) > 0 {
			queueId = p[0].QueueId

			if err = w.UpdateSectorPieceStatus(sector.Sector.Miner, sector.Sector.Number, define.QueueSectorStatusFinish); err != nil {
				return err
			}
		} else {

			detail, _ := w.GetSectorQueueDetail(sector.Sector.Miner, sector.Sector.Number)
			queueId = detail.QueueId

			if err = w.UpdateSectorQueueDetailStatus(sector.Sector.Miner, sector.Sector.Number, define.QueueSectorStatusFinish); err != nil {
				return err
			}
		}

		if err = w.AddCompleteCountByID(queueId); err != nil {
			return err
		}
	}
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Sector.Miner)
	sqlparam := ""
	if sector.Type > 0 {
		sqlparam += fmt.Sprintf(",sector_type = %d", sector.Type)
	}

	var sectorCount int64
	if err = global.ZC_DB.Table(table).Where("actor=? and sector_id=?", sector.Sector.Miner, sector.Sector.Number).Count(&sectorCount).Error; err != nil {
		return err
	}

	if sectorCount > 0 {
		sql := `UPDATE ` + table + ` SET sector_status = ?` + sqlparam + `,updated_at = NOW(),finish_at = NOW() WHERE sector_id = ? `
		if err = global.ZC_DB.Exec(sql, sector.Status, sector.Sector.Number).Error; err != nil {
			return err
		}
	} else {

		sql := `INSERT INTO ` + table + ` (sector_status,sector_id,actor,sector_type,sector_size)VALUES(?,?,?,?,?)`
		if err = global.ZC_DB.Exec(sql, sector.Status, sector.Sector.Number, sector.Sector.Miner, sector.Type, sector.Size).Error; err != nil {
			return err
		}
	}
	return
}

// AddSectorTicket
// @author: nathan
// @function: AddSectorTicket
// @description: The sector added ticket information
// @param: sector *pb.SectorTicket
// @return: error
func (w *DispatchService) AddSectorTicket(sector *pb.SectorTicket) error {
	if !w.IsExistNode(sector.Sector.Miner) {
		return define.ErrorNodeNotFound
	}
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Sector.Miner)
	sql := `UPDATE ` + table + ` SET ticket = ?,ticket_h = ?,sector_status = ? WHERE sector_id = ? `
	var err error
	if global.ZC_DB.Exec(sql, sector.Ticket, sector.TicketH, sector.Status, sector.Sector.Number).RowsAffected == 0 {
		info, _ := w.GetSectorInfo(sector.Sector.Miner, sector.Sector.Number)
		if info.ID == 0 {
			sql = `INSERT INTO ` + table + ` (ticket,ticket_h,sector_id,actor,sector_status)VALUES(?,?,?,?,?)`
			err = global.ZC_DB.Exec(sql, sector.Ticket, sector.TicketH, sector.Sector.Number, sector.Sector.Miner, sector.Status).Error
		}
	}
	return err
}

// AddSectorCommDR
// @author: nathan
// @function: AddSectorCommDR
// @description: Add PC2 information to sectors
// @param: sector *pb.SectorCommDR
// @return: error
func (w *DispatchService) AddSectorCommDR(sector *pb.SectorCommDR) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Sector.Miner)
	sql := `UPDATE ` + table + ` SET cid_comm_r = ?,cid_comm_d = ?,sector_status = ? WHERE sector_id = ? `
	return global.ZC_DB.Exec(sql, sector.CommR, sector.CommD, sector.Status, sector.Sector.Number).Error
}

// AddSectorWaitSeed
// @author: nathan
// @function: AddSectorWaitSeed
// @description: Adds WaitSeed information to sectors
// @param: sector *pb.SectorSeed
// @return: error
func (w *DispatchService) AddSectorWaitSeed(sector *pb.SectorSeed) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Sector.Miner)
	sql := `UPDATE ` + table + ` SET seed = ?,seed_h = ?,sector_status = ? WHERE sector_id = ? `
	return global.ZC_DB.Exec(sql, sector.Seed, sector.SeedH, sector.Status, sector.Sector.Number).Error
}

// AddSectorCommit2
// @author: nathan
// @function: AddSectorCommit2
// @description: Add c2 information to the sector
// @param: sector *pb.SectorCommit
// @return: error
func (w *DispatchService) AddSectorCommit2(sector *pb.SectorProof) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Sector.Miner)
	sql := `UPDATE ` + table + ` SET proof = ?,sector_status = ? WHERE sector_id = ? `
	return global.ZC_DB.Exec(sql, sector.Proof, sector.Status, sector.Sector.Number).Error
}

// AddSectorPreCID
// @author: nathan
// @function: AddSectorPreCID
// @description: Sector added p2 message ID
// @param: sector *pb.SectorCID
// @return: error
func (w *DispatchService) AddSectorPreCID(sector *pb.SectorCID) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Sector.Miner)
	sql := `UPDATE ` + table + ` SET pre_cid = ?,sector_status = ? WHERE sector_id = ? `
	return global.ZC_DB.Exec(sql, sector.Cid, sector.Status, sector.Sector.Number).Error
}

// AddSectorCommitCID
// @author: nathan
// @function: AddSectorCommitCID
// @description: Sector adds C2 message ID
// @param: sector *pb.SectorCID
// @return: error
func (w *DispatchService) AddSectorCommitCID(sector *pb.SectorCID) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Sector.Miner)
	sql := `UPDATE ` + table + ` SET commit_cid = ?,sector_status = ? WHERE sector_id = ? `
	return global.ZC_DB.Exec(sql, sector.Cid, sector.Status, sector.Sector.Number).Error
}

// CreateSectorTable
// @author: nathan
// @function: CreateSectorTable
// @description: Create a sector table
// @param: minerId string
// @return: error
func (w *DispatchService) CreateSectorTable(minerId string) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), minerId)
	sql := `CREATE TABLE IF NOT EXISTS ` + table + ` SELECT * FROM ` + model.LotusSectorInfo{}.TableName() + ";"
	if err := global.ZC_DB.Exec(sql).Error; err != nil {
		return err
	}

	sql = "ALTER TABLE " + table +
		" MODIFY COLUMN `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT FIRST," +
		"ADD PRIMARY KEY (`id`)," +
		"ADD UNIQUE INDEX `sector_id`(`sector_id`) USING BTREE," +
		"ADD INDEX `idx_lotus_sector_info_sector_id`(`sector_id`) USING BTREE;"
	global.ZC_DB.Exec(sql)
	return nil
}

// CreateSectorLogTable
// @author: nathan
// @function: CreateSectorLogTable
// @description: Create a sector log table
// @param: minerId string
// @return: error
func (w *DispatchService) CreateSectorLogTable(minerId string) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorLog{}.TableName(), minerId)
	sql := `CREATE TABLE IF NOT EXISTS ` + table + ` SELECT * FROM ` + model.LotusSectorLog{}.TableName() + ";"
	if err := global.ZC_DB.Exec(sql).Error; err != nil {
		return err
	}

	sql = "ALTER TABLE " + table +
		" ADD PRIMARY KEY (`id`)," +
		"ADD INDEX `idx_lotus_sector_log_sector_id`(`sector_id`) USING BTREE," +
		"ADD INDEX `idx_sector_status`(`sector_status`) USING BTREE;"
	global.ZC_DB.Exec(sql)
	return nil
}

// CreateSectorPieceTable
// @author: nathan
// @function: CreateSectorPieceTable
// @description: Create the sector order table
// @param: minerId string
// @return: error
func (w *DispatchService) CreateSectorPieceTable(minerId string) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorPiece{}.TableName(), minerId)

	var total int64
	err := global.ZC_DB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_NAME='%s'", table)).Count(&total).Error
	if err != nil {
		return err
	}
	if total == 0 {

		sql := `CREATE TABLE IF NOT EXISTS ` + table + ` SELECT * FROM ` + model.LotusSectorPiece{}.TableName() + ";"
		if err := global.ZC_DB.Exec(sql).Error; err != nil {
			return err
		}

		sql = "ALTER TABLE " + table +
			" MODIFY COLUMN `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT FIRST," +
			"ADD PRIMARY KEY (`id`)," +
			"ADD UNIQUE INDEX `cid_and_sector_id`(`sector_id`, `piece_cid`)," +
			"ADD INDEX `idx_lotus_sector_piece_sector_id`(`sector_id`) USING BTREE," +
			"ADD INDEX `idx_lotus_sector_piece_queue_id`(`queue_id`) USING BTREE;"
		return global.ZC_DB.Exec(sql).Error
	}

	return nil
}

// CreateSectorQueueDetailTable
// @author: nathan
// @function: CreateSectorQueueDetailTable
// @description: Create the sector task queue details table
// @param: minerId string
// @return: error
func (w *DispatchService) CreateSectorQueueDetailTable(minerId string) error {
	table := fmt.Sprintf("%s_%s", model.LotusSectorQueueDetail{}.TableName(), minerId)
	sql := `CREATE TABLE IF NOT EXISTS ` + table + ` SELECT * FROM ` + model.LotusSectorQueueDetail{}.TableName() + ";"
	if err := global.ZC_DB.Exec(sql).Error; err != nil {
		return err
	}

	sql = "ALTER TABLE " + table +
		" MODIFY COLUMN `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT FIRST," +
		"ADD PRIMARY KEY (`id`)," +
		"ADD INDEX `idx_lotus_sector_queue_detail_sector_id`(`sector_id`) USING BTREE," +
		"ADD INDEX `idx_lotus_sector_queue_detail_queue_id`(`queue_id`) USING BTREE;"
	global.ZC_DB.Exec(sql)
	return nil
}

// GetSectorsList
// @author: nathan
// @function: GetSectorsList
// @description: Gets the sector list
// @param: sector *request1.SectorPage
// @return: error
func (w *DispatchService) GetSectorsList(sector request1.SectorPage) (list interface{}, total int64, err error) {
	if !w.IsExistNode(sector.Actor) {
		return nil, 0, define.ErrorNodeNotFound
	}
	sqlparam := fmt.Sprintf(` WHERE actor = '%s'`, sector.Actor)
	if sector.Number != 0 {
		sqlparam += fmt.Sprintf(` AND sector_id = %d `, sector.Number)
	}
	if sector.SectorType != 0 {
		sqlparam += fmt.Sprintf(` AND sector_type = %d `, sector.SectorType)
	}
	if sector.Status != "" {
		sqlparam += fmt.Sprintf(` AND sector_status = '%s' `, sector.Status)
	}

	table := fmt.Sprintf("%s_%s", model.LotusSectorInfo{}.TableName(), sector.Actor)

	sqltotal := `select count(1) from ` + table + sqlparam
	if err = global.ZC_DB.Raw(sqltotal).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sql := fmt.Sprintf(`select id,sector_id,actor,sector_status,sector_size,sector_type,created_at,finish_at from %s `, table)
	var res []response.SectorInfo
	sqlparam += utils.LimitAndOrder("sector_id", "desc", sector.Page, sector.PageSize)
	return res, total, global.ZC_DB.Raw(sql + sqlparam).Scan(&res).Error
}

// IsExistNode
// @author: nathan
// @function: IsExistNode
// @description: Modify the sector state
// @param: miner string
// @return: error
func (w *DispatchService) IsExistNode(miner string) bool {

	time.Sleep(time.Second)
	sql := fmt.Sprintf(`SHOW TABLES LIKE '%s_%s'`, model.LotusSectorInfo{}.TableName(), miner)
	var str string
	if global.ZC_DB.Raw(sql).Scan(&str).RowsAffected > 0 {
		return true
	}
	return false
}

// AddSectorLog
// @author: nathan
// @function: AddSectorLog
// @description: Sectors add order information
// @param: sector *pb.SectorLog
// @return: error

func (w *DispatchService) AddSectorLog(sector *pb.SectorLog) error {
	if !w.IsExistNode(sector.Sector.Miner) {
		return define.ErrorNodeNotFound
	}

	table := fmt.Sprintf("%s_%s", model.LotusSectorLog{}.TableName(), sector.Sector.Miner)
	sql := `INSERT INTO ` + table + ` (id,sector_id,actor,sector_status,error_msg,worker_id,worker_ip)VALUES(?,?,?,?,?,?,?)`
	global.ZC_DB.Exec(sql, sector.ID, sector.Sector.Number, sector.Sector.Miner, sector.SectorStatus, sector.ErrorMsg, sector.WorkerId, sector.WorkerIp)
	return nil
}

// EndSectorLog
// @author: nathan
// @function: EndSectorLog
// @description: Records the sector phase completion time
// @param: id, miner string
// @return: error

func (w *DispatchService) EndSectorLog(id, miner, errMsg string) error {
	if !w.IsExistNode(miner) {
		return define.ErrorNodeNotFound
	}

	table := fmt.Sprintf("%s_%s", model.LotusSectorLog{}.TableName(), miner)
	sql := `UPDATE ` + table + ` SET finish_at = now(),updated_at = now(),error_msg=? WHERE id = ?`
	global.ZC_DB.Exec(sql, errMsg, id)
	return nil
}

// GetSectorDetails
// @author: nathan
// @function: GetSectorDetails
// @description: Gets the sector list
// @param: sector *request1.SectorPage
// @return: error
func (w *DispatchService) GetSectorDetails(sector request1.SectorPage) (res *response.SectorDetails, total int64, err error) {
	if !w.IsExistNode(sector.Actor) {
		return nil, 0, define.ErrorNodeNotFound
	}

	sqltotal := fmt.Sprintf("select count(1) from `%s_%s` where actor = '%s' and sector_id = %d ", model.LotusSectorLog{}.TableName(), sector.Actor, sector.Actor, sector.Number)
	if err = global.ZC_DB.Raw(sqltotal).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	res = new(response.SectorDetails)

	res.SectorInfo, err = w.GetSectorInfo(sector.Actor, sector.Number)
	if err != nil {
		return nil, 0, err
	}

	if res.SectorInfo.SectorType == define.SectorTypeDC {
		pieces, err := w.GetSectorPiece(sector.Actor, sector.Number)
		if err != nil {
			return nil, 0, err
		}
		for _, v := range pieces {
			res.Piece = append(res.Piece, v.PieceCid)
		}
	}

	var logs []model.LotusSectorLog

	sql := fmt.Sprintf("select * from `%s_%s` where actor = '%s' and sector_id = %d ", model.LotusSectorLog{}.TableName(), sector.Actor, sector.Actor, sector.Number)
	if sector.PageSize != 0 && sector.Page != 0 {
		sql += utils.LimitAndOrder("created_at", "desc", sector.Page, sector.PageSize)
		if err = global.ZC_DB.Raw(sql).Scan(&logs).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err = global.ZC_DB.Raw(sql + " order by created_at desc").Scan(&logs).Error; err != nil {
			return nil, 0, err
		}
	}
	res.SectorLog = logs

	return res, total, nil
}

// GetSectorInfo
// @author: nathan
// @function: GetSectorInfo
// @description: Get sector information
// @param: actor string, sectorId uint64
// @return: sector *model.LotusSectorInfo, err error
func (w *DispatchService) GetSectorInfo(actor string, sectorId uint64) (sector response.SectorInfo, err error) {
	if !w.IsExistNode(actor) {
		return sector, define.ErrorNodeNotFound
	}
	sql := fmt.Sprintf("SELECT id,created_at,updated_at,sector_id,actor,sector_status,sector_type,sector_size,IF(sector_status='Proving',finish_at,NOW()) finish_at FROM `%s_%s` WHERE sector_id = %d", model.LotusSectorInfo{}.TableName(), actor, sectorId)

	return sector, global.ZC_DB.Raw(sql).Scan(&sector).Error
}
