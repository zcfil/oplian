package oplocal

import (
	"context"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/lotusrpc"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"time"
)

var waitClear bool
var StartClear time.Time

// SealingAbort Stop the task
func (p *OpLotusService) SealingAbort(args *pb.FilRestWorker) error {
	tasks := Tasking.GetRunArray()
	log.Println(tasks, "Start clearing tasks：", args.Host, args.Worker, args.NewHost)
	for _, tk := range tasks {
		switch tk.TType {
		case define.AddPiece.String(), define.SealPreCommit1.String():
			//判断扇区状态
			var actor, _ = strconv.ParseUint(tk.MinerId, 10, 64)
			log.Println("SealingAbort()", args.Host.Token, args.Host.Ip, actor, tk.SectorId, tk.Tid)
			if err := lotusrpc.FullApi.SealingAbort(args.Host.Token, args.Host.Ip, actor, tk.SectorId, tk.Tid); err != nil {
				log.Println("SealingAbort:", err.Error())
				continue
			}
			Tasking.Remove(tk)
		}
	}
	var errmsg error

	var remove = make(map[uint64]struct{})

	sectors := StorageService{}.RangeSectors(define.PathIpfsWorker, args.Host.Param)

	for _, number := range sectors {

		info, err := global.OpToGatewayClient.GetSectorStatus(context.Background(), &pb.SectorID{Miner: args.Host.Param, Number: number})
		if err != nil {
			errmsg = err
			log.Println("GetSectorStatus:", err.Error())
			continue
		}

		//All P2 were removed before
		if define.ReturnType(info.Status).BeforeP2(info.Status) {
			log.Println("SectorRemove()", args.Host.Token, args.Host.Ip, number)
			if err = lotusrpc.FullApi.SectorRemove(args.Host.Token, args.Host.Ip, number); err != nil {
				log.Println("SectorRemove:", err.Error())
			}
			remove[number] = struct{}{}
		}
	}
	if !waitClear && errmsg == nil {
		waitClear = true
		go func() {
			var minute = time.NewTicker(time.Second)
			var clear = time.NewTicker(time.Hour * 3)
			StartClear = time.Now()
			defer func() {
				waitClear = false
				StartClear = time.Time{}
			}()
			param := &pb.WorkerInfo{
				Id:         args.Worker.Id,
				Ip:         args.Worker.Ip,
				MinerToken: args.NewHost.Token,
				MinerIp:    args.NewHost.Ip,
			}
		isClear:
			for {
				select {
				case <-minute.C:
					sectors = StorageService{}.RangeSectors(define.PathIpfsWorker, args.Host.Param)
					leave := utils.ArrayNotRepeat(sectors, remove)
					log.Println("Waiting to clear cache, currently remaining sectors：", leave, "，The remaining time is forcibly cleared：", (time.Hour*3 - time.Since(StartClear)).String())
					if leave == 0 {
						// Start a new service
						if err := p.RunNewWorker(param); err != nil {
							log.Println("RunNewWorker error:", err.Error())
						}
						break isClear
					}
					minute = time.NewTicker(time.Minute)
				case <-clear.C:
					sectors = StorageService{}.RangeSectors(define.PathIpfsWorker, args.Host.Param)
					log.Println("To clear the cache, the current remaining sectors:", sectors, "，Start forcible clearing")
					for _, tk := range tasks {

						var actor, _ = strconv.ParseUint(tk.MinerId, 10, 64)
						log.Println("SealingAbort()", args.Host.Token, args.Host.Ip, actor, tk.SectorId, tk.Tid)
						if err := lotusrpc.FullApi.SealingAbort(args.Host.Token, args.Host.Ip, actor, tk.SectorId, tk.Tid); err != nil {
							log.Println("SealingAbort:", err.Error())
						}
						Tasking.Remove(tk)
						if err := lotusrpc.FullApi.SectorRemove(args.Host.Token, args.Host.Ip, tk.SectorId); err != nil {
							log.Println("SectorRemove:", err.Error())
						}
					}
					// Start a new service
					if err := p.RunNewWorker(param); err != nil {
						log.Println("RunNewWorker error:", err.Error())
					}
					break isClear
				}
			}
		}()
	} else {
		log.Println("waitClear:", waitClear, ",errmsg:", errmsg)
	}
	return nil
}
