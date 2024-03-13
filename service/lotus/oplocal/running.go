package oplocal

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"oplian/define"
	"oplian/global"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

type Running struct {
	Run  map[string][]*pb.Task
	Lock sync.RWMutex
}

func (r *Running) GetRunCount(ss string) int {
	r.Lock.RLock()
	defer r.Lock.RUnlock()
	return len(r.Run[ss])
}

// GetRunList worker task list
func (r *Running) GetRunList(tt string) *pb.TaskInfoList {

	r.Lock.RLock()
	defer r.Lock.RUnlock()
	ts := r.Run[tt]
	infos := make([]*pb.TaskInfo, len(ts))

	apSectorMap := make(map[string]uint64)
	pSectorMap := make(map[string]uint64)
	//log.Println(fmt.Sprintf("GetRunList ts:%+v", ts))
	for _, v := range ts {
		switch v.TType {
		case define.AddPiece.String():
			apSectorMap[utils.SectorNumString(v.MinerId, v.SectorId)] = v.SectorId
			break
		case define.SealPreCommit1.String(), define.SealPreCommit2.String():
			pSectorMap[utils.SectorNumString(v.MinerId, v.SectorId)] = v.SectorId
			break
		}
	}

	var err error
	apSectorPath := make(map[uint64]string)
	pSectorPath := make(map[uint64]string)
	if len(apSectorMap) > 0 {
		apSectorPath, err = StorageService{}.FindSectorsPath(define.PathIpfsWorker, apSectorMap, define.FTUnsealed)
		if err != nil {
			log.Println("apSectorPath FindSectorPath err:", err)
		}
		log.Println("GetRunList apSectorPath:", len(apSectorPath))
	}
	if len(pSectorMap) > 0 {

		pSectorPath, err = StorageService{}.FindSectorsPath(define.PathIpfsWorker, pSectorMap, define.FTCache)
		if err != nil {
			log.Println("p1SectorPath FindSectorPath err:", err)
		}
		log.Println("GetRunList p1SectorPath:", len(pSectorPath))
	}

	var wait sync.WaitGroup
	for k, v := range ts {
		wait.Add(1)
		go func(i int, task *pb.Task) {
			defer wait.Done()
			// Get consumption time
			timeLength := time.Since(task.StartTime.AsTime())
			// Get the progress
			progress := 0
			is64G := task.SectorSize > define.Ss32GiB
			switch task.TType {
			case define.AddPiece.String():

				spath, ok := apSectorPath[task.SectorId]
				if !ok {
					log.Println(fmt.Sprintf("GetRunList %d,ap找不到扇区", task.SectorId))
				}
				file, err := os.Stat(spath)
				if err != nil {
					log.Println("GetRunList Stat:", err, spath)
				}
				if file != nil {
					progress, _ = strconv.Atoi(fmt.Sprintf("%.0f", float64(file.Size())/float64(task.SectorSize)*100))
				}

			case define.SealPreCommit1.String():

				spath, ok := pSectorPath[task.SectorId]
				if !ok {
					log.Println(fmt.Sprintf("GetRunList %d,p1找不到扇区", task.SectorId))
				}
				progress = define.CacheFile.ProgressP1(spath, task.SectorSize)

			case define.SealPreCommit2.String():

				var trees []string
				var ssize = define.CacheFile.TreeCSize()
				var rlastSize = define.CacheFile.TreeRLastSize()
				if is64G {
					trees = define.CacheFile.Trees64G()
				} else {
					trees = define.CacheFile.Trees32G()
				}

				spath, ok := pSectorPath[task.SectorId]
				if !ok {
					log.Println(fmt.Sprintf("GetRunList %d,p2找不到扇区", task.SectorId))
				}
				for i, treeName := range trees {
					tpath := path.Join(spath, treeName)
					file, err := os.Stat(tpath)
					if err != nil {
						log.Println("treeName:", err)
						break
					}
					if i >= len(trees)/2 {
						ssize = rlastSize
					}
					// Complete one layer
					if uint64(file.Size()) >= ssize {
						progress += 100 / len(trees)
						continue
					}
					p, _ := strconv.Atoi(fmt.Sprintf("%.0f", float64(file.Size())*100/float64(ssize)))
					progress += p
					break
				}

				//C1 C2 cannot be judged by the file, but can only be estimated
			case define.SealCommit1.String():
				progress = rand.Intn(99)
			case define.SealCommit2.String():
				//15 minutes is the baseline
				progress, _ = strconv.Atoi(fmt.Sprintf("%.0f", timeLength.Seconds()*100/(time.Minute.Seconds()*5)))
				progress = progress * 110 / 100
			}
			if progress >= 100 {
				progress = 99
			}
			info := &pb.TaskInfo{
				Ip:         task.Ip,
				OpId:       global.OpUUID.String(),
				TType:      task.TType,
				Progress:   int32(progress),
				MinerId:    task.MinerId,
				TimeLength: utils.WholeSecond(timeLength),
				SectorId:   task.SectorId,
			}
			infos[i] = info
		}(k, v)
	}
	wait.Wait()
	return &pb.TaskInfoList{Tasks: infos}
}
func (r *Running) Push(task *pb.Task) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	task.StartTime = timestamppb.New(time.Now())
	Tasking.Run[task.TType] = append(Tasking.Run[task.TType], task)
	return
}

// Remove Remove target task
func (r *Running) Remove(task *pb.Task) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	var NewRunning []*pb.Task
	for _, v := range Tasking.Run[task.TType] {

		if v.Tid == task.Tid {
			continue
		}
		NewRunning = append(NewRunning, v)
	}

	Tasking.Run[task.TType] = NewRunning
	return
}

// GetRunArray worker task list
func (r *Running) GetRunArray() []*pb.Task {
	r.Lock.RLock()
	defer r.Lock.RUnlock()
	var res []*pb.Task
	for _, v := range r.Run {
		res = append(res, v...)
	}

	return res
}
