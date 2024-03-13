package slot

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"oplian/define"
	"oplian/global"
	"oplian/lotusrpc"
	"oplian/model/lotus"
	"oplian/model/slot"
	"oplian/model/slot/request"
	response2 "oplian/model/slot/response"
	"oplian/service/lotus/dispatch"
	"oplian/service/pb"
	"oplian/utils"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var WorkerCarServiceApi = new(WorkerCarService)

type WorkerCarService struct {
}

// GetWorkerCarTaskList peter
// @function: GetWorkerCarTaskList
// @description: 获取workerCar任务列表
// @param: info request.QueryWorkerCarReq
// @return: list interface{}, total int64, err error
func (w *WorkerCarService) GetWorkerCarTaskList(info request.QueryWorkerCarReq) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&slot.WorkerCarTask{})
	var workerCluster []slot.WorkerCarTask
	sql := "1=1"
	var param []interface{}
	if info.KeyWord != "" {
		info.KeyWord = "%" + info.KeyWord + "%"
		sql += " and (task_name like ? or id like ?)"
		param = append(param, info.KeyWord, info.KeyWord)
	}
	if info.TaskStatus != "" {
		sql += " and task_status = ?"
		param = append(param, info.TaskStatus)
	}
	if info.MinerId != "" {
		sql += " and miner_id=?"
		param = append(param, info.MinerId)
	}
	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("id desc").Find(&workerCluster).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return workerCluster, total, err
}

// ModifyTaskStatus peter
// @function: ModifyTaskStatus
// @description: 修改workerCar任务状态 0进行中,2暂停,3已终止
// @param: info request.ModifyWorkerCarReq
// @return: err error
func (w *WorkerCarService) ModifyTaskStatus(info request.ModifyWorkerCarReq) error {

	if !strings.Contains("0,1,2,3,4", info.TaskStatus) {
		return errors.New("任务状态异常:" + info.TaskStatus)
	}
	var main slot.WorkerCarTask
	err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", info.Id).Find(&main).Error
	if err != nil {
		return err
	}
	if (main == slot.WorkerCarTask{}) {
		return errors.New(fmt.Sprintf("ID:%s,找不到对应的数据", info.Id))
	}

	err = global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", info.Id).Update("task_status", info.TaskStatus).Error
	if err != nil {
		return err
	}

	if main.TaskType == define.CarTaskTypeAuto {
		err = global.ZC_DB.Model(&slot.WorkerCarTaskNo{}).Where("task_id=? and task_status not in(1,5)", info.Id).Update("task_status", info.TaskStatus).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// ModifyWorkerNum peter
// @function: ModifyWorkerNum
// @description: 修改workerCar worker任务数
// @param: info request.ModifyWorkerCarReq
// @return: err error
func (w *WorkerCarService) ModifyWorkerNum(info request.ModifyWorkerCarReq) error {

	if !utils.IsNumber(info.WorkerTaskNum) {
		return errors.New("worker任务数异常:" + info.WorkerTaskNum)
	}

	var main slot.WorkerCarTask
	err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", info.Id).Find(&main).Error
	if err != nil {
		return err
	}
	if (main == slot.WorkerCarTask{}) {
		return errors.New(fmt.Sprintf("ID:%s,找不到对应的数据", info.Id))
	}

	return global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", info.Id).Update("worker_task_num", info.WorkerTaskNum).Error
}

// GetWorkerCarTaskDetail peter
// @function: GetWorkerCarTaskDetail
// @description: workerCar任务详情
// @param: info request.QueryWorkerCarDetailReq
// @return: err error
func (w *WorkerCarService) GetWorkerCarTaskDetail(info request.QueryWorkerCarDetailReq) (map[string]interface{}, error) {

	dataMap := make(map[string]interface{})
	var main slot.WorkerCarTask
	err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", info.Id).Find(&main).Error
	if err != nil {
		return dataMap, err
	}
	if (main != slot.WorkerCarTask{}) {
		dataMap["main"] = main
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), main.MinerId)
	db := global.ZC_DB.Table(tableName)
	var detail []lotus.LotusSectorPiece
	var total int64
	err = db.Where("queue_id", info.Id).Count(&total).Error
	if err != nil {
		return dataMap, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("run_index").Find(&detail).Error
		if err != nil {
			return dataMap, err
		}
	}

	resDetail := make([]response2.WorkerCarTaskDetail, 0)
	for _, v := range detail {

		paramDetail := &response2.WorkerCarTaskDetail{

			TaskName:       main.TaskName,
			CreatedAt:      utils.TimeToFormat(v.CreatedAt, utils.YearMonthDayHMS),
			DealId:         v.DealUuid,
			DealExpireDate: utils.TimeToFormat(v.ExpirationTime, utils.YearMonthDayHMS),
			CarName:        path.Join(v.CarPath, fmt.Sprintf("%s.car", v.PieceCid)),
			SectorSize:     fmt.Sprintf("%d GiB", v.PieceSize/define.Ss1GiB),
			MinerId:        v.Actor,
			TaskNo:         v.RunIndex,
			TaskStatus:     v.JobStatus,
			SectorId:       strconv.Itoa(int(v.SectorId)),
		}

		if v.JobStatus == 1 {
			paramDetail.WaitTime = time.Now().Sub(v.CreatedAt).String()
		} else {
			paramDetail.WaitTime = ""
		}
		resDetail = append(resDetail, *paramDetail)
	}

	if len(detail) > 0 {
		detailData := &response2.WorkerCarDetail{
			Page:     info.Page,
			PageSize: info.PageSize,
			Total:    total,
			Detail:   resDetail,
		}

		dataMap["detail"] = detailData
	}

	return dataMap, nil
}

// AddWorkerCarTask peter
// @function: AddWorkerCarTask
// @description: 新增workerCar任务
// @param: info lotus.WorkerCarTask
// @return: err error
func (w *WorkerCarService) AddWorkerCarTask(info slot.WorkerCarTask) (uint, error) {
	err := global.ZC_DB.Create(&info).Error
	return info.ID, err
}

// AddWorkerCarTaskDetail peter
// @function: AddWorkerCarTaskDetail
// @description: 新增workerCar任务明细
// @param: info []lotus.WorkerCarTaskDetail
// @return: err error
func (w *WorkerCarService) AddWorkerCarTaskDetail(minerId string, info []lotus.LotusSectorPiece) error {

	//创建表
	err := dispatch.DispatchServiceApi.CreateSectorPieceTable(minerId)
	if err != nil {
		return err
	}
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)

	var sectorPiece lotus.LotusSectorPiece
	if err := global.ZC_DB.Table(tableName).Order("id desc").Limit(1).Find(&sectorPiece).Error; err != nil {
		return err
	}

	index := 1
	if (sectorPiece != lotus.LotusSectorPiece{}) {
		if sectorPiece.RunIndex > 0 {
			index = sectorPiece.RunIndex
		}
	}

	infoData := make([]lotus.LotusSectorPiece, 0)
	for _, v := range info {
		index++
		v.RunIndex = index
		infoData = append(infoData, v)
	}

	return global.ZC_DB.Table(tableName).CreateInBatches(&infoData, len(infoData)).Error

}

// CheckWorkerCarParam peter
// @function: CheckWorkerCarParam
// @description: 检查workerCar任务参数
// @param: c *gin.Context
// @return: lotus.WorkerCarTask, error
func (w *WorkerCarService) CheckWorkerCarParam(c *gin.Context) (slot.WorkerCarTask, error) {

	var info slot.WorkerCarTask
	taskName, _ := c.GetPostForm("taskName")
	taskType, _ := c.GetPostForm("taskType")
	sectorType, _ := c.GetPostForm("sectorType")
	minerId, _ := c.GetPostForm("minerId")
	quotaWallet, _ := c.GetPostForm("quotaWallet")
	dataSourcePath, _ := c.GetPostForm("dataSourcePath")
	orderRange, _ := c.GetPostForm("orderRange")
	workerTaskNum, _ := c.GetPostForm("workerTaskNum")
	originalOpId, _ := c.GetPostForm("originalOpId")
	originalDir, _ := c.GetPostForm("originalDir")
	sectorSize, _ := c.GetPostForm("sectorSize")
	validityDays, _ := c.GetPostForm("validityDays")

	if utils.IsNull(taskName) {
		return info, errors.New("任务名称不能为空")
	}
	if utils.IsNull(taskType) {
		return info, errors.New("任务类型不能为空")
	} else {
		if !strings.Contains("1,2", taskType) {
			return info, errors.New(fmt.Sprintf("任务类型异常:%s", taskType))
		}
	}
	if utils.IsNull(sectorType) {
		return info, errors.New("扇区类型不能为空")
	}
	if utils.IsNull(minerId) {
		return info, errors.New("节点不能为空")
	}
	if utils.IsNull(sectorSize) {
		return info, errors.New("扇区大小不能为空")
	} else {
		if !strings.Contains("32,64", sectorSize) {
			return info, errors.New(fmt.Sprintf("扇区大小异常:%s", sectorSize))
		}
	}

	if utils.IsNull(workerTaskNum) {
		return info, errors.New("worker任务数限制不能为空")
	} else {
		b, err := regexp.MatchString("^[0-9]+$", workerTaskNum)
		if err != nil {
			return info, err
		}
		if !b {
			return info, errors.New(fmt.Sprintf("worker任务数异常:%s,请重新输入", workerTaskNum))
		}
	}

	if taskType == "2" {
		if utils.IsNull(originalOpId) {
			return info, errors.New("源值主机opId不能为空")
		}
		if utils.IsNull(originalDir) {
			return info, errors.New("源值主机目录不能为空")
		}
	} else if taskType == "1" {

		if utils.IsNull(quotaWallet) {
			return info, errors.New("DC份额钱包不能为空")
		}

		if utils.IsNull(orderRange) {
			return info, errors.New("订单数量范围不能为空")
		} else {
			if !strings.Contains(orderRange, ",") {
				return info, errors.New(fmt.Sprintf("订单数量范围异常:%s", orderRange))
			}
			orderRange := strings.Split(orderRange, ",")
			rangA, _ := strconv.Atoi(orderRange[0])
			rangB, _ := strconv.Atoi(orderRange[1])
			if len(orderRange) < 2 || rangA >= rangB {
				return info, errors.New(fmt.Sprintf("订单数量范围异常:%s", orderRange))
			}
		}

		if utils.IsNull(dataSourcePath) {
			return info, errors.New("数据源地址不能为空")
		}

		if utils.IsNull(validityDays) {
			return info, errors.New("订单生命周期不能为空")
		}

		if !strings.Contains("530,180", validityDays) {
			return info, fmt.Errorf("订单生命周期异常:%s", validityDays)
		}

	}

	st, _ := strconv.Atoi(sectorType)
	wt, _ := strconv.Atoi(workerTaskNum)
	ty, _ := strconv.Atoi(taskType)
	vd, _ := strconv.Atoi(validityDays)
	info.TaskName = taskName
	info.TaskType = ty
	info.SectorType = st
	info.MinerId = minerId
	info.QuotaWallet = quotaWallet
	info.DataSourcePath = dataSourcePath
	info.OrderRange = orderRange
	info.WorkerTaskNum = wt
	info.OriginalOpId = originalOpId
	info.OriginalDir = originalDir
	info.SectorSize = sectorSize
	info.ValidityDays = vd

	return info, nil
}

// CheckImportFile peter
// @function: CheckImportFile
// @description: 检查导入订单数据
// @param: c *gin.Context
// @return: []lotus.WorkerCarTaskDetail, error
func (w *WorkerCarService) CheckImportFile(c *gin.Context, minerId string) ([]lotus.LotusSectorPiece, error) {

	var detail []lotus.LotusSectorPiece

	file, fileHead, err := c.Request.FormFile("dcInfo")
	if err != nil {
		return detail, err
	}
	defer file.Close()
	//解析订单数据
	buf := make([]byte, fileHead.Size)
	file.Read(buf)
	//检查格式
	str := utils.TrimBlankSpace(string(buf))
	deals := strings.Split(str, ";")
	if len(deals) == 0 {
		return detail, errors.New("文件导入数据为空")
	}

	//创建表
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)
	err = dispatch.DispatchServiceApi.CreateSectorPieceTable(minerId)
	if err != nil {
		return detail, err
	}
	var resData []lotus.LotusSectorPiece
	err = global.ZC_DB.Table(tableName).Where("actor=? and job_status not in(1,2)", minerId).Find(&resData).Error
	if err != nil {
		return detail, err
	}

	dealMap := make(map[string]struct{})
	for _, v := range resData {
		if v.DealUuid != "" {
			dealMap[v.DealUuid] = struct{}{}
		}
	}

	for k, v := range deals {

		if utils.IsNull(v) {
			continue
		}
		k += 1
		format := strings.Split(v, ",")
		if len(format) < 3 {
			return detail, errors.New(fmt.Sprintf("第%d行,数据格式有误:%s", k, v))
		}
		if _, err = uuid.FromString(format[0]); err != nil {
			return detail, errors.New(fmt.Sprintf("第%d行,订单号异常:%s", k, v))
		} else {
			if _, ok := dealMap[format[0]]; ok {
				return detail, errors.New(fmt.Sprintf("第%d行,订单号已存在:%s", k, v))
			}
		}
		if !strings.HasPrefix(format[1], "baga6ea4seaq") {
			return detail, errors.New(fmt.Sprintf("第%d行,carId异常:%s", k, v))
		}
		exp, err := strconv.ParseUint(format[2], 10, 64)
		if err != nil {
			return detail, errors.New(fmt.Sprintf("第%d行,订单高度异常:%s", k, v))
		}
		if exp+120 <= utils.BlockHeight() {
			return detail, errors.New(fmt.Sprintf("第%d行,订单已过期:%s", k, v))
		}
	}

	for k, v := range deals {

		if utils.IsNull(v) {
			continue
		}
		k++
		format := strings.Split(v, ",")
		t, _ := strconv.Atoi(format[2])
		data := &lotus.LotusSectorPiece{
			DealUuid:       format[0],
			ExpirationTime: utils.BlockHeightToTime(int64(t) + define.DealDdValidity),
			RunIndex:       k,
			PieceCid:       format[1],
			JobStatus:      define.CarTaskDetailCopy,
		}
		detail = append(detail, *data)

	}

	if len(detail) == 0 {
		return detail, errors.New("文件导入数据为空")
	}

	return detail, nil
}

// GetExecuteWorkerCarTask peter
// @function: GetExecuteWorkerCarTask
// @description: 获取可执行的workerCar任务
// @param: minerId
// @return: lotus.WorkerCarTask, error
func (w *WorkerCarService) GetExecuteWorkerCarTask(minerId string, taskType int) (slot.WorkerCarTask, error) {

	var wt slot.WorkerCarTask
	err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("miner_id=? and task_type=? and task_status=?", minerId, taskType, define.CarTaskStatusNotStarted).Order("id").Limit(1).Find(&wt).Error
	if err != nil {
		return wt, err
	}

	return wt, nil
}

// GetDoWorkerList peter
// @function: GetDoWorkerList
// @description: 获取可做任务的worker机
// @param: minerId
// @return: []lotus.WorkerCarTask, error
func (w *WorkerCarService) GetDoWorkerList(minerId string) ([]lotus.LotusWorkerInfo, error) {

	var wt []lotus.LotusWorkerInfo
	var resWt []lotus.LotusWorkerInfo
	sql := ` SELECT f.ip FROM lotus_miner_info w 
			INNER JOIN lotus_worker_info f ON w.id=f.miner_id AND f.deleted_at IS null
            INNER JOIN lotus_worker_config cg ON f.op_id=cg.op_id AND cg.on_off1=1 and cg.deleted_at IS null
			WHERE w.actor=? AND f.worker_type=0 and f.deploy_status=2 and f.run_status=1 AND w.deleted_at IS NULL `
	err := global.ZC_DB.Raw(sql, minerId).Scan(&wt).Error
	if err != nil {
		return nil, err
	}

	//已接任务worker
	var wct []slot.WorkerCarTaskNo
	err = global.ZC_DB.Model(&slot.WorkerCarTaskNo{}).Where("miner_id=? and task_status in(0,2)", minerId).Find(&wct).Error
	if err != nil {
		return nil, err
	}

	wctMap := make(map[string]struct{})
	for _, v := range wct {
		wctMap[v.WorkerIp] = struct{}{}
	}

	for _, v := range wt {
		if _, ok := wctMap[v.Ip]; !ok {
			resWt = append(resWt, v)
		}
	}

	return resWt, nil
}

// GetManualWorkerList peter
// @function: GetManualWorkerList
// @description: 获取手动可做任务的worker机
// @param: minerId
// @return: []lotus.WorkerCarTask, error
func (w *WorkerCarService) GetManualWorkerList(minerId string) ([]lotus.LotusWorkerInfo, error) {

	var wt []lotus.LotusWorkerInfo
	sql := ` SELECT f.ip FROM lotus_miner_info w 
			INNER JOIN lotus_worker_info f ON w.id=f.miner_id AND f.deleted_at IS null
            INNER JOIN lotus_worker_config cg ON f.op_id=cg.op_id AND cg.on_off1=1 and cg.deleted_at IS null
			WHERE w.actor=? AND f.worker_type=0 and f.deploy_status=2 and f.run_status=1 AND w.deleted_at IS NULL `
	err := global.ZC_DB.Raw(sql, minerId).Scan(&wt).Error
	if err != nil {
		return nil, err
	}

	return wt, nil
}

// AddCarWorkerTaskNo peter
// @function: AddCarWorkerTaskNo
// @description: 新增worker机任务编号
// @param: info *pb.CarWorkerTaskNoInfo
// @return: []lotus.WorkerCarTask, error
func (w *WorkerCarService) AddCarWorkerTaskNo(info slot.WorkerCarTaskNo) error {

	var total int64
	err := global.ZC_DB.Model(&slot.WorkerCarTaskNo{}).Where("miner_id=? and car_no=? and input_dir=?", info.MinerId, info.CarNo, info.InputDir).Count(&total).Error
	if err != nil {
		return err
	}

	if total == 0 {
		return global.ZC_DB.Create(&info).Error
	}

	return nil
}

// GetRunWorkerTask peter
// @function: GetRunWorkerTask
// @description: 获取在跑任务
// @param: info *pb.String
// @return: error
func (w *WorkerCarService) GetRunWorkerTask(info *pb.String) (slot.WorkerCarTaskNo, error) {

	var wt slot.WorkerCarTaskNo
	err := global.ZC_DB.Model(&slot.WorkerCarTaskNo{}).Where("worker_ip=? and task_status=?", info.Value, define.CarTaskNoStatusStarted).Order("id").Limit(1).Find(&wt).Error
	if err != nil {
		return wt, err
	}

	return wt, nil
}

// ModifyCarTaskNo peter
// @function: GetRunWorkerTask
// @description: 更新car任务编号
// @param: info *pb.String
// @return: error
func (w *WorkerCarService) ModifyCarTaskNo(info *pb.CarWorkerTaskNoInfo) error {

	if info.StartNo > info.EndNo {
		return errors.New(fmt.Sprintf("StartNo,EndNo参数异常:%d,%d", info.StartNo, info.EndNo))
	}
	updateMap := make(map[string]interface{})
	updateMap["start_no"] = info.StartNo
	updateMap["end_no"] = info.EndNo
	if info.TaskStatus > 0 {
		updateMap["task_status"] = info.TaskStatus
	}
	err := global.ZC_DB.Model(&slot.WorkerCarTaskNo{}).Where("worker_ip = ? and id = ?", info.WorkerIp, info.Id).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

// ModifyCarTaskStatus peter
// @function: ModifyCarTaskStatus
// @description: 更改任务状态
// @param: workerIp string, id, taskStatus int
// @return: error
func (w *WorkerCarService) ModifyCarTaskStatus(info slot.WorkerCarTask) error {

	if !strings.Contains("0,1,2,3", strconv.Itoa(int(info.TaskStatus))) {
		return errors.New(fmt.Sprintf("info.TaskStatus 状态异常:%d", info.TaskStatus))
	}
	updateMap := make(map[string]interface{})
	if info.TaskStatus > 0 {
		updateMap["task_status"] = info.TaskStatus
	}
	updateMap["current_no"] = info.CurrentNo
	err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id = ?", info.ID).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

// GetRunCarTaskDetail 获取在跑workerCar任务明细
func (w *WorkerCarService) GetRunCarTaskDetail(workerIp, minerId string, tasKType int) (response2.CarTaskDetailInfo, error) {

	var detail response2.CarTaskDetailInfo
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)
	err := dispatch.DispatchServiceApi.CreateSectorPieceTable(minerId)
	if err != nil {
		return detail, err
	}

	sql := ` SELECT p.id,p.piece_cid,p.piece_size,p.car_size,p.data_cid,p.job_status,
                    t.original_op_id,t.original_dir,t.quota_wallet,t.validity_days FROM ` + tableName + ` p
				INNER JOIN worker_car_task t ON p.queue_id=t.id AND t.deleted_at IS null `

	var param []interface{}
	param = append(param, minerId)
	condition := "WHERE p.actor=? AND p.deleted_at IS NULL "
	if tasKType == define.CarTaskTypeAuto {
		condition += " AND p.job_status=? "
		param = append(param, define.CarTaskTypeAuto)
	} else {
		condition += " AND p.job_status=? "
		param = append(param, define.CarTaskDetailCopy)
	}
	if workerIp != "" {
		condition += " and p.worker_ip=? "
		param = append(param, workerIp)
	}
	if tasKType > 0 {
		condition += " and t.task_type=? "
		param = append(param, tasKType)
	}

	condition += " ORDER BY id LIMIT 1 "
	sql += condition

	err = global.ZC_DB.Raw(sql, param...).Scan(&detail).Error
	if err != nil {
		return detail, err
	}

	return detail, err
}

// GetWaitCarTaskDetail peter
// @function: GetWaitCarTaskDetail
// @description: 获取待创建任务明细
// @param: info *pb.WorkerCarTaskDetail
// @return: error
func (w *WorkerCarService) GetWaitCarTaskDetail(workerIp, minerId string) (slot.WorkerCarTaskDetail, error) {

	var detail slot.WorkerCarTaskDetail
	err := global.ZC_DB.Model(&slot.WorkerCarTaskDetail{}).Where("worker_ip=? and miner_id = ? and task_status=?", workerIp, minerId, define.CarTaskDetailWait).Order("id").Limit(1).Find(&detail).Error
	if err != nil {
		return detail, err
	}

	return detail, err
}

// GetAllCarTaskDetail peter
// @function: GetAllCarTaskDetail
// @description: 获取所有待创建任务明细
// @param: info *pb.WorkerCarTaskDetail
// @return: error
func (w *WorkerCarService) GetAllCarTaskDetail(id int, minerId string) ([]lotus.LotusSectorPiece, error) {

	var detail []lotus.LotusSectorPiece
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)
	err := global.ZC_DB.Table(tableName).Where("queue_id=? and actor = ? and job_status=?", id, minerId, define.CarTaskDetailCopy).Find(&detail).Error
	if err != nil {
		return detail, err
	}

	return detail, err
}

// GetWorkerTaskCount peter
// @function: GetWorkerTaskCount
// @description: 获取worker机任务数量
// @param: info *pb.WorkerCarTaskDetail
// @return: error
func (w *WorkerCarService) GetWorkerTaskCount(minerId string) ([]response2.WorkerTask, error) {

	var detail []response2.WorkerTask
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)
	sql := fmt.Sprintf("SELECT s.worker_ip,COUNT(s.worker_ip) AS total FROM %s s WHERE s.actor=? and s.worker_ip <> '' AND s.job_status in (1,2,7) group by s.worker_ip", tableName)
	err := global.ZC_DB.Raw(sql, minerId).Scan(&detail).Error
	if err != nil {
		return detail, err
	}

	return detail, err
}

// ModifyCarTaskDetailInfo peter
// @function: ModifyCarTaskDetailInfo
// @description:  更改任务明细信息
// @param: info *pb.CarWorkerTaskDetailInfo
// @return: error
func (w *WorkerCarService) ModifyCarTaskDetailInfo(info *pb.CarWorkerTaskDetailInfo) error {

	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), info.MinerId)
	updateMap := make(map[string]interface{})

	if info.DealId != "" {
		updateMap["deal_uuid"] = info.DealId
	}
	if info.CarName != "" {
		updateMap["car_path"] = info.CarName
	}
	if info.PieceCid != "" {
		updateMap["piece_cid"] = info.PieceCid
	}
	if info.WorkerIp != "" {
		updateMap["worker_ip"] = info.WorkerIp
	}
	if info.TaskStatus > 0 {
		updateMap["job_status"] = info.TaskStatus
	}

	err := global.ZC_DB.Table(tableName).Where("id = ?", info.Id).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

// GetRand peter
// @function: GetRand
// @description: 获取随机数
// @param: info *pb.String
// @return: error
func (w *WorkerCarService) GetRand(info *pb.String) ([]slot.WorkerCarRand, error) {

	var wr []slot.WorkerCarRand
	err := global.ZC_DB.Model(&slot.WorkerCarRand{}).Find(&wr).Error
	if err != nil {
		return nil, err
	}

	return wr, nil
}

// CarFileExist peter
// @function: CarFileExist
// @description: 判断car文件是否存在
// @param: sectorId string
// @return: bool, error
func (w *WorkerCarService) CarFileExist(minerId, pieceCid string) (string, error) {

	carPath := ""
	if minerId != "" && pieceCid != "" {

		var resWt lotus.LotusSectorPiece
		tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)
		err := global.ZC_DB.Table(tableName).Where("actor=? and piece_cid=? and job_status IN(1,2)", minerId, pieceCid).Order("id desc").Limit(1).Find(&resWt).Error
		if err != nil {
			return carPath, err
		}

		if (resWt != lotus.LotusSectorPiece{}) {
			carPath = resWt.CarPath
		}
	}

	return carPath, nil
}

// GetCarUrlList peter
// @function: GetCarUrlList
// @description: 获取car文件路径
// @param: sectorId string
// @return: bool, error
func (w *WorkerCarService) GetCarUrlList(minerId string) ([]lotus.LotusSectorPiece, error) {

	var resWt []lotus.LotusSectorPiece
	if minerId != "" {
		//createdAt := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
		tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)
		err := global.ZC_DB.Table(tableName).Where("actor=? and job_status IN(1,2) and car_path <> '' and piece_cid <> '' ", minerId).Find(&resWt).Error
		if err != nil {
			return resWt, err
		}
	}

	time.Sleep(time.Second * 5)
	return resWt, nil
}

// GetDcQuotaWalletList peter
// @function: GetDcQuotaWalletList
// @description: 获取DC额度钱包
// @param: minerId string
// @return: []response.QuotaInfo, error
func (w *WorkerCarService) GetDcQuotaWalletList(minerId string) ([]response2.WorkerCarMiner, error) {

	if minerId == "" {
		return nil, errors.New("minerId 不能为空")
	}
	var configInfo response2.ConfigInfo
	var wt []response2.WorkerCarMiner
	sql := `SELECT b.dc_quota_wallet,l.token AS lotus_token,l.ip AS lotus_ip,w.actor,w.token AS miner_token,w.ip AS miner_ip
		   FROM lotus_miner_info w 
			LEFT JOIN lotus_boost_info b ON w.id=b.miner_id
			LEFT JOIN lotus_info l ON w.lotus_id=l.id
			WHERE w.actor=? AND w.deleted_at IS NULL AND b.deleted_at IS NULL `

	err := global.ZC_DB.Raw(sql, minerId).Scan(&configInfo).Error
	if err != nil {
		return nil, err
	}

	if (configInfo != response2.ConfigInfo{}) {

		maxSecotrId := "0"
		sectorSize := "32GiB"
		if !utils.IsNull(configInfo.MinerToken) && !utils.IsNull(configInfo.MinerIp) {
			num, err := lotusrpc.FullApi.MaxSectorNumber(configInfo.MinerToken, configInfo.MinerIp)
			if err == nil {
				maxSecotrId = strconv.Itoa(int(num))
			}

			if !utils.IsNull(configInfo.Actor) {
				mi, err := lotusrpc.FullApi.StateMinerInfo(configInfo.LotusToken, configInfo.LotusIp, configInfo.Actor)
				if err == nil {
					sectorSize = fmt.Sprintf("%dGiB", mi.SectorSize/define.Ss1GiB)
				}
			}
		}
		if configInfo.DcQuotaWallet != "" {
			walletAr := strings.Split(configInfo.DcQuotaWallet, ",")
			if len(walletAr) > 0 {
				for _, v := range walletAr {
					vAr := strings.Split(v, "|")
					switch len(vAr) {
					case 1:
						wt = append(wt, response2.WorkerCarMiner{Wallet: vAr[0], MaxSectorId: maxSecotrId, SectorSize: sectorSize})
						break
					case 2:
						wt = append(wt, response2.WorkerCarMiner{Wallet: vAr[0], Remark: vAr[1], MaxSectorId: maxSecotrId, SectorSize: sectorSize})
						break
					}
				}
			}
		}
	}

	return wt, nil
}

// DistributeWorkerTask peter
// @function: DistributeWorkerTask
// @description: worker任务分配
// @param: minerId string
// @return: error
func (w *WorkerCarService) DistributeWorkerTask() error {

	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:

			var minerList []lotus.LotusMinerInfo
			err := global.ZC_DB.Model(&lotus.LotusMinerInfo{}).Where("deploy_status =2 and run_status=1").Find(&minerList).Error
			if err != nil {
				return err
			}

			if len(minerList) > 0 {

				var wg sync.WaitGroup
				wg.Add(len(minerList))
				for _, v := range minerList {

					go func(minerInfo lotus.LotusMinerInfo) {
						defer func() {
							wg.Done()
						}()

						if err := w.WorkerCarTask(minerInfo.Actor); err != nil {
							log.Println("worker任务分配 WorkerCarTask err:", err)
						}

					}(v)
				}
				wg.Wait()
			}
		}
	}

	return nil
}

// WorkerCarTask peter
// @function: WorkerCarTask
// @description: worker任务分配
// @param: minerId string
// @return: error
func (w *WorkerCarService) WorkerCarTask(minerId string) error {

	resData, err := w.GetExecuteWorkerCarTask(minerId, 2)
	if err != nil {
		return err
	}

	if (resData == slot.WorkerCarTask{}) {
		return w.AutoWorkerCarTask(minerId)
	} else {
		return w.ManualWorkerCarTask(minerId)
	}

	return nil
}

// AutoWorkerCarTask peter
// @function: AutoWorkerCarTask
// @description: worker自动任务分配
// @param: minerId string
// @return: error
func (w *WorkerCarService) AutoWorkerCarTask(minerId string) error {

	//任务编号
	resData, err := w.GetExecuteWorkerCarTask(minerId, define.CarTaskTypeAuto)
	if err != nil {
		return err
	}

	if (resData != slot.WorkerCarTask{}) {

		//默认10个
		if resData.WorkerTaskNum == 0 {
			resData.WorkerTaskNum = 1
		}

		//可接任务worker
		resWorkerList, err := w.GetDoWorkerList(minerId)
		if err != nil {
			return err
		}

		if len(resWorkerList) > 0 {

			totalTask := 0
			rangAr := strings.Split(resData.OrderRange, ",")
			startNo, endNo := 0, 0
			if len(rangAr) != 2 {
				return errors.New(fmt.Sprintf("编号异常:%s", resData.OrderRange))
			}
			if len(rangAr) == 2 {

				//分配编号
				rangA, _ := strconv.Atoi(rangAr[0])
				rangB, _ := strconv.Atoi(rangAr[1])
				if resData.CurrentNo > 0 {
					rangA = resData.CurrentNo
				}
				if rangB <= rangA {

					//标记任务完成状态
					err = w.ModifyWorkerTaskStatus(minerId, resData, define.CarTaskTypeAuto)
					if err != nil {
						return err
					}
					return nil
				}

				rangeMap := make(map[int]int)
				if rangB > resData.WorkerTaskNum {

					rangLen := (rangB - resData.CurrentNo) / resData.WorkerTaskNum
					if rangLen*resData.WorkerTaskNum < (rangB - resData.CurrentNo) {
						rangLen++
					}

					for i := 0; i < rangLen; i++ {

						if i == len(resWorkerList) {
							break
						}
						totalTask += resData.WorkerTaskNum
						if startNo == 0 && endNo == 0 {
							st := 0
							if resData.CurrentNo > 0 {
								startNo = resData.CurrentNo
								st = resData.WorkerTaskNum + resData.CurrentNo
							} else {
								startNo = rangA
								st = startNo + resData.WorkerTaskNum
							}
							endNo = st
							if endNo >= rangB {
								endNo = rangB
							}

						} else {
							startNo = endNo
							endNo = startNo + resData.WorkerTaskNum
							if endNo >= rangB {
								endNo = rangB
							}
						}
						if startNo != endNo {
							rangeMap[startNo] = endNo
						}
					}

				} else {
					st := 0
					if resData.CurrentNo > 0 {
						startNo = resData.CurrentNo
						st = resData.WorkerTaskNum + resData.CurrentNo
					} else {
						startNo = rangA
						if rangB < resData.WorkerTaskNum {
							st = rangB
						} else {
							st = resData.WorkerTaskNum
						}

					}
					endNo = st
					if startNo != endNo {
						rangeMap[startNo] = endNo
					}
				}

				//排序
				var taskAr []int
				for k := range rangeMap {
					taskAr = append(taskAr, k)
				}

				sort.Ints(taskAr)

				//增加worker任务
				j := 0
				workerLen := len(resWorkerList)
				carFilesWorkerAr := make([]slot.WorkerCarTaskNo, 0)
				carNoMap := make(map[string]struct{})
				for _, v := range taskAr {

					if j < workerLen {

						res := resWorkerList[j]
						if v == rangeMap[v] {
							continue
						}

						if _, ok := carNoMap[fmt.Sprintf("%d,%d", v, rangeMap[v])]; !ok {

							paramData := &slot.WorkerCarTaskNo{
								TaskId:   int(resData.ID),
								MinerId:  minerId,
								WorkerIp: res.Ip,
								CarNo:    fmt.Sprintf("%d,%d", v, rangeMap[v]),
								InputDir: resData.DataSourcePath,
							}

							carFilesWorkerAr = append(carFilesWorkerAr, *paramData)
							carNoMap[paramData.CarNo] = struct{}{}

							//更新任务编号
							paramCarTask := &slot.WorkerCarTask{}
							paramCarTask.ID = resData.ID
							if rangeMap[v] >= rangB {
								paramCarTask.CurrentNo = rangeMap[v]
							} else {
								paramCarTask.CurrentNo = rangeMap[v]
							}

							err = w.ModifyCarTaskStatus(*paramCarTask)
							if err != nil {
								return err
							}
						}
						j++
					}
				}

				//写入car编号任务
				for _, v := range carFilesWorkerAr {
					err = w.AddCarWorkerTaskNo(v)
					if err != nil {
						return err
					}
				}
			}
		}

		//标记任务完成状态
		err = w.ModifyWorkerTaskStatus(minerId, resData, define.CarTaskTypeAuto)
		if err != nil {
			return err
		}
	}

	return nil
}

// ManualWorkerCarTask peter
// @function: ManualWorkerCarTask
// @description: worker手动任务分配
// @param: minerId string
// @return: error
func (w *WorkerCarService) ManualWorkerCarTask(minerId string) error {

	//暂停自动任务
	err := w.StopAutoCarTask(minerId, define.CarTaskTypeAuto)
	if err != nil {
		return err
	}

	//待分配手动任务
	resCarTask, err := w.GetExecuteWorkerCarTask(minerId, define.CarTaskTypeManual)
	if err != nil {
		return err
	}

	if (resCarTask != slot.WorkerCarTask{}) {

		//待跑手动订单
		resCarTaskDetail, err := w.GetAllCarTaskDetail(int(resCarTask.ID), minerId)
		if err != nil {
			return err
		}

		if len(resCarTaskDetail) > 0 {

			//worker已接任务数
			resWorkerTask, err := w.GetWorkerTaskCount(minerId)
			if err != nil {
				return err
			}

			workerTaskMap := make(map[string]int)
			for _, v := range resWorkerTask {
				if v.WorkerIp != "" {
					//worker实际任务数
					if v.Total >= resCarTask.WorkerTaskNum {
						v.Total = 0
					} else {
						v.Total = resCarTask.WorkerTaskNum - v.Total
					}
					workerTaskMap[v.WorkerIp] = v.Total
				}
			}

			//可接任务worker
			resWorker, err := w.GetManualWorkerList(minerId)
			if err != nil {
				return err
			}

			//分配任务
			carOrderAr := make([]*pb.CarWorkerTaskDetailInfo, 0)
			if len(resWorker) > 0 {

				for _, v := range resWorker {

					//补齐worker任务数
					if _, ok := workerTaskMap[v.Ip]; !ok {
						workerTaskMap[v.Ip] = resCarTask.WorkerTaskNum
					}

					//任务数满跳出
					if count, ok := workerTaskMap[v.Ip]; ok {
						if count == 0 {
							continue
						}
					}

					for _, v1 := range resCarTaskDetail {

						if v1.WorkerIp != "" {
							continue
						}
						v1.WorkerIp = v.Ip
						detailParam := &pb.CarWorkerTaskDetailInfo{
							Id:       uint64(v1.ID),
							MinerId:  minerId,
							WorkerIp: v1.WorkerIp,
						}
						carOrderAr = append(carOrderAr, detailParam)
						//减少任务数
						if count, ok := workerTaskMap[v.Ip]; ok {
							count--
							if count <= 0 {
								break
							} else {
								workerTaskMap[v.Ip] = count
							}
						}
					}
				}

				//订单分配指定worker机
				if len(carOrderAr) > 0 {

					for _, v := range carOrderAr {
						err = w.ModifyCarTaskDetailInfo(v)
						if err != nil {
							return err
						}
					}
				}
			}
		}

		//标记任务完成状态
		err = w.ModifyWorkerTaskStatus(minerId, resCarTask, define.CarTaskTypeManual)
		if err != nil {
			return err
		}
	} else {
		//开启自动任务
		err := w.StartAutoCarTask(minerId, define.CarTaskTypeAuto)
		if err != nil {
			return err
		}
	}

	return nil
}

// StopAutoCarTask peter
// @function: StopAutoCarTask
// @description: 停止手动任务
// @param: minerId string
// @return: error
func (w *WorkerCarService) StopAutoCarTask(minerId string, taskType int) error {

	var wct []slot.WorkerCarTask
	err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("miner_id=? and task_type=? and task_status=0", minerId, taskType).Scan(&wct).Error
	if err != nil {
		return err
	}

	for _, v := range wct {

		err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", v.ID).Update("task_status", 2).Error
		if err != nil {
			return err
		}

		//err = global.ZC_DB.Model(&slot.WorkerCarTaskNo{}).Where("task_id=? and task_status=0", v.ID).Update("task_status", 2).Error
		//if err != nil {
		//	return err
		//}
	}

	return nil
}

// StartAutoCarTask peter
// @function: StartAutoCarTask
// @description: 开始手动任务
// @param: minerId string
// @return: error
func (w *WorkerCarService) StartAutoCarTask(minerId string, taskType int) error {

	var wct slot.WorkerCarTask
	err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("miner_id=? and task_type=? and task_status=2", minerId, taskType).Order("id").Limit(1).Scan(&wct).Error
	if err != nil {
		return err
	}

	if (wct != slot.WorkerCarTask{}) {

		err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", wct.ID).Update("task_status", 0).Error
		if err != nil {
			return err
		}

		err = global.ZC_DB.Model(&slot.WorkerCarTaskNo{}).Where("task_id=? and task_status=2", wct.ID).Update("task_status", 0).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// GetBoostConfig peter
// @function: GetBoostConfig
// @description: 获取boost配置
// @param: minerId string
// @return: string, error
func (w *WorkerCarService) GetBoostConfig(minerId string) (*pb.String, error) {

	str := &pb.String{}
	if minerId == "" {
		return str, errors.New("minerId 不能为空")
	}

	var boostConfig response2.BoostConfig
	sql := ` SELECT b.lan_ip,b.lan_port,b.token FROM lotus_miner_info m
			INNER JOIN lotus_boost_info b ON m.id=b.miner_id AND b.deleted_at IS null
			 WHERE m.actor=? AND m.deleted_at IS null `

	err := global.ZC_DB.Raw(sql, minerId).Scan(&boostConfig).Error
	if err != nil {
		return str, err
	}

	if (boostConfig != response2.BoostConfig{}) {
		if boostConfig.LanIp == "" {
			return str, errors.New("IP不能为空")
		}
		if boostConfig.LanPort == "" {
			return str, errors.New("端口号不能为空")
		}
		if boostConfig.Token == "" {
			return str, errors.New("token 不能为空")
		}

		str.Value = fmt.Sprintf("%s:/ip4/%s/tcp/%s/http", boostConfig.Token, boostConfig.LanIp, boostConfig.LanPort)
	}

	return str, err
}

// ModifyWorkerTaskStatus peter
// @function: ModifyWorkerTaskStatus
// @description: 标记完成状态
// @param: minerId string
// @return: error
func (w *WorkerCarService) ModifyWorkerTaskStatus(minerId string, wct slot.WorkerCarTask, taskType int) error {

	err := dispatch.DispatchServiceApi.CreateSectorPieceTable(minerId)
	if err != nil {
		return err
	}

	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), minerId)
	checkStatus := false
	finishNum := 0

	if taskType == define.CarTaskTypeAuto {

		rangeAr := strings.Split(wct.OrderRange, ",")
		rangeA, _ := strconv.Atoi(rangeAr[0])
		rangeB, _ := strconv.Atoi(rangeAr[1])
		finishNum = wct.CurrentNo - rangeA
		if wct.CurrentNo == rangeB {
			checkStatus = true
		}

	} else {

		waitTotal := 0
		var sectorPiece []lotus.LotusSectorPiece
		err = global.ZC_DB.Table(tableName).Where("queue_id=?", wct.ID).Find(&sectorPiece).Error
		if err != nil {
			return err
		}

		for _, v := range sectorPiece {
			if v.WorkerIp != "" {
				finishNum++
			} else {
				waitTotal++
			}
		}
		if waitTotal == 0 {
			checkStatus = true
		}
	}

	if finishNum > 0 && wct.FinishNum != finishNum {

		err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", wct.ID).Update("finish_num", finishNum).Error
		if err != nil {
			return err
		}
	}

	if checkStatus {

		updateMap := make(map[string]interface{})
		updateMap["task_status"] = define.CarTaskStatusNoFinish
		updateMap["finish_time"] = time.Now()
		err := global.ZC_DB.Model(&slot.WorkerCarTask{}).Where("id", wct.ID).Updates(updateMap).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// AddCarRand peter
// @function: AddCarRand
// @description: 增加随机数
// @param:
// @return: error
func (w *WorkerCarService) AddCarRand() error {

	var total int64
	err := global.ZC_DB.Model(&slot.WorkerCarRand{}).Limit(1).Count(&total).Error
	if err != nil {
		return err
	}

	numMap := make(map[int]struct{})
	numParam := make([]slot.WorkerCarRand, 0)
	if total == 0 {

		for i := 0; i < 20000; i++ {

			num := rand.Int()
			if _, ok := numMap[num]; !ok {
				numParam = append(numParam, slot.WorkerCarRand{NumIndex: i, Number: num})
				numMap[num] = struct{}{}
			}
			if len(numParam) == 100 {

				err := global.ZC_DB.Model(&slot.WorkerCarRand{}).CreateInBatches(numParam, len(numParam)).Error
				if err != nil {
					return err
				}

				numParam = nil
			}
		}
	}

	return nil
}

// AddCarFiles peter
// @function: AddCarFiles
// @description: 增加Car文件
// @param:
// @return: error
func (w *WorkerCarService) AddCarFiles(info *pb.CarFiles) error {

	carFile := &slot.WorkerCarFiles{
		RelationId:  int(info.RelationId),
		FileName:    info.FileName,
		FileIndex:   int(info.FileIndex),
		FileStr:     info.FileStr,
		CarFileName: info.CarFileName,
		PieceCid:    info.PieceCid,
		PieceSize:   int(info.PieceSize),
		CarSize:     int(info.CarSize),
		DataCid:     info.DataCid,
		InputDir:    info.InputDir,
	}

	var total int64
	if err := global.ZC_DB.Model(&slot.WorkerCarFiles{}).Where("car_file_name", info.CarFileName).Count(&total).Error; err != nil {
		return err
	}

	if total > 0 {
		return nil
	}

	return global.ZC_DB.Create(&carFile).Error
}

func (w *WorkerCarService) GetCarFilesTable(minerId string) (string, error) {

	if utils.IsNull(minerId) {
		return "", fmt.Errorf("minerId is null")
	}

	table := fmt.Sprintf("%s_%s", slot.WorkerCarFiles{}.TableName(), minerId)
	var total int64
	err := global.ZC_DB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_NAME='%s'", table)).Count(&total).Error
	if err != nil {
		return "", err
	}
	if total == 0 {

		sql := `CREATE TABLE IF NOT EXISTS ` + table + ` SELECT * FROM ` + slot.WorkerCarFiles{}.TableName() + ";"
		if err := global.ZC_DB.Exec(sql).Error; err != nil {
			return "", err
		}
		//拷贝没有索引，手动增加
		sql = "ALTER TABLE " + table +
			" MODIFY COLUMN `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT FIRST," +
			"ADD PRIMARY KEY (`id`)," +
			"ADD INDEX `piece_cid`(`piece_cid`) USING BTREE," +
			"ADD INDEX `car_file_name`(`car_file_name`) USING BTREE;"
		return table, global.ZC_DB.Exec(sql).Error
	}

	return table, nil
}
