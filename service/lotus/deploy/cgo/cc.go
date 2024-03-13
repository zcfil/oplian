package cgo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/service/lotus/oplocal"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"time"
)

const (
	UnSealed = 1 << iota
	Cache
	Sealed
)

var fileTypes = []int{UnSealed, Cache, Sealed}

// RunSectorTaskToCc run CC task
func RunSectorTaskToCc(param *pb.SectorsTask, size uint64) {

	go func(info *pb.SectorsTask) {

		preCommit1Run := make(chan request.P1Run, 1024)
		preCommit1Finish := make(chan struct{}, 1024)
		preCommit2Finish := make(chan struct{}, 1)
		preCommit2Run := make(chan request.P2Run, 1024)
		apPieceRun := make(chan request.SectorTicket, 1024)
		storageRun := make(chan request.SectorRef, 1024)

		//WaitDeals
		go func() {

			for _, st := range info.SectorRefTask {

				if string(st.Ticket) == "" {
					return
				}

				miner, err := utils.FileCoinStrToUint64(st.Id.Miner)
				if err != nil {
					return
				}

				s := request.SectorTicket{
					Sector: request.SectorRef{
						ID: request.SectorID{
							Miner:  miner,
							Number: st.Id.Number,
						},
						ProofType:       int64(size),
						TaskDetailId:    st.TaskDetailId,
						SectorRecoverId: st.SectorRecoverId,
					},
					Ticket: st.Ticket,
				}
				apPieceRun <- s
			}

		}()

		//AddPiece
		go func() {
			for {
				select {
				case sid := <-apPieceRun:

					minerId := strconv.Itoa(int(sid.Sector.ID.Miner))
					tk := pb.Task{TType: string(define.TTAddPiece), MinerId: minerId, SectorId: uint64(sid.Sector.ID.Number)}
					tk.TType = define.AddPiece.String()
					err := oplocal.WorkerRunServiceApi.AddRunning(&tk)
					if err != nil {
						continue
					}

					go func(ticket request.SectorTicket) {

						defer oplocal.WorkerRunServiceApi.SubRunning(&tk)

						var sRef request.SectorRef
						sRef.ID = sid.Sector.ID
						sRef.ProofType = sid.Sector.ProofType
						p1info, err := runAp(sRef, size, sid.Sector.TaskDetailId, sid.Sector.SectorRecoverId)
						if err != nil {
							return
						}
						p1 := request.P1Run{
							PieceInfo:       p1info,
							St:              ticket,
							TaskDetailId:    sid.Sector.TaskDetailId,
							SectorRecoverId: sid.Sector.SectorRecoverId,
						}
						preCommit1Run <- p1

					}(sid)
				}
			}
		}()

		//PreCommit1
		go func() {

			preCommit1Finish <- struct{}{}
			for {
				select {
				case p1 := <-preCommit1Run:
					go func(run request.P1Run) {

						minerId := strconv.Itoa(int(p1.St.Sector.ID.Miner))
						tk := pb.Task{Wid: define.RedoSectorsTask.String(), TType: string(define.TTPreCommit1), MinerId: minerId, SectorId: uint64(p1.St.Sector.ID.Number)}
						for {
							checkTask := oplocal.WorkerRunServiceApi.Ok(&tk)

							if !checkTask {

								time.Sleep(time.Minute)
								continue
							}
							break
						}
						tk.TType = define.SealPreCommit1.String()
						oplocal.WorkerRunServiceApi.AddRunning(&tk)

						p1out, err := RunP1(run.St.Sector, run.St.Ticket, run.PieceInfo, p1.TaskDetailId, p1.SectorRecoverId)
						defer func() {
							<-preCommit1Finish
							oplocal.WorkerRunServiceApi.SubRunning(&tk)
						}()
						if err != nil {
							log.Println(fmt.Sprintf("P1 error: %s Sector: %+v", err.Error(), run.St.Sector.ID))
							return
						}
						p2 := request.P2Run{
							PreCommit1Out:   p1out,
							Sector:          p1.St.Sector,
							TaskDetailId:    p1.TaskDetailId,
							SectorRecoverId: p1.SectorRecoverId,
						}
						preCommit2Run <- p2

					}(p1)
				}
			}
		}()

		//PreCommit2
		go func() {
			preCommit2Finish <- struct{}{}
			for {
				select {
				case p2 := <-preCommit2Run:

					minerId := strconv.Itoa(int(p2.Sector.ID.Miner))
					tk := pb.Task{Wid: define.RedoSectorsTask.String(), TType: string(define.TTPreCommit2), MinerId: minerId, SectorId: uint64(p2.Sector.ID.Number)}
					for {

						checkTask := oplocal.WorkerRunServiceApi.Ok(&tk)
						if !checkTask {
							time.Sleep(time.Minute)
							continue
						}
						break
					}
					tk.TType = define.SealPreCommit2.String()
					oplocal.WorkerRunServiceApi.AddRunning(&tk)

					err := RunP2(p2.Sector, p2.PreCommit1Out, p2.TaskDetailId, p2.SectorRecoverId)
					<-preCommit2Finish
					oplocal.WorkerRunServiceApi.SubRunning(&tk)

					if err != nil {
						log.Println(fmt.Sprintf("P2 error：%s,SectorRef: %+v ", err.Error(), p2.Sector.ID))
						continue
					}
					storageRun <- p2.Sector

				}
			}
		}()

		//FinalizeSector
		for {
			select {
			case si := <-storageRun:
				go func() {

					minerId := strconv.Itoa(int(si.ID.Miner))
					tk := pb.Task{TType: string(define.MoveStorage), MinerId: minerId, SectorId: uint64(si.ID.Number)}
					oplocal.WorkerRunServiceApi.AddRunning(&tk)
					defer oplocal.WorkerRunServiceApi.SubRunning(&tk)
					if err := moveStorage(info.WorkerOpId, info.StorageOpIp, si, UnSealed|Cache|Sealed, si.TaskDetailId, si.SectorRecoverId); err != nil {
						return
					}

				}()
			}
		}

	}(param)

}

// runAp run AP
func runAp(sector request.SectorRef, pieceSize uint64, taskDetailId, sectorRecoverId uint64) ([]request.PieceInfo, error) {

	RunTaskParam := &request.RunTaskParam{
		Miner:     sector.ID.Miner,
		Number:    sector.ID.Number,
		ProofType: sector.ProofType,
		PieceSize: int(pieceSize),
	}

	var pieces []request.PieceInfo
	checkApTask := true
	go func() {

		_, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), RunCcApTask, "", RunTaskParam, 0)
		if err != nil {
			checkApTask = false
			_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.Failed.Uint64()})
			if err != nil {
				checkApTask = false
				return
			}
		}
	}()

	for {

		time.Sleep(time.Minute)
		if !checkApTask {
			return pieces, fmt.Errorf("the CC failed to run the AP task error:%v", sector.ID)
		}
		res, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), GetCcApTaskRes, "", RunTaskParam, 0)
		if err != nil {
			return pieces, err
		}

		var result response.ApResData
		err = json.Unmarshal(res, &result)
		if err != nil {
			return pieces, err
		}

		if result.Code != http.StatusOK {
			return pieces, fmt.Errorf("the CC failed to run the AP task. Procedure %+v,%d", sector.ID, result.Code)
		}

		if result.Data.Code == 0 {
			continue
		}

		pieces = append(pieces, result.Data.PInfo...)
		break
	}

	_, err := global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.AP.Uint64()})
	if err != nil {
		return nil, err
	}

	return pieces, nil
}
