package gateway

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/lotus/request"
	"oplian/service/lotus/deploy"
	"oplian/service/lotus/deploy/cgo"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"path"
)

// Commit2TaskResult C2 Task Results
func (p *GateWayServiceImpl) Commit2TaskResult(ctx context.Context, args *pb.SectorID) (*pb.String, error) {

	res := &pb.String{}
	gateWayFilePathC2 := define.PathIpfsData + "/c2task" + define.GateWayCsPathC2

	param := request.C2TaskInfo{}
	param.DelType = utils.ONE
	if utils.IsNull(args.Miner) {
		err := deploy.WorkerClusterServiceApi.DelC2TaskInfo(param)
		if err != nil {
			return res, err
		} else {
			os.RemoveAll(gateWayFilePathC2)
		}

	} else {

		miner, _ := utils.FileCoinStrToUint64(args.Miner)
		fileName := path.Join(gateWayFilePathC2, fmt.Sprintf("s-t0%d-%d.json", miner, args.Number))
		param.Miner = fmt.Sprintf("t0%d", miner)
		param.Number = int(args.Number)
		err := deploy.WorkerClusterServiceApi.DelC2TaskInfo(param)
		if err != nil {
			return res, err
		} else {
			os.Remove(fileName)
		}
	}

	return res, nil

}

// CompleteCommit2 Complete C2 task
func (p *GateWayServiceImpl) CompleteCommit2(ctx context.Context, args *pb.FileInfo) (*pb.ResponseMsg, error) {

	res := &pb.ResponseMsg{}
	//更改c2任务状态
	return res, deploy.WorkerClusterServiceApi.ModifyC2TaskStatus(args)

}

// ModifySectorStatus Change task status
func (g *GateWayServiceImpl) ModifySectorStatus(ctx context.Context, args *pb.TaskStatus) (*pb.String, error) {

	log.Println("gateway ModifySectorStatus success")
	err := cgo.ModifySectorStatus(args)
	if err != nil {
		return &pb.String{}, err
	}

	return &pb.String{}, nil
}

// C2FileSynStatus C2 task document synchronization status
func (g *GateWayServiceImpl) C2FileSynStatus(ctx context.Context, args *pb.C2SectorID) (*pb.String, error) {

	res := &pb.String{}
	str, err := deploy.WorkerClusterServiceApi.C2FileSynStatus(args)
	if err != nil {
		return res, err
	}
	return &pb.String{Value: str}, nil
}

// AddC2Task Add C2 task
func (g *GateWayServiceImpl) AddC2Task(ctx context.Context, args *pb.SectorID) (*pb.String, error) {

	str, err := lotusService.AddC2Task(request.C2TaskInfo{Miner: args.Miner, Number: int(args.Number)})
	if err != nil {
		return &pb.String{}, err
	}

	return &pb.String{Value: str}, nil
}

// DelC2Task Update C2 task
func (g *GateWayServiceImpl) DelC2Task(ctx context.Context, args *pb.SectorID) (*pb.String, error) {
	return &pb.String{}, lotusService.DelC2TaskInfo(request.C2TaskInfo{Miner: args.Miner, Number: int(args.Number)})
}

// C2FileSynLotus C2 task file synchronization to lotus
func (g *GateWayServiceImpl) C2FileSynLotus(ctx context.Context, args *pb.FileInfo) (*pb.String, error) {

	f := make([]*pb.FileInfo, 0)
	f = append(f, &pb.FileInfo{FileName: args.FileName})
	res := &pb.String{}
	filePath := path.Join(define.PathIpfsData+"/c2task"+define.GateWayCsPathC2, args.FileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return res, nil
	}

	minerId := fmt.Sprintf("f0%d", args.Miner)
	filePath = "/ipfs/data/c2task/proof"
	err = deploy.WorkerClusterServiceApi.SynC2Result(minerId, filePath, args.FileName, data)
	if err != nil {
		log.Println("C2FileSynLotus file copy exception:", err)
		return res, err
	} else {
		// Mark C2 task file copy completed
		par := &pb.C2SectorID{ResType: 1, Miner: fmt.Sprintf("t0%d", args.Miner), Number: args.Number}
		_, err := deploy.WorkerClusterServiceApi.C2FileSynStatus(par)
		if err != nil {
			return res, err
		}
		os.Remove(path.Join(define.PathIpfsData+"/c2task"+define.GateWayCsPathC2, args.FileName))
	}

	return res, nil
}

// RunOpC2 Start OPC2
func (g *GateWayServiceImpl) RunOpC2(ctx context.Context, args *pb.String) (*pb.String, error) {

	client, dis := global.OpClinets.GetOpClient(args.Value)
	if dis {
		return nil, errors.New(fmt.Sprintf("opClient Connection failed:%s", args.Value))
	}
	return client.RunOpC2(ctx, args)
}

// StopOpC2 Stop OPC2
func (g *GateWayServiceImpl) StopOpC2(ctx context.Context, args *pb.String) (*pb.String, error) {

	client, dis := global.OpClinets.GetOpClient(args.Value)
	if dis {
		return nil, errors.New(fmt.Sprintf("opClient Connection failed:%s", args.Value))
	}
	return client.StopOpC2(ctx, args)
}

// RedoC2Task Redo C2 task
func (g *GateWayServiceImpl) RedoC2Task(ctx context.Context, args *pb.String) (*pb.String, error) {
	return &pb.String{}, lotusService.RedoC2Task(args.Value)
}
