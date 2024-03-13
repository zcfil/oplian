package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/go-units"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/model/system"
	"oplian/service/lotus/deploy/cgo"
	"oplian/service/lotus/oplocal"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var SectorsRecoverServiceApi = new(SectorsRecoverService)

type SectorsRecoverService struct {
}

// GetSectorsRecoverList
// @function: GetWorkerClusterList
// @description: Gets the sector list
// @param: info request.SectorsRecover
// @return: list interface{}, total int64, err error
func (s *SectorsRecoverService) GetSectorsRecoverList(info request.SectorsRecover) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&lotus.LotusSectorRecover{})
	var sectorRecover []lotus.LotusSectorRecover
	sql := "1=1"
	var param []interface{}
	if info.BelongingNode != "" {
		sql += " and miner_id = ?"
		param = append(param, info.BelongingNode)
	}
	if info.SectorId != "" {
		info.SectorId = "%" + info.SectorId + "%"
		sql += " and sector_id like ?"
		param = append(param, info.SectorId)
	}
	if info.SectorSize > 0 {
		sql += " and sector_size = ? "
		param = append(param, info.SectorSize)
	}
	if info.SectorType > 0 {
		sql += " and sector_type = ?"
		param = append(param, info.SectorType)
	}
	if info.SectorStatus > 0 {
		sql += " and sector_status = ?"
		param = append(param, info.SectorStatus)
	}
	err = db.Where(sql, param...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&sectorRecover).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return sectorRecover, total, err
}

func (s *SectorsRecoverService) AddSectorsRecoverTask(info request.SectorsRecoverTask) (bool, error) {

	err := global.ZC_DB.Transaction(func(tx *gorm.DB) error {

		if len(info.Ids) == utils.ZERO {
			return errors.New("ids[] data is nil ")
		}
		if len(info.WorkerOp) == utils.ZERO {
			return errors.New("op_relations[] data is nil ")
		}

		//保存扇区任务
		lotusSectorTask := &request.LotusSectorTask{
			TaskName:          info.TaskName,
			Actor:             info.Actor,
			SectorType:        info.SectorType,
			SectorSize:        info.SectorSize,
			SectorTotal:       info.SectorTotal,
			OriginalValueOpId: info.OriginalValueOpId,
			OriginalValueDir:  info.OriginalValueDir,
			StorageOpId:       info.StorageOpId,
			StorageOpIp:       info.StorageOpIp,
			StorageOpName:     info.StorageOpName,
			TaskStatus:        utils.ZERO,
		}
		err := tx.Save(lotusSectorTask).Error
		if err != nil {
			return err
		}

		//保存扇区任务明细
		idsStr := ""
		for _, v := range info.Ids {
			idsStr += `'` + v.Id + `',`
		}
		if idsStr != "" {
			idsStr = utils.SubStr(idsStr, 0, len(idsStr)-1)
		}

		var sectorRecover []lotus.LotusSectorRecover
		err = tx.Model(&lotus.LotusSectorRecover{}).Where("id in(" + idsStr + ")").Find(&sectorRecover).Error
		if err != nil {
			return err
		}

		if len(sectorRecover) == utils.ZERO {
			return errors.New("sectorRecover[] data is nil")
		}

		//更新扇区状态
		err = tx.Model(&lotus.LotusSectorRecover{}).Where("id in("+idsStr+")").Update("sector_status", utils.FOUR).Error
		if err != nil {
			return err
		}

		sectorTaskDetail := make([]lotus.LotusSectorTaskDetail, 0)
		for _, v := range sectorRecover {

			data := &lotus.LotusSectorTaskDetail{
				SectorRecoverId: int(v.ZC_MODEL.ID),
				RelationId:      strconv.Itoa(int(lotusSectorTask.ZC_MODEL.ID)),
				MinerId:         v.MinerId,
				SectorId:        v.SectorId,
				Ticket:          v.Ticket,
				SectorStatus:    0,
				SectorSize:      info.SectorSize,
			}
			sectorTaskDetail = append(sectorTaskDetail, *data)
		}

		err = tx.Create(&sectorTaskDetail).Error
		if err != nil {
			return err
		}

		//保存worker主机
		idsStr = ""
		for _, v := range info.WorkerOp {
			idsStr += `'` + v.Id + `',`
		}
		if idsStr != "" {
			idsStr = utils.SubStr(idsStr, 0, len(idsStr)-1)
		}

		var hostRecord []system.SysHostRecord
		err = global.ZC_DB.Model(&system.SysHostRecord{}).Where("uuid in(" + idsStr + ")").Find(&hostRecord).Error
		if err != nil {
			return err
		}

		if len(hostRecord) == utils.ZERO {
			return errors.New("hostRecord[] data is nil")
		}

		opRelationsAr := make([]system.SysOpRelations, 0)
		for _, v := range hostRecord {

			data := &system.SysOpRelations{
				RelationId:    strconv.Itoa(int(lotusSectorTask.ZC_MODEL.ID)),
				RelationType:  utils.FOUR,
				GateWayId:     v.GatewayId,
				OpId:          v.UUID,
				Ip:            v.IntranetIP,
				ServerName:    v.HostName,
				AssetsNum:     v.AssetNumber,
				DeviceSn:      v.DeviceSN,
				OperateSystem: v.OperatingSystem,
				RoomId:        v.RoomId,
				RoomName:      v.RoomName,
				Status:        utils.ZERO,
				BeginTime:     time.Now(),
			}
			opRelationsAr = append(opRelationsAr, *data)
		}

		err = tx.Create(&opRelationsAr).Error
		if err != nil {
			return err
		}

		return nil
	})

	return true, err
}

func (s *SectorsRecoverService) RedoSectorsTask() {

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {

		select {
		case <-ticker.C:
			err := s.DoSectorsTask()
			if err != nil {
				log.Println("DoSectorsTask err:", err)
			}
		}
	}
}

func (s *SectorsRecoverService) RunSectorSealTask() {

	time.Sleep(time.Second * 30)
	script := "killall -9 oplian-sectors-task"
	_, err := utils.ExecuteScript(script)
	if err != nil {
		log.Println("kill sectorTask ExecuteScript err:", err)
	}
	log.Println(script)

	script = fmt.Sprintf(define.PathIpfsScriptRunSectorTask+" %s %s", define.OpSectorTaskPort, define.MainDisk)
	_, err = utils.ExecuteScript(script)
	if err != nil {
		log.Println("run sectorTask ExecuteScript err:", err)
	} else {
		log.Println("sectorTask Successful startup:", script)
	}

}

func (s *SectorsRecoverService) DoSectorsTask() error {

	var lotusSectorTask []lotus.LotusSectorTask
	err := global.ZC_DB.Model(&lotus.LotusSectorTask{}).Where("task_status in(0,1)").Find(&lotusSectorTask).Error
	if err != nil {
		return err
	}

	if len(lotusSectorTask) > utils.ZERO {

		for _, t := range lotusSectorTask {

			var sysOpRelations []system.SysOpRelations
			err = global.ZC_DB.Model(&system.SysOpRelations{}).Where("relation_id=? and relation_type=?", t.ID, utils.FOUR).Find(&sysOpRelations).Error
			if err != nil {
				return err
			}

			if len(sysOpRelations) == utils.ZERO {
				return errors.New(fmt.Sprintf("DoSectorsTask ID:%d 无可执行worker机!", t.ID))
			}

			var ls lotus.LotusSectorTaskDetail
			err = global.ZC_DB.Model(&lotus.LotusSectorTaskDetail{}).Where("relation_id= ? and task_order > ?", t.ID, utils.ZERO).Order("task_order desc").Limit(utils.ONE).Find(&ls).Error
			if err != nil {
				ls = lotus.LotusSectorTaskDetail{TaskOrder: utils.ZERO}
			}

			var lotusSectorDetail []lotus.LotusSectorTaskDetail

			for _, v1 := range sysOpRelations {

				client, dis := global.OpClinets.GetOpClient(v1.OpId)
				if dis {
					log.Println(fmt.Sprintf("DoSectorsTask opClient Connection failed:%s", v1.OpId))
					continue
				}

				sectorRefTask := make([]*pb.SectorRefTask, 0)
				err = global.ZC_DB.Model(&lotus.LotusSectorTaskDetail{}).Where(" sector_status = ? and relation_id= ?", 0, t.ID).Find(&lotusSectorDetail).Error
				if err != nil {
					return err
				}

				if len(lotusSectorDetail) > utils.ZERO {

					for i := 0; i < len(lotusSectorDetail); i++ {

						v := lotusSectorDetail[i]
						tk := pb.Task{Wid: define.RedoSectorsTask.String(), TType: string(define.TTAddPiece), MinerId: v.MinerId, SectorId: uint64(v.SectorId)}
						checkAp, err := client.Ok(context.TODO(), &tk)
						if checkAp != nil {

							if !checkAp.Value {

								time.Sleep(time.Minute)
								break
							}
						} else {

							time.Sleep(time.Minute)
							break
						}

						sectorRefTask = append(sectorRefTask, &pb.SectorRefTask{
							Id: &pb.SectorID{
								Miner:  v.MinerId,
								Number: uint64(v.SectorId),
							},
							SectorSize:      uint64(v.SectorSize),
							Ticket:          []byte(v.Ticket),
							TaskDetailId:    uint64(v.ID),
							SectorRecoverId: uint64(v.SectorRecoverId),
						})

						updateMap := make(map[string]interface{})
						updateMap["sector_status"] = utils.ONE
						if ls.TaskOrder == utils.ZERO {
							updateMap["task_order"] = i + utils.ONE
						} else {
							updateMap["task_order"] = ls.TaskOrder + utils.ONE
						}
						err = global.ZC_DB.Model(&lotus.LotusSectorTaskDetail{}).Where("id =?", v.ID).Updates(updateMap).Error
						if err != nil {
							return err
						}
					}

					if len(sectorRefTask) > 0 {

						param := &pb.SectorsTask{
							SectorRefTask:     sectorRefTask,
							SectorType:        uint64(t.SectorType),
							SectorSize:        uint64(t.SectorSize),
							WorkerOpId:        v1.OpId,
							WorkerOpIp:        v1.Ip,
							StorageOpIp:       t.StorageOpIp,
							OriginalValueOpId: t.OriginalValueOpId,
							OriginalValueDir:  t.OriginalValueDir,
						}

						log.Println(fmt.Sprintf("DoSectorsTask param:%+v", param))
						_, err = client.RedoSectorsTask(context.TODO(), param)
						if err != nil {
							return err
						}
					}
				}
			}

			if t.TaskStatus == utils.ZERO {

				updateMap := make(map[string]interface{})
				updateMap["task_status"] = utils.ONE
				updateMap["begin_time"] = time.Now()

				err = global.ZC_DB.Model(&lotus.LotusSectorTask{}).Where("id", t.ID).Updates(updateMap).Error
				if err != nil {
					return err
				}
			}

			if t.TaskStatus != utils.THREE {

				runCount := utils.ZERO
				finishCount := utils.ZERO
				err = global.ZC_DB.Model(&lotus.LotusSectorTaskDetail{}).Where("sector_status <> 5 and relation_id=?", t.ID).Find(&lotusSectorDetail).Error
				if err != nil {
					return err
				}

				if len(lotusSectorDetail) > utils.ZERO {

					for _, v := range lotusSectorDetail {
						if v.SectorStatus == utils.ZERO || v.SectorStatus == utils.ONE || v.SectorStatus == utils.TWO || v.SectorStatus == utils.THREE {
							runCount++
						} else if v.SectorStatus == utils.FOUR {
							finishCount++
						}
					}

					modifyCheck := false
					updateMap := make(map[string]interface{})
					if runCount > utils.ZERO && runCount != t.RunCount {
						//更新正在跑的数量
						updateMap["run_count"] = runCount
						modifyCheck = true
					} else if finishCount > utils.ZERO && finishCount != t.FinishCount {
						//更新完成的数量
						updateMap["finish_count"] = finishCount
						modifyCheck = true
					} else if finishCount == t.SectorTotal {
						//更新任务状态
						updateMap["run_count"] = utils.ZERO
						updateMap["task_status"] = utils.THREE
						updateMap["end_time"] = time.Now()
						modifyCheck = true
					}
					if modifyCheck {
						err = global.ZC_DB.Model(&lotus.LotusSectorTask{}).Where("id", t.ID).Updates(updateMap).Error
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func (s *SectorsRecoverService) WorkerRedoTask(info *pb.SectorsTask) (bool, error) {

	sectorSizeInfo := fmt.Sprintf("%dGiB", info.SectorSize)
	sectorSizeInt, err := units.RAMInBytes(sectorSizeInfo)
	if err != nil {
		return false, err
	}

	switch info.SectorType {
	case define.CC.Uint64():
		cgo.RunSectorTaskToCc(info, uint64(sectorSizeInt))
		break
	case define.DC.Uint64():
		cgo.RunSectorTaskToDc(info, uint64(sectorSizeInt))
		break
	default:
		break
	}

	return true, nil
}

func GetDirList(pathName string, files *[]string, index int) error {

	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		return err
	}
	*files = append(*files, pathName)
	for _, fileInfo := range rd {
		if fileInfo.IsDir() {
			if index >= 2000 {
				break
			}
			index++
			_ = GetDirList(pathName+"/"+fileInfo.Name(), files, index)
		}
	}
	return err
}

func (s *SectorsRecoverService) ReadDirFile(info *pb.DirFileReq) (file []*pb.DirFile, err error) {

	res := make([]*pb.DirFile, 0)
	fileMap := make(map[string]int)

	var files []string
	err = GetDirList(info.Path, &files, utils.ZERO)
	if err != nil {
		return res, err
	}

	for _, v := range files {

		fileInfoList, err := ioutil.ReadDir(v)
		if err != nil {
			log.Println("ReadDir err:", err.Error())
			continue
		}

		for i := range fileInfoList {
			if utils.IsFile(v+"/"+fileInfoList[i].Name()) && strings.Contains(fileInfoList[i].Name(), ".car") {
				if n, ok := fileMap[v]; ok {
					fileMap[v] = n + 1
				} else {
					fileMap[v] = 1
				}
			}
		}
	}

	for k, v := range fileMap {

		t := &pb.DirFile{Path: k, FileNum: uint64(v)}
		res = append(res, t)
	}

	if len(res) == utils.ZERO {
		t := &pb.DirFile{Path: info.Path, FileNum: 0}
		res = append(res, t)
	}

	return res, nil
}

func (s *SectorsRecoverService) GetCarFilePath(info *pb.CarFile) (res *pb.CarFile, err error) {

	r := &pb.CarFile{}
	var files []string
	err = GetDirList(info.Path, &files, utils.ZERO)
	if err != nil {
		return r, err
	}

	checkPath := false
	for _, v := range files {

		if checkPath {
			break
		}
		fileInfoList, err := ioutil.ReadDir(v)
		if err != nil {
			log.Println("ReadDir err:", err.Error())
			continue
		}

		for i := range fileInfoList {
			if utils.IsFile(v+"/"+fileInfoList[i].Name()) && fileInfoList[i].Name() == info.FileName {
				r.Path = v
				r.FileName = fileInfoList[i].Name()
				r.OpId = info.OpId
				checkPath = true
				break
			}
		}
	}

	return r, nil
}

var lastPath string

func (s *SectorsRecoverService) GetOpFilePath(info *pb.OpFilePath) (string, error) {

	filePath := fmt.Sprintf("%s/%s", define.PathIpfsStorage, "storage.json")

	moveStorage := oplocal.Tasking.GetRunCount(define.MoveStorage.String())
	toPath := ""
	var err error
	if !utils.ExistFileOrDir(filePath) {
		//worker存储
		toPath, _, err = oplocal.StorageService{}.FindStorage(filePath, moveStorage, "store")
		if err != nil {
			return "", err
		}
	} else {
		//NFS存储
		disks := utils.GetOpDiskInfo()

		for i, disk := range disks {
			storagePath := filepath.Join(disk.Mounted, define.StoragePath)
			//是否可存储
			b, err := os.ReadFile(path.Join(storagePath, define.SectorStoreConfig))
			if err != nil {
				continue
			}

			var ss oplocal.LocalStorageMeta
			err = json.Unmarshal(b, &ss)
			if err != nil {
				return "", err
			}

			log.Printf("SectorStore:%+v", ss)
			if ss.CanStore {
				if utils.DiskSpaceSufficient(storagePath, define.Ss32GiB, 1) {
					if lastPath == storagePath && i < len(disks)-1 {
						continue
					}
					lastPath = storagePath
					toPath = storagePath
					break
				}
			}
		}
	}

	return toPath, nil
}

func (s *SectorsRecoverService) GetWorkerOpList(info request.WorkerInfo) (data interface{}, err error) {

	var lotusWorkerInfo []lotus.LotusWorkerInfo

	condition := " deploy_status=2 and run_status=1 "
	var param []interface{}
	if info.GateWayId != "" {
		condition += " and gate_id=? "
		param = append(param, info.GateWayId)
	}
	if info.WorkerType != "" {
		condition += " and worker_type=? "
		param = append(param, info.WorkerType)
	}
	if info.KeyWord != "" {
		info.KeyWord = "%" + info.KeyWord + "%"
		condition += " and ip like ? "
		param = append(param, info.KeyWord)
	}
	err = global.ZC_DB.Model(&lotus.LotusWorkerInfo{}).Where(condition, param...).Find(&lotusWorkerInfo).Error
	if err != nil {
		return lotusWorkerInfo, err
	}

	w := make([]response.WorkerOpList, 0)
	for _, v := range lotusWorkerInfo {
		w = append(w, response.WorkerOpList{OpId: v.OpId, Ip: v.Ip, Id: v.ID})
	}

	return w, nil
}

func (s *SectorsRecoverService) AddBadSector(info lotus.LotusSectorRecover) error {

	var total int64
	err := global.ZC_DB.Model(&lotus.LotusSectorRecover{}).Where("miner_id = ? and sector_id = ? and sector_status in(3,4)", info.MinerId, info.SectorId).Count(&total).Error
	if err != nil {
		return err
	}

	if total == 0 {

		r, err := SectorsRecoverServiceApi.GetSectorAr(info)
		if err != nil {
			return err
		}
		if len(r) > 0 {
			info.Ticket = r[0].Ticket
			info.SectorType = r[0].SectorType
			info.SectorSize = 32
		}
		if info.SectorType != 1 && info.SectorType != 2 {
			info.SectorType = 2
		}

		err = global.ZC_DB.Model(&lotus.LotusSectorRecover{}).Create(&info).Error
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *SectorsRecoverService) GetSectorAr(info lotus.LotusSectorRecover) (lf []lotus.LotusSectorInfo, e error) {

	var lsi []lotus.LotusSectorInfo
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorInfo{}.TableName(), info.MinerId)
	err := global.ZC_DB.Table(tableName).Where("sector_id = ?", info.SectorId).Find(&lsi).Error
	if err != nil {
		return lsi, err
	}
	return lsi, nil
}

func (s *SectorsRecoverService) GetBadSectorCount(info lotus.LotusSectorRecover) (int, error) {

	var total int64
	err := global.ZC_DB.Model(&lotus.LotusSectorAbnormal{}).Where("miner_id = ? and sector_id = ?", info.MinerId, info.SectorId).Count(&total).Error
	if err != nil {
		return utils.ZERO, err
	}

	if total == 0 {
		lb := &lotus.LotusSectorAbnormal{
			MinerId:      info.MinerId,
			SectorId:     info.SectorId,
			AbnormalTime: utils.TimeToFormat(time.Now(), utils.YearMonthDay),
			Count:        utils.ONE,
		}
		err = global.ZC_DB.Model(&lotus.LotusSectorAbnormal{}).Create(lb).Error
		if err != nil {
			return utils.ZERO, err
		}
		total = utils.ONE
	} else {
		err = global.ZC_DB.Model(&lotus.LotusSectorAbnormal{}).Where("miner_id = ? and sector_id = ?", info.MinerId, info.SectorId).Update("count", utils.TWO).Error
		if err != nil {
			return utils.ZERO, err
		}
		total = utils.TWO
	}

	return int(total), err
}

func (s *SectorsRecoverService) GetHostType(uuid string) (string, error) {

	var res system.SysHostRecord
	err := global.ZC_DB.Model(&system.SysHostRecord{}).Where("uuid", uuid).Find(&res).Error
	if err != nil {
		return "", err
	}

	if (res == system.SysHostRecord{}) {
		return "", errors.New(fmt.Sprintf("未找到主机:%s", uuid))
	}

	return strconv.Itoa(res.HostClassify), nil
}

func (s *SectorsRecoverService) GetCarFileList(info *pb.SectorID) ([]response.CarFile, error) {

	var res []response.CarFile
	tableName := fmt.Sprintf("%s_%s", lotus.LotusSectorPiece{}.TableName(), info.Miner)
	sql := fmt.Sprintf(" SELECT distinct r.piece_cid,r.piece_size FROM %s r WHERE r.actor=? AND r.sector_id=? AND r.deleted_at IS null", tableName)
	err := global.ZC_DB.Raw(sql, info.Miner, info.Number).Scan(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *SectorsRecoverService) ModifySectorTaskStatus(info request.SectorStatus) error {

	if info.Id == utils.ZERO || info.Status == utils.ZERO {
		return errors.New("id,status param is err")
	}

	if info.Status != define.Pause.Int() && info.Status != define.Stop.Int() && info.Status != define.Start.Int() {
		return errors.New("status data is err")
	}

	global.ZC_DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&lotus.LotusSectorTask{}).Where("id", info.Id).Update("task_status", info.Status).Error
		if err != nil {
			return err
		}

		taskDetailStatus := 0
		switch info.Status {
		case define.Pause.Int():
			taskDetailStatus = 6
			break
		case define.Stop.Int():
			taskDetailStatus = 7
			break
		case define.Start.Int():
			taskDetailStatus = -1
			break
		}

		err = tx.Model(&lotus.LotusSectorTaskDetail{}).Where("sector_status in(-1,0,6) and relation_id=?", info.Id).Update("sector_status", taskDetailStatus).Error
		if err != nil {
			return err
		}

		if info.Status == define.Stop.Int() {

			var lsd []lotus.LotusSectorTaskDetail
			err = tx.Model(&lotus.LotusSectorTaskDetail{}).Where("sector_status in(-1,0,6) and relation_id=?", info.Id).Find(&lsd).Error
			if err != nil {
				return err
			}

			var ids []int
			for _, v := range lsd {
				ids = append(ids, v.SectorRecoverId)
			}

			if len(ids) > utils.ZERO {
				err = tx.Model(&lotus.LotusSectorRecover{}).Where("id in (?)", ids).Update("sector_status", utils.THREE).Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return nil
}

func (s *SectorsRecoverService) GetCarFileParam(pieceCid string) (response.WorkerCarFiles, error) {

	var carFile response.WorkerCarFiles
	if utils.IsNull(pieceCid) {
		return carFile, fmt.Errorf("pieceCid is null")
	}

	err := global.ZC_DB.Table("worker_car_files").Where("piece_cid", pieceCid).Find(&carFile).Error
	if err != nil {
		return carFile, err
	}

	return carFile, nil
}
