package oplocal

import (
	"context"
	"fmt"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/lotusrpc"
	"oplian/service/pb"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MinerService struct{}

// NewPledgeTask
// @author: nathan
// @function: NewPledgeTask
// @description: Create a new pledge task
// @param: ctx context.Context, actor string
// @return: error
func (m *MinerService) NewPledgeTask(ctx context.Context, actor string, info *pb.OpenWindow) error {

	if info != nil {
		OpWorkers.Push(info)
	}
	tasks, err := global.OpToGatewayClient.GetActorTaskQueue(ctx, &pb.String{Value: actor})
	if err != nil {
		return err
	}

	for _, task := range tasks.Queues {

		if task.RunCount+task.CompleteCount >= task.JobTotal {
			continue
		}
		// Set the AP or the number of concurrent order imports
		APCount = int(task.ConcurrentImport)

		switch task.SectorType {
		case define.SectorTypeCC:
			if WaitRecord {
				log.Println("The task failed！")
				return nil
			}
			sector, err := lotusrpc.FullApi.PledgeSector(task.MinerToken, task.MinerIp)
			log.Println("Release task:", sector, err)
			if err != nil {
				return err
			}
			// Prevent duplicate creation
			WaitRecord = true
			go func() {
				defer func() {
					WaitRecord = false
				}()
				// Write must be guaranteed
				for {
					if _, err = global.OpToGatewayClient.AddRunCountByID(context.Background(), &pb.TaskQueue{ID: task.ID}); err != nil {
						global.ZC_LOG.Error("AddRunCountByID：" + err.Error())
						time.Sleep(time.Second)
						continue
					}
					break
				}
				// Write must be guaranteed
				for {
					detail := &pb.SectorQueueDetail{
						Sid:          &pb.SectorID{Number: sector.Number, Miner: actor},
						QueueId:      task.ID,
						SectorStatus: define.QueueSectorStatusCreate,
					}
					if _, err = global.OpToGatewayClient.AddSectorQueueDetail(context.Background(), detail); err != nil {
						global.ZC_LOG.Error("AddSectorQueueDetail：" + err.Error())
						time.Sleep(time.Second)
						continue
					}
					break
				}
			}()
		case define.SectorTypeDC:
			// Get order data
			deals, err := global.OpToGatewayClient.GetWaitImportDeal(ctx, &pb.DealParam{Actor: task.Actor, QueueId: task.ID, Count: task.ConcurrentImport})
			if err != nil {
				global.ZC_LOG.Error("GetWaitImportDeal：" + err.Error())
				return err
			}

			if len(deals.Deals) == 0 {
				continue
			}
			fileAr := make([]*pb.FileInfo, 1)
			var dealWait sync.WaitGroup
			dealWait.Add(len(deals.Deals))
			for _, v := range deals.Deals {
				// Concurrent import
				go func(deal *pb.DealInfo) {
					defer dealWait.Done()
					fileAr[0] = &pb.FileInfo{FileName: deal.PieceCid + ".car"}

					// File copy
					in := pb.SynFileInfo{
						FromPath: deal.CarPath,
						ToPath:   define.PathIpfsDataCar,
						Ip:       global.LocalIP,
						Port:     define.OpPort,
						OpId:     deal.CarOpId,
						FileInfo: fileAr,
					}
					res, err := global.OpToGatewayClient.SysFilePoint(ctx, &in)
					if err != nil {
						global.ZC_LOG.Error("SysFilePoint：" + err.Error())
						return
					}

					var stat = &pb.EditStatus{
						Id:     deal.Id,
						Actor:  task.Actor,
						Status: int32(define.QueueSectorStatusCreate),
					}
					if res.Code != 200 {
						log.Println(res.Msg)
						stat.Status = int32(define.QueueSectorStatusFitFail)
						if _, err = global.OpToGatewayClient.EditQueueDetailStatus(ctx, stat); err != nil {
							global.ZC_LOG.Error("EditQueueDetailStatus：" + err.Error())
						}
						return
					}
					if atomic.LoadInt32(&ImportCarCount) >= task.ConcurrentImport {
						return
					}
					atomic.AddInt32(&ImportCarCount, 1)
					// Order import
					impcmd := fmt.Sprintf(`%s/boostd import-data %s %s`, define.PathIpfsProgram, deal.DealUuid, filepath.Join(define.PathIpfsDataCar, deal.PieceCid+".car"))
					log.Println("Ready to import:", impcmd, ImportCarCount)
					go func() {
						out, err := exec.Command("bash", "-c", impcmd).CombinedOutput()
						if err != nil {
							global.ZC_LOG.Error("import：" + string(out) + "," + err.Error())
							stat.Status = define.QueueSectorStatusFail
						}
						if strings.Contains(strings.ToLower(string(out)), "error") {
							global.ZC_LOG.Error("import：" + string(out))
							stat.Status = define.QueueSectorStatusFail
						}
						if _, err = global.OpToGatewayClient.EditQueueDetailStatus(context.Background(), stat); err != nil {
							global.ZC_LOG.Error("EditQueueDetailStatus：" + err.Error())
						}
						atomic.AddInt32(&ImportCarCount, -1)
					}()

				}(v)
			}
			dealWait.Wait()
		}
		break
	}
	return nil
}
