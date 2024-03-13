package deploy

import (
	"context"
	"fmt"
	"log"
	"oplian/define"
	"oplian/global"
	model "oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/lotus/response"
	"oplian/model/system"
	sysrequest "oplian/model/system/request"
	"oplian/service/pb"
	"oplian/utils"
	"time"
)

func (deploy *DeployService) AddMiner(miner *model.LotusMinerInfo) error {
	if miner.ID == 0 {
		return global.ZC_DB.Save(miner).Error
	}
	return global.ZC_DB.Updates(miner).Error
}

func (deploy *DeployService) UpdateMinerStatusAndLink(ID uint, status int, linkId uint64) error {
	var miner model.LotusMinerInfo
	db := global.ZC_DB.Model(model.LotusMinerInfo{})
	if err := db.Where("id = ?", ID).First(&miner).Error; err != nil {
		return err
	}
	miner.RunStatus = status
	if linkId != 0 {
		miner.LotusId = linkId
	}
	return db.Save(&miner).Error
}

func (deploy *DeployService) GetMiner(id uint64) (model.LotusMinerInfo, error) {
	var miner model.LotusMinerInfo
	return miner, global.ZC_DB.Model(model.LotusMinerInfo{}).Where("id = ?", id).First(&miner).Error
}

func (deploy *DeployService) GetMinerByOpId(opId string) (model.LotusMinerInfo, error) {
	var miner model.LotusMinerInfo
	return miner, global.ZC_DB.Model(model.LotusMinerInfo{}).Where("op_id = ?", opId).First(&miner).Error
}

func (deploy *DeployService) GetMinerToken(opId string) (string, error) {

	var miner model.LotusMinerInfo
	sql := `SELECT m.token,m.ip FROM lotus_worker_info w
				INNER JOIN lotus_miner_info m ON w.miner_id=m.id WHERE w.op_id=?`

	if err := global.ZC_DB.Raw(sql, opId).Scan(&miner).Error; err != nil {
		return "", err
	}

	token := fmt.Sprintf("%s:/ip4/%s/tcp/%s/http", miner.Token, miner.Ip, define.MinerPort)
	return token, nil
}

// @author: lex
// @function: GetMinerByActor
// @description: 获取miner
// @param: actor string
// @return: model.LotusMinerInfo, error

func (deploy *DeployService) GetMinerByActor(actor string) (model.LotusMinerInfo, error) {
	var miner model.LotusMinerInfo
	return miner, global.ZC_DB.Model(model.LotusMinerInfo{}).Where("actor = ?", actor).First(&miner).Error
}

func (deploy *DeployService) GetMinerRun(id uint64) (miner response.RunMiner, err error) {
	sql := `SELECT lm.ip,lm.port,lm.deploy_status,lm.run_status,lm.is_manage,lm.is_wdpost,lm.is_wnpost,lm.partitions,l.token lotus_token,l.ip lotus_ip,lm.gate_id,lm.op_id
			FROM lotus_miner_info lm
			LEFT JOIN lotus_info l on lm.lotus_id = l.id
			where lm.id = ?`
	if err = global.ZC_DB.Raw(sql, id).Scan(&miner).Error; err != nil {
		return response.RunMiner{}, err
	}
	return miner, nil
}

func (deploy *DeployService) GetMinerList(info request.MinerTypePage) (list interface{}, total int64, err error) {
	var minerList []*response.MinerInfo
	sqlparam := "WHERE 1=1 "
	if info.Keyword != "" {
		sqlparam += fmt.Sprintf(" AND (lm.actor like '%%%s%%' OR lm.ip like '%%%s%%') ", info.Keyword, info.Keyword)
	}
	if info.ColonyType != 0 {
		sqlparam += fmt.Sprintf(" AND c.colony_type = %d ", info.ColonyType)
	}
	if info.GateId != "" {
		sqlparam += fmt.Sprintf(" AND lm.gate_id = '%s' ", info.GateId)
	}
	if info.IsManage {
		sqlparam += fmt.Sprintf(" AND lm.is_manage = 1 ")
	}

	sql := `SELECT b.id bid,lm.id,lm.op_id,lm.gate_id,r.room_id,r.room_name,r.host_name,r.internet_ip,device_sn,lm.ip,lm.port,lm.actor,lm.start_at,lm.finish_at,lm.deploy_status,lm.sector_size,lm.add_type,
				lm.partitions,lm.is_manage,lm.is_wdpost,lm.is_wnpost,worker_count,storage_count,l.id lotus_id,l.ip lotus_ip,l.host_name lotus_host_name,l.token lotus_token,lm.err_msg,c.colony_type
			FROM lotus_miner_info lm
			LEFT JOIN (
				SELECT count(1) worker_count,miner_id
						FROM lotus_worker_info
						GROUP BY miner_id
			)w on lm.id = w.miner_id 
			LEFT JOIN (
					select count(1) storage_count,colony_name 
						FROM lotus_storage_info
						GROUP BY colony_name
				)a ON lm.actor = a.colony_name
			LEFT JOIN sys_host_records r ON lm.op_id = r.uuid
			LEFT JOIN (
				SELECT op_id,ip,host_name,token,l.id
				FROM lotus_info l
				LEFT JOIN sys_host_records h ON l.op_id = h.uuid
			)l ON lm.lotus_id = l.id
			LEFT JOIN sys_colony c ON lm.actor = c.colony_name
			LEFT JOIN lotus_boost_info b ON lm.id = b.miner_id AND b.deploy_status = 2
			` + sqlparam

	sqlTotal := `SELECT COUNT(1) FROM lotus_miner_info lm LEFT JOIN sys_colony c ON lm.actor = c.colony_name `
	err = global.ZC_DB.Model(&model.LotusMinerInfo{}).Raw(sqlTotal + sqlparam).Count(&total).Error
	sql += utils.LimitAndOrder("lm.actor", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&minerList).Error; err != nil {
		return nil, 0, err
	}

	for _, v := range minerList {
		httpToken := v.LotusToken
		v.LotusToken = ""
		gcli := global.GateWayClinets.GetGateWayClinet(v.GateId)
		if gcli == nil {
			global.ZC_LOG.Error(v.GateId + "not exist！")
			continue
		}

		if online, err := gcli.OpOnline(context.Background(), &pb.String{Value: v.OpId}); err == nil {
			v.Online = online.Value
		} else {
			global.ZC_LOG.Error(v.GateId + ",OpOnline:" + err.Error())
		}
		if v.Actor != "" && httpToken != "" && v.Online {
			param := &pb.FilParam{Token: httpToken, Ip: v.LotusIp, Param: v.Actor}
			wallets, err := deploy.MinerWalletInfo(gcli, param)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				continue
			}
			v.Wallets = wallets
			v.WalletCount = len(wallets)
		}
	}

	return minerList, total, err
}

func (deploy *DeployService) GetMinerMonitorList(info sysrequest.HostMonitorReq) (list []response.MinerMonitorInfo, total int64, err error) {
	sqlparam := " WHERE 1=1 "
	if info.Keyword != "" {
		sqlparam += fmt.Sprintf(" AND (lm.ip like '%%%s%%') ", info.Keyword)
	}
	if info.GateId != "" {
		sqlparam += fmt.Sprintf(" AND lm.gate_id = '%s' ", info.GateId)
	}
	sql := `SELECT * FROM lotus_miner_info lm `
	//求总数
	sqlTotal := `SELECT COUNT(1) FROM lotus_miner_info lm `
	err = global.ZC_DB.Model(&model.LotusMinerInfo{}).Raw(sqlTotal + sqlparam).Count(&total).Error
	sql += sqlparam + utils.LimitAndOrder("actor", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return
}

func (deploy *DeployService) MinerWalletInfo(gcli pb.GateServiceClient, param *pb.FilParam) ([]response.Wallet, error) {
	actor, err := gcli.StateMinerInfo(context.Background(), param)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return nil, err
	}

	var wallets []response.Wallet

	if actor.Owner != "" {
		param.Param = actor.Owner
		addr, err := gcli.StateAccountKey(context.Background(), param)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}
		param.Param = addr.Value
		balance, err := gcli.WalletBalance(context.Background(), param)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
		}
		wallets = append(wallets, response.Wallet{Address: addr.Value, Balance: balance.Balance, Attribute: "owner钱包"})
	}

	if actor.Worker != "" {
		param.Param = actor.Worker
		addr, err := gcli.StateAccountKey(context.Background(), param)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return nil, err
		}
		param.Param = addr.Value
		balance, err := gcli.WalletBalance(context.Background(), param)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
		}
		wallets = append(wallets, response.Wallet{Address: addr.Value, Balance: balance.Balance, Attribute: "worker钱包"})
	}

	if len(actor.Control) > 0 {
		for _, control := range actor.Control {
			param.Param = control
			addr, err := gcli.StateAccountKey(context.Background(), param)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				continue
			}
			param.Param = addr.Value
			balance, err := gcli.WalletBalance(context.Background(), param)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
			}
			wallets = append(wallets, response.Wallet{Address: addr.Value, Balance: balance.Balance, Attribute: "post钱包"})
		}
	}
	//log.Println("wallets:", wallets)
	return wallets, err
}

func (deploy *DeployService) GetNodeList(info sysrequest.ColonyPageInfo) (list interface{}, total int64, err error) {
	var nodeList []*response.NodeInfo
	sqlparam := "WHERE 1=1"
	if info.ColonyType != 0 {
		sqlparam = fmt.Sprintf(" AND colony_type = %d", info.ColonyType)
	}
	if info.Keyword != "" {
		sqlparam += fmt.Sprintf(" AND (c.colony_name like '%%%s%%') ", info.Keyword)
	}
	//求总数
	sqlTotal := `select count(1) from sys_colony ` + sqlparam
	err = global.ZC_DB.Model(&system.SysColony{}).Raw(sqlTotal).Count(&total).Error

	sql := `select c.*,m.actor,IFNULL(w.wcount,0) wcount,IFNULL(m.mcount,0)mcount,IFNULL(s.scount,0)scount,m.sector_size from sys_colony c
				LEFT JOIN (SELECT count(1) mcount,actor,sector_size FROM lotus_miner_info WHERE deploy_status = 2 AND actor <> '' GROUP BY actor ) m on c.colony_name = m.actor
				LEFT JOIN (SELECT count(1) wcount,actor FROM lotus_worker_info w 
									LEFT JOIN lotus_miner_info m ON w.miner_id = m.id  
									WHERE worker_type = 0 AND m.deploy_status = 2 AND IFNULL(actor,'') <> '' GROUP BY actor 
							) w on c.colony_name = w.actor
				LEFT JOIN (SELECT count(1) scount,colony_name FROM lotus_storage_info s 
									WHERE s.deploy_status = 2 AND IFNULL(colony_name,'') <> '' GROUP BY colony_name 
							) s on c.colony_name = s.colony_name ` + sqlparam

	sql += utils.LimitAndOrder("w.wcount", "desc", info.Page, info.PageSize)
	if err = global.ZC_DB.Raw(sql).Scan(&nodeList).Error; err != nil {
		return nil, 0, err
	}
	//查找可用lotus
	lotusList, _, err := deploy.GetLotusList(request.LotusInfoPage{SyncStatus: define.SyncFinish.Int(), DeployStatus: define.DeployFinish.Int()})
	if err != nil || len(lotusList) == 0 {

		return nodeList, total, nil
	}

	for _, v := range nodeList {
		for _, lotusInfo := range lotusList {
			if lotusInfo.Online {
				gclient := global.GateWayClinets.GetGateWayClinet(lotusInfo.GateId)
				if gclient == nil {
					global.ZC_LOG.Error(lotusInfo.GateId + "不存在！")
					continue
				}
				param := &pb.FilParam{Token: lotusInfo.Token, Ip: lotusInfo.Ip, Param: v.ColonyName}
				wallets, err := deploy.MinerWalletInfo(gclient, param)
				if err != nil {
					global.ZC_LOG.Error(err.Error())
					continue
				}
				v.Wallets = wallets
				v.WalletCount = len(v.Wallets)
				param.Param = v.ColonyName
				sectors, err := gclient.StateMinerSectorCount(context.Background(), param)
				if err != nil {
					global.ZC_LOG.Error(err.Error())
					continue
				}
				power, err := gclient.StateMinerPower(context.Background(), param)
				if err != nil {
					global.ZC_LOG.Error(err.Error())
					continue
				}
				v.Power = power.MinerPower
				v.TotalPower = power.TotalPower
				v.Live = sectors.Live
				v.Faulty = sectors.Faulty
				v.Active = sectors.Active
				break
			}
		}

	}
	return nodeList, total, nil
}

func (deploy *DeployService) GetRelationMinerList(id request.IDActor) (list interface{}, err error) {
	sqlparam := fmt.Sprintf(" WHERE actor ='%s'", id.Actor)
	if id.Actor == "" {
		sqlparam = fmt.Sprintf(" WHERE lotus_id = %d", id.ID)
	}
	var minerList []response.RelationMiner
	sql := `SELECT lm.id,lm.op_id,lm.gate_id,r.room_id,r.room_name,host_name,device_sn,lm.ip,lm.actor,lm.sector_size
			FROM lotus_miner_info lm 
			LEFT JOIN sys_host_records r ON lm.op_id = r.uuid
			` + sqlparam
	return minerList, global.ZC_DB.Raw(sql).Scan(&minerList).Error
}

func (deploy *DeployService) GetPledge(opId string) (model.LotusPledgeConfig, error) {
	var pledge model.LotusPledgeConfig
	return pledge, global.ZC_DB.Model(model.LotusPledgeConfig{}).Where("op_id = ?", opId).First(&pledge).Error
}

func (deploy *DeployService) AddUndoneP1(OpId string, num int) error {
	var pledge model.LotusPledgeConfig
	if err := global.ZC_DB.Model(model.LotusPledgeConfig{}).Where("op_id = ?", OpId).First(&pledge).Error; err != nil {
		return err
	}
	pledge.UndonePre += num
	return global.ZC_DB.Save(&pledge).Error
}

func (deploy *DeployService) GetNodesNum(gateId string) (list []response.ActorInfo, err error) {
	var actor []response.ActorInfo
	sqlparam := ""
	if gateId != "" {
		sqlparam += fmt.Sprintf(" WHERE gate_id = '%s'", gateId)
	}
	sql := `SELECT actor,gate_id,sector_size,c.colony_type FROM lotus_miner_info m
					LEFT JOIN sys_colony c ON m.actor = c.colony_name
					` + sqlparam + `
					GROUP BY actor`
	if err = global.ZC_DB.Raw(sql).Scan(&actor).Error; err != nil {
		return nil, err
	}
	return actor, err
}

func (deploy *DeployService) GetManageMiners(keyword string) (list []*response.ManageInfo, err error) {
	where := "is_manage = ?"
	if keyword != "" {
		where += ` and (actor LIKE '%` + keyword + `%' or ip LIKE '%` + keyword + `%')`
	}
	if err = global.ZC_DB.Model(model.LotusMinerInfo{}).Where(where, true).Find(&list).Error; err != nil {
		return nil, err
	}

	for _, v := range list {
		info, err := deploy.GetLotus(v.LotusId)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			continue
		}
		gcli := global.GateWayClinets.GetGateWayClinet(v.GateId)
		if gcli == nil {
			global.ZC_LOG.Error(v.GateId + "not exist！")
			continue
		}
		param := &pb.FilParam{Token: info.Token, Ip: v.Ip, Param: v.Actor}
		actor, err := gcli.StateMinerInfo(context.Background(), param)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			continue
		}

		if actor.Worker != "" {
			param.Param = actor.Worker
			addr, err := gcli.StateAccountKey(context.Background(), param)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				continue
			}
			v.WorkerWallet = addr.Value
			param.Param = addr.Value
			balance, err := gcli.WalletBalance(context.Background(), param)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
			}
			v.WorkerBalance = balance.Balance
		}
	}
	return
}

func (deploy *DeployService) GateMinerList(gateId string) (list []model.LotusMinerInfo, err error) {
	sql := `SELECT m.* FROM lotus_miner_info m
					WHERE m.deploy_status = ? AND m.deleted_at IS NULL AND m.gate_id = ?`
	if err = global.ZC_DB.Raw(sql, define.DeployFinish, gateId).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (deploy *DeployService) GateMinerSelectList(gateId string) (list []response.MinerSelectInfoResp, err error) {
	sql := `SELECT m.* FROM lotus_miner_info m
					WHERE m.deploy_status = ? AND m.deleted_at IS NULL AND m.gate_id = ?`
	if err = global.ZC_DB.Raw(sql, define.DeployFinish, gateId).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}


func (deploy *DeployService) GetMinerListByActor(actor string) (list []response.MinerInfo, err error) {
	sql := `SELECT m.*,l.token lotus_token,l.ip lotus_ip FROM lotus_miner_info m
					LEFT JOIN lotus_info l ON m.lotus_id = l.id
					WHERE m.actor = ? AND m.deploy_status = ? AND m.deleted_at IS NULL`
	if err = global.ZC_DB.Raw(sql, actor, define.DeployFinish).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (deploy *DeployService) CheckMinerRole(param request.MinerParam) (msg string, err error) {
	var partitions, wnpostIp, manageIp string

	actors, err := deploy.GetMinerListByActor(param.Actor)
	if err != nil {
		return "", fmt.Errorf("failed to get node information：%v", err)
	}
	for _, miner := range actors {
		if miner.Id == param.Id {
			continue
		}
		if miner.IsWnpost {
			wnpostIp = miner.Ip
			if param.IsWnpost {
				return "", fmt.Errorf("existing deployed successful winning post machine：%s", miner.Ip)
			}
		}
		if miner.IsManage {
			manageIp = miner.Ip
			if param.IsManage {
				return "", fmt.Errorf("existing scheduling machine deployment is successful：%s", miner.Ip)
			}
		}
		if miner.IsWdpost {
			partitions += miner.Partitions
		}
	}
	msg = fmt.Sprintf(`Has set the partitions to：%s
The wnpost server has been configured：%s
The manage server has been configured：%s
			`, partitions, wnpostIp, manageIp)
	return msg, nil
}

func (deploy *DeployService) CheckMinerRolezcjs(param request.MinerParam) (msg string, err error) {
	var partitions, wnpostIp, manageIp string
	if param.LotusId != 0 {
		lotusInfo, err := deploy.GetLotus(param.LotusId)
		if err != nil {
			return "", fmt.Errorf("lotus information failure：%v", err)
		}
		gclient := global.GateWayClinets.GetGateWayClinet(lotusInfo.GateId)
		if gclient == nil {
			return "", fmt.Errorf("%s not exist！", lotusInfo.GateId)
		}
		actorWallet, err := gclient.StateMinerInfo(context.Background(), &pb.FilParam{Ip: lotusInfo.Ip, Token: lotusInfo.Token, Param: param.Actor})
		if err != nil {
			return "", fmt.Errorf("miner information failure：%v", err)
		}
		wallets, err := gclient.GetWalletList(context.Background(), &pb.RequestOp{OpId: lotusInfo.OpId, Ip: lotusInfo.Ip, Token: lotusInfo.Token})
		if err != nil {
			return "", fmt.Errorf("failed to get purse lists：%v", err)
		}
		haveWdpost := !param.IsWdpost
		haveWorker := false
		var wdpostAddr, workerAddr string

		addr, err := gclient.StateAccountKey(context.Background(), &pb.FilParam{Ip: lotusInfo.Ip, Token: lotusInfo.Token, Param: actorWallet.Worker})
		workerAddr = addr.Value
		if err != nil {
			return "", fmt.Errorf("failed to obtain the worker wallet：%v", err)
		}
		//未设置wdpost钱包默认是worker钱包
		if len(actorWallet.Control) == 0 {
			actorWallet.Control = append(actorWallet.Control, actorWallet.Worker)
		} else {
			addr, err = gclient.StateAccountKey(context.Background(), &pb.FilParam{Ip: lotusInfo.Ip, Token: lotusInfo.Token, Param: actorWallet.Control[0]})
			if err != nil {
				return "", fmt.Errorf("failed to obtain the wdpost walle：%v", err)
			}
		}
		wdpostAddr = addr.Value

		for _, wallet := range wallets.Wallets {
			if !haveWorker {

				if workerAddr == wallet.Address {
					haveWorker = true

					if len(actorWallet.Control) == 0 {
						haveWdpost = true
					}
				}
			}
			if !haveWdpost {
				if addr.Value == wallet.Address {
					haveWdpost = true
				}
			}

			if haveWdpost && haveWorker {
				break
			}
		}

		if !haveWorker {
			log.Println("The worker wallet does not exist on this lotus node：", workerAddr)
		}

		if !haveWdpost {
			log.Println("The wdpost wallet does not exist on this lotus node：", wdpostAddr)
		}
	}

	actors, err := deploy.GetMinerListByActor(param.Actor)
	if err != nil {
		return "", fmt.Errorf("failed to get node information：%v", err)
	}
	for _, miner := range actors {
		if miner.Id == param.Id {
			continue
		}
		if miner.IsWnpost {
			wnpostIp = miner.Ip
			if param.IsWnpost {
				return "", fmt.Errorf("existing deployed successful winning post machine：%s", miner.Ip)
			}
		}
		if miner.IsManage {
			manageIp = miner.Ip
			if param.IsManage {
				return "", fmt.Errorf("existing scheduling machine deployment is successful：%s", miner.Ip)
			}
		}
		if miner.IsWdpost {
			partitions += miner.Partitions
		}
	}
	msg = fmt.Sprintf(`partitions have been set to：%s
The wnpost server has been configured：%s
The manage server has been configured：%s
			`, partitions, wnpostIp, manageIp)
	return msg, nil
}

func (deploy *DeployService) ModifyMinerRole(param request.MinerParam) (err error) {
	minerParam := map[string]interface{}{
		"is_wdpost":  param.IsWdpost,
		"is_wnpost":  param.IsWnpost,
		"is_manage":  param.IsManage,
		"partitions": param.Partitions,
	}
	return global.ZC_DB.Model(model.LotusMinerInfo{}).Where("id=?", param.Id).Updates(minerParam).Error
}

func (deploy *DeployService) ModifyMinerStatus(miner map[string]int) error {

	ipStr := ""
	for k, _ := range miner {
		if utils.IsNull(ipStr) {
			ipStr = fmt.Sprintf("'%s'", k)
		} else {
			ipStr += fmt.Sprintf(",'%s'", k)
		}
	}

	if ipStr != "" {

		modifyMap := make(map[uint]int)
		var minerList []model.LotusMinerInfo
		err := global.ZC_DB.Model(&model.LotusMinerInfo{}).Where("ip in(" + ipStr + ")").Find(&minerList).Error
		if err != nil {
			return err
		}

		if len(minerList) > 0 {

			for _, v := range minerList {
				if val, ok := miner[v.Ip]; ok {
					if val != v.RunStatus {
						modifyMap[v.ID] = val
					}
				}
			}

			if len(modifyMap) > 0 {
				for k, v := range modifyMap {

					err = global.ZC_DB.Model(&model.LotusMinerInfo{}).Where("id", k).Update("run_status", v).Error
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
