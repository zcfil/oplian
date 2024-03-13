package commit2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"oplian/define"
	"oplian/global"
	"oplian/model/gateway"
	"oplian/service/lotus/deploy"
	"oplian/service/pb"
	"oplian/utils"
	"path"
	"sync"
	"time"
)

var WorkerInfoImpl = new(WorkerInfo)
var WorkersClientLock sync.RWMutex
var WorkersClient = make(map[string]*WorkerInfo)

// var WorkersRL sync.RWMutex
//var RunningRL sync.RWMutex
//var Running = make(map[gateway.SectorID]time.Time)
//var Commit2RL sync.RWMutex
//var Commit2 = make(map[gateway.SectorID]*gateway.ProofResult)

const (
	LocalUrl    = "127.0.0.1"
	GetC2Result = "/getRunC2Result"
	RunC2Task   = "/runC2Task"
	Port        = "50063"
)

type SealerC2Param struct {
	Phase1Out []byte
	Miner     uint64
	Number    uint64
}

type ResData struct {
	Code int      `json:"Code"`
	Data C2Result `json:"Data"`
}

type C2Result struct {
	Proof []byte `json:"Proof"`
	Code  int    `json:"Code"`
}

type WorkerInfo struct {
	GpuUse     bool
	Disconnect bool
	Host       string
	OpC2Id     string
	TimeOut    time.Time
	WorkerRL   sync.RWMutex
}

func SetWorkersClient(info *WorkerInfo) {
	WorkersClientLock.Lock()
	defer WorkersClientLock.Unlock()
	WorkersClient[info.OpC2Id] = NewWorkerInfo(info.Host, info.OpC2Id, info.GpuUse)
}

func DelWorkersClient(opC2Id string) {
	WorkersClientLock.Lock()
	defer WorkersClientLock.Unlock()
	delete(WorkersClient, opC2Id)
}

func GetWorkersClient(opC2Id string) (bool, *WorkerInfo) {
	WorkersClientLock.RLock()
	defer WorkersClientLock.RUnlock()
	if _, ok := WorkersClient[opC2Id]; ok {
		return true, WorkersClient[opC2Id]
	} else {
		return false, nil
	}
}

func GetWorkersClientList() []*WorkerInfo {
	WorkersClientLock.RLock()
	defer WorkersClientLock.RUnlock()

	var res []*WorkerInfo
	for _, v := range WorkersClient {
		res = append(res, v)
	}

	return res
}

func NewWorkerInfo(host, opC2Id string, gpuUse bool) *WorkerInfo {

	w := &WorkerInfo{
		GpuUse:     gpuUse,
		Disconnect: false,
		Host:       host,
		OpC2Id:     opC2Id,
		TimeOut:    time.Now(),
	}

	return w
}

func ImplRunCommit2(sp *pb.SealerParam) error {

	miner, err := utils.FileCoinStrToUint64(sp.Sector.Id.Miner)
	if err != nil {
		return err
	}

	data := &gateway.SealerParam{
		ID: struct {
			Miner  uint64
			Number uint64
		}{Miner: miner, Number: sp.Sector.Id.Number},
		Phase1Out: sp.Phase1Out,
	}

	err = SealCommitPhase2(*data, sp.OpC2Id, sp.Host, sp.SealPort)
	if err != nil {
		return err
	}

	return nil
}

func SealCommitPhase2(args gateway.SealerParam, opId, host, sealPort string) error {

	go func(str, str0, port string) {

		start := time.Now()
		opFilePathC1 := "/ipfs/data/c2task" + define.OpCsPathC1
		gateWayFilePathC2 := "/ipfs/data/c2task" + define.GateWayCsPathC2

		//跑C2任务
		var proof []byte
		sp := &SealerC2Param{
			Miner:     args.ID.Miner,
			Number:    args.ID.Number,
			Phase1Out: args.Phase1Out,
		}

		checkC2Task := true
		go func() {
			_, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, port), RunC2Task, "", sp, time.Minute*10)
			if err != nil {
				checkC2Task = false
				log.Println("SealCommitPhase2 C2任务请求失败:", sp.Miner, sp.Number, err)
				return
			}
		}()

		for {

			time.Sleep(time.Second * 5)
			if !checkC2Task {
				break
			}
			res, err := GetRunC2Result(sp, port)
			if err != nil || res.Code != http.StatusOK {
				break
			} else {
				if res.Data.Code == 0 {
					continue
				} else {
					proof = res.Data.Proof
					break
				}
			}
		}

		dur := time.Now().Sub(start)
		if len(proof) > 0 {
			log.Println(fmt.Sprintf("SealCommitPhase2 sector :%d,C2 Mission success", sp.Number))
		} else {
			log.Println(fmt.Sprintf("SealCommitPhase2 sector :%d,C2 task failed", sp.Number))

			_, err := global.OpC2ToOp.CompleteCommit2(context.TODO(), &pb.FileInfo{Miner: args.ID.Miner, Number: args.ID.Number, OpId: str, Host: str0, TaskStatus: utils.TWO})
			if err != nil {
				log.Println("SealCommitPhase2 CompleteCommit2 err:", err.Error())
				return
			}
			return
		}

		for {

			global.OpC2ToOp.DelGateWayFile(context.TODO(), &pb.FileInfo{FileName: path.Join(gateWayFilePathC2, fmt.Sprintf("s-t0%d-%d", args.ID.Miner, args.ID.Number))})
			t := &pb.FileInfo{Path: gateWayFilePathC2, FileName: fmt.Sprintf("s-t0%d-%d", args.ID.Miner, args.ID.Number), FileData: proof}
			_, err := deploy.WorkerClusterServiceApi.FileToByteGateWay(t)
			if err != nil {
				log.Println("SealCommitPhase2 FileToByteOp err:", err.Error())
				time.Sleep(time.Second * 5)
				continue
			} else {

				fileName := path.Join(opFilePathC1, fmt.Sprintf("t0%d-%d.json", args.ID.Miner, args.ID.Number))
				log.Println("SealCommitPhase2 Remove the task source file:", fileName)
				_, err = global.OpC2ToOp.DelOpFile(context.TODO(), &pb.FileInfo{FileName: fileName})
				if err != nil {
					log.Println("SealCommitPhase2 CreateOpFile err:", err.Error())
					time.Sleep(time.Second * 5)
					continue
				}

				_, err = global.OpC2ToOp.CompleteCommit2(context.TODO(), &pb.FileInfo{Miner: args.ID.Miner, Number: args.ID.Number, OpId: str, Host: str0, TimeLength: dur.String(), TaskStatus: utils.ONE})
				if err != nil {
					log.Println("SealCommitPhase2 CompleteCommit2 err:", err.Error())
					time.Sleep(time.Second * 5)
					continue
				}

				break
			}
		}

	}(opId, host, sealPort)

	return nil
}

func GetRunC2Result(sp *SealerC2Param, port string) (c2 ResData, er error) {

	res, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, port), GetC2Result, "", sp, 15*time.Second)
	if err != nil {
		fmt.Println("RequestDo err:", err)
	}

	var result ResData
	err = json.Unmarshal(res, &result)
	if err != nil {
		fmt.Println("Unmarshal err:", err)
		return result, err
	}

	return result, nil

}
