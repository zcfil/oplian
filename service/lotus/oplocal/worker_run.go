package oplocal

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path/filepath"
	"strings"
)

var WorkerRunServiceApi = new(WorkerRunService)

type WorkerRunService struct{}

// Ok
// @author: nathan
// @function: OK
// @description: Determine whether the task is available
// @param: task *pb.Task
// @return: bool
func (run WorkerRunService) Ok(task *pb.Task) bool {

	if OpWorkers.GetOpenWindow(task.Wid) == nil && task.Wid != define.RedoSectorsTask.String() {
		return false
	}
	Preselect.Lock.Lock()
	defer Preselect.Lock.Unlock()

	switch define.TaskType(task.TType) {
	case define.TTAddPiece:

		//Check whether the Car file exists
		if task.CarPath != "" {
			if _, err := os.Stat(task.CarPath); err != nil {
				//log.Println("CarFileExist Stat err:", err)
				return false
			}
		}

		//// Check whether the number of tasks being executed exceeds the set number
		curP1 := Tasking.GetRunCount(define.SealPreCommit1.String())
		if curP1 >= PreCount1 {
			log.Println("curP1 >= PreCount1:", curP1, PreCount1)
			return false
		}
		// Check whether the number of concurrent aps exceeds the threshold
		if Tasking.GetRunCount(define.AddPiece.String()) >= APCount*2 {
			log.Println("Tasking.GetRunCount(define.AddPiece.String()) >= APCount*2:", Tasking.GetRunCount(define.AddPiece.String()), APCount*2)
			return false
		}
		//Number of aps
		apCount, err := StorageService{}.FindSectorCount(define.PathIpfsWorker, define.FTUnsealed, task.MinerId)
		if err != nil {
			log.Println(err)
			return false
		}
		preCount, err := StorageService{}.FindSectorCount(define.PathIpfsWorker, define.FTCache, task.MinerId)
		if err != nil {
			log.Println(err)
			return false
		}
		_, meta, _ := StorageService{}.FindSectorFile(define.PathIpfsWorker, task.MinerId, task.SectorId, define.FTUnsealed)
		if !meta.CanSeal {
			// Number of AP,P1, and P2 sectors on the hard drive > Number of P1,P2 sectors on the hard drive + Number of allowed concurrent tasks + (Set number - Number of tasks that have been connected) Number of tasks that can be connected
			if apCount > preCount+APCount+PreCount1-curP1 {
				log.Println("apCount > preCount+APCount+PreCount1-curP1:", apCount, preCount, APCount, PreCount1, curP1)
				return false
			}
		}

		_, _, err = StorageService{}.FindStorage(define.PathIpfsWorker, 0, "seal")
		if err != nil {
			log.Println(err)
			return false
		}
		b := Preselect.Ok(define.TTAddPiece)
		log.Println("Whether to accept a task", define.TTAddPiece, b)
		return b
	case define.TTPreCommit1:
		//if task.Ip != OpWorkers.Info[task.Wid].Ip {
		//	return false
		//}
		curP1 := Tasking.GetRunCount(define.SealPreCommit1.String())
		//设置任务是否超过正在执行任务
		if curP1 >= PreCount1 {
			return false
		}

		spath, meta, err := StorageService{}.FindSectorFile(define.PathIpfsWorker, task.MinerId, task.SectorId, define.FTUnsealed)
		if err != nil {
			log.Println(err)
			return false
		}

		//Check whether the storage matches
		if !meta.CanSeal {
			return false
		}
		if !utils.DiskSpaceSufficient(spath, task.SectorSize, PreCount1*2) {
			return false
		}
		b := Preselect.Ok(define.TTPreCommit1)
		log.Println("Whether to accept a task", define.TTPreCommit1, b)
		return b
	case define.TTPreCommit2:
		//Set whether the number of tasks exceeds the number of tasks being performed
		if Tasking.GetRunCount(define.SealPreCommit2.String()) >= PreCount2 {
			return false
		}
		_, meta, err := StorageService{}.FindSectorFile(define.PathIpfsWorker, task.MinerId, task.SectorId, define.FTSealed)
		if err != nil {
			log.Println(err)
			return false
		}
		//Check whether the storage matches
		if !meta.CanSeal {
			return false
		}
		b := Preselect.Ok(define.TTPreCommit2)
		log.Println("Whether to accept a task", define.TTPreCommit2, b)
		return b
	case define.TTCommit1:
		_, meta, err := StorageService{}.FindSectorFile(define.PathIpfsWorker, task.MinerId, task.SectorId, define.FTSealed)
		if err != nil {
			log.Println(err)
			return false
		}
		//Check whether the storage matches
		if !meta.CanSeal {
			return false
		}
		log.Println("Receiving task", define.TTCommit1)
	case define.TTCommit2:
		window := OpWorkers.GetOpenWindow(task.Wid)
		if window == nil {
			return false
		}
		log.Println("Whether to accept a task", define.TTCommit2, window.RunC2)
		return window.RunC2
	case define.TTFetchLong:
		window := OpWorkers.GetOpenWindow(task.Wid)
		if window == nil {
			return false
		}
		colony, err := global.OpToGatewayClient.GetColony(context.Background(), &pb.Actor{MinerId: task.MinerId})
		if err != nil {
			log.Println("Failed to obtain cluster information. Procedure：", err)
			return false
		}
		//非worker存储，直接找文件在不在本地
		if colony.ColonyType != define.StorageTypeWorker {
			str, meta, err := StorageService{}.FindSectorFile(define.PathIpfsWorker, task.MinerId, task.SectorId, define.FTSealed)
			if err != nil {
				return false
			}
			if meta.CanSeal {
				return true
			}
			log.Println("NFS storage,", str)
		} else {
			//worker存储
			if !window.Storage {
				return false
			}
			moveStorage := Tasking.GetRunCount(define.MoveStorage.String())
			//判断硬盘大小
			path, _, err := StorageService{}.FindStorage(define.PathIpfsStorage, moveStorage, "store")
			if err != nil {
				log.Println(err)
				return false
			}
			log.Println("Find storage：", path)
		}
		return true
	case define.TTFinalize:
		_, meta, err := StorageService{}.FindSectorFile(define.PathIpfsWorker, task.MinerId, task.SectorId, define.FTSealed)
		if err != nil {
			return false
		}
		//Check whether the storage matches
		if !meta.CanSeal {
			return false
		}
		log.Println("Receiving task", define.TTFinalize)
	case define.TTFinalizeUnsealed:
		_, meta, err := StorageService{}.FindSectorFile(define.PathIpfsWorker, task.MinerId, task.SectorId, define.FTUnsealed)
		if err != nil {
			return false
		}
		//Check whether the storage matches
		if !meta.CanSeal {
			return false
		}
		log.Println("Receiving task", define.TTFinalizeUnsealed)
	}
	return true
}

// OkNew
// @author: nathan
// @function: OkNew
// @description: Determine whether the task is available
// @param: task *pb.Task
// @return: bool
func (run WorkerRunService) OkNew(miner *pb.MinerSize) *pb.TaskCan {
	var cantask pb.TaskCan
	Preselect.Lock.Lock()
	b1, b2 := Preselect.OkNew()
	Preselect.Lock.Unlock()
	if b1 {

		cantask.CanAp = int32(APCount*2 - Tasking.GetRunCount(define.AddPiece.String()))
		sectorType := make([]string, 0)
		sectorType = append(sectorType, define.FTUnsealed.String(), define.FTCache.String())
		sectorCount, sectorPath, err := StorageService{}.FindAllSector(define.PathIpfsWorker, sectorType, miner.Actor)
		apCount := sectorCount["apCount"]
		preCount := sectorCount["preCount"]
		//apCount, err := StorageService{}.FindSectorCount(define.PathIpfsWorker, define.FTUnsealed, miner.Actor)
		if err != nil {
			log.Println(err)
			cantask.CanAp = 0
			return &cantask
		}

		var apRedundancy = int32(APCount*5 + preCount - apCount)
		if cantask.CanAp > apRedundancy {
			cantask.CanAp = apRedundancy
		}

		cur1 := Tasking.GetRunCount(define.SealPreCommit1.String())
		cantask.CanP1 = int32(PreCount1) - int32(cur1)

		var sectorCan = make([]*pb.SectorPathCan, len(sectorPath))
		var i = 0

		for k, _ := range sectorPath {

			sectorCan[i] = &pb.SectorPathCan{CanCount: 0, Number: k}
			i++

		}

		if cantask.CanP1 > int32(DiskCount) {
			cantask.CanP1 = int32(DiskCount)
		}

		cantask.Sectors = sectorCan

		filepath.Walk(define.PathIpfsDataWorkerCar, func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if strings.HasPrefix(info.Name(), "baga6ea4seaq") {
				fileName := strings.Split(info.Name(), ".")
				if len(fileName) == 2 {
					cantask.Cars = append(cantask.Cars, fileName[0])
				}
			}
			return nil
		})

	}
	if b2 {
		cantask.CanP2 = int32(PreCount2 - Tasking.GetRunCount(define.SealPreCommit2.String()))
	}
	if cantask.CanP1 > 1 {
		cantask.CanP1 = 1
	}
	if cantask.CanP2 > 1 {
		cantask.CanP2 = 1
	}
	if cantask.CanAp > 1 {
		cantask.CanAp = 1
	}
	if cantask.CanAp > 0 || cantask.CanP1 > 0 || cantask.CanP2 > 0 {
		log.Println("start task：", cantask.CanAp, cantask.CanP1, cantask.CanP2, "硬盘剩余容量：", DiskCount)
	}
	//cantask.Cars

	return &cantask
}

// AddRunning
// @author: nathan
// @function: AddRunning
// @description: Add a record of executing a task to the OP
// @param: pb.Task
// @return: error
func (p WorkerRunService) AddRunning(args *pb.Task) error {
	if args.TType == define.SealPreCommit1.String() {

		if _, ok := IsDiskCount[args.SectorId]; !ok {
			DiskCount--
			IsDiskCount[args.SectorId] = struct{}{}
		}

		PathSealCount.Push(SealSectorPath.Get(utils.SectorNumString(args.MinerId, args.SectorId)), utils.SectorNumString(args.MinerId, args.SectorId), true)
	}
	Tasking.Push(args)
	return nil
}

// SubRunning
// @author: nathan
// @function: SubRunning
// @description: The task has been completed. The OP terminal registration has been revoked
// @param: pb.Task
// @return: error
func (p WorkerRunService) SubRunning(args *pb.Task) error {
	switch args.TType {
	case define.MoveStorage.String():
		delete(IsDiskCount, args.SectorId)
		DiskCount++
		PathSealCount.Sub(SealSectorPath.Get(utils.SectorNumString(args.MinerId, args.SectorId)), utils.SectorNumString(args.MinerId, args.SectorId))
		SealSectorPath.Remove(utils.SectorNumString(args.MinerId, args.SectorId))
	case define.SealPreCommit1.String():
		PathSealCount.Push(SealSectorPath.Get(utils.SectorNumString(args.MinerId, args.SectorId)), utils.SectorNumString(args.MinerId, args.SectorId), false)
	}
	Tasking.Remove(args)
	return nil
}

// ResetWorkerRunning
// @author: nathan
// @function: ResetWorkerRunning
// @description: Update an ongoing task from the worker
// @param: pb.Request
// @return: error
func (p WorkerRunService) ResetWorkerRunning(args *pb.WorkerTasks) error {

	if args.Info == nil {
		return fmt.Errorf("nil")
	}
	if args.Info.WorkerId == "00000000-0000-0000-0000-000000000000" || args.Info.Ip == "" {
		return fmt.Errorf("workerId：%s,ip：%s", args.Info.WorkerId, args.Info.Ip)
	}
	OpWorkers.Push(args.Info)

	NewRunning := make(map[string][]*pb.Task)
	for _, task := range args.Tasks {
		PathSealCount.Push(SealSectorPath.Get(utils.SectorNumString(task.MinerId, task.SectorId)), utils.SectorNumString(task.MinerId, task.SectorId), true)
		NewRunning[task.TType] = append(NewRunning[task.TType], task)
	}

	Tasking.Lock.Lock()
	Tasking.Run = NewRunning
	Tasking.Lock.Unlock()
	return nil
}
