package system

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/system"
	"oplian/model/system/request"
	"oplian/model/system/response"
	"oplian/service/pb"
	"oplian/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

var WarnManageServiceApp = new(WarnManageService)

type WarnManageService struct {
}

// GetWarnList
// @function: GetWarnList
// @description: Alarm list
// @param: info request.WarnReq
// @return: list interface{}, total int64, err error`
func (ws *WarnManageService) GetWarnList(info request.WarnReq) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysWarnManage{})
	var warnManageList []system.SysWarnManage
	sql := "1=1"
	var param []interface{}
	if info.WarnKeyWord != "" {
		info.WarnKeyWord = "%" + info.WarnKeyWord + "%"
		sql += " and (warn_id like ? or warn_name like ?)"
		param = append(param, info.WarnKeyWord, info.WarnKeyWord)
	}
	if info.WarnType > 0 {
		sql += " and warn_type=? "
		param = append(param, info.WarnType)
	}
	if info.WarnStatus > 0 {
		sql += " and warn_status=?"
		param = append(param, info.WarnStatus)
	}
	if info.Ip != "" {
		info.Ip = "%" + info.Ip + "%"
		sql += " and ip like ?"
		param = append(param, info.Ip)
	}
	if info.AssetsKeyWord != "" {
		info.AssetsKeyWord = "%" + info.AssetsKeyWord + "%"
		sql += " and (assets_num like ? or sn like ?)"
		param = append(param, info.AssetsKeyWord, info.AssetsKeyWord)
	}
	if info.StrategiesId != "" {
		info.StrategiesId = "%" + info.StrategiesId + "%"
		sql += " and strategies_id like ?"
		param = append(param, info.StrategiesId)
	}
	if info.ComputerType > 0 {
		sql += " and computer_type =?"
		param = append(param, info.ComputerType)
	}
	if info.RoomId != "" {
		sql += " and computer_room_no = ?"
		param = append(param, info.RoomId)
	}

	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&warnManageList).Error
		if err != nil {
			return nil, 0, err
		}
	}

	//处理日期格式数据
	var warnManageListNew []system.SysWarnManage
	for _, v := range warnManageList {
		v.NotifyTime = time.Time{}
		warnManageListNew = append(warnManageListNew, v)
	}

	return warnManageList, total, err
}

// ModifyWarnStatus
// @function: ModifyWarnStatus
// @description: Change alarm status
// @param: info  request.WarnReq
// @return: bool, error
func (ws *WarnManageService) ModifyWarnStatus(info request.WarnReq) (bool, error) {

	updateMap := make(map[string]interface{})
	updateMap["warn_status"] = info.WarnStatus
	updateMap["finish_time"] = utils.GetNowStr()
	err := global.ZC_DB.Model(&system.SysWarnManage{}).Where("id", info.Id).Updates(updateMap).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetWarnTotal
// @function: GetWarnTotal
// @description: Alarm center statistics
// @param: info request.WarnReq
// @return: response.WarnAllRes, error
func (ws *WarnManageService) GetWarnTotal(info request.WarnReq) (data interface{}, e error) {

	var warnManageList []system.SysWarnManage
	var wl response.WarnAllRes
	err := global.ZC_DB.Model(&system.SysWarnManage{}).Scan(&warnManageList).Error
	if err != nil {
		return wl, err
	}

	// Alarm type 1 Indicator alarm 2 Log alarm 3 Associated alarm 4 Event alarm (policy alarm) 5 Inspection alarm 6 Service alarm
	// Alarm status 1 is cleared,2 is in the alarm state, and 3 is disabled
	if len(warnManageList) > 0 {

		for _, v := range warnManageList {

			warnType := strconv.Itoa(int(v.WarnType))
			status := strconv.Itoa(int(v.WarnStatus))
			if status == "2" {
				if strings.Contains("1,2,3,4", warnType) {
					wl.PolicyTotal++
				}
				if warnType == "5" {
					wl.PatrolTotal++
				}
				if warnType == "6" {
					wl.BusinessTotal++
				}
			}
			if strings.Contains("1,3", status) {
				wl.ProcessedTotal++
			}
		}

		wl.WarnTotal = wl.PolicyTotal + wl.PatrolTotal + wl.BusinessTotal + wl.ProcessedTotal
	}

	return wl, nil
}

// GetWarnTrend
// @function: GetWarnTrend
// @description: Status diagram of the alarm center
// @param: info request.WarnReq
// @return: response.WarnAllRes, error
func (ws *WarnManageService) GetWarnTrend(info request.WarnReq) (res []response.WarnTrendRes, e error) {

	if info.BeginTime == "" && info.EndTime == "" {
		info.BeginTime = time.Now().Add(-time.Hour * 1).Format("2006-01-02 15:04:05")
		info.EndTime = utils.GetNowStr()
	}

	var re []response.WarnTrendRes
	var warnManageList []system.SysWarnManage
	db := global.ZC_DB.Model(&system.SysWarnManage{})
	err := db.Where("warn_status=2 and created_at>? and created_at<=?", info.BeginTime, info.EndTime).Scan(&warnManageList).Error
	if err != nil {
		return nil, err
	}

	totalMap, _ := GetTimeMap(info)
	if len(warnManageList) > 0 {

		for k, _ := range totalMap {
			total := 0
			for _, v := range warnManageList {
				if v.CreatedAt.Before(utils.StrToTime(k)) {
					total++
					totalMap[k] = total
				}
			}
		}
	}

	var str []string
	for k, _ := range totalMap {
		str = append(str, k)
	}
	sort.Strings(str)

	for _, v := range str {
		data := &response.WarnTrendRes{
			TimeStr: v,
			Number:  totalMap[v],
		}
		re = append(re, *data)
	}

	return re, nil
}

// SaveStrategy
// @function: SaveStrategy
// @description: Save or modify the alarm policy
// @param: info request.WarnStrategies
// @return: error
func (ws *WarnManageService) SaveStrategy(info request.WarnStrategiesReq) error {

	var err error
	db := global.ZC_DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	if info.Ws.ZC_MODEL.ID == 0 {
		info.Ws.ZC_MODEL.CreatedAt = time.Now()
		info.Ws.ZC_MODEL.UpdatedAt = time.Now()
		info.Ws.StrategiesStatus = define.Start.Int()
		err = db.Create(&info.Ws).Error
		if err != nil {
			return err
		}

		for _, v := range info.Sr {
			v.RelationId = strconv.Itoa(int(info.Ws.ZC_MODEL.ID))
			v.RelationType = 3
			v.BeginTime = time.Now()
			v.ZC_MODEL.CreatedAt = time.Now()
			v.ZC_MODEL.UpdatedAt = time.Now()
			err = db.Create(&v).Error
			if err != nil {
				return err
			}
		}
	} else {
		err = db.Where("id", info.Ws.ZC_MODEL.ID).Updates(info.Ws).Error
		if err != nil {
			return err
		}

		err = db.Delete(&system.SysOpRelations{}, "relation_id=?", info.Ws.ZC_MODEL.ID).Error
		if err != nil {
			return err
		}

		for _, v := range info.Sr {
			if v.ID > 0 {
				sql := `update sys_op_relations set deleted_at=null,updated_at=now() where id=?`
				err = db.Exec(sql, v.ID).Error
			} else {
				v.RelationId = strconv.Itoa(int(info.Ws.ZC_MODEL.ID))
				v.RelationType = 3
				v.BeginTime = time.Now()
				v.ZC_MODEL.CreatedAt = time.Now()
				v.ZC_MODEL.UpdatedAt = time.Now()
				err = db.Create(&v).Error
			}

			if err != nil {
				return err
			}
		}
	}

	return err
}

// DelStrategy
// @function: DelStrategy
// @description: Deleting an Alarm policy
// @param: info request.WarnStrategies
// @return: error
func (ws *WarnManageService) DelStrategy(info request.StrategyReq) error {
	return global.ZC_DB.Delete(&system.SysWarnStrategies{}, "id", info.Id).Error
}

// GetStrategyDetail
// @function: GetStrategyDetail
// @description: Alarm Policy Details
// @param: info request.WarnReq
// @return: data interface{}, err error
func (ws *WarnManageService) GetStrategyDetail(info request.StrategyReq) (data interface{}, err error) {

	dataMap := make(map[string]interface{})
	var strategiesInfo system.SysWarnStrategies
	var ipInfo []response.IpInfo
	err = global.ZC_DB.Model(&system.SysWarnStrategies{}).Where("id", info.Id).Find(&strategiesInfo).Error
	if err != nil {
		return nil, err
	}

	err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("relation_type=3 and relation_id=?", info.Id).Find(&ipInfo).Error
	if err != nil {
		return nil, err
	}

	dataMap["strategies_info"] = strategiesInfo
	dataMap["ip_info"] = ipInfo

	return dataMap, nil
}

// ModifyStrategyStatus
// @function: ModifyStrategyStatus
// @description: Change the alert policy state
// @param: info request.WarnReq
// @return: data interface{}, err error
func (ws *WarnManageService) ModifyStrategyStatus(info request.StrategyReq) error {

	return global.ZC_DB.Model(&system.SysWarnStrategies{}).Where("id", info.Id).Update("strategies_status", info.StrategyStatus).Error

}

// GetStrategyList
// @function: GetStrategyList
// @description: Alarm policy list
// @param: info request.WarnReq
// @return: list interface{}, total int64, err error
func (ws *WarnManageService) GetStrategyList(info request.StrategyReq) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysWarnStrategies{})
	var strategiesList []response.SysWarnStrategies

	sql := "1=1"
	var param []interface{}
	if info.StrategyKeyWord != "" {
		info.StrategyKeyWord = "%" + info.StrategyKeyWord + "%"
		sql += " and (strategies_id like ? or strategies_name like ?)"
		param = append(param, info.StrategyKeyWord, info.StrategyKeyWord)
	}
	if info.StrategyType > 0 {
		sql += " and strategies_type=? "
		param = append(param, info.StrategyType)

	}
	if info.Ip != "" {
		info.Ip = "%" + info.Ip + "%"
		sql += " and ip_list like ? "
		param = append(param, info.Ip)

	}

	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&strategiesList).Error
		if err != nil {
			return nil, 0, err
		}
	}

	ItemVal := "1|cpu使用率,2|磁盘空间使用率,3|程序内存使用占比,4|程序内存使用占比,5|内存使用量,6|丢包率(外网及内网),"
	ItemVal += "7|最大延迟(外网及内网),8|系统重启时间,9|GPU显示存使用占比,10|GPU功耗,"
	ItemVal += "11|业务日志,12|主机日志,"
	ItemVal += "13|agent 心跳丢失,14|磁盘IO,15|磁盘写满,16|ping 不可达,17|异常事件,18|主机重启,19|主机宕机"
	ItemValAr := strings.Split(ItemVal, ",")
	itemMap := make(map[string]string)
	for _, v := range ItemValAr {
		vAr := strings.Split(v, "|")
		itemMap[vAr[0]] = vAr[1]
	}

	var resStrategies []response.SysWarnStrategies
	for _, v := range strategiesList {
		v.MonitorItem = itemMap[v.MonitorItem]
		if utils.IsNull(v.MonitorItem) {
			v.MonitorItem = ""
		}
		resStrategies = append(resStrategies, v)

	}

	return resStrategies, total, err
}

// StrategyWarnConfig
// @function: StrategyWarnConfig
// @description: Policy alarm configuration
// @param: info request.WarnReq
// @return: bool, error
func (ws *WarnManageService) StrategyWarnConfig() ([]response.WarnStrategiesRes, error) {

	warnStrategiesRes := make([]response.WarnStrategiesRes, 0)
	var warnStrategies []system.SysWarnStrategies
	err := global.ZC_DB.Model(&system.SysWarnStrategies{}).Where("strategies_status", 1).Scan(&warnStrategies).Error
	if err != nil {
		return nil, err
	}

	var ipInfo []response.IpInfo
	err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("relation_type", 3).Order("relation_id").Scan(&ipInfo).Error
	if err != nil {
		return nil, err
	}

	warnMap := make(map[string]int)
	for k, v := range warnStrategies {
		warnMap[v.StrategiesId] = k
		res := &response.WarnStrategiesRes{
			WarnStrategies: v,
		}
		warnStrategiesRes = append(warnStrategiesRes, *res)
	}

	for _, v := range ipInfo {
		if index, ok := warnMap[v.RelationId]; ok {

			ipData := &response.IpInfo{
				OpId:       v.OpId,
				GateWayId:  v.GateWayId,
				Ip:         v.Ip,
				ServerName: v.ServerName,
				RelationId: v.RelationId,
				RoomId:     v.RoomId,
			}
			warnStrategiesRes[index].IpInfo = append(warnStrategiesRes[index].IpInfo, *ipData)
		}
	}
	return warnStrategiesRes, nil
}

// StrategyProcessing
// @function: StrategyProcessing
// @description: Policy processing
// @param: info request.WarnReq
// @return: bool, error
func (ws *WarnManageService) StrategyProcessing(ctx context.Context, res []response.WarnStrategiesRes) error {

	var client pb.GateServiceClient
	for _, v := range res {

		if v.WarnStrategies.EffectivePeriod != "" {
			effAr := strings.Split(v.WarnStrategies.EffectivePeriod, "-")
			beginTime := time.Now().Format("2006-01-02") + " " + effAr[0] + ":00"
			endTime := time.Now().Format("2006-01-02") + " " + effAr[1] + ":00"
			if utils.StrToTime(beginTime).Before(time.Now()) && utils.StrToTime(endTime).Before(time.Now()) {
				continue
			}
		}

		for _, v1 := range v.IpInfo {
			if client == nil {
				client = global.GateWayClinets.GetGateWayClinet(v1.GateWayId)
			}
			if client == nil {
				global.ZC_LOG.Error("", zap.String("err", "GateWayClient is nil"))
				return errors.New("GateWayClient is nil")
			}
			_, err := client.StrategyProcess(context.TODO(), &pb.StrategyInfo{GateWayId: v1.GateWayId, OpId: v1.OpId, RoomId: v1.RoomId, StrategiesId: v.WarnStrategies.StrategiesId})
			if err != nil {
				global.ZC_LOG.Error("StrategyProcess", zap.String("err", err.Error()))
				return err
			}
		}
	}

	return nil
}

// StrategyProcessType
// @function: StrategyProcessType
// @description: Policy type processing
// @param: info request.WarnReq
// @return: bool, error
func (ws *WarnManageService) StrategyProcessType(opId, roomId, StrategiesId string) (bool, error) {

	//str := ""
	//resMsg := ""
	var sd system.SysWarnStrategies
	err := global.ZC_DB.Model(&system.SysWarnStrategies{}).Where("strategies_id", StrategiesId).Scan(&sd).Error
	if err != nil {
		return false, err
	}

	//策略类型 1指标 2日志 3关联 4事件
	switch sd.StrategiesType {
	case utils.ONE:
		IndexStrategyProcess(utils.ONE, opId, roomId, sd)
		IndexStrategyProcess(utils.TWO, opId, roomId, sd)
		break
	case utils.TWO:
		LogStrategyProcess(utils.ONE, opId, roomId, sd)
		LogStrategyProcess(utils.TWO, opId, roomId, sd)
		LogStrategyProcess(utils.THREE, opId, roomId, sd)
		break
	case utils.THREE:
		RelateStrategyProcess(opId, roomId, sd)
		break
	case utils.FOUR:
		EventStrategyProcess(utils.ONE, opId, sd)
		EventStrategyProcess(utils.TWO, opId, sd)
		break
	}

	return true, nil
}

// IndexStrategyProcess
// @function: IndexStrategyProcess
// @description: Indicator policy alarm
// @param: info request.WarnReq
// @return: bool, error
func IndexStrategyProcess(queryType int, opId, roomId string, sd system.SysWarnStrategies) (bool, error) {

	msgMap := make(map[string]string)
	msgMap["001"] = "cpu使用率"
	msgMap["002"] = "磁盘空间使用率"
	msgMap["003"] = "程序内存使用占比"
	msgMap["004"] = "程序内存使用占比"
	msgMap["005"] = "内存使用量"
	msgMap["006"] = "丢包率(外网及内网)"
	msgMap["007"] = "最大延迟(外网及内网)"
	msgMap["008"] = "系统重启时间"
	msgMap["009"] = "GPU显示存使用占比"
	msgMap["010"] = "GPU功耗"
	resData := make(map[string]string)
	scriptMap := make(map[string]string)
	//cpu使用率
	scriptMap["001"] = `top -b -n1 | fgrep "Cpu" | awk '{print 100-$8,"%"}'`
	//磁盘空间使用率
	//scriptMap["002"] = ""
	////程序内存使用占比
	//scriptMap["003"] = ""
	////主机重启次数
	//scriptMap["004"] = ""
	////内存使用量
	//scriptMap["005"] = ""
	////丢包率(外网及内网)
	//scriptMap["006"] = ""
	////最大延迟(外网及内网)
	//scriptMap["007"] = ""
	////系统重启时间
	//scriptMap["008"] = ""
	////GPU显示存使用占比
	//scriptMap["009"] = ""
	////GPU功耗
	//scriptMap["010"] = ""

	for k, v := range scriptMap {
		m := ResIndexType(queryType, utils.ONE, k, v, sd)
		k = k + "|" + strconv.Itoa(queryType)
		if _, ok := m[k]; ok {
			resData[k] = m[k]
		}
	}

	for k, v := range resData {

		kAr := strings.Split(k, "|")
		if kAr[1] == "1" {
			data := &system.SysWarnManage{
				WarnName:       msgMap[kAr[0]],
				WarnType:       1,
				ComputerId:     opId,
				ComputerRoomId: roomId,
				WarnInfo:       v,
			}
			err := WarnManageServiceApp.AddWarnInfo(*data)
			if err != nil {
				return false, err
			}
		} else {
			updateMap := make(map[string]interface{})
			updateMap["warn_status"] = 1
			updateMap["warn_info"] = v
			updateMap["finish_time"] = utils.GetNowStr()
			err := global.ZC_DB.Model(&system.SysWarnManage{}).Where("computer_id=? and warn_type=1", opId).Updates(updateMap).Error
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func ResIndexType(queryType int, logType int, scriptKey, scriptVal string, sd system.SysWarnStrategies) map[string]string {

	resMap := make(map[string]string)
	i, closeFor, timeNum, warnNum, recoverNum := 0, 0, 0, 0, 0
	//1 Trigger condition,2 recover condition,3 No data
	if sd.CycleType == 1 {
		if queryType == 1 {
			timeNum = sd.CycleLength * sd.TriggerNum / sd.TriggerCheck
			closeFor = sd.TriggerCheck
		} else if queryType == 2 {
			timeNum = sd.CycleLength
			closeFor = sd.RecoverNum
		} else if queryType == 3 && sd.LostEnable == 1 {
			timeNum = sd.CycleLength
			closeFor = sd.LostNum
		}
	} else {
		closeFor = 1
	}

	resMsg := ""
	key := scriptKey + "|" + strconv.Itoa(queryType)

	for {

		if i == closeFor {
			break
		}

		select {

		case <-time.After(time.Duration(timeNum) * time.Second):

			b, err := utils.ExecScript(scriptVal)
			if err != nil {
				log.Println("ExecScript err:", err)
				//continue
			}
			str := string(b)

			//测试
			if queryType == 1 {
				str = "70"
			} else if queryType == 2 {
				str = "63"
			}

			if logType == 1 {
				resMsg = str
				curNum, _ := strconv.Atoi(str)
				sysNum, _ := strconv.Atoi(sd.KeyWord)
				if curNum > sysNum && queryType == 1 {
					warnNum++
				} else if curNum < sysNum && queryType == 2 {
					recoverNum++
				}
			} else {

			}
		}
		i++
	}

	if logType == 1 {
		if warnNum == sd.TriggerCheck && queryType == 1 {
			resMap[key] = resMsg
		}
		if recoverNum == sd.RecoverNum && queryType == 2 {
			resMap[key] = resMsg
		}
	} else {

	}

	return resMap
}

// LogStrategyProcess
// @function: LogStrategyProcess
// @description: Log policy alarm
// @param: info request.WarnReq
// @return: bool, error
func LogStrategyProcess(queryType int, opId, roomId string, sd system.SysWarnStrategies) (bool, error) {

	rMap := make(map[int]string)
	rMap[1] = "触发"
	rMap[2] = "恢复"
	rMap[3] = "无数据"

	resMsg := ""
	rowMap := make(map[string]struct{})
	//001miner机,002worker机,003lotus机,004存储机,005其它
	opType := "001"
	opTypeMap := make(map[string]string)
	//opTypeMap["001"] = " /mnt/md0/ipfs/logs/worker-c2-1-10.0.16.64.log"
	i, closeFor, timeNum, warnNum, recoverNum := 0, 0, 0, 0, 0
	//1触发条件,2恢复条件,3无数据
	if sd.CycleType == 1 {
		if queryType == 1 {
			timeNum = sd.CycleLength * sd.TriggerNum / sd.TriggerCheck
			closeFor = sd.TriggerCheck
		} else if queryType == 2 {
			timeNum = sd.CycleLength
			closeFor = sd.RecoverNum
		} else if queryType == 3 && sd.LostEnable == 1 {
			timeNum = sd.CycleLength
			closeFor = sd.LostNum
		}
	} else {
		closeFor = 1
	}

	beginTime := time.Now()
	for {

		if i == closeFor {
			break
		}

		select {

		case <-time.After(time.Duration(timeNum) * time.Second):

			logPath := opTypeMap[opType]
			var df []system.SysDictionaryConfig
			err := global.ZC_DB.Model(&system.SysDictionaryConfig{}).Where("relation_id", opId).Scan(&df).Error
			if err != nil {
				return false, err
			}

			dfMap := make(map[string]int)
			for k, v := range df {
				str := opId + v.DictionaryKey
				dfMap[str] = k
			}

			id := 0
			dataCheck := false
			startRow, endRow := "1", "20"

			key := opId + opType
			if index, ok := dfMap[key]; ok {

				dataCheck = true
				id = int(df[index].ZC_MODEL.ID)
				startRow = df[index].DictionaryValue
				s, _ := strconv.Atoi(startRow)
				e, _ := strconv.Atoi(df[index].DefaultValue)
				endRow = strconv.Itoa(s + e)
			}

			script := `sed -n '` + startRow + `,` + endRow + `p' ` + logPath
			log.Println("script:" + rMap[queryType] + "," + script)
			str, err := utils.ExecScript(script)
			if err != nil {
				global.ZC_LOG.Error("err", zap.String("sed", err.Error()))
				return false, err
			}

			if string(str) != "" {
				resMsg = string(str)
				if strings.Contains(string(str), sd.KeyWord) && queryType == 1 {
					warnNum++
				} else if !strings.Contains(string(str), sd.KeyWord) && queryType == 2 {
					recoverNum++
				}
				if !dataCheck {
					rowMap[endRow] = struct{}{}
					t := system.SysDictionaryConfig{
						RelationId:      opId,
						DictionaryKey:   opType,
						DefaultValue:    "20",
						DictionaryValue: endRow,
					}
					err = global.ZC_DB.Create(&t).Error
				} else {
					rowMap[endRow] = struct{}{}
					dataMap := make(map[string]interface{})
					dataMap["dictionary_value"] = endRow
					err = global.ZC_DB.Model(&system.SysDictionaryConfig{}).Where("id", id).Updates(dataMap).Error
				}
				if err != nil {
					return false, err
				}
			}
		}
		i++
	}

	endTime := time.Now()
	log.Println("timeLength:", endTime.Sub(beginTime))

	resStr := ""
	if warnNum == sd.TriggerCheck && queryType == 1 {
		resStr = "010"
	} else if recoverNum == sd.RecoverNum && queryType == 2 {
		resStr = "011"
	} else if len(rowMap) == 1 && queryType == 3 {
		resStr = "012"
	}

	if resStr != "" {
		if strings.Contains("010,012", resStr) {
			resMsg = utils.SubStr(resMsg, 0, strings.Index(resMsg, sd.KeyWord))
			data := &system.SysWarnManage{
				WarnName:       "日志告警",
				WarnType:       2,
				ComputerId:     opId,
				ComputerRoomId: roomId,
				WarnInfo:       resMsg,
			}
			err := WarnManageServiceApp.AddWarnInfo(*data)
			if err != nil {
				return false, err
			}
		} else {
			updateMap := make(map[string]interface{})
			updateMap["warn_status"] = 1
			updateMap["warn_info"] = ""
			updateMap["finish_time"] = utils.GetNowStr()
			err := global.ZC_DB.Model(&system.SysWarnManage{}).Where("computer_id=? and warn_type=2", opId).Updates(updateMap).Error
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// EventStrategyProcess
// @function: EventStrategyProcess
// @description: Event policy alarm
// @param: info request.WarnReq
// @return: bool, error
func EventStrategyProcess(queryType int, opId string, sd system.SysWarnStrategies) (bool, error) {

	resData := make(map[string]string)
	scriptMap := make(map[string]string)
	scriptMap["011"] = `top -b -n1 | fgrep "Cpu" | awk '{print 100-$8,"%"}'`
	scriptMap["012"] = ""
	scriptMap["013"] = ""
	scriptMap["014"] = ""
	scriptMap["015"] = ""
	scriptMap["016"] = ""
	scriptMap["017"] = ""

	for k, v := range scriptMap {
		m := ResIndexType(queryType, utils.TWO, k, v, sd)
		k = k + "||" + strconv.Itoa(queryType)
		if _, ok := m[k]; ok {
			resData[k] = m[k]
		}
	}

	return true, nil
}

// RelateStrategyProcess
// @function: RelateStrategyProcess
// @description: Associated policy alarm
// @param: info request.WarnReq
// @return: bool, error
func RelateStrategyProcess(opId, roomId string, sd system.SysWarnStrategies) (bool, error) {

	if sd.KeyWord == "" {
		return false, errors.New("data keyWord is nil")
	}

	idStr := make([]string, 0)
	idAr := strings.Split(sd.KeyWord, ",")
	for _, v := range idAr {
		idStr = append(idStr, v)
	}

	timeNum := 0
	if sd.CycleType == 1 {
		timeNum = sd.CycleLength
	}

	i := 0
	for {

		if i == 1 {
			break
		}
		select {
		case <-time.After(time.Duration(timeNum) * time.Second):
			var warData []system.SysWarnManage
			err := global.ZC_DB.Model(&system.SysWarnManage{}).Where("computer_id=? and warn_status=2 and marked=0 and warn_type in(1,2,3,4) and strategies_id in(?) ", opId, idStr).Scan(&warData).Error
			if err != nil {
				return false, err
			}

			if len(warData) > 0 {
				ids := make([]int, 0)
				for _, v := range warData {
					ids = append(ids, int(v.ZC_MODEL.ID))
				}

				if len(ids) > 0 {
					err = global.ZC_DB.Model(&system.SysWarnManage{}).Where("id in(?)", ids).Update("marked", 1).Error
					if err != nil {
						return false, err
					}
				}

				data := &system.SysWarnManage{
					WarnName:       "关联日志",
					WarnType:       3,
					ComputerId:     opId,
					ComputerRoomId: roomId,
					WarnInfo:       "",
				}
				err := WarnManageServiceApp.AddWarnInfo(*data)
				if err != nil {
					return false, err
				}

			}
		}
		i++
	}

	return true, nil
}

// AddWarnInfo
// @function: AddWarnInfo
// @description: Add alarm
// @param: info system.SysWarnManage
// @return: error
func (ws *WarnManageService) AddWarnInfo(info system.SysWarnManage) error {

	beginDate := utils.TimeToFormat(time.Now(), utils.YearMonthDay)
	endDate := utils.TimeAddDay(utils.ONE)
	var sysWarnManage []system.SysWarnManage
	err := global.ZC_DB.Model(&system.SysWarnManage{}).Where("warn_status=2 and created_at>? and created_at<? and warn_name=?", beginDate, endDate, info.WarnName).Find(&sysWarnManage).Error
	if err != nil {
		return err
	}

	if len(sysWarnManage) == utils.ZERO {

		var resHost system.SysHostRecord
		err = global.ZC_DB.Model(&system.SysHostRecord{}).Where("uuid", info.ComputerId).Scan(&resHost).Error
		if err != nil {
			return err
		}
		var resRoom system.SysMachineRoomRecord
		err = global.ZC_DB.Model(&system.SysMachineRoomRecord{}).Where("room_id", resHost.RoomId).Scan(&resRoom).Error
		if err != nil {
			return err
		}

		info.WarnId = utils.GetUid(100000000)
		if (resHost != system.SysHostRecord{}) {
			info.Ip = resHost.IntranetIP
			info.AssetsNum = resHost.AssetNumber
			info.Sn = resHost.DeviceSN
			info.ComputerType = resHost.HostClassify
		}
		if (resRoom != system.SysMachineRoomRecord{}) {
			info.ComputerRoomName = resRoom.RoomName
			info.ComputerRoomNo = resRoom.RoomId
		}
		info.WarnTime = time.Now()

		return global.ZC_DB.Model(&system.SysWarnManage{}).Create(&info).Error
	}

	return nil
}

// StrategyType
// @function: StrategyType
// @description: Policy Type 1 Indicator 2 Log 3 Associated with 4 events
// @param: info request.WarnReq
// @return: bool, error
func (ws *WarnManageService) StrategyType(strategyType int64) (data []response.ItemPrj, err error) {

	ItemVal := ""
	switch strategyType {
	case utils.ONE:

		ItemVal = "001|cpu使用率,002|磁盘空间使用率,003|程序内存使用占比,004|程序内存使用占比,005|内存使用量,006|丢包率(外网及内网),"
		ItemVal += "007|最大延迟(外网及内网),008|系统重启时间,009|GPU显示存使用占比,010|GPU功耗"
		break
	case utils.TWO:

		ItemVal = "011|业务日志,012|主机日志"
		break
	case utils.THREE:

		ItemVal = ""
		break
	case utils.FOUR:

		ItemVal = "013|agent 心跳丢失,014|磁盘IO,015|磁盘写满,016|ping 不可达,017|异常事件,018|主机重启,019|主机宕机"
		break
	}

	var res []response.ItemPrj
	if ItemVal != "" {
		ItemValAr := strings.Split(ItemVal, ",")
		for _, v := range ItemValAr {
			vAr := strings.Split(v, "|")
			res = append(res, response.ItemPrj{Key: vAr[0], Value: vAr[1]})
		}
	}

	return res, nil
}

// GetStrategyId
// @function: GetStrategyId
// @description: Policy Id
// @param: info request.WarnReq
// @return: bool, error
func (ws *WarnManageService) GetStrategyId() interface{} {
	return utils.GetUid(100000000)
}

func GetTimeMap(info request.WarnReq) (map[string]int, string) {

	t1 := utils.StrToTime(info.BeginTime)
	t2 := utils.StrToTime(info.EndTime)
	minuteNum := int(t2.Sub(t1)) / 1000000000 / 60
	timeLen := 0
	timeType := ""
	if minuteNum <= 5 {
		timeType = ""
		timeLen = minuteNum * 60 / 10
	} else {
		timeType = "minute"
		timeLen = minuteNum / 10
	}

	t := ""
	timeMap := make(map[string]int)
	for i := 0; i <= 10; i++ {
		timeLength := timeLen * i
		if minuteNum <= 5 {
			t = utils.TimeToFormat(t1.Add(time.Duration(timeLength)*time.Second), "")
		} else {
			t = utils.TimeToFormat(t1.Add(time.Duration(timeLength)*time.Minute), utils.YearMonthDayHMS)
		}
		timeMap[t] = 0
	}

	return timeMap, timeType
}

// GetWarnNumTotal
// @function: GetWarnNumTotal
// @description: Obtain the quantity in the alarm
// @param: nil
// @return: int64, error
func (ws *WarnManageService) GetWarnNumTotal() (total int64, err error) {
	db := global.ZC_DB.Model(&system.SysWarnManage{})
	err = db.Where("warn_status = 2").Count(&total).Error
	if err != nil {
		return 0, err
	}
	return
}
