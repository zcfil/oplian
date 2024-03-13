package cgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/service/lotus/oplocal"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// RunP1 run P1
func RunP1(param request.SectorRef, ticketPreimage []byte, piece []request.PieceInfo, taskDetailId, sectorRecoverId uint64) ([]byte, error) {

	log.Println(fmt.Sprintf("[%d] runP1 Running ...", param.ID.Number))
	pt := param.ProofType
	minerId := strconv.Itoa(int(param.ID.Miner))
	sid := &pb.SectorID{Miner: minerId, Number: param.ID.Number}

	sectorPath, err := oplocal.StorageService{}.FindSealStorage(define.PathIpfsWorker, &pb.SectorRef{Id: sid, ProofType: uint64(pt)})
	if err == nil {
		oplocal.SealSectorPath.Push(utils.SectorNumString(minerId, param.ID.Number), sectorPath.DiskPath)
	}

	RunTaskParam := &request.RunTaskParam{
		Miner:       param.ID.Miner,
		Number:      param.ID.Number,
		ProofType:   param.ProofType,
		Unsealed:    sectorPath.Unsealed,
		Sealed:      sectorPath.Sealed,
		Cache:       sectorPath.Cache,
		Update:      sectorPath.Update,
		UpdateCache: sectorPath.UpdateCache,
		Piece:       piece,
		Ticket:      ticketPreimage,
	}

	_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.PC1.Uint64()})
	if err != nil {
		return nil, err
	}

	checkP1Task := true
	go func() {

		_, err = utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), RunP1Task, "", RunTaskParam, time.Hour*6)
		if err != nil {
			checkP1Task = false
			_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Failed.Uint64()})
			if err != nil {
				return
			}
		}
	}()

	var out []byte
	for {

		time.Sleep(time.Minute)
		if !checkP1Task {

			return nil, fmt.Errorf("running P1 task failed:%v", param.ID)
		}
		res, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), GetP1TaskRes, "", RunTaskParam, time.Minute)
		if err != nil {

			_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Failed.Uint64()})
			if err != nil {
				return nil, err
			}
			return nil, err
		}

		var result response.P1ResData
		if err := json.Unmarshal(res, &result); err != nil {
			return nil, err
		}

		if result.Code != http.StatusOK {
			return nil, err
		}

		if result.Data.Code == 0 {
			continue
		}

		out = result.Data.Out
		break
	}

	return out, nil
}

// RunP2 run P2
func RunP2(param request.SectorRef, pc1o []byte, taskDetailId, sectorRecoverId uint64) error {

	pt := param.ProofType
	minerId := strconv.Itoa(int(param.ID.Miner))
	sid := &pb.SectorID{Miner: minerId, Number: param.ID.Number}

	sectorPath, err := oplocal.StorageService{}.FindSealStorage(define.PathIpfsWorker, &pb.SectorRef{Id: sid, ProofType: uint64(pt)})
	if err == nil {
		oplocal.SealSectorPath.Push(utils.SectorNumString(minerId, param.ID.Number), sectorPath.DiskPath)
	}

	RunTaskParam := &request.RunTaskParam{
		Miner:     param.ID.Miner,
		Number:    param.ID.Number,
		ProofType: param.ProofType,
		Phase1Out: pc1o,
		Cache:     sectorPath.Cache,
		Sealed:    sectorPath.Sealed,
	}

	_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.PC1.Uint64()})
	if err != nil {
		return err
	}

	checkP2Task := true
	go func() {

		_, err = utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), RunP2Task, "", RunTaskParam, time.Hour)
		if err != nil {
			_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Failed.Uint64()})
			if err != nil {
				return
			}
			checkP2Task = false
			return
		}
	}()

	for {

		time.Sleep(time.Minute)
		if !checkP2Task {
			return fmt.Errorf("failed to run the P2 task:%+v", param.ID)
		}
		res, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), GetP2TaskRes, "", RunTaskParam, time.Minute)
		if err != nil {
			_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Failed.Uint64()})
			if err != nil {
				return err
			}
			return err
		}

		var result response.P2ResData
		if err := json.Unmarshal(res, &result); err != nil {
			return err
		}

		if result.Code != http.StatusOK {
			return fmt.Errorf("p2 result failed %+v,%d", param.ID, result.Code)
		}

		if result.Data.Code == 0 {
			continue
		}

		break
	}

	return nil
}

// moveStorage Copy file
func moveStorage(opId, ip string, param request.SectorRef, flag int, taskDetailId, sectorRecoverId uint64) error {

	var si request.SectorRef
	si.ID = param.ID
	si.ProofType = param.ProofType
	minerId := strconv.Itoa(int(param.ID.Miner))
	sid := &pb.SectorID{Miner: minerId, Number: param.ID.Number}
	sectorPath, err := oplocal.StorageService{}.FindSealStorage(define.PathIpfsWorker, &pb.SectorRef{Id: sid, ProofType: uint64(param.ProofType)})

	RunTaskParam := &request.RunTaskParam{
		Cache:  sectorPath.Cache,
		Size:   si.ProofType,
		Number: param.ID.Number,
	}

	checkMoveTask := true
	go func() {

		_, err = utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), RunMoveTask, "", RunTaskParam, time.Hour)
		if err != nil {
			_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Failed.Uint64()})
			if err != nil {
				return
			}
			checkMoveTask = false
			return
		}
	}()

	for {

		time.Sleep(time.Minute)
		if !checkMoveTask {
			return fmt.Errorf("move task failed:%v", param.ID)
		}
		res, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), GetMoveTaskRes, "", RunTaskParam, time.Minute)
		if err != nil {
			_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Failed.Uint64()})
			if err != nil {
				return err
			}
			return err
		}

		var result response.MoveResData
		if err := json.Unmarshal(res, &result); err != nil {
			return err
		}

		if result.Code != http.StatusOK {
			return fmt.Errorf("move result failed %+v,%d", param.ID, result.Code)
		}

		if result.Data.Code == 0 {
			continue
		}

		break
	}

	conn, err := utils.GrpcConnect(ip, define.OpPort)
	if err != nil {
		return err
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	path, err := pb.NewOpServiceClient(conn).GetOpFilePath(context.TODO(), &pb.OpFilePath{SectorSize: uint64(si.ProofType)})
	if err != nil {
		return err
	}

	toPath := path.Value
	if toPath == "" {
		return errors.New("the storage path cannot be found on the storage machine")
	}

	for _, v := range fileTypes {
		switch flag & v {
		case UnSealed:

			unsealed := filepath.Join(sectorPath.Unsealed, "..")
			unsealedStorage := filepath.Join(toPath, "/unsealed")
			fileAr := make([]*pb.FileInfo, 0)
			fileAr = append(fileAr, &pb.FileInfo{
				FileName: SectorNumString(si.ID),
			})
			t := &pb.SynFileInfo{
				FromPath: unsealed,
				ToPath:   unsealedStorage,
				OpId:     opId,
				Ip:       ip,
				Port:     define.OpPort,
				FileInfo: fileAr,
			}

			res, _ := global.OpToGatewayClient.SysFilePoint(context.TODO(), t)
			if res != nil {
				if res.Code != 200 {
					return fmt.Errorf(res.Msg)
				}
			}

			os.Remove(sectorPath.Unsealed)

		case Cache:

			cache := filepath.Join(sectorPath.Cache, "..")
			cacheStorage := filepath.Join(toPath, "/cache")
			FileName := SectorNumString(si.ID)
			t := &pb.SynFileInfo{
				FromPath:    cache,
				ToPath:      cacheStorage,
				OpId:        opId,
				Ip:          ip,
				Port:        define.OpPort,
				ZipFileName: fmt.Sprintf("%s.zip", FileName),
			}

			res, _ := global.OpToGatewayClient.SysFilePoint(context.TODO(), t)
			if res != nil {
				if res.Code != 200 {
					return fmt.Errorf(res.Msg)
				}
			}

			os.Remove(sectorPath.Cache)

		case Sealed:

			sealed := filepath.Join(sectorPath.Sealed, "..")
			sealedStorage := filepath.Join(toPath, "/sealed")
			fileAr := make([]*pb.FileInfo, 0)
			fileAr = append(fileAr, &pb.FileInfo{
				FileName: SectorNumString(si.ID),
			})
			t := &pb.SynFileInfo{
				FromPath: sealed,
				ToPath:   sealedStorage,
				OpId:     opId,
				Ip:       ip,
				Port:     define.OpPort,
				FileInfo: fileAr,
			}

			res, _ := global.OpToGatewayClient.SysFilePoint(context.TODO(), t)
			if res != nil {
				if res.Code != 200 {
					return fmt.Errorf(res.Msg)
				}
			}

			os.Remove(sectorPath.Sealed)
		}
	}

	_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Finish.Uint64()})
	if err != nil {
		return err
	}

	return nil
}

func SectorNumString(id request.SectorID) string {
	return fmt.Sprintf("s-t0%d-%d", id.Miner, id.Number)
}

// ModifySectorStatus Update the sector status
func ModifySectorStatus(info *pb.TaskStatus) error {

	err := global.ZC_DB.Model(&lotus.LotusSectorTaskDetail{}).Where("id", info.TaskDetailId).Update("sector_status", info.Status).Error
	if err != nil {
		return err
	}

	//Task status 0 Wait,1 start,2ap,3p1,4p2,5 Completed,6 Recovery failed
	if strings.Contains("2,5,6", strconv.Itoa(int(info.Status))) {

		status := utils.ZERO
		switch info.Status {
		case define.AP.Uint64():
			status = int(define.Recovering.Uint64())
			break
		case define.Finish.Uint64():
			status = int(define.Recover.Uint64())
			break
		case define.Failed.Uint64():
			status = int(define.RecoverFailed.Uint64())
			break
		}

		upDateMap := make(map[string]interface{})
		upDateMap["sector_status"] = status
		if info.Status == define.Finish.Uint64() {
			upDateMap["recover_time"] = time.Now()
		}
		err = global.ZC_DB.Model(&lotus.LotusSectorRecover{}).Where("id", info.SectorRecoverId).Updates(upDateMap).Error
		if err != nil {
			return err
		}
	}

	return nil
}
