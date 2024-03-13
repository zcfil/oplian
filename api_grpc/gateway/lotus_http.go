package gateway

import (
	"context"
	"fmt"
	"oplian/lotusrpc"
	"oplian/service/pb"
)

// StateMinerInfo Obtain information on the miner chain
func (g *GateWayServiceImpl) StateMinerInfo(ctx context.Context, args *pb.FilParam) (*pb.ActorWallet, error) {
	actor, err := lotusrpc.FullApi.StateMinerInfo(args.Token, args.Ip, args.Param)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain node information：%s", err.Error())
	}
	return &pb.ActorWallet{Owner: actor.Owner, Worker: actor.Worker, Control: actor.ControlAddresses}, nil
}

func (g *GateWayServiceImpl) StateAccountKey(ctx context.Context, args *pb.FilParam) (*pb.String, error) {
	addr, err := lotusrpc.FullApi.StateAccountKey(args.Token, args.Ip, args.Param)
	if err != nil {
		return &pb.String{Value: addr}, err
	}
	return &pb.String{Value: addr}, nil
}

// WalletBalance Wallet balance
func (g *GateWayServiceImpl) WalletBalance(ctx context.Context, args *pb.FilParam) (*pb.Wallet, error) {
	balance, err := lotusrpc.FullApi.WalletBalance(args.Token, args.Ip, args.Param)
	if err != nil {
		return nil, err
	}
	return &pb.Wallet{Address: args.Param, Balance: balance}, nil
}

// SealingSchedDiag Obtain miner ready task information
func (g *GateWayServiceImpl) SealingSchedDiag(ctx context.Context, args *pb.FilParam) (*pb.SchedDiagRequestInfo, error) {
	diags, err := lotusrpc.FullApi.SealingSchedDiag(args.Token, args.Ip)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain node ready task information：%s", err.Error())
	}
	var requests []*pb.SchedDiagInfo
	for _, v := range diags.Requests {
		s := &pb.SchedDiagInfo{
			SchedId:  v.SchedId,
			SectorID: v.Sector.Number,
			TaskType: v.TaskType,
			Priority: int32(v.Priority),
		}
		requests = append(requests, s)
	}
	return &pb.SchedDiagRequestInfo{Requests: requests}, nil
}

// QueryAsk Obtain order data
func (g *GateWayServiceImpl) QueryAsk(ctx context.Context, args *pb.QueryParam) (*pb.AskInfo, error) {
	actorInfo, err := lotusrpc.FullApi.StateMinerInfo(args.LotusToken, args.LotusIp, args.Param)
	if err != nil {
		return nil, err
	}
	res, err := lotusrpc.FullApi.ClientQueryAsk(args.LotusToken, args.LotusIp, actorInfo.PeerId, args.Param)
	if err != nil {
		return nil, err
	}
	return &pb.AskInfo{Miner: res.Miner, Price: res.Price, VerifiedPrice: res.VerifiedPrice, MinPieceSize: res.MinPieceSize, MaxPieceSize: res.MaxPieceSize}, nil
}

// QueryDataCap Obtain DC balance
func (g *GateWayServiceImpl) QueryDataCap(ctx context.Context, args *pb.QueryParam) (*pb.String, error) {
	res, err := lotusrpc.FullApi.StateVerifiedClientStatus(args.LotusToken, args.LotusIp, args.Param)
	if err != nil {
		return nil, err
	}
	return &pb.String{Value: res}, nil
}

// StateMinerSectorCount Obtain mining sector information
func (g *GateWayServiceImpl) StateMinerSectorCount(ctx context.Context, args *pb.FilParam) (*pb.MinerSectors, error) {
	sectors, err := lotusrpc.FullApi.StateMinerSectorCount(args.Token, args.Ip, args.Param)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain node sector information：%s", err.Error())
	}
	return &pb.MinerSectors{Live: sectors.Live, Active: sectors.Active, Faulty: sectors.Faulty}, nil
}

// StateMinerPower Obtain miner computing power
func (g *GateWayServiceImpl) StateMinerPower(ctx context.Context, args *pb.FilParam) (*pb.Power, error) {
	power, err := lotusrpc.FullApi.StateMinerPower(args.Token, args.Ip, args.Param)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain node sector information：%s", err.Error())
	}
	return &pb.Power{MinerPower: power.MinerPower, TotalPower: power.TotalPower}, nil
}
