package cgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/service/lotus/deploy/cgo/car/tools"
	"oplian/service/lotus/deploy/cgo/car/util"
	"oplian/service/lotus/oplocal"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LocalUrl       = "127.0.0.1"
	RunCcApTask    = "/runCcApTask"
	GetCcApTaskRes = "/getCcApTaskRes"
	RunDcApTask    = "/runDcApTask"
	GetDcApTaskRes = "/getDcApTaskRes"
	RunP1Task      = "/runP1Task"
	GetP1TaskRes   = "/getP1TaskRes"
	RunP2Task      = "/runP2Task"
	GetP2TaskRes   = "/getP2TaskRes"
	RunMoveTask    = "/runMoveTask"
	GetMoveTaskRes = "/getMoveTaskRes"
	Port           = "50059"
)

func RunSectorTaskToDc(param *pb.SectorsTask, size uint64) {

	go func(info *pb.SectorsTask) {

		preCommit1Run := make(chan request.P1Run, 1024)
		preCommit1Finish := make(chan struct{}, 1024)
		preCommit2Finish := make(chan struct{}, 1)
		preCommit2Run := make(chan request.P2Run, 1024)
		apPieceRun := make(chan request.SectorTicket, 1024)
		storageRun := make(chan request.SectorRef, 1024)

		if err := utils.AddDir(define.PathIpfsData + "/car"); err != nil {
			log.Println("RunSectorTaskToDc 创建目录失败:", define.PathIpfsData+"/car")
		}

		carTotal := 0
		//拷贝原值文件
		for i, st := range info.SectorRefTask {

			carAr, err := global.OpToGatewayClient.CarFileList(context.TODO(), &pb.SectorID{Miner: st.Id.Miner, Number: st.Id.Number})
			if err != nil {
				break
			}

			if len(carAr.CarInfo) == utils.ZERO {
				continue
			}

			param.SectorRefTask[i].PieceCid = carAr.CarInfo[0].PieceCid
			param.SectorRefTask[i].PieceSize = carAr.CarInfo[0].PieceSize

			//获取car文件路径
			for _, v := range carAr.CarInfo {

				p, err := global.OpToGatewayClient.CarFilePath(context.TODO(), &pb.CarFile{OpId: param.OriginalValueOpId, FileName: v.FileName, Path: param.OriginalValueDir})
				if err != nil && strings.Contains(err.Error(), "opClient Connection failed") {
					continue
				}

				addCar := true
				if p != nil {
					if p.Path != "" {
						addCar = false
					}
				}

				//拷贝
				if !addCar {

					log.Println(fmt.Sprintf("Car %s.car 进入服务器查找...", carAr.CarInfo[0].PieceCid))
					//复制car文件到worker机
					fileAr := make([]*pb.FileInfo, 0)
					fileAr = append(fileAr, &pb.FileInfo{
						FileName: v.FileName,
					})
					t := &pb.SynFileInfo{
						FromPath: p.Path,
						ToPath:   define.PathIpfsData + "/car",
						OpId:     param.OriginalValueOpId,
						Ip:       param.WorkerOpIp,
						Port:     define.OpPort,
						FileInfo: fileAr,
					}

					fmt.Printf("SynFileInfo t:%+v", t)
					res, err := global.OpToGatewayClient.SysFilePoint(context.TODO(), t)
					if err != nil || res.Code != 200 {
						break
					}
					carTotal++

				} else {

					log.Println(fmt.Sprintf("Car %s.car 进入重新生成...", carAr.CarInfo[0].PieceCid))
					carFileParam, err := global.OpToGatewayClient.CarFileParam(context.TODO(), &pb.String{Value: carAr.CarInfo[0].PieceCid})
					if err != nil {
						break
					}

					if carFileParam.PieceCid != "" {

						log.Println(fmt.Sprintf("%s.car 找到生成参数...", carAr.CarInfo[0].PieceCid))
						os.Remove(path.Join(define.PathIpfsData+"/car", v.FileName))
						err = CreateCarFile(carFileParam)
						if err != nil {
							break
						}

						for {

							time.Sleep(time.Second * 30)
							p, err := global.OpToGatewayClient.CarFilePath(context.TODO(), &pb.CarFile{OpId: global.OpUUID.String(), FileName: v.FileName, Path: define.PathIpfsData + "/car"})
							if err != nil {
								continue
							}
							if p != nil {
								if p.Path != "" {
									carTotal++
									break
								}
							}
						}
					}
				}
			}
		}

		if carTotal > 0 {

			//WaitDeals
			go func() {

				for _, st := range info.SectorRefTask {
					if string(st.Ticket) == "" {
						log.Println("ticket err：", st.Id.Number)
						return
					}

					miner, err := utils.FileCoinStrToUint64(st.Id.Miner)
					if err != nil {
						log.Println("miner err：", st.Id.Miner)
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
						Ticket:    st.Ticket,
						PieceCid:  st.PieceCid,
						PieceSize: int(st.PieceSize),
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
						err := oplocal.WorkerRunServiceApi.AddRunning(&tk)
						if err != nil {
							return
						}

						go func(ticket request.SectorTicket) {

							defer oplocal.WorkerRunServiceApi.SubRunning(&tk)
							p1info, err := addPiece(ticket.Sector, ticket.PieceCid, ticket.PieceSize, sid.Sector.TaskDetailId, sid.Sector.SectorRecoverId)
							if err != nil {
								log.Println("AP error：", err.Error(), ",", ticket.Sector.ID)
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
								oplocal.WorkerRunServiceApi.SubRunning(&tk)
								<-preCommit1Finish
							}()
							if err != nil {
								log.Println(fmt.Sprintf("P1 error: %s SectorRef: %+v", err.Error(), run.St.Sector.ID))
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
				for {
					preCommit2Finish <- struct{}{}
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
						if err := RunP2(p2.Sector, p2.PreCommit1Out, p2.TaskDetailId, p2.SectorRecoverId); err != nil {
							log.Println(fmt.Sprintf("P2 error：%s,SectorRef: %+v ", err.Error(), p2.Sector.ID))
							continue
						}
						<-preCommit2Finish
						oplocal.WorkerRunServiceApi.SubRunning(&tk)

						storageRun <- p2.Sector

					}
				}
			}()

			//Finish
			for {
				select {
				case si := <-storageRun:
					go func() {

						minerId := strconv.Itoa(int(si.ID.Miner))
						tk := pb.Task{TType: string(define.MoveStorage), MinerId: minerId, SectorId: uint64(si.ID.Number)}
						oplocal.WorkerRunServiceApi.AddRunning(&tk)
						defer oplocal.WorkerRunServiceApi.SubRunning(&tk)
						err := moveStorage(info.WorkerOpId, info.StorageOpIp, si, UnSealed|Cache|Sealed, si.TaskDetailId, si.SectorRecoverId)
						if err != nil {
							return
						}
					}()
				}
			}

		} else {

			if len(info.SectorRefTask) > 0 {

				_, err := global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: info.SectorRefTask[0].TaskDetailId, SectorRecoverId: info.SectorRefTask[0].SectorRecoverId, Status: define.Failed.Uint64()})
				if err != nil {
					return
				}
			}
		}

	}(param)

}

func CreateCarFile(info *pb.CarInfo) error {

	f, err := ioutil.ReadDir(info.InPutDir)
	if err != nil {
		return err
	}

	if len(f) == 0 {
		return fmt.Errorf("the source file directory file is empty")
	}

	resCarFile := make([]*pb.CarInfo, 0)
	resCarFile = append(resCarFile, info)

	var wg sync.WaitGroup
	limiter := make(chan struct{}, 1)
	for _, v := range resCarFile {

		wg.Add(1)
		limiter <- struct{}{}
		var cars []util.Finfo
		var ff response.Files
		if v.FileStr != "" {
			err := json.Unmarshal([]byte(v.FileStr), &ff)
			if err != nil {
				if err != nil {
					return err
				}
			}

			for _, v1 := range ff.Files {

				fileStrAr := strings.Split(v1.FileStr, "|")
				for _, v2 := range fileStrAr {

					v2Ar := strings.Split(v2, "=")
					indexStrAr := strings.Split(v2Ar[1], ",")
					s, _ := strconv.Atoi(indexStrAr[0])
					e, _ := strconv.Atoi(indexStrAr[1])
					cars = append(cars, util.Finfo{

						Name:  v1.FilePath,
						Path:  v1.FilePath,
						Size:  v1.FileSize,
						Start: int64(s),
						End:   int64(e),
					})
				}
			}
		}

		go func(record *pb.CarInfo, files []util.Finfo) {
			defer func() {
				<-limiter
				wg.Done()
			}()

			car, err := tools.GenerateCar(context.TODO(), define.PathIpfsData+"/car", "", record.InPutDir, files)
			if err != nil {
				return
			}

			if record.FileName == car.PieceCid+".car" && record.PieceCid == car.PieceCid &&
				record.CarSize == car.CarSize && record.PieceSize == int64(car.PieceSize) && record.DataCid == car.DataCid {

				log.Println(fmt.Sprintf("The car file is successfully recreated. Procedure:%s", car.PieceCid+".car"))
			} else {
				log.Println(fmt.Sprintf("The car file is recreated abnormally:%s", record.FileName))
			}

		}(v, cars)

	}

	wg.Wait()

	return nil
}

// run AP
func addPiece(param request.SectorRef, pieceCid string, pieceSize int, taskDetailId, sectorRecoverId uint64) ([]request.PieceInfo, error) {

	var pieces []request.PieceInfo
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
		PieceCid:    pieceCid,
		PieceSize:   pieceSize,
		CarPath:     define.PathIpfsData,
		Unsealed:    sectorPath.Unsealed,
		Sealed:      sectorPath.Sealed,
		Cache:       sectorPath.Cache,
		Update:      sectorPath.Update,
		UpdateCache: sectorPath.UpdateCache,
	}

	_, err = global.OpToGatewayClient.ModifySectorStatus(context.TODO(), &pb.TaskStatus{TaskDetailId: taskDetailId, SectorRecoverId: sectorRecoverId, Status: define.AP.Uint64()})
	if err != nil {
		return pieces, err
	}

	checkApTask := true
	go func() {

		_, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), RunDcApTask, "", RunTaskParam, time.Hour)
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
			return pieces, fmt.Errorf("the DC runs AP tasks %s Request failed", RunDcApTask)
		}
		res, err := utils.RequestDo(fmt.Sprintf("%s:%s", LocalUrl, Port), GetDcApTaskRes, "", RunTaskParam, time.Minute)
		if err != nil {
			return pieces, err
		}

		var result response.ApResData
		err = json.Unmarshal(res, &result)
		if err != nil {
			return pieces, err
		}

		if result.Code != http.StatusOK {
			return pieces, err
		}

		if result.Data.Code == 0 {
			continue
		}

		pieces = append(pieces, result.Data.PInfo...)
		break
	}

	return pieces, nil

}
