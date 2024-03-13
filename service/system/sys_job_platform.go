package system

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus"
	"oplian/model/system"
	"oplian/model/system/request"
	"oplian/model/system/response"
	"oplian/service/gateway"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

var JobPlatformServiceApp = new(JobPlatformService)

type JobPlatformService struct {
	IsStop bool
}

// ExecuteScript
// @function: ExecuteScript
// @description: Execute script
// @param: s request.ScriptInfo
// @return: bool, error
func (jp *JobPlatformService) ExecuteScript(s request.JobReq) error {

	go func() {

		if s.TimeLength < utils.ZERO || s.TimeLength > 86400 {
			global.ZC_LOG.Error("time_length data is error")
			return
		}

		scriptName := fmt.Sprintf("script_execute_%s.sh", utils.TimeToFormat(time.Now(), utils.YMDHMS))

		id, err := AddExecuteRecords(request.DistributeReq{GateWayId: s.GateWayId, TaskName: s.TaskName, TaskType: s.TaskType, SysOpRelations: s.SysOpRelations, UserName: s.UserName, ScriptName: scriptName})
		if err != nil {
			global.ZC_LOG.Error("AddExecuteRecords err:", zap.Error(err))
			return
		}
		s.ID = id

		gatewayTimeStart := time.Now()
		var executeRelation []system.SysOpRelations
		var executeRecord system.SysJobExecuteRecords
		m, _ := RecordsDetail(s.ID)
		if len(m) > 0 {

			if err := mapstructure.Decode(m["execute_records"], &executeRecord); err != nil {
				global.ZC_LOG.Error("AddExecuteRecords err:", zap.Error(err))
				return
			}
			if err := mapstructure.Decode(m["execute_relation"], &executeRelation); err != nil {
				global.ZC_LOG.Error("AddExecuteRecords err:", zap.Error(err))
				return
			}
		} else {
			global.ZC_LOG.Error("RecordsDetail data is nil")
			return
		}

		var sys sync.WaitGroup
		sys.Add(len(executeRelation))
		client := global.GateWayClinets.GetGateWayClinet(s.GateWayId)
		if client == nil {
			global.ZC_LOG.Error(fmt.Sprintf("GateWayClient Connection failed: %s", s.GateWayId))
			return
		}

		for i, v := range executeRelation {

			go func(index int, d system.SysOpRelations) {

				defer sys.Done()
				opTimeStart := time.Now()
				opStatus := define.TaskSuccess.Int()
				resMsg := ""
				re, err := client.ExecuteScript(context.TODO(), &pb.ScriptInfo{OpId: d.OpId, FileName: scriptName, Script: s.ScriptContent, TimeLength: int64(s.TimeLength)})
				if re != nil {
					resMsg = re.Value
					if strings.Contains(resMsg, "超时") {
						opStatus = define.TaskFailed.Int()
					}
				}
				if err != nil {
					resMsg = err.Error()
					opStatus = define.TaskFailed.Int()
				}

				updateMap := make(map[string]interface{})
				updateMap["status"] = opStatus
				updateMap["time_length"] = time.Now().Sub(opTimeStart).String()
				updateMap["res_msg"] = resMsg
				err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("id=? and status=?", d.ZC_MODEL.ID, define.TaskInProgress.Int()).Updates(updateMap).Error
				if err != nil {
					global.ZC_LOG.Error("SysOpRelations err:", zap.Error(err))
					return
				}

			}(i, v)
		}

		sys.Wait()

		var sysOpRelations []system.SysOpRelations
		err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("relation_id", s.ID).Find(&sysOpRelations).Error
		if err != nil {
			global.ZC_LOG.Error("SysOpRelations err:", zap.Error(err))
			return
		}

		allStatus := ""
		for _, v := range sysOpRelations {
			allStatus += strconv.Itoa(v.Status) + ","
		}

		// status 1 succeeded,2 failed,3 Executing,4 parts succeeded
		if !utils.IsNull(allStatus) {

			allStatus = utils.SubStr(allStatus, utils.ZERO, len(allStatus)-utils.ZERO)
			status := utils.ZERO
			if !strings.Contains(allStatus, strconv.Itoa(define.TaskFailed.Int())) {
				status = define.TaskSuccess.Int()
			} else if !strings.Contains(allStatus, strconv.Itoa(define.TaskSuccess.Int())) {
				status = define.TaskFailed.Int()
			} else if strings.Contains(allStatus, strconv.Itoa(define.TaskSuccess.Int())) && strings.Contains(allStatus, strconv.Itoa(define.TaskFailed.Int())) {
				status = define.TaskPartialSuccess.Int()
			}

			updateMap := make(map[string]interface{})
			updateMap["status"] = status
			updateMap["time_length"] = time.Now().Sub(gatewayTimeStart).String()
			err = global.ZC_DB.Model(&system.SysJobExecuteRecords{}).Where("id=? and status=?", s.ID, define.TaskInProgress.Int()).Updates(updateMap).Error
			if err != nil {
				global.ZC_LOG.Error("SysOpRelations err:", zap.Error(err))
				return
			}
		}

	}()

	return nil
}

// AddExecuteRecords
// @function: AddExecuteRecords
// @description: Add execution record
// @param: s request.ScriptInfo
// @return: string, error
func AddExecuteRecords(s request.DistributeReq) (id int, err error) {

	uId := utils.GetUid(100000000)
	jobId := 0

	tx := global.ZC_DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	t := &system.SysJobExecuteRecords{
		TaskName:    s.TaskName,
		TaskType:    s.TaskType,
		TaskNum:     uId,
		OpNumber:    len(s.SysOpRelations),
		OperateUser: s.UserName,
	}
	if !utils.IsNull(s.ScriptName) {
		t.ScriptName = s.ScriptName
	}
	err = tx.Create(&t).Error
	if err != nil {
		return id, err
	}

	jobId = int(t.ZC_MODEL.ID)
	for _, v := range s.SysOpRelations {

		v.RelationId = strconv.Itoa(jobId)
		v.RelationType = s.TaskType
		v.GateWayId = s.GateWayId

		err := tx.Create(&v).Error
		if err != nil {
			return id, err
		}
	}

	return jobId, err
}

// ExecuteRecordsDetail
// @function: ExecuteRecordsDetail
// @description: Execution record details
// @param: s request.ScriptInfo
// @return: string, error
func (jp *JobPlatformService) ExecuteRecordsDetail(id int) (data map[string]interface{}, err error) {
	return RecordsDetail(id)
}

// RecordsDetail Execution record details
func RecordsDetail(id int) (data map[string]interface{}, err error) {

	res := make(map[string]interface{})
	var executeRecord system.SysJobExecuteRecords
	err = global.ZC_DB.Model(&system.SysJobExecuteRecords{}).Where("id=?", id).Find(&executeRecord).Error
	if err != nil {
		return nil, err
	}

	if executeRecord.TaskNum == "" {
		return nil, errors.New("sys_execute_records data is nil")
	}

	var executeRelation []system.SysOpRelations
	err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("relation_id=? and relation_type=?", id, executeRecord.TaskType).Order("created_at desc").Find(&executeRelation).Error
	if err != nil {
		return nil, err
	}

	if executeRecord.TaskType == utils.TWO {
		var sysFileManageUpload []system.SysFileManageUpload
		err = global.ZC_DB.Model(&system.SysFileManageUpload{}).Where("relation_id=?", id).Order("created_at desc").Find(&sysFileManageUpload).Error
		if err != nil {
			return nil, err
		}
		res["execute_file"] = sysFileManageUpload
	}

	res["execute_records"] = executeRecord
	res["execute_relation"] = executeRelation

	return res, nil
}

// GetExecuteRecordsList
// @function: GetExecuteRecordsList
// @description: Execution record list
// @param: info request.ExecuteRecordsReq
// @return: list interface{}, total int64, err error
func (jp *JobPlatformService) GetExecuteRecordsList(info request.ExecuteRecordsReq) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysJobExecuteRecords{})
	var ExecuteRecordList []system.SysJobExecuteRecords
	sql := "1=1"
	var param []interface{}
	if info.TaskType != "" {
		sql += " and task_type=?"
		param = append(param, info.TaskType)
	}
	if info.Keyword != "" {
		info.Keyword = "%" + info.Keyword + "%"
		sql += " and (task_num like ? or task_name like ?)"
		param = append(param, info.Keyword, info.Keyword)
	}
	if info.Status != "" {
		sql += " and status=?"
		param = append(param, info.Status)
	}
	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&ExecuteRecordList).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return ExecuteRecordList, total, err
}

// ExecuteResult
// @function: ExecuteResult
// @description: Result of executing script
// @param: script string
// @return: string, error
func (jp *JobPlatformService) ExecuteResult(info *pb.ScriptInfo) (string, error) {

	err := utils.CreateFile(utils.FileInfo{Path: define.FileRootDir + "/script", FileName: info.FileName, FileData: []byte(info.Script)})
	if err != nil {
		return "", err
	}

	timeLen := utils.ZERO
	if info.TimeLength > utils.ZERO {
		timeLen = int(info.TimeLength)
	} else {
		timeLen = 100000
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLen)*time.Second)
	defer cancel()

	script := fmt.Sprintf(path.Join(define.FileRootDir+"/script", info.FileName))
	log.Println("script:", script)
	exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("chmod 777 %s", script)).CombinedOutput()
	cmd := exec.CommandContext(ctx, "bash", "-c", script)
	b, err := cmd.CombinedOutput()
	msg := ""
	if ctx.Err() == context.DeadlineExceeded {
		msg = fmt.Sprintf("Program timeout termination:%s", ctx.Err())
	} else {
		msg = string(b)
	}

	if err == nil {
		os.Remove(path.Join(define.FileRootDir+"/script", info.FileName))
	}

	return msg, err
}

// BatchUploadFile
// @function: BatchUploadFile
// @description: Batch upload files
// @param: r *http.Request, f request.Distribute
// @return: string, error
func (jp *JobPlatformService) BatchUploadFile(c *gin.Context) (request.File, error) {

	var fl request.File
	//上传文件
	m, err := utils.FileUpLoad(c, define.FileUploadDir)
	if err != nil {
		log.Println("utils.FileUpLoad err: ", err.Error())
		return fl, err
	}

	for k, v := range m {

		fs := &request.UpLoadFile{
			FileName: k,
			FileSize: v,
		}

		fl.UpLoadFile = append(fl.UpLoadFile, *fs)
	}

	return fl, nil
}

// FileDistribution
// @function: FileDistribution
// @description: File distribution
// @param: r *http.Request, f request.Distribute
// @return: string, error
func (jp *JobPlatformService) FileDistribution(info request.DistributeReq) error {

	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New(fmt.Sprintf(" gateWayClient Connection failed: %s", info.GateWayId))
	}

	for _, v := range info.OpList {
		if !utils.IsNull(v.OpId) {
			index := strings.LastIndex(v.FilePath, "/")
			filePath := utils.SubStr(v.FilePath, utils.ZERO, index)
			fileName := utils.SubStr(v.FilePath, index+utils.ONE, len(v.FilePath))
			_, err := client.CheckOpPath(context.TODO(), &pb.DirFileReq{OpId: v.OpId, Path: filePath, FileName: fileName})
			if err != nil {
				resMsg := fmt.Sprintf("主机Ip: %s 文件路径: %s 不存在", v.Ip, path.Join(filePath, fileName))
				return errors.New(resMsg)
			}
		}
	}

	go func() {

		now := utils.GetNowStr()

		id, err := AddExecuteRecords(info)
		if err != nil {
			global.ZC_LOG.Error("AddExecuteRecords err:", zap.Error(err))
			return
		}
		info.Id = id

		dataAr := make([]system.SysFileManageUpload, 0)
		if len(info.UpLoadFile) > utils.ZERO {
			for _, v := range info.UpLoadFile {
				if utils.IsNull(v.FileName) {
					continue
				}
				dataAr = append(dataAr, system.SysFileManageUpload{RelationId: strconv.Itoa(info.Id), FileType: define.FileLocalUpload.Int(), FilePath: v.FileName, FileName: v.FileName, FileSize: int(v.FileSize)})
			}
		}

		if len(info.OpList) > utils.ZERO {
			for _, v := range info.OpList {

				if !utils.IsNull(v.OpId) {
					v.FileType = define.FileHostUpload.Int()
				} else {
					v.FileType = define.FileDistributeUpload.Int()
				}
				v.RelationId = strconv.Itoa(info.Id)
				dataAr = append(dataAr, v)

			}
		}

		if len(dataAr) > utils.ZERO {
			err = global.ZC_DB.Model(&system.SysFileManageUpload{}).Create(&dataAr).Error
			if err != nil {
				global.ZC_LOG.Error("SysFileManageUpload err:", zap.Error(err))
				return
			}
		}

		var wg0 sync.WaitGroup
		wg0.Add(len(info.OpList))
		if len(info.OpList) > utils.ZERO {
			for _, v := range info.OpList {

				go func(sys system.SysFileManageUpload) {

					defer wg0.Done()
					if !utils.IsNull(sys.OpId) {
						index := strings.LastIndex(sys.FilePath, "/")
						fileName := utils.SubStr(sys.FilePath, index+utils.ONE, len(sys.FilePath))
						path := utils.SubStr(sys.FilePath, utils.ZERO, index)
						param := &pb.AddFileInfo{
							OpId:     sys.OpId,
							FromPath: path,
							FileName: fileName,
							ToPath:   define.FileGateWayDir,
						}
						_, err := client.OpFileToGateWay(context.TODO(), param)
						if err != nil {
							global.ZC_LOG.Error("OpFileToGateWay:", zap.Error(err))
							return
						} else {
							info.UpLoadFile = append(info.UpLoadFile, request.UpLoadFile{FileName: fileName})
						}
					}

					//文件管理文件分发过来的
					if utils.IsNull(sys.OpId) && !utils.IsNull(sys.FilePath) {
						info.UpLoadFile = append(info.UpLoadFile, request.UpLoadFile{FileName: sys.FilePath, FileType: utils.ONE})
					}

				}(v)
			}
		}
		wg0.Wait()

		fileMap := make(map[string]interface{})
		for _, v := range info.UpLoadFile {
			if utils.IsNull(v.FileName) {
				continue
			}
			fileMap[v.FileName] = struct{}{}
		}
		fileInfoList, err := ioutil.ReadDir(define.FileUploadDir)
		if err != nil {
			global.ZC_LOG.Error("ReadDir err:", zap.Error(err))
			return
		}

		var wg1 sync.WaitGroup
		for i := range fileInfoList {

			fileName := fileInfoList[i].Name()
			if _, ok := fileMap[fileName]; !ok {
				continue
			} else {
				wg1.Add(1)
			}

			go func(name string) {

				defer wg1.Done()
				sy := &request.SynFileReq{
					FromPath: path.Join(define.FileUploadDir, name),
					ToPath:   define.FileGateWayDir,
					FileName: name,
				}

				_, err = GateWayFileToByte(sy, client)
				if err != nil {
					global.ZC_LOG.Error("GateWayFileToByte", zap.String("err", err.Error()))
				}

			}(fileName)
		}

		wg1.Wait()

		for k, _ := range fileMap {
			os.Remove(path.Join(define.FileUploadDir, k))
		}

		fileArr := make([]*pb.FileInfo, 0)
		for _, v := range info.UpLoadFile {
			if utils.IsNull(v.FileName) {
				continue
			}
			fileArr = append(fileArr, &pb.FileInfo{FileName: v.FileName, FileSize: v.FileSize, FileType: int64(v.FileType)})
		}

		param := &pb.FileSynOp{
			Id:         uint64(info.Id),
			TimeLength: int64(info.TimeLength),
			Enable:     int64(info.Enable),
			LimitSpeed: int64(info.LimitSpeed),
			FilePath:   info.FilePath,
			SendType:   int64(info.SendType),
			FileInfo:   fileArr,
			TimeOut:    now,
		}
		_, err = client.FileSynOpHost(context.TODO(), param)
		if err != nil {
			global.ZC_LOG.Error("FileSynOpHost err:", zap.Error(err))
			return
		}
	}()

	return nil
}

// ForcedToStop
// @function: ForcedToStop
// @description: Forced termination
// @param: r *http.Request, f request.Distribute
// @return: string, error
func (jp *JobPlatformService) ForcedToStop(info request.DistributeReq) error {

	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New("gateway client Connection failed:" + info.GateWayId)
	}

	client.SetJobPlatformStop(context.TODO(), &pb.JobPlatform{IsStop: true})

	return nil
}

// ExecuteScriptStop
// @function: ExecuteScriptStop
// @description: Script execution termination
// @param: r *http.Request, f request.Distribute
// @return: string, error
func (jp *JobPlatformService) ExecuteScriptStop(info request.JobReq) error {

	var sysJobExecuteRecords []system.SysJobExecuteRecords
	err := global.ZC_DB.Model(&system.SysJobExecuteRecords{}).Where("id=? and status=?", info.ID, define.TaskInProgress.Int()).Find(&sysJobExecuteRecords).Error
	if err != nil {
		return err
	}
	if len(sysJobExecuteRecords) == utils.ZERO {
		return errors.New("执行记录数据不存在")
	}

	info.ScriptContent = sysJobExecuteRecords[0].ScriptName
	var sysOpRelations []system.SysOpRelations
	err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("relation_id=? and status=?", info.ID, define.TaskInProgress.Int()).Find(&sysOpRelations).Error
	if err != nil {
		return err
	}

	if len(sysOpRelations) == utils.ZERO {
		return errors.New("无可终止任务主机")
	}

	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New("gateway client Connection failed:" + info.GateWayId)
	}

	var wg sync.WaitGroup
	wg.Add(len(sysOpRelations))
	for _, v := range sysOpRelations {

		go func(sy system.SysOpRelations) {

			defer wg.Done()

			resMsg := ""
			_, err := client.ScriptStop(context.TODO(), &pb.ScriptInfo{OpId: v.OpId, Script: info.ScriptContent})
			if err != nil {
				resMsg = err.Error()
				global.ZC_LOG.Error("ScriptStop err:", zap.Error(err))
			}

			//更新状态
			resMsg += "程序强行终止 "
			updateMap := make(map[string]interface{})
			updateMap["res_msg"] = resMsg
			updateMap["status"] = define.TaskFailed.Int()
			err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("id", sy.ID).Updates(updateMap).Error
			if err != nil {
				global.ZC_LOG.Error("modify SysOpRelations err:", zap.Error(err))
				return
			}

		}(v)
	}

	wg.Wait()

	err = global.ZC_DB.Model(&system.SysJobExecuteRecords{}).Where("id", info.ID).Update("status", define.TaskFailed.Int()).Error
	if err != nil {
		return err
	}

	return nil
}

// ExecuteFileSynOp
// @function: ExecuteFileSynOp
// @description: File synchronization to OP
// @param: r *http.Request, f request.Distribute
// @return: string, error
func (jp *JobPlatformService) ExecuteFileSynOp(info *pb.FileSynOp) error {

	var executeRelation []system.SysOpRelations
	var executeRecord system.SysJobExecuteRecords
	reMap, _ := RecordsDetail(int(info.Id))
	if len(reMap) > 0 {

		if err := mapstructure.Decode(reMap["execute_records"], &executeRecord); err != nil {
			return err
		}
		if err := mapstructure.Decode(reMap["execute_relation"], &executeRelation); err != nil {
			return err
		}
	} else {
		return errors.New("RecordsDetail data is nil")
	}

	m := make(map[string]int64)
	gateWayFileName := ""
	for _, v := range info.FileInfo {
		if v.FileType == utils.ONE {
			gateWayFileName = v.FileName
		}
		m[v.FileName] = v.FileType
	}

	fileInfoList, err := ioutil.ReadDir(define.FileGateWayDir)
	if err != nil {
		log.Println("ReadDir", zap.String("err", err.Error()))
		return err
	}

	if len(fileInfoList) == utils.ZERO {
		return errors.New("数据文件已丢失")
	}

	dirMsg := ""
	allStatus := ""
	var wg sync.WaitGroup
	wg.Add(len(executeRelation))
	start := time.Now()
	for j, v := range executeRelation {

		go func(index int, filePath string, wr system.SysOpRelations) {

			defer wg.Done()

			checkError := false
			resMsg := ""
			singleStartTime := time.Now()
			if JobPlatformServiceApp.IsStop && !checkError {
				checkError = true
				resMsg = fmt.Sprintf("主机op手动强制终止,Ip: %s", v.Ip)
			}

			client, dis := global.OpClinets.GetOpClient(wr.OpId)
			if dis && !checkError {
				checkError = true
				resMsg = fmt.Sprintf("主机op未启动,Ip: %s", v.Ip)
			}

			//Transmission mode 1 mandatory,0 rigorous
			if info.SendType == int64(define.StrictMode.Int()) && !checkError {
				_, err := client.CheckOpPath(context.TODO(), &pb.DirFileReq{OpId: v.OpId, Path: filePath})
				if err != nil {
					checkError = true
					resMsg = fmt.Sprintf("主机Ip: %s 文件路径: %s 不存在", v.Ip, filePath)
				}
			}

			//Check whether the file management file exists
			if !checkError && !utils.IsNull(gateWayFileName) {

				checkFile := false
				for i := range fileInfoList {
					if gateWayFileName == fileInfoList[i].Name() {
						checkFile = true
						break
					}
				}

				if !checkFile {
					checkError = true
					resMsg = fmt.Sprintf("文件管理,%s不存在", gateWayFileName)
				}
			}

			if !checkError {
				for i := range fileInfoList {

					if !utils.IsFile(filePath + "/" + fileInfoList[i].Name()) {
						continue
					}
					if _, ok := m[fileInfoList[i].Name()]; !ok && len(m) > utils.ZERO {
						continue
					}

					sy := &request.SynFileReq{
						FromPath: define.FileGateWayDir + "/" + fileInfoList[i].Name(),
						ToPath:   filePath,
						FileName: fileInfoList[i].Name(),
						OpId:     wr.OpId,
					}
					if info.Enable == utils.ONE {
						sy.LimitSpeed = int(info.LimitSpeed)
					}
					if info.TimeLength > utils.ZERO {
						sy.TimeOut = info.TimeOut
						sy.TimeLength = info.TimeLength
					}

					_, err = DistributeFileToByte(sy, client)
					if err != nil {
						resMsg = err.Error()
						log.Println("FileToByte", zap.String("err", err.Error()))
						break
					}
				}
			}

			updateMap := make(map[string]interface{})
			updateMap["status"] = define.TaskSuccess.Int()
			if !utils.IsNull(resMsg) {
				allStatus += strconv.Itoa(define.TaskFailed.Int()) + ","
				updateMap["status"] = define.TaskFailed.Int()
				updateMap["res_msg"] = resMsg
			} else {
				allStatus += strconv.Itoa(define.TaskSuccess.Int()) + ","
				updateMap["res_msg"] = "文件分发成功"
			}
			updateMap["time_length"] = time.Now().Sub(singleStartTime).String()
			err = global.ZC_DB.Model(system.SysOpRelations{}).Where("id=? and status=?", wr.ZC_MODEL.ID, define.TaskInProgress.Int()).Updates(updateMap).Error
			if err != nil {
				log.Println("SysOpRelations data update", zap.String("err", err.Error()))
				return
			}

		}(j, info.FilePath, v)
	}

	wg.Wait()

	//Update the execution status. 1 Succeeded. 2 Failed. 3 During the execution, four parts succeeded
	end := time.Now()
	updateMap := make(map[string]interface{})
	if !utils.IsNull(allStatus) {
		if !strings.Contains(allStatus, strconv.Itoa(define.TaskFailed.Int())) {
			updateMap["status"] = define.TaskSuccess.Int()
		} else if strings.Contains(allStatus, strconv.Itoa(define.TaskSuccess.Int())) && strings.Contains(allStatus, strconv.Itoa(define.TaskFailed.Int())) {
			updateMap["status"] = define.TaskPartialSuccess.Int()
		} else if !strings.Contains(allStatus, strconv.Itoa(define.TaskSuccess.Int())) {
			updateMap["status"] = define.TaskFailed.Int()
		}
	} else if !utils.IsNull(dirMsg) {
		updateMap["status"] = define.TaskFailed.Int()
	}

	updateMap["time_length"] = end.Sub(start).String()
	err = global.ZC_DB.Model(system.SysJobExecuteRecords{}).Where("id", executeRecord.ZC_MODEL.ID).Updates(updateMap).Error
	if err != nil {
		log.Println("SysJobExecuteRecords data update", zap.String("err", err.Error()))
		return err
	}

	JobPlatformServiceApp.IsStop = false

	if updateMap["status"] == define.TaskSuccess.Int() {
		for k, _ := range m {
			//排除gateway文件
			if k == gateWayFileName {
				continue
			}
			os.Remove(path.Join(define.FileGateWayDir, k))
		}
	}

	return nil
}

// OpFileToGateWay
// @function: OpFileToGateWay
// @description: Synchronize files from the host to the gateway
// @param: f request.DistributeReq
// @return: string, error
func (jp *JobPlatformService) OpFileToGateWay(info request.OpFileSync) (request.File, error) {

	var fl request.File
	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return fl, errors.New("gateway client Connection failed:" + info.GateWayId)
	}

	for _, v := range info.OpListReq {

		l := strings.LastIndex(v.Path, "/")
		if l == -1 {
			l = strings.LastIndex(v.Path, "\\")
		}
		filePath := utils.SubStr(v.Path, 0, l)
		fileName := utils.SubStr(v.Path, l+1, len(v.Path))

		_, err := client.OpFileToGateWay(context.TODO(), &pb.AddFileInfo{OpId: v.OpId, FileName: fileName, FromPath: filePath, ToPath: define.FileGateWayDir})
		if err != nil {
			log.Println("OpFileToGateWay", zap.String("err", err.Error()))
		}

		fs := &request.UpLoadFile{
			FileName: fileName,
		}

		fl.UpLoadFile = append(fl.UpLoadFile, *fs)

	}

	return fl, nil
}

// CreateFile
// @function: CreateFile
// @description: Generate file
// @param: script string
// @return: string, error
func (jp *JobPlatformService) CreateFile(p *pb.FileInfo) (bool, error) {

	ok, _ := utils.PathExists(p.Path)
	if !ok {
		err := os.MkdirAll(p.Path, os.ModePerm)
		if err != nil {
			return false, err
		}
	}

	fileUrl := p.Path + "/" + p.FileName
	f, err := os.OpenFile(fileUrl, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return false, err
	}

	_, err = f.Write(p.FileData)
	if err != nil {
		return false, err
	}
	defer f.Close()

	return true, nil
}

// PointFileToByte
// @function: PointFileToByte
// @description: Point-to-point replication files are converted to byte
// @param: sy *request.SynFileReq, opClient pb.OpServiceClient
// @return: chunks []byte, err error
func PointFileToByte(sy *request.SynFileReq, client pb.OpServiceClient) (chunks []byte, err error) {

	fileName := sy.ToPath + "/" + sy.FileName
	client.DelOpFile(context.TODO(), &pb.FileInfo{FileName: fileName})
	f, err := os.Open(sy.FromPath)
	defer f.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(f)
	beginTime := time.Now()
	for {

		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}

		f := &pb.FileInfo{
			FileData: dataByte[:n],
			FileName: sy.FileName,
			Path:     sy.ToPath,
		}

		_, err = client.FileDistribution(context.TODO(), f)
		if err != nil {
			log.Println("PointFileToByte", zap.String("err", err.Error()))
			break
		}
		chunks = append(chunks, dataByte[:n]...)
	}
	isEOF := strings.Compare(err.Error(), "EOF")
	if isEOF == 0 {

		err = nil
		fmt.Println(fmt.Printf("read %s success, len=%v, Total time spent: %s", sy.FromPath, len(chunks), time.Now().Sub(beginTime).String()))
		return chunks, nil
	}

	return
}

// GateWayFileToByte
// @function: GateWayFileToByte
// @description: Synchronize to the gateway file and convert it to byte
// @param: sy *request.SynFileReq
// @return: chunks []byte, err error
func GateWayFileToByte(sy *request.SynFileReq, client pb.GateServiceClient) (chunks []byte, err error) {

	fileName := sy.ToPath + "/" + sy.FileName
	client.DelGateWayFile(context.TODO(), &pb.FileInfo{FileName: fileName})

	f, err := os.Open(sy.FromPath)
	defer f.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(f)

	if client == nil {
		return nil, errors.New("GateWayFileToByte OpToGatewayClient Connection failed")
	}

	beginTime := time.Now()
	for {
		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}

		f := &pb.FileInfo{
			FileData: dataByte[:n],
			FileName: sy.FileName,
			Path:     sy.ToPath,
		}

		_, err = client.AddGateWayFile(context.TODO(), f)
		if err != nil {
			log.Println("GateWayFileToByte", zap.String("err", err.Error()))
			break
		}
		chunks = append(chunks, dataByte[:n]...)
	}
	isEOF := strings.Compare(err.Error(), "EOF")
	if isEOF == 0 {

		err = nil
		fmt.Println(fmt.Printf("read %s success, len=%v ,Total time spent %s", sy.FromPath, len(chunks), time.Now().Sub(beginTime).String()))
		return chunks, nil
	}

	return
}

// DistributeFileToByte
// @function: DistributeFileToByte
// @description: Convert the distributed file to byte
// @param: sy *request.SynFileReq, client pb.GateServiceClient
// @return: chunks []byte, err error
func DistributeFileToByte(sy *request.SynFileReq, client pb.OpServiceClient) (chunks []byte, err error) {

	client.DelOpFile(context.TODO(), &pb.FileInfo{OpId: sy.OpId, FileName: fmt.Sprintf(path.Join(sy.ToPath, sy.FileName))})
	f, err := os.Open(sy.FromPath)
	defer f.Close()
	if err != nil {
		return
	}
	reader := bufio.NewReader(f)

	allSize := sy.LimitSpeed * define.Ss1MB
	dzSize := utils.ZERO
	beginTime := time.Now()
	errMsg := ""

	for {

		if JobPlatformServiceApp.IsStop {
			errMsg = "Forced program termination"
			break
		}

		if sy.TimeLength > utils.ZERO && utils.StrToTime(sy.TimeOut).Add(time.Second*time.Duration(sy.TimeLength)).Before(time.Now()) {
			errMsg = "Program timeout termination"
			break
		}

		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}

		f := &pb.FileInfo{
			Path:     sy.ToPath,
			FileData: dataByte[:n],
			FileName: sy.FileName,
		}

		_, err = client.FileDistribution(context.TODO(), f)
		if err != nil {
			global.ZC_LOG.Error("FileDistribution", zap.String("err", err.Error()))
			break
		}

		//限速
		dzSize += define.FileByte.Int()
		if dzSize > allSize && allSize > utils.ZERO {
			dzSize = 0
			time.Sleep(time.Second)
		}

		chunks = append(chunks, dataByte[:n]...)
	}

	if utils.IsNull(errMsg) {
		isEOF := strings.Compare(err.Error(), "EOF")
		if isEOF == 0 {
			err = nil
			log.Println(fmt.Printf("read %s success: len=%v ,Total time spent: %s", sy.FromPath, len(chunks), time.Now().Sub(beginTime).String()))
			return chunks, nil
		} else {
			return chunks, err
		}
	} else {
		return chunks, errors.New(errMsg)
	}

	return
}

// SaveFileHost
// @function: SaveFileHost
// @description: Setup file host
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) SaveFileHost(info request.AddFileReq) error {

	var res system.SysDictionaryConfig
	err := global.ZC_DB.Model(&system.SysDictionaryConfig{}).Where("relation_id=? and default_value=? and dictionary_value=?", "000001", info.GateWayId, info.OpId).Scan(&res).Error
	if err != nil {
		return err
	}

	if (res == system.SysDictionaryConfig{}) {

		res.RelationId = "000001"
		res.DictionaryKey = "host"
		res.DefaultValue = info.GateWayId
		res.DictionaryValue = info.OpId
		err = global.ZC_DB.Model(&system.SysDictionaryConfig{}).Create(&res).Error
	} else {

		updateMap := make(map[string]interface{})
		updateMap["default_value"] = info.GateWayId
		updateMap["dictionary_value"] = info.OpId
		err = global.ZC_DB.Model(&system.SysDictionaryConfig{}).Where("relation_id", "000001").Updates(updateMap).Error
	}

	if err != nil {
		return err
	}

	return nil
}

// GetFileManageList
// @function: GetFileManageList
// @description: File management list
// @param: info request.ExecuteRecordsReq
// @return: list interface{}, total int64, err error
func (jp *JobPlatformService) GetFileManageList(info request.FileManageReq) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysFileManage{})
	var sysFileManageList []system.SysFileManage
	sql := "1=1"
	var param []interface{}
	if info.FileName != "" {
		info.FileName = "%" + info.FileName + "%"
		sql += " and file_name like ?"
		param = append(param, info.FileName)
	}
	if info.FileStatus != 0 {
		sql += " and file_status=?"
		param = append(param, info.FileStatus)
	}
	if info.FileType != 0 {
		sql += " and file_type=?"
		param = append(param, info.FileType)
	}
	if info.KeyWord != "" {
		sql += " and room_id=? "
		param = append(param, info.KeyWord)
	}
	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&sysFileManageList).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return sysFileManageList, total, err
}

// GetFileName
// @function: GetFileName
// @description: Get file name
// @param: info request.ExecuteRecordsReq
// @return: list interface{}, total int64, err error
func (jp *JobPlatformService) GetFileName(info *pb.FileNameInfo) (fileName string, err error) {

	log.Println(fmt.Sprintf("GetFileName GateWayId: %s,FileType: %d", info.GateWayId, info.FileType))
	if info.FileType <= utils.ZERO {
		return fileName, errors.New("fileType param is error")
	}

	var sysFileManage system.SysFileManage
	err = global.ZC_DB.Model(&system.SysFileManage{}).Where("gate_way_id=? and file_type=? and file_status=?", info.GateWayId, info.FileType, define.FileFinish.Int()).Order("updated_at desc").Limit(utils.ONE).Find(&sysFileManage).Error
	if err != nil {
		return fileName, err
	}

	if (sysFileManage == system.SysFileManage{}) {
		return fileName, errors.New("no documentation could be found")
	}

	return sysFileManage.FileName, nil
}

// DelFileManage
// @function: DelFileManage
// @description: Delete file management
// @param: info request.ExecuteRecordsReq
// @return: error
func (jp *JobPlatformService) DelFileManage(info request.FileManageReq) error {

	var sysFileManage system.SysFileManage
	err := global.ZC_DB.Model(&system.SysFileManage{}).Where("id", info.Id).Find(&sysFileManage).Error
	if err != nil {
		return err
	}

	if (sysFileManage == system.SysFileManage{}) {
		return errors.New(fmt.Sprintf("File data does not exist,id:%d", info.Id))
	}

	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New("DelFileManage GateWayClient Connection failed" + info.GateWayId)
	}

	t := &pb.FileInfo{
		FileName: define.FileGateWayDir + "/" + sysFileManage.FileName,
	}

	_, err = client.DelGateWayFile(context.TODO(), t)
	if err != nil {
		global.ZC_LOG.Error("DelFile", zap.String("err", err.Error()))
	}

	err = global.ZC_DB.Delete(&system.SysFileManage{}, "id", info.Id).Error
	if err != nil {
		return err
	}

	return nil
}

// ModifyFileManage
// @function: ModifyFileManage
// @description: Update file management
// @param: info request.ExecuteRecordsReq
// @return: error
func (jp *JobPlatformService) ModifyFileManage(info request.FileManageReq) error {

	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New("ModifyFileManage GateWayClient Connection failed" + info.GateWayId)
	}

	_, err := client.ModifyOnlineFile(context.TODO(), &pb.AddFileInfo{Id: strconv.Itoa(info.Id), FileUrl: info.FileUrl, FileType: int64(info.FileType)})
	if err != nil {
		return err
	}

	return nil
}

// OnlineDownload
// @function: OnlineDownload
// @description: Download online
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) OnlineDownload(info *pb.AddFileInfo) error {

	if info.FileType != define.GenerallyFile.Int64() && info.FileType != define.ProveFile.Int64() && info.FileType != define.SnapshotFile.Int64() {
		return errors.New("file_type is error")
	}

	if utils.IsNull(info.FileUrl) {
		return errors.New("file_url is nil")
	}

	if info.FileType == define.ProveFile.Int64() {
		var fileCount int64
		err := global.ZC_DB.Model(&system.SysFileManage{}).Where("file_type=? and gate_way_id=?", info.FileType, info.GateWayId).Count(&fileCount).Error
		if err != nil {
			return err
		}
		if fileCount > utils.ZERO {
			return errors.New("only one certificate file can be downloaded from each equipment room")
		}
	}

	var hostRecord system.SysHostRecord
	err := global.ZC_DB.Model(&system.SysHostRecord{}).Where("uuid", info.OpId).Scan(&hostRecord).Error
	if err != nil {
		return err
	}

	res, err := new(gateway.DownloadService).DowloadFile(info.FileUrl, define.FileGateWayDir, "")
	if err != nil {
		log.Println("DownloadFile err:", err.Error())
	}

	r := &system.SysFileManage{
		RoomId:     info.RoomId,
		RoomName:   info.RoomName,
		GateWayId:  info.GateWayId,
		FileType:   int(info.FileType),
		FileUrl:    info.FileUrl,
		FileStatus: define.FileGetting.Int(),
	}
	if (hostRecord != system.SysHostRecord{}) {
		r.OpId = info.OpId
		r.ServerName = hostRecord.HostName
		r.Ip = hostRecord.IntranetIP
	}
	if res != nil {
		beginIndex := strings.LastIndex(res.FilePath, `/`)
		if beginIndex == -1 {
			beginIndex = strings.LastIndex(res.FilePath, `\`)
		}
		fileName := utils.SubStr(res.FilePath, beginIndex+1, len(res.FilePath))
		r.FileName = fileName
		r.FileSize = int(res.Total)
	} else {
		r.FileName = "未知文件"
		r.FileSize = utils.ZERO
		r.FileStatus = define.FileError.Int()
	}

	err = global.ZC_DB.Model(&system.SysFileManage{}).Create(r).Error
	if err != nil {
		return err
	}

	return nil
}

// ModifyOnlineFile
// @function: ModifyOnlineFile
// @description: Update online download
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) ModifyOnlineFile(info *pb.AddFileInfo) error {

	if info.FileType != define.GenerallyFile.Int64() && info.FileType != define.ProveFile.Int64() && info.FileType != define.SnapshotFile.Int64() {
		return errors.New("file_type is error")
	}

	if utils.IsNull(info.FileUrl) {
		return errors.New("file_url is nil")
	}

	var sysFileManage system.SysFileManage
	err := global.ZC_DB.Model(&system.SysFileManage{}).Where("id", info.Id).Find(&sysFileManage).Error
	if err != nil {
		return err
	}

	if info.FileType == define.ProveFile.Int64() {
		var fileCount int64
		err = global.ZC_DB.Model(&system.SysFileManage{}).Where("id<>? and file_type=?", info.Id, info.FileType).Count(&fileCount).Error
		if err != nil {
			return err
		}
		if fileCount > utils.ZERO {
			return errors.New("Only one certificate file can be downloaded from each equipment room")
		}
	}

	if (sysFileManage == system.SysFileManage{}) {
		return errors.New("Update file does not exist")
	}

	res, err := new(gateway.DownloadService).DowloadFile(info.FileUrl, define.FileGateWayDir, "")
	if err != nil {
		log.Println("DownloadFile err:", err.Error())
	}

	updateMap := make(map[string]interface{})
	if res != nil {
		beginIndex := strings.LastIndex(res.FilePath, `/`)
		if beginIndex == -1 {
			beginIndex = strings.LastIndex(res.FilePath, `\`)
		}
		fileName := utils.SubStr(res.FilePath, beginIndex+1, len(res.FilePath))
		updateMap["file_name"] = fileName
		updateMap["file_size"] = res.Total
		updateMap["file_status"] = define.FileGetting.Int()
	} else {
		updateMap["file_name"] = "未知文件"
		updateMap["file_size"] = utils.ZERO
		updateMap["file_status"] = define.FileError.Int()
	}

	updateMap["file_url"] = info.FileUrl
	err = global.ZC_DB.Model(&system.SysFileManage{}).Where("id", info.Id).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

// AddFileGrpc
// @function: AddFileGrpc
// @description: New file
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) AddFileGrpc(info request.AddFileReq) error {

	fileAr := make([]*pb.FileInfo, 0)
	for _, v := range info.UpLoadFile {

		f := &pb.FileInfo{
			FileName: v.FileName,
			FileSize: v.FileSize,
		}
		fileAr = append(fileAr, f)
	}

	data := &pb.AddFileInfo{
		RoomId:      info.RoomId,
		RoomName:    info.RoomName,
		GateWayId:   info.GateWayId,
		FileType:    info.FileType,
		AddType:     info.AddType,
		OpId:        info.OpId,
		Ip:          info.Ip,
		FileUrl:     info.FileUrl,
		FileInfo:    fileAr,
		ZipFileName: info.ZipFileName,
	}

	if len(info.OpList) > 0 {
		data.Ip = info.OpList[0].Ip
		data.Port = info.OpList[0].Port
	}

	if info.AddType == define.AddLocalUpload.Int64() {
		err := JobPlatformServiceApp.FileLocalSynGateWay(data)
		if err != nil {
			return err
		}
	} else {

		client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
		if client == nil {
			return errors.New("AddFileGrpc GateWayClient Connection failed:" + info.GateWayId)
		}
		result, _ := client.FileOpSynGateWay(context.TODO(), data)
		if result.Code != 200 {
			return errors.New(result.Msg)
		}
	}

	return nil
}

// FileNodeCopy
// @function: FileNodeCopy
// @description: File node copy
// @param: info *pb.AddFileInfo
// @return: error
func (jp *JobPlatformService) FileNodeCopy(info *pb.AddFileInfo) (int, error) {

	if info.FileType != define.HeightFile.Int64() && info.FileType != define.MinerFile.Int64() {
		return 0, errors.New("file_type is err")
	}
	ipAddress := utils.SubStr(info.Ip, strings.LastIndex(info.Ip, ".")+1, len(info.Ip))
	if len(info.FileInfo) == 0 {

		if info.ZipFileName == "" {
			if info.FileType == define.HeightFile.Int64() {
				info.ZipFileName = "lotus" + ipAddress
			} else {
				info.ZipFileName = "miner" + ipAddress
			}
		}

		f := &pb.FileInfo{
			FileName: info.ZipFileName,
		}
		info.FileInfo = append(info.FileInfo, f)
		info.FileName = info.ZipFileName + ".zip"
	} else {
		if len(info.FileInfo) == 1 {
			info.FileName = info.FileInfo[0].FileName
		} else {
			return 0, errors.New("batch file copying is not supported")
		}
	}

	b, err := JobPlatformServiceApp.GetFileStatus(info.ZipFileName)
	if err != nil {
		global.ZC_LOG.Error("GetFileStatus", zap.String("err:", err.Error()))
		return 0, err
	}
	if b {

		return 0, errors.New("during file synchronization, please repeat the operation")
	}

	id, err := JobPlatformServiceApp.AddFileManage(info)
	if err != nil {
		global.ZC_LOG.Error("AddFileManage", zap.String("err:", err.Error()))
		return 0, err
	}

	return id, nil
}

// FileOpSynGateWay
// @function: FileOpSynGateWay
// @description: Synchronize file op to gateway
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) FileOpSynGateWay(info *pb.AddFileInfo) error {

	RunStopService(info.FileType, define.StopService.String())
	fileMap := make(map[string]struct{})

	if info.FileType == define.HeightFile.Int64() {
		info.FromPath = define.PathIpfsLotusDatastore
	} else if info.FileType == define.MinerFile.Int64() {
		info.FromPath = define.PathIpfsMiner
	}
	info.ToPath = define.FileGateWayDir

	for _, v := range info.FileInfo {
		fileMap[v.FileName] = struct{}{}
	}

	dirPath, fileName := "", ""
	if info.ZipFileName != "" {
		fileName = info.ZipFileName + ".zip"
		if info.FileType == 3 {
			dirPath = define.PathIpfsLotusDatastore
		} else if info.FileType == 5 {
			dirPath = define.PathIpfsMiner
		}
		if dirPath != "" {
			bool, err := ZipFiles(dirPath, fileName, nil)
			if err != nil {
				return err
			}

			if bool && len(info.FileInfo) > 0 {
				fileMap[fileName] = struct{}{}
			}
		}
	}

	log.Println("fromPath:" + info.FromPath)
	log.Println("toPath:" + info.ToPath)
	fileInfoList, err := ioutil.ReadDir(info.FromPath)
	if err != nil {
		return err
	}

	if info.ZipFileName == "" {
		err = FileExist(fileInfoList, fileMap)
		if err != nil {
			return err
		}
	}

	fileSize := 0
	for i := range fileInfoList {

		if !utils.IsFile(info.FromPath + "/" + fileInfoList[i].Name()) {
			continue
		}

		if len(fileMap) > 0 {
			if _, ok := fileMap[fileInfoList[i].Name()]; !ok {
				continue
			}
		}

		fileSize = int(fileInfoList[i].Size())
		sy := &request.SynFileReq{
			FromPath:  info.FromPath + "/" + fileInfoList[i].Name(),
			ToPath:    info.ToPath,
			FileName:  fileInfoList[i].Name(),
			GateWayId: info.GateWayId,
			OpId:      info.OpId,
			Ip:        info.Ip,
			Port:      info.Port,
		}

		_, err = GateWayFileToByte(sy, global.OpToGatewayClient)
		if err != nil {
			return err
		}
	}

	id, _ := strconv.Atoi(info.Id)
	global.OpToGatewayClient.ModifyFileStatus(context.TODO(), &pb.FileManage{Id: int64(id), FileSize: int64(fileSize)})

	os.Remove(dirPath + "/" + fileName)

	RunStopService(info.FileType, define.StartService.String())

	return nil
}

// OpLocalFileSynGateWay
// @function: OpLocalFileSynGateWay
// @description: op Synchronize local files to the gateway
// @param: info info *pb.AddFileInfo
// @return: error
func (jp *JobPlatformService) OpLocalFileSynGateWay(info *pb.AddFileInfo) (chunks []byte, err error) {

	//先删除再同步
	global.OpToGatewayClient.DelGateWayFile(context.TODO(), &pb.FileInfo{FileName: fmt.Sprintf("%s/%s", info.ToPath, info.FileName)})

	var reader *bytes.Reader
	if len(info.FileData) == utils.ZERO {
		fileData, err := ioutil.ReadFile(info.FromPath + "/" + info.FileName)
		if err != nil {
			return chunks, err
		}
		reader = bytes.NewReader(fileData)
	} else {
		reader = bytes.NewReader(info.FileData)
	}

	beginTime := time.Now()
	for {

		var n int
		dataByte := make([]byte, define.FileByte.Int())
		n, err = reader.Read(dataByte)
		if err != nil || 0 == n {
			break
		}

		_, err := global.OpToGatewayClient.AddGateWayFile(context.TODO(), &pb.FileInfo{Path: info.ToPath, FileName: info.FileName, FileData: dataByte[:n]})
		if err != nil {
			log.Println("AddGateWayFile err:", err.Error())
			return chunks, err
		}
		chunks = append(chunks, dataByte[:n]...)
	}
	isEOF := strings.Compare(err.Error(), "EOF")
	if isEOF == 0 {
		err = nil
		fmt.Println(fmt.Printf("read %s success, len=%v, Total time spent: %s", info.ToPath, len(chunks), time.Now().Sub(beginTime).String()))
		return chunks, nil
	}

	return
}

// FileLocalSynGateWay
// @function: FileLocalSynGateWay
// @description: Synchronize files locally to the gateway
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) FileLocalSynGateWay(info *pb.AddFileInfo) error {

	if info.FileType != define.GenerallyFile.Int64() && info.FileType != define.ProveFile.Int64() && info.FileType != define.SnapshotFile.Int64() {
		return errors.New("file_type is error")
	}
	if len(info.FileInfo) == 0 {
		return errors.New("up_load_file[] is nil")
	}

	fileMap := make(map[string]struct{})
	info.FromPath = define.FileUploadDir
	info.ToPath = define.FileGateWayDir

	for _, v := range info.FileInfo {
		fileMap[v.FileName] = struct{}{}
	}

	log.Println("fromPath:" + info.FromPath)
	log.Println("toPath:" + info.ToPath)
	fileInfoList, err := ioutil.ReadDir(info.FromPath)
	if err != nil {
		return err
	}

	err = FileExist(fileInfoList, fileMap)
	if err != nil {
		return err
	}

	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New("FileLocalSynGateWay GateWayClient Connection failed" + info.GateWayId)
	}

	var wg sync.WaitGroup
	wg.Add(len(fileInfoList))
	for i := range fileInfoList {

		if !utils.IsFile(info.FromPath + "/" + fileInfoList[i].Name()) {
			continue
		}

		if len(fileMap) > 0 {
			if _, ok := fileMap[fileInfoList[i].Name()]; !ok {
				continue
			}
		}

		go func(index int) {

			defer wg.Done()

			r := &system.SysFileManage{
				RoomId:     info.RoomId,
				RoomName:   info.RoomName,
				GateWayId:  info.GateWayId,
				FileName:   fileInfoList[index].Name(),
				FileSize:   int(fileInfoList[index].Size()),
				FileType:   int(info.FileType),
				FileUrl:    info.FileUrl,
				FileStatus: utils.ONE,
			}

			err = global.ZC_DB.Model(&system.SysFileManage{}).Create(r).Error
			if err != nil {
				global.ZC_LOG.Error("SysFileManage", zap.String("sql", err.Error()))
			}

			sy := &request.SynFileReq{
				FromPath: info.FromPath + "/" + fileInfoList[index].Name(),
				ToPath:   info.ToPath,
				FileName: fileInfoList[index].Name(),
			}

			_, err = GateWayFileToByte(sy, client)
			if err != nil {
				global.ZC_LOG.Error("GateWayFileToByte", zap.String("err", err.Error()))
			}

			err = global.ZC_DB.Model(&system.SysFileManage{}).Where("id", r.ZC_MODEL.ID).Update("file_status", utils.TWO).Error
			if err != nil {
				global.ZC_LOG.Error("SysFileManage", zap.String("sql", err.Error()))
			}

			os.Remove(sy.FromPath)

		}(i)

	}

	return nil
}

func FileExist(fileInfoList []fs.FileInfo, fileMap map[string]struct{}) error {

	existMap := make(map[string]struct{})
	if len(fileMap) > 0 {
		for i := range fileInfoList {
			if _, ok := fileMap[fileInfoList[i].Name()]; ok {
				existMap[fileInfoList[i].Name()] = struct{}{}
			}
		}
		if len(existMap) == 0 || len(existMap) != len(fileMap) {
			fileAr := ""
			for k, _ := range fileMap {
				if _, ok := existMap[k]; !ok {
					fileAr += k + "File not found,"
				}
			}
			if fileAr != "" {
				fileAr = utils.SubStr(fileAr, 0, len([]rune(fileAr))-1)
				return errors.New(fileAr)
			}
		}
	}
	return nil
}

// DownLoadFileToOp
// @function: DownLoadFileToOp
// @description: Download the file gateway to Op
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) DownLoadFileToOp(info *pb.AddFileInfo) error {

	global.ZC_LOG.Info(fmt.Sprintf("DownLoadFileToOp param:%+v", info))
	fileMap := make(map[string]struct{})
	for _, v := range info.FileInfo {
		fileMap[v.FileName] = struct{}{}
	}

	log.Println("fromPath:" + info.FromPath)
	log.Println("toPath:" + info.ToPath)
	fileInfoList, err := ioutil.ReadDir(info.FromPath)
	if err != nil {
		return err
	}

	err = FileExist(fileInfoList, fileMap)
	if err != nil {
		return err
	}

	var client pb.OpServiceClient
	if info.Ip != "" && info.Port != "" {

		conn, err := utils.GrpcConnect(info.Ip, info.Port)
		if err != nil {
			return err
		}

		client = pb.NewOpServiceClient(conn)
		defer func() {
			if conn != nil {
				conn.Close()
			}
		}()

	} else {
		clientService, dis := global.OpClinets.GetOpClient(info.OpId)
		if dis {
			return errors.New("GetOpClient Connection failed:" + info.OpId)
		}
		client = clientService
	}
	if client == nil {
		return errors.New("GrpcConnect Connection failed:" + info.Ip + ":" + info.Port)
	}

	for i := range fileInfoList {

		if !utils.IsFile(info.FromPath + "/" + fileInfoList[i].Name()) {
			continue
		}

		if len(fileMap) > 0 {
			if _, ok := fileMap[fileInfoList[i].Name()]; !ok {
				continue
			}
		}

		sy := &request.SynFileReq{
			FromPath:  info.FromPath + "/" + fileInfoList[i].Name(),
			ToPath:    info.ToPath,
			FileName:  fileInfoList[i].Name(),
			GateWayId: info.GateWayId,
			OpId:      info.OpId,
			Ip:        info.Ip,
			Port:      info.Port,
		}

		_, err = PointFileToByte(sy, client)
		if err != nil {
			return err
		}

	}

	return nil
}

// RunStopService Disable or enable a service
func RunStopService(fileType int64, str string) {

	command := ""
	if fileType == define.HeightFile.Int64() {
		command = fmt.Sprintf("supervisorctl %s %s", str, define.ProgramLotus.String())
	} else {
		command = fmt.Sprintf("supervisorctl %s %s", str, define.ProgramMiner.String())
	}

	b, err := utils.ExecScript(command)
	if err != nil {
		global.ZC_LOG.Error("ExecScript", zap.String("err", err.Error()))
	}
	log.Println("ExecScript", zap.String("b", command+"|"+string(b)))

}

// AddFileManage
// @function: AddFileManage
// @description: Add file management
// @param: info *pb.AddFileInfo
// @return: error
func (jp *JobPlatformService) AddFileManage(info *pb.AddFileInfo) (int, error) {

	t := &system.SysFileManage{
		RoomId:     info.RoomId,
		RoomName:   info.RoomName,
		GateWayId:  info.GateWayId,
		OpId:       info.OpId,
		FileName:   info.FileName,
		FileSize:   0,
		FileType:   int(info.FileType),
		FileStatus: 1,
	}

	if utils.IsNull(t.RoomName) {

		var res system.SysMachineRoomRecord
		err := global.ZC_DB.Model(&system.SysMachineRoomRecord{}).Where("gateway_id", info.GateWayId).Find(&res).Limit(utils.ONE).Error
		if err != nil {
			return 0, err
		}
		t.RoomId = res.RoomId
		t.RoomName = res.RoomName
	}

	return int(t.ZC_MODEL.ID), global.ZC_DB.Model(&system.SysFileManage{}).Create(&t).Error
}

// ModifyFileStatus
// @function: ModifyFileStatus
// @description: Update file status
// @param: id int
// @return: error
func (jp *JobPlatformService) ModifyFileStatus(id, fileSize int) error {

	updateMap := make(map[string]interface{})
	updateMap["file_status"] = 2
	updateMap["file_size"] = fileSize
	err := global.ZC_DB.Model(&system.SysFileManage{}).Where("id=?", id).Updates(updateMap).Error
	if err != nil {
		global.ZC_LOG.Error("SysFileManage", zap.String("err", err.Error()))
		return err
	}
	return nil
}

// GetFileStatus
// @function: GetFileStatus
// @description:File status
// @param: info *pb.AddFileInfo
// @return: error
func (jp *JobPlatformService) GetFileStatus(fileName string) (bool, error) {

	var res []system.SysFileManage
	err := global.ZC_DB.Model(&system.SysFileManage{}).Where(" file_status=1 and file_name=?", fileName+".zip").Scan(&res).Error
	if err != nil {
		return false, err
	}

	if len(res) > 0 {
		return true, nil
	}

	return false, nil
}

// SysFilePoint
// @function: SysFilePoint
// @description: op endpoint synchronizes files from point to point
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) SysFilePoint(info request.SynFileReq) error {

	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New("sysFilePoint GateWayClient Connection failed" + info.GateWayId)
	}

	fileMap := make([]*pb.FileInfo, 0)
	for _, v := range info.Files {
		file := &pb.FileInfo{
			FileName: v.FileName,
		}
		fileMap = append(fileMap, file)
	}

	t := &pb.SynFileInfo{
		OpId:     info.OpId,
		ToOpId:   info.ToOpId,
		FromPath: info.FromPath,
		ToPath:   info.ToPath,
		Ip:       info.Ip,
		Port:     info.Port,
		FileInfo: fileMap,
	}

	res, _ := client.SysFilePoint(context.TODO(), t)
	if res.Code != 200 {
		return errors.New(res.Msg)
	}

	return nil
}

// FilePointProcess
// @function: FilePointProcess
// @description: op endpoint synchronizes files from point to point
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) FilePointProcess(info *pb.SynFileInfo) error {

	fileMap := make(map[string]struct{})
	for _, v := range info.FileInfo {
		fileMap[v.FileName] = struct{}{}
	}

	log.Println("fromPath:" + info.FromPath)
	log.Println("toPath:" + info.ToPath)

	checkZip := false
	fileName := "new_zip.zip"
	if info.ZipFileName != "" {
		fileName = info.ZipFileName
	}

	if len(info.FileInfo) == 0 {

		os.Remove(info.FromPath + "/" + fileName)

		bool, err := ZipFiles(info.FromPath, fileName, info.FileInfo)
		if err != nil {
			return err
		}
		checkZip = bool
		fileMap[fileName] = struct{}{}
	}

	fileInfoList, err := ioutil.ReadDir(info.FromPath)
	if err != nil {
		return err
	}

	err = FileExist(fileInfoList, fileMap)
	if err != nil {
		return err
	}

	conn, err := utils.GrpcConnect(info.Ip, info.Port)
	if err != nil {
		return err
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	client := pb.NewOpServiceClient(conn)
	if client == nil {
		return errors.New("GrpcConnect Connection failed:" + info.Ip + ":" + info.Port)
	}

	for i := range fileInfoList {

		if !utils.IsFile(info.FromPath + "/" + fileInfoList[i].Name()) {
			continue
		}

		if len(fileMap) > 0 {
			if _, ok := fileMap[fileInfoList[i].Name()]; !ok {
				continue
			}
		}

		sy := &request.SynFileReq{
			FromPath: info.FromPath + "/" + fileInfoList[i].Name(),
			ToPath:   info.ToPath,
			FileName: fileInfoList[i].Name(),
		}

		_, err = PointFileToByte(sy, client)
		if err != nil {
			return err
		}
	}

	if checkZip {

		os.Remove(info.FromPath + "/" + fileName)
		t := &pb.FileInfo{
			FileName: info.ToPath + "/" + fileName,
			Path:     info.ToPath,
		}
		_, err = pb.NewOpServiceClient(conn).UnZipSynFile(context.TODO(), t)
		if err != nil {
			log.Println("conn", zap.String("err", err.Error()))
		}
	}

	return nil
}

// UnZipSynFile
// @function: UnZipSynFile
// @description: Unzip file
// @param: info request.AddFileReq
// @return: error
func (jp *JobPlatformService) UnZipSynFile(info *pb.FileInfo) error {

	err := utils.UnzipLocal(info.FileName, info.Path)
	if err != nil {
		return err
	}
	err = JobPlatformServiceApp.DelFile(info)
	return err
}

// DelFile
// @function: DelFile
// @description: Delete file
// @param: info *pb.FileInfo
// @return: error
func (jp *JobPlatformService) DelFile(info *pb.FileInfo) error {
	log.Println("DelFile:", info.FileName)
	return os.Remove(info.FileName)
}

// GetFileListByType
// @function: GetFileListByType
// @description: File management list
// @param: info request.FileTypeReq
// @return: list interface{}, err error
func (jp *JobPlatformService) GetFileListByType(info request.FileTypeReq) (list interface{}, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysFileManage{}).Where("file_type = ?", info.FileType).Where("gate_way_id = ?", info.GateId)
	var res []response.FileInfo

	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, err
}

// DownLoadFile
// @function: DownLoadFile
// @description: File download
// @param: info request.FileTypeReq
// @return: list interface{}, err error
func (jp *JobPlatformService) DownLoadFile(info request.DownLoadReq) error {

	if len(info.Files) == 0 {
		return errors.New("files[] is nil")
	}
	client := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
	if client == nil {
		return errors.New("gatewayClient Connection failed" + info.GateWayId)
	}

	opMap := make([]*pb.OpInfo, 0)
	fileMap := make([]*pb.FileInfo, 0)
	for _, v := range info.OpList {
		t := &pb.OpInfo{
			Ip:   v.Ip,
			Port: v.Port,
		}
		opMap = append(opMap, t)
	}
	for _, v := range info.Files {
		t := &pb.FileInfo{
			FileName: v.FileName,
		}
		fileMap = append(fileMap, t)
	}

	t := &pb.DownLoadInfo{
		DownloadPath: info.DownloadPath,
		OpInfo:       opMap,
		FileInfo:     fileMap,
	}
	res, _ := client.DownLoadFiles(context.TODO(), t)
	if res.Code != 200 {
		return errors.New(res.Msg)
	}

	t1 := &pb.FileInfo{
		FileName: define.FileGateWayDir + "/new_zip.zip",
		Path:     define.FileGateWayDir,
	}
	client.DelGateWayFile(context.TODO(), t1)
	return nil
}

// ReadGateWayFile
// @function: ReadGateWayFile
// @description: Reading the GateWay file
// @param: info request.FileTypeReq
// @return: list interface{}, err error
func (jp *JobPlatformService) ReadGateWayFile(info *pb.DownLoadInfo) error {

	filePath := define.FileGateWayDir
	if len(info.GateWayPath) != 0 {
		filePath = info.GateWayPath
	}

	fileName := "new_zip.zip"
	bool, err := ZipFiles(filePath, fileName, info.FileInfo)
	if err != nil {
		return err
	}

	unZipAr := make([]string, 0)
	if bool {
		info.FileInfo = append(info.FileInfo, &pb.FileInfo{FileName: fileName})
		unZipAr = append(unZipAr, fileName)
	}

	errMap := make(map[int]string)
	var wg sync.WaitGroup
	wg.Add(len(info.OpInfo))
	for j, v := range info.OpInfo {

		go func(p *pb.OpInfo, p1 *pb.DownLoadInfo, index int, m map[int]string) {

			defer wg.Done()
			param := &pb.AddFileInfo{
				Ip:       p.Ip,
				Port:     p.Port,
				FromPath: filePath,
				ToPath:   p1.DownloadPath,
				FileInfo: p1.FileInfo,
				OpId:     p.OpId,
			}

			if len(p.OpId) != 0 {
				param.Port = ""
			}

			err := JobPlatformServiceApp.DownLoadFileToOp(param)
			if err != nil {
				log.Println("DownLoadFileToOp", zap.String("err", err.Error()))
				errMap[index] = "ip:" + p.Ip + err.Error()
			}

		}(v, info, j, errMap)
	}

	wg.Wait()

	if len(errMap) > 0 {
		errMsg := ""
		for _, v := range errMap {
			errMsg += v + ","
		}
		if errMsg != "" {
			errMsg = utils.SubStr(errMsg, 0, len([]rune(errMsg))-1)
			return errors.New(errMsg)
		}
	}

	if len(unZipAr) > utils.ZERO {

		var wg1 sync.WaitGroup
		wg1.Add(len(info.OpInfo))
		for _, v := range info.OpInfo {

			go func(p *pb.OpInfo) {

				defer wg1.Done()
				conn, err := utils.GrpcConnect(p.Ip, p.Port)
				if err != nil {
					global.ZC_LOG.Error("GrpcConnect", zap.String("err", "Connection failed"+p.Ip+":"+p.Port))
				}
				defer func() {
					if conn != nil {
						conn.Close()
					}
				}()

				if utils.IsNull(fileName) {
					return
				}

				t := &pb.FileInfo{
					FileName: info.DownloadPath + "/" + fileName,
					Path:     info.DownloadPath,
				}
				_, err = pb.NewOpServiceClient(conn).UnZipSynFile(context.TODO(), t)
				if err != nil {
					log.Println("conn", zap.String("err", err.Error()))
				}
			}(v)
		}
		wg1.Wait()
	}

	os.Remove(path.Join(define.FileGateWayDir, fileName))

	return nil
}

// GetFileHostList
// @function: GetFileHostList
// @description: File host list
// @param:nil
// @return: res system.SysDictionaryConfig, err error
func (jp *JobPlatformService) GetFileHostList() (res []system.SysDictionaryConfig, err error) {

	var r []system.SysDictionaryConfig
	err = global.ZC_DB.Model(&system.SysDictionaryConfig{}).Where("relation_id", "000001").Scan(&r).Error
	return r, err
}

// GetLotusHeightList
// @function: GetLotusHeightList
// @description: Gets a list of lotus height files
// @param:nil
// @return: res system.SysDictionaryConfig, err error
func (jp *JobPlatformService) GetLotusHeightList(info request.GateWayReq) (res interface{}, err error) {

	var lotusInfo []lotus.LotusInfo
	if info.Id == "" {
		return lotusInfo, errors.New("id is nil")
	}
	err = global.ZC_DB.Model(&lotus.LotusInfo{}).Where("gate_id", info.Id).Scan(&lotusInfo).Error
	if err != nil {
		return nil, err
	}

	if len(lotusInfo) == 0 {
		return lotusInfo, nil
	}

	type LotusHeight struct {
		OpId   string `json:"op_id"`
		IP     string `json:"ip"`
		Port   string `json:"port"`
		Height int    `json:"height"`
	}

	client := global.GateWayClinets.GetGateWayClinet(lotusInfo[0].GateId)
	if client == nil {
		return nil, errors.New("gatewayClient Connection failed:" + lotusInfo[0].GateId)
	}

	resAr := make([]LotusHeight, 0)
	for _, v := range lotusInfo {

		h, err := client.LotusHeight(context.TODO(), &pb.RequestOp{Ip: v.Ip, OpId: v.OpId, Token: v.Token})
		if err != nil {
			continue
		}

		t := &LotusHeight{
			OpId:   v.OpId,
			IP:     v.Ip,
			Port:   v.Port,
			Height: int(h.Height),
		}
		resAr = append(resAr, *t)
	}

	return resAr, err
}

// ZipFiles Compressed file
func ZipFiles(dir, fileName string, info []*pb.FileInfo) (bool, error) {

	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}

	fileMap := make(map[string]struct{})
	for i := range fileInfoList {
		fileMap[fileInfoList[i].Name()] = struct{}{}
	}

	zipMap := make(map[string]string)
	if len(info) > 0 {
		for _, v := range info {
			filePath := dir + "/" + v.FileName
			if _, ok := fileMap[v.FileName]; ok && !utils.IsFile(filePath) {
				zipMap[v.FileName] = filePath
			}
		}
	} else {
		for i := range fileInfoList {
			filePath := dir + "/" + fileInfoList[i].Name()
			if fileName != "new_zip.zip" {
				nameAr := strings.Split(fileName, ".")
				if fileInfoList[i].Name() != nameAr[0] {
					continue
				}
			}
			zipMap[fileInfoList[i].Name()] = filePath
		}
	}

	skipFile := "hot.badger"
	zipBool := false
	if len(zipMap) > 0 {
		zipAr := make([]string, 0)
		for _, v := range zipMap {
			zipAr = append(zipAr, v)
		}
		err = utils.Zip(dir+"/"+fileName, skipFile, zipAr...)
		if err != nil {
			return false, err
		}
		zipBool = true
	}

	return zipBool, nil
}

// GetMinerList
// @function: GetMinerList
// @description: Get a list of miner hosts
// @param:nil
// @return: res system.SysDictionaryConfig, err error
func (jp *JobPlatformService) GetMinerList(info request.GateWayReq) (res []response.MinerInfoRes, err error) {

	var minerList []response.MinerInfoRes
	if info.Id == "" {
		return minerList, errors.New("id is nil")
	}
	err = global.ZC_DB.Model(&lotus.LotusMinerInfo{}).Where("gate_id", info.Id).Scan(&minerList).Error
	return minerList, err
}
