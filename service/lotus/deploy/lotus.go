package deploy

import (
	"context"
	"fmt"
	"log"
	"oplian/global"
	model "oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"time"
)

type DeployService struct{}

func (deploy *DeployService) AddLotus(lotus *model.LotusInfo) error {
	if lotus.ID == 0 {
		return global.ZC_DB.Save(lotus).Error
	}
	return global.ZC_DB.Updates(lotus).Error
}

func (deploy *DeployService) GetLotusList(info request.LotusInfoPage) (lotusList []*response.LotusInfo, total int64, err error) {
	sqlparam := " where 1=1 "
	sql := `SELECT l.id,l.op_id,r.room_id,l.gate_id,r.room_name,host_name,device_sn,l.ip,l.port,l.actor,mc.miner_count,l.start_at,l.finish_at,snapshot_at,l.sync_status,l.deploy_status,l.token
			FROM lotus_info l
			LEFT JOIN (
				SELECT lotus_id,count(1) miner_count FROM lotus_miner_info GROUP BY lotus_id
			)mc ON l.id = mc.lotus_id
			LEFT JOIN sys_host_records r ON l.op_id = r.uuid`
	if info.Keyword != "" {
		sqlparam += ` and (l.ip like '%` + info.Keyword + `%' or r.host_name like '%` + info.Keyword + `%' or r.room_name like '%` + info.Keyword + `%')`
	}
	if info.SyncStatus != 0 {
		sqlparam += ` and l.sync_status = '` + strconv.Itoa(info.SyncStatus) + `'`
	}
	if info.DeployStatus != 0 {
		sqlparam += ` and l.deploy_status = '` + strconv.Itoa(info.DeployStatus) + `'`
	}
	if info.GateId != "" {
		sqlparam += ` and l.gate_id = '` + info.GateId + `'`
	}

	sql += sqlparam
	sql += utils.LimitAndOrder("host_name", "desc", info.Page, info.PageSize)
	err = global.ZC_DB.Raw(sql).Scan(&lotusList).Error
	if err != nil {
		return nil, 0, err
	}

	sqlTotal := `SELECT count(1) FROM lotus_info l LEFT JOIN sys_host_records r ON l.op_id = r.uuid `
	err = global.ZC_DB.Model(&model.LotusInfo{}).Raw(sqlTotal + sqlparam).Count(&total).Error

	for _, v := range lotusList {
		gclient := global.GateWayClinets.GetGateWayClinet(v.GateId)
		if gclient == nil {
			log.Println("not exist！" + v.GateId)
			continue
		}
		v.Online = true
		walletList, err := gclient.GetWalletList(context.Background(), &pb.RequestOp{OpId: v.OpId, GateId: v.GateId, Ip: v.Ip, Token: v.Token})
		if err != nil {
			v.WalletCount = 0
			log.Println("Fetch wallet error！" + v.GateId + err.Error())
			continue
		}
		log.Println(v.Ip, "Number of wallets：", len(walletList.Wallets))
		v.WalletCount = len(walletList.Wallets)
	}

	return lotusList, total, err
}

func (deploy *DeployService) GetRoomAllLotus(param request.LotusInfoPage) ([]*response.RelationLotusInfo, error) {
	var lotusList []*response.RelationLotusInfo
	sqlparam := `WHERE l.gate_id = '` + param.GateId + `' `
	sql := `SELECT l.id,l.op_id,l.gate_id,r.room_id,r.room_name,host_name,l.ip,l.actor,l.token
			FROM lotus_info l
			LEFT JOIN sys_host_records r ON l.op_id = r.uuid
			 `
	if param.SyncStatus != 0 {
		sqlparam += " and l.sync_status = " + strconv.Itoa(param.SyncStatus)
	}
	if param.DeployStatus != 0 {
		sqlparam += " and l.deploy_status = " + strconv.Itoa(param.DeployStatus)
	}
	if param.Keyword != "" {
		sqlparam += ` and l.ip like '%` + param.Keyword + `%' or host_name like '%` + param.Keyword + `%'`
	}
	if err := global.ZC_DB.Raw(sql + sqlparam).Scan(&lotusList).Error; err != nil {
		return nil, err
	}

	for _, v := range lotusList {
		gclient := global.GateWayClinets.GetGateWayClinet(v.GateId)
		if gclient == nil {
			global.ZC_LOG.Warn(v.GateId + "not exist！")
			continue
		}
		ws, err := gclient.GetWalletList(context.Background(), &pb.RequestOp{GateId: v.GateId, OpId: v.OpId, Ip: v.Ip, Token: v.Token})
		if err != nil {
			global.ZC_LOG.Warn(v.Ip + "Fetch wallet error：" + err.Error())
			continue
		}
		for _, w := range ws.Wallets {
			v.Wallets = append(v.Wallets, response.Wallet{Address: w.Address, Balance: w.Balance})
		}
	}

	return lotusList, nil
}

func (deploy *DeployService) GetLotus(id uint64) (model.LotusInfo, error) {
	var lotus model.LotusInfo
	return lotus, global.ZC_DB.Model(model.LotusInfo{}).Where("id = ?", id).First(&lotus).Error
}

func (deploy *DeployService) UpdateLotus(lotus model.LotusInfo) error {
	return global.ZC_DB.Save(&lotus).Error
}

func (deploy *DeployService) UpdateLotusStatus(ID uint, status int) error {
	var lotus model.LotusInfo
	db := global.ZC_DB.Model(model.LotusInfo{})
	if err := db.Where("id = ?", ID).First(&lotus).Error; err != nil {
		return err
	}
	lotus.RunStatus = status
	return global.ZC_DB.Save(&lotus).Error
}

func (deploy *DeployService) GetLotusByOpID(opId string) (model.LotusInfo, error) {
	var lotus model.LotusInfo
	return lotus, global.ZC_DB.Model(model.LotusInfo{}).Where("op_id = ?", opId).First(&lotus).Error
}

func (deploy *DeployService) ModifyLotusStatus(lotus map[string]int) error {

	ipStr := ""
	for k, _ := range lotus {
		if utils.IsNull(ipStr) {
			ipStr = fmt.Sprintf("'%s'", k)
		} else {
			ipStr += fmt.Sprintf(",'%s'", k)
		}
	}

	if ipStr != "" {

		modifyMap := make(map[uint]int)
		var minerList []model.LotusInfo
		err := global.ZC_DB.Model(&model.LotusInfo{}).Where("ip in(" + ipStr + ")").Find(&minerList).Error
		if err != nil {
			return err
		}

		if len(minerList) > 0 {

			for _, v := range minerList {
				if val, ok := lotus[v.Ip]; ok {
					if val != v.RunStatus {
						modifyMap[v.ID] = val
					}
				}
			}

			if len(modifyMap) > 0 {
				for k, v := range modifyMap {

					err = global.ZC_DB.Model(&model.LotusInfo{}).Where("id", k).Update("run_status", v).Error
					if err != nil {
						return err
					}
					time.Sleep(time.Second)
				}
			}
		}
	}

	return nil
}
