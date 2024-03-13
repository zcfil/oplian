package gateway

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"oplian/config"
	"oplian/define"
	"oplian/global"
	"oplian/lotusrpc"
	request1 "oplian/model/common/request"
	"oplian/model/lotus"
	model "oplian/model/lotus"
	"oplian/model/lotus/request"
	"oplian/model/system"
	"oplian/service/pb"
	"oplian/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ClearWorker Clear local data for workers
func (g *GateWayServiceImpl) ClearWorker(ctx context.Context, args *pb.RequestOp) (*pb.ResponseMsg, error) {
	opClient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return &pb.ResponseMsg{Code: 404, Msg: "Op not online"}, errors.New("Op：" + args.OpId + " not online！")
	}
	res, err := opClient.ClearWorker(ctx, args)
	if err != nil {
		return res, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, err
}

// SealingAbort Stop the task
func (g *GateWayServiceImpl) SealingAbort(ctx context.Context, args *pb.ResetWorker) (*emptypb.Empty, error) {
	server, err := lotusService.GetWorkerMiner(args.Id)
	if err != nil {
		return nil, fmt.Errorf("GetWorkerMiner：%s", err.Error())
	}
	opClient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return nil, errors.New("Op：" + args.OpId + " not online！")
	}
	if server.Id == args.LinkId {
		return nil, fmt.Errorf("%d Link ID %d has not changed：%d", args.Id, server.Id, args.LinkId)
	}
	if err = lotusService.UpdateWorkerStatusAndLink(uint(args.Id), -1, define.DeployReset.Int(), args.LinkId); err != nil {
		return nil, fmt.Errorf("UpdateWorkerStatusAndLink：%s", err.Error())
	}
	var fil = &pb.FilParam{Token: server.Token, Ip: server.Ip, Param: server.Actor}
	newMiner, err := lotusService.GetMiner(args.LinkId)
	if err != nil {
		return nil, fmt.Errorf("GetMiner：%s", err.Error())
	}
	newHost := &pb.FilParam{Token: newMiner.Token, Ip: newMiner.Ip, Param: newMiner.Actor}
	return opClient.SealingAbort(ctx, &pb.FilRestWorker{Worker: args, Host: fil, NewHost: newHost})
}

// AddLotus Run lotus
func (g *GateWayServiceImpl) AddLotus(ctx context.Context, args *pb.LotusInfo) (*emptypb.Empty, error) {
	opClient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return &emptypb.Empty{}, errors.New("Op：" + args.OpId + " not online！")
	}

	mlotus := &model.LotusInfo{
		GateId:     args.GateId,
		OpId:       args.OpId,
		Ip:         args.Ip,
		Port:       define.LotusPort,
		Token:      "",
		SyncStatus: int(define.SyncFinish),
		StartAt:    time.Now(),
	}
	// Modify synchronization status
	if args.ImportMode != int32(define.LotusInitRunModel) {
		mlotus.SyncStatus = int(define.Synchronizing)
		mlotus.DeployStatus = int(define.DeployRunning)
	}
	mlotus.ID = uint(args.LotusId)
	err := lotusService.DeployService.AddLotus(mlotus)
	if err != nil {
		err = fmt.Errorf("error adding lotus：%s", err.Error())
		global.ZC_LOG.Error(err.Error())
		if strings.Contains(err.Error(), "Duplicate") {
			err = fmt.Errorf("the server has repeatedly created lotus")
		}
		return &emptypb.Empty{}, err
	}
	//run
	sys, err := systemService.GetSysHostRecord(args.OpId)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	args.ReRaid = !sys.IsGroupArray

	err = systemService.UpdateSysHostRecordClassify(&system.SysHostRecord{UUID: args.OpId, HostClassify: config.HostLotusType, IsGroupArray: true})
	if err != nil {
		err = fmt.Errorf("error modifying server type：%s", err.Error())
		global.ZC_LOG.Error(err.Error())
		return &emptypb.Empty{}, err
	}

	args.LotusId = uint64(mlotus.ID)
	_, err = opClient.RunNewLotus(ctx, args)
	if err != nil {
		err = fmt.Errorf("starting lotus error：%s", err.Error())
		global.ZC_LOG.Error(err.Error())
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// GetWalletList Get wallet list
func (g *GateWayServiceImpl) GetWalletList(ctx context.Context, args *pb.RequestOp) (*pb.WalletList, error) {
	opClient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return nil, errors.New("Op：" + args.OpId + " not online！")
	}

	return opClient.GetWalletList(ctx, &pb.RequestConnect{Ip: args.Ip, Token: args.Token})
	//return &pb.WalletList{Wallets: []*pb.Wallet{{Address: "f1sdklhfdsklfhsofsdfsdhsdkl", Balance: 10}, {Address: "f3sdklhfdsdsfgsdfdsfdsnhkfjkyhsadfkasgfjkdsagfsadgjkfklfhsofsdfsdhsdkl", Balance: 0.456456441}}}, nil
}

// GetRoomWalletList Get wallet list
func (g *GateWayServiceImpl) GetRoomWalletList(ctx context.Context, args *pb.RequestOp) (*pb.WalletList, error) {
	var param request.LotusInfoPage
	param.GateId = args.GateId
	ls, err := lotusService.GetRoomAllLotus(param)
	if err != nil {
		return nil, err
	}
	var wallet pb.WalletList
	fmt.Println("Obtain the number of lotus：", len(ls))
	for _, v := range ls {
		_, dis := global.OpClinets.GetOpClient(v.OpId)
		if dis {
			continue
		}
		fmt.Println("Obtain the number of wallets：", v.Ip, len(v.Wallets))
		for _, w := range v.Wallets {
			wallet.Wallets = append(wallet.Wallets, &pb.Wallet{Address: w.Address, OpId: v.OpId, Balance: w.Balance})
		}
	}
	return &wallet, nil
	//return &pb.WalletList{Wallets: []*pb.Wallet{{Address: "f1sdklhfdsklfhsofsdfsdhsdkl", Balance: 10}, {Address: "f3sdklhfdsdsfgsdfdsfdsnhkfjkyhsadfkasgfjkdsagfsadgjkfklfhsofsdfsdhsdkl", Balance: 0.456456441}}}, nil
}

// CreateSectorTable Create sector table
func (g *GateWayServiceImpl) CreateSectorTable(ctx context.Context, args *pb.Actor) (*emptypb.Empty, error) {
	err := lotusService.CreateSectorTable(args.MinerId)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// CreateSectorPieceTable Create sector order table
func (g *GateWayServiceImpl) CreateSectorPieceTable(ctx context.Context, args *pb.Actor) (*emptypb.Empty, error) {
	err := lotusService.CreateSectorPieceTable(args.MinerId)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// CreateSectorQueueDetailTable Create sector task queue details table
func (g *GateWayServiceImpl) CreateSectorQueueDetailTable(ctx context.Context, args *pb.Actor) (*emptypb.Empty, error) {
	err := lotusService.CreateSectorQueueDetailTable(args.MinerId)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// CreateSectorLogTable Create sector log table
func (g *GateWayServiceImpl) CreateSectorLogTable(ctx context.Context, args *pb.Actor) (*emptypb.Empty, error) {
	err := lotusService.CreateSectorLogTable(args.MinerId)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// AddMiner add miner
func (g *GateWayServiceImpl) AddMiner(ctx context.Context, args *pb.MinerInfo) (*pb.ResponseMsg, error) {
	opClient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		global.ZC_LOG.Error("err.Error()")
		return nil, errors.New("Op：" + args.OpId + " not online！")
	}

	miner := &model.LotusMinerInfo{
		GateId:     global.GateWayID.String(),
		OpId:       args.OpId,
		LotusId:    args.LotusId,
		Ip:         args.Ip,
		Port:       define.MinerPort,
		SectorSize: args.SectorSize,
		Partitions: args.Partitions,
		Actor:      args.Actor,
		IsWdpost:   args.IsWdpost,
		IsWnpost:   args.IsWnpost,
		IsManage:   args.IsManage,
		AddType:    int(args.AddType),
	}

	miner.ID = uint(args.MinerId)
	miner.DeployStatus = int(define.DeployRunning)
	err := lotusService.DeployService.AddMiner(miner)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			err = fmt.Errorf("the server has repeatedly created miners")
		}
		global.ZC_LOG.Info(err.Error())
		return nil, err
	}

	//run
	sys, err := systemService.GetSysHostRecord(args.OpId)
	if err != nil {
		global.ZC_LOG.Info(err.Error())
		return &pb.ResponseMsg{Code: 500, Msg: "operation failed！"}, err
	}

	args.ReRaid = !sys.IsGroupArray

	err = systemService.UpdateSysHostRecordClassify(&system.SysHostRecord{UUID: args.OpId, HostClassify: config.HostMinerType, IsGroupArray: true})
	if err != nil {
		return nil, err
	}
	args.Actor = miner.Actor
	args.MinerId = uint64(miner.ID)
	_, err = opClient.RunNewMiner(ctx, args)
	if err != nil {
		global.ZC_LOG.Info(err.Error())
		return nil, err
	}

	var list []lotus.LotusStorageInfo
	list, err = lotusService.GetStorageInfoByActor(miner.Actor)
	if err != nil {
		global.ZC_LOG.Info("GetStorageInfoByActor error: " + err.Error())
	} else {

		if len(list) > 0 {
			for _, val := range list {
				if val.ColonyType != define.StorageTypeNFS {
					continue
				}
				opDir := utils.RestoreMountDir(val.NFSDisk, val.Ip)
				mountInfo := &pb.MountDiskInfo{OpIP: val.Ip, OpDir: opDir}
				if _, err = opClient.NodeMountDisk(context.Background(), mountInfo); err != nil {
					global.ZC_LOG.Info("NodeMountDisk error：" + err.Error())
					continue
				}
				storageInfo := &pb.StorageInfo{MountDir: opDir, NodeIp: miner.Ip, StorageIp: val.Ip}
				if _, err = opClient.AddNodeStorage(context.Background(), storageInfo); err != nil {
					global.ZC_LOG.Info("NodeMountDisk error：" + err.Error())
					continue
				}
			}
		}
	}

	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, err
}

// UpdateLotus Update lotus status
func (g *GateWayServiceImpl) UpdateLotus(ctx context.Context, args *pb.ConnectInfo) (*emptypb.Empty, error) {
	full, err := lotusService.GetLotus(args.Id)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	full.FinishAt = time.Now()
	full.DeployStatus = int(args.DeployStatus)
	full.SyncStatus = int(args.SyncStatus)
	full.RunStatus = int(args.RunStatus)
	full.SnapshotAt = time.Now()
	full.Token = args.Token
	full.ErrMsg = args.ErrMsg
	if err = lotusService.UpdateLotus(full); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateMiner Update miner status
func (g *GateWayServiceImpl) UpdateMiner(ctx context.Context, args *pb.ConnectInfo) (*emptypb.Empty, error) {
	miner, err := lotusService.GetMiner(args.Id)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if miner.DeployStatus != define.DeployFinish.Int() && args.DeployStatus == define.DeployFinish.Int32() {
		miner.FinishAt = time.Now()
	}

	miner.Token = args.Token
	miner.DeployStatus = int(args.DeployStatus)
	miner.RunStatus = int(args.RunStatus)
	miner.ErrMsg = args.ErrMsg
	if args.Actor != "" {
		miner.Actor = args.Actor
	}

	err = lotusService.AddMiner(&miner)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateBoost 更新boost状态
func (g *GateWayServiceImpl) UpdateBoost(ctx context.Context, args *pb.ConnectInfo) (*emptypb.Empty, error) {
	var boost lotus.LotusBoostInfo
	if args.Id > 0 {
		boost, _ = lotusService.GetBoost(args.Id)
	}
	if boost.DeployStatus != define.DeployFinish.Int() && args.DeployStatus == define.DeployFinish.Int32() {
		boost.FinishAt = time.Now()
	}

	boost.Token = args.Token
	boost.DeployStatus = int(args.DeployStatus)
	boost.RunStatus = int(args.RunStatus)
	boost.ErrMsg = args.ErrMsg
	boost.ID = uint(args.Id)

	if err := lotusService.AddBoost(&boost); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateWorker 更新worker状态
func (g *GateWayServiceImpl) UpdateWorker(ctx context.Context, args *pb.ConnectInfo) (*emptypb.Empty, error) {
	worker, err := lotusService.GetWorker(args.Id)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if worker.DeployStatus != define.DeployFinish.Int() && args.DeployStatus == define.DeployFinish.Int32() {
		worker.FinishAt = time.Now()
	}

	worker.DeployStatus = int(args.DeployStatus)
	worker.RunStatus = int(args.RunStatus)
	worker.ErrMsg = args.ErrMsg
	err = lotusService.UpdateWorker(&worker)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, err
}

// UpdateStorage 更新worker状态
func (g *GateWayServiceImpl) UpdateStorage(ctx context.Context, args *pb.ConnectInfo) (*emptypb.Empty, error) {
	worker, err := lotusService.GetStorage(args.Id)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if worker.DeployStatus != define.DeployFinish.Int() && args.DeployStatus == define.DeployFinish.Int32() {
		worker.FinishAt = time.Now()
	}

	worker.DeployStatus = int(args.DeployStatus)
	worker.ErrMsg = args.ErrMsg
	err = lotusService.UpdateStorage(&worker)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, err
}

// SetWorkerTask Update Worker Task Configuration
func (g *GateWayServiceImpl) SetWorkerTask(ctx context.Context, args *pb.WorkerConfig) (*pb.ResponseMsg, error) {
	var err error
	opClient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return nil, errors.New("Op：" + args.OpId + " not online！")
	}
	if _, err = opClient.SetPreNumber(ctx, &pb.PreNumber{P1: args.PreCount1, P2: args.PreCount2}); err != nil {
		return nil, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddWorker add worker
func (g *GateWayServiceImpl) AddWorker(ctx context.Context, args *pb.BatchWorker) (*emptypb.Empty, error) {
	if args == nil {
		return &emptypb.Empty{}, fmt.Errorf("BatchWorker is nil")
	}
	var wait sync.WaitGroup
	wait.Add(len(args.Host))
	var errmsg error
	for i := 0; i < len(args.Host); i++ {
		log.Println("add host：", args.Host[i].Ip)
		go func(j int) {
			defer func() {
				wait.Done()
			}()

			opClient, dis := global.OpClinets.GetOpClient(args.Host[j].OpId)
			if dis {
				errmsg = fmt.Errorf("Op：" + args.Host[j].OpId + " not online！")
				global.ZC_LOG.Error(errmsg.Error())
				return
			}
			info, err := lotusService.GetMiner(args.MinerId)
			if err != nil {
				errmsg = err
				global.ZC_LOG.Error(errmsg.Error())
				return
			}

			worker := &model.LotusWorkerInfo{
				GateId:       global.GateWayID.String(),
				OpId:         args.Host[j].OpId,
				Ip:           args.Host[j].Ip,
				DeployStatus: int(define.DeployRunning),
				MinerId:      args.MinerId,
			}
			worker.ID = uint(args.Host[j].Id)
			if err = lotusService.DeployService.AddWorker(worker); err != nil {
				if strings.Contains(err.Error(), "Duplicate") {
					errmsg = fmt.Errorf("this server creates duplicate workers %s", worker.Ip)
				}
				log.Println(err)
				return
			}

			err = systemService.UpdateSysHostRecordClassify(&system.SysHostRecord{UUID: args.Host[j].OpId, HostClassify: config.HostWorkerType, IsGroupArray: true})
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				return
			}
			w := &pb.WorkerInfo{
				Id:         uint64(worker.ID),
				Ip:         args.Host[j].Ip,
				MinerToken: info.Token,
				MinerIp:    info.Ip,
			}
			_, err = opClient.RunNewWorker(context.Background(), w)
			if err != nil {
				errmsg = fmt.Errorf(args.Host[j].OpId + " not online！" + err.Error())
				global.ZC_LOG.Error(err.Error())
				return
			}

			l := lotus.LoutsWorkerConfig{
				OpId:      args.Host[j].OpId,
				GateId:    args.GateId,
				WorkerId:  worker.ID,
				PreCount1: 13,
				PreCount2: 1,
				Actor:     info.Actor,
				IP:        args.Host[j].Ip,
				Port:      define.WorkerPort,
			}
			if err = lotusService.DispatchService.AddWrokerConfig(&l); err != nil {
				global.ZC_LOG.Error(err.Error())
				return
			}

			minerInfo, err := lotusService.GetMiner(args.MinerId)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
			} else {
				var list []lotus.LotusStorageInfo
				list, err = lotusService.GetStorageInfoByActor(minerInfo.Actor)
				if err != nil {
					global.ZC_LOG.Error(err.Error())
				} else {
					if len(list) > 0 {
						for _, val := range list {
							if val.ColonyType != define.StorageTypeNFS {
								continue
							}
							opDir := utils.RestoreMountDir(val.NFSDisk, val.Ip)
							mountInfo := &pb.MountDiskInfo{OpIP: val.Ip, OpDir: opDir}
							if _, err = opClient.NodeMountDisk(context.Background(), mountInfo); err != nil {
								global.ZC_LOG.Info("NodeMountDisk error：" + err.Error())
								continue
							}
						}
					}
				}
			}
		}(i)
	}
	wait.Wait()
	time.Sleep(time.Millisecond * 100)
	return &emptypb.Empty{}, errmsg
}

// WorkerMountNFS Worker mounting NFS
func (g *GateWayServiceImpl) WorkerMountNFS(ctx context.Context, args *pb.OpHostUUID) (*emptypb.Empty, error) {
	log.Println("WorkerMountNFS Begin")
	var err error
	var minerInfo lotus.LotusMinerInfo

	opClient, dis := global.OpClinets.GetOpClient(args.HostUUID)
	if dis {
		errmsg := fmt.Errorf("Op：" + args.HostUUID + " not online！")
		log.Println(errmsg.Error())
		return &emptypb.Empty{}, errmsg
	}
	if args.HostType == strconv.Itoa(config.HostWorkerType) {
		workerInfo, err := lotusService.DeployService.GetWorkerByOPId(args.HostUUID)
		if err != nil {
			log.Println(err.Error())
			return &emptypb.Empty{}, err
		}

		minerInfo, err = lotusService.GetMiner(workerInfo.MinerId)
	} else if args.HostType == strconv.Itoa(config.HostMinerType) {
		minerInfo, err = lotusService.DeployService.GetMinerByOpId(args.HostUUID)
		if err != nil {
			log.Println(err.Error())
			return &emptypb.Empty{}, err
		}
	} else {
		return &emptypb.Empty{}, errors.New("args hostType is error")
	}

	if err != nil {
		global.ZC_LOG.Error(err.Error())
	} else {
		var list []lotus.LotusStorageInfo
		list, err = lotusService.GetStorageInfoByActor(minerInfo.Actor)
		if err != nil {
			log.Println(err.Error())
		} else {
			if len(list) > 0 {
				for _, val := range list {
					go func(storage lotus.LotusStorageInfo) {
						if storage.ColonyType != define.StorageTypeNFS {
							return
						}
						opDir := utils.RestoreMountDir(storage.NFSDisk, storage.Ip)
						mountInfo := &pb.MountDiskInfo{OpIP: storage.Ip, OpDir: opDir}
						if _, err = opClient.NodeMountDisk(context.Background(), mountInfo); err != nil {
							log.Println("NodeMountDisk error：" + err.Error())
							return
						}
					}(val)
				}
			}
		}
	}
	return &emptypb.Empty{}, nil
}

// AddBoost Add boost
func (g *GateWayServiceImpl) AddBoost(ctx context.Context, args *pb.BoostInfo) (*emptypb.Empty, error) {
	opClient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return &emptypb.Empty{}, errors.New("Op：" + args.OpId + " not online！")
	}

	miner, err := lotusService.GetMiner(args.MinerId)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	lot, err := lotusService.GetLotus(args.LotusId)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	args.MinerIp = miner.Ip
	args.MinerToken = miner.Token
	args.Actor = miner.Actor
	args.LotusIp = lot.Ip
	args.LotusToken = lot.Token

	boost := &model.LotusBoostInfo{
		GateId:        args.GateId,
		OpId:          args.OpId,
		LotusId:       args.LotusId,
		MinerId:       args.MinerId,
		LanIp:         args.LanIp,
		LanPort:       args.LanPort,
		InternetIp:    args.InternetIp,
		InternetPort:  args.InternetPort,
		DeployStatus:  define.DeployRunning.Int(),
		RunStatus:     define.RunStatusStop.Int(),
		NetworkType:   int(args.NetworkType),
		DcQuotaWallet: args.DcQuotaWallet,
	}
	boost.ID = uint(args.Id)

	if err = lotusService.DeployService.AddBoost(boost); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			err = fmt.Errorf("the server has repeatedly created Boost")
		}
		return &emptypb.Empty{}, err
	}
	args.Id = uint64(boost.ID)
	return opClient.RunBoost(ctx, args)
}

// UpdateSectorStatus Modify sector status
func (g *GateWayServiceImpl) UpdateSectorStatus(ctx context.Context, args *pb.SectorStatus) (*pb.ResponseMsg, error) {
	if err := lotusService.UpdateSectorStatus(args); err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddSectorTicket Add sector ticket
func (g *GateWayServiceImpl) AddSectorTicket(ctx context.Context, args *pb.SectorTicket) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorTicket(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddSectorCommDR Add sector P2 information
func (g *GateWayServiceImpl) AddSectorCommDR(ctx context.Context, args *pb.SectorCommDR) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorCommDR(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddSectorWaitSeed Add sector WaitSeed information
func (g *GateWayServiceImpl) AddSectorWaitSeed(ctx context.Context, args *pb.SectorSeed) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorWaitSeed(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddSectorCommit2 Add sector C2 completion result
func (g *GateWayServiceImpl) AddSectorCommit2(ctx context.Context, args *pb.SectorProof) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorCommit2(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddSectorPreCID Add sector P2 message ID
func (g *GateWayServiceImpl) AddSectorPreCID(ctx context.Context, args *pb.SectorCID) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorPreCID(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddSectorCommitCID Add sector C2 message ID
func (g *GateWayServiceImpl) AddSectorCommitCID(ctx context.Context, args *pb.SectorCID) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorCommitCID(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// AddSectorLog Add sector log information
func (g *GateWayServiceImpl) AddSectorLog(ctx context.Context, args *pb.SectorLog) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorLog(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// UpdateSectorLog Modifying sector log information
func (g *GateWayServiceImpl) UpdateSectorLog(ctx context.Context, args *pb.SectorLog) (*pb.ResponseMsg, error) {
	err := lotusService.EndSectorLog(args.ID, args.Sector.Miner, args.ErrorMsg)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// GetRunningCount Obtain the number of running tasks
func (g *GateWayServiceImpl) GetRunningCount(ctx context.Context, args *pb.OpTask) (*pb.TaskCount, error) {
	cli, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return nil, errors.New("Op：" + args.OpId + " not online！")
	}
	return cli.GetRunningCount(ctx, &pb.String{Value: args.TType})
}

// GetRunningList Get Run Task Data
func (g *GateWayServiceImpl) GetRunningList(ctx context.Context, args *pb.OpTask) (*pb.TaskInfoList, error) {
	if args.OpId != "" {
		cli, dis := global.OpClinets.GetOpClient(args.OpId)
		if dis {
			return nil, errors.New("Op：" + args.OpId + " not online！")
		}
		return cli.GetRunningList(ctx, &pb.String{Value: args.TType})
	}
	var err error

	ops := global.OpClinets.OnLineList()
	taskList := make([]*pb.TaskInfoList, len(ops))
	var wait sync.WaitGroup
	for i, op := range ops {
		wait.Add(1)
		go func(index int, client pb.OpServiceClient) {
			defer wait.Done()
			taskList[index], err = client.GetRunningList(ctx, &pb.String{Value: args.TType})
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				return
			}
		}(i, op)
	}
	wait.Wait()
	var taskInfos []*pb.TaskInfo

	for _, task := range taskList {
		if task == nil {
			continue
		}
		taskInfos = append(taskInfos, task.Tasks...)
	}
	return &pb.TaskInfoList{Tasks: taskInfos}, nil
}

// GetStorageByActor Get storage
func (g *GateWayServiceImpl) GetStorageByActor(ctx context.Context, args *pb.String) (*pb.LinkList, error) {
	//var res *pb.LinkList
	colony, err := systemService.GetColony(args.Value)
	if err != nil {
		return nil, err
	}
	var list []*pb.FilParam
	if colony.ColonyType == define.StorageTypeWorker {
		list, err = lotusService.GetStorageByActor(args.Value)
		if err != nil {
			return nil, err
		}
	} else {
		list, err = lotusService.GetWorkerByActor(args.Value)
		if err != nil {
			return nil, err
		}
	}
	return &pb.LinkList{Links: list, StorageType: int32(colony.ColonyType)}, err
}

// GetActorTaskQueue Get task queue
func (g *GateWayServiceImpl) GetActorTaskQueue(ctx context.Context, args *pb.String) (*pb.TaskQueues, error) {
	//var res *pb.LinkList
	queues, err := lotusService.GetRunningTaskQueue(args.Value)
	if err != nil {
		return nil, err
	}

	return &pb.TaskQueues{Queues: queues}, err
}

// AddCompleteCountByID Add task queue completion quantity
func (g *GateWayServiceImpl) AddCompleteCountByID(ctx context.Context, args *pb.TaskQueue) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.AddCompleteCountByID(args.ID)
}

// AddRunCountByID Get the number of tasks in the queue being executed
func (g *GateWayServiceImpl) AddRunCountByID(ctx context.Context, args *pb.TaskQueue) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.AddRunCountByID(args.ID)
}

// AddSectorQueueDetail Add queue details data
func (g *GateWayServiceImpl) AddSectorQueueDetail(ctx context.Context, args *pb.SectorQueueDetail) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.AddSectorQueueDetail(args)
}

// AddSectorPiece Add sector order information
func (g *GateWayServiceImpl) AddSectorPiece(ctx context.Context, args *pb.SectorPiece) (*pb.ResponseMsg, error) {
	err := lotusService.AddSectorPiece(args)
	if err != nil {
		return &pb.ResponseMsg{Code: 400, Msg: "operation failed！"}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful！"}, nil
}

// GetWaitImportDeal Obtain order data
func (g *GateWayServiceImpl) GetWaitImportDeal(ctx context.Context, args *pb.DealParam) (*pb.DealList, error) {
	res, err := lotusService.GetWaitImportDeal(args)
	if err != nil {
		return nil, err
	}
	return &pb.DealList{Deals: res}, nil
}

// EditQueueDetailStatus Get modified order status
func (g *GateWayServiceImpl) EditQueueDetailStatus(ctx context.Context, args *pb.EditStatus) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.EditTaskQueueDetailStatus(request1.ActorIdStatus{Status: int(args.Status), Actor: args.Actor, ID: uint(args.Id)}, "")
}

// RestartMiner Start stop service
func (g *GateWayServiceImpl) RestartMiner(ctx context.Context, args *pb.MinerRun) (*emptypb.Empty, error) {

	oclient, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return &emptypb.Empty{}, errors.New("Op：" + args.OpId + " not online！")
	}
	return oclient.RunMiner(ctx, args)
}

// SectorStorage Sector storage path
func (g *GateWayServiceImpl) SectorStorage(ctx context.Context, args *pb.SectorActorID) (*pb.SectorPaths, error) {
	strs, err := lotusrpc.FullApi.SectorStorage(args.Token, args.Ip, args.Miner, args.Number)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.String, len(strs))
	for i, str := range strs {
		res[i] = &pb.String{Value: str}
	}

	return &pb.SectorPaths{Strs: res}, nil
}

func (g *GateWayServiceImpl) AddColony(ctx context.Context, args *pb.Colony) (*pb.Colony, error) {
	var colony system.SysColony
	colony.ColonyName = args.ColonyName
	colony.ColonyType = int(args.ColonyType)
	return &pb.Colony{ID: uint64(colony.ID)}, systemService.AddColony(&colony)
}

func (g *GateWayServiceImpl) GetColony(ctx context.Context, args *pb.Actor) (*pb.Colony, error) {
	sysCol, err := systemService.GetColony(args.MinerId)
	if err != nil {
		return nil, err
	}
	return &pb.Colony{ColonyName: sysCol.ColonyName, ID: uint64(sysCol.ID), ColonyType: int32(sysCol.ColonyType)}, nil
}

func (g *GateWayServiceImpl) AddStorage(ctx context.Context, args *pb.BatchStroage) (*emptypb.Empty, error) {
	if args == nil {
		return &emptypb.Empty{}, fmt.Errorf("BatchStroage is nil")
	}

	minerList, err := lotusService.DeployService.GetMinerListByActor(args.Node)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return nil, err
	}

	workerList, err := lotusService.DeployService.GetWorkerByActor(args.Node)
	if err != nil {
		global.ZC_LOG.Error(err.Error())
		return nil, err
	}

	for i := 0; i < len(args.Host); i++ {
		go func(j int) {
			var errMsg string

			var storage = &model.LotusStorageInfo{
				OpId:         args.Host[j].OpId,
				GateId:       global.GateWayID.String(),
				Ip:           args.Host[j].Ip,
				DeployStatus: define.DeployRunning.Int(),
				ColonyName:   args.Node,
				ColonyType:   int(args.StorageType),
			}
			defer func() {

				storageInfo, err := lotusService.DeployService.GetStorageByOpID(args.Host[j].OpId)
				if err != nil && storageInfo.ID == 0 {
					log.Println("GetStorageByOpID: ", err)
					storage.ErrMsg = errMsg
					if err := lotusService.DeployService.AddStorage(storage); err != nil {
						log.Println("AddStorage: ", err)
						return
					}
					fmt.Println("storage info:", storage)
				} else {
					storage.ID = storageInfo.ID
				}

				var deployStatus int32
				if errMsg == "" {
					deployStatus = define.DeployFinish.Int32()
				} else {
					deployStatus = define.DeployFail.Int32()
				}

				if args.StorageType == define.StorageTypeNFS {
					if _, err = g.UpdateStorage(context.Background(), &pb.ConnectInfo{Id: uint64(storage.ID), Token: "", DeployStatus: deployStatus, RunStatus: define.RunStatusRunning.Int32()}); err != nil {
						log.Println("UpdateStorage: ", err)
						return
					}
				}
				if deployStatus == define.DeployFinish.Int32() {
					if err := systemService.UpdateSysHostRecordClassify(&system.SysHostRecord{UUID: args.Host[j].OpId, HostClassify: config.HostStorageType, IsGroupArray: false}); err != nil {
						errMsg = "UpdateSysHostRecordClassify error：" + err.Error()
						log.Println(errMsg)
						return
					}
				}
			}()
			opClient, dis := global.OpClinets.GetOpClient(args.Host[j].OpId)
			if dis {
				errMsg = "Op：" + args.Host[j].OpId + " not online！"
				log.Println(errMsg)
				return
			}
			switch args.StorageType {
			case define.StorageTypeWorker:
				id, err := strconv.ParseUint(args.Node, 10, 64)
				if err != nil {
					errMsg = "minerId error：" + args.Node
					log.Println(errMsg)
					return
				}
				info, err := lotusService.GetMiner(id)
				if err != nil {
					errMsg = "Get Miner error：" + err.Error()
					log.Println(errMsg)
					return
				}
				storage.ColonyName = info.Actor
				if err := lotusService.DeployService.AddStorage(storage); err != nil {
					log.Println("AddStorage: ", err)
					return
				}
				worker := &pb.WorkerInfo{
					Id:         args.Host[j].Id,
					Ip:         args.Host[j].Ip,
					MinerToken: info.Token,
					MinerIp:    info.Ip,
				}
				_, err = opClient.RunNewStorage(context.Background(), worker)
				if err != nil {
					errMsg = "RunNewStorage：" + err.Error()
					log.Println(errMsg)
					return
				}
			case define.StorageTypeNFS:
				mountInfo := &pb.MountDiskInfo{OpIP: args.Host[j].Ip}
				mountDir, err := opClient.NodeAddShareDir(context.Background(), mountInfo)
				if err != nil {
					errMsg = "NodeAddShareDir error：" + err.Error()
					log.Println(errMsg)
					return
				}
				if len(mountDir.Value) == 0 {
					errMsg = "Storage machine disk sharing initialization failed " + args.Host[j].OpId
					log.Println(errMsg)
					return
				}
				storage.NFSDisk = utils.DealMountDir(mountDir.Value, args.Host[j].Ip)
				if len(minerList) > 0 {
					if err := lotusService.DeployService.AddStorage(storage); err != nil {
						log.Println("AddStorage: ", err)
						return
					}
					for _, miner := range minerList {
						nodeClient, dis := global.OpClinets.GetOpClient(miner.OpId)
						if dis {
							errMsg = "Association Node：" + miner.OpId + "  not online！"
							log.Println(errMsg)
							return
						}
						mountInfo.OpDir = mountDir.Value
						if _, err = nodeClient.NodeMountDisk(context.Background(), mountInfo); err != nil {
							errMsg = "NodeMountDisk error：" + err.Error()
							log.Println(errMsg)
							return
						}
						storageInfo := &pb.StorageInfo{MountDir: mountDir.Value, NodeIp: miner.Ip, StorageIp: args.Host[j].Ip}
						if _, err = nodeClient.AddNodeStorage(context.Background(), storageInfo); err != nil {
							errMsg = "NodeMountDisk error：" + err.Error()
							log.Println(errMsg)
							return
						}
					}
				}
				if len(workerList) > 0 {
					for _, worker := range workerList {
						// 链接worker主机
						nodeClient, dis := global.OpClinets.GetOpClient(worker.Opid)
						if dis {
							errMsg = "Association Node：" + worker.Opid + "  not online！"
							log.Println(errMsg)
							return
						}
						// 节点主机挂载存储机nfs
						mountInfo.OpDir = mountDir.Value
						if _, err = nodeClient.NodeMountDisk(context.Background(), mountInfo); err != nil {
							errMsg = "NodeMountDisk error：" + err.Error()
							log.Println(errMsg)
							return
						}
					}
				}
			}
		}(i)
	}
	return &emptypb.Empty{}, nil
}

func (g *GateWayServiceImpl) GetSectorStatus(ctx context.Context, args *pb.SectorID) (*pb.SectorStatus, error) {
	sector, err := lotusService.GetSectorInfo(args.Miner, args.Number)
	if err != nil {
		return nil, err
	}
	return &pb.SectorStatus{Sector: args, Status: sector.SectorStatus, Type: int32(sector.SectorType), Size: sector.SectorSize}, nil
}

func (g *GateWayServiceImpl) OpLocalSectors(ctx context.Context, args *pb.OpMiner) (*pb.SectorCount, error) {
	client, dis := global.OpClinets.GetOpClient(args.OpId)
	if dis {
		return nil, errors.New("Op：" + args.OpId + " not online！")
	}
	list, err := client.LocalSectors(ctx, &pb.Actor{MinerId: args.Miner})
	if err != nil {
		return nil, err
	}
	var statusMap = make(map[string]int32)
	var statusStore = make(map[string]int32)
	var statusMsg = make(map[string]int32)
	miner, merr := lotusService.GetMinerByActor(args.Miner)
	if merr != nil {
		log.Println(merr)
	}

	for _, sector := range list.Sectors {
		info, err := lotusService.GetSectorInfo(args.Miner, sector.Number)
		if err != nil {
			return nil, err
		}
		statusMap[info.SectorStatus]++
		if merr == nil {
			exist, _ := lotusrpc.FullApi.SectorStoreFile(miner.Token, miner.Ip, utils.MinerActorID(args.Miner), sector.Number)
			if exist >= (4 | 2) { //1 unsealed 2cache,4sealed
				statusStore[info.SectorStatus]++
			}
		}
		if info.PreCid != "" {
			statusMsg[info.SectorStatus]++
		}
	}
	var counts []*pb.StatusCount
	for status, count := range statusMap {
		counts = append(counts, &pb.StatusCount{Status: status, Count: count, Store: statusStore[status], PreMsg: statusMsg[status]})
	}

	return &pb.SectorCount{Sectors: counts}, nil
}

func (g *GateWayServiceImpl) CheckLotusHeart(ctx context.Context, args *pb.String) (*pb.String, error) {
	define.Lh.LotusLock.Lock()
	defer define.Lh.LotusLock.Unlock()

	define.Lh.LotusMap[args.Value] = time.Now()
	return &pb.String{}, nil
}

func (g *GateWayServiceImpl) CheckMinerHeart(ctx context.Context, args *pb.String) (*pb.String, error) {
	define.Lh.MinerLock.Lock()
	defer define.Lh.MinerLock.Unlock()

	define.Lh.MinerMap[args.Value] = time.Now()
	return &pb.String{}, nil
}

func (g *GateWayServiceImpl) CheckOWorkerHeart(ctx context.Context, args *pb.String) (*pb.String, error) {
	define.Lh.WorkerLock.Lock()
	defer define.Lh.WorkerLock.Unlock()

	define.Lh.WorkerMap[args.Value] = time.Now()
	return &pb.String{}, nil
}
