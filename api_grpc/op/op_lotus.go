package op

import (
	"context"
	"errors"
	"github.com/multiformats/go-base32"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/service/lotus/oplocal"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path/filepath"
)

// ClearWorker Clear task
func (g *OpServiceImpl) ClearWorker(ctx context.Context, args *pb.RequestOp) (*pb.ResponseMsg, error) {
	lotusService.ClearWorker(define.PathIpfsWorker)
	return &pb.ResponseMsg{Code: 200, Msg: "Operation successful"}, nil
}

// RunNewLotus Run lotus
func (p *OpServiceImpl) RunNewLotus(ctx context.Context, args *pb.LotusInfo) (*emptypb.Empty, error) {
	err := lotusService.RunNewLotus(args)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	log.Println("Copy wallet mode：", args.WalletNewMode, args.Wallets)
	if define.WalletNewMode(args.WalletNewMode) == define.WalletCopy {
		filemap := make(map[string]map[string]struct{})
		for _, v := range args.Wallets {
			fileName := base32.RawStdEncoding.EncodeToString([]byte(define.WalletPrefix + v.Address))
			if _, ok := filemap[v.OpId]; !ok {
				filemap[v.OpId] = make(map[string]struct{})
			}
			filemap[v.OpId][fileName] = struct{}{}
			log.Println("Copy wallet：", fileName, v.Address)
		}
		//Copy wallet files
		for OpId, nameMap := range filemap {
			var FileInfos []*pb.FileInfo
			for name, _ := range nameMap {
				FileInfos = append(FileInfos, &pb.FileInfo{FileName: name})
			}
			in := pb.SynFileInfo{
				FromPath: define.PathIpfsLotusKeystore,
				ToPath:   define.PathIpfsLotusKeystore,
				Ip:       args.Ip,
				Port:     define.OpPort,
				OpId:     OpId,
				ToOpId:   args.OpId,
				FileInfo: FileInfos,
			}
			res, err := global.OpToGatewayClient.SysFilePoint(ctx, &in)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				continue
			}
			if res.Code != 200 {
				global.ZC_LOG.Error(res.Msg)
				continue
			}
			for _, v := range FileInfos {
				if err = os.Chmod(filepath.Join(define.PathIpfsLotusKeystore, v.FileName), 0600); err != nil {
					global.ZC_LOG.Error(err.Error())
				}
			}
		}
	}
	return &emptypb.Empty{}, nil
}

// GetWalletList Run to obtain wallet list
func (p *OpServiceImpl) GetWalletList(ctx context.Context, args *pb.RequestConnect) (*pb.WalletList, error) {
	if args.Token == "" {
		return nil, errors.New("token cannot be empty！")
	}
	if args.Ip == "" {
		return nil, errors.New("ip cannot be empty！")
	}
	return lotusService.GetWalletList(args)
}

// RunNewMiner Run a new miner
func (p *OpServiceImpl) RunNewMiner(ctx context.Context, args *pb.MinerInfo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.RunNewMiner(args)
}

// RunMiner Run miner
func (p *OpServiceImpl) RunMiner(ctx context.Context, args *pb.MinerRun) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.RunMiner(args)
}

// RunNewWorker run worker
func (p *OpServiceImpl) RunNewWorker(ctx context.Context, args *pb.WorkerInfo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.RunNewWorker(args)
}

// RunNewStorage run storage
func (p *OpServiceImpl) RunNewStorage(ctx context.Context, args *pb.WorkerInfo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.RunNewStorage(args)
}

// AddNodeStorage Add the corresponding disk directory to the storage configuration file
func (p *OpServiceImpl) AddNodeStorage(ctx context.Context, args *pb.StorageInfo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.AddNodeStorage(args)
}

// RunWorker  run worker
func (p *OpServiceImpl) RunWorker(ctx context.Context, args *pb.FilParam) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.RunWorker(args)
}

// RunBoost run boost
func (p *OpServiceImpl) RunBoost(ctx context.Context, args *pb.BoostInfo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.RunBoost(args)
}

// WorkerTaksRunList Get the worker executing the task
func (p *OpServiceImpl) WorkerTaksRunList(ctx context.Context, args *pb.String) (*pb.TaskList, error) {
	var reply pb.TaskList
	oplocal.Tasking.Lock.RLock()
	for _, ts := range oplocal.Tasking.Run {
		for _, v := range ts {
			reply.Tasks = append(reply.Tasks, v)
		}
	}
	oplocal.Tasking.Lock.RUnlock()
	return &reply, nil
}

// GetRunningCount Obtain the number of tasks that the worker is running
func (p *OpServiceImpl) GetRunningCount(ctx context.Context, args *pb.String) (*pb.TaskCount, error) {
	count := oplocal.Tasking.GetRunCount(args.Value)
	return &pb.TaskCount{TCount: int32(count), TType: args.Value}, nil
}

// GetRunningList Get task data that the worker is running
func (p *OpServiceImpl) GetRunningList(ctx context.Context, args *pb.String) (*pb.TaskInfoList, error) {
	return oplocal.Tasking.GetRunList(args.Value), nil
}

// Ok Okay, determine if the task can be completed
func (p *OpServiceImpl) Ok(ctx context.Context, args *pb.Task) (*wrapperspb.BoolValue, error) {
	return &wrapperspb.BoolValue{Value: lotusService.Ok(args)}, nil
}

// OkNew Determine how many tasks can be completed
func (p *OpServiceImpl) OkNew(ctx context.Context, args *pb.MinerSize) (*pb.TaskCan, error) {
	return lotusService.OkNew(args), nil
}

// AddRunning Increase the number of received tasks
func (p *OpServiceImpl) AddRunning(ctx context.Context, args *pb.Task) (*pb.ResponseMsg, error) {
	var res pb.ResponseMsg
	if err := lotusService.AddRunning(args); err != nil {
		res.Msg = err.Error()
		return nil, err
	}
	res.Code = 200
	res.Msg = "Added successfully！"
	return &res, nil
}

// SubRunning Reduce accepted tasks
func (p *OpServiceImpl) SubRunning(ctx context.Context, args *pb.Task) (*pb.ResponseMsg, error) {
	var res pb.ResponseMsg
	if err := lotusService.SubRunning(args); err != nil {
		res.Msg = err.Error()
		return nil, err
	}
	res.Code = 200
	res.Msg = "Reduce Success！"
	return &res, nil
}

// ResetWorkerRunning Reset Received Tasks
func (p *OpServiceImpl) ResetWorkerRunning(ctx context.Context, args *pb.WorkerTasks) (*emptypb.Empty, error) {
	if err := lotusService.ResetWorkerRunning(args); err != nil {
		return nil, err
	}

	_, err := global.OpToGatewayClient.CheckOWorkerHeart(ctx, &pb.String{Value: args.Info.Ip})
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

// SetPreNumber Update available task configurations
func (p *OpServiceImpl) SetPreNumber(ctx context.Context, args *pb.PreNumber) (*pb.ResponseMsg, error) {
	var res pb.ResponseMsg
	if args.P1 >= 0 {
		oplocal.PreCount1 = int(args.P1)
	}
	if args.P2 >= 0 {
		oplocal.PreCount2 = int(args.P2)
	}
	res.Code = 200
	res.Msg = "Update success！"
	return &res, nil
}

// UpdateSectorStatus Modify sector status
func (p *OpServiceImpl) UpdateSectorStatus(ctx context.Context, args *pb.SectorStatus) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.UpdateSectorStatus(ctx, args)
}

// AddSectorTicket Add sector ticket
func (p *OpServiceImpl) AddSectorTicket(ctx context.Context, args *pb.SectorTicket) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorTicket(ctx, args)
}

// AddSectorCommDR Add sector P2 information
func (p *OpServiceImpl) AddSectorCommDR(ctx context.Context, args *pb.SectorCommDR) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorCommDR(ctx, args)
}

// AddSectorWaitSeed Add sector WaitSeed information
func (p *OpServiceImpl) AddSectorWaitSeed(ctx context.Context, args *pb.SectorSeed) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorWaitSeed(ctx, args)
}

// AddSectorCommit2 Add sector C2 completion result
func (p *OpServiceImpl) AddSectorCommit2(ctx context.Context, args *pb.SectorProof) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorCommit2(ctx, args)
}

// AddSectorPreCID Add sector P2 message ID
func (p *OpServiceImpl) AddSectorPreCID(ctx context.Context, args *pb.SectorCID) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorPreCID(ctx, args)
}

// AddSectorCommitCID Add sector C2 message ID
func (p *OpServiceImpl) AddSectorCommitCID(ctx context.Context, args *pb.SectorCID) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorCommitCID(ctx, args)
}

// AddSectorPiece Add sector order information
func (p *OpServiceImpl) AddSectorPiece(ctx context.Context, args *pb.SectorPiece) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorPiece(ctx, args)
}

// AddSectorLog Add sector log information
func (p *OpServiceImpl) AddSectorLog(ctx context.Context, args *pb.SectorLog) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.AddSectorLog(ctx, args)
}

// UpdateSectorLog Modifying sector log information
func (p *OpServiceImpl) UpdateSectorLog(ctx context.Context, args *pb.SectorLog) (*pb.ResponseMsg, error) {
	return global.OpToGatewayClient.UpdateSectorLog(ctx, args)
}

// GetStorageByActor Get storage
func (p *OpServiceImpl) GetStorageByActor(ctx context.Context, args *pb.String) (*pb.LinkList, error) {
	return global.OpToGatewayClient.GetStorageByActor(ctx, args)
}

// EditMinerApCount Miner heartbeat check
func (p *OpServiceImpl) EditMinerApCount(ctx context.Context, actor *pb.PledgeParam) (*emptypb.Empty, error) {
	log.Println("EditMinerApCount:", actor.Actor, actor.ApCount, actor.IsManage, actor.Info)
	if actor.ApCount > 2 {
		return &emptypb.Empty{}, nil
	}
	if actor.IsManage {
		if err := lotusService.NewPledgeTask(ctx, actor.Actor, actor.Info); err != nil {
			return nil, err
		}
	}

	_, err := global.OpToGatewayClient.CheckMinerHeart(ctx, &pb.String{Value: actor.Info.Ip})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (p *OpServiceImpl) AcquireSector(ctx context.Context, sector *pb.SectorRef) (*pb.SectorPath, error) {
	sectorPath, err := lotusService.FindSealStorage(define.PathIpfsWorker, sector)
	if err == nil {
		oplocal.SealSectorPath.Push(utils.SectorNumString(sector.Id.Miner, sector.Id.Number), sectorPath.DiskPath)
	}
	log.Println("AcquireSector:", sector.Id.Miner, sector.Id.Number, sector.ProofType, sectorPath)
	return sectorPath, err
}

func (p *OpServiceImpl) GetColony(ctx context.Context, args *pb.Actor) (*pb.Colony, error) {
	return global.OpToGatewayClient.GetColony(ctx, args)
}

func (p *OpServiceImpl) SealingAbort(ctx context.Context, args *pb.FilRestWorker) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, lotusService.SealingAbort(args)
}

func (p *OpServiceImpl) LocalSectors(ctx context.Context, args *pb.Actor) (*pb.SectorList, error) {
	sectors := lotusService.RangeSectors(define.PathIpfsWorker, args.MinerId)
	var sector = make([]*pb.SectorID, len(sectors))
	for i, v := range sectors {
		sector[i] = &pb.SectorID{Miner: args.MinerId, Number: v}
	}
	return &pb.SectorList{Sectors: sector}, nil
}

func (p *OpServiceImpl) GetGateWayFile(ctx context.Context, args *pb.String) (*pb.String, error) {
	return global.OpToGatewayClient.GetGateWayFile(ctx, args)
}

func (p *OpServiceImpl) CheckLotusHeart(ctx context.Context, args *pb.String) (*pb.String, error) {
	return global.OpToGatewayClient.CheckLotusHeart(ctx, args)
}
