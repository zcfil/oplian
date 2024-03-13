package system

import (
	"context"
	"fmt"
	"log"
	"oplian/define"
	"oplian/global"
	model "oplian/model/lotus"
	"oplian/model/system/request"
	"oplian/model/system/response"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

var MonitorCenterServerApp = new(MonitorCenterServer)

type MonitorCenterServer struct {
}

// GetBusinessMonitor
// @function: GetBusinessMonitor
// @description: Service monitoring
// @param: info request.MonitorCenterReq
// @return: list []response.BusinessListReq, err error
func (m *MonitorCenterServer) GetBusinessMonitor(info request.MonitorCenterReq) (list []response.BusinessListReq, err error) {
	res := make([]response.BusinessListReq, 5)
	sectorStatus := make([]string, 0)
	//State of 001 / ap, 002 / p1 / p2, 003, 004 / c1, 005 / c2
	sectorStatus = append(sectorStatus, "001", "002", "003", "004", "005")
	if info.GateWayId != "" {
		gclient := global.GateWayClinets.GetGateWayClinet(info.GateWayId)
		if gclient == nil {
			return
		}
		var schedList *pb.SchedDiagRequestInfo

		miner, err := m.GetMinerManage(info.MinerId)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
		} else {
			//miner队列中
			schedList, err = gclient.SealingSchedDiag(context.Background(), &pb.FilParam{Token: miner.Token, Ip: miner.Ip})
			if err != nil {
				global.ZC_LOG.Error(err.Error())
			}
		}

		avgTime := make([]response.TimeTotal, 5)
		for i := 0; i < 5; i++ {
			//ws
			if i == define.WSMonitor.Int() {
				avgTime[i], err = m.SectorStatusAvgTime(info.MinerId, define.SealMonitorType(i).String())
				if err != nil {
					global.ZC_LOG.Error(err.Error())
				}
				avgTime[i].AvgSecond = int64(time.Second.Seconds()) * 75 * 60
				continue
			}
			//c2
			if i == define.C2Monitor.Int() {
				avgTime[i], err = m.Commit2AvgTime(info.MinerId)
				if err != nil {
					global.ZC_LOG.Error(err.Error())
				}
				continue
			}
			//ap,p1,p2
			avgTime[i], err = m.SectorStatusAvgTime(info.MinerId, define.SealMonitorType(i).String())
			if err != nil {
				global.ZC_LOG.Error(err.Error())
			}
		}
		var wait sync.WaitGroup
		wait.Add(len(res))
		for i, v := range sectorStatus {
			go func(index int, val string) {
				defer wait.Done()

				var runList *pb.TaskInfoList
				pendTotal := 0

				if index == define.WSMonitor.Int() {
					//ws不是从worker上执行任务，直接从数据库取
					sectors, err1 := m.SectorRunList(info.MinerId, define.WSMonitor.String())
					tasks := make([]*pb.TaskInfo, len(sectors))
					log.Println("secotrs:", len(sectors), err1)
					for j, sect := range sectors {
						timeLength := time.Since(sect.CreatedAt)
						progress, _ := strconv.Atoi(fmt.Sprintf("%.0f", timeLength.Seconds()*100/(75*time.Minute.Seconds())))
						tasks[j] = &pb.TaskInfo{
							MinerId:    sect.Actor,
							SectorId:   sect.SectorId,
							Progress:   int32(progress),
							TimeLength: utils.WholeSecond(timeLength),
							Ip:         miner.Ip,
						}
					}
					runList = &pb.TaskInfoList{Tasks: tasks}
				} else {

					runList, err = gclient.GetRunningList(context.Background(), &pb.OpTask{TType: define.SealMonitorType(index).String()})
					if err != nil {
						return
					}
					if schedList != nil {
						for _, v := range schedList.Requests {
							switch v.TaskType {
							case define.TTAddPiece.String():
								if index == define.ApMonitor.Int() {
									pendTotal++
								}
							case define.TTPreCommit1.String():
								if index == define.P1Monitor.Int() {
									pendTotal++
								}
							case define.TTPreCommit2.String():
								if index == define.P2Monitor.Int() {
									pendTotal++
								}
							case define.TTCommit1.String():
								if index == define.WSMonitor.Int() {
									pendTotal++
								}
							case define.TTCommit2.String():
								if index == define.C2Monitor.Int() {
									pendTotal++
								}
							}
						}
					}
				}

				opRes := make([]response.OpListRes, 0)
				//var timeCount time.Duration
				timeOut := 0
				for _, t := range runList.GetTasks() {

					if utils.FMinerID(t.MinerId) != info.MinerId {
						continue
					}
					op := response.OpListRes{
						MinerId:    info.MinerId,
						SectorID:   t.SectorId,
						Ip:         t.Ip,
						Progress:   int(t.Progress),
						TimeLength: t.TimeLength,
					}
					td, _ := time.ParseDuration(t.TimeLength)

					opRes = append(opRes, op)

					if td > define.SealMonitorType(index).TimeOut(t.SectorSize) {
						timeOut++
					}
				}
				aTime, _ := time.ParseDuration(fmt.Sprintf("%ds", avgTime[index].AvgSecond))
				res[index] = response.BusinessListReq{
					SectorStatus:   val,
					ProcessTotal:   len(opRes),
					PendTotal:      pendTotal,
					TimeLength:     aTime.String(),
					TimeOutTotal:   timeOut,
					OpListRes:      opRes,
					TotalCompleted: avgTime[index].Total,
				}
			}(i, v)
		}
		wait.Wait()
	}

	return res, nil
}

// SectorRunList Gets a list of sectors being run
func (m *MonitorCenterServer) SectorRunList(minerId, status string) (sectors []model.LotusSectorLog, err error) {
	table := fmt.Sprintf("%s_%s", model.LotusSectorLog{}.TableName(), minerId)
	sql := `SELECT * FROM ` + table + ` WHERE sector_status = ? AND IF(sector_status = 'WaitSeed',DATE_ADD(created_at,interval 75 MINUTE) >= NOW(),finish_at IS NULL)`
	log.Println(sql)
	return sectors, global.ZC_DB.Model(model.LotusSectorInfo{}).Raw(sql, status).Scan(&sectors).Error
}

// GetMinerManage Node scheduling miner
func (m *MonitorCenterServer) GetMinerManage(minerId string) (miner model.LotusMinerInfo, err error) {
	return miner, global.ZC_DB.Model(model.LotusMinerInfo{}).Where("actor = ? AND is_manage = 1", minerId).First(&miner).Error
}

// SectorStatusAvgTime Average time of sector phase
func (m *MonitorCenterServer) SectorStatusAvgTime(minerId string, status string) (res response.TimeTotal, err error) {
	table := fmt.Sprintf("%s_%s", model.LotusSectorLog{}.TableName(), minerId)

	sql := `SELECT IFNULL(ROUND(AVG(TIMESTAMPDIFF(second,created_at,finish_at)),0),0)avg_second,count(1)total FROM ` + table + ` 
			WHERE (finish_at IS NOT NULL or (sector_status = 'WaitSeed' AND DATE_ADD(created_at,interval 75 MINUTE) < NOW())) AND error_msg = '' AND sector_status = ?`
	if err = global.ZC_DB.Raw(sql, status).Scan(&res).Error; err != nil {
		return res, err
	}

	//avgTime, _ := time.ParseDuration(fmt.Sprintf("%ds", res.AvgTime))
	return res, nil
}

// Commit2AvgTime C2 average time
func (m *MonitorCenterServer) Commit2AvgTime(minerId string) (res response.TimeTotal, err error) {
	minerId = strings.Replace(minerId, "f0", "t0", 1)
	table := fmt.Sprintf("%s_%s", model.LotusWorkerRelations{}.TableName(), time.Now().Format("200601"))
	sql := `SELECT ROUND(AVG(TIMESTAMPDIFF(second,begin_time,end_time)),0)avg_second,count(1)total FROM ` + table + ` WHERE end_time IS NOT NULL AND miner = ?`
	if err = global.ZC_DB.Raw(sql, minerId).Scan(&res).Error; err != nil || res.Total == 0 {

		table = fmt.Sprintf("%s_%s", model.LotusWorkerRelations{}.TableName(), time.Now().AddDate(0, -1, 0).Format("200601"))
		if err = global.ZC_DB.Raw(sql, minerId).Scan(&res).Error; err != nil {
			return res, err
		}
	}
	//avgTime, _ := time.ParseDuration(fmt.Sprintf("%ds", res.AvgTime))
	return res, nil
}
