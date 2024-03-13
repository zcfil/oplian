package deploy

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"oplian/lotusrpc"
	"oplian/model/lotus/response"
	"sync"
	"time"

	"gorm.io/gorm"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/system"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path"
	"strconv"
	"strings"
)

var WorkerClusterServiceApi = new(WorkerClusterService)

var WorkerCluster = WorkerClusterService{
	C2TaskMap: make(map[uint]response.C2TaskInfo),
}

type WorkerClusterService struct {
	C2TaskMap map[uint]response.C2TaskInfo
}

func (w *WorkerClusterService) GetWorkerClusterList(info request.WorkerCluster) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&lotus.LotusWorkerCluster{})
	var workerCluster []lotus.LotusWorkerCluster
	sql := "1=1"
	var param []interface{}
	if info.OpKeyWord != "" {
		info.OpKeyWord = "%" + info.OpKeyWord + "%"
		sql += " and (ip like ? or server_name like ?)"
		param = append(param, info.OpKeyWord, info.OpKeyWord)
	}
	if info.RoomKeyWord != "" {
		info.RoomKeyWord = "%" + info.RoomKeyWord + "%"
		sql += " and (room_num like ? or room_name like ?)"
		param = append(param, info.RoomKeyWord, info.RoomKeyWord)
	}
	if info.SectorSize > 0 {
		sql += " and sector_size=?"
		param = append(param, info.SectorSize)
	}
	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&workerCluster).Error
		if err != nil {
			return nil, 0, err
		}
	}

	ids := ""
	var workerClusterList []lotus.LotusWorkerCluster
	for _, v := range workerCluster {
		if ids == "" {
			ids = fmt.Sprintf("%d", v.ID)
		} else {
			ids += fmt.Sprintf(",%d", v.ID)
		}
	}

	if ids != "" {
		resMap, err := w.GetC2FinishCount(ids)
		if err != nil {
			return nil, 0, err
		}

		if len(resMap) > 0 {
			for _, v := range workerCluster {
				if total, ok := resMap[fmt.Sprintf("%d", v.ID)]; ok {
					v.TaskNum = total
				} else {
					v.TaskNum = 0
				}
				workerClusterList = append(workerClusterList, v)
			}
		}
	}

	return workerClusterList, total, err
}

func (w *WorkerClusterService) GetC2FinishCount(ids string) (map[string]int, error) {

	if utils.IsNull(ids) {
		return nil, errors.New("ids Cannot be empty")
	}

	resMap := make(map[string]int)
	type FinishCount struct {
		RelationId string `json:"relationId"`
		Total      int    `json:"total"`
	}

	tableNameList := make([]string, 0)
	yearMonth := w.GetTableYearMonth()

	if len(yearMonth) > 0 {
		for _, v := range yearMonth {

			tableName := fmt.Sprintf("%s_%s", lotus.LotusWorkerRelations{}.TableName(), v)
			var total int64
			sql := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_NAME='%s'", tableName)
			err := global.ZC_DB.Raw(sql).Count(&total).Error
			if err != nil {
				return nil, err
			}
			if total > 0 {
				sql := `SHOW TABLES LIKE '%` + tableName + `%'`
				var str string
				if global.ZC_DB.Raw(sql).Scan(&str).RowsAffected > 0 {
					tableNameList = append(tableNameList, tableName)
				}
			}
		}
	}

	if len(tableNameList) > 0 {

		idsList := make([]int, 0)
		if strings.Contains(ids, ",") {
			idsAr := strings.Split(ids, ",")
			for _, v := range idsAr {
				id, _ := strconv.Atoi(v)
				idsList = append(idsList, id)
			}
		} else {
			id, _ := strconv.Atoi(ids)
			idsList = append(idsList, id)
		}

		for _, v := range idsList {

			sqlStr := ""
			for index, v1 := range tableNameList {
				sql := fmt.Sprintf("SELECT r.relation_id,COUNT(1) AS total FROM %s r WHERE r.relation_id=%s AND r.task_status=1 AND r.sync_status=1 AND r.is_remove=1", v1, strconv.Itoa(v))
				if index == len(tableNameList)-1 {
					sqlStr += sql
				} else {
					sqlStr += sql + " UNION ALL "
				}
			}

			if sqlStr != "" {

				var fc []FinishCount
				err := global.ZC_DB.Raw(sqlStr).Scan(&fc).Error
				if err != nil {
					return nil, err
				}

				for _, v2 := range fc {

					total, ok := resMap[v2.RelationId]
					if !ok {
						resMap[v2.RelationId] = v2.Total
					} else {
						resMap[v2.RelationId] = v2.Total + total
					}
				}
			}
		}
	}

	return resMap, nil
}

func (w *WorkerClusterService) GetTableYearMonth() []string {

	defaultYearMonth := "202312"
	var yearMonth []string
	beginIndex := 0
	beginYear, _ := strconv.Atoi(time.Now().AddDate(-99, 0, 0).Format("2006"))
	endYear, _ := strconv.Atoi(time.Now().Format("2006"))
	curMonth, _ := strconv.Atoi(time.Now().Format("01"))
	for i := beginYear; i <= endYear; i++ {
		j := 0
		for {
			j++
			monthStr := ""
			if j < 10 {
				monthStr = fmt.Sprintf("0%d", j)
			} else {
				monthStr = fmt.Sprintf("%d", j)
			}
			yearMonthStr := fmt.Sprintf("%d", i)
			yearMonthStr += monthStr
			yearMonth = append(yearMonth, yearMonthStr)
			if i == endYear && j == curMonth {
				break
			}
			if j == 12 {
				break
			}
		}
	}

	for i, v := range yearMonth {
		if v == defaultYearMonth {
			beginIndex = i
		}
	}

	var resYearMonth []string
	for i, v := range yearMonth {
		if i >= beginIndex {
			resYearMonth = append(resYearMonth, v)
		}
	}

	return resYearMonth
}

// AddWorkerOp
// @function: AddWorkerOp
// @description: Add Worker
// @param: info request.WorkerCluster
// @return: error
func (w *WorkerClusterService) AddWorkerOp(info request.AddWorkerCluster) error {

	opIdStr := ""
	param := make([]lotus.LotusWorkerCluster, 0)
	for _, v := range info.WorkerCluster {
		param = append(param, v)
		opIdStr += `'` + v.OpId + `',`
	}
	if len(param) == 0 {
		return errors.New("worker_cluster[] data is nil")
	}

	global.ZC_DB.Transaction(func(tx *gorm.DB) error {

		if opIdStr != "" {

			opIdStr = utils.SubStr(opIdStr, utils.ZERO, len(opIdStr)-utils.ONE)
			var workerList []lotus.LotusWorkerCluster
			err := tx.Model(&lotus.LotusWorkerCluster{}).Where("op_id in(" + opIdStr + ")").Find(&workerList).Error
			if err != nil {
				return err
			}

			errMsg := ""
			for _, v := range workerList {
				if !strings.Contains(errMsg, v.Ip) {
					errMsg += fmt.Sprintf("%s/%s,", v.Ip, v.ServerName)
				}
			}

			if errMsg != "" {
				errMsg = fmt.Sprintf("worker主机:%s已存在,不可重复增加", utils.SubStr(errMsg, utils.ZERO, len(errMsg)-utils.ONE))
				return errors.New(errMsg)
			}

			var workerCount int64
			err = global.ZC_DB.Model(&lotus.LotusWorkerInfo{}).Where("op_id in(" + opIdStr + ")").Count(&workerCount).Error
			if err != nil {
				return err
			}

			if workerCount > 0 {

				err = tx.Model(&system.SysHostRecord{}).Where("uuid in("+opIdStr+")").Update("host_classify", 5).Error
				if err != nil {
					return err
				}
			}
		}

		err := tx.Model(&lotus.LotusWorkerCluster{}).CreateInBatches(&param, len(param)).Error
		if err != nil {
			return err
		}

		return nil
	})

	var hostRecord system.SysHostRecord
	if len(param) > 0 {

		err := global.ZC_DB.Model(&system.SysHostRecord{}).Where("uuid", param[0].OpId).Find(&hostRecord).Error
		if err != nil {
			return err
		}
	}

	if (hostRecord == system.SysHostRecord{}) {
		return errors.New(fmt.Sprintf("opId:%s,Abnormal, can not find the relevant data", param[0].OpId))
	}

	client := global.GateWayClinets.GetGateWayClinet(hostRecord.GatewayId)
	if client == nil {
		return errors.New(fmt.Sprintf("GetGateWayClinet Connection failed:%s", global.GateWayID.String()))
	}

	var wg sync.WaitGroup
	wg.Add(len(param))
	for _, v := range param {

		go func(OpId string) {
			defer func() {
				wg.Done()
			}()
			if OpId != "" {
				_, err := client.RunOpC2(context.TODO(), &pb.String{Value: OpId})
				if err != nil {
					return
				}
			}
		}(v.OpId)
	}
	wg.Wait()

	return nil
}

// DelWorkerOp
// @function: DelWorkerOp
// @description: Delete worker
// @param: info lotus.LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) DelWorkerOp(info request.WorkerOp) error {

	if info.ID == 0 {
		return errors.New("id is nil")
	}

	var cluster lotus.LotusWorkerCluster
	err := global.ZC_DB.Model(&lotus.LotusWorkerCluster{}).Where("id", info.ID).Find(&cluster).Error
	if err != nil {
		return err
	}

	if (cluster == lotus.LotusWorkerCluster{}) {
		return errors.New(fmt.Sprintf("%d 异常,找不到相关数据", info.ID))
	}

	var hostRecord system.SysHostRecord
	err = global.ZC_DB.Model(&system.SysHostRecord{}).Where("uuid", cluster.OpId).Find(&hostRecord).Error
	if err != nil {
		return err
	}

	if (hostRecord == system.SysHostRecord{}) {
		return errors.New(fmt.Sprintf("opId:%s,异常,找不到相关数据", cluster.OpId))
	}

	client := global.GateWayClinets.GetGateWayClinet(hostRecord.GatewayId)
	if client == nil {
		return errors.New(fmt.Sprintf("GetGateWayClinet Connection failed:%s", global.GateWayID.String()))
	}

	client.StopOpC2(context.TODO(), &pb.String{Value: cluster.OpId})

	var workerCount int64
	err = global.ZC_DB.Model(&lotus.LotusWorkerInfo{}).Where("op_id", cluster.OpId).Count(&workerCount).Error
	if err != nil {
		return err
	}

	if workerCount > 0 {

		err = global.ZC_DB.Model(&system.SysHostRecord{}).Where("uuid=?", cluster.OpId).Update("host_classify", 2).Error
		if err != nil {
			return err
		}
	}

	return global.ZC_DB.Model(&lotus.LotusWorkerCluster{}).Delete("id", info.ID).Error
}

// WorkerTaskDetail
// @function: WorkerTaskDetail
// @description: worker task details
// @param: info LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) WorkerTaskDetail(info request.WorkerOp) (list []lotus.LotusWorkerRelations, total int64, err error) {

	if info.YearMonth == "" {
		info.YearMonth = utils.TimeToFormat(time.Now(), utils.YearMonth)
	} else {
		info.YearMonth = strings.Replace(info.YearMonth, "-", "", -1)
	}

	tableName, err := w.GetWorkerRelationsTable(info.YearMonth)
	if err != nil {
		return nil, 0, err
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Table(tableName)

	var param []interface{}
	sql := `1=1 `
	if info.ID > 0 {
		sql += " and relation_id = ? "
		param = append(param, info.ID)
	}
	if info.OpKeyWord != "" {
		info.OpKeyWord = "%" + info.OpKeyWord + "%"
		sql += " and task_id like ? "
		param = append(param, info.OpKeyWord)
	}
	if info.RoomId != "" {
		sql += "and room_id = ? "
		param = append(param, info.RoomId)
	}
	if info.SectorSize != "" {
		sql += "and sector_size = ? "
		param = append(param, info.SectorSize)
	}
	if info.TaskStatus != "" {
		sql += "and task_status = ? "
		param = append(param, info.TaskStatus)
	}

	var lotusWorkerRelations []lotus.LotusWorkerRelations
	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		if strings.Contains(err.Error(), "lotus_worker_relations") {
			return nil, 0, nil
		}
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&lotusWorkerRelations).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return lotusWorkerRelations, total, nil
}

// AddC2Task
// @function: AddC2Task
// @description: Add C2 Task
// @param: info LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) AddC2Task(info request.C2TaskInfo) (string, error) {

	tableName, err := w.GetWorkerRelationsTable("")
	if err != nil {
		return "", err
	}

	var workerRelations lotus.LotusWorkerRelations
	err = global.ZC_DB.Table(tableName).Where("miner=? and number=?", info.Miner, info.Number).Find(&workerRelations).Error
	if err != nil {
		return "", err
	}

	if (workerRelations == lotus.LotusWorkerRelations{}) {

		param := &lotus.LotusWorkerRelations{
			TaskId:       utils.GetUid(100000000),
			GateId:       global.GateWayID.String(),
			Miner:        info.Miner,
			Number:       info.Number,
			SectorSize:   32,
			RelationType: utils.ONE,
			TaskStatus:   utils.FOUR,
		}

		err = global.ZC_DB.Table(tableName).Create(&param).Error
		if err != nil {
			return "", err
		}

		return strconv.Itoa(int(param.ZC_MODEL.ID)), err
	}

	return strconv.Itoa(int(workerRelations.ID)), err
}

func (w *WorkerClusterService) GetWorkerRelationsTable(yearMonth string) (string, error) {

	if yearMonth == "" {
		yearMonth = utils.TimeToFormat(time.Now(), utils.YearMonth)
	}
	tableName := fmt.Sprintf("%s_%s", lotus.LotusWorkerRelations{}.TableName(), yearMonth)
	var total int64
	sql := fmt.Sprintf("select count(1) from %s limit 1", tableName)
	err := global.ZC_DB.Raw(sql).Count(&total).Error
	if err != nil {
		col := "id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,"
		col += "INDEX deleted_at (deleted_at) USING BTREE,"
		col += "INDEX relation_id (relation_id) USING BTREE,"
		col += "INDEX op_id (op_id) USING BTREE,"
		col += "INDEX room_id (room_id) USING BTREE,"
		col += "INDEX task_status (task_status) USING BTREE,"
		col += "UNIQUE INDEX miner_number (miner,number) USING BTREE"
		sql := fmt.Sprintf("CREATE TABLE %s (%s) SELECT * FROM %s ", tableName, col, lotus.LotusWorkerRelations{}.TableName())
		err = global.ZC_DB.Exec(sql).Error
		if err != nil {
			return "", nil
		}
	}

	return tableName, nil
}

func FileToByte(fileName, filePath, opId string, client pb.GateServiceClient) (chunks []byte, err error) {

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(f)
	opFilePathC1 := define.PathIpfsData + "/c2task" + define.OpCsPathC1

	client.DelOpFile(context.TODO(), &pb.FileInfo{OpId: opId, FileName: opFilePathC1 + "/" + fileName})

	for {

		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}

		f := &pb.ScriptInfo{
			OpId:     opId,
			FileData: dataByte[:n],
			FileName: fileName,
			Path:     opFilePathC1,
		}

		_, err = client.FileDistribution(context.TODO(), f)
		if err != nil {
			break
		}
		chunks = append(chunks, dataByte[:n]...)
	}
	isEOF := strings.Compare(err.Error(), "EOF")
	if isEOF == 0 {
		err = nil
		return chunks, nil
	}

	return
}

func (w *WorkerClusterService) FileToByteOp(info *pb.FileInfo) (chunks []byte, err error) {

	reader := bytes.NewReader(info.FileData)
	if global.OpC2ToOp != nil {
		global.OpC2ToOp.DelOpFile(context.TODO(), &pb.FileInfo{FileName: info.Path})
	}

	beginTime := time.Now()
	for {

		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}

		_, err := global.OpC2ToOp.CreateOpFile(context.TODO(), &pb.FileInfo{Path: info.Path, FileName: info.FileName, FileData: dataByte[:n]})
		if err != nil {
			break
		}
		chunks = append(chunks, dataByte[:n]...)
	}
	isEOF := strings.Compare(err.Error(), "EOF")
	if isEOF == 0 {
		err = nil
		fmt.Println(fmt.Printf("FileToByteOp read %s success, len=%v, 共耗时: %s", info.Path, len(chunks), time.Now().Sub(beginTime).String()))
		return chunks, nil
	}

	return
}

func (w *WorkerClusterService) FileToByteOpC2(info *pb.FileInfo, client pb.OpServiceClient) (chunks []byte, err error) {

	reader := bytes.NewReader(info.FileData)
	//删除OP文件
	client.DelOpFile(context.TODO(), &pb.FileInfo{FileName: path.Join(info.Path, info.FileName)})
	beginTime := time.Now()
	var er error
	for {

		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, er = reader.Read(dataByte)
		if er != nil || 0 == n {
			break
		}

		_, er = client.CreateOpFile(context.TODO(), &pb.FileInfo{Path: info.Path, FileName: info.FileName, FileData: dataByte[:n]})
		if er != nil {
			break
		}
		chunks = append(chunks, dataByte[:n]...)
	}
	isEOF := strings.Compare(er.Error(), "EOF")
	if isEOF == 0 {
		er = nil
		fmt.Println(fmt.Printf("FileToByteOpC2 read %s success, len=%v, 共耗时: %s", info.Path, len(chunks), time.Now().Sub(beginTime).String()))
		return chunks, nil
	}

	return
}

func (w *WorkerClusterService) FileToByteGateWay(info *pb.FileInfo) (chunks []byte, err error) {

	reader := bytes.NewReader(info.FileData)
	beginTime := time.Now()
	for {

		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}

		_, err := global.OpC2ToOp.AddGateWayFile(context.TODO(), &pb.FileInfo{Path: info.Path, FileName: info.FileName, FileData: dataByte[:n]})
		if err != nil {
			return chunks, err
		}
		chunks = append(chunks, dataByte[:n]...)
	}
	isEOF := strings.Compare(err.Error(), "EOF")
	if isEOF == 0 {
		err = nil
		fmt.Println(fmt.Printf("FileToByteGateWay read %s success, len=%v,共耗时: %s", info.Path, len(chunks), time.Now().Sub(beginTime).String()))
		return chunks, nil
	}

	return
}

// ExportWorkerTask
// @function: ExportWorkerTask
// @description: Export the worker task list
// @param: info LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) ExportWorkerTask(info request.WorkerOp) (m map[int]map[string]string, e error) {

	info.Page = 1
	info.PageSize = 1000
	result, _, err := WorkerClusterServiceApi.WorkerTaskDetail(info)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("暂无数据导出")
	}

	resMap := make(map[int]map[string]string)
	statusMap := make(map[int]string)
	statusMap[0] = "未执行"
	statusMap[1] = "完成"
	statusMap[2] = "失败"
	statusMap[3] = "进行中"
	statusMap[4] = "排队中"

	for i, v := range result {

		rowMap := make(map[string]string)
		rowMap["taskId"] = v.TaskId
		rowMap["taskAgent"] = fmt.Sprintf("%s/%s", v.TaskAgent, v.TaskAgentNo)
		rowMap["sectorSize"] = strconv.Itoa(v.SectorSize) + "G"
		rowMap["serverName"] = v.ServerName + "/" + v.Ip
		rowMap["taskStatus"] = statusMap[v.TaskStatus]
		rowMap["beginTime"] = utils.TimeToFormat(v.BeginTime, "")
		rowMap["timeLength"] = v.TimeLength
		rowMap["endTime"] = utils.TimeToFormat(v.EndTime, "")
		resMap[i] = rowMap
	}

	return resMap, nil
}

// RunC2Task
// @function: RunC2Task
// @description: Run C2 task
// @param: info LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) RunC2Task() {

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {

		select {
		case <-ticker.C:

			err := w.DoC2Task()
			if err != nil {
				log.Println("DoC2Task err:", err)
			}

		}
	}
}

// DoC2Task Run C2 task
func (w *WorkerClusterService) DoC2Task() error {

	var lotusWorkerCluster []lotus.LotusWorkerCluster
	sql := `SELECT cl.id,cl.op_id,cl.ip,cl.room_num,cl.room_name,cl.server_name,cl.task_num 
				FROM lotus_worker_cluster cl
				INNER join sys_host_records rec ON cl.op_id=rec.uuid AND cl.deleted_at IS NULL
				WHERE rec.gateway_id=? AND cl.deleted_at IS null`

	err := global.ZC_DB.Raw(sql, global.GateWayID.String()).Scan(&lotusWorkerCluster).Error
	if err != nil {
		return err
	}

	if len(lotusWorkerCluster) > utils.ZERO {

		delFile := make(map[string]string)
		tableName, err := w.GetWorkerRelationsTable("")
		if err != nil {
			return err
		}

		//Task failure retry
		var errRelations []lotus.LotusWorkerRelations
		err = global.ZC_DB.Table(tableName).Where("gate_id=? and task_status=2", global.GateWayID.String()).Find(&errRelations).Error
		if err != nil {
			return err
		}

		if len(errRelations) > 0 {

			sql := fmt.Sprintf("delete from %s where gate_id=? and task_status=2", tableName)
			err = global.ZC_DB.Exec(sql, global.GateWayID.String()).Error
			if err != nil {
				return err
			}

			for _, v := range errRelations {
				delFile[v.OpId] = fmt.Sprintf("%s-%d.json", v.Miner, v.Number)
			}
		}

		//Retry the task for 10 minutes
		var timeOUtRelations []lotus.LotusWorkerRelations
		err = global.ZC_DB.Table(tableName).Where("gate_id=? and task_status=3", global.GateWayID.String()).Find(&timeOUtRelations).Error
		if err != nil {
			return err
		}

		if len(timeOUtRelations) > 0 {

			for _, v := range timeOUtRelations {
				if time.Now().Sub(v.BeginTime).Minutes() > 10 {

					log.Println(fmt.Sprintf("%d,C2 task time out to redo", v.Number))
					sql := fmt.Sprintf("delete from %s where id=? and gate_id=? and task_status=3", tableName)
					err = global.ZC_DB.Exec(sql, v.ID, global.GateWayID.String()).Error
					if err != nil {
						return err
					} else {
						delFile[v.OpId] = fmt.Sprintf("%s-%d.json", v.Miner, v.Number)
					}
				}
			}
		}

		//Obtaining the task result expires 10 minutes later
		var resRelations []lotus.LotusWorkerRelations
		err = global.ZC_DB.Table(tableName).Where("gate_id=? and is_remove=0 and task_status=1", global.GateWayID.String()).Find(&resRelations).Error
		if err != nil {
			return err
		}

		if len(resRelations) > 0 {

			for _, v := range resRelations {
				if time.Now().Sub(v.BeginTime).Minutes() > 15 {

					log.Println(fmt.Sprintf("%d,C2 Get the result timeout redo", v.Number))
					sql := fmt.Sprintf("delete from %s where id=? and gate_id=? and is_remove=0 and task_status=1", tableName)
					err = global.ZC_DB.Exec(sql, v.ID, global.GateWayID.String()).Error
					if err != nil {
						return err
					} else {
						delFile[v.OpId] = fmt.Sprintf("%s-%d.json", v.Miner, v.Number)
					}
				}
			}
		}

		if len(delFile) > 0 {

			var delWg sync.WaitGroup
			delWg.Add(len(delFile))
			for k, v := range delFile {

				go func(opId, fileName string) {

					defer func() { delWg.Done() }()

					client, dis := global.OpClinets.GetOpClient(opId)
					if dis {
						log.Println(fmt.Sprintf("DoC2Task GetOpClient Connection failed:%s", opId))
						return
					}

					filePath := path.Join("/ipfs/data/c2task"+define.OpCsPathC1, fileName)
					log.Println("DoC2Task C2任务异常删除文件:", filePath)
					ct, cancel := context.WithTimeout(context.Background(), time.Second*15)
					defer func() {
						cancel()
					}()
					client.DelOpFile(ct, &pb.FileInfo{FileName: filePath})

				}(k, v)
			}
			delWg.Wait()
		}

		for _, v := range lotusWorkerCluster {

			client, dis := global.OpClinets.GetOpClient(v.OpId)
			if dis {
				continue
			}

			//Obtain the number of C2 tasks that can be executed
			opC2Client, err := client.GetOpC2Client(context.TODO(), &pb.OpC2Client{})
			if err != nil {
				continue
			}

			if len(opC2Client.OpInfo) > utils.ZERO {

				// Number of running tasks
				var runTaskCount int64
				err = global.ZC_DB.Table(tableName).Where("gate_id=? and op_id=? and task_status=3", global.GateWayID.String(), v.OpId).Count(&runTaskCount).Error
				if err != nil {
					return err
				}

				if int(runTaskCount) >= len(opC2Client.OpInfo) {
					continue
				}

				// Obtain the number of free C2 tasks
				c2IdIndex := make(map[int]int)
				taskCount := utils.ZERO
				for k1, v1 := range opC2Client.OpInfo {
					//log.Println("DoC2Task OpInfo:", v.Ip, v1.OpId, v1.GpuUse)
					if !v1.GpuUse {
						c2IdIndex[taskCount] = k1
						taskCount++
					}
				}

				waitTaskCount := len(opC2Client.OpInfo) - int(runTaskCount)
				if taskCount > waitTaskCount {
					taskCount = waitTaskCount
				}
				if taskCount <= utils.ZERO {
					continue
				}

				// C2 task to be completed
				var workerRelations []lotus.LotusWorkerRelations
				err = global.ZC_DB.Table(tableName).Where("gate_id=? and task_status=4", global.GateWayID.String()).Order("id").Find(&workerRelations).Error
				if err != nil {
					return err
				}

				if len(workerRelations) > utils.ZERO {

					log.Println(fmt.Sprintf("Host IP:%s,Number of free:%d", v.Ip, taskCount))
					// The number of sectors is allocated based on the number of tasks that can be executed
					var wg sync.WaitGroup
					for i := 0; i < taskCount; i++ {

						if i == len(workerRelations) {
							break
						} else {
							wg.Add(1)
						}

						go func(index int, lwc lotus.LotusWorkerCluster) {
							defer func() {
								wg.Done()
							}()

							wr := workerRelations[index]
							// Bind the c2-worker machine
							updateMap := make(map[string]interface{})
							updateMap["relation_id"] = lwc.ID
							updateMap["op_id"] = lwc.OpId
							updateMap["ip"] = lwc.Ip
							updateMap["server_name"] = lwc.ServerName
							updateMap["room_id"] = lwc.RoomNum
							updateMap["room_name"] = lwc.RoomName
							updateMap["begin_time"] = utils.GetNowStr()
							updateMap["task_status"] = utils.THREE

							// Distribute documents
							filePath := define.PathIpfsData + "/c2task" + define.OpCsPathC1
							fileName := fmt.Sprintf("%s-%d.json", wr.Miner, wr.Number)
							b, err := ioutil.ReadFile(path.Join(filePath, fileName))
							if err != nil || len(b) == 0 {
								sql := fmt.Sprintf("delete from %s where id=? and gate_id=?", tableName)
								err = global.ZC_DB.Exec(sql, v.ID, global.GateWayID.String()).Error
								if err != nil {
									return
								}
								return
							}
							f, _ := os.Stat(path.Join(filePath, fileName))
							if f.Size() < 51000000 || f.Size() > 55000000 {
								os.Remove(path.Join(filePath, fileName))
								return
							}

							opC2FilePath := "/ipfs/data/c2task" + define.OpCsPathC1
							_, err = WorkerClusterServiceApi.FileToByteOpC2(&pb.FileInfo{Path: opC2FilePath, FileName: fileName, FileData: b}, client)
							if err != nil {
								return
							}

							err = global.ZC_DB.Table(tableName).Where("id", wr.ID).Updates(updateMap).Error
							if err != nil {
								return
							}

							t2 := &pb.SectorID{
								Miner:  wr.Miner,
								Number: uint64(wr.Number),
							}
							t1 := &pb.SectorRef{
								Id: t2,
							}

							param := &pb.SealerParam{}
							if num, ok := c2IdIndex[index]; ok {
								param.OpC2Id = opC2Client.OpInfo[num].OpId
								param.Sector = t1
								param.Host = fmt.Sprintf("%s:%s", opC2Client.OpInfo[num].Ip, opC2Client.OpInfo[num].Port)
							} else {
								return
							}

							// Execute the C2 task
							ct, cancel := context.WithTimeout(context.Background(), time.Second*30)
							defer func() {
								cancel()
							}()
							_, err = client.Commit2TaskRun(ct, param)
							if err != nil {
								return
							}

						}(i, v)
					}
					wg.Wait()
				}
			}
		}
	}

	return nil
}

// ModifyC2TaskStatus
// @function: ModifyC2TaskStatus
// @description: Change the C2 task status
// @param: info LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) ModifyC2TaskStatus(info *pb.FileInfo) error {

	minerId := fmt.Sprintf("t0%d", info.Miner)
	if utils.IsNull(minerId) || info.Number < utils.ZERO {
		return errors.New("miner or number param is error")
	}

	tableName, err := w.GetWorkerRelationsTable("")
	if err != nil {
		return err
	}
	updateMap := make(map[string]interface{})
	if info.TaskStatus == utils.ONE {
		updateMap["task_status"] = info.TaskStatus
		updateMap["time_length"] = info.TimeLength
		updateMap["end_time"] = utils.GetNowStr()
	} else {
		updateMap["task_status"] = info.TaskStatus
	}
	return global.ZC_DB.Table(tableName).Where("miner=? and number=? and task_status=3", minerId, info.Number).Updates(updateMap).Error
}

// DelC2TaskInfo
// @function: DelC2TaskInfo
// @description: Example Delete C2 task data
// @param: info LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) DelC2TaskInfo(info request.C2TaskInfo) error {

	sql := ""
	tableName, err := w.GetWorkerRelationsTable("")
	if err != nil {
		return err
	}
	var param []interface{}
	if info.DelType == 1 {

		sql = fmt.Sprintf("update %s set is_remove=?,updated_at=now() where task_status=? and deleted_at is null", tableName)
		param = append(param, utils.ONE, utils.ONE)

		if !utils.IsNull(info.Miner) {
			sql += " and miner=?"
			param = append(param, info.Miner)
		}
		if info.Number > 0 {
			sql += " and number=?"
			param = append(param, info.Number)
		}

	} else {

		sql = fmt.Sprintf("update %s set deleted_at=now() where deleted_at is null and task_status=? and miner=? and number=?", tableName)
		param = append(param, utils.FOUR, info.Miner, info.Number)

	}

	return global.ZC_DB.Table(tableName).Exec(sql, param...).Error
}

// C2FileSynStatus
// @function: C2FileSynStatus
// @description: Synchronize C2 status
// @param: info LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) C2FileSynStatus(info *pb.C2SectorID) (string, error) {

	tableName, err := w.GetWorkerRelationsTable("")
	if err != nil {
		return "", err
	}

	if info.ResType == 1 {

		err = global.ZC_DB.Table(tableName).Where("miner=? and number=? and task_status=1", info.Miner, info.Number).Update("sync_status", utils.ONE).Error
		if err != nil {
			return "", err
		}

	} else {

		var lotusWorkerRelations lotus.LotusWorkerRelations
		err := global.ZC_DB.Table(tableName).Where("id", info.Id).Find(&lotusWorkerRelations).Error
		if err != nil {
			return "", err
		} else {
			return strconv.Itoa(lotusWorkerRelations.SyncStatus), err
		}
	}
	return "", nil
}

// GetC2workerInfo
// @function: GetC2workerInfo
// @description: Get the C2 worker
// @param: info lotus.LotusWorkerCluster
// @return: error
func (w *WorkerClusterService) GetC2workerInfo(opId string) (int, error) {

	if opId == "" {
		return 0, errors.New("opId is nil")
	}

	var total int64
	err := global.ZC_DB.Model(&lotus.LotusWorkerCluster{}).Where("op_id=?", opId).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return int(total), nil
}

// SynC2Result
// @function: SynC2Result
// @description: Get Miner information
// @param: minerId string
// @return: error
func (w *WorkerClusterService) SynC2Result(minerId string, filePath, filename string, fileData []byte) error {

	var minerInfo lotus.LotusMinerInfo
	err := global.ZC_DB.Model(&lotus.LotusMinerInfo{}).Where("actor=? and deploy_status=2 and is_manage=1", minerId).Find(&minerInfo).Error
	if err != nil {
		return err
	}

	if (minerInfo != lotus.LotusMinerInfo{}) {
		err = lotusrpc.FullApi.SynC2Result(minerInfo.Token, minerInfo.Ip, filePath, filename, fileData)
		if err != nil {
			return err
		}
	} else {
		return errors.New("SynC2Result miner is not online, please redeploy")
	}

	return nil
}

// RedoC2Task
// @function: RedoC2Task
// @description: Redo the C2 task
// @param: opId string
// @return: error
func (w *WorkerClusterService) RedoC2Task(opId string) error {

	if utils.IsNull(opId) {
		return fmt.Errorf("The opId cannot be empty")
	}

	tableName, err := w.GetWorkerRelationsTable("")
	if err != nil {
		return err
	}

	checkNumber := true
	number, err := strconv.Atoi(opId)
	if err != nil {
		checkNumber = false
	}

	var redoCount int64
	db := global.ZC_DB.Table(tableName)
	if !checkNumber {
		err = db.Where("op_id", opId).Count(&redoCount).Error
	} else {
		err = db.Where("number", number).Count(&redoCount).Error
	}
	if err != nil {
		return err
	}

	if redoCount > 0 {

		if !checkNumber {
			sql := fmt.Sprintf("delete from %s where op_id=? and task_status=3", tableName)
			err = global.ZC_DB.Exec(sql, opId).Error
		} else {
			sql := fmt.Sprintf("delete from %s where number=?", tableName)
			err = global.ZC_DB.Exec(sql, number).Error
		}
		if err != nil {
			return err
		}
	}

	return nil
}
